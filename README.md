# aws-creds-oidc

This project is an AWS SAM application that provides short-lived AWS (STS) credentials via OIDC authentication.  A matching OIDC provider and roles with matching trust policies are required in the AWS account.


## Goals

* OIDC authentication
* no AWS credentials stored on disk (as AWSCLI does)
* no other credentials stored unencrypted on clients
* token `UserId` is not forgeable, unless OIDC client secrets are exposed
* no wrappers for AWSCLI

## Operation

The supplied credential provider CLI tool can be hooked into AWSCLI via the `credential_process` option.  It connects to the SAM lambda for authentication and credential retrieval.  While this could be done entirely locally, e.g. [aws-cli-oidc](https://github.com/stensonb/aws-cli-oidc), it requires distributing client credentials that are stored on disk unencrypted.  Also, role ARNs have to be communicated and managed for each account.

It exposes two endpoints via API Gateway:

- `/auth`: Constructs the OIDC authentication URL.
- `/creds`: Receives the code, verifies the state, exchanges the token for AWS credentials and returns them.

## Prerequisites

Ensure the following build dependencies are installed:

- golang
- aws-sam-cli

## Getting Started

1. Build the project:

   ```sh
   make
   ```

2. Deploy to AWS:

   ```sh
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

   ```console
   $ aws sts get-caller-identity --profile oidc-test:administrator
   {
      "UserId": "AROAY6QNGSHIVDFKWHO3G:user@example.com",
      "Account": "1234567890",
      "Arn": "arn:aws:sts::1234567890:assumed-role/oidc-administrator-access/user@example.com"
   }
   ```

## Development

### Local Testing

1. Create `env.json`:

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

2. Run the following commands from the source directory:

   ```sh
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

   ```console
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

### Pre-commit Hooks

This project uses [pre-commit](https://pre-commit.com/) to enforce code quality and security checks before each commit.

To set up pre-commit hooks:

```sh
pre-commit install
```

This will install pre-commit (if not already installed) and set up the git hooks. Hooks include formatting, linting, security checks, and more. You can run all hooks manually with:

```sh
pre-commit run --all-files
```

## Alternatives

- Use [synfinatic/aws-sso-cli](https://github.com/synfinatic/aws-sso-cli) if you can use AWS SSO.  It does not require deploying additional AWS resources (lambda, OIDC providers, IAM roles with trust policies for the OIDC provider, etc.).  The tool is generally very feature-complete, and can even generate and manage all accessible SSO profiles in the AWS config file.

- Use [chanzuckerberg/aws-oidc](https://github.com/chanzuckerberg/aws-oidc), if you do not mind that the role session name is forgeable.

- Use [stensonb/aws-cli-oidc](https://github.com/stensonb/aws-cli-oidc) if you do not mind that client credentials are exposed on disk, and that the role session name is forgeable.

## License
MIT
