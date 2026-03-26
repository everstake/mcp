package mcp

import (
	"context"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *MCPServer) getApiDocs(context.Context, *sdkmcp.CallToolRequest) (*sdkmcp.CallToolResult, error) {
	return newTextResult(s.mcpConfig.GetApiDocs.StaticResponse), nil
}

func (s *MCPServer) getContactInformation(context.Context, *sdkmcp.CallToolRequest) (*sdkmcp.CallToolResult, error) {
	return newTextResult(s.mcpConfig.GetContactInformation.StaticResponse), nil
}
