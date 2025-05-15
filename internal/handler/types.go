package handler

import (
	"time"
)

// AwsCredsRequest is the input for /auth.
type AwsCredsRequest struct {
	Challenge   string `json:"challenge"`
	State       string `json:"state"`
	RedirectURI string `json:"redirect_uri"`
}

// AwsCredsResponse is the output for /auth.
type AwsCredsResponse struct {
	Version         int
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
	Expiration      time.Time
}
