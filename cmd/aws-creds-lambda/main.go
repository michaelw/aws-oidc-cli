// Package main provides the Lambda entrypoint for the aws-creds-oidc service.
package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	oidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/michaelw/aws-creds-oidc/internal/awsutils"
	"github.com/michaelw/aws-creds-oidc/internal/handler"
)

func main() {
	log.SetFlags(log.Lshortfile) // Disable timestamp and other prefixes
	ctx := context.Background()

	issuer := os.Getenv("OIDC_ISSUER")
	clientID := os.Getenv("OIDC_CLIENT_ID")
	clientSecret := os.Getenv("OIDC_CLIENT_SECRET")
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		log.Fatalf("failed to initialize OIDC provider: %v", err)
	}
	// OIDC client will be constructed in the handler with redirectURL from the request

	stsClient, err := awsutils.NewSTSClient(ctx)
	if err != nil {
		log.Fatalf("failed to initialize STS client: %v", err)
	}

	h := handler.NewAwsCredsHandler(provider, clientID, clientSecret, stsClient)
	lambda.Start(h.Serve)
}
