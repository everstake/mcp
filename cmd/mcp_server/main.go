package main

import (
	"context"
	"fmt"
	"mcp-server/internal/config"
	"mcp-server/internal/server"
	"mcp-server/internal/server/mcp"
	"mcp-server/pkg/everstake/dashboard"
	"mcp-server/pkg/log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	svcCfg, err := config.LoadServiceConfig()
	if err != nil {
		msg := fmt.Sprintf("failed to load service config: %s", err.Error())
		log.Logger.Fatal(msg, log.E(err))
	}

	mcpCfg, err := config.LoadMCPConfig()
	if err != nil {
		log.Logger.Fatal("failed to load mcp config", log.E(err))
	}

	dbClient := dashboard.NewDashboard(svcCfg.DashboardURL, svcCfg.DashboardAPIKey)

	mcps, err := mcp.New(mcpCfg, dbClient)
	if err != nil {
		log.Logger.Fatal("failed to create mcp server", log.E(err))
	}

	go func() {
		err := runTransport(ctx, svcCfg, mcps)
		if err != nil {
			log.Logger.Error("mcp server: run", log.E(err))
		}
		stop()
	}()

	<-ctx.Done()
	<-time.After(time.Second * 10)
	log.Logger.Info("Terminated by timeout")
}

func runTransport(ctx context.Context, svcCfg config.ServiceConfig, mcps *mcp.MCPServer) error {
	switch svcCfg.Transport {
	case config.TransportStdio:
		log.Logger.Info("starting MCP server with stdio transport")
		return mcps.RunStdio(ctx)
	default:
		return server.New(svcCfg, mcps).Run(ctx)
	}
}
