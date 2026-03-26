# Everstake MCP Server

An [MCP (Model Context Protocol)](https://modelcontextprotocol.io) server that makes Everstake accessible to AI agents. Built with [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk).

## Structure

```
.
├── tools.yaml                  # Tool definitions (name, description, args)
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── config/loader.go        # YAML loader + ToTools()
│   └── tools/everstake.go      # Tool handlers
└── go.mod
```

Tools are defined in `tools.yaml` — add a new tool there and wire up its handler in `everstake.go`. No boilerplate elsewhere.

## Tools

| Tool | Description |
|---|---|
| `get_api_docs` | Returns the URL to Everstake's developer documentation |

## Run

```bash
go mod tidy
go run ./cmd/mcp_server
```

Server starts on `:8080` by default (SSE transport on `/sse`).

## Add a Tool

**1. Define it in `tools.yaml`:**

```yaml
tools:
  - name: get_networks
    description: Returns all networks supported by Everstake.
    args:
      - name: ecosystem
        type: string
        required: false
        description: Filter by ecosystem (e.g. "evm", "cosmos").
```

**2. Add a handler in `internal/tools/everstake.go`:**

```go
func HandleGetNetworks(_ context.Context, _ *mcp.ServerSession, _ *mcp.CallToolParamsFor[map[string]any]) (*mcp.CallToolResultFor[any], error) {
    // ...
    return &mcp.CallToolResultFor[any]{
        Content: []mcp.Content{&mcp.TextContent{Text: result}},
    }, nil
}
```

**3. Register it in `cmd/mcp_server/main.go`:**

```go
handlers := map[string]mcp.ToolHandler{
    "get_api_docs": tools.HandleGetAPIDocs,
    "get_networks": tools.HandleGetNetworks,
}
```

## Dependencies

- [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk) — official MCP Go SDK
- [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) — YAML config parsing
