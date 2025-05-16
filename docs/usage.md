# Usage

## Prerequisites

Ensure the following build dependencies are installed:

- golang
- aws-sam-cli

## Steps

1. **Build the project:**

   ```sh
   make
   ```

2. **Deploy to AWS:**

   ```sh
   sam deploy --guided
   ```

3. **Create `~/.config/aws-oidc/oidc-providers.json`:**

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

4. **Add a profile to `~/.aws/config`:**

   ```
   [profile oidc-test:administrator]
   credential_process = /path/to/aws-oidc process --provider=test-provider --role=oidc-administrator-access --account=1234567890
   ```

5. **Test with AWS CLI:**

   ```console
   $ aws sts get-caller-identity --profile oidc-test:administrator
   {
      "UserId": "AROAY6QNGSHIVDFKWHO3G:user@example.com",
      "Account": "1234567890",
      "Arn": "arn:aws:sts::1234567890:assumed-role/oidc-administrator-access/user@example.com"
   }
   ```
