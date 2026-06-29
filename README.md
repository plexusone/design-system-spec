# Design System Spec

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

A declarative, machine-readable specification for defining design systems as code, built for AI-native development.

## Primary Goals

DSS exists to solve two problems in AI-assisted development:

### 1. AI Agents Build Compliant UI

When an AI agent (Claude, Copilot, etc.) generates UI code, it should automatically follow the design system—using correct colors, spacing, component variants, and patterns.

**How DSS helps:** The `dss generate --llm` command produces a structured context document containing:
- Design principles and constraints
- Token values with semantic meanings
- Component specs with `allowedContexts` and `antiPatterns`
- Accessibility requirements

Include this in your AI context, and generated code follows the design system.

### 2. AI Agents Review Code for Compliance

After code is written (by humans or AI), an agent should be able to review it and report design system violations.

**How DSS helps:** The `dss validate` command checks implementations against the spec:
- Hardcoded colors → should use CSS variables
- Invalid variant values → must match component spec
- Missing accessibility attributes → required by spec
- Anti-pattern violations → flagged in component spec

JSON output (`--json`) enables programmatic integration with CI or agent workflows.

## Assessment: How Well Does DSS Achieve These Goals?

| Goal | Effectiveness | Notes |
|------|---------------|-------|
| **AI builds compliant UI** | **High** | LLM context with intent, examples, and anti-patterns significantly improves AI output quality |
| **AI reviews for compliance** | **Medium-High** | Catches token violations and accessibility issues; cannot verify visual correctness or behavioral logic |

### What DSS Catches

- Hardcoded colors, spacing values
- Invalid component variants
- Missing `alt`, `aria-label` attributes
- Anti-patterns (multiple primary buttons, nested cards)
- Non-standard token values

### What DSS Cannot Catch

- Visual correctness (does it *look* right?)
- Interaction behavior (does hover/focus work?)
- Semantic appropriateness (is this the *right* component for the use case?)
- Layout and responsive issues

For visual validation, combine DSS with visual regression testing (Chromatic, Percy).

## Installation

```bash
# CLI
go install github.com/plexusone/design-system-spec/cmd/dss@latest

# Go SDK
go get github.com/plexusone/design-system-spec
```

## Quick Start

### 1. Create a Spec

```
my-design-system/
├── meta.json
├── foundations/
│   └── colors.json
└── components/
    └── button.json
```

**meta.json:**
```json
{
  "name": "My Design System",
  "version": "1.0.0"
}
```

**components/button.json:**
```json
{
  "id": "Button",
  "name": "Button",
  "variants": [
    { "id": "default", "isDefault": true },
    { "id": "secondary" },
    { "id": "destructive" }
  ],
  "llm": {
    "intent": "Trigger user actions",
    "allowedContexts": ["forms", "dialogs", "toolbars"],
    "antiPatterns": ["Multiple primary buttons per view"]
  }
}
```

### 2. Generate LLM Context

```bash
dss generate --llm ./DESIGN_CONTEXT.md
```

Add `DESIGN_CONTEXT.md` to your AI assistant's context (Claude Projects, Copilot instructions, etc.).

### 3. Validate Implementations

```bash
# Human-readable
dss validate ./src/components

# JSON for CI/agents
dss validate --json ./src/components
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `dss info` | Display design system metadata |
| `dss generate` | Generate CSS, TypeScript types, LLM prompt |
| `dss validate <dir>` | Validate component implementations |
| `dss bind` | Generate theme bindings CSS/TypeScript/SCSS |
| `dss contract show` | Display component theming contract |
| `dss contract validate` | Validate all theming contracts |

### Generate Options

```bash
dss generate --llm ./DESIGN_CONTEXT.md  # LLM context (primary use case)
dss generate --css ./src/index.css      # Tailwind v4 @theme block
dss generate --types ./src/lib/types.ts # TypeScript interfaces
dss generate --package ./dist           # NPM package with framework presets
```

### Generating NPM Packages

Generate publishable NPM packages with framework-specific outputs:

```bash
# Generate package with all targets
dss generate --package ./dist --targets all

# Generate with specific targets
dss generate --package ./dist --targets css,tailwind,shadcn

# Custom scope and name
dss generate --package ./dist --scope @myorg --name tokens
```

**Available Targets:**

| Target | Output | Use Case |
|--------|--------|----------|
| `css` | CSS custom properties | Vanilla CSS/HTML projects |
| `tailwind` | Tailwind v4 preset and theme | Tailwind CSS v4+ projects |
| `shadcn` | ShadCN/UI theme variables | ShadCN/UI component library |
| `mkdocs-material` | MkDocs Material theme | Documentation sites |
| `scss` | SCSS variables | Sass/SCSS projects |
| `json` | Raw JSON tokens | Build tool pipelines |
| `w3c` | W3C Design Tokens format | Standards-compliant token exchange |

**Using Generated Packages:**

Tailwind v4:

```javascript
import tokens from '@myorg/design-tokens/tailwind';

export default {
  presets: [tokens],
  content: ['./src/**/*.{js,jsx,ts,tsx}'],
};
```

ShadCN:

```css
@import '@myorg/design-tokens/shadcn';
```

MkDocs Material:

```yaml
extra_css:
  - https://unpkg.com/@myorg/design-tokens/mkdocs/theme.css
```

### Theme Bindings

```bash
# Generate CSS bindings from themeBindings configuration
dss bind --output ./theme.css

# Generate TypeScript constants
dss bind --format typescript --output ./theme.ts

# Use semantic auto-mapping strategy
dss bind --strategy semantic --output ./theme.css
```

### Contract Commands

```bash
# Show a component's theming contract
dss contract show button

# Validate all theming contracts
dss contract validate
```

### Validate Output

```
✓ Passed:
  Button: validated against spec

⚠ Warnings (2):
  [no-hardcoded-colors] ./src/components/Card.tsx:45
    Hardcoded color '#333' - use CSS variable from design system
  [button-accessible-name] ./src/components/IconButton.tsx:12
    Icon-only button should have aria-label

Summary: 3 checks (1 passed, 2 warnings, 0 errors)
```

## LLM Context Structure

The generated `DESIGN_CONTEXT.md` includes:

```markdown
# My Design System

## Design Principles
- Clarity Over Complexity: ...

## Design Tokens
| Token | Value | Usage |
|-------|-------|-------|
| `primary` | `hsl(222 47% 11%)` | Primary actions and CTAs |

## Components

### Button
**Intent:** Trigger user actions
**Use in:** forms, dialogs, toolbars
**Don't use in:** inline-text, navigation

**Anti-patterns:**
- Multiple primary buttons per view
- Using button for navigation (use Link)

**Examples:**
<Button>Save</Button>
<Button variant="destructive">Delete</Button>
```

This structure is optimized for LLM comprehension and consistent code generation.

## Theming Contracts

DSS supports formal theming contracts between component libraries and consuming applications:

**Component Side (themingContract):**
```json
{
  "themingContract": {
    "prefix": "--btn",
    "tokens": [
      {
        "id": "background",
        "cssProperty": "--btn-background",
        "semantic": "primary",
        "defaultLight": "#0066CC",
        "defaultDark": "#3399FF"
      }
    ]
  }
}
```

**Application Side (themeBindings):**
```json
{
  "themeBindings": [
    {
      "component": "button",
      "mappings": [
        { "from": "brand-primary", "to": "background" }
      ]
    }
  ]
}
```

**Generate CSS:**
```bash
dss bind --output ./theme.css
```

See [Theming Specification](docs/specification/theming.md) for details.

## Figma Integration

> **Note:** Figma integration is for teams transitioning from traditional design workflows. For AI-native development, DSS specs are authored directly—Figma is not required.

For teams that still use Figma:

- **Tokens Studio** can sync Figma variables ↔ JSON tokens
- Future: `dss import figma` / `dss export figma` commands

If your workflow is AI-first (specs authored in JSON, UI generated by AI), skip Figma entirely.

## Canonical Layers

DSS defines 9 layers (most projects only need 3):

| Layer | Required | Purpose |
|-------|----------|---------|
| **Meta** | Yes | Name, version |
| **Foundations** | Yes | Tokens (colors, typography, spacing) |
| **Components** | Yes | UI elements with LLM context |
| Principles | No | Design philosophy |
| Patterns | No | Multi-component solutions |
| Templates | No | Page layouts |
| Content | No | Voice & tone |
| Accessibility | No | WCAG requirements |
| Governance | No | Versioning policies |

## Go SDK

```go
import dss "github.com/plexusone/design-system-spec/sdk/go"

ds, _ := dss.LoadDesignSystem("./my-design-system/")
ds.Validate()

// Generate for AI context
prompt, _ := ds.GenerateLLMPrompt(dss.DefaultLLMPromptOptions())

// Generate for build
css, _ := ds.GenerateCSS(dss.DefaultCSSOptions())
types, _ := ds.GenerateReactTypes(dss.DefaultReactOptions())
```

## Project Structure

```
design-system-spec/
├── cmd/dss/              # CLI tool
│   └── cmd/
│       ├── generate.go   # Generate CSS, TypeScript, LLM
│       ├── validate.go   # Validate implementations
│       ├── bind.go       # Theme bindings generation
│       ├── contract.go   # Theming contract commands
│       └── info.go
├── sdk/go/               # Go SDK
│   ├── loader.go         # Load design systems
│   ├── theming.go        # Theming contract types
│   ├── contract_validate.go  # Contract validation
│   ├── gen_bindings.go   # Theme bindings generator
│   ├── gen_css.go        # CSS generator
│   ├── gen_react.go      # TypeScript generator
│   ├── gen_llm.go        # LLM context generator
│   ├── gen_mermaid.go    # Mermaid diagram generator
│   ├── gen_d2.go         # D2 diagram generator
│   ├── gen_w3c_tokens.go # W3C Design Tokens export
│   └── gen_docs.go       # Documentation generator
├── schema/               # JSON Schemas (generated)
├── ui/                   # Web component viewer (Lit)
└── docs/                 # MkDocs documentation
```

## Roadmap

- [x] Core SDK (9 canonical layers)
- [x] CLI (`generate`, `validate`, `info`)
- [x] Code generators (CSS, TypeScript, LLM)
- [x] Documentation (MkDocs)
- [x] Theming contracts and bindings
- [x] Diagram generators (Mermaid, D2)
- [x] W3C Design Tokens export
- [x] Web component viewer (ui/)
- [ ] `dss init` scaffolding
- [ ] `dss lint` for spec completeness
- [ ] Advanced validation (color contrast, cross-references)
- [ ] Figma tokens import/export (for transitioning teams)

## License

MIT

 [go-ci-svg]: https://github.com/plexusone/design-system-spec/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/design-system-spec/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/design-system-spec/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/design-system-spec/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/design-system-spec/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/design-system-spec/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/plexusone/design-system-spec
 [goreport-url]: https://goreportcard.com/report/github.com/plexusone/design-system-spec
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/design-system-spec
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/design-system-spec
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/design-system-spec/blob/master/LICENSE
