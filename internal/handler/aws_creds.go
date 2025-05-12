// handleAwsCreds handles the /aws-creds endpoint.
// It starts the OIDC flow, stores session, and long-polls for completion.
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/michaelw/aws-creds-oidc/internal/awscreds"
	"github.com/michaelw/aws-creds-oidc/internal/oidc"
	"github.com/michaelw/aws-creds-oidc/internal/session"
)

// AwsCredsRequest is the input for /aws-creds.
type AwsCredsRequest struct {
	ClientID  string `json:"client_id"`
	AccountID string `json:"account_id"`
	RoleName  string `json:"role_name"`
}

// AwsCredsResponse is the output for /aws-creds.
type AwsCredsResponse struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	SessionToken    string `json:"session_token"`
	Expiration      int64  `json:"expiration"`
}

// Handler dependencies for DI.
type AwsCredsHandler struct {
	SessionStore session.Store
	OIDCClient   oidc.OIDCClient
	STSClient    awscreds.STSClient
}

// HandleAwsCreds is the Lambda handler for /aws-creds as a method.
func (h *AwsCredsHandler) HandleAwsCreds(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var input AwsCredsRequest
	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "invalid input"}, nil
	}

	sessionID := fmt.Sprintf("sess-%d", time.Now().UnixNano())
	state := fmt.Sprintf("state-%d", time.Now().UnixNano())
	_, _, err := h.OIDCClient.StartAuth(ctx, input.ClientID, state)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "oidc start error"}, nil
	}

	sess := &session.Session{
		SessionID:   sessionID,
		ClientID:    input.ClientID,
		AccountID:   input.AccountID,
		RoleName:    input.RoleName,
		State:       state,
		Status:      "pending",
		AccessToken: "",
		ExpiresAt:   time.Now().Add(10 * time.Minute).Unix(),
	}
	if err := h.SessionStore.Put(ctx, sess); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "session store error"}, nil
	}

	// Respond with auth URL and poll for completion
	pollTimeout := 60 * time.Second
	pollInterval := 2 * time.Second
	start := time.Now()
	for {
		if time.Since(start) > pollTimeout {
			return events.APIGatewayProxyResponse{StatusCode: 408, Body: "timeout"}, nil
		}
		s, err := h.SessionStore.Get(ctx, sessionID)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "session get error"}, nil
		}
		if s != nil && s.Status == "complete" && s.AccessToken != "" {
			// Call STS
			roleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", input.AccountID, input.RoleName)
			ak, sk, st, err := h.STSClient.AssumeRoleWithWebIdentity(ctx, roleArn, sessionID, s.AccessToken, 900)
			if err != nil {
				return events.APIGatewayProxyResponse{StatusCode: 500, Body: "sts error"}, nil
			}
			resp := AwsCredsResponse{
				AccessKeyID:     ak,
				SecretAccessKey: sk,
				SessionToken:    st,
				Expiration:      time.Now().Add(15 * time.Minute).Unix(),
			}
			b, _ := json.Marshal(resp)
			return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(b)}, nil
		}
		time.Sleep(pollInterval)
	}
}

// HandleCallback handles the /callback endpoint for OIDC redirect as a method of AwsCredsHandler.
func (h *AwsCredsHandler) HandleCallback(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	code := req.QueryStringParameters["code"]
	state := req.QueryStringParameters["state"]
	if code == "" || state == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "missing code or state"}, nil
	}
	// Find session by state
	sess, err := h.findSessionByState(ctx, state)
	if err != nil || sess == nil {
		return events.APIGatewayProxyResponse{StatusCode: 404, Body: "session not found"}, nil
	}
	// OIDC client
	var oidcClient oidc.OIDCClient
	if h.OIDCClient != nil {
		oidcClient = h.OIDCClient
	} else {
		oidcClient, err = oidc.NewOIDCClient(ctx)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "oidc error"}, nil
		}
	}
	// For demo, code_verifier is not persisted; in production, store it in session
	accessToken, err := oidcClient.ExchangeCode(ctx, code, "")
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "oidc exchange error"}, nil
	}
	sess.AccessToken = accessToken
	sess.Status = "complete"
	err = h.SessionStore.Update(ctx, sess)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "session update error"}, nil
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "session complete"}, nil
}

// findSessionByState finds a session by state (inefficient scan; optimize for production).
func (h *AwsCredsHandler) findSessionByState(ctx context.Context, state string) (*session.Session, error) {
	ds, ok := h.SessionStore.(*session.DynamoStore)
	if !ok {
		return nil, nil
	}
	table := ds.TableName
	if table == "" {
		return nil, nil
	}
	client := ds.Client
	out, err := client.Scan(ctx, &dynamodb.ScanInput{
		TableName: &table,
	})
	if err != nil {
		return nil, err
	}
	for _, item := range out.Items {
		var s session.Session
		err := attributevalue.UnmarshalMap(item, &s)
		if err == nil && strings.EqualFold(s.State, state) {
			return &s, nil
		}
	}
	return nil, nil
}

// Serve routes API Gateway requests to the appropriate handler method.
func (h *AwsCredsHandler) Serve(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.Path {
	case "/aws-creds":
		return h.HandleAwsCreds(ctx, req)
	case "/callback":
		return h.HandleCallback(ctx, req)
	default:
		return events.APIGatewayProxyResponse{StatusCode: 404}, nil
	}
}

// NewAwsCredsHandler constructs a handler with injected dependencies.
func NewAwsCredsHandler(store session.Store, oidcClient oidc.OIDCClient, stsClient awscreds.STSClient) *AwsCredsHandler {
	return &AwsCredsHandler{
		SessionStore: store,
		OIDCClient:   oidcClient,
		STSClient:    stsClient,
	}
}
