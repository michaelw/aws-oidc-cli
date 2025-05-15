package handler

import (
	"time"
)

// AwsOidcAuthRequest is the input for /auth.
type AwsOidcAuthRequest struct {
	Challenge   string `json:"challenge"`
	State       string `json:"state"`
	RedirectURI string `json:"redirect_uri"`
}

// AwsOidcCredsRequest is the input for /creds POST endpoint.
type AwsOidcCredsRequest struct {
	Code        string `json:"code"`
	Verifier    string `json:"verifier"`
	Account     string `json:"account"`
	Role        string `json:"role"`
	RedirectURI string `json:"redirect_uri"`
}

// AwsOidcCredsResponse is the output for /auth.
type AwsOidcCredsResponse struct {
	Version         int
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
	Expiration      time.Time
}
