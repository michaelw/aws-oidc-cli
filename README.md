# aws-creds-oidc

This project is an AWS SAM serverless Go application that provides AWS credentials via OIDC authentication. It exposes two endpoints via API Gateway:

- `/aws-creds`: Long-polls for OIDC authentication (client credential flow with PKCE), then uses the access token with AWS AssumeRoleWithWebIdentity to return AWS credentials for the requested account and role.
- `/callback`: Handles the OIDC redirect, receives the code, verifies the state, and completes the session for the long-polling request.

## Features
- Written in Go, following best practices for modularity, dependency injection, and testability.
- All business logic is unit tested.
- Uses AWS SDK for Go v2 and existing Go libraries for OIDC and DynamoDB.
- All sensitive and secret data is encrypted at rest.
- Session state is stored in DynamoDB for durability and signaling between endpoints.

## Getting Started

1. Install the AWS SAM CLI and Go 1.x.
2. Build the project:
   ```zsh
   sam build
   ```
3. Deploy to AWS:
   ```zsh
   sam deploy --guided
   ```

## Project Structure
- `hello-world/`: Main Go Lambda function (to be refactored for endpoints and business logic)
- `template.yaml`: AWS SAM template defining resources and API Gateway endpoints

## Development
- Follow Go best practices for modularity and testability.
- Use dependency injection for all components.
- Add unit tests for all business logic.

## License
MIT
