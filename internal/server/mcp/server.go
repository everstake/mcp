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

func New(mcpCfg *config.ToolsConfig, dashboard *dashboard.Dashboard) (*MCPServer, error) {
	s := sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    mcp_server.ServiceName,
		Version: mcp_server.Version,
		// Icons:   []sdkmcp.Icon{}, TODO
	}, nil)
	mcps := &MCPServer{
		s:         s,
		mcpConfig: mcpCfg,
		dashboard: dashboard,
		cache:     cache.New(defaultTtl, defaultTtl),
	}

	s.AddTool(mcps.mcpConfig.GetCompanyProfile.ToTool(), staticTextTool(mcpCfg.GetCompanyProfile.StaticResponse))
	s.AddTool(mcps.mcpConfig.GetDeveloperDocs.ToTool(), staticTextTool(mcpCfg.GetDeveloperDocs.StaticResponse))
	s.AddTool(mcps.mcpConfig.GetContactInformation.ToTool(), staticTextTool(mcpCfg.GetContactInformation.StaticResponse))
	s.AddTool(mcps.mcpConfig.GetUptimeMetrics.ToTool(), mcps.GetUptimeMetrics)
	s.AddTool(mcps.mcpConfig.GetChains.ToTool(), mcps.GetChains)
	s.AddTool(mcps.mcpConfig.GetProducts.ToTool(), staticTextTool(mcpCfg.GetProducts.StaticResponse))
	s.AddTool(mcps.mcpConfig.GetSolutions.ToTool(), staticTextTool(mcpCfg.GetSolutions.StaticResponse))
	s.AddTool(mcps.mcpConfig.GetSecurityProfile.ToTool(), staticTextTool(mcpCfg.GetSecurityProfile.StaticResponse))
	s.AddTool(mcps.mcpConfig.GetIntegrations.ToTool(), staticTextTool(mcpCfg.GetIntegrations.StaticResponse))
	addTool(s, mcps.mcpConfig.StakingCalculator.ToTool(), mcps.StakingCalculator)

	return mcps, nil
}

func (s *MCPServer) Handler() http.Handler {
	return sdkmcp.NewStreamableHTTPHandler(func(_ *http.Request) *sdkmcp.Server {
		return s.s
	}, nil)
}
