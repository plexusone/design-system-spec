package dss

import (
	"context"
	"testing"
)

func TestNewService(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:    "Test System",
			Version: "1.0.0",
		},
	}

	service := NewService(ds)

	if service == nil {
		t.Fatal("NewService returned nil")
	}

	if service.DesignSystem() != ds {
		t.Error("DesignSystem() did not return expected value")
	}
}

func TestServiceGetComponent(t *testing.T) {
	ds := &DesignSystem{
		Components: []Component{
			{ID: "button", Name: "Button", Description: "A clickable button"},
			{ID: "input", Name: "Input", Description: "A text input"},
		},
	}
	service := NewService(ds)
	ctx := context.Background()

	t.Run("existing component", func(t *testing.T) {
		comp, err := service.GetComponent(ctx, "button")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if comp.ID != "button" {
			t.Errorf("expected ID 'button', got '%s'", comp.ID)
		}
		if comp.Name != "Button" {
			t.Errorf("expected Name 'Button', got '%s'", comp.Name)
		}
	})

	t.Run("non-existent component", func(t *testing.T) {
		_, err := service.GetComponent(ctx, "nonexistent")
		if err == nil {
			t.Error("expected error for non-existent component")
		}
	})
}

func TestServiceListComponents(t *testing.T) {
	ds := &DesignSystem{
		Components: []Component{
			{ID: "button", Name: "Button", Category: "actions"},
			{ID: "input", Name: "Input", Category: "forms"},
			{ID: "modal", Name: "Modal", Category: "overlays"},
		},
	}
	service := NewService(ds)
	ctx := context.Background()

	summaries := service.ListComponents(ctx)

	if len(summaries) != 3 {
		t.Fatalf("expected 3 components, got %d", len(summaries))
	}

	// Check first component
	if summaries[0].ID != "button" {
		t.Errorf("expected first ID 'button', got '%s'", summaries[0].ID)
	}
	if summaries[0].Category != "actions" {
		t.Errorf("expected first category 'actions', got '%s'", summaries[0].Category)
	}
}

func TestServiceGetComponentVariants(t *testing.T) {
	ds := &DesignSystem{
		Components: []Component{
			{
				ID:   "button",
				Name: "Button",
				Variants: []Variant{
					{ID: "primary", Name: "Primary", IsDefault: true},
					{ID: "secondary", Name: "Secondary"},
					{ID: "outline", Name: "Outline"},
				},
			},
		},
	}
	service := NewService(ds)
	ctx := context.Background()

	variants, err := service.GetComponentVariants(ctx, "button")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(variants) != 3 {
		t.Fatalf("expected 3 variants, got %d", len(variants))
	}

	if !variants[0].IsDefault {
		t.Error("expected first variant to be default")
	}
}

func TestServiceGetToken(t *testing.T) {
	ds := &DesignSystem{
		Foundations: Foundations{
			Colors: []ColorToken{
				{ID: "primary-500", Value: "#3B82F6", Semantic: "primary"},
				{ID: "danger-500", Value: "#EF4444", Semantic: "danger"},
			},
			Spacing: &SpacingScale{
				BaseUnit: "4px",
				Scale: []SpacingToken{
					{ID: "0", Value: "0px"},
					{ID: "1", Value: "4px"},
					{ID: "4", Value: "16px"},
				},
			},
		},
	}
	service := NewService(ds)
	ctx := context.Background()

	t.Run("get color token", func(t *testing.T) {
		token, err := service.GetToken(ctx, "color", "primary-500")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		colorToken, ok := token.(ColorToken)
		if !ok {
			t.Fatal("expected ColorToken type")
		}
		if colorToken.Value != "#3B82F6" {
			t.Errorf("expected value '#3B82F6', got '%s'", colorToken.Value)
		}
	})

	t.Run("get spacing token", func(t *testing.T) {
		token, err := service.GetToken(ctx, "spacing", "4")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		spacingToken, ok := token.(SpacingToken)
		if !ok {
			t.Fatal("expected SpacingToken type")
		}
		if spacingToken.Value != "16px" {
			t.Errorf("expected value '16px', got '%s'", spacingToken.Value)
		}
	})

	t.Run("unknown token type", func(t *testing.T) {
		_, err := service.GetToken(ctx, "unknown", "foo")
		if err == nil {
			t.Error("expected error for unknown token type")
		}
	})

	t.Run("non-existent token", func(t *testing.T) {
		_, err := service.GetToken(ctx, "color", "nonexistent")
		if err == nil {
			t.Error("expected error for non-existent token")
		}
	})
}

func TestServiceListTokens(t *testing.T) {
	ds := &DesignSystem{
		Foundations: Foundations{
			Colors: []ColorToken{
				{ID: "primary-500", Value: "#3B82F6"},
				{ID: "danger-500", Value: "#EF4444"},
			},
		},
	}
	service := NewService(ds)
	ctx := context.Background()

	tokens, err := service.ListTokens(ctx, "color")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}

	if tokens[0].Type != "color" {
		t.Errorf("expected type 'color', got '%s'", tokens[0].Type)
	}
}

func TestServiceGetPattern(t *testing.T) {
	ds := &DesignSystem{
		Patterns: []Pattern{
			{ID: "form-validation", Name: "Form Validation", Description: "Pattern for validating forms"},
			{ID: "data-table", Name: "Data Table", Description: "Pattern for displaying tabular data"},
		},
	}
	service := NewService(ds)
	ctx := context.Background()

	t.Run("existing pattern", func(t *testing.T) {
		pattern, err := service.GetPattern(ctx, "form-validation")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if pattern.Name != "Form Validation" {
			t.Errorf("expected name 'Form Validation', got '%s'", pattern.Name)
		}
	})

	t.Run("non-existent pattern", func(t *testing.T) {
		_, err := service.GetPattern(ctx, "nonexistent")
		if err == nil {
			t.Error("expected error for non-existent pattern")
		}
	})
}

func TestServiceListPatterns(t *testing.T) {
	ds := &DesignSystem{
		Patterns: []Pattern{
			{ID: "form-validation", Name: "Form Validation", Category: "forms"},
			{ID: "data-table", Name: "Data Table", Category: "data"},
		},
	}
	service := NewService(ds)
	ctx := context.Background()

	patterns := service.ListPatterns(ctx)

	if len(patterns) != 2 {
		t.Fatalf("expected 2 patterns, got %d", len(patterns))
	}

	if patterns[0].ID != "form-validation" {
		t.Errorf("expected first ID 'form-validation', got '%s'", patterns[0].ID)
	}
}

func TestServiceGetMeta(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:        "My Design System",
			Version:     "2.1.0",
			Description: "A comprehensive design system",
		},
	}
	service := NewService(ds)
	ctx := context.Background()

	meta := service.GetMeta(ctx)

	if meta.Name != "My Design System" {
		t.Errorf("expected name 'My Design System', got '%s'", meta.Name)
	}
	if meta.Version != "2.1.0" {
		t.Errorf("expected version '2.1.0', got '%s'", meta.Version)
	}
}

func TestServiceGetComponentAntiPatterns(t *testing.T) {
	ds := &DesignSystem{
		Components: []Component{
			{
				ID:   "button",
				Name: "Button",
				LLM: &LLMContext{
					AntiPatterns: []string{
						"Don't use for destructive actions",
						"Avoid multiple primary buttons",
					},
				},
			},
			{
				ID:   "input",
				Name: "Input",
				// No LLM context
			},
		},
	}
	service := NewService(ds)
	ctx := context.Background()

	t.Run("component with anti-patterns", func(t *testing.T) {
		antiPatterns, err := service.GetComponentAntiPatterns(ctx, "button")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(antiPatterns) != 2 {
			t.Fatalf("expected 2 anti-patterns, got %d", len(antiPatterns))
		}
	})

	t.Run("component without anti-patterns", func(t *testing.T) {
		antiPatterns, err := service.GetComponentAntiPatterns(ctx, "input")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if antiPatterns != nil {
			t.Errorf("expected nil anti-patterns, got %v", antiPatterns)
		}
	})
}

func TestServiceGenerateLLMPrompt(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:        "Test System",
			Version:     "1.0.0",
			Description: "A test design system",
		},
		Components: []Component{
			{ID: "button", Name: "Button", Description: "A clickable button"},
		},
	}
	service := NewService(ds)
	ctx := context.Background()

	t.Run("default options", func(t *testing.T) {
		prompt, err := service.GenerateLLMPrompt(ctx, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if prompt == "" {
			t.Error("expected non-empty prompt")
		}
		if !contains(prompt, "Test System") {
			t.Error("expected prompt to contain 'Test System'")
		}
	})

	t.Run("custom options", func(t *testing.T) {
		opts := &PromptOptions{
			Format:            "markdown",
			IncludeComponents: true,
		}
		prompt, err := service.GenerateLLMPrompt(ctx, opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !contains(prompt, "Button") {
			t.Error("expected prompt to contain 'Button'")
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && searchSubstring(s, substr)))
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
