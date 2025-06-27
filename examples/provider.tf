terraform {
  required_providers {
    openrouter = {
      source = "standujar/openrouter"
      version = "~> 0.1"
    }
  }
}

# Configure the OpenRouter Provider
provider "openrouter" {
  # API key can be set via OPENROUTER_API_KEY environment variable
  # api_key = "your-api-key"
  
  # Optional: Override the default endpoint
  # endpoint = "https://openrouter.ai/api/v1"
}