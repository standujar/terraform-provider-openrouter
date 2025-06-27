# Terraform Provider for OpenRouter

A Terraform provider for managing OpenRouter resources.

## TODO List

### Initial Setup
- [ ] Initialize Go module (`go mod init github.com/[username]/terraform-provider-openrouter`)
- [ ] Add Terraform Plugin SDK v2 dependency
- [ ] Create basic project structure
- [ ] Set up GitHub repository
- [ ] Configure .gitignore for Go and Terraform

### Provider Implementation
- [ ] Create `main.go` entry point
- [ ] Implement `provider.go` with schema and configuration
- [ ] Create OpenRouter API client (`internal/client/client.go`)
- [ ] Add authentication handling (API key)
- [ ] Implement error handling and retries

### Resources
- [ ] `openrouter_api_key` resource
  - [ ] Create
  - [ ] Read
  - [ ] Update
  - [ ] Delete
  - [ ] Import
- [ ] `openrouter_route` resource (if applicable)
- [ ] `openrouter_model_preference` resource (if applicable)

### Data Sources
- [ ] `openrouter_models` data source (list available models)
- [ ] `openrouter_account` data source (account information)
- [ ] `openrouter_usage` data source (usage statistics)

### Testing
- [ ] Unit tests for client
- [ ] Unit tests for resources
- [ ] Unit tests for data sources
- [ ] Acceptance tests setup
- [ ] Acceptance tests for all resources
- [ ] Acceptance tests for all data sources

### Documentation
- [ ] Provider documentation
- [ ] Resource documentation
- [ ] Data source documentation
- [ ] Examples directory with use cases
- [ ] Contributing guide

### CI/CD
- [ ] GitHub Actions workflow for testing
- [ ] Linting workflow (golangci-lint)
- [ ] Release workflow
- [ ] Goreleaser configuration

### Publishing
- [ ] GPG key for signing releases
- [ ] Terraform Registry manifest
- [ ] Initial release (v0.1.0)
- [ ] Announce on Terraform Registry

### Future Enhancements
- [ ] Rate limiting handling
- [ ] Comprehensive error messages
- [ ] Debug logging
- [ ] Performance optimizations
- [ ] Support for bulk operations