package server

import (
	"context"
	"errors"
	"fmt"
	mcp_server "mcp-server"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"mcp-server/internal/config"
	servermcp "mcp-server/internal/server/mcp"
	"mcp-server/internal/server/middleware"
	"mcp-server/pkg/log"
)

type Server struct {
	mcpServer *servermcp.MCPServer
	serverCfg config.ServiceConfig
}

func New(serverCfg config.ServiceConfig, mcpServer *servermcp.MCPServer) *Server {
	return &Server{
		serverCfg: serverCfg,
		mcpServer: mcpServer,
	}
}

func (s *Server) Run(ctx context.Context) error {
	r := s.initRouter(s.mcpServer)

	addr := fmt.Sprintf(":%d", s.serverCfg.Port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Logger.Error("graceful shutdown failed", log.E(err))
		}
	}()

	log.Logger.Info("starting server", log.V("addr", addr))

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (*Server) initRouter(mcpSrv *servermcp.MCPServer) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(middleware.RateLimit(
		mcp_server.RateLimitRefillPerSecond,
		mcp_server.RateLimitMaxRequestsPerMinute,
	))

	r.Any("/", gin.WrapH(mcpSrv.Handler()))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return r
}
