package handler

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	coreosoidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	awsutils "github.com/michaelw/aws-creds-oidc/internal/awsutils"
	"github.com/michaelw/aws-creds-oidc/internal/oidc"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func newTestHandler(stsErr error, token *oauth2.Token, tokenErr error) *AwsCredsHandler {
	return NewAwsCredsHandler(
		&oidc.MockOIDCClient{
			OIDCClient: oidc.NewOIDCClient(
				&coreosoidc.Provider{},
				"clientid",
				"secret",
			),
			ExchangeCodeFunc: func(ctx context.Context, code, verifier, redirectURI string) (*oauth2.Token, error) {
				return token, tokenErr
			},
		},
		&awsutils.MockSTSClient{
			AssumeRoleWithWebIdentityFunc: func(ctx context.Context, roleArn, roleSessionName, webIdentityToken string, durationSeconds int32) (string, string, string, *time.Time, error) {
				if stsErr != nil {
					return "", "", "", nil, stsErr
				}
				exp := time.Now().Add(1 * time.Hour)
				return "AKIA", "SK", "ST", &exp, nil
			},
		},
	)
}

func TestHandleAuth_MissingParams(t *testing.T) {
	h := newTestHandler(nil, nil, nil)
	cases := []struct {
		name   string
		params map[string]string
		errMsg string
	}{
		{"missing state", map[string]string{"challenge": "c", "redirect_uri": "r"}, "missing state"},
		{"missing challenge", map[string]string{"state": "s", "redirect_uri": "r"}, "missing challenge"},
		{"missing redirect_uri", map[string]string{"state": "s", "challenge": "c"}, "missing redirect_uri"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := events.APIGatewayProxyRequest{QueryStringParameters: c.params}
			resp, _ := h.HandleAuth(context.Background(), req)
			assert.Equal(t, 400, resp.StatusCode)
			assert.Contains(t, resp.Body, c.errMsg)
		})
	}
}

func TestHandleCreds_MissingFields(t *testing.T) {
	h := newTestHandler(nil, nil, nil)
	base := AwsOidcCredsRequest{
		Code:        "c",
		Verifier:    "v",
		Account:     "a",
		Role:        "r",
		RedirectURI: "u",
	}
	fields := []struct {
		name   string
		modify func(*AwsOidcCredsRequest)
		errMsg string
	}{
		{"missing code", func(b *AwsOidcCredsRequest) { b.Code = "" }, "missing code"},
		{"missing verifier", func(b *AwsOidcCredsRequest) { b.Verifier = "" }, "missing verifier"},
		{"missing account", func(b *AwsOidcCredsRequest) { b.Account = "" }, "missing account ID"},
		{"missing role", func(b *AwsOidcCredsRequest) { b.Role = "" }, "missing role"},
		{"missing redirect_uri", func(b *AwsOidcCredsRequest) { b.RedirectURI = "" }, "missing redirect_uri"},
	}
	for _, f := range fields {
		t.Run(f.name, func(t *testing.T) {
			b := base
			f.modify(&b)
			data, _ := json.Marshal(b)
			req := events.APIGatewayProxyRequest{Body: string(data)}
			resp, _ := h.HandleCreds(context.Background(), req)
			assert.Equal(t, 400, resp.StatusCode)
			assert.Contains(t, resp.Body, f.errMsg)
		})
	}
}

func TestHandleCreds_InvalidJSON(t *testing.T) {
	h := newTestHandler(nil, nil, nil)
	req := events.APIGatewayProxyRequest{Body: "notjson"}
	resp, _ := h.HandleCreds(context.Background(), req)
	assert.Equal(t, 400, resp.StatusCode)
	assert.Contains(t, resp.Body, "invalid JSON body")
}

func TestHandleCreds_STSError(t *testing.T) {
	tok := &oauth2.Token{}
	tok = tok.WithExtra(map[string]any{"id_token": createTestJWT(t, "foo@bar.com")})
	h := newTestHandler(errors.New("sts error"), tok, nil)
	b := AwsOidcCredsRequest{
		Code:        "c",
		Verifier:    "v",
		Account:     "a",
		Role:        "r",
		RedirectURI: "u",
	}
	data, _ := json.Marshal(b)
	req := events.APIGatewayProxyRequest{Body: string(data)}
	resp, _ := h.HandleCreds(context.Background(), req)
	assert.Equal(t, 400, resp.StatusCode)
	assert.Contains(t, resp.Body, "sts error")
}

func TestHandleCreds_ValidFlow(t *testing.T) {
	tok := &oauth2.Token{}
	tok = tok.WithExtra(map[string]any{"id_token": createTestJWT(t, "foo@bar.com")})
	h := newTestHandler(nil, tok, nil)
	b := AwsOidcCredsRequest{
		Code:        "c",
		Verifier:    "v",
		Account:     "a",
		Role:        "r",
		RedirectURI: "u",
	}
	data, _ := json.Marshal(b)
	req := events.APIGatewayProxyRequest{Body: string(data)}
	resp, _ := h.HandleCreds(context.Background(), req)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Contains(t, resp.Body, "AccessKeyId")
}

func createTestJWT(t *testing.T, email string) string {
	claims := jwt.MapClaims{"email": email}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("failed to sign test jwt: %v", err)
	}
	return s
}

func TestServe_UnknownPath(t *testing.T) {
	h := newTestHandler(nil, nil, nil)
	req := events.APIGatewayProxyRequest{Path: "/unknown"}
	resp, _ := h.Serve(context.Background(), req)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestParseIDTokenClaimsUnverified(t *testing.T) {
	// Valid JWT with email claim
	claims := jwt.MapClaims{"email": "foo@bar.com"}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := tok.SignedString([]byte("secret"))
	parsed, err := parseIDTokenClaimsUnverified(s)
	assert.NoError(t, err)
	assert.Equal(t, "foo@bar.com", parsed.Email)

	// Invalid JWT
	_, err = parseIDTokenClaimsUnverified("notatoken")
	assert.Error(t, err)
}
