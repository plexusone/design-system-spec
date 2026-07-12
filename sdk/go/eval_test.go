package dss

import (
	"context"
	"testing"

	"github.com/plexusone/structured-evaluation/rubric"
)

func TestEvaluate(t *testing.T) {
	ctx := context.Background()

	t.Run("complete design system", func(t *testing.T) {
		ds := createCompleteDesignSystem()
		service := NewService(ds)

		result, err := service.Evaluate(ctx, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// IntScore should be Good (4) or Excellent (5) for a complete system
		if result.IntScore < rubric.ScoreGood {
			t.Errorf("expected IntScore >= 4 (Good) for complete system, got %d", result.IntScore)
		}

		if len(result.Categories) != 4 {
			t.Errorf("expected 4 categories, got %d", len(result.Categories))
		}
	})

	t.Run("minimal design system", func(t *testing.T) {
		ds := &DesignSystem{
			Meta: Meta{
				Name:    "Minimal System",
				Version: "1.0.0",
			},
		}
		service := NewService(ds)

		result, err := service.Evaluate(ctx, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// IntScore should be low for minimal system
		if result.IntScore > rubric.ScoreAcceptable {
			t.Errorf("expected low IntScore for minimal system, got %d", result.IntScore)
		}

		if len(result.Findings) == 0 {
			t.Error("expected findings for minimal system")
		}
	})

	t.Run("with options", func(t *testing.T) {
		ds := createCompleteDesignSystem()
		service := NewService(ds)

		opts := &EvalOptions{
			IncludeIssues: false,
			Categories:    []string{"completeness", "documentation"},
		}

		result, err := service.Evaluate(ctx, opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result.Categories) != 2 {
			t.Errorf("expected 2 categories, got %d", len(result.Categories))
		}
	})
}

func TestEvaluateCompleteness(t *testing.T) {
	ctx := context.Background()

	t.Run("missing meta fields", func(t *testing.T) {
		ds := &DesignSystem{
			Meta: Meta{
				Name: "", // Missing
			},
		}
		service := NewService(ds)

		result, err := service.Evaluate(ctx, &EvalOptions{
			Categories:    []string{"completeness"},
			IncludeIssues: true,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should have findings for missing name, version, etc.
		hasNameFinding := false
		for _, finding := range result.Findings {
			if finding.ID == "completeness-meta-name" {
				hasNameFinding = true
				break
			}
		}
		if !hasNameFinding {
			t.Error("expected completeness-meta-name finding")
		}
	})

	t.Run("missing foundations", func(t *testing.T) {
		ds := &DesignSystem{
			Meta: Meta{
				Name:    "Test",
				Version: "1.0.0",
			},
			Foundations: Foundations{
				Colors: nil, // Missing
			},
		}
		service := NewService(ds)

		result, err := service.Evaluate(ctx, &EvalOptions{
			Categories:    []string{"completeness"},
			IncludeIssues: true,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		hasColorFinding := false
		for _, finding := range result.Findings {
			if finding.ID == "completeness-colors" {
				hasColorFinding = true
				break
			}
		}
		if !hasColorFinding {
			t.Error("expected completeness-colors finding")
		}
	})
}

func TestEvaluateAgentReadiness(t *testing.T) {
	ctx := context.Background()

	t.Run("component without LLM context", func(t *testing.T) {
		ds := &DesignSystem{
			Meta: Meta{
				Name:    "Test",
				Version: "1.0.0",
			},
			Components: []Component{
				{
					ID:   "button",
					Name: "Button",
					LLM:  nil, // Missing LLM context
				},
			},
		}
		service := NewService(ds)

		result, err := service.Evaluate(ctx, &EvalOptions{
			Categories:    []string{"agent-readiness"},
			IncludeIssues: true,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		hasLLMFinding := false
		for _, finding := range result.Findings {
			if finding.ID == "agent-readiness-button-llm-context" {
				hasLLMFinding = true
				break
			}
		}
		if !hasLLMFinding {
			t.Error("expected agent-readiness-button-llm-context finding")
		}
	})

	t.Run("component with complete LLM context", func(t *testing.T) {
		ds := &DesignSystem{
			Meta: Meta{
				Name:    "Test",
				Version: "1.0.0",
			},
			Components: []Component{
				{
					ID:   "button",
					Name: "Button",
					LLM: &LLMContext{
						Intent:       "Primary action trigger",
						AntiPatterns: []string{"Multiple primary buttons"},
						ExampleUsage: []string{"<Button variant='primary'>Submit</Button>"},
					},
				},
			},
			Principles: []Principle{
				{ID: "simplicity", Name: "Simplicity"},
				{ID: "consistency", Name: "Consistency"},
			},
		}
		service := NewService(ds)

		result, err := service.Evaluate(ctx, &EvalOptions{
			Categories:    []string{"agent-readiness"},
			IncludeIssues: true,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should have high score with complete LLM context
		for _, cat := range result.Categories {
			if cat.Category == "agent-readiness" && cat.IntScore < rubric.ScoreAcceptable {
				t.Errorf("expected agent-readiness IntScore >= 3 (Acceptable), got %d", cat.IntScore)
			}
		}
	})
}

func TestEvaluateAccessibility(t *testing.T) {
	ctx := context.Background()

	t.Run("missing accessibility section", func(t *testing.T) {
		ds := &DesignSystem{
			Meta: Meta{
				Name:    "Test",
				Version: "1.0.0",
			},
			Accessibility: nil,
		}
		service := NewService(ds)

		result, err := service.Evaluate(ctx, &EvalOptions{
			Categories:    []string{"accessibility"},
			IncludeIssues: true,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		hasA11yFinding := false
		for _, finding := range result.Findings {
			if finding.ID == "accessibility-section" {
				hasA11yFinding = true
				break
			}
		}
		if !hasA11yFinding {
			t.Error("expected accessibility-section finding")
		}
	})

	t.Run("component without accessibility", func(t *testing.T) {
		ds := &DesignSystem{
			Meta: Meta{
				Name:    "Test",
				Version: "1.0.0",
			},
			Accessibility: &Accessibility{
				WCAGLevel: "AA",
			},
			Components: []Component{
				{
					ID:            "button",
					Name:          "Button",
					Accessibility: nil, // Missing
				},
			},
		}
		service := NewService(ds)

		result, err := service.Evaluate(ctx, &EvalOptions{
			Categories:    []string{"accessibility"},
			IncludeIssues: true,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		hasCompA11yFinding := false
		for _, finding := range result.Findings {
			if finding.ID == "accessibility-component-button" {
				hasCompA11yFinding = true
				break
			}
		}
		if !hasCompA11yFinding {
			t.Error("expected accessibility-component-button finding")
		}
	})
}

func TestEvaluateDocumentation(t *testing.T) {
	ctx := context.Background()

	t.Run("missing descriptions", func(t *testing.T) {
		ds := &DesignSystem{
			Meta: Meta{
				Name:        "Test",
				Version:     "1.0.0",
				Description: "", // Missing
			},
			Components: []Component{
				{
					ID:          "button",
					Name:        "Button",
					Description: "", // Missing
					Props: []Prop{
						{
							Name:        "variant",
							Type:        "string",
							Description: "", // Missing
						},
					},
				},
			},
		}
		service := NewService(ds)

		result, err := service.Evaluate(ctx, &EvalOptions{
			Categories:    []string{"documentation"},
			IncludeIssues: true,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Check for documentation findings
		hasMetaDescFinding := false
		hasCompDescFinding := false
		for _, finding := range result.Findings {
			if finding.ID == "documentation-meta-description" {
				hasMetaDescFinding = true
			}
			if finding.ID == "documentation-component-button-description" {
				hasCompDescFinding = true
			}
		}
		if !hasMetaDescFinding {
			t.Error("expected documentation-meta-description finding")
		}
		if !hasCompDescFinding {
			t.Error("expected documentation-component-button-description finding")
		}
	})
}

func TestDefaultEvalOptions(t *testing.T) {
	opts := DefaultEvalOptions()

	if !opts.IncludeIssues {
		t.Error("expected IncludeIssues to be true by default")
	}

	if opts.MinScore != 4 { // Changed from 80 to 4 (IntegerScore)
		t.Errorf("expected MinScore to be 4, got %d", opts.MinScore)
	}

	if len(opts.Categories) != 4 {
		t.Errorf("expected 4 categories, got %d", len(opts.Categories))
	}
}

func TestPercentageToIntScore(t *testing.T) {
	// Thresholds: >= 90 Excellent, >= 75 Good, >= 50 Acceptable, >= 25 MajorRevisions, < 25 Unacceptable
	tests := []struct {
		passed   int
		checks   int
		expected rubric.IntegerScore
	}{
		{10, 10, rubric.ScoreExcellent},     // 100%
		{9, 10, rubric.ScoreExcellent},      // 90%
		{8, 10, rubric.ScoreGood},           // 80%
		{7, 10, rubric.ScoreAcceptable},     // 70% - below 75 threshold
		{6, 10, rubric.ScoreAcceptable},     // 60%
		{5, 10, rubric.ScoreAcceptable},     // 50%
		{4, 10, rubric.ScoreMajorRevisions}, // 40%
		{3, 10, rubric.ScoreMajorRevisions}, // 30%
		{2, 10, rubric.ScoreUnacceptable},   // 20%
		{1, 10, rubric.ScoreUnacceptable},   // 10%
		{0, 10, rubric.ScoreUnacceptable},   // 0%
		{0, 0, rubric.ScoreExcellent},       // Edge case: no checks
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			score := percentageToIntScore(tt.passed, tt.checks)
			if score != tt.expected {
				t.Errorf("percentageToIntScore(%d, %d) = %d, expected %d", tt.passed, tt.checks, score, tt.expected)
			}
		})
	}
}

// createCompleteDesignSystem creates a well-populated design system for testing.
func createCompleteDesignSystem() *DesignSystem {
	return &DesignSystem{
		Meta: Meta{
			Name:        "Complete System",
			Version:     "1.0.0",
			Description: "A complete design system for testing",
		},
		Principles: []Principle{
			{ID: "simplicity", Name: "Simplicity", Description: "Keep it simple"},
			{ID: "consistency", Name: "Consistency", Description: "Be consistent"},
			{ID: "accessibility", Name: "Accessibility", Description: "Accessible to all"},
		},
		Foundations: Foundations{
			Colors: []ColorToken{
				{ID: "primary", Value: "#0066cc", Usage: "Primary actions"},
				{ID: "secondary", Value: "#6c757d", Usage: "Secondary actions"},
			},
			Spacing: &SpacingScale{
				BaseUnit: "4px",
				Scale: []SpacingToken{
					{ID: "1", Value: "4px"},
					{ID: "2", Value: "8px"},
					{ID: "4", Value: "16px"},
				},
			},
			Typography: &Typography{
				FontFamilies: []FontFamily{
					{ID: "sans", Name: "Sans Serif", Stack: "Inter, system-ui, sans-serif"},
				},
				FontSizes: []FontSize{
					{ID: "sm", Value: "0.875rem"},
					{ID: "base", Value: "1rem"},
				},
				FontWeights: []FontWeight{
					{ID: "normal", Value: 400},
					{ID: "bold", Value: 700},
				},
				LineHeights: []LineHeight{
					{ID: "normal", Value: "1.5"},
				},
			},
			Elevation: []ElevationToken{
				{ID: "sm", Value: "0 1px 2px rgba(0,0,0,0.1)"},
			},
		},
		Components: []Component{
			{
				ID:          "button",
				Name:        "Button",
				Description: "A clickable button component",
				Props: []Prop{
					{
						Name:        "variant",
						Type:        "enum",
						Description: "Visual style of the button",
						EnumValues:  []string{"primary", "secondary", "outline"},
					},
					{
						Name:        "size",
						Type:        "enum",
						Description: "Size of the button",
						EnumValues:  []string{"sm", "md", "lg"},
					},
				},
				Variants: []Variant{
					{ID: "primary", Name: "Primary", Description: "Main action button"},
					{ID: "secondary", Name: "Secondary", Description: "Secondary action"},
				},
				Accessibility: &ComponentA11y{
					Role: "button",
					KeyboardSupport: []KeyboardInteraction{
						{Key: "Enter", Action: "Activate button"},
						{Key: "Space", Action: "Activate button"},
					},
				},
				LLM: &LLMContext{
					Intent:       "Trigger user actions",
					AntiPatterns: []string{"Multiple primary buttons per view"},
					ExampleUsage: []string{"<Button variant='primary'>Submit</Button>"},
				},
			},
		},
		Patterns: []Pattern{
			{
				ID:          "form",
				Name:        "Form Pattern",
				Description: "Standard form layout with validation",
				Components: []PatternComponent{
					{ComponentID: "button"},
				},
			},
		},
		Accessibility: &Accessibility{
			WCAGLevel:   "AA",
			WCAGVersion: "2.1",
			Keyboard: &KeyboardRequirements{
				FocusVisible:   true,
				FocusOrder:     true,
				NoKeyboardTrap: true,
			},
			ScreenReader: &ScreenReaderRequirements{
				SemanticHTML: true,
				ARIALabels:   true,
			},
		},
	}
}
