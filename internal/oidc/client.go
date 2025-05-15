package oidc

import (
	"context"

	coreosoidc "github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// OIDCClient allows for mocking oidcClient in tests
// Used by handler for dependency injection
// Renamed from OIDCClientIface
type OIDCClient interface {
	NewConfig(redirectURI string) *oauth2.Config
	ExchangeCode(ctx context.Context, code, verifier, redirectURI string) (*oauth2.Token, error)
}

// oidcClient holds OIDC provider and client credentials
// Implements OIDCClient (interface)
type oidcClient struct {
	Provider     *coreosoidc.Provider
	ClientID     string
	ClientSecret string
}

// NewOIDCClient constructs a new oidcClient and returns it as OIDCClient
func NewOIDCClient(provider *coreosoidc.Provider, clientID, clientSecret string) OIDCClient {
	return &oidcClient{
		Provider:     provider,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

func (c *oidcClient) NewConfig(redirectURI string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Endpoint:     c.Provider.Endpoint(),
		RedirectURL:  redirectURI,
		Scopes:       []string{coreosoidc.ScopeOpenID, "profile", "email"},
	}
}

func (c *oidcClient) ExchangeCode(ctx context.Context, code, verifier, redirectURI string) (*oauth2.Token, error) {
	return c.NewConfig(redirectURI).Exchange(ctx, code, oauth2.VerifierOption(verifier))
}
