package oidc

import (
	"context"
)

// MockOIDCClient is a mock implementation of OIDCClient for testing.
type MockOIDCClient struct {
	StartAuthFunc    func(ctx context.Context, clientID, state string) (string, string, error)
	ExchangeCodeFunc func(ctx context.Context, code, codeVerifier string) (string, error)
}

func (m *MockOIDCClient) StartAuth(ctx context.Context, clientID, state string) (string, string, error) {
	if m.StartAuthFunc != nil {
		return m.StartAuthFunc(ctx, clientID, state)
	}
	return "mockAuthURL", "mockVerifier", nil
}

func (m *MockOIDCClient) ExchangeCode(ctx context.Context, code, codeVerifier string) (string, error) {
	if m.ExchangeCodeFunc != nil {
		return m.ExchangeCodeFunc(ctx, code, codeVerifier)
	}
	return "mockAccessToken", nil
}
