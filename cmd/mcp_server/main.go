package main

import (
	"context"
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
		log.Logger.Fatal("failed to load service config", log.E(err))
	}

	mcpCfg, err := config.LoadMCPConfig()
	if err != nil {
		log.Logger.Fatal("failed to load mcp config", log.E(err))
	}

	dbClient := dashboard.NewDashboard(svcCfg.DashboardURL)

	mcps, err := mcp.New(mcpCfg, dbClient)
	if err != nil {
		log.Logger.Fatal("failed to create mcp server", log.E(err))
	}

	svc := server.New(svcCfg, mcps)

	go func() {
		err := svc.Run(ctx)
		if err != nil {
			log.Logger.Error("http server: ServeAPI", log.E(err))
		}
		stop()
	}()

	<-ctx.Done()
	<-time.After(time.Second * 10)
	log.Logger.Info("Terminated by timeout")
}
