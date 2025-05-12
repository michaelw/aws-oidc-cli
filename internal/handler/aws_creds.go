// handleAwsCreds handles the /auth endpoint.
// It starts the OIDC flow, stores session, and long-polls for completion.
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
	awscreds "github.com/michaelw/aws-creds-oidc/internal/awsutils"
	"github.com/michaelw/aws-creds-oidc/internal/oidc"
)

// AwsCredsRequest is the input for /auth.
type AwsCredsRequest struct {
	State     string `json:"state"`
	Challenge string `json:"challenge"`
}

// AwsCredsResponse is the output for /auth.
type AwsCredsResponse struct {
	Version         int
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
	Expiration      time.Time
}

// Handler dependencies for DI.
type AwsCredsHandler struct {
	OIDCClient oidc.OIDCClient
	STSClient  awscreds.STSClient
}

// NewAwsCredsHandler constructs a handler with injected dependencies.
func NewAwsCredsHandler(oidcClient oidc.OIDCClient, stsClient awscreds.STSClient) *AwsCredsHandler {
	return &AwsCredsHandler{
		OIDCClient: oidcClient,
		STSClient:  stsClient,
	}
}

// HandleAuth is the Lambda handler for /auth as a method.
func (h *AwsCredsHandler) HandleAuth(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	state := req.QueryStringParameters["state"]
	challenge := req.QueryStringParameters["challenge"]
	if state == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "missing state"}, nil
	}
	if challenge == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "missing challenge"}, nil
	}

	authURL, err := h.OIDCClient.StartAuth(ctx, challenge, state)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("OIDC start error: %v", err)}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 302,
		Headers: map[string]string{
			"Location": authURL,
		},
	}, nil
}

// HandleCreds handles the /creds endpoint for OIDC redirect as a method of AwsCredsHandler.
// Now expects POST with JSON body: { code, verifier, account, role }
func (h *AwsCredsHandler) HandleCreds(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body struct {
		Code     string `json:"code"`
		Verifier string `json:"verifier"`
		Account  string `json:"account"`
		Role     string `json:"role"`
	}
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "invalid JSON body"}, nil
	}
	if body.Code == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "missing code"}, nil
	}
	if body.Verifier == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "missing verifier"}, nil
	}
	if body.Account == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "missing account ID"}, nil
	}
	if body.Role == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "missing role"}, nil
	}

	// Exchange code for token
	token, err := h.OIDCClient.ExchangeCode(ctx, body.Code, body.Verifier)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: err.Error()}, nil
	}
	idToken, ok := token.Extra("id_token").(string)
	if !ok || idToken == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "no id_token in token response"}, nil
	}
	log.Printf("ID Token: %s", idToken)

	// Parse email from idToken
	claims, err := parseIDTokenClaims(idToken)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: fmt.Sprintf("failed to parse id_token: %v", err)}, nil
	}
	if claims.Email == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "email claim not found in id_token"}, nil
	}
	email := claims.Email
	log.Printf("Email: %s", email)

	// Call STS
	roleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", body.Account, body.Role)
	duration := 30 * time.Minute
	exp := time.Now().Add(duration)
	ak, sk, st, err := h.STSClient.AssumeRoleWithWebIdentity(ctx, roleArn, email, idToken, int32(duration.Seconds()))
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: err.Error()}, nil
	}

	// Return credentials in AWS credential_process format
	resp := AwsCredsResponse{
		Version:         1,
		AccessKeyId:     ak,
		SecretAccessKey: sk,
		SessionToken:    st,
		Expiration:      exp,
	}
	b, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{StatusCode: 200,
		Body:    string(b),
		Headers: map[string]string{"Content-Type": "application/json"},
	}, nil
}

// IDTokenClaims holds the claims we care about from the ID token
// (expand as needed for more claims)
type IDTokenClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// parseIDTokenClaims parses a JWT and returns the claims (without verifying signature)
func parseIDTokenClaims(idToken string) (*IDTokenClaims, error) {
	claims := &IDTokenClaims{}
	_, _, err := new(jwt.Parser).ParseUnverified(idToken, claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// Serve routes API Gateway requests to the appropriate handler method.
func (h *AwsCredsHandler) Serve(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.Path {
	case "/auth":
		return h.HandleAuth(ctx, req)
	case "/creds":
		return h.HandleCreds(ctx, req)
	default:
		return events.APIGatewayProxyResponse{StatusCode: 404}, nil
	}
}
