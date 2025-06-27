# Example: Get Specific OpenRouter API Key

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

# Retrieve information about a specific API key by its hash ID
data "openrouter_api_key" "existing_key" {
  # Replace with your actual API key hash
  id = "your-api-key-hash-here"
}

# Output the key information
output "key_name" {
  description = "Name of the API key"
  value       = data.openrouter_api_key.existing_key.name
}

output "key_limit" {
  description = "Spending limit of the API key"
  value       = data.openrouter_api_key.existing_key.limit
}

output "key_usage" {
  description = "Current usage of the API key"
  value       = data.openrouter_api_key.existing_key.usage
}

output "key_disabled" {
  description = "Whether the API key is disabled"
  value       = data.openrouter_api_key.existing_key.is_disabled
}

output "key_created_at" {
  description = "When the API key was created"
  value       = data.openrouter_api_key.existing_key.created_at
}

output "is_provisioner" {
  description = "Whether this is a provisioner key"
  value       = data.openrouter_api_key.existing_key.is_provisioner
}