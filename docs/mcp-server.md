# MCP Server

The `dss-mcp` command exposes design system operations as an MCP (Model Context Protocol) server, enabling AI assistants like Claude to query and validate against your design system specification.

## Overview

The MCP server provides 21 tools organized into four categories:

| Category | Tools | Purpose |
|----------|-------|---------|
| Spec Reading | 7 | Query components, tokens, patterns, and metadata |
| Guidance | 4 | Generate prompts, get variants, props, and anti-patterns |
| Validation | 4 | Validate files against the design system |
| Fix | 6 | Auto-fix design system violations |

## Installation

```bash
# Build from source
go install github.com/plexusone/design-system-spec/cmd/dss-mcp@latest

# Or build locally
go build -o dss-mcp ./cmd/dss-mcp
```

## Usage

### Basic Usage

```bash
# Start MCP server with your design system spec
dss-mcp --spec ./design-system/

# Or with a single spec file
dss-mcp --spec ./design-system.yaml
```

### With Browser Validation (w3pilot)

```bash
# Enable w3pilot browser tools for visual validation
dss-mcp --spec ./design-system/ --browser
```

This adds 169 additional browser automation tools from w3pilot for:

- Screenshot comparison
- Computed style verification
- Accessibility tree inspection
- Visual regression testing

## Claude Desktop Integration

Add to your Claude Desktop configuration:

### macOS

Edit `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "design-system": {
      "command": "dss-mcp",
      "args": ["--spec", "/path/to/your/design-system"]
    }
  }
}
```

### With Browser Tools

```json
{
  "mcpServers": {
    "design-system": {
      "command": "dss-mcp",
      "args": ["--spec", "/path/to/your/design-system", "--browser"]
    }
  }
}
```

## Available Tools

### Spec Reading Tools

#### `get_component`

Get a component definition including variants, props, states, events, and accessibility requirements.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | string | Yes | Component ID (e.g., "button", "input") |

**Example:**

```json
{
  "id": "button"
}
```

#### `list_components`

List all components with ID, name, category, and description.

**Parameters:** None

#### `get_token`

Get a design token value.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `type` | string | Yes | Token type: color, spacing, typography, elevation, borderRadius, breakpoint |
| `name` | string | Yes | Token name/ID (e.g., "primary-500") |

#### `list_tokens`

List all tokens of a given type.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `type` | string | Yes | Token type |

#### `get_pattern`

Get a pattern definition including components, layout, and variations.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | string | Yes | Pattern ID (e.g., "form-validation") |

#### `list_patterns`

List all patterns with ID, name, category, and description.

**Parameters:** None

#### `get_meta`

Get design system metadata including name, version, and description.

**Parameters:** None

### Guidance Tools

#### `generate_prompt`

Generate an LLM context prompt for the design system.

**Parameters:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `format` | string | No | "markdown" | Output format: markdown, xml |
| `include_foundations` | boolean | No | true | Include design tokens |
| `include_components` | boolean | No | true | Include component definitions |
| `include_patterns` | boolean | No | true | Include pattern recommendations |
| `include_accessibility` | boolean | No | true | Include accessibility requirements |
| `include_anti_patterns` | boolean | No | true | Include anti-patterns to avoid |

#### `get_variants`

Get all valid variants for a component.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `component_id` | string | Yes | Component ID |

#### `get_props`

Get prop definitions for a component including types, defaults, and constraints.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `component_id` | string | Yes | Component ID |

#### `get_anti_patterns`

Get anti-patterns to avoid for a component.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `component_id` | string | Yes | Component ID |

### Validation Tools

#### `validate_file`

Validate a file against the design system spec.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | Yes | Path to the file to validate |
| `rules` | array | No | Specific rules to check (default: all) |
| `include_context` | boolean | No | Include code snippets in violations |

**Available Rules:**

- `no-hardcoded-colors` - Check for hardcoded hex/rgb/hsl colors
- `use-spacing-scale` - Check for non-standard spacing values
- `img-alt-required` - Check for images without alt attributes
- `button-accessible-name` - Check for icon buttons without aria-label
- `valid-variant` - Check for unknown variant values
- `single-primary-button` - Check for multiple primary buttons
- `no-nested-cards` - Check for nested card components

#### `validate_directory`

Validate all files in a directory.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | Yes | Directory path to validate |
| `extensions` | array | No | File extensions to check (default: .tsx, .jsx, .ts, .js, .css) |
| `rules` | array | No | Specific rules to check |

#### `check_colors`

Check a file for hardcoded colors that should use design tokens.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | Yes | Path to the file to check |

#### `check_spacing`

Check a file for hardcoded spacing values that should use the spacing scale.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | Yes | Path to the file to check |

### Fix Tools

#### `fix_file`

Fix design system violations in a file. Replaces hardcoded colors with CSS variables, fixes spacing values, and adds missing accessibility attributes.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | Yes | Path to the file to fix |
| `rules` | array | No | Specific rules to fix (default: all) |
| `dry_run` | boolean | No | If true, return fixes without applying them |

#### `suggest_fixes`

Suggest fixes for design system violations without applying them. Returns the original content, suggested fixes, and what the fixed content would look like.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | Yes | Path to the file to analyze |
| `rules` | array | No | Specific rules to check (default: all) |

#### `fix_colors`

Fix hardcoded colors in a file by replacing them with CSS variables from the design system.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | Yes | Path to the file to fix |
| `dry_run` | boolean | No | If true, return fixes without applying them |

#### `fix_spacing`

Fix hardcoded spacing values in a file by replacing them with spacing tokens.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | Yes | Path to the file to fix |
| `dry_run` | boolean | No | If true, return fixes without applying them |

#### `fix_accessibility`

Fix accessibility issues in a file. Adds missing alt attributes to images and aria-label attributes to icon-only buttons.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | Yes | Path to the file to fix |
| `dry_run` | boolean | No | If true, return fixes without applying them |

#### `fix_directory`

Fix design system violations in all files in a directory.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | Yes | Path to the directory to fix |
| `rules` | array | No | Specific rules to fix (default: all) |
| `dry_run` | boolean | No | If true, return fixes without applying them |

## Example Workflows

### Component Implementation

1. Ask Claude to implement a button component
2. Claude uses `get_component` to get the button spec
3. Claude uses `get_variants` to understand available variants
4. Claude implements the component
5. Claude uses `validate_file` to check compliance

### Design Token Lookup

```
User: "What color should I use for primary actions?"

Claude: [calls get_token with type="color", name="primary-500"]
"The primary action color is #3B82F6 (primary-500)"
```

### Code Review

```
User: "Review my Button.tsx for design system compliance"

Claude: [calls validate_file with path="./src/components/Button.tsx"]
"Found 2 warnings:
- Line 15: Hardcoded color '#ff0000' - use CSS variable
- Line 23: Missing aria-label on icon button"
```

### Auto-Fix Violations

```
User: "Fix the design system violations in my components folder"

Claude: [calls fix_directory with path="./src/components", dry_run=true]
"Found 5 fixable issues across 3 files:
- Button.tsx: 2 color fixes, 1 accessibility fix
- Card.tsx: 1 spacing fix
- Input.tsx: 1 color fix

Would you like me to apply these fixes?"

User: "Yes, apply them"

Claude: [calls fix_directory with path="./src/components"]
"Applied 5 fixes:
- Replaced #ff0000 with var(--color-error)
- Added aria-label to icon button
- Replaced 15px with var(--spacing-4)
..."
```

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     AI Agents (Claude, etc.)                     │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ MCP Protocol (stdio)
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    dss-mcp MCP Server                            │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                  omniskill Runtime                          ││
│  │  ┌──────────────────────────┐  ┌──────────────────────────┐││
│  │  │ designsystem skill       │  │ w3pilot skill            │││
│  │  │ (21 native Go tools)     │  │ (169 browser tools)      │││
│  │  └──────────┬───────────────┘  └────────────┬─────────────┘││
│  └─────────────┼───────────────────────────────┼───────────────┘│
└────────────────┼───────────────────────────────┼────────────────┘
                 │                               │
                 ▼                               ▼
┌────────────────────────────────────┐  ┌───────────────────────┐
│      sdk/go/ (Service Layer)       │  │  w3pilot MCP Server   │
│  - Component operations            │  │  (subprocess)         │
│  - Token operations                │  └───────────────────────┘
│  - Validation                      │
│  - Prompt generation               │
└────────────────────────────────────┘
```

## Troubleshooting

### Server Won't Start

```bash
# Check if spec path is correct
ls -la /path/to/your/design-system

# Run with verbose output
dss-mcp --spec ./design-system/ 2>&1 | head -20
```

### Tool Errors

```bash
# Test with MCP inspector
npx @anthropic/mcp-inspector dss-mcp --spec ./design-system/
```

### Browser Tools Not Available

Ensure w3pilot is installed and in your PATH:

```bash
which w3pilot
w3pilot --version
```

## Embedding Design System Specs

For distribution, you can create a custom MCP server binary that embeds your design system spec using Go's `embed` package. This eliminates the need for external spec files.

### Example: Embedded MCP Server

```go
package main

import (
    "context"
    "embed"
    "fmt"
    "io/fs"
    "os"

    "github.com/modelcontextprotocol/go-sdk/mcp"
    dss "github.com/plexusone/design-system-spec/sdk/go"
    "github.com/plexusone/design-system-spec/skills/designsystem"
    "github.com/plexusone/design-system-spec/internal/omniskill/mcp/server"
)

//go:embed spec/*
var specFS embed.FS

func main() {
    if err := run(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func run() error {
    ctx := context.Background()

    // Load from embedded filesystem
    sub, err := fs.Sub(specFS, "spec")
    if err != nil {
        return fmt.Errorf("failed to get spec subdirectory: %w", err)
    }

    ds, err := dss.LoadDesignSystemFromFS(sub)
    if err != nil {
        return fmt.Errorf("failed to load design system: %w", err)
    }

    // Create service and skill
    service := dss.NewService(ds)
    dsSkill := designsystem.New(service)
    if err := dsSkill.Init(ctx); err != nil {
        return fmt.Errorf("failed to initialize skill: %w", err)
    }
    defer dsSkill.Close()

    // Create and run MCP server
    rt := server.New(&mcp.Implementation{
        Name:    "my-design-system",
        Version: "1.0.0",
    }, nil)
    rt.RegisterSkill(dsSkill)

    return rt.ServeStdio(ctx)
}
```

### Project Structure

```
my-design-system-mcp/
├── go.mod
├── main.go
└── spec/
    ├── meta.json
    ├── foundations/
    │   └── colors.json
    └── components/
        └── button.json
```

### Building and Distributing

```bash
# Build single binary with embedded spec
go build -o my-design-system-mcp .

# Distribute just the binary - no external files needed
./my-design-system-mcp
```

### Claude Desktop Configuration

No `--spec` flag needed since the spec is embedded:

```json
{
  "mcpServers": {
    "my-design-system": {
      "command": "/path/to/my-design-system-mcp"
    }
  }
}
```

## Related

- [CLI Reference](cli.md) - Command-line interface
- [SDK Reference](sdk.md) - Go SDK for programmatic access
- [Getting Started](getting-started.md) - Quick start guide
