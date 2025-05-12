// Package main provides the Lambda entrypoint for the aws-creds-oidc service.
package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/michaelw/aws-creds-oidc/internal/awscreds"
	"github.com/michaelw/aws-creds-oidc/internal/handler"
	"github.com/michaelw/aws-creds-oidc/internal/oidc"
	"github.com/michaelw/aws-creds-oidc/internal/session"
)

func main() {
	ctx := context.Background()
	store, err := session.NewDynamoStore(ctx)
	if err != nil {
		panic("failed to initialize session store: " + err.Error())
	}
	// OIDC client
	oidcClient, err := oidc.NewOIDCClient(ctx)
	if err != nil {
		panic("failed to initialize OIDC client: " + err.Error())
	}
	stsClient, err := awscreds.NewSTSClient(ctx)
	if err != nil {
		panic("failed to initialize STS client: " + err.Error())
	}

	h := handler.NewAwsCredsHandler(store, oidcClient, stsClient)
	lambda.Start(h.Serve)
}
