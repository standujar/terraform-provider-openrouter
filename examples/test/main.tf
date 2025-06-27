# Test Configuration with Real API Key

terraform {
  required_providers {
    openrouter = {
      source  = "standujar/openrouter"
      version = "~> 0.1.0"
    }
  }
}

provider "openrouter" {
  api_key = "your-api-key"
}

# Create a test API key
resource "openrouter_api_key" "test_key" {
  name  = "Terraform Test Key"
  limit = 10.0
}

output "api_key_id" {
  description = "The hash ID of the created API key"
  value       = openrouter_api_key.test_key.id
}

output "api_key_value" {
  description = "The actual API key value (only shown during creation)"
  value       = openrouter_api_key.test_key.key
  sensitive   = true
}
