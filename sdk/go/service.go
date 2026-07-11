package dss

import (
	"context"
	"fmt"
)

// Service provides design system operations for CLI, MCP, and SDK consumers.
// It wraps a loaded DesignSystem and provides a unified API for querying
// and validating against the specification.
type Service struct {
	ds *DesignSystem
}

// NewService creates a service from a loaded design system.
func NewService(ds *DesignSystem) *Service {
	return &Service{ds: ds}
}

// DesignSystem returns the underlying design system.
func (s *Service) DesignSystem() *DesignSystem {
	return s.ds
}

// ComponentSummary provides a brief overview of a component.
type ComponentSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
}

// TokenSummary provides a brief overview of a design token.
type TokenSummary struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// PatternSummary provides a brief overview of a pattern.
type PatternSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
}

// PromptOptions configures LLM prompt generation.
type PromptOptions struct {
	// Format: "markdown" (default), "xml"
	Format string `json:"format,omitempty"`

	// ComponentIDs limits output to specific components (empty = all)
	ComponentIDs []string `json:"componentIds,omitempty"`

	// IncludeFoundations includes design tokens
	IncludeFoundations bool `json:"includeFoundations,omitempty"`

	// IncludeComponents includes component definitions
	IncludeComponents bool `json:"includeComponents,omitempty"`

	// IncludePatterns includes pattern definitions
	IncludePatterns bool `json:"includePatterns,omitempty"`

	// IncludeAccessibility includes a11y requirements
	IncludeAccessibility bool `json:"includeAccessibility,omitempty"`

	// IncludeAntiPatterns includes anti-pattern warnings
	IncludeAntiPatterns bool `json:"includeAntiPatterns,omitempty"`

	// MaxExamples limits code examples per component
	MaxExamples int `json:"maxExamples,omitempty"`
}

// DefaultPromptOptions returns comprehensive defaults.
func DefaultPromptOptions() *PromptOptions {
	return &PromptOptions{
		Format:               "markdown",
		IncludeFoundations:   true,
		IncludeComponents:    true,
		IncludePatterns:      true,
		IncludeAccessibility: true,
		IncludeAntiPatterns:  true,
		MaxExamples:          3,
	}
}

// GetComponent returns a component by ID.
func (s *Service) GetComponent(_ context.Context, id string) (*Component, error) {
	for i := range s.ds.Components {
		if s.ds.Components[i].ID == id {
			return &s.ds.Components[i], nil
		}
	}
	return nil, fmt.Errorf("component not found: %s", id)
}

// ListComponents returns summaries of all components.
func (s *Service) ListComponents(_ context.Context) []ComponentSummary {
	result := make([]ComponentSummary, len(s.ds.Components))
	for i, c := range s.ds.Components {
		result[i] = ComponentSummary{
			ID:          c.ID,
			Name:        c.Name,
			Category:    c.Category,
			Description: c.Description,
		}
	}
	return result
}

// GetComponentVariants returns the variants for a component.
func (s *Service) GetComponentVariants(_ context.Context, componentID string) ([]Variant, error) {
	for _, c := range s.ds.Components {
		if c.ID == componentID {
			return c.Variants, nil
		}
	}
	return nil, fmt.Errorf("component not found: %s", componentID)
}

// GetComponentProps returns the props for a component.
func (s *Service) GetComponentProps(_ context.Context, componentID string) ([]Prop, error) {
	for _, c := range s.ds.Components {
		if c.ID == componentID {
			return c.Props, nil
		}
	}
	return nil, fmt.Errorf("component not found: %s", componentID)
}

// GetComponentAntiPatterns returns anti-patterns for a component.
func (s *Service) GetComponentAntiPatterns(_ context.Context, componentID string) ([]string, error) {
	for _, c := range s.ds.Components {
		if c.ID == componentID {
			if c.LLM != nil {
				return c.LLM.AntiPatterns, nil
			}
			return nil, nil
		}
	}
	return nil, fmt.Errorf("component not found: %s", componentID)
}

// GetToken returns a design token by type and name.
func (s *Service) GetToken(_ context.Context, tokenType, name string) (any, error) {
	f := s.ds.Foundations

	switch tokenType {
	case "color":
		for _, c := range f.Colors {
			if c.ID == name {
				return c, nil
			}
		}
	case "spacing":
		if f.Spacing != nil {
			for _, sp := range f.Spacing.Scale {
				if sp.ID == name {
					return sp, nil
				}
			}
		}
	case "typography":
		if f.Typography != nil {
			for _, ts := range f.Typography.TypeScale {
				if ts.ID == name {
					return ts, nil
				}
			}
		}
	case "elevation":
		for _, e := range f.Elevation {
			if e.ID == name {
				return e, nil
			}
		}
	case "borderRadius":
		for _, br := range f.BorderRadius {
			if br.ID == name {
				return br, nil
			}
		}
	case "breakpoint":
		for _, bp := range f.Breakpoints {
			if bp.ID == name {
				return bp, nil
			}
		}
	default:
		return nil, fmt.Errorf("unknown token type: %s", tokenType)
	}

	return nil, fmt.Errorf("token not found: %s/%s", tokenType, name)
}

// ListTokens returns all tokens of a given type.
func (s *Service) ListTokens(_ context.Context, tokenType string) ([]TokenSummary, error) {
	f := s.ds.Foundations
	var result []TokenSummary

	switch tokenType {
	case "color":
		for _, c := range f.Colors {
			result = append(result, TokenSummary{
				ID:    c.ID,
				Type:  "color",
				Value: c.Value,
			})
		}
	case "spacing":
		if f.Spacing != nil {
			for _, sp := range f.Spacing.Scale {
				result = append(result, TokenSummary{
					ID:    sp.ID,
					Type:  "spacing",
					Value: sp.Value,
				})
			}
		}
	case "typography":
		if f.Typography != nil {
			for _, ts := range f.Typography.TypeScale {
				result = append(result, TokenSummary{
					ID:    ts.ID,
					Type:  "typography",
					Value: fmt.Sprintf("%s %s", ts.FontSize, ts.FontWeight),
				})
			}
		}
	case "elevation":
		for _, e := range f.Elevation {
			result = append(result, TokenSummary{
				ID:    e.ID,
				Type:  "elevation",
				Value: e.Value,
			})
		}
	case "borderRadius":
		for _, br := range f.BorderRadius {
			result = append(result, TokenSummary{
				ID:    br.ID,
				Type:  "borderRadius",
				Value: br.Value,
			})
		}
	case "breakpoint":
		for _, bp := range f.Breakpoints {
			result = append(result, TokenSummary{
				ID:    bp.ID,
				Type:  "breakpoint",
				Value: bp.MinWidth,
			})
		}
	default:
		return nil, fmt.Errorf("unknown token type: %s", tokenType)
	}

	return result, nil
}

// GetPattern returns a pattern by ID.
func (s *Service) GetPattern(_ context.Context, id string) (*Pattern, error) {
	for i := range s.ds.Patterns {
		if s.ds.Patterns[i].ID == id {
			return &s.ds.Patterns[i], nil
		}
	}
	return nil, fmt.Errorf("pattern not found: %s", id)
}

// ListPatterns returns summaries of all patterns.
func (s *Service) ListPatterns(_ context.Context) []PatternSummary {
	result := make([]PatternSummary, len(s.ds.Patterns))
	for i, p := range s.ds.Patterns {
		result[i] = PatternSummary{
			ID:          p.ID,
			Name:        p.Name,
			Category:    p.Category,
			Description: p.Description,
		}
	}
	return result
}

// GetMeta returns the design system metadata.
func (s *Service) GetMeta(_ context.Context) Meta {
	return s.ds.Meta
}

// GenerateLLMPrompt generates a complete LLM context prompt.
func (s *Service) GenerateLLMPrompt(_ context.Context, opts *PromptOptions) (string, error) {
	if opts == nil {
		opts = DefaultPromptOptions()
	}

	llmOpts := LLMPromptOptions{
		Format:               opts.Format,
		IncludeFoundations:   opts.IncludeFoundations,
		IncludeComponents:    opts.IncludeComponents,
		IncludePatterns:      opts.IncludePatterns,
		IncludeAccessibility: opts.IncludeAccessibility,
		IncludeAntiPatterns:  opts.IncludeAntiPatterns,
		MaxExamples:          opts.MaxExamples,
	}

	return s.ds.GenerateLLMPrompt(llmOpts)
}

// GenerateComponentPrompt generates an LLM prompt for a specific component.
func (s *Service) GenerateComponentPrompt(ctx context.Context, componentID string) (string, error) {
	comp, err := s.GetComponent(ctx, componentID)
	if err != nil {
		return "", err
	}

	// Create a temporary DesignSystem with just this component
	tempDS := &DesignSystem{
		Meta:        s.ds.Meta,
		Foundations: s.ds.Foundations,
		Components:  []Component{*comp},
	}

	return tempDS.GenerateLLMPrompt(LLMPromptOptions{
		Format:              "markdown",
		IncludeFoundations:  true,
		IncludeComponents:   true,
		IncludeAntiPatterns: true,
		MaxExamples:         5,
	})
}

// GetAccessibility returns the system-wide accessibility requirements.
func (s *Service) GetAccessibility(_ context.Context) *Accessibility {
	return s.ds.Accessibility
}
