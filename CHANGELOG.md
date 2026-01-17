# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- Migrated MCP server to use domain.BuildMCPServer() from wetwire-core-go v1.13.0
- Updated wetwire-core-go from v1.12.0 to v1.13.0
- Removed manual MCP implementation in favor of auto-generated MCP server

### Added

#### Phase 1: Foundation (Complete)
- Prometheus Duration type with serialization (30s, 5m, 1h30m format)
- Prometheus GlobalConfig for scrape/evaluation intervals
- Prometheus ScrapeConfig for job configuration
- Prometheus StaticConfig for static target groups
- Supporting types: BasicAuth, TLSConfig, RelabelConfig
- YAML serialization with Serialize() and SerializeToFile()
- AST-based resource discovery (internal/discover package)
- CLI commands: build, lint, list
- Table and JSON output formats for list command

#### Infrastructure
- Initial repository setup
- CI/CD workflows (build, test, release)
- Project documentation (README, CLAUDE.md, CHANGELOG)
