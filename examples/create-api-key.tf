# Example: Create OpenRouter API Key

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

# Create a new API key with spending limit
resource "openrouter_api_key" "my_key" {
  name = "My Application Key"
  
  # Optional: Set a spending limit in USD
  limit = 50.0
  
  # Optional: Set a time limit in minutes (1440 = 24 hours)
  # limit_minutes = 1440
}

# Output the created key details
output "api_key_id" {
  description = "The hash ID of the created API key"
  value       = openrouter_api_key.my_key.id
}

output "api_key_value" {
  description = "The actual API key value (only shown during creation)"
  value       = openrouter_api_key.my_key.key
  sensitive   = true
}

output "api_key_usage" {
  description = "Current usage of the API key"
  value       = openrouter_api_key.my_key.usage
}