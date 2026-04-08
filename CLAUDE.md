# Everstake MCP Server — Session Context

## Project

MCP server exposing Everstake data to AI agents. Built in Go using the official `modelcontextprotocol/go-sdk`.

**Repo root:** `mcp-server` module
**Entry point:** `cmd/mcp_server/main.go`

---

## Stack

| Concern | Choice |
|---|---|
| MCP SDK | `github.com/modelcontextprotocol/go-sdk v1.4.1` |
| HTTP router | `github.com/gin-gonic/gin v1.10.0` |
| MCP transport | Streamable HTTP only (`/` via `NewStreamableHTTPHandler`) |
| Config | `github.com/caarlos0/env` — env vars, no file |
| Cache | `github.com/patrickmn/go-cache` — in-memory, 10 min TTL |
| Logging | `github.com/sirupsen/logrus` |
| Go version | 1.26.1 |

---

## Structure

```
cmd/mcp_server/main.go          signal handling, wires deps, starts server
internal/
  config/
    service_config.go           ServiceConfig{Port, DashboardUrl} via caarlos0/env
    mcp_config.go               ToolsConfig, custom UnmarshalYAML, LoadMCPConfig
    config_unmarshal_test.go    tests for YAML unmarshaling
  server/
    server.go                   Gin router, /health, mounts MCP handler, graceful shutdown
    mcp/
      server.go                 MCPServer struct, New(), Handler(), tool registration
      dashboard.go              GetUptimeMetrics, GetChains handlers (live API + cache)
      utils.go                  newTextResult(), newJsonResult(), staticTextTool(), ErrFailedToFetchDashboard
    middleware/
      ratelimit.go              rate limiting middleware
pkg/
  everstake/dashboard/          HTTP client for dashboard-api.everstake.one
    dashboard.go                Dashboard client, base HTTP methods
    chains.go                   Chain types and API endpoints
    datafeed.go                 Data feed types and API endpoints
    leads.go                    Leads types and API endpoints
  log/
    log.go                      logrus wrapper
mcp.go                          embed tools.yaml, Version, ServiceName, DashboardCacheTtl consts
tools.yaml                      tool definitions (name derived from map key, not field)
Dockerfile                      multi-stage, golang:1.26.1-alpine builder + alpine:3.21 runtime
Makefile                        build automation (lint)
```

---

## Tools

Defined in `tools.yaml` as a map. The map key IS the tool name — injected automatically via `ToolsConfig.UnmarshalYAML` using reflection on yaml struct tags. No manual name-setting needed when adding tools.

| Tool | Handler | Source |
|---|---|---|
| `get_company_profile` | `staticTextTool` | static_response from yaml |
| `get_products` | `staticTextTool` | static_response from yaml |
| `get_solutions` | `staticTextTool` | static_response from yaml |
| `get_developer_docs` | `staticTextTool` | static_response from yaml |
| `get_contact_information` | `staticTextTool` | static_response from yaml |
| `get_security_profile` | `staticTextTool` | static_response from yaml |
| `get_integrations` | `staticTextTool` | static_response from yaml |
| `get_uptime_metrics` | `dashboard.go` | live API + cache |
| `get_chains` | `dashboard.go` | live API + cache |

**Static tools (7):** Company profile, products, solutions, developer docs, contact info, security profile, integrations
**Dynamic tools (2):** Uptime metrics, chains (from dashboard API with 30min cache)

---

## Key Design Decisions

**Transport: Streamable HTTP only**
SSE (legacy 2024-11-05 spec) was considered but dropped. Streamable HTTP (2025-03-26 spec) is a superset — single endpoint handles GET/POST/DELETE, proxy-friendly. SSE is deprecated. Claude Desktop will migrate; new clients already use Streamable HTTP.

**Tool name injection via reflection**
YAML uses map keys as tool names (`get_api_docs:` not `name: get_api_docs`). `ToolsConfig.UnmarshalYAML` reads each field's `yaml` struct tag and sets `ToolConfig.Name` automatically. Adding a new tool to `ToolsConfig` gets its name for free.

**`static_response` in yaml**
Tools with purely static responses (docs URL, contact links, company profile, etc.) declare `static_response` in yaml. No separate handler file needed. Registered as a closure in `MCPServer.New()` via `staticTextTool()` helper in `utils.go` that returns a `ToolHandler` reading the static text from config.

**Error handling**
Tool-level errors (API failures) → `IsError: true` in `Content` — LLM sees it and can self-correct. Protocol errors → return Go `error`. Sensitive details logged via logrus; sanitized message returned to LLM.

**`StructuredContent` must be a JSON object**
SDK validates this. `[]Chain` is an array → wrapped as `{"data": [...]}` in `newJsonResult`. `StructuredContent` is for the generic `ToolHandlerFor[In, Out]` typed path; current handlers use low-level `ToolHandler`.

**Gin + `gin.WrapH`**
Both MCP handlers implement `http.Handler`. Chi was considered (native `http.Handler`, zero friction) but Gin was chosen. `gin.WrapH` wraps the streamable handler with one line. `gin.New()` + `gin.Recovery()` only — no Gin logger (logrus handles it).

---

## MCP SDK Notes (v1.4.1)

```go
// ToolHandler signature (low-level)
type ToolHandler = func(context.Context, *CallToolRequest) (*CallToolResult, error)

// Server init
s := sdkmcp.NewServer(&sdkmcp.Implementation{Name: ..., Version: ...}, nil)
s.AddTool(tool, handler)

// Streamable HTTP — single endpoint, handles GET + POST + DELETE
sdkmcp.NewStreamableHTTPHandler(func(*http.Request) *sdkmcp.Server { return s }, nil)

// Tool definition
&sdkmcp.Tool{
    Name:        name,
    Description: desc,
    InputSchema: map[string]any{"type": "object"},
}

// Results
newTextResult(str)  → Content: []Content{&TextContent{Text: str}}
newJsonResult(obj)  → StructuredContent: map[string]any{"data": obj}  // must be object
```

---

## Config (env vars)

| Var | Default | Required |
|---|---|---|
| `PORT` | `8080` | no |
| `DASHBOARD_URL` | `https://dashboard-api.everstake.one` | yes |
| `GIN_MODE` | — | no (set to `release` in Dockerfile) |

---

## Makefile

Available commands:

| Command | Description |
|---|---|
| `make lint` | Run golangci-lint on the codebase. Requires golangci-lint to be installed (`go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`). |

Linting uses `.golangci.yml` configuration in the project root.

---

## Dockerfile

Two-stage: `golang:1.26.1-alpine` builder → `alpine:3.21` runtime.
`CGO_ENABLED=0` — no C deps.
`-ldflags="-s -w"` — strips debug symbols.
`ca-certificates` only runtime dep (outbound HTTPS to dashboard API).
Non-root user `app`.
