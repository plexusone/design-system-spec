package dss

import (
	"encoding/json"
	"testing"
)

func TestThemingContractSerialization(t *testing.T) {
	contract := ThemingContract{
		Prefix:      "--btn",
		Description: "Button theming tokens",
		Tokens: []ThemeToken{
			{
				ID:           "background",
				CSSProperty:  "--btn-background",
				Semantic:     "primary",
				Description:  "Button background color",
				DefaultLight: "#0066CC",
				DefaultDark:  "#3399FF",
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(contract)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	// Unmarshal back
	var decoded ThemingContract
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	// Verify
	if decoded.Prefix != contract.Prefix {
		t.Errorf("prefix: got %q, want %q", decoded.Prefix, contract.Prefix)
	}
	if decoded.Description != contract.Description {
		t.Errorf("description: got %q, want %q", decoded.Description, contract.Description)
	}
	if len(decoded.Tokens) != 1 {
		t.Errorf("tokens length: got %d, want 1", len(decoded.Tokens))
	}
	if decoded.Tokens[0].ID != "background" {
		t.Errorf("token id: got %q, want %q", decoded.Tokens[0].ID, "background")
	}
	if decoded.Tokens[0].Semantic != "primary" {
		t.Errorf("token semantic: got %q, want %q", decoded.Tokens[0].Semantic, "primary")
	}
}

func TestThemeBindingsSerialization(t *testing.T) {
	bindings := ThemeBindings{
		Component: "button",
		SpecURL:   "https://example.com/button.json",
		ThemeMode: "dark",
		Strategy:  "semantic",
		Mappings: []TokenMapping{
			{
				From:      "brand-primary",
				To:        "background",
				Transform: "rgb",
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(bindings)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	// Unmarshal back
	var decoded ThemeBindings
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	// Verify
	if decoded.Component != bindings.Component {
		t.Errorf("component: got %q, want %q", decoded.Component, bindings.Component)
	}
	if decoded.SpecURL != bindings.SpecURL {
		t.Errorf("specUrl: got %q, want %q", decoded.SpecURL, bindings.SpecURL)
	}
	if decoded.ThemeMode != bindings.ThemeMode {
		t.Errorf("themeMode: got %q, want %q", decoded.ThemeMode, bindings.ThemeMode)
	}
	if decoded.Strategy != bindings.Strategy {
		t.Errorf("strategy: got %q, want %q", decoded.Strategy, bindings.Strategy)
	}
	if len(decoded.Mappings) != 1 {
		t.Errorf("mappings length: got %d, want 1", len(decoded.Mappings))
	}
	if decoded.Mappings[0].Transform != "rgb" {
		t.Errorf("mapping transform: got %q, want %q", decoded.Mappings[0].Transform, "rgb")
	}
}

func TestIsValidSemantic(t *testing.T) {
	tests := []struct {
		semantic string
		want     bool
	}{
		{"", true},              // Empty is allowed
		{"primary", true},       // Valid
		{"secondary", true},     // Valid
		{"danger", true},        // Valid
		{"text-muted", true},    // Valid
		{"invalid", false},      // Not in list
		{"PRIMARY", false},      // Case sensitive
		{"text_muted", false},   // Wrong format
	}

	for _, tt := range tests {
		got := IsValidSemantic(tt.semantic)
		if got != tt.want {
			t.Errorf("IsValidSemantic(%q) = %v, want %v", tt.semantic, got, tt.want)
		}
	}
}
