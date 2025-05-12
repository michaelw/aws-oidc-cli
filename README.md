# aws-creds-oidc

This project is an AWS SAM serverless Go application that provides AWS credentials via OIDC authentication. It exposes two endpoints via API Gateway:

- `/auth`: Constructs the OIDC authentication URL.
- `/creds`: Receives the code, verifies the state, exchanges the token for AWS credentials and returns them.

## Getting Started

1. Install the AWS SAM CLI and Go 1.x.
2. Build the project:

   ```zsh
   make
   ```

3. Deploy to AWS:

   ```zsh
   sam deploy --guided
   ```

## Project Structure
- `aws-creds-lambda/`: Main Go Lambda function (to be refactored for endpoints and business logic)
- `template.yaml`: AWS SAM template defining resources and API Gateway endpoints

## Development
- Follow Go best practices for modularity and testability.
- Add unit tests for all business logic.

## License
MIT
