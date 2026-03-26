<Q>
whats the best for multiline text in config? i am making a configurable mcp description and i am thinking how to make this // Add tool tool := mcp.NewTool("hello_world", mcp.WithDescription("Say hello to someone"), mcp.WithString("name", mcp.Required(), mcp.Description("Name of the person to greet"), ), ) more configurable like how to pass the description? the struct im aiming to fill is smth like struct MCPToolConfig { Description string Args []Args{ Description Required .... } } struct MCPConfig { HelloWorld MCPToolConfig LaunchDBSelfDestructionAndDropProduction MCPToolConfig } smth like this and im thinking putting this in huge const string chunks in code is literally a crime - i would hate to edit and look at that code so why not make it configurable? and what do we use for configs? json json is utter garbage for text yaml is also retarded for text md files? not config format best thing i came up with is xml i also been thinking of toml but that is wishful thinking - ive never used it and seems like there is a reason why i havent so what are my options? is there an approach for this? i bet im not the first to encounter this issue
</Q>

<A>

Yeah, stuffing long human-readable descriptions into Go string literals is exactly the kind of thing that turns a codebase into a crime scene over time 😄
You’re thinking in the right direction: this is a content/config separation problem, not just “which format is less painful.”

Let’s break it down properly.

🧠 First: what problem are you actually solving?

You have two different concerns:

Structured config (args, required flags, names, etc.)
Long-form text content (descriptions, docs, prompts)

Trying to force both into one format (JSON/YAML/XML) is what creates pain.

👉 The real solution is usually:

Split structure from text

🏆 Best practical approaches
✅ Option 1 — YAML + multiline blocks (most common in real systems)

YAML actually does handle multiline text well if you use it correctly:

hello_world:
  description: |
    This tool says hello to a user.

    It supports multiple languages and formats.
    You can extend it later with more features.
  args:
    - name: name
      required: true
      description: |
        Name of the person to greet.
        Can include unicode characters.
Why this works
| preserves formatting exactly
readable in editors
widely supported
maps cleanly to Go structs
In Go:
type Arg struct {
    Name        string `yaml:"name"`
    Required    bool   `yaml:"required"`
    Description string `yaml:"description"`
}

type MCPToolConfig struct {
    Description string `yaml:"description"`
    Args        []Arg `yaml:"args"`
}

type MCPConfig struct {
    HelloWorld MCPToolConfig `yaml:"hello_world"`
}

👉 This is probably the least painful + most standard solution.

✅ Option 2 — External Markdown files (cleanest for large text)

If descriptions get long, don’t even keep them in config.

configs/
  tools.yaml
descriptions/
  hello_world.md
  launch_db.md
YAML:
hello_world:
  description_file: descriptions/hello_world.md
  args:
    - name: name
      required: true
      description_file: descriptions/hello_world_name.md
Why this is 🔥
editors are optimized for .md
no escaping nonsense
easy to version, diff, review
non-devs can edit it

👉 This is what many serious systems do (OpenAI prompts, Terraform docs, etc.)

⚖️ Option 3 — TOML (actually solid, underrated)

TOML supports multiline strings cleanly:

[hello_world]
description = """
This tool says hello.

Supports multiple formats.
"""

[[hello_world.args]]
name = "name"
required = true
description = """
Name of the person to greet.
"""
Pros
cleaner than YAML (less footguns)
better than JSON for humans
Cons
weaker ecosystem vs YAML

👉 If you like strictness → TOML is actually a very good choice.

❌ What to avoid (you’re right here)
JSON
"description": "line1\nline2\nline3"

💀 unreadable, miserable to edit

XML

You can do:

<description>
This is text
with multiple lines
</description>

But:

verbose as hell
annoying parsing
nobody enjoys maintaining XML in 2026

👉 XML is not a serious improvement over YAML here.

💡 Advanced pattern (what experienced teams do)

👉 Combine both:

Structure in YAML
Text in Markdown

Then load like:

func loadDescription(path string) string {
    b, _ := os.ReadFile(path)
    return string(b)
}
🧩 Bonus: if you want max flexibility

Instead of:

mcp.WithDescription("...")

Do:

mcp.WithDescription(cfg.Description)

Where cfg is loaded from file.

👉 Your builder stays clean, config becomes data-driven.

🧠 Recommendation (based on your use case)

Given you're building MCP tools + likely evolving descriptions:

👉 Best setup:

YAML for structure
Markdown files for descriptions

If you want simpler:

👉 Just use YAML with | blocks

⚡ TL;DR
JSON → garbage for multiline text ❌
XML → verbose, not worth it ❌
YAML → ✅ best balance (use |)
TOML → ✅ clean alternative
Markdown files → ✅ best for large text (recommended)

If you want, I can sketch a full loader + builder pattern in Go that plugs directly into your mcp.NewTool(...) flow cleanly.
</A>