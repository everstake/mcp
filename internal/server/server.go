package server

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"mcp-server/internal/config"
	"mcp-server/internal/tools"
	"mcp-server/pkg/log"
)

type Server struct {
	serverCfg config.ServiceConfig
	mcpConfig *config.ToolsConfig
}

func New(serverCfg config.ServiceConfig, mcpConfig *config.ToolsConfig) *Server {
	return &Server{
		serverCfg: serverCfg,
		mcpConfig: mcpConfig,
	}
}

func (s *Server) Run() error {
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "Everstake MCP",
		Version: "1.0.0",
	}, nil)

	handlers := map[string]mcp.ToolHandler{
		"get_api_docs": tools.HandleGetAPIDocs,
	}

	for _, tool := range s.mcpConfig.ToTools() {
		h, ok := handlers[tool.Name]
		if !ok {
			slog.Warn("no handler registered for tool", "tool", tool.Name)
			continue
		}
		mcpServer.AddTool(tool, h)
	}

	addr := s.serverCfg.Addr()
	slog.Info("starting Everstake MCP server", "addr", addr)

	handler := mcp.NewSSEHandler(func(_ *http.Request) *mcp.Server {
		return mcpServer
	})

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Logger.Error("server error", log.E(err))
		os.Exit(1)
	}

	return nil
}
