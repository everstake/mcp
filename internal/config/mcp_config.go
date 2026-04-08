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
		GetCompanyProfile     ToolConfig `yaml:"get_company_profile"`
		GetDeveloperDocs      ToolConfig `yaml:"get_developer_docs"`
		GetContactInformation ToolConfig `yaml:"get_contact_information"`
		GetUptimeMetrics      ToolConfig `yaml:"get_uptime_metrics"`
		GetChains             ToolConfig `yaml:"get_chains"`
		GetProducts           ToolConfig `yaml:"get_products"`
		GetSolutions          ToolConfig `yaml:"get_solutions"`
		GetSecurityProfile    ToolConfig `yaml:"get_security_profile"`
		GetIntegrations       ToolConfig `yaml:"get_integrations"`
		StakingCalculator     ToolConfig `yaml:"staking_calculator"`
		RequestIntegration    ToolConfig `yaml:"request_integration"`
	}

	ToolAnnotations struct {
		DestructiveHint bool `yaml:"destructive_hint"`
		IdempotentHint  bool `yaml:"idempotent_hint"`
		ReadOnlyHint    bool `yaml:"readonly_hint"`
		OpenWorld       bool `yaml:"open_world"`
	}

	ToolConfig struct {
		InputSchema     map[string]interface{} `yaml:"input_schema"`
		Name            string                 `yaml:"-"`
		Description     string                 `yaml:"description"`
		StaticResponse  string                 `yaml:"static_response"`
		ToolAnnotations `yaml:",inline"`
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
	schema := map[string]any{"type": "object"}
	if tc.InputSchema != nil {
		schema = tc.InputSchema
	}

	return &mcp.Tool{
		Name:        tc.Name,
		Description: tc.Description,
		InputSchema: schema,
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: &tc.DestructiveHint,
			IdempotentHint:  tc.IdempotentHint,
			OpenWorldHint:   &tc.OpenWorld,
			ReadOnlyHint:    tc.ReadOnlyHint,
		},
	}
}
