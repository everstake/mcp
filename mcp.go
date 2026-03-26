package mcp_server

import (
	_ "embed"
	"time"
)

//go:embed tools.yaml
var MCPConfig []byte

const (
	Version     = "1.0.0"
	ServiceName = "Everstake MCP"

	// mcp server cache
	DashboardCacheTtl = 10 * time.Minute
)
