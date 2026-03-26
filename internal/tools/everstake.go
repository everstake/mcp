package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func HandleGetAPIDocs(_ context.Context, _ *mcp.ServerSession, _ *mcp.CallToolParamsFor[map[string]any]) (*mcp.CallToolResultFor[any], error) {
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "https://docs.everstake.one/"},
		},
	}, nil
}
