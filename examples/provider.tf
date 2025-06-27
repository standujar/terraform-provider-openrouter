# Basic OpenRouter Provider Configuration

terraform {
  required_version = ">= 1.0"

  required_providers {
    openrouter = {
      source  = "standujar/openrouter"
      version = "~> 0.1.0"
    }
  }
}

# Configure the OpenRouter Provider
provider "openrouter" {
  # API key can be provided here or via OPENROUTER_API_KEY environment variable
  # api_key = "sk-or-v1-your-api-key-here"

  # Optional: Override the default API endpoint
  # endpoint = "https://openrouter.ai/api/v1"
}
