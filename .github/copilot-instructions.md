<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

# Project: aws-creds-oidc

- This is an AWS SAM serverless Go project for OIDC-based AWS credential vending.
- Emphasize modularity, dependency injection, and testability in all code.
- All business logic must be unit tested.
- Use latest AWS SDK for Go (v2), and prefer existing libraries for OIDC and DynamoDB.
- All secrets and sensitive data must be encrypted at rest.
- Structure code for easy mocking and testability.
- Endpoints: `/aws-creds` (long-poll for OIDC, then AssumeRoleWithWebIdentity) and `/callback` (OIDC redirect, completes session).
- Session state is stored in DynamoDB.
- Ensure that all sensitive data is encrypted at rest.
