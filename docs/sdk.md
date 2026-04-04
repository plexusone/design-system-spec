# Go SDK

The Go SDK provides programmatic access to DSS for loading, validating, and generating code from design system specifications.

## Installation

```bash
go get github.com/plexusone/design-system-spec/sdk/go
```

## Quick Start

```go
package main

import (
    "fmt"
    dss "github.com/plexusone/design-system-spec/sdk/go"
)

func main() {
    // Load from directory
    ds, err := dss.LoadDesignSystem("./my-design-system/")
    if err != nil {
        panic(err)
    }

    // Validate
    if err := ds.Validate(); err != nil {
        panic(err)
    }

    fmt.Printf("Loaded: %s v%s\n", ds.Meta.Name, ds.Meta.Version)
}
```

## Loading Design Systems

### From Directory

```go
ds, err := dss.LoadDesignSystem("./my-design-system/")
```

The loader scans for:

- `meta.json` (required)
- `principles.json`, `accessibility.json`, etc.
- `foundations/`, `components/`, `patterns/` subdirectories

### From Single File

```go
ds, err := dss.LoadDesignSystem("./design-system.json")
```

## Code Generation

### Generate CSS

```go
opts := dss.DefaultCSSOptions()
opts.Format = "tailwind4"  // or "css-vars", "scss"
opts.IncludeComments = true

css, err := ds.GenerateCSS(opts)
if err != nil {
    panic(err)
}
fmt.Println(css)
```

### CSS Formats

| Format | Output |
|--------|--------|
| `tailwind4` | Tailwind v4 `@theme` block |
| `css-vars` | Standard `:root` custom properties |
| `scss` | SCSS variables |

### Generate TypeScript Types

```go
opts := dss.DefaultReactOptions()

types, err := ds.GenerateReactTypes(opts)
if err != nil {
    panic(err)
}
fmt.Println(types)
```

Output example:

```typescript
export type ButtonVariant = 'default' | 'secondary' | 'destructive';

export interface ButtonProps {
  variant?: ButtonVariant;
  disabled?: boolean;
  children: React.ReactNode;
}
```

### Generate LLM Prompt

```go
opts := dss.DefaultLLMPromptOptions()
opts.Format = "markdown"  // or "xml"

prompt, err := ds.GenerateLLMPrompt(opts)
if err != nil {
    panic(err)
}
fmt.Println(prompt)
```

## Type Reference

### DesignSystem

```go
type DesignSystem struct {
    Meta          Meta
    Principles    []Principle
    Foundations   Foundations
    Components    []Component
    Patterns      []Pattern
    Templates     []Template
    Content       *Content
    Accessibility *Accessibility
    Governance    *Governance
}
```

### Component

```go
type Component struct {
    ID          string
    Name        string
    Description string
    Category    string
    Variants    []Variant
    States      []State
    Props       []Prop
    Slots       []Slot
    Anatomy     []AnatomyPart
    Constraints *Constraints
    LLM         *LLMContext
}
```

### LLMContext

```go
type LLMContext struct {
    Intent            string
    AllowedContexts   []string
    ForbiddenContexts []string
    ExampleUsage      []string
    AntiPatterns      []string
    SemanticMeaning   string
    RelatedElements   []string
    PriorityScore     int
}
```

## Validation

```go
if err := ds.Validate(); err != nil {
    fmt.Printf("Validation failed: %v\n", err)
}
```

Validation checks:

- Required fields present (name, version)
- Component IDs are unique
- Variant references are valid
- Color values are parseable

## JSON Schema Generation

```go
schema := dss.GenerateSchema()
data, _ := json.MarshalIndent(schema, "", "  ")
fmt.Println(string(data))
```

## Building a Custom Tool

```go
package main

import (
    "flag"
    "fmt"
    "os"

    dss "github.com/plexusone/design-system-spec/sdk/go"
)

func main() {
    dir := flag.String("dir", ".", "Design system directory")
    flag.Parse()

    ds, err := dss.LoadDesignSystem(*dir)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    // Custom processing
    for _, c := range ds.Components {
        fmt.Printf("Component: %s (%d variants)\n", c.Name, len(c.Variants))
    }
}
```
