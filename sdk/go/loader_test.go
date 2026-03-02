package dss

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestLoadMinimalExample(t *testing.T) {
	// Get the path to the examples directory relative to this test file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get current file path")
	}
	sdkDir := filepath.Dir(filename)
	exampleDir := filepath.Join(sdkDir, "..", "..", "examples", "minimal-system")

	ds, err := LoadDesignSystem(exampleDir)
	if err != nil {
		t.Fatalf("LoadDesignSystem failed: %v", err)
	}

	// Verify meta
	if ds.Meta.Name != "Minimal Design System" {
		t.Errorf("Meta.Name = %q, want %q", ds.Meta.Name, "Minimal Design System")
	}
	if ds.Meta.Version != "1.0.0" {
		t.Errorf("Meta.Version = %q, want %q", ds.Meta.Version, "1.0.0")
	}

	// Verify foundations
	if len(ds.Foundations.Colors) != 7 {
		t.Errorf("len(Colors) = %d, want 7", len(ds.Foundations.Colors))
	}
	if ds.Foundations.Typography == nil {
		t.Error("Typography is nil")
	} else {
		if len(ds.Foundations.Typography.FontFamilies) != 2 {
			t.Errorf("len(FontFamilies) = %d, want 2", len(ds.Foundations.Typography.FontFamilies))
		}
	}
	if ds.Foundations.Spacing == nil {
		t.Error("Spacing is nil")
	} else {
		if ds.Foundations.Spacing.BaseUnit != "4px" {
			t.Errorf("Spacing.BaseUnit = %q, want %q", ds.Foundations.Spacing.BaseUnit, "4px")
		}
	}

	// Verify components
	if len(ds.Components) != 2 {
		t.Errorf("len(Components) = %d, want 2", len(ds.Components))
	}

	// Find button component
	var button *Component
	for i := range ds.Components {
		if ds.Components[i].ID == "button" {
			button = &ds.Components[i]
			break
		}
	}
	if button == nil {
		t.Fatal("Button component not found")
	}
	if len(button.Variants) != 5 {
		t.Errorf("len(Button.Variants) = %d, want 5", len(button.Variants))
	}
	if button.LLM == nil {
		t.Error("Button.LLM is nil")
	} else {
		if button.LLM.PriorityScore != 90 {
			t.Errorf("Button.LLM.PriorityScore = %d, want 90", button.LLM.PriorityScore)
		}
	}

	// Validate the loaded system
	if err := ds.Validate(); err != nil {
		t.Errorf("Validate failed: %v", err)
	}
}
