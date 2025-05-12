package oidc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockOIDCClient(t *testing.T) {
	mock := &MockOIDCClient{}
	authURL, verifier, err := mock.StartAuth(context.Background(), "cid", "state")
	assert.NoError(t, err)
	assert.Equal(t, "mockAuthURL", authURL)
	assert.Equal(t, "mockVerifier", verifier)

	token, err := mock.ExchangeCode(context.Background(), "code", "verifier")
	assert.NoError(t, err)
	assert.Equal(t, "mockAccessToken", token)
}
