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
    ThemeBindings []ThemeBindings  // Application token mappings
}
```

### Component

```go
type Component struct {
    ID              string
    Name            string
    Description     string
    Category        string
    Variants        []Variant
    States          []State
    Props           []Prop
    Events          []ComponentEvent   // Custom events
    Slots           []Slot
    Uses            []string           // Component dependencies
    Constraints     *Constraints
    ThemingContract *ThemingContract   // External theming API
    Accessibility   *ComponentA11y
    LLM             *LLMContext
}
```

### ThemingContract

```go
type ThemingContract struct {
    Prefix      string       // CSS custom property prefix (e.g., "--btn")
    Description string
    Tokens      []ThemeToken
}

type ThemeToken struct {
    ID           string  // Token identifier
    CSSProperty  string  // Full CSS property name
    Semantic     string  // Semantic category for auto-mapping
    Description  string
    DefaultLight string  // Light mode default
    DefaultDark  string  // Dark mode default
}
```

### ThemeBindings

```go
type ThemeBindings struct {
    Component string          // Component ID to theme
    SpecURL   string          // Optional URL to fetch component spec
    ThemeMode string          // "light", "dark", or empty for both
    Strategy  string          // "explicit", "semantic", "inherit"
    Mappings  []TokenMapping
}

type TokenMapping struct {
    From      string  // Application token ID
    To        string  // Component token ID
    Transform string  // Optional CSS transform function
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

## Theming Contracts

### Validate Contracts

```go
// Validate a single component's contract
report := dss.ValidateContract(&component)
if !report.Passed {
    for _, err := range report.Errors {
        fmt.Printf("Error: %s: %s\n", err.Field, err.Message)
    }
}

// Validate all contracts in a design system
reports := dss.ValidateAllContracts(ds)
fmt.Printf("Errors: %d, Warnings: %d\n",
    dss.TotalErrors(reports),
    dss.TotalWarnings(reports))
```

### Generate Theme Bindings

```go
opts := dss.BindingOptions{
    Format:          dss.FormatCSS,  // or FormatTypeScript, FormatSCSS
    DefaultStrategy: "semantic",      // or "explicit", "inherit"
}

bindings, err := dss.GenerateBindings(ds, opts)
if err != nil {
    panic(err)
}

// Write to file
f, _ := os.Create("theme.css")
dss.WriteBindings(f, bindings)
```

## Diagram Generation

### Mermaid Diagrams

```go
opts := dss.DefaultMermaidOptions()
opts.Direction = "LR"  // Left to Right
opts.ShowSlots = true

// Component relationship diagram
mermaid, err := ds.GenerateMermaid(opts)

// Single component focus diagram
diagram, err := ds.GenerateMermaidComponentDiagram("button", opts)

// Token usage diagram
tokenDiagram, err := ds.GenerateMermaidTokenDiagram(opts)
```

### D2 Diagrams

```go
opts := dss.DefaultD2Options()
opts.Sketch = true  // Hand-drawn look

// Full architecture diagram
d2, err := ds.GenerateD2(opts)

// Single component diagram
diagram, err := ds.GenerateD2ComponentDiagram("button", opts)
```

## W3C Design Tokens Export

```go
opts := dss.DefaultW3CTokensOptions()
opts.IncludeDescriptions = true
opts.IncludeExtensions = true

tokens, err := ds.GenerateW3CTokens(opts)
// Output: W3C Design Tokens Community Group format JSON
```

## Documentation Generation

```go
opts := dss.DefaultDocsOptions()
opts.OutputDir = "docs/spec"
opts.IncludeMermaid = true
opts.IncludeD2 = true
opts.MkDocsCompatible = true

output, err := ds.GenerateDocs(opts)
// output.Files contains map of path -> content
for path, content := range output.Files {
    // Write each file
}
```

## NPM Package Generation

Generate publishable NPM packages with framework-specific presets and configurations.

### GeneratePackage Method

```go
opts := dss.PackageGeneratorOptions{
    OutputDir:   "./dist",
    Scope:       "@myorg",
    PackageName: "design-tokens",
    Targets: []dss.PackageTarget{
        dss.TargetCSS,
        dss.TargetTailwind,
        dss.TargetShadCN,
    },
    IncludeReadme: true,
}

err := ds.GeneratePackage(opts)
if err != nil {
    panic(err)
}
```

### PackageGeneratorOptions

```go
type PackageGeneratorOptions struct {
    // OutputDir is the directory to write the package to
    OutputDir string

    // Targets specifies which framework outputs to generate
    Targets []PackageTarget

    // Scope is the NPM scope (e.g., "@plexusone")
    // If empty, derived from meta
    Scope string

    // PackageName is the package name (default: "design-tokens")
    PackageName string

    // DryRun previews output without writing files
    DryRun bool

    // IncludeReadme generates a README.md file
    IncludeReadme bool
}
```

### PackageTarget Constants

```go
const (
    TargetCSS            PackageTarget = "css"
    TargetTailwind       PackageTarget = "tailwind"
    TargetShadCN         PackageTarget = "shadcn"
    TargetMkDocsMaterial PackageTarget = "mkdocs-material"
    TargetSCSS           PackageTarget = "scss"
    TargetJSON           PackageTarget = "json"
    TargetW3C            PackageTarget = "w3c"
)
```

### Example: Generate All Targets

```go
ds, _ := dss.LoadDesignSystem("./my-design-system/")

opts := dss.DefaultPackageOptions()
opts.OutputDir = "./dist"
opts.Targets = []dss.PackageTarget{
    dss.TargetCSS,
    dss.TargetTailwind,
    dss.TargetShadCN,
    dss.TargetMkDocsMaterial,
    dss.TargetSCSS,
    dss.TargetJSON,
    dss.TargetW3C,
}

err := ds.GeneratePackage(opts)
```

### Example: Preview Without Writing

```go
opts := dss.DefaultPackageOptions()
opts.OutputDir = "./dist"
opts.DryRun = true

err := ds.GeneratePackage(opts)
// No files written, validates package structure
```

### Generated Package Structure

```
dist/
├── package.json           # NPM manifest with exports
├── index.js               # CommonJS entry point
├── index.mjs              # ES module entry point
├── index.d.ts             # TypeScript declarations
├── README.md              # Usage documentation
├── css/
│   └── tokens.css
├── tailwind/
│   ├── preset.js
│   └── theme.json
├── shadcn/
│   └── theme.css
├── mkdocs/
│   └── theme.css
├── scss/
│   └── tokens.scss
├── json/
│   └── tokens.json
└── w3c/
    └── tokens.json
```

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
