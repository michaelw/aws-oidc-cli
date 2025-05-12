package awscreds

import (
	"context"
)

// MockSTSClient is a mock implementation of STSClient for testing.
type MockSTSClient struct {
	AssumeRoleWithWebIdentityFunc func(ctx context.Context, roleArn, roleSessionName, webIdentityToken string, durationSeconds int32) (string, string, string, error)
}

func (m *MockSTSClient) AssumeRoleWithWebIdentity(ctx context.Context, roleArn, roleSessionName, webIdentityToken string, durationSeconds int32) (string, string, string, error) {
	if m.AssumeRoleWithWebIdentityFunc != nil {
		return m.AssumeRoleWithWebIdentityFunc(ctx, roleArn, roleSessionName, webIdentityToken, durationSeconds)
	}
	return "mockAccessKey", "mockSecretKey", "mockSessionToken", nil
}
