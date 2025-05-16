# Development Guide


## Project Structure

- `./cmd/aws-creds-lambda/`: Main Go Lambda function (to be refactored for endpoints and business logic)
- `./cmd/aws-oidc/`: Command line utility to interface with AWS CLI, via `credential_process`
- `template.yaml`: AWS SAM template defining resources and API Gateway endpoints

## Local Testing

1. **Create `env.json`:**

   ```json
   {
      "AwsCredsFunction": {
         "OIDC_ISSUER": "https://...",
         "OIDC_CLIENT_ID": "<...>",
         "OIDC_CLIENT_SECRET": "<...>"
      }
   }
   ```

   Ensure that `${OIDC_ISSUER}/.well-known/openid-configuration` exists and is accessible, and has a corresponding client credential configured.

2. **Start the local API:**

   ```sh
   make
   sam local start-api --env-vars env.json
   ```

3. **Create `oidc-providers.json`:**

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

4. **Run the CLI tool in another terminal:**

   ```sh
   ./aws-oidc process --config=oidc-providers.json --provider=test-provider --role=oidc-administrator-access --account=1234567890
   ```

   Example output:
   ```json
   {
      "Version": 1,
      "AccessKeyId": "ASIA...",
      "SecretAccessKey": "...",
      "SessionToken": "...",
      "Expiration": "2025-05-15T17:15:30Z"
   }
   ```

## Pre-commit Hooks

This project uses [pre-commit](https://pre-commit.com/) to enforce code quality and security checks before each commit.

To set up pre-commit hooks:

```sh
pre-commit install
```

This will install pre-commit (if not already installed) and set up the git hooks. Hooks include formatting, linting, security checks, and more. You can run all hooks manually with:

```sh
pre-commit run --all-files
```
