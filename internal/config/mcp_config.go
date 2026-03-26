package config

import (
	"fmt"

	mcp_server "mcp-server"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"gopkg.in/yaml.v3"
)

type ToolsConfig struct {
	Tools []ToolConfig `yaml:"tools"`
}

type ToolConfig struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	Args        []ArgConfig `yaml:"args"`
}

type ArgConfig struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Required    bool   `yaml:"required"`
	Description string `yaml:"description"`
}

func LoadMCPConfig() (*ToolsConfig, error) {
	cfg := ToolsConfig{}
	if err := yaml.Unmarshal(mcp_server.MCPConfig, &cfg); err != nil {
		return nil, fmt.Errorf("parse mcp config: %w", err)
	}
	return &cfg, nil
}

func (cfg *ToolsConfig) ToTools() []*mcp.Tool {
	tools := make([]*mcp.Tool, 0, len(cfg.Tools))
	for i := range cfg.Tools {
		tools = append(tools, toTool(&cfg.Tools[i]))
	}
	return tools
}

func toTool(tc *ToolConfig) *mcp.Tool {
	return &mcp.Tool{
		Name:        tc.Name,
		Description: tc.Description,
		InputSchema: buildInputSchema(tc.Args),
	}
}

func buildInputSchema(args []ArgConfig) *jsonschema.Schema {
	properties := make(map[string]*jsonschema.Schema)
	var required []string

	for _, arg := range args {
		prop := &jsonschema.Schema{Description: arg.Description}
		switch arg.Type {
		case "number":
			prop.Type = "number"
		case "boolean":
			prop.Type = "boolean"
		default: // "string"
			prop.Type = "string"
		}
		properties[arg.Name] = prop
		if arg.Required {
			required = append(required, arg.Name)
		}
	}

	schema := &jsonschema.Schema{Type: "object"}
	if len(properties) > 0 {
		schema.Properties = properties
	}
	if len(required) > 0 {
		schema.Required = required
	}
	return schema
}
