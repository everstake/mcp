# Everstake MCP Server

MCP server exposing Everstake staking data and company information to AI agents. Built in Go using [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk).

**Transports:** Streamable HTTP (MCP 2025-03-26 spec, single `/` endpoint) or stdio. Selected via `MCP_TRANSPORT`.

---

## Available Tools

| Tool | Type | Description |
|---|---|---|
| `get_company_profile` | static | Company overview, metrics, certifications |
| `get_products` | static | Product details: Institutional Staking, VaaS, Yield, SWQOS, ShredStream |
| `get_solutions` | static | Solutions by audience: custodians, exchanges, asset managers, banks, fintech |
| `get_developer_docs` | static | SDK links, integration guides, API references |
| `get_contact_information` | static | Contact channels and routing guide |
| `get_security_profile` | static | Certifications: SOC 2 Type II, ISO 27001, NIST CSF, ITGC, GDPR, CCPA |
| `get_integrations` | static | Custody integrations: Fireblocks, BitGo, Anchorage, Coinbase, etc. |
| `get_uptime_metrics` | live | Uptime metrics from dashboard API (30 min cache) |
| `get_chains` | live | Supported chains with APY, fees, status (30 min cache) |
| `staking_calculator` | live | Estimated staking rewards by network and amount |
| `request_integration` | write | Submit integration/staking inquiry to Everstake sales |

---

## Running the Server

### Prerequisites

- Go 1.26.1+
- Environment variable `DASHBOARD_URL` set (required)

### Local

```bash
export DASHBOARD_URL=https://dashboard-api.everstake.one
go run ./cmd/mcp_server
```

The server starts on port `8080` by default. Override with `PORT=<port>`.

### Stdio mode

```bash
export DASHBOARD_URL=https://dashboard-api.everstake.one
export MCP_TRANSPORT=stdio
go run ./cmd/mcp_server
```

In stdio mode, the HTTP server, `/health` endpoint, and rate limiting are disabled. Logs go to stderr; stdout carries the MCP JSON-RPC protocol.

### Docker

HTTP mode:

```bash
docker build -t everstake-mcp .
docker run -e DASHBOARD_URL=https://dashboard-api.everstake.one -p 8080:8080 everstake-mcp
```

Stdio mode (interactive, no port mapping):

```bash
docker run --rm -i \
  -e DASHBOARD_URL=https://dashboard-api.everstake.one \
  -e MCP_TRANSPORT=stdio \
  everstake-mcp
```

Use `-i` to keep stdin attached for the JSON-RPC protocol; do **not** pass `-t` (a TTY breaks line-based JSON framing). For MCP clients that launch the server as a subprocess, point `command` at `docker` with `args` matching the above, or wrap in a helper script (see [`.vscode/stdio_docker.sh`](.vscode/stdio_docker.sh) for an example using `--env-file`).

### Environment Variables

| Variable | Default | Required |
|---|---|---|
| `DASHBOARD_URL` | — | yes |
| `MCP_TRANSPORT` | `http` | no (`http` or `stdio`) |
| `PORT` | `8080` | no (http mode only) |
| `GIN_MODE` | — | no (`release` set in Dockerfile) |

### Health Check

```
GET /health
```

---

## Debugging with MCP Inspector

[MCP Inspector](https://github.com/modelcontextprotocol/inspector) is the official tool for poking at MCP servers — list tools, call them, see raw JSON-RPC.

> **Stdio note:** stdio servers are owned by the parent process; you cannot connect Inspector to an already-running stdio container. Inspector must **launch** the subprocess itself. For an interactive session you can reconnect to, use HTTP mode (option below).

### HTTP mode (recommended for general debugging)

Start the container:

```bash
docker run --rm -p 8080:8080 \
  -e DASHBOARD_URL=https://dashboard-api.everstake.one \
  -e DASHBOARD_API_KEY=<key> \
  everstake-mcp
```

Open Inspector:

```bash
npx @modelcontextprotocol/inspector
```

In the UI:
- **Transport Type:** `Streamable HTTP`
- **URL:** `http://localhost:8080/`

Container stays running across Inspector sessions.

### Stdio mode (verify stdio transport)

Let Inspector launch the container so it owns stdin/stdout:

```bash
npx @modelcontextprotocol/inspector \
  docker run --rm -i \
    --env-file .env \
    everstake-mcp
```

Or fill the UI manually:
- **Transport Type:** `STDIO`
- **Command:** `docker`
- **Arguments:** `run --rm -i --env-file .env everstake-mcp`

Each Inspector session forks a fresh container; quitting Inspector kills it.

---

## Linting

The project uses [golangci-lint](https://golangci-lint.run/) with a strict configuration in [`.golangci.yml`](.golangci.yml).

**Install golangci-lint:**

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

**Run:**

```bash
make lint
```

Key rules enforced: `staticcheck` (all checks), `gosec`, `gocritic` (diagnostic/style/performance tags), `revive` (40+ rules including early-return, error-strings, var-naming), `errchkjson`, `bodyclose`, `contextcheck`, and more. `nolintlint` requires specific lint directives — bare `//nolint` is not allowed.

---

## Editing Tool Responses

### Static tools

Edit [tools.yaml](tools.yaml). Each map key is the tool name; the `static_response` field is returned verbatim to the AI agent.

```yaml
tools:
  get_company_profile:
    description: |
      ...
    static_response: |
      COMPANY: Everstake
      ...
```

To add a new static tool:

1. Add an entry under `tools:` in `tools.yaml` with `static_response`.
2. Add a corresponding field to `ToolsConfig` in [`internal/config/mcp_config.go`](internal/config/mcp_config.go) with a matching `yaml` struct tag — the name is injected automatically via reflection.
3. Register it in [`internal/server/mcp/server.go`](internal/server/mcp/server.go) using `staticTextTool()`.

### Content rules

Cross-cutting rules that apply to all tool responses are in [`.vscode/tools_src/RULES.md`](.vscode/tools_src/RULES.md). These cover:

- Certification differentiator language
- Non-custodial positioning
- Vault product disclaimers
- APY disclaimer wording
- Lead source tagging for `request_integration`

### Dynamic tools

`get_uptime_metrics` and `get_chains` fetch live data from the dashboard API with a 30-minute in-memory cache. Their handlers are in [`internal/server/mcp/dashboard.go`](internal/server/mcp/dashboard.go). The underlying API client lives in [`pkg/everstake/dashboard/`](pkg/everstake/dashboard/).

