package oidc

import (
	"context"

	"golang.org/x/oauth2"
)

type MockOIDCClient struct {
	OIDCClient
	ExchangeCodeFunc func(ctx context.Context, code, verifier, redirectURI string) (*oauth2.Token, error)
}

var _ OIDCClient = (*MockOIDCClient)(nil)

func (m *MockOIDCClient) ExchangeCode(ctx context.Context, code, verifier, redirectURI string) (*oauth2.Token, error) {
	if m.ExchangeCodeFunc != nil {
		return m.ExchangeCodeFunc(ctx, code, verifier, redirectURI)
	}
	tok := &oauth2.Token{AccessToken: "mockAccessToken"}
	tok = tok.WithExtra(map[string]any{"id_token": "mockIDToken"})
	return tok, nil
}
