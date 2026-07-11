package dss

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"testing/fstest"
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

func TestLoadDesignSystemFromFS(t *testing.T) {
	// Create an in-memory filesystem with a design system
	fsys := fstest.MapFS{
		"meta.json": &fstest.MapFile{
			Data: []byte(`{"name": "Test System", "version": "1.0.0"}`),
		},
		"foundations/colors.json": &fstest.MapFile{
			Data: []byte(`[{"id": "primary", "value": "#3B82F6"}]`),
		},
		"components/button.json": &fstest.MapFile{
			Data: []byte(`{"id": "button", "name": "Button", "variants": [{"id": "default"}]}`),
		},
	}

	ds, err := LoadDesignSystemFromFS(fsys)
	if err != nil {
		t.Fatalf("LoadDesignSystemFromFS failed: %v", err)
	}

	// Verify meta
	if ds.Meta.Name != "Test System" {
		t.Errorf("Meta.Name = %q, want %q", ds.Meta.Name, "Test System")
	}
	if ds.Meta.Version != "1.0.0" {
		t.Errorf("Meta.Version = %q, want %q", ds.Meta.Version, "1.0.0")
	}

	// Verify foundations
	if len(ds.Foundations.Colors) != 1 {
		t.Errorf("len(Colors) = %d, want 1", len(ds.Foundations.Colors))
	}
	if ds.Foundations.Colors[0].ID != "primary" {
		t.Errorf("Colors[0].ID = %q, want %q", ds.Foundations.Colors[0].ID, "primary")
	}

	// Verify components
	if len(ds.Components) != 1 {
		t.Errorf("len(Components) = %d, want 1", len(ds.Components))
	}
	if ds.Components[0].ID != "button" {
		t.Errorf("Components[0].ID = %q, want %q", ds.Components[0].ID, "button")
	}
}

func TestLoadDesignSystemFromFSSingleFile(t *testing.T) {
	// Create an in-memory filesystem with a single design-system.json file
	fsys := fstest.MapFS{
		"design-system.json": &fstest.MapFile{
			Data: []byte(`{
				"meta": {"name": "Single File System", "version": "2.0.0"},
				"components": [{"id": "card", "name": "Card"}]
			}`),
		},
	}

	ds, err := LoadDesignSystemFromFS(fsys)
	if err != nil {
		t.Fatalf("LoadDesignSystemFromFS failed: %v", err)
	}

	if ds.Meta.Name != "Single File System" {
		t.Errorf("Meta.Name = %q, want %q", ds.Meta.Name, "Single File System")
	}
	if len(ds.Components) != 1 {
		t.Errorf("len(Components) = %d, want 1", len(ds.Components))
	}
}

func TestLoadDesignSystemFromFSWithRealDirectory(t *testing.T) {
	// Load from the real minimal-system example using os.DirFS
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get current file path")
	}
	sdkDir := filepath.Dir(filename)
	exampleDir := filepath.Join(sdkDir, "..", "..", "examples", "minimal-system")

	fsys := os.DirFS(exampleDir)

	ds, err := LoadDesignSystemFromFS(fsys)
	if err != nil {
		t.Fatalf("LoadDesignSystemFromFS failed: %v", err)
	}

	// Verify meta
	if ds.Meta.Name != "Minimal Design System" {
		t.Errorf("Meta.Name = %q, want %q", ds.Meta.Name, "Minimal Design System")
	}

	// Verify foundations
	if len(ds.Foundations.Colors) != 7 {
		t.Errorf("len(Colors) = %d, want 7", len(ds.Foundations.Colors))
	}

	// Verify components
	if len(ds.Components) != 2 {
		t.Errorf("len(Components) = %d, want 2", len(ds.Components))
	}
}
