package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

const yt = `
# tool annotations
static_readonly_tool: &static_readonly_tool
  destructive_hint: true
  idempotent_hint: true
  readonly_hint: true
  open_world: true
 
val:
  <<: *static_readonly_tool
  var1: value1
  var2: value2
`

type val struct {
	Val testType `yaml:"val"`
}
type testType struct {
	ToolAnnotations `yaml:",inline"`
	Var1            string `yaml:"var1"`
	Var2            string `yaml:"var2"`
}

func TestUnmarshalYAML(t *testing.T) {
	var yval val
	if err := yaml.Unmarshal([]byte(yt), &yval); err != nil {
		t.Fatalf("failed to unmarshal yaml: %v", err)
	}
	cfg := yval.Val

	t.Logf("cfg = %+v", cfg)
	t.Logf("DestructiveHint=%v IdempotentHint=%v ReadOnlyHint=%v OpenWorld=%v Var1=%q Var2=%q",
		cfg.DestructiveHint, cfg.IdempotentHint, cfg.ReadOnlyHint, cfg.OpenWorld, cfg.Var1, cfg.Var2)

	if !cfg.DestructiveHint {
		t.Fatalf("expected DestructiveHint to be true, got false")
	}
	if !cfg.IdempotentHint {
		t.Fatalf("expected IdempotentHint to be true, got false")
	}
	if !cfg.ReadOnlyHint {
		t.Fatalf("expected ReadOnlyHint to be true, got false")
	}
	if !cfg.OpenWorld {
		t.Fatalf("expected OpenWorld to be true, got false")
	}
}
