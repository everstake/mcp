package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"mcp-server/internal/config"
	"mcp-server/internal/tools"
	"mcp-server/pkg/log"
)

func main() {
	svcCfg, err := config.LoadServiceConfig()
	if err != nil {
		log.Logger.Fatal("failed to load service config", log.E(err))
	}

	mcpCfg, err := config.LoadMCPConfig()
	if err != nil {
		log.Logger.Fatal("failed to load mcp config", log.E(err))
	}

	s := mcp.NewServer(&mcp.Implementation{
		Name:    "Everstake MCP",
		Version: "1.0.0",
	}, nil)

	handlers := map[string]mcp.ToolHandler{
		"get_api_docs": tools.HandleGetAPIDocs,
	}

	for _, tool := range mcpCfg.ToTools() {
		h, ok := handlers[tool.Name]
		if !ok {
			slog.Warn("no handler registered for tool", "tool", tool.Name)
			continue
		}
		s.AddTool(tool, h)
	}

	addr := svcCfg.Addr()
	slog.Info("starting Everstake MCP server", "addr", addr)

	handler := mcp.NewSSEHandler(func(_ *http.Request) *mcp.Server {
		return s
	})

	if err := http.ListenAndServe(addr, handler); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}
