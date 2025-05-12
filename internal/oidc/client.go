package oidc

import (
	"context"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// OIDCClient defines the interface for OIDC operations.
type OIDCClient interface {
	StartAuth(ctx context.Context, codeChallengeS256, state string) (authURL string, err error)
	ExchangeCode(ctx context.Context, code, codeVerifier string) (*oauth2.Token, error)
}

// oidcClient implements OIDCClient using go-oidc and oauth2.
type oidcClient struct {
	Provider    *oidc.Provider
	OAuthConfig *oauth2.Config
}

func NewOIDCClient(ctx context.Context) (OIDCClient, error) {
	issuer := os.Getenv("OIDC_ISSUER")
	clientID := os.Getenv("OIDC_CLIENT_ID")
	clientSecret := os.Getenv("OIDC_CLIENT_SECRET")
	redirectURL := os.Getenv("OIDC_REDIRECT_URL")
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, err
	}
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  redirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	return &oidcClient{Provider: provider, OAuthConfig: config}, nil
}

func (c *oidcClient) StartAuth(ctx context.Context, challengeS256, state string) (string, error) {
	url := c.OAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", challengeS256),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"))
	return url, nil
}

func (c *oidcClient) ExchangeCode(ctx context.Context, code, codeVerifier string) (*oauth2.Token, error) {
	tok, err := c.OAuthConfig.Exchange(ctx, code, oauth2.VerifierOption(codeVerifier))
	if err != nil {
		return nil, err
	}
	return tok, nil
}
