package dss

import (
	"bytes"
	"strings"
	"testing"
)

func TestGenerateBindings(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
		Foundations: Foundations{
			Colors: []ColorToken{
				{ID: "brand-primary", Value: "#8B5CF6"},
				{ID: "brand-secondary", Value: "#06B6D4"},
				{ID: "text-dark", Value: "#0F172A"},
			},
		},
		Components: []Component{
			{
				ID:   "button",
				Name: "Button",
				ThemingContract: &ThemingContract{
					Prefix: "--btn",
					Tokens: []ThemeToken{
						{ID: "background", CSSProperty: "--btn-background", Semantic: "primary", DefaultLight: "#0066CC", DefaultDark: "#3399FF"},
						{ID: "text", CSSProperty: "--btn-text", Semantic: "text", DefaultLight: "#FFFFFF", DefaultDark: "#0A0E1A"},
					},
				},
			},
		},
		ThemeBindings: []ThemeBindings{
			{
				Component: "button",
				Strategy:  "explicit",
				Mappings: []TokenMapping{
					{From: "brand-primary", To: "background"},
					{From: "text-dark", To: "text"},
				},
			},
		},
	}

	opts := BindingOptions{
		Format:          FormatCSS,
		DefaultStrategy: "explicit",
	}

	bindings, err := GenerateBindings(ds, opts)
	if err != nil {
		t.Fatalf("GenerateBindings failed: %v", err)
	}

	if len(bindings) != 1 {
		t.Fatalf("Expected 1 binding, got %d", len(bindings))
	}

	b := bindings[0]
	if b.Component != "button" {
		t.Errorf("Component = %q, want %q", b.Component, "button")
	}

	// Check CSS output contains expected values
	if !strings.Contains(b.CSS, "--btn-background: #8B5CF6") {
		t.Errorf("CSS missing --btn-background mapping; got: %s", b.CSS)
	}
	if !strings.Contains(b.CSS, "--btn-text: #0F172A") {
		t.Errorf("CSS missing --btn-text mapping; got: %s", b.CSS)
	}
}

func TestGenerateBindings_SemanticStrategy(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
		Foundations: Foundations{
			Colors: []ColorToken{
				{ID: "primary-500", Value: "#3B82F6"},
				{ID: "secondary-500", Value: "#64748B"},
			},
		},
		Components: []Component{
			{
				ID:   "card",
				Name: "Card",
				ThemingContract: &ThemingContract{
					Prefix: "--card",
					Tokens: []ThemeToken{
						{ID: "background", CSSProperty: "--card-background", Semantic: "primary", DefaultLight: "#FFF", DefaultDark: "#000"},
						{ID: "border", CSSProperty: "--card-border", Semantic: "secondary", DefaultLight: "#EEE", DefaultDark: "#333"},
					},
				},
			},
		},
		ThemeBindings: []ThemeBindings{
			{
				Component: "card",
				Strategy:  "semantic",
				Mappings:  []TokenMapping{}, // Empty - rely on semantic matching
			},
		},
	}

	opts := BindingOptions{
		Format:          FormatCSS,
		DefaultStrategy: "semantic",
	}

	bindings, err := GenerateBindings(ds, opts)
	if err != nil {
		t.Fatalf("GenerateBindings failed: %v", err)
	}

	if len(bindings) != 1 {
		t.Fatalf("Expected 1 binding, got %d", len(bindings))
	}

	// Semantic strategy should auto-map by semantic field
	if !strings.Contains(bindings[0].CSS, "--card-background") {
		t.Errorf("CSS missing --card-background; got: %s", bindings[0].CSS)
	}
}

func TestGenerateBindings_InheritStrategy(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
		Components: []Component{
			{
				ID:   "badge",
				Name: "Badge",
				ThemingContract: &ThemingContract{
					Prefix: "--badge",
					Tokens: []ThemeToken{
						{ID: "bg", CSSProperty: "--badge-bg", DefaultLight: "#EEE", DefaultDark: "#333"},
					},
				},
			},
		},
		ThemeBindings: []ThemeBindings{
			{
				Component: "badge",
				ThemeMode: "dark",
				Strategy:  "inherit",
				Mappings:  []TokenMapping{},
			},
		},
	}

	opts := BindingOptions{Format: FormatCSS}

	bindings, err := GenerateBindings(ds, opts)
	if err != nil {
		t.Fatalf("GenerateBindings failed: %v", err)
	}

	// Inherit strategy should use component defaults
	if !strings.Contains(bindings[0].CSS, "--badge-bg: #333") {
		t.Errorf("CSS should use dark default; got: %s", bindings[0].CSS)
	}
}

func TestGenerateBindings_TypeScriptFormat(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
		Foundations: Foundations{
			Colors: []ColorToken{
				{ID: "brand", Value: "#FF0000"},
			},
		},
		Components: []Component{
			{
				ID:   "alert",
				Name: "Alert",
				ThemingContract: &ThemingContract{
					Prefix: "--alert",
					Tokens: []ThemeToken{
						{ID: "bg-color", CSSProperty: "--alert-bg-color", DefaultLight: "#FFF", DefaultDark: "#000"},
					},
				},
			},
		},
		ThemeBindings: []ThemeBindings{
			{
				Component: "alert",
				Strategy:  "explicit",
				Mappings: []TokenMapping{
					{From: "brand", To: "bg-color"},
				},
			},
		},
	}

	opts := BindingOptions{Format: FormatTypeScript}

	bindings, err := GenerateBindings(ds, opts)
	if err != nil {
		t.Fatalf("GenerateBindings failed: %v", err)
	}

	// Check TypeScript output
	if !strings.Contains(bindings[0].CSS, "export const alertTheme") {
		t.Errorf("TypeScript missing theme export; got: %s", bindings[0].CSS)
	}
	if !strings.Contains(bindings[0].CSS, "bgColor: '#FF0000'") {
		t.Errorf("TypeScript missing camelCase property; got: %s", bindings[0].CSS)
	}
}

func TestGenerateBindings_SCSSFormat(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
		Foundations: Foundations{
			Colors: []ColorToken{
				{ID: "accent", Value: "#00FF00"},
			},
		},
		Components: []Component{
			{
				ID:   "tag",
				Name: "Tag",
				ThemingContract: &ThemingContract{
					Prefix: "--tag",
					Tokens: []ThemeToken{
						{ID: "fill", CSSProperty: "--tag-fill", DefaultLight: "#FFF", DefaultDark: "#000"},
					},
				},
			},
		},
		ThemeBindings: []ThemeBindings{
			{
				Component: "tag",
				Strategy:  "explicit",
				Mappings: []TokenMapping{
					{From: "accent", To: "fill"},
				},
			},
		},
	}

	opts := BindingOptions{Format: FormatSCSS}

	bindings, err := GenerateBindings(ds, opts)
	if err != nil {
		t.Fatalf("GenerateBindings failed: %v", err)
	}

	// Check SCSS output has both variables and CSS
	if !strings.Contains(bindings[0].CSS, "$tag-fill: #00FF00") {
		t.Errorf("SCSS missing variable; got: %s", bindings[0].CSS)
	}
	if !strings.Contains(bindings[0].CSS, "--tag-fill: #00FF00") {
		t.Errorf("SCSS missing CSS property; got: %s", bindings[0].CSS)
	}
}

func TestGenerateBindings_NoContract(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
		Components: []Component{
			{ID: "button", Name: "Button"}, // No theming contract
		},
		ThemeBindings: []ThemeBindings{
			{Component: "button", Mappings: []TokenMapping{}},
		},
	}

	opts := BindingOptions{Format: FormatCSS}

	bindings, err := GenerateBindings(ds, opts)
	if err != nil {
		t.Fatalf("GenerateBindings failed: %v", err)
	}

	// Should return binding with warning
	if len(bindings) != 1 {
		t.Fatalf("Expected 1 binding, got %d", len(bindings))
	}
	if len(bindings[0].Warnings) == 0 {
		t.Error("Expected warning about missing contract")
	}
}

func TestWriteBindings(t *testing.T) {
	bindings := []GeneratedBinding{
		{Component: "button", CSS: "/* Button */\n:root { --btn-bg: #fff; }\n"},
		{Component: "card", CSS: "/* Card */\n:root { --card-bg: #eee; }\n"},
		{Component: "empty", CSS: ""},
	}

	var buf bytes.Buffer
	if err := WriteBindings(&buf, bindings); err != nil {
		t.Fatalf("WriteBindings failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "/* Button */") {
		t.Error("Output missing button CSS")
	}
	if !strings.Contains(output, "/* Card */") {
		t.Error("Output missing card CSS")
	}
	// Empty binding should be skipped
	if strings.Count(output, ":root") != 2 {
		t.Errorf("Expected 2 :root blocks, got output: %s", output)
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"background", "background"},
		{"bg-color", "bgColor"},
		{"text-primary-dark", "textPrimaryDark"},
		{"a-b-c", "aBC"},
		{"", ""},
	}

	for _, tt := range tests {
		got := toCamelCase(tt.input)
		if got != tt.want {
			t.Errorf("toCamelCase(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestResolveTokenValue(t *testing.T) {
	ds := &DesignSystem{
		Foundations: Foundations{
			Colors: []ColorToken{
				{ID: "primary-500", Value: "#3B82F6"},
				{ID: "colors.accent", Value: "#8B5CF6"},
			},
		},
	}

	tests := []struct {
		tokenRef string
		want     string
	}{
		{"primary-500", "#3B82F6"},
		{"colors.primary-500", "#3B82F6"},          // dot notation with lookup
		{"colors.accent", "#8B5CF6"},               // exact match with dot
		{"unknown", "var(--unknown)"},              // not found
		{"colors.unknown", "var(--colors-unknown)"}, // not found with dot
	}

	for _, tt := range tests {
		got := resolveTokenValue(ds, tt.tokenRef)
		if got != tt.want {
			t.Errorf("resolveTokenValue(%q) = %q, want %q", tt.tokenRef, got, tt.want)
		}
	}
}
