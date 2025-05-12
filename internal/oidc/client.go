package oidc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// OIDCClient defines the interface for OIDC operations.
type OIDCClient interface {
	StartAuth(ctx context.Context, clientID, state string) (authURL, codeVerifier string, err error)
	ExchangeCode(ctx context.Context, code, codeVerifier string) (accessToken string, err error)
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

func generatePKCE() (string, string, error) {
	codeVerifier := make([]byte, 32)
	_, err := rand.Read(codeVerifier)
	if err != nil {
		return "", "", err
	}
	verifier := base64.RawURLEncoding.EncodeToString(codeVerifier)
	challenge := verifier // For simplicity, use plain (not S256) for now
	return verifier, challenge, nil
}

func (c *oidcClient) StartAuth(ctx context.Context, clientID, state string) (string, string, error) {
	verifier, challenge, err := generatePKCE()
	if err != nil {
		return "", "", err
	}
	url := c.OAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("code_challenge", challenge), oauth2.SetAuthURLParam("code_challenge_method", "plain"))
	return url, verifier, nil
}

func (c *oidcClient) ExchangeCode(ctx context.Context, code, codeVerifier string) (string, error) {
	tok, err := c.OAuthConfig.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		return "", err
	}
	return tok.AccessToken, nil
}
