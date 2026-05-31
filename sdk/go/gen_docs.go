package dss

import (
	"fmt"
	"path/filepath"
	"strings"
)

// DocsOptions configures documentation output.
type DocsOptions struct {
	// OutputDir is the base directory for generated docs
	OutputDir string

	// ProjectName for titles and headers
	ProjectName string

	// ProjectURL for external links (e.g., GitHub repo)
	ProjectURL string

	// IncludeMermaid embeds Mermaid diagrams inline
	IncludeMermaid bool

	// IncludeD2 generates D2 diagram files
	IncludeD2 bool

	// IncludeTokens generates W3C tokens file
	IncludeTokens bool

	// MkDocsCompatible adds MkDocs-specific formatting (admonitions, etc.)
	MkDocsCompatible bool
}

// DefaultDocsOptions returns sensible defaults.
func DefaultDocsOptions() DocsOptions {
	return DocsOptions{
		OutputDir:        "docs/spec",
		ProjectName:      "",
		ProjectURL:       "",
		IncludeMermaid:   true,
		IncludeD2:        true,
		IncludeTokens:    true,
		MkDocsCompatible: true,
	}
}

// DocsOutput contains all generated documentation files.
type DocsOutput struct {
	// Files maps relative paths to content
	Files map[string]string
}

// GenerateDocs generates complete MkDocs-compatible documentation.
func (ds *DesignSystem) GenerateDocs(opts DocsOptions) (*DocsOutput, error) {
	output := &DocsOutput{
		Files: make(map[string]string),
	}

	projectName := opts.ProjectName
	if projectName == "" {
		projectName = ds.Meta.Name
	}

	// Generate index page
	index, err := ds.generateDocsIndex(opts, projectName)
	if err != nil {
		return nil, fmt.Errorf("generating index: %w", err)
	}
	output.Files["index.md"] = index

	// Generate component pages
	if len(ds.Components) > 0 {
		// Components index
		compIndex, err := ds.generateComponentsIndex(opts)
		if err != nil {
			return nil, fmt.Errorf("generating components index: %w", err)
		}
		output.Files["components/index.md"] = compIndex

		// Individual component pages
		for _, c := range ds.Components {
			page, err := ds.generateComponentPage(c, opts)
			if err != nil {
				return nil, fmt.Errorf("generating component page for %s: %w", c.ID, err)
			}
			output.Files[filepath.Join("components", c.ID+".md")] = page
		}
	}

	// Generate tokens documentation
	if !ds.Foundations.isEmpty() {
		tokensIndex, err := ds.generateTokensIndex(opts)
		if err != nil {
			return nil, fmt.Errorf("generating tokens index: %w", err)
		}
		output.Files["tokens/index.md"] = tokensIndex

		// W3C tokens JSON
		if opts.IncludeTokens {
			w3cOpts := DefaultW3CTokensOptions()
			tokens, err := ds.GenerateW3CTokens(w3cOpts)
			if err != nil {
				return nil, fmt.Errorf("generating W3C tokens: %w", err)
			}
			output.Files["tokens/tokens.json"] = tokens
		}
	}

	// Generate diagrams
	if opts.IncludeMermaid && len(ds.Components) > 0 {
		mermaidOpts := DefaultMermaidOptions()
		mermaid, err := ds.GenerateMermaid(mermaidOpts)
		if err == nil { // Don't fail if no relationships
			output.Files["diagrams/component-graph.mmd"] = mermaid
		}

		// Token usage diagram
		tokenDiagram, err := ds.GenerateMermaidTokenDiagram(mermaidOpts)
		if err == nil {
			output.Files["diagrams/token-usage.mmd"] = tokenDiagram
		}
	}

	if opts.IncludeD2 && len(ds.Components) > 0 {
		d2Opts := DefaultD2Options()
		d2, err := ds.GenerateD2(d2Opts)
		if err == nil {
			output.Files["diagrams/architecture.d2"] = d2
		}
	}

	// Generate patterns documentation
	if len(ds.Patterns) > 0 {
		patternsIndex, err := ds.generatePatternsIndex(opts)
		if err != nil {
			return nil, fmt.Errorf("generating patterns index: %w", err)
		}
		output.Files["patterns/index.md"] = patternsIndex
	}

	return output, nil
}

func (ds *DesignSystem) generateDocsIndex(opts DocsOptions, projectName string) (string, error) {
	var b strings.Builder

	b.WriteString("# ")
	b.WriteString(projectName)
	b.WriteString(" Design System\n\n")

	if ds.Meta.Description != "" {
		b.WriteString(ds.Meta.Description)
		b.WriteString("\n\n")
	}

	// Version badge
	if ds.Meta.Version != "" {
		b.WriteString(fmt.Sprintf("**Version:** %s\n\n", ds.Meta.Version))
	}

	// Overview diagram
	if opts.IncludeMermaid && len(ds.Components) > 0 {
		b.WriteString("## Architecture Overview\n\n")
		mermaidOpts := DefaultMermaidOptions()
		mermaid, err := ds.GenerateMermaid(mermaidOpts)
		if err == nil {
			b.WriteString("```mermaid\n")
			b.WriteString(mermaid)
			b.WriteString("```\n\n")
		}
	}

	// Quick stats
	b.WriteString("## Quick Stats\n\n")
	b.WriteString("| Category | Count |\n")
	b.WriteString("|----------|-------|\n")
	b.WriteString(fmt.Sprintf("| Components | %d |\n", len(ds.Components)))
	b.WriteString(fmt.Sprintf("| Patterns | %d |\n", len(ds.Patterns)))
	b.WriteString(fmt.Sprintf("| Color Tokens | %d |\n", len(ds.Foundations.Colors)))
	if ds.Foundations.Spacing != nil {
		b.WriteString(fmt.Sprintf("| Spacing Tokens | %d |\n", len(ds.Foundations.Spacing.Scale)))
	}
	b.WriteString("\n")

	// Navigation
	b.WriteString("## Documentation\n\n")
	if len(ds.Components) > 0 {
		b.WriteString("- [Components](components/index.md) - UI component specifications\n")
	}
	if !ds.Foundations.isEmpty() {
		b.WriteString("- [Tokens](tokens/index.md) - Design tokens and theming\n")
	}
	if len(ds.Patterns) > 0 {
		b.WriteString("- [Patterns](patterns/index.md) - Common usage patterns\n")
	}
	b.WriteString("\n")

	// External links
	if opts.ProjectURL != "" {
		b.WriteString("## Links\n\n")
		b.WriteString(fmt.Sprintf("- [Source Code](%s)\n", opts.ProjectURL))
	}

	return b.String(), nil
}

func (ds *DesignSystem) generateComponentsIndex(opts DocsOptions) (string, error) {
	var b strings.Builder

	b.WriteString("# Components\n\n")

	// Component graph
	if opts.IncludeMermaid {
		b.WriteString("## Component Relationships\n\n")
		mermaidOpts := DefaultMermaidOptions()
		mermaid, err := ds.GenerateMermaid(mermaidOpts)
		if err == nil {
			b.WriteString("```mermaid\n")
			b.WriteString(mermaid)
			b.WriteString("```\n\n")
		}
	}

	// Group by category
	categories := make(map[string][]Component)
	for _, c := range ds.Components {
		cat := c.Category
		if cat == "" {
			cat = "Other"
		}
		categories[cat] = append(categories[cat], c)
	}

	b.WriteString("## Component List\n\n")

	for cat, components := range categories {
		b.WriteString(fmt.Sprintf("### %s\n\n", formatD2Label(cat)))
		b.WriteString("| Component | Description |\n")
		b.WriteString("|-----------|-------------|\n")
		for _, c := range components {
			desc := c.Description
			if len(desc) > 60 {
				desc = desc[:57] + "..."
			}
			b.WriteString(fmt.Sprintf("| [%s](%s.md) | %s |\n", c.Name, c.ID, desc))
		}
		b.WriteString("\n")
	}

	return b.String(), nil
}

func (ds *DesignSystem) generateComponentPage(c Component, opts DocsOptions) (string, error) {
	var b strings.Builder

	// Title
	b.WriteString("# ")
	b.WriteString(c.Name)
	b.WriteString("\n\n")

	// Description
	if c.Description != "" {
		b.WriteString(c.Description)
		b.WriteString("\n\n")
	}

	// Category badge
	if c.Category != "" {
		b.WriteString(fmt.Sprintf("**Category:** %s\n\n", c.Category))
	}

	// Component diagram
	if opts.IncludeMermaid {
		mermaidOpts := DefaultMermaidOptions()
		diagram, err := ds.GenerateMermaidComponentDiagram(c.ID, mermaidOpts)
		if err == nil {
			b.WriteString("## Component Diagram\n\n")
			b.WriteString("```mermaid\n")
			b.WriteString(diagram)
			b.WriteString("```\n\n")
		}
	}

	// Props
	if len(c.Props) > 0 {
		b.WriteString("## Props\n\n")
		b.WriteString("| Name | Type | Required | Default | Description |\n")
		b.WriteString("|------|------|----------|---------|-------------|\n")
		for _, p := range c.Props {
			required := ""
			if p.Required {
				required = "Yes"
			}
			defaultVal := ""
			if p.Default != nil {
				defaultVal = fmt.Sprintf("`%v`", p.Default)
			}
			desc := p.Description
			if len(desc) > 40 {
				desc = desc[:37] + "..."
			}
			b.WriteString(fmt.Sprintf("| `%s` | `%s` | %s | %s | %s |\n",
				p.Name, p.Type, required, defaultVal, desc))
		}
		b.WriteString("\n")

		// Detailed prop descriptions
		for _, p := range c.Props {
			if p.Description != "" || len(p.EnumValues) > 0 || p.Constraints != nil {
				b.WriteString(fmt.Sprintf("### `%s`\n\n", p.Name))
				if p.Description != "" {
					b.WriteString(p.Description)
					b.WriteString("\n\n")
				}
				if len(p.EnumValues) > 0 {
					b.WriteString("**Allowed values:**\n\n")
					for _, v := range p.EnumValues {
						b.WriteString(fmt.Sprintf("- `%s`\n", v))
					}
					b.WriteString("\n")
				}
				if p.Constraints != nil {
					b.WriteString("**Constraints:**\n\n")
					if p.Constraints.Min != nil {
						b.WriteString(fmt.Sprintf("- Min: `%v`\n", *p.Constraints.Min))
					}
					if p.Constraints.Max != nil {
						b.WriteString(fmt.Sprintf("- Max: `%v`\n", *p.Constraints.Max))
					}
					if p.Constraints.Pattern != "" {
						b.WriteString(fmt.Sprintf("- Pattern: `%s`\n", p.Constraints.Pattern))
					}
					b.WriteString("\n")
				}
			}
		}
	}

	// Events
	if len(c.Events) > 0 {
		b.WriteString("## Events\n\n")
		b.WriteString("| Name | Bubbles | Description |\n")
		b.WriteString("|------|---------|-------------|\n")
		for _, e := range c.Events {
			bubbles := ""
			if e.Bubbles {
				bubbles = "Yes"
			}
			desc := e.Description
			if len(desc) > 50 {
				desc = desc[:47] + "..."
			}
			b.WriteString(fmt.Sprintf("| `%s` | %s | %s |\n", e.Name, bubbles, desc))
		}
		b.WriteString("\n")

		// Event details
		for _, e := range c.Events {
			if len(e.Detail) > 0 {
				b.WriteString(fmt.Sprintf("### `%s` Detail\n\n", e.Name))
				b.WriteString("| Field | Type | Description |\n")
				b.WriteString("|-------|------|-------------|\n")
				for _, d := range e.Detail {
					b.WriteString(fmt.Sprintf("| `%s` | `%s` | %s |\n", d.Name, d.Type, d.Description))
				}
				b.WriteString("\n")
			}
		}
	}

	// Slots
	if len(c.Slots) > 0 {
		b.WriteString("## Slots\n\n")
		b.WriteString("| Name | Required | Description |\n")
		b.WriteString("|------|----------|-------------|\n")
		for _, s := range c.Slots {
			required := ""
			if s.Required {
				required = "Yes"
			}
			b.WriteString(fmt.Sprintf("| `%s` | %s | %s |\n", s.Name, required, s.Description))
		}
		b.WriteString("\n")

		// Allowed components
		for _, s := range c.Slots {
			if len(s.AllowedComponents) > 0 {
				b.WriteString(fmt.Sprintf("**Slot `%s` accepts:**\n\n", s.Name))
				for _, ac := range s.AllowedComponents {
					b.WriteString(fmt.Sprintf("- [%s](%s.md)\n", ac, ac))
				}
				b.WriteString("\n")
			}
		}
	}

	// States
	if len(c.States) > 0 {
		b.WriteString("## States\n\n")
		b.WriteString("| State | Description |\n")
		b.WriteString("|-------|-------------|\n")
		for _, s := range c.States {
			b.WriteString(fmt.Sprintf("| `%s` | %s |\n", s.ID, s.Description))
		}
		b.WriteString("\n")
	}

	// Variants
	if len(c.Variants) > 0 {
		b.WriteString("## Variants\n\n")
		b.WriteString("| Variant | Default | Description |\n")
		b.WriteString("|---------|---------|-------------|\n")
		for _, v := range c.Variants {
			isDefault := ""
			if v.IsDefault {
				isDefault = "Yes"
			}
			b.WriteString(fmt.Sprintf("| `%s` | %s | %s |\n", v.ID, isDefault, v.Description))
		}
		b.WriteString("\n")
	}

	// Dependencies
	if len(c.Uses) > 0 {
		b.WriteString("## Dependencies\n\n")
		b.WriteString("This component uses:\n\n")
		for _, dep := range c.Uses {
			b.WriteString(fmt.Sprintf("- [%s](%s.md)\n", dep, dep))
		}
		b.WriteString("\n")
	}

	// Tokens used
	if len(c.TokensUsed) > 0 {
		b.WriteString("## Tokens Used\n\n")
		b.WriteString("| Token |\n")
		b.WriteString("|-------|\n")
		for _, t := range c.TokensUsed {
			b.WriteString(fmt.Sprintf("| `%s` |\n", t))
		}
		b.WriteString("\n")
	}

	// Accessibility
	if c.Accessibility != nil {
		b.WriteString("## Accessibility\n\n")
		if c.Accessibility.Role != "" {
			b.WriteString(fmt.Sprintf("**ARIA Role:** `%s`\n\n", c.Accessibility.Role))
		}
		if len(c.Accessibility.RequiredAttributes) > 0 {
			b.WriteString("**Required ARIA Attributes:**\n\n")
			for _, attr := range c.Accessibility.RequiredAttributes {
				b.WriteString(fmt.Sprintf("- `%s`\n", attr))
			}
			b.WriteString("\n")
		}
		if len(c.Accessibility.KeyboardSupport) > 0 {
			b.WriteString("**Keyboard Support:**\n\n")
			b.WriteString("| Key | Action |\n")
			b.WriteString("|-----|--------|\n")
			for _, k := range c.Accessibility.KeyboardSupport {
				b.WriteString(fmt.Sprintf("| `%s` | %s |\n", k.Key, k.Action))
			}
			b.WriteString("\n")
		}
	}

	// External links
	if c.DocumentationURL != "" || c.FigmaURL != "" {
		b.WriteString("## Links\n\n")
		if c.DocumentationURL != "" {
			b.WriteString(fmt.Sprintf("- [Documentation](%s)\n", c.DocumentationURL))
		}
		if c.FigmaURL != "" {
			b.WriteString(fmt.Sprintf("- [Figma](%s)\n", c.FigmaURL))
		}
		b.WriteString("\n")
	}

	return b.String(), nil
}

func (ds *DesignSystem) generateTokensIndex(opts DocsOptions) (string, error) {
	var b strings.Builder

	b.WriteString("# Design Tokens\n\n")

	b.WriteString("Design tokens are the visual design atoms of the design system.\n\n")

	if opts.IncludeTokens {
		b.WriteString("**Download:** [tokens.json](tokens.json) (W3C Design Tokens format)\n\n")
	}

	f := ds.Foundations

	// Colors
	if len(f.Colors) > 0 {
		b.WriteString("## Colors\n\n")
		b.WriteString("| Token | Value | Usage |\n")
		b.WriteString("|-------|-------|-------|\n")
		for _, c := range f.Colors {
			usage := c.Usage
			if len(usage) > 40 {
				usage = usage[:37] + "..."
			}
			b.WriteString(fmt.Sprintf("| `%s` | `%s` | %s |\n", c.ID, c.Value, usage))
		}
		b.WriteString("\n")
	}

	// Spacing
	if f.Spacing != nil && len(f.Spacing.Scale) > 0 {
		b.WriteString("## Spacing\n\n")
		b.WriteString("| Token | Value |\n")
		b.WriteString("|-------|-------|\n")
		for _, s := range f.Spacing.Scale {
			b.WriteString(fmt.Sprintf("| `%s` | `%s` |\n", s.ID, s.Value))
		}
		b.WriteString("\n")
	}

	// Typography
	if f.Typography != nil {
		if len(f.Typography.FontFamilies) > 0 {
			b.WriteString("## Font Families\n\n")
			b.WriteString("| Token | Value |\n")
			b.WriteString("|-------|-------|\n")
			for _, ff := range f.Typography.FontFamilies {
				value := ff.Stack
				if value == "" {
					value = ff.Value
				}
				b.WriteString(fmt.Sprintf("| `%s` | `%s` |\n", ff.ID, value))
			}
			b.WriteString("\n")
		}

		if len(f.Typography.FontSizes) > 0 {
			b.WriteString("## Font Sizes\n\n")
			b.WriteString("| Token | Value |\n")
			b.WriteString("|-------|-------|\n")
			for _, fs := range f.Typography.FontSizes {
				b.WriteString(fmt.Sprintf("| `%s` | `%s` |\n", fs.ID, fs.Value))
			}
			b.WriteString("\n")
		}
	}

	// Border Radius
	if len(f.BorderRadius) > 0 {
		b.WriteString("## Border Radius\n\n")
		b.WriteString("| Token | Value |\n")
		b.WriteString("|-------|-------|\n")
		for _, br := range f.BorderRadius {
			b.WriteString(fmt.Sprintf("| `%s` | `%s` |\n", br.ID, br.Value))
		}
		b.WriteString("\n")
	}

	// Shadows
	if len(f.Elevation) > 0 {
		b.WriteString("## Shadows\n\n")
		b.WriteString("| Token | Value | Usage |\n")
		b.WriteString("|-------|-------|-------|\n")
		for _, e := range f.Elevation {
			usage := e.Usage
			if len(usage) > 40 {
				usage = usage[:37] + "..."
			}
			b.WriteString(fmt.Sprintf("| `%s` | `%s` | %s |\n", e.ID, e.Value, usage))
		}
		b.WriteString("\n")
	}

	// Token usage diagram
	if opts.IncludeMermaid && len(ds.Components) > 0 {
		b.WriteString("## Token Usage\n\n")
		mermaidOpts := DefaultMermaidOptions()
		diagram, err := ds.GenerateMermaidTokenDiagram(mermaidOpts)
		if err == nil {
			b.WriteString("```mermaid\n")
			b.WriteString(diagram)
			b.WriteString("```\n\n")
		}
	}

	return b.String(), nil
}

func (ds *DesignSystem) generatePatternsIndex(opts DocsOptions) (string, error) {
	var b strings.Builder

	b.WriteString("# Patterns\n\n")

	b.WriteString("Patterns are common UI solutions built from multiple components.\n\n")

	for _, p := range ds.Patterns {
		b.WriteString(fmt.Sprintf("## %s\n\n", p.Name))

		if p.Description != "" {
			b.WriteString(p.Description)
			b.WriteString("\n\n")
		}

		if len(p.Components) > 0 {
			b.WriteString("**Components used:**\n\n")
			for _, pc := range p.Components {
				role := pc.Role
				if role == "" {
					role = "component"
				}
				b.WriteString(fmt.Sprintf("- [%s](../components/%s.md) - %s\n", pc.ComponentID, pc.ComponentID, role))
			}
			b.WriteString("\n")
		}

		if p.Layout != nil {
			b.WriteString("**Layout:**\n\n")
			b.WriteString(fmt.Sprintf("- Type: `%s`\n", p.Layout.Type))
			if p.Layout.Direction != "" {
				b.WriteString(fmt.Sprintf("- Direction: `%s`\n", p.Layout.Direction))
			}
			if p.Layout.Spacing != "" {
				b.WriteString(fmt.Sprintf("- Spacing: `%s`\n", p.Layout.Spacing))
			}
			b.WriteString("\n")
		}

		if p.Behavior != "" {
			b.WriteString("**Behavior:**\n\n")
			b.WriteString(p.Behavior)
			b.WriteString("\n\n")
		}
	}

	return b.String(), nil
}
