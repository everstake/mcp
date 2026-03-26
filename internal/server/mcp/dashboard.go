package mcp

import (
	"context"
	"fmt"
	mcp_server "mcp-server"
	"mcp-server/pkg/everstake/dashboard"
	"mcp-server/pkg/log"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	globalStatsCacheKey = "dashboard-global-stats"
	chainsCacheKey      = "dashboard-chains"
)

func (s *MCPServer) GetUptimeMetrics(context.Context, *sdkmcp.CallToolRequest) (*sdkmcp.CallToolResult, error) {
	if cached, found := s.cache.Get(globalStatsCacheKey); found {
		val := cached.(*dashboard.GlobalStatsValue)
		uptimeStr := fmt.Sprint(val.ServiceUptime)
		return newTextResult(uptimeStr), nil
	}

	globalStats, err := s.dashboard.GetGlobalStats()
	if err != nil {
		log.Logger.Error("failed to get global stats from dashboard", log.E(err))
		return nil, ErrFailedToFetchDashboard
	}

	s.cache.Set(globalStatsCacheKey, globalStats, mcp_server.DashboardCacheTtl)
	return newTextResult(fmt.Sprint(globalStats.ServiceUptime)), nil
}

func (s *MCPServer) GetChains(context.Context, *sdkmcp.CallToolRequest) (*sdkmcp.CallToolResult, error) {
	if cached, found := s.cache.Get(chainsCacheKey); found {
		val := cached.([]dashboard.Chain)
		return newJsonResult(val), nil
	}

	chains, err := s.dashboard.GetChains()
	if err != nil {
		log.Logger.Error("failed to get chains from dashboard", log.E(err))
		return nil, ErrFailedToFetchDashboard
	}

	// clear big nothing
	for i := range chains {
		chains[i].LogoBlack = ""
		chains[i].LogoWhite = ""
	}

	s.cache.Set(chainsCacheKey, chains, mcp_server.DashboardCacheTtl)
	return newJsonResult(chains), nil
}
