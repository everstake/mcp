package mcp

import (
	"context"
	"fmt"
	"slices"
	"strings"

	mcp_server "mcp-server"
	"mcp-server/pkg/everstake/dashboard"
	"mcp-server/pkg/log"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	globalStatsCacheKey = "dashboard-global-stats"
	chainsCacheKey      = "dashboard-chains"

	percentDivisor = 100.0
	monthsPerYear  = 12.0
)

func (s *MCPServer) GetUptimeMetrics(ctx context.Context, _ *sdkmcp.CallToolRequest) (*sdkmcp.CallToolResult, error) {
	if cached, found := s.cache.Get(globalStatsCacheKey); found {
		val, ok := cached.(*dashboard.GlobalStatsValue)
		if ok {
			uptimeStr := fmt.Sprint(val.ServiceUptime)
			return newTextResult(uptimeStr), nil
		}
		log.Logger.Warn("invalid cache type for global stats")
	}

	globalStats, err := s.dashboard.GetGlobalStats(ctx)
	if err != nil {
		log.Logger.Error("failed to get global stats from dashboard", log.E(err))
		return nil, ErrFailedToFetchDashboard
	}

	return newTextResult(fmt.Sprint(globalStats.ServiceUptime)), nil
}

func (s *MCPServer) GetChains(ctx context.Context, _ *sdkmcp.CallToolRequest) (*sdkmcp.CallToolResult, error) {
	chains, err := s.fetchChains(ctx)
	if err != nil {
		log.Logger.Error("failed to get chains from dashboard", log.E(err))
		return nil, ErrFailedToFetchDashboard
	}

	return newJSONResult(chains), nil
}

type StakingCalculatorRequest struct {
	Network  string  `json:"network" description:"The blockchain network to stake on, e.g. Ethereum, Solana, etc." required:"false"`
	Currency string  `json:"currency" description:"The ticker symbol of the cryptocurrency to stake, e.g. ETH, SOL, etc." required:"false"`
	Amount   float64 `json:"amount" description:"The amount of cryptocurrency to stake."`
}

type StakingCalculatorResponse struct {
	Network                 string  `json:"network"`
	Ticker                  string  `json:"ticker"`
	CurrentAPYUsed          string  `json:"current_apy_used"`
	Currency                string  `json:"currency"`
	StakingURL              string  `json:"staking_url"`
	Disclaimer              string  `json:"disclaimer"`
	AmountStaked            float64 `json:"amount_staked"`
	EstimatedAnnualRewards  float64 `json:"estimated_annual_rewards"`
	EstimatedMonthlyRewards float64 `json:"estimated_monthly_rewards"`
}

func (s *MCPServer) StakingCalculator(ctx context.Context, _ *sdkmcp.CallToolRequest, input StakingCalculatorRequest) (*sdkmcp.CallToolResult, any, error) {
	chains, err := s.fetchChains(ctx)
	if err != nil {
		log.Logger.Error("failed to get chains for staking calculator", log.E(err))
		return nil, nil, ErrFailedToFetchDashboard
	}

	// currency/chain matching either currency/chain in chains
	chainIx := slices.IndexFunc(chains, func(ch dashboard.Chain) bool {
		return strings.EqualFold(ch.CurrencyCode, input.Currency) || strings.EqualFold(ch.Chain, input.Network) ||
			strings.EqualFold(ch.Chain, input.Currency) || strings.EqualFold(ch.CurrencyCode, input.Network)
	})
	if chainIx == -1 {
		log.Logger.Warn("unsupported currency for staking calculator", log.V("currency", input.Currency))
		return nil, nil, fmt.Errorf("unsupported currency: %s", input.Currency)
	}
	chain := chains[chainIx]

	apr := chain.Apr
	annualRewards := input.Amount * (apr / percentDivisor)
	monthlyRewards := annualRewards / monthsPerYear

	stakingURL := fmt.Sprintf("https://stake.everstake.one/dashboard/stake/%s/", strings.ToLower(chain.Chain))

	response := StakingCalculatorResponse{
		Network:                 chain.Chain,
		Ticker:                  strings.ToUpper(chain.CurrencyCode),
		AmountStaked:            input.Amount,
		CurrentAPYUsed:          fmt.Sprintf("%.2f%%", apr),
		EstimatedAnnualRewards:  annualRewards,
		EstimatedMonthlyRewards: monthlyRewards,
		Currency:                strings.ToUpper(chain.CurrencyCode),
		StakingURL:              stakingURL,
		Disclaimer:              "Estimates are based on current network APY and are not guaranteed. APY fluctuates based on network conditions, total stake, and validator performance.",
	}

	return nil, response, nil
}

func (s *MCPServer) RequestIntegration(ctx context.Context, _ *sdkmcp.CallToolRequest, input *dashboard.PDLead) (*sdkmcp.CallToolResult, any, error) {
	if input.LeadSource == "" {
		input.LeadSource = "MCP service (source not specified)"
	}

	err := s.dashboard.CreatePDLead(ctx, input)
	if err != nil {
		log.Logger.Error("failed to create pd lead", log.E(err))
		return &sdkmcp.CallToolResult{
			IsError: true,
			Content: []sdkmcp.Content{
				&sdkmcp.TextContent{Text: "Submission failed. Please try again or contact Everstake directly at https://everstake.one/contact-us"},
			},
		}, nil, nil
	}

	return newTextResult("Your inquiry has been submitted. Everstake's team will be in touch shortly."), nil, nil
}

func (s *MCPServer) fetchChains(ctx context.Context) ([]dashboard.Chain, error) {
	if cached, found := s.cache.Get(chainsCacheKey); found {
		if chains, ok := cached.([]dashboard.Chain); ok {
			return chains, nil
		}
		log.Logger.Warn("invalid cache type for chains")
	}

	chains, err := s.dashboard.GetChains(ctx)
	if err != nil {
		return nil, err
	}

	// Clear unnecessary data
	for i := range chains {
		chains[i].LogoBlack = ""
		chains[i].LogoWhite = ""
	}

	s.cache.Set(chainsCacheKey, chains, mcp_server.DashboardCacheTTL)
	return chains, nil
}
