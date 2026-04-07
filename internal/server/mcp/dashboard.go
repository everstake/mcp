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

	return newTextResult(fmt.Sprint(globalStats.ServiceUptime)), nil
}

func (s *MCPServer) GetChains(context.Context, *sdkmcp.CallToolRequest) (*sdkmcp.CallToolResult, error) {
	chains, err := s.getChains()
	if err != nil {
		log.Logger.Error("failed to get chains from dashboard", log.E(err))
		return nil, ErrFailedToFetchDashboard
	}

	return newJsonResult(chains), nil
}

type StakingCalculatorRequest struct {
	Network  string  `json:"network" description:"The blockchain network to stake on, e.g. Ethereum, Solana, etc." required:"false"`
	Amount   float64 `json:"amount" description:"The amount of cryptocurrency to stake."`
	Currency string  `json:"currency" description:"The ticker symbol of the cryptocurrency to stake, e.g. ETH, SOL, etc." required:"false"`
}

type StakingCalculatorResponse struct {
	Network                 string  `json:"network"`
	Ticker                  string  `json:"ticker"`
	AmountStaked            float64 `json:"amount_staked"`
	CurrentAPYUsed          string  `json:"current_apy_used"`
	EstimatedAnnualRewards  float64 `json:"estimated_annual_rewards"`
	EstimatedMonthlyRewards float64 `json:"estimated_monthly_rewards"`
	Currency                string  `json:"currency"`
	StakingURL              string  `json:"staking_url"`
	Disclaimer              string  `json:"disclaimer"`
}

func (s *MCPServer) StakingCalculator(_ context.Context, _ *sdkmcp.CallToolRequest, input StakingCalculatorRequest) (*sdkmcp.CallToolResult, any, error) {
	chains, err := s.getChains()
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
	annualRewards := input.Amount * (apr / 100)
	monthlyRewards := annualRewards / 12

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

func (s *MCPServer) RequestIntegration(_ context.Context, _ *sdkmcp.CallToolRequest, input dashboard.PDLead) (*sdkmcp.CallToolResult, any, error) {
	if input.LeadSource == "" {
		input.LeadSource = "MCP service (source not specified)"
	}

	err := s.dashboard.CreatePDLead(input)
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

func (s *MCPServer) getChains() ([]dashboard.Chain, error) {
	if cached, found := s.cache.Get(chainsCacheKey); found {
		return cached.([]dashboard.Chain), nil
	}

	chains, err := s.dashboard.GetChains()
	if err != nil {
		return nil, err
	}

	// Clear unnecessary data
	for i := range chains {
		chains[i].LogoBlack = ""
		chains[i].LogoWhite = ""
	}

	s.cache.Set(chainsCacheKey, chains, mcp_server.DashboardCacheTtl)
	return chains, nil
}
