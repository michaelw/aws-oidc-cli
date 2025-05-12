package awscreds

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// STSClient defines the interface for AWS STS operations.
type STSClient interface {
	AssumeRoleWithWebIdentity(ctx context.Context, roleArn, roleSessionName, webIdentityToken string, durationSeconds int32) (accessKeyID, secretAccessKey, sessionToken string, err error)
}

// stsClient implements STSClient using AWS SDK v2.
type stsClient struct {
	Client *sts.Client
}

func NewSTSClient(ctx context.Context) (STSClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &stsClient{Client: sts.NewFromConfig(cfg)}, nil
}

func (r *stsClient) AssumeRoleWithWebIdentity(ctx context.Context, roleArn, roleSessionName, webIdentityToken string, durationSeconds int32) (string, string, string, error) {
	out, err := r.Client.AssumeRoleWithWebIdentity(ctx, &sts.AssumeRoleWithWebIdentityInput{
		RoleArn:          aws.String(roleArn),
		RoleSessionName:  aws.String(roleSessionName),
		WebIdentityToken: aws.String(webIdentityToken),
		DurationSeconds:  aws.Int32(durationSeconds),
	})
	if err != nil {
		return "", "", "", err
	}
	return aws.ToString(out.Credentials.AccessKeyId), aws.ToString(out.Credentials.SecretAccessKey), aws.ToString(out.Credentials.SessionToken), nil
}
