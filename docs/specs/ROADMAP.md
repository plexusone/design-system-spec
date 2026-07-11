# Design System Spec Roadmap

## MCP Server Implementation

This section tracks the MCP (Model Context Protocol) server implementation for AI-assisted design system workflows.

### Completed

- [x] SDK Service Layer (`sdk/go/service.go`)
- [x] File Validation Logic (`sdk/go/validate_file.go`)
- [x] Skill Definition (`skills/designsystem/`)
- [x] MCP Server Entry Point (`cmd/dss-mcp/main.go`)
- [x] Internal omniskill stubs (`internal/omniskill/`)

### In Progress

#### Phase 1: Refactor CLI to Use Service ✅

Unify CLI and MCP by having both use the service layer.

| File | Status | Description |
|------|--------|-------------|
| `cmd/dss/cmd/validate.go` | ✅ Done | Use `service.ValidateDirectory()` instead of inline logic |
| `cmd/dss/cmd/info.go` | ✅ Done | Use `service.GetMeta()`, `service.ListComponents()` |
| `cmd/dss/cmd/generate.go` | ✅ Done | Use `service.GenerateLLMPrompt()` |

#### Phase 2: Skill Tests ✅

Add comprehensive tests for the skill package.

| File | Status | Description |
|------|--------|-------------|
| `skills/designsystem/skill_test.go` | ✅ Done | Test tool registration and execution |

#### Phase 3: Integration Testing

End-to-end testing of the MCP server.

| Test | Status | Description |
|------|--------|-------------|
| Build verification | ✅ Done | MCP server builds and CLI works correctly |
| Spec loading | ✅ Done | Loads minimal-system example successfully |
| MCP Inspector | Pending | `npx @anthropic/mcp-inspector dss-mcp --spec ./examples/minimal-system` |
| Claude Desktop | Pending | Test with real Claude Desktop configuration |

#### Phase 5: Embedded Filesystem Support ✅

Enable loading design systems from embedded filesystems for single-binary distribution.

| File | Status | Description |
|------|--------|-------------|
| `sdk/go/loader.go` | ✅ Done | Added `LoadDesignSystemFromFS()` for fs.FS support |
| `sdk/go/loader_test.go` | ✅ Done | Tests for embedded filesystem loading |
| `docs/mcp-server.md` | ✅ Done | Documentation for embedded MCP servers |
| `docs/sdk.md` | ✅ Done | Documentation for `LoadDesignSystemFromFS` |

#### Phase 6: Fix Tools ✅

Auto-fix design system violations for agent-ready workflows.

| File | Status | Description |
|------|--------|-------------|
| `sdk/go/fix_file.go` | ✅ Done | Fix logic for colors, spacing, accessibility |
| `sdk/go/fix_file_test.go` | ✅ Done | Comprehensive tests for fix operations |
| `skills/designsystem/tools_fix.go` | ✅ Done | MCP tools: fix_file, suggest_fixes, fix_colors, fix_spacing, fix_accessibility, fix_directory |
| `docs/mcp-server.md` | ✅ Done | Documentation for fix tools |

#### Phase 4: Documentation ✅

| File | Status | Description |
|------|--------|-------------|
| `docs/mcp-server.md` | ✅ Done | MCP server usage guide |
| `README.md` | ✅ Done | Add MCP server section |
| `mkdocs.yml` | ✅ Done | Added to navigation |

### Architecture

```
dss-mcp (MCP Server)
├── designsystem skill (15 tools)
│   ├── Spec reading: get_component, list_components, get_token, etc.
│   ├── Guidance: generate_prompt, get_variants, get_props, get_anti_patterns
│   └── Validation: validate_file, validate_directory, check_colors, check_spacing
│
└── w3pilot skill (optional, via --browser)
    └── 169 browser automation tools (auto-discovered)
```

### Usage

```bash
# Start MCP server
dss-mcp --spec ./design-system/

# With browser validation
dss-mcp --spec ./design-system/ --browser
```

### Claude Desktop Configuration

```json
{
  "mcpServers": {
    "design-system": {
      "command": "dss-mcp",
      "args": ["--spec", "/path/to/spec"]
    }
  }
}
```

---

# NPM Package Generation Roadmap

This document specifies the NPM package generation feature for design-system-spec.

## Overview

The `dss generate --package` command generates a complete, publishable NPM package containing design tokens and framework-specific presets from a design system specification.

## Motivation

Currently, design systems like PlexusOne, ProductBuildersHQ, and AIStandardsIO define tokens in JSON specs but must manually create NPM packages for distribution. This feature automates that process, enabling:

1. **Consistent token distribution** across projects
2. **Framework-specific presets** (Tailwind, ShadCN, MkDocs Material)
3. **TypeScript type safety** for token consumers
4. **Version synchronization** between spec and published package

## CLI Interface

### Basic Usage

```bash
# Generate NPM package to ./dist
dss generate --package ./dist

# Generate with specific targets
dss generate --package ./dist --targets tailwind,shadcn,mkdocs-material

# Generate with custom scope
dss generate --package ./dist --scope @myorg
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--package` | `-p` | - | Output directory for NPM package |
| `--targets` | `-t` | `css,tailwind` | Comma-separated list of targets |
| `--scope` | `-s` | From meta | NPM scope (e.g., `@plexusone`) |
| `--name` | `-n` | From meta | Package name (default: `design-tokens`) |
| `--dry-run` | | `false` | Preview without writing files |

### Available Targets

| Target | Description | Output Files |
|--------|-------------|--------------|
| `css` | CSS custom properties | `css/tokens.css` |
| `tailwind` | Tailwind CSS v4 preset | `tailwind/preset.js`, `tailwind/theme.css` |
| `shadcn` | ShadCN/UI theme variables | `shadcn/theme.css`, `shadcn/colors.json` |
| `mkdocs-material` | MkDocs Material theme | `mkdocs/extra.css`, `mkdocs/palette.yml` |
| `scss` | SCSS variables | `scss/_variables.scss` |
| `json` | Raw JSON tokens | `tokens.json` |
| `w3c` | W3C Design Tokens format | `tokens.w3c.json` |

## Output Structure

```
dist/
├── package.json              # Generated from meta
├── README.md                 # Auto-generated documentation
├── index.js                  # CommonJS entry
├── index.mjs                 # ESM entry
├── index.d.ts                # TypeScript declarations
│
├── css/
│   └── tokens.css            # CSS custom properties
│
├── tailwind/
│   ├── preset.js             # Tailwind v4 preset (theme.extend)
│   ├── preset.d.ts           # TypeScript types for preset
│   └── theme.css             # @theme block for Tailwind v4
│
├── shadcn/
│   ├── theme.css             # ShadCN theme variables
│   ├── colors.json           # Color palette for shadcn init
│   └── components.json       # ShadCN configuration
│
├── mkdocs/
│   ├── extra.css             # MkDocs Material theme overrides
│   └── palette.yml           # Color scheme configuration
│
├── scss/
│   └── _variables.scss       # SCSS variables
│
├── tokens.json               # Raw token data
└── tokens.w3c.json           # W3C Design Tokens format
```

## Generated package.json

```json
{
  "name": "@plexusone/design-tokens",
  "version": "0.1.0",
  "description": "Design tokens for PlexusOne",
  "main": "index.js",
  "module": "index.mjs",
  "types": "index.d.ts",
  "exports": {
    ".": {
      "import": "./index.mjs",
      "require": "./index.js",
      "types": "./index.d.ts"
    },
    "./css": "./css/tokens.css",
    "./tailwind": {
      "import": "./tailwind/preset.js",
      "types": "./tailwind/preset.d.ts"
    },
    "./shadcn": "./shadcn/theme.css",
    "./mkdocs": "./mkdocs/extra.css",
    "./scss": "./scss/_variables.scss"
  },
  "files": ["css", "tailwind", "shadcn", "mkdocs", "scss", "*.js", "*.mjs", "*.d.ts", "*.json"],
  "keywords": ["design-tokens", "design-system", "css", "tailwind"],
  "repository": {
    "type": "git",
    "url": "https://github.com/plexusone/plexusone-design-system"
  },
  "license": "MIT",
  "generatedBy": "design-system-spec"
}
```

## Framework-Specific Outputs

### Tailwind CSS v4

**tailwind/preset.js:**
```javascript
/** @type {import('tailwindcss').Config} */
export default {
  theme: {
    extend: {
      colors: {
        primary: 'var(--color-primary)',
        'primary-light': 'var(--color-primary-light)',
        // ...
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'ui-monospace', 'monospace'],
      },
      spacing: {
        // Uses CSS variables for Tailwind v4
      },
    },
  },
}
```

**tailwind/theme.css:**
```css
@import "tailwindcss";

@theme {
  --color-primary: #06b6d4;
  --color-primary-light: #22d3ee;
  /* ... */
}
```

### ShadCN/UI

**shadcn/theme.css:**
```css
@layer base {
  :root {
    --background: 0 0% 100%;
    --foreground: 222.2 84% 4.9%;
    --primary: 187 94% 43%;
    --primary-foreground: 210 40% 98%;
    /* HSL format for ShadCN */
  }

  .dark {
    --background: 222.2 84% 4.9%;
    --foreground: 210 40% 98%;
    /* Dark mode overrides */
  }
}
```

**shadcn/colors.json:**
```json
{
  "primary": {
    "DEFAULT": "hsl(187, 94%, 43%)",
    "foreground": "hsl(210, 40%, 98%)"
  }
}
```

### MkDocs Material

**mkdocs/extra.css:**
```css
[data-md-color-scheme="slate"] {
  --md-primary-fg-color: #06b6d4;
  --md-accent-fg-color: #8b5cf6;
  /* ... */
}
```

**mkdocs/palette.yml:**
```yaml
palette:
  - scheme: slate
    primary: custom
    accent: custom
    toggle:
      icon: material/brightness-4
      name: Switch to light mode
  - scheme: default
    primary: custom
    accent: custom
    toggle:
      icon: material/brightness-7
      name: Switch to dark mode
```

## Implementation Plan

### Phase 1: Core Package Generator

1. Add `PackageGeneratorOptions` struct to `sdk/go/`
2. Implement `GeneratePackage()` method on `DesignSystem`
3. Generate `package.json` from `Meta`
4. Generate `index.js`, `index.mjs`, `index.d.ts` exports
5. Add `--package` flag to `dss generate` command

### Phase 2: Framework Targets

1. **css** - Reuse existing `GenerateCSS()` with `css-vars` format
2. **tailwind** - New `GenerateTailwindPreset()` method
3. **shadcn** - New `GenerateShadCNTheme()` method
4. **mkdocs-material** - Reuse existing `generateMkDocsMaterialCSS()`
5. **scss** - Reuse existing `GenerateCSS()` with `scss` format

### Phase 3: TypeScript Types

1. Generate `index.d.ts` with token type definitions
2. Generate `tailwind/preset.d.ts` for Tailwind preset types
3. Export token constants as typed objects

### Phase 4: Documentation

1. Auto-generate `README.md` with usage examples
2. Document available exports and imports
3. Include version and generation timestamp

## Testing

```bash
# Generate package from minimal-system example
dss generate -d ./examples/minimal-system --package ./test-output

# Verify package structure
ls -la ./test-output/

# Validate package.json
node -e "require('./test-output/package.json')"

# Test Tailwind preset
npx tailwindcss -i ./test-output/tailwind/theme.css -o /dev/null
```

## Usage Examples

### In a Tailwind Project

```javascript
// tailwind.config.js
import preset from '@plexusone/design-tokens/tailwind'

export default {
  presets: [preset],
  content: ['./src/**/*.{js,ts,jsx,tsx}'],
}
```

### In a ShadCN Project

```bash
# Copy theme to your project
cp node_modules/@plexusone/design-tokens/shadcn/theme.css ./src/styles/
```

### In MkDocs

```yaml
# mkdocs.yml
extra_css:
  - node_modules/@plexusone/design-tokens/mkdocs/extra.css
```

### Direct CSS Import

```css
@import '@plexusone/design-tokens/css/tokens.css';

.my-button {
  background: var(--color-primary);
  color: var(--color-text);
}
```

## Success Criteria

1. Generated package passes `npm publish --dry-run`
2. Tailwind preset works with `npx tailwindcss`
3. ShadCN theme integrates without errors
4. MkDocs Material theme renders correctly
5. TypeScript types provide autocomplete for all tokens
6. Package size is reasonable (< 50KB uncompressed)

## Related

- [TASKS.md](../../TASKS.md) - Project task tracking
- [cli.md](../cli.md) - CLI reference
- [W3C Design Tokens](https://design-tokens.github.io/community-group/format/) - Standard format
