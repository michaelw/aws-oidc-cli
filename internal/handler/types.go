package handler

import (
	"time"
)

// AwsCredsRequest is the input for /auth.
type AwsCredsRequest struct {
	State     string `json:"state"`
	Challenge string `json:"challenge"`
}

// AwsCredsResponse is the output for /auth.
type AwsCredsResponse struct {
	Version         int
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
	Expiration      time.Time
}
