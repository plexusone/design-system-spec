package dss

import (
	"context"
	"testing"
)

func TestGetAccessibilityRequirements(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
		Accessibility: &Accessibility{
			WCAGLevel: "AA",
			ColorContrast: &ColorContrastRequirements{
				NormalTextRatio: 4.5,
				LargeTextRatio:  3.0,
				NonTextRatio:    3.0,
			},
		},
		Components: []Component{
			{
				ID:   "button",
				Name: "Button",
				Accessibility: &ComponentA11y{
					Role:               "button",
					RequiredAttributes: []string{"aria-label"},
					KeyboardSupport: []KeyboardInteraction{
						{Key: "Enter", Action: "activates button"},
						{Key: "Space", Action: "activates button"},
					},
					FocusManagement: "visible",
				},
			},
		},
	}
	service := NewService(ds)
	ctx := context.Background()

	reqs, err := service.GetAccessibilityRequirements(ctx, "button")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if reqs.ComponentID != "button" {
		t.Errorf("expected component button, got %s", reqs.ComponentID)
	}

	if reqs.Role != "button" {
		t.Errorf("expected role button, got %s", reqs.Role)
	}

	if reqs.Keyboard == nil {
		t.Fatal("expected keyboard requirements")
	}

	if reqs.Keyboard.Interactions["Enter"] != "activates button" {
		t.Errorf("expected Enter interaction, got %v", reqs.Keyboard.Interactions)
	}

	if reqs.ColorContrast == nil {
		t.Fatal("expected color contrast requirements")
	}

	if reqs.ColorContrast.Text != "4.5:1" {
		t.Errorf("expected 4.5:1 contrast, got %s", reqs.ColorContrast.Text)
	}

	if len(reqs.WCAGCriteria) == 0 {
		t.Error("expected WCAG criteria")
	}
}

func TestGetAccessibilityRequirementsNotFound(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
	}
	service := NewService(ds)
	ctx := context.Background()

	_, err := service.GetAccessibilityRequirements(ctx, "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent component")
	}
}

func TestGetAntiPatterns(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
		Components: []Component{
			{
				ID:   "button",
				Name: "Button",
				LLM: &LLMContext{
					AntiPatterns: []string{"Don't use divs as buttons"},
				},
			},
		},
	}
	service := NewService(ds)
	ctx := context.Background()

	// Get all anti-patterns
	result, err := service.GetAntiPatterns(ctx, "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.AntiPatterns) == 0 {
		t.Fatal("expected anti-patterns")
	}

	// Get anti-patterns for button
	result, err = service.GetAntiPatterns(ctx, "button", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should include built-in and custom anti-patterns for button
	found := false
	for _, ap := range result.AntiPatterns {
		if ap.Description == "Don't use divs as buttons" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected custom anti-pattern from LLM context")
	}

	// Get anti-patterns for color-contrast rule
	result, err = service.GetAntiPatterns(ctx, "", "color-contrast")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.AntiPatterns) == 0 {
		t.Error("expected anti-patterns for color-contrast rule")
	}
}

func TestSuggestContrastToken(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
		Foundations: Foundations{
			Colors: []ColorToken{
				{ID: "white", Value: "#FFFFFF"},
				{ID: "black", Value: "#000000"},
				{ID: "gray-500", Value: "#6B7280"},
				{ID: "primary", Value: "#3B82F6"},
			},
		},
	}
	service := NewService(ds)
	ctx := context.Background()

	// Suggest tokens that contrast with white background
	suggestions, err := service.SuggestContrastToken(ctx, "#FFFFFF", 4.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(suggestions) == 0 {
		t.Fatal("expected contrast suggestions")
	}

	// Black should definitely be suggested for white background
	found := false
	for _, s := range suggestions {
		if s.Token == "black" {
			found = true
			if !s.MeetsAA {
				t.Error("black on white should meet AA")
			}
			if !s.MeetsAAA {
				t.Error("black on white should meet AAA")
			}
			break
		}
	}
	if !found {
		t.Error("expected black to be suggested for white background")
	}
}

func TestSuggestContrastTokenInvalidColor(t *testing.T) {
	ds := &DesignSystem{
		Meta:        Meta{Name: "Test", Version: "1.0.0"},
		Foundations: Foundations{},
	}
	service := NewService(ds)
	ctx := context.Background()

	_, err := service.SuggestContrastToken(ctx, "invalid", 4.5)
	if err == nil {
		t.Error("expected error for invalid color")
	}
}

func TestGetComponentFixContext(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
		Components: []Component{
			{
				ID:   "button",
				Name: "Button",
				Uses: []string{"icon"},
			},
		},
		Foundations: Foundations{
			Colors: []ColorToken{
				{ID: "primary", Value: "#3B82F6"},
			},
		},
	}
	service := NewService(ds)
	ctx := context.Background()

	context, err := service.GetComponentFixContext(ctx, "button", "color-contrast")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if context.ComponentID != "button" {
		t.Errorf("expected button, got %s", context.ComponentID)
	}

	if context.IssueType != "color-contrast" {
		t.Errorf("expected color-contrast, got %s", context.IssueType)
	}

	if len(context.RelatedComponents) == 0 {
		t.Error("expected related components")
	}

	if context.TokensAvailable["colors"] == nil {
		t.Error("expected color tokens")
	}

	if len(context.StylesToCheck) == 0 {
		t.Error("expected styles to check for color-contrast")
	}
}

func TestGetComponentFixContextNotFound(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
	}
	service := NewService(ds)
	ctx := context.Background()

	_, err := service.GetComponentFixContext(ctx, "nonexistent", "color-contrast")
	if err == nil {
		t.Error("expected error for nonexistent component")
	}
}
