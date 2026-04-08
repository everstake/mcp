package mcp

import (
	"context"
	"fmt"
	"mcp-server/pkg/log"
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
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

func newJSONResult(obj any) *sdkmcp.CallToolResult {
	return &sdkmcp.CallToolResult{
		StructuredContent: map[string]any{"data": obj},
	}
}

func staticTextTool(response string) sdkmcp.ToolHandler {
	return func(context.Context, *sdkmcp.CallToolRequest) (*sdkmcp.CallToolResult, error) {
		return newTextResult(response), nil
	}
}

// sdkmcp.AddTool wrapper.
// Parses type schemas for tools with typed handlers
func addTool[In any, Out any](s *sdkmcp.Server, tool *sdkmcp.Tool, handler sdkmcp.ToolHandlerFor[In, Out]) {
	forOpts := &jsonschema.ForOptions{
		IgnoreInvalidTypes: false, // err if unknown type
	}
	inType := reflect.TypeFor[In]()
	if inType.Kind() == reflect.Pointer {
		inType = inType.Elem()
	}

	inputSchema, err := jsonschema.ForType(inType, forOpts)
	if err != nil {
		var empty In
		log.Logger.Fatal("failed to generate json schema for type", log.V("type", fmt.Sprintf("%T", empty)), log.E(err))
	}

	// outputSchema, err := jsonschema.For[Out](forOpts)
	// if err != nil {
	// 	var empty Out
	// 	log.Logger.Fatal("failed to generate json schema for type", log.V("type", fmt.Sprintf("%T", empty)), log.E(err))
	// }

	tool.InputSchema = inputSchema
	// tool.OutputSchema = outputSchema

	sdkmcp.AddTool(s, tool, handler)
}
