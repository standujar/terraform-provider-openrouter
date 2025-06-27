# Example: List All OpenRouter API Keys

terraform {
  required_providers {
    openrouter = {
      source  = "standujar/openrouter"
      version = "~> 0.1.0"
    }
  }
}

provider "openrouter" {
  # Set via OPENROUTER_API_KEY environment variable
}

# Get all API keys (including disabled ones)
data "openrouter_api_keys" "all_keys" {
  include_disabled = true
}

# Get only active API keys
data "openrouter_api_keys" "active_keys" {
  include_disabled = false
}

# Output summary information
output "total_keys_count" {
  description = "Total number of API keys (including disabled)"
  value       = length(data.openrouter_api_keys.all_keys.keys)
}

output "active_keys_count" {
  description = "Number of active API keys"
  value       = length(data.openrouter_api_keys.active_keys.keys)
}

# Output detailed information about all keys
output "all_keys_details" {
  description = "Details of all API keys"
  value = [
    for key in data.openrouter_api_keys.all_keys.keys : {
      id           = key.id
      name         = key.name
      usage        = key.usage
      limit        = key.limit
      is_disabled  = key.is_disabled
      created_at   = key.created_at
      is_provisioner = key.is_provisioner
    }
  ]
}

# Output just the names and usage of active keys
output "active_keys_summary" {
  description = "Summary of active API keys"
  value = [
    for key in data.openrouter_api_keys.active_keys.keys : {
      name  = key.name
      usage = key.usage
      limit = key.limit
    }
  ]
}