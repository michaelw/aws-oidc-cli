package oidc

import (
	"context"

	"golang.org/x/oauth2"
)

// MockOIDCClient is a mock implementation of OIDCClient for testing.
type MockOIDCClient struct {
	StartAuthFunc    func(ctx context.Context, clientID, sessionID, verifier, state string) (string, error)
	ExchangeCodeFunc func(ctx context.Context, code, codeVerifier string) (*oauth2.Token, error)
}

func (m *MockOIDCClient) StartAuth(ctx context.Context, clientID, sessionID, verifier, state string) (string, error) {
	if m.StartAuthFunc != nil {
		return m.StartAuthFunc(ctx, clientID, sessionID, verifier, state)
	}
	return "mockAuthURL", nil
}

func (m *MockOIDCClient) ExchangeCode(ctx context.Context, code, codeVerifier string) (*oauth2.Token, error) {
	if m.ExchangeCodeFunc != nil {
		return m.ExchangeCodeFunc(ctx, code, codeVerifier)
	}
	tok := &oauth2.Token{AccessToken: "mockAccessToken"}
	tok = tok.WithExtra(map[string]any{"id_token": "mockIDToken"})
	return tok, nil
}
