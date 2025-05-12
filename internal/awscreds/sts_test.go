package awscreds

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockSTSClient(t *testing.T) {
	mock := &MockSTSClient{}
	ak, sk, st, err := mock.AssumeRoleWithWebIdentity(context.Background(), "arn", "sess", "token", 900)
	assert.NoError(t, err)
	assert.Equal(t, "mockAccessKey", ak)
	assert.Equal(t, "mockSecretKey", sk)
	assert.Equal(t, "mockSessionToken", st)
}
