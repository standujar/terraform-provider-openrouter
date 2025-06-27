# Create a new API key
resource "openrouter_api_key" "example" {
  name = "My Terraform API Key"
  
  # Optional: Set a spend limit in USD
  limit = 100.0
  
  # Optional: Set a time limit in minutes
  # limit_minutes = 1440  # 24 hours
}

# Output the API key (sensitive)
output "api_key" {
  value     = openrouter_api_key.example.key
  sensitive = true
}

# Data source to fetch an existing API key
data "openrouter_api_key" "existing" {
  id = "sk-or-v1-1234567890abcdef"  # Replace with actual key hash
}

# Data source to list all API keys
data "openrouter_api_keys" "all" {
  include_disabled = true
}

# Output the list of API keys
output "all_keys" {
  value = [for key in data.openrouter_api_keys.all.keys : {
    name = key.name
    id   = key.id
    usage = key.usage
    disabled = key.is_disabled
  }]
}