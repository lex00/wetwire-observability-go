<picture>
  <source media="(prefers-color-scheme: dark)" srcset="docs/wetwire-dark.svg">
  <img src="docs/wetwire-light.svg" width="100" height="67">
</picture>

Thank you for your interest in contributing to wetwire-observability-go!

## Getting Started

See the [Developer Guide](docs/DEVELOPERS.md) for:
- Development environment setup
- Project structure
- Running tests

## Code Style

- **Formatting**: Use `gofmt` (automatic in most editors)
- **Linting**: Use `go vet` and `golangci-lint`
- **Imports**: Use `goimports` for automatic import management

```bash
# Format code
gofmt -w .

# Lint
go vet ./...
golangci-lint run ./...

# Check for common issues
go build ./...
```

## Commit Messages

Follow conventional commits:

```
feat(prometheus): Add support for remote_write configuration
fix(grafana): Correct panel position calculation
docs: Update dashboard examples
test: Add tests for alertmanager routing
chore: Update dependencies
```

## Pull Request Process

1. Create feature branch: `git checkout -b feature/my-feature`
2. Make changes with tests
3. Run tests: `go test ./...`
4. Run linter: `golangci-lint run ./...`
5. Commit with clear messages
6. Push and open Pull Request
7. Address review comments

## Adding a New Lint Rule

1. Add rule to `internal/lint/rules.go`
2. Implement the check function
3. Add test case in `internal/lint/rules_test.go`
4. Update docs/LINT_RULES.md with the new rule

Lint rules use the `WOB` prefix (Wetwire OBservability). See [docs/LINT_RULES.md](docs/LINT_RULES.md) for the complete rule reference and category ranges.

## Adding a New Resource Type

When adding support for a new observability resource (e.g., new Grafana panel type):

1. Add type definition in appropriate package (`grafana/`, `prometheus/`, etc.)
2. Add discovery support in `internal/discover/`
3. Add serialization in `internal/serialize/`
4. Add lint rules if applicable
5. Add tests and examples

## Reporting Issues

- Use GitHub Issues for bug reports and feature requests
- Include reproduction steps for bugs
- Check existing issues before creating new ones

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
