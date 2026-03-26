package mcp

import (
	"fmt"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// errs exposable to LLM / public
var (
	ErrFailedToFetchDashboard = fmt.Errorf("failed to fetch dashboard")
)

func newTextResult(response string) *sdkmcp.CallToolResult {
	return &sdkmcp.CallToolResult{
		Content: []sdkmcp.Content{
			&sdkmcp.TextContent{Text: response},
		},
	}
}

func newJsonResult(obj any) *sdkmcp.CallToolResult {
	return &sdkmcp.CallToolResult{
		StructuredContent: obj,
	}
}
