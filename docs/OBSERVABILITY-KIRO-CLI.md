# Kiro CLI Integration

Use Kiro CLI with wetwire-observability for AI-assisted monitoring configuration in enterprise environments.

## Prerequisites

- Go 1.23+ installed
- Kiro CLI installed ([installation guide](https://kiro.dev/docs/cli/installation/))
- AWS Builder ID or GitHub/Google account (for Kiro authentication)

---

## Step 1: Install wetwire-obs

### Option A: Using Go (recommended)

```bash
go install github.com/lex00/wetwire-observability-go/cmd/wetwire-obs@latest
```

### Option B: Pre-built binaries

Download from [GitHub Releases](https://github.com/lex00/wetwire-observability-go/releases):

```bash
# macOS (Apple Silicon)
curl -LO https://github.com/lex00/wetwire-observability-go/releases/latest/download/wetwire-obs-darwin-arm64
chmod +x wetwire-obs-darwin-arm64
sudo mv wetwire-obs-darwin-arm64 /usr/local/bin/wetwire-obs

# Linux (x86-64)
curl -LO https://github.com/lex00/wetwire-observability-go/releases/latest/download/wetwire-obs-linux-amd64
chmod +x wetwire-obs-linux-amd64
sudo mv wetwire-obs-linux-amd64 /usr/local/bin/wetwire-obs
```

### Verify installation

```bash
wetwire-obs --version
```

---

## Step 2: Install Kiro CLI

```bash
# Install Kiro CLI
curl -fsSL https://cli.kiro.dev/install | bash

# Verify installation
kiro-cli --version

# Sign in (opens browser)
kiro-cli login
```

---

## Step 3: Configure Kiro for wetwire-obs

Run the design command with `--provider kiro` to auto-configure:

```bash
# Create a project directory
mkdir my-monitoring && cd my-monitoring

# Initialize Go module
go mod init my-monitoring

# Run design with Kiro provider (auto-installs configs on first run)
wetwire-obs design --provider kiro "Create alerts for my API"
```

This automatically installs:

| File | Purpose |
|------|---------|
| `~/.kiro/agents/wetwire-obs-runner.json` | Kiro agent configuration |
| `.kiro/mcp.json` | Project MCP server configuration |

### Manual configuration (optional)

**~/.kiro/agents/wetwire-obs-runner.json:**
```json
{
  "name": "wetwire-obs-runner",
  "description": "Observability config generator using wetwire-observability",
  "prompt": "You are an observability configuration assistant...",
  "model": "claude-sonnet-4",
  "mcpServers": {
    "wetwire": {
      "command": "wetwire-obs",
      "args": ["mcp"],
      "cwd": "/path/to/your/project"
    }
  },
  "tools": ["*"]
}
```

---

## Step 4: Run Kiro with wetwire design

### Using the wetwire-obs CLI

```bash
# Start Kiro design session
wetwire-obs design --provider kiro "Create monitoring for a microservices API"
```

### Using Kiro CLI directly

```bash
# Start chat with wetwire-obs-runner agent
kiro-cli chat --agent wetwire-obs-runner

# Or with an initial prompt
kiro-cli chat --agent wetwire-obs-runner "Create error rate alerts"
```

---

## Available MCP Tools

The wetwire-obs MCP server exposes these tools to Kiro:

| Tool | Description | Example |
|------|-------------|---------|
| `wetwire_init` | Initialize a new project | `wetwire_init(path="./monitoring")` |
| `wetwire_lint` | Lint configs for issues | `wetwire_lint(path="./...")` |
| `wetwire_build` | Generate config files | `wetwire_build(path="./...", mode="standalone")` |

---

## Example Session

```
$ wetwire-obs design --provider kiro "Create monitoring for an API with error rate and latency alerts"

Installed Kiro agent config: ~/.kiro/agents/wetwire-obs-runner.json
Installed project MCP config: .kiro/mcp.json
Starting Kiro CLI design session...

> I'll help you create monitoring for your API with error rate and latency alerts.

Let me create the configuration files.

[Calling wetwire_init...]
[Calling wetwire_lint...]
[Calling wetwire_build...]

I've created the following files:
- monitoring/scrapes.go - Prometheus scrape configuration
- monitoring/alerts.go - Error rate and latency alerts
- monitoring/dashboard.go - Grafana dashboard

The alerts include:
- HighErrorRate - Fires when error rate exceeds 5%
- HighLatency - Fires when p99 latency exceeds 500ms

Would you like me to add any additional alerts or dashboards?
```

---

## Workflow

The Kiro agent follows this workflow:

1. **Explore** - Understand your monitoring requirements
2. **Plan** - Design the alert and dashboard structure
3. **Implement** - Generate Go code using wetwire-observability patterns
4. **Lint** - Run `wetwire_lint` to check for issues
5. **Build** - Run `wetwire_build` to generate configs

---

## Deploying Generated Configs

After Kiro generates your monitoring configuration:

```bash
# Build standalone configs
wetwire-obs build ./monitoring --mode=standalone -o ./output/

# Deploy Prometheus config
cp output/prometheus.yml /etc/prometheus/prometheus.yml
systemctl reload prometheus

# Deploy Alertmanager config
cp output/alertmanager.yml /etc/alertmanager/alertmanager.yml
systemctl reload alertmanager

# Import Grafana dashboards
# (Use Grafana provisioning or API)
```

### For Kubernetes with Prometheus Operator

```bash
# Build Operator CRDs
wetwire-obs build ./monitoring --mode=operator -o ./manifests/

# Apply to cluster
kubectl apply -f manifests/
```

---

## Troubleshooting

### MCP server not found

```
Mcp error: -32002: No such file or directory
```

**Solution:** Ensure `wetwire-obs` is in your PATH:

```bash
which wetwire-obs

# If not found, add to PATH or reinstall
go install github.com/lex00/wetwire-observability-go/cmd/wetwire-obs@latest
```

### Kiro CLI not found

```
kiro-cli not found in PATH
```

**Solution:** Install Kiro CLI:

```bash
curl -fsSL https://cli.kiro.dev/install | bash
```

### Authentication issues

```
Error: Not authenticated
```

**Solution:** Sign in to Kiro:

```bash
kiro-cli login
```

---

## Known Limitations

### Automated Testing

When using `wetwire-obs test --provider kiro`, tests run in non-interactive mode. For true persona simulation with multi-turn conversations, use the Anthropic provider:

```bash
wetwire-obs test --provider anthropic --persona expert "Create alerts"
```

### Interactive Design Mode

Interactive design mode (`wetwire-obs design --provider kiro`) works fully as expected with real-time conversation.

---

## See Also

- [CLI Reference](CLI.md) - Full wetwire-obs CLI documentation
- [Quick Start](QUICK_START.md) - Getting started with wetwire-observability
- [Kiro CLI Installation](https://kiro.dev/docs/cli/installation/) - Official installation guide
- [Kiro CLI Docs](https://kiro.dev/docs/cli/) - Official Kiro documentation
