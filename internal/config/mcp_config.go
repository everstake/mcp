package config

import (
	"fmt"

	mcp_server "mcp-server"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"gopkg.in/yaml.v3"
)

type MCPConfig struct {
	Tools ToolConfig `yaml:"tools"`
}

type ToolsConfig struct {
	GetApiDocs            ToolConfig `yaml:"get_api_docs"`
	GetContactInformation ToolConfig `yaml:"get_contact_information"`
}

type ToolConfig struct {
	Name           string `yaml:"name"`
	Description    string `yaml:"description"`
	StaticResponse string `yaml:"static_response"`
}

func LoadMCPConfig() (*ToolsConfig, error) {
	cfg := ToolsConfig{}
	if err := yaml.Unmarshal(mcp_server.MCPConfig, &cfg); err != nil {
		return nil, fmt.Errorf("parse mcp config: %w", err)
	}
	return &cfg, nil
}

func (tc *ToolConfig) ToTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        tc.Name,
		Description: tc.Description,
		InputSchema: map[string]any{"type": "object"}, // todo
	}
}

// func (tc *ToolConfig) GetInputSchema() interface{} {}
