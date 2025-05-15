// Package main provides the Lambda entrypoint for the aws-creds-oidc service.
package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	coreosoidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/michaelw/aws-creds-oidc/internal/awsutils"
	handler "github.com/michaelw/aws-creds-oidc/internal/handler"
	"github.com/michaelw/aws-creds-oidc/internal/oidc"
)

func main() {
	log.SetFlags(log.Lshortfile) // Disable timestamp and other prefixes
	ctx := context.Background()

	issuer := os.Getenv("OIDC_ISSUER")
	clientID := os.Getenv("OIDC_CLIENT_ID")
	clientSecret := os.Getenv("OIDC_CLIENT_SECRET")
	provider, err := coreosoidc.NewProvider(ctx, issuer)
	if err != nil {
		log.Fatalf("failed to initialize OIDC provider: %v", err)
	}
	oidcClient := oidc.NewOIDCClient(
		provider,
		clientID,
		clientSecret,
	)

	stsClient, err := awsutils.NewSTSClient(ctx)
	if err != nil {
		log.Fatalf("failed to initialize STS client: %v", err)
	}

	h := handler.NewAwsCredsHandler(oidcClient, stsClient)
	lambda.Start(h.Serve)
}
