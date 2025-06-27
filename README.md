# Terraform Provider for OpenRouter

[![Release](https://img.shields.io/github/v/release/standujar/terraform-provider-openrouter)](https://github.com/standujar/terraform-provider-openrouter/releases)
[![License](https://img.shields.io/github/license/standujar/terraform-provider-openrouter)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/standujar/terraform-provider-openrouter)](https://goreportcard.com/report/github.com/standujar/terraform-provider-openrouter)

A Terraform provider for managing [OpenRouter](https://openrouter.ai/) resources, enabling infrastructure-as-code management of your OpenRouter API keys and configurations.

## Features

- 🔑 **API Key Management**: Create, read, update, and delete OpenRouter API keys
- 💰 **Spending Controls**: Set spending limits and time-based restrictions
- 📊 **Usage Monitoring**: Track API key usage and status
- 🔒 **Secure**: Sensitive values are properly marked and handled
- 📋 **Complete**: Full CRUD operations with import support

## Quick Start

### Installation

#### Terraform 0.13+

```hcl
terraform {
  required_providers {
    openrouter = {
      source  = "standujar/openrouter"
      version = "~> 0.1.0"
    }
  }
}
```

#### Manual Installation

1. Download the latest release from [GitHub Releases](https://github.com/standujar/terraform-provider-openrouter/releases)
2. Extract and place in your Terraform plugins directory

### Authentication

Set your OpenRouter API key as an environment variable:

```bash
export OPENROUTER_API_KEY="sk-or-v1-your-api-key-here"
```

Or configure it directly in your Terraform configuration:

```hcl
provider "openrouter" {
  api_key = "sk-or-v1-your-api-key-here"
}
```

> **Note**: Using environment variables is recommended for security.

## Usage

### Basic Example

```hcl
# Configure the provider
provider "openrouter" {
  # API key via OPENROUTER_API_KEY environment variable
}

# Create an API key with spending limit
resource "openrouter_api_key" "my_app" {
  name  = "My Application"
  limit = 100.0  # $100 USD spending limit
}

# Output the key (sensitive)
output "api_key" {
  value     = openrouter_api_key.my_app.key
  sensitive = true
}
```

### Advanced Example

```hcl
# Create API key with time and spending limits
resource "openrouter_api_key" "limited_key" {
  name          = "Limited Access Key"
  limit         = 50.0    # $50 spending limit
  limit_minutes = 1440    # 24 hours time limit
}

# Get information about existing keys
data "openrouter_api_keys" "all" {
  include_disabled = true
}

# Get specific key details
data "openrouter_api_key" "existing" {
  id = "your-key-hash-here"
}

# Output key summary
output "key_summary" {
  value = {
    total_keys    = length(data.openrouter_api_keys.all.keys)
    active_keys   = length([for k in data.openrouter_api_keys.all.keys : k if !k.is_disabled])
    total_usage   = sum([for k in data.openrouter_api_keys.all.keys : k.usage])
  }
}
```

## Resources

### `openrouter_api_key`

Manages an OpenRouter API key.

#### Arguments

- `name` (String, Required) - The name of the API key
- `limit` (Number, Optional) - Spending limit in USD
- `limit_minutes` (Number, Optional) - Time limit in minutes
- `is_disabled` (Boolean, Optional) - Whether the key is disabled (default: false)

#### Attributes

- `id` (String) - The unique hash identifier of the API key
- `key` (String, Sensitive) - The actual API key value (only available during creation)
- `usage` (Number) - Current usage in USD
- `created_at` (String) - Creation timestamp

#### Import

```bash
terraform import openrouter_api_key.example your-key-hash-here
```

## Data Sources

### `openrouter_api_key`

Retrieves information about a specific API key.

#### Arguments

- `id` (String, Required) - The hash identifier of the API key

#### Attributes

- `name` (String) - The name of the API key
- `limit` (Number) - Spending limit in USD
- `limit_minutes` (Number) - Time limit in minutes
- `usage` (Number) - Current usage in USD
- `is_disabled` (Boolean) - Whether the key is disabled
- `is_provisioner` (Boolean) - Whether this is a provisioner key
- `created_at` (String) - Creation timestamp

### `openrouter_api_keys`

Retrieves a list of all API keys.

#### Arguments

- `include_disabled` (Boolean, Optional) - Include disabled keys in results (default: false)

#### Attributes

- `keys` (List of Objects) - List of API keys with the same attributes as `openrouter_api_key`

## Configuration Reference

### Provider Configuration

```hcl
provider "openrouter" {
  api_key  = "sk-or-v1-your-api-key-here"  # Optional, can use OPENROUTER_API_KEY env var
  endpoint = "https://openrouter.ai/api/v1" # Optional, defaults to official API
}
```

#### Arguments

- `api_key` (String, Optional) - OpenRouter API key. Can also be set via `OPENROUTER_API_KEY` environment variable
- `endpoint` (String, Optional) - Custom API endpoint URL. Defaults to `https://openrouter.ai/api/v1`

## Examples

See the [`examples/`](examples/) directory for complete usage examples:

- [`create-api-key.tf`](examples/create-api-key.tf) - Creating API keys with limits
- [`get-specific-key.tf`](examples/get-specific-key.tf) - Retrieving specific key information
- [`list-all-keys.tf`](examples/list-all-keys.tf) - Listing and filtering API keys

## Development

### Requirements

- [Go](https://golang.org/doc/install) >= 1.19
- [Terraform](https://www.terraform.io/downloads.html) >= 1.0

### Building

```bash
go build -o terraform-provider-openrouter
```

### Testing

```bash
# Unit tests
go test ./...

# Acceptance tests (requires OPENROUTER_API_KEY)
TF_ACC=1 go test ./... -v
```

### Local Development

1. Build the provider:
   ```bash
   go build -o terraform-provider-openrouter
   ```

2. Create local plugin directory:
   ```bash
   mkdir -p ~/.terraform.d/plugins/registry.terraform.io/standujar/openrouter/0.1.0/linux_amd64
   ```

3. Copy the binary:
   ```bash
   cp terraform-provider-openrouter ~/.terraform.d/plugins/registry.terraform.io/standujar/openrouter/0.1.0/linux_amd64/
   ```

4. Use in your Terraform configuration with version `0.1.0`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Guidelines

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure your code:
- Follows Go best practices
- Includes appropriate tests
- Updates documentation as needed
- Passes all CI checks

## Support

- 📚 [Documentation](https://registry.terraform.io/providers/standujar/openrouter/latest/docs)
- 🐛 [Issues](https://github.com/standujar/terraform-provider-openrouter/issues)
- 💡 [Feature Requests](https://github.com/standujar/terraform-provider-openrouter/issues/new?template=feature_request.md)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [OpenRouter](https://openrouter.ai/) for providing the API
- [HashiCorp](https://www.hashicorp.com/) for Terraform and the Provider SDK
- The Terraform community for inspiration and best practices

---

Made with ❤️ for the Terraform and OpenRouter communities.