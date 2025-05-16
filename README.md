# aws-creds-oidc

This project is an AWS SAM application that provides short-lived AWS (STS) credentials via OIDC authentication.  A matching OIDC provider and roles with matching trust policies are required in the AWS account.

## Goals

* Obtain AWS STS credentials via OIDC authentication
* No AWS credentials stored on disk (as AWSCLI does)
* No other credentials stored unencrypted on clients
* For audit purposes, the token's `UserId` is not forgeable, unless OIDC client secrets are exposed
* No wrappers for AWSCLI

## Operation

The supplied credential provider CLI tool can be hooked into AWSCLI via the `credential_process` option.  It connects to the SAM lambda for authentication and credential retrieval.  While this could be done entirely locally, e.g. [aws-cli-oidc](https://github.com/stensonb/aws-cli-oidc), it would require distributing client credentials that are stored on disk unencrypted, or a public client.  Also, role ARNs would have to be communicated and managed for each account.

It exposes two endpoints via API Gateway:

- `/auth`: Constructs the OIDC authentication URL.
- `/creds`: Receives the code, verifies the state, exchanges the token for AWS credentials and returns them.

See [docs/architecture.md](docs/architecture.md) for architecture diagrams (rendered with Mermaid).

## Prerequisites

Ensure the following build dependencies are installed:

- golang
- aws-sam-cli

## Getting Started

See [docs/usage.md](docs/usage.md) for setup and usage instructions.

## Development

See [docs/dev.md](docs/dev.md) for development instructions.

## Alternatives

- Use [synfinatic/aws-sso-cli](https://github.com/synfinatic/aws-sso-cli) if you can use AWS SSO.  It does not require deploying additional AWS resources (lambda, OIDC providers, IAM roles with trust policies for the OIDC provider, etc.).  The tool is generally very feature-complete, and can even generate and manage all accessible SSO profiles in the AWS config file.

- Use [99designs/aws-vault](https://github.com/99designs/aws-vault) if you do not mind wrapping your AWS commands into `aws-vault exec <profile> -- ...`, or duplicating your AWS profiles when using `aws-vault` as `credential_process`.

- Use [chanzuckerberg/aws-oidc](https://github.com/chanzuckerberg/aws-oidc), if you do not mind that the role session name is forgeable.

- Use [stensonb/aws-cli-oidc](https://github.com/stensonb/aws-cli-oidc) if you do not mind that client credentials are exposed on disk, and that the role session name is forgeable.

## License
MIT
