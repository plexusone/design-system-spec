package dss

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGeneratePackage(t *testing.T) {
	// Load the minimal-system example
	ds, err := LoadDesignSystem("../../examples/minimal-system")
	if err != nil {
		t.Fatalf("failed to load design system: %v", err)
	}

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "dss-package-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("failed to remove temp dir: %v", err)
		}
	}()

	// Generate package
	opts := DefaultPackageOptions()
	opts.OutputDir = tmpDir
	opts.Targets = []PackageTarget{TargetCSS, TargetTailwind, TargetShadCN}

	if err := ds.GeneratePackage(opts); err != nil {
		t.Fatalf("failed to generate package: %v", err)
	}

	// Verify expected files exist
	expectedFiles := []string{
		"package.json",
		"index.js",
		"index.mjs",
		"index.d.ts",
		"README.md",
		"css/tokens.css",
		"tailwind/preset.js",
		"tailwind/theme.css",
		"tailwind/preset.d.ts",
		"shadcn/theme.css",
		"shadcn/colors.json",
	}

	for _, f := range expectedFiles {
		path := filepath.Join(tmpDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file not found: %s", f)
		}
	}
}

func TestParseTargets(t *testing.T) {
	tests := []struct {
		input    string
		expected []PackageTarget
	}{
		{"css", []PackageTarget{TargetCSS}},
		{"css,tailwind", []PackageTarget{TargetCSS, TargetTailwind}},
		{"shadcn,mkdocs", []PackageTarget{TargetShadCN, TargetMkDocsMaterial}},
		{"", []PackageTarget{TargetCSS, TargetTailwind}},
		{"all", []PackageTarget{
			TargetCSS, TargetTailwind, TargetShadCN,
			TargetMkDocsMaterial, TargetSCSS, TargetJSON, TargetW3C,
		}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseTargets(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("ParseTargets(%q) = %v, want %v", tt.input, result, tt.expected)
				return
			}
			for i, target := range result {
				if target != tt.expected[i] {
					t.Errorf("ParseTargets(%q)[%d] = %v, want %v", tt.input, i, target, tt.expected[i])
				}
			}
		})
	}
}

func TestHexToHSL(t *testing.T) {
	tests := []struct {
		hex      string
		expected string
	}{
		{"#ffffff", "0.0 0.0% 100.0%"},
		{"#000000", "0.0 0.0% 0.0%"},
		{"#ff0000", "0.0 100.0% 50.0%"},
	}

	for _, tt := range tests {
		t.Run(tt.hex, func(t *testing.T) {
			result := hexToHSL(tt.hex)
			if result != tt.expected {
				t.Errorf("hexToHSL(%q) = %q, want %q", tt.hex, result, tt.expected)
			}
		})
	}
}
