package handler

import (
	"time"
)

// AuthRequest is the input for /auth.
type AuthRequest struct {
	Challenge   string `json:"challenge"`
	State       string `json:"state"`
	RedirectURI string `json:"redirect_uri"`
}

// CredsRequest is the input for /creds POST endpoint.
type CredsRequest struct {
	Code        string `json:"code"`
	Verifier    string `json:"verifier"`
	Account     string `json:"account"`
	Role        string `json:"role"`
	RedirectURI string `json:"redirect_uri"`
}

// CredsResponse is the output for /auth.
type CredsResponse struct {
	Version         int
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
	Expiration      time.Time
}
