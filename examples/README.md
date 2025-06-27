# OpenRouter Terraform Provider Examples

This directory contains example Terraform configurations showing how to use the OpenRouter provider.

## Prerequisites

1. Install the OpenRouter provider
2. Set your API key as an environment variable:
   ```bash
   export OPENROUTER_API_KEY="sk-or-v1-your-api-key-here"
   ```

## Examples

### Basic Provider Configuration
See `provider.tf` for basic provider setup.

### API Key Management
See `api-key-management/` for examples of:
- Creating API keys with limits
- Reading existing API keys
- Listing all API keys
- Updating API key properties

### Complete Setup
See `complete/` for a full example with all resources and data sources.

## Running Examples

1. Navigate to an example directory
2. Initialize Terraform: `terraform init`
3. Plan changes: `terraform plan`
4. Apply changes: `terraform apply`

## Cleanup

To remove all created resources:
```bash
terraform destroy
```