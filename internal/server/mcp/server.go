package mcp

import (
	"net/http"
	"time"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	mcp_server "mcp-server"
	"mcp-server/internal/config"
	"mcp-server/pkg/everstake/dashboard"

	"github.com/patrickmn/go-cache"
)

const (
	defaultTtl = 30 * time.Minute
)

// Server wraps the MCP SDK server and exposes an http.Handler.
type MCPServer struct {
	s         *sdkmcp.Server
	mcpConfig *config.ToolsConfig
	dashboard *dashboard.Dashboard

	cache *cache.Cache
}

func New(cfg *config.ToolsConfig, dashboard *dashboard.Dashboard) (*MCPServer, error) {
	s := sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    mcp_server.ServiceName,
		Version: mcp_server.Version,
	}, nil)
	mcps := &MCPServer{
		s:         s,
		mcpConfig: cfg,
		dashboard: dashboard,
		cache:     cache.New(defaultTtl, defaultTtl),
	}

	s.AddTool(mcps.mcpConfig.GetApiDocs.ToTool(), mcps.getApiDocs)

	return mcps, nil
}

func (s *MCPServer) Handler() http.Handler {
	return sdkmcp.NewStreamableHTTPHandler(func(_ *http.Request) *sdkmcp.Server {
		return s.s
	}, nil)
}
