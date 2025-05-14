package awsutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockSTSClient(t *testing.T) {
	var stsClient STSClient = &MockSTSClient{}
	ak, sk, st, _, err := stsClient.AssumeRoleWithWebIdentity(context.Background(), "arn", "sess", "token", 900)
	assert.NoError(t, err)
	assert.Equal(t, "mockAccessKey", ak)
	assert.Equal(t, "mockSecretKey", sk)
	assert.Equal(t, "mockSessionToken", st)
}
