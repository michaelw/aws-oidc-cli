# aws-creds-oidc

This project is an AWS SAM serverless Go application that provides AWS credentials via OIDC authentication. It exposes two endpoints via API Gateway:

- `/auth`: Constructs the OIDC authentication URL.
- `/creds`: Receives the code, verifies the state, exchanges the token for AWS credentials and returns them.

## Prerequisites

Ensure the following build dependencies are installed:

- golang
- aws-sam-cli

## Getting Started

1. Build the project:

   ```zsh
   make
   ```

2. Deploy to AWS:

   ```zsh
   sam deploy --guided
   ```

3. Create `~/.config/aws-oidc/oidc-providers.json`:

   ```json
   {
      "providers": [
         {
               "name": "test-provider",
               "api_url": "<API endpoint from deployment step>"
         }
      ]
   }
   ```

4. Add this to `~/.aws/config`:

   ```
   [profile oidc-test:administrator]
   credential_process = /path/to/aws-oidc process --provider=test-provider --role=oidc-administrator-access --account=1234567890
   ```

5. Test with AWS CLI:

   ```zsh
   aws sts get-caller-identity --profile oidc-test:administrator
   ```

## Development

### Local Testing

1. Create `env.json`:

   ```json
   {
      "AwsCredsFunction": {
         "OIDC_ISSUER": "https://.../",
         "OIDC_CLIENT_ID": "<...>",
         "OIDC_CLIENT_SECRET": "<...>"
      }
   }
   ```

   Ensure that `${OIDC_ISSUER}/.well-known/openid-configuration` exists and is accessible, and has a corresponding client credential configured.

2. Run the following commands from the source directory:

   ```zsh
   make
   sam local start-api --env-vars env.json
   ```

3. `oidc-providers.json`:

   ```json
   {
      "providers": [
         {
               "name": "test-provider",
               "api_url": "http://127.0.0.1:3000/"
         }
      ]
   }
   ```

4. In another terminal session, run the following command:

   ```shell
   $ ./aws-oidc process --config=oidc-providers.json --provider=test-provider --role=oidc-administrator-access --account=1234567890
   [...]
   {
      "Version": 1,
      "AccessKeyId": "ASIA...",
      "SecretAccessKey": "...",
      "SessionToken": "...",
      "Expiration": "2025-05-15T17:15:30Z"
   }
   ```

### Project Structure
- `./cmd/aws-creds-lambda/`: Main Go Lambda function (to be refactored for endpoints and business logic)
- `./cmd/aws-oidc/`: Command line utility to interface with AWS CLI, via `credential_process`
- `template.yaml`: AWS SAM template defining resources and API Gateway endpoints

## Pre-commit Hooks

This project uses [pre-commit](https://pre-commit.com/) to enforce code quality and security checks before each commit.

To set up pre-commit hooks:

```sh
./scripts/install-pre-commit.sh
```

This will install pre-commit (if not already installed) and set up the git hooks. Hooks include formatting, linting, security checks, and more. You can run all hooks manually with:

```sh
pre-commit run --all-files
```

## Alternatives

- Use [synfinatic/aws-sso-cli](https://github.com/synfinatic/aws-sso-cli) if you can use AWS SSO.  It does not require deploying additional AWS resources (lambda, OIDC providers, IAM roles with trust policies for the OIDC provider, etc.)

## License
MIT
