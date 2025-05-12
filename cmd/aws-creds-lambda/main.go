// Package main provides the Lambda entrypoint for the aws-creds-oidc service.
package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/michaelw/aws-creds-oidc/internal/awsutils"
	"github.com/michaelw/aws-creds-oidc/internal/handler"
	"github.com/michaelw/aws-creds-oidc/internal/oidc"
)

func main() {
	log.SetFlags(log.Lshortfile) // Disable timestamp and other prefixes
	ctx := context.Background()

	// OIDC client
	oidcClient, err := oidc.NewOIDCClient(ctx)
	if err != nil {
		panic("failed to initialize OIDC client: " + err.Error())
	}
	stsClient, err := awsutils.NewSTSClient(ctx)
	// stsClient, err := &awsutils.MockSTSClient{}, nil
	if err != nil {
		panic("failed to initialize STS client: " + err.Error())
	}

	h := handler.NewAwsCredsHandler(oidcClient, stsClient)
	lambda.Start(h.Serve)
}
