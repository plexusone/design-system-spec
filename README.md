# Design System Spec

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

A declarative, machine-readable specification for defining complete design systems as code.

## Overview

Design System Spec (DSS) provides a canonical framework for expressing design systems (like Material Design, Carbon, Fluent) in a structured, version-controlled, LLM-optimized format. It enables:

- 📋 **Declarative definitions** - Define tokens, components, patterns as data
- 📁 **Multi-format support** - JSON/YAML for tokens, Markdown for documentation
- 🤖 **LLM optimization** - Explicit intent, contexts, and constraints for AI code generation
- ✅ **Validation** - Schema-based validation at multiple layers
- 🔷 **Go-first approach** - Go structs as source of truth, generating JSON Schema

## Installation

```bash
go get github.com/plexusone/design-system-spec/sdk/go
```

## Quick Start

```go
package main

import (
    "encoding/json"
    "fmt"

    dss "github.com/plexusone/design-system-spec/sdk/go"
)

func main() {
    // Define a minimal design system
    ds := dss.DesignSystem{
        Meta: dss.Meta{
            Name:        "My Design System",
            Version:     "1.0.0",
            Description: "A custom design system",
        },
        Foundations: dss.Foundations{
            Colors: []dss.ColorToken{
                {
                    ID:       "primary-500",
                    Value:    "#0066CC",
                    Semantic: "primary",
                    Usage:    "Primary brand color for CTAs and key actions",
                },
            },
        },
        Components: []dss.Component{
            {
                ID:          "button",
                Name:        "Button",
                Description: "Interactive button for triggering actions",
                Variants: []dss.Variant{
                    {ID: "primary", Name: "Primary", IsDefault: true},
                    {ID: "secondary", Name: "Secondary"},
                },
                LLM: &dss.LLMContext{
                    Intent:            "Trigger user actions",
                    AllowedContexts:   []string{"forms", "dialogs", "toolbars"},
                    ForbiddenContexts: []string{"inline-text"},
                    AntiPatterns:      []string{"Multiple primary buttons in one view"},
                },
            },
        },
    }

    // Validate
    if err := ds.Validate(); err != nil {
        panic(err)
    }

    // Output as JSON
    data, _ := json.MarshalIndent(ds, "", "  ")
    fmt.Println(string(data))
}
```

## Canonical Layers

DSS defines **9 canonical layers** for complete design system specification:

| Layer | Purpose | Types |
|-------|---------|-------|
| **Meta** | System metadata | `Meta`, `Maintainer` |
| **Principles** | Design philosophy | `Principle`, `PrincipleExample` |
| **Foundations** | Design tokens | `ColorToken`, `Typography`, `SpacingScale`, `ElevationToken`, `MotionSystem`, `GridSystem`, `Breakpoint` |
| **Components** | UI elements | `Component`, `Variant`, `State`, `Prop`, `Slot`, `Constraints` |
| **Patterns** | Multi-component solutions | `Pattern`, `PatternComponent`, `PatternLayout` |
| **Templates** | Page layouts | `Template`, `TemplateRegion`, `TemplateGrid` |
| **Content** | Voice & tone | `Content`, `VoiceGuidelines`, `ToneGuideline`, `Terminology`, `MicrocopyGuideline` |
| **Accessibility** | WCAG compliance | `Accessibility`, `ColorContrastRequirements`, `KeyboardRequirements` |
| **Governance** | Policies | `Governance`, `VersioningPolicy`, `DeprecationPolicy` |

## LLM Optimization

Each component, pattern, and template can include `LLMContext` for AI code generation:

```go
LLM: &dss.LLMContext{
    Intent:            "Primary action button for form submissions",
    AllowedContexts:   []string{"form-submit", "modal-confirm", "primary-cta"},
    ForbiddenContexts: []string{"destructive-action", "navigation"},
    ExampleUsage:      []string{"<Button variant='primary'>Submit</Button>"},
    AntiPatterns:      []string{"Don't use for destructive actions", "Avoid multiple primary buttons"},
    SemanticMeaning:   "Signals the main action users should take",
    RelatedElements:   []string{"button-secondary", "button-outline", "link"},
    PriorityScore:     90,
}
```

## Loading Specs

Load from a single file or directory structure:

```go
// Single file
ds, err := dss.LoadDesignSystem("design-system.json")

// Directory structure
ds, err := dss.LoadDesignSystem("./my-system/")
// Expected structure:
//   my-system/
//     meta.json
//     foundations/
//       colors.json
//       typography.json
//     components/
//       button.json
//       input.json
```

## JSON Schema

Generate JSON Schema for validation:

```go
schema := dss.GenerateSchema()
data, _ := json.MarshalIndent(schema, "", "  ")
// Use with JSON Schema validators
```

Or use the pre-generated schemas in `schema/`:

- `schema/design-system.schema.json` - Root schema
- `schema/foundations/foundations.schema.json`
- `schema/components/component.schema.json`

## Building

```bash
# Build all
make build

# Run tests
make test

# Generate schemas
make generate-schema

# Verify schemas
make verify-schema
```

## Project Structure

```
design-system-spec/
├── sdk/go/                 # Go SDK (nested module)
│   ├── meta.go             # Meta types
│   ├── foundations.go      # Token types
│   ├── components.go       # Component types
│   ├── patterns.go         # Pattern types
│   ├── templates.go        # Template types
│   ├── content.go          # Content guideline types
│   ├── accessibility.go    # Accessibility types
│   ├── governance.go       # Governance types
│   ├── llm.go              # LLM optimization types
│   ├── designsystem.go     # Root DesignSystem type
│   ├── loader.go           # File loaders
│   └── jsonschema.go       # Schema generation
├── schema/                 # Generated JSON Schemas
├── tools/generate/         # Schema generator tool
├── Makefile
└── TASKS.md                # Implementation tracking
```

## Roadmap

See [TASKS.md](TASKS.md) for detailed implementation status.

- [x] Phase 1: Core SDK types
- [x] Phase 2: Schema generation
- [ ] Phase 3: CLI tool (`dss validate`, `dss lint`, `dss generate`)
- [ ] Phase 4: Agent team for validation
- [ ] Phase 5: Examples & documentation

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
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fdesign-system-spec
 [loc-svg]: https://tokei.rs/b1/github/plexusone/design-system-spec
 [repo-url]: https://github.com/plexusone/design-system-spec
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/design-system-spec/blob/master/LICENSE
