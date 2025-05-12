package oidc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockOIDCClient(t *testing.T) {
	oidcClient := &MockOIDCClient{}
	authURL, err := oidcClient.StartAuth(context.Background(), "cid", "sessid", "verifier", "state")
	assert.NoError(t, err)
	assert.Equal(t, "mockAuthURL", authURL)

	token, err := oidcClient.ExchangeCode(context.Background(), "code", "verifier")
	assert.NoError(t, err)
	assert.Equal(t, "mockAccessToken", token.AccessToken)
	assert.Equal(t, "mockIDToken", token.Extra("id_token"))
}
