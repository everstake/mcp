package config

import (
	"fmt"
	"reflect"
	"strings"

	mcp_server "mcp-server"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"gopkg.in/yaml.v3"
)

type (
	mcpConfigFile struct {
		Tools ToolsConfig `yaml:"tools"`
	}

	ToolsConfig struct {
		GetApiDocs            ToolConfig `yaml:"get_api_docs"`
		GetContactInformation ToolConfig `yaml:"get_contact_information"`
		GetUptimeMetrics      ToolConfig `yaml:"get_uptime_metrics"`
		GetChains             ToolConfig `yaml:"get_chains"`
	}

	ToolConfig struct {
		Name           string `yaml:"-"`
		Description    string `yaml:"description"`
		StaticResponse string `yaml:"static_response"`
	}
)

// unmarshal and set name to yaml tag
func (t *ToolsConfig) UnmarshalYAML(value *yaml.Node) error {
	type plain ToolsConfig
	if err := value.Decode((*plain)(t)); err != nil {
		return err
	}

	rv := reflect.ValueOf(t).Elem()
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		tag := rt.Field(i).Tag.Get("yaml")
		if tag == "" {
			continue
		}
		name := strings.SplitN(tag, ",", 2)[0]
		rv.Field(i).FieldByName("Name").SetString(name)
	}
	return nil
}

func LoadMCPConfig() (*ToolsConfig, error) {
	var raw mcpConfigFile
	if err := yaml.Unmarshal(mcp_server.MCPConfig, &raw); err != nil {
		return nil, fmt.Errorf("parse mcp config: %w", err)
	}
	return &raw.Tools, nil
}

func (tc *ToolConfig) ToTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        tc.Name,
		Description: tc.Description,
		InputSchema: map[string]any{"type": "object"},
	}
}
