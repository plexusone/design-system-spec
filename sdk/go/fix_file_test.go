package dss

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFixHardcodedColors(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
		Foundations: Foundations{
			Colors: []ColorToken{
				{ID: "primary", Value: "#3B82F6"},
				{ID: "secondary", Value: "#64748B"},
				{ID: "error", Value: "#EF4444"},
			},
		},
	}
	service := NewService(ds)

	// Create temp file with hardcoded colors
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.css")
	content := `.button {
  background: #3B82F6;
  border-color: #64748B;
  color: #EF4444;
}`
	if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	// Run fix with dry run
	opts := &FixOptions{DryRun: true, IncludeOriginal: true}
	result, err := service.FixFile(context.Background(), testFile, opts)
	if err != nil {
		t.Fatalf("FixFile failed: %v", err)
	}

	if len(result.Fixes) != 3 {
		t.Errorf("expected 3 fixes, got %d", len(result.Fixes))
	}

	if result.Summary.ColorFixes != 3 {
		t.Errorf("expected 3 color fixes, got %d", result.Summary.ColorFixes)
	}

	// Verify replacements
	if !strings.Contains(result.FixedContent, "var(--color-primary)") {
		t.Error("expected fixed content to contain var(--color-primary)")
	}
	if !strings.Contains(result.FixedContent, "var(--color-secondary)") {
		t.Error("expected fixed content to contain var(--color-secondary)")
	}
	if !strings.Contains(result.FixedContent, "var(--color-error)") {
		t.Error("expected fixed content to contain var(--color-error)")
	}

	// Verify original file unchanged (dry run)
	originalContent, _ := os.ReadFile(testFile)
	if !strings.Contains(string(originalContent), "#3B82F6") {
		t.Error("dry run should not modify original file")
	}
}

func TestFixHardcodedSpacing(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
		Foundations: Foundations{
			Spacing: &SpacingScale{
				Scale: []SpacingToken{
					{ID: "4", Value: "16px", PixelValue: 16},
					{ID: "6", Value: "24px", PixelValue: 24},
				},
			},
		},
	}
	service := NewService(ds)

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.css")
	content := `.card {
  padding: 16px;
  margin: 24px;
}`
	if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	opts := &FixOptions{DryRun: true}
	result, err := service.FixFile(context.Background(), testFile, opts)
	if err != nil {
		t.Fatalf("FixFile failed: %v", err)
	}

	if len(result.Fixes) != 2 {
		t.Errorf("expected 2 fixes, got %d", len(result.Fixes))
	}

	if result.Summary.SpacingFixes != 2 {
		t.Errorf("expected 2 spacing fixes, got %d", result.Summary.SpacingFixes)
	}

	if !strings.Contains(result.FixedContent, "var(--spacing-4)") {
		t.Error("expected fixed content to contain var(--spacing-4)")
	}
	if !strings.Contains(result.FixedContent, "var(--spacing-6)") {
		t.Error("expected fixed content to contain var(--spacing-6)")
	}
}

func TestFixMissingAlt(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
	}
	service := NewService(ds)

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.tsx")
	content := `function App() {
  return (
    <div>
      <img src="/images/hero-banner.png" />
      <img src="/logo.svg" alt="Logo" />
    </div>
  );
}`
	if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	opts := &FixOptions{DryRun: true, Rules: []string{"img-alt-required"}}
	result, err := service.FixFile(context.Background(), testFile, opts)
	if err != nil {
		t.Fatalf("FixFile failed: %v", err)
	}

	// Should only fix the first image (second already has alt)
	if len(result.Fixes) != 1 {
		t.Errorf("expected 1 fix, got %d", len(result.Fixes))
	}

	if result.Summary.AccessibilityFixes != 1 {
		t.Errorf("expected 1 accessibility fix, got %d", result.Summary.AccessibilityFixes)
	}

	// Check that alt was derived from filename
	if !strings.Contains(result.FixedContent, `alt="hero banner"`) {
		t.Errorf("expected alt text derived from filename, got: %s", result.FixedContent)
	}
}

func TestFixMissingAriaLabel(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
	}
	service := NewService(ds)

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.tsx")
	content := `function Toolbar() {
  return (
    <div>
      <Button size="icon" icon={IconTrash}>
      <Button size="icon" aria-label="Delete">
    </div>
  );
}`
	if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	opts := &FixOptions{DryRun: true, Rules: []string{"button-accessible-name"}}
	result, err := service.FixFile(context.Background(), testFile, opts)
	if err != nil {
		t.Fatalf("FixFile failed: %v", err)
	}

	// Should only fix the first button
	if len(result.Fixes) != 1 {
		t.Errorf("expected 1 fix, got %d", len(result.Fixes))
	}

	if !strings.Contains(result.FixedContent, `aria-label=`) {
		t.Error("expected aria-label to be added")
	}
}

func TestSuggestFixes(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
		Foundations: Foundations{
			Colors: []ColorToken{
				{ID: "primary", Value: "#3B82F6"},
			},
		},
	}
	service := NewService(ds)

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.css")
	content := `.button { background: #3B82F6; }`
	if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := service.SuggestFixes(context.Background(), testFile, nil)
	if err != nil {
		t.Fatalf("SuggestFixes failed: %v", err)
	}

	// Should include original content
	if result.OriginalContent == "" {
		t.Error("expected original content to be included")
	}

	// File should not be modified
	currentContent, _ := os.ReadFile(testFile)
	if !strings.Contains(string(currentContent), "#3B82F6") {
		t.Error("SuggestFixes should not modify file")
	}
}

func TestFixFileAppliesChanges(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
		Foundations: Foundations{
			Colors: []ColorToken{
				{ID: "primary", Value: "#3B82F6"},
			},
		},
	}
	service := NewService(ds)

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.css")
	content := `.button { background: #3B82F6; }`
	if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	// Run fix without dry run
	opts := &FixOptions{DryRun: false}
	_, err := service.FixFile(context.Background(), testFile, opts)
	if err != nil {
		t.Fatalf("FixFile failed: %v", err)
	}

	// File should be modified
	currentContent, _ := os.ReadFile(testFile)
	if strings.Contains(string(currentContent), "#3B82F6") {
		t.Error("FixFile should have replaced hardcoded color")
	}
	if !strings.Contains(string(currentContent), "var(--color-primary)") {
		t.Error("FixFile should have added CSS variable")
	}
}

func TestFindClosestColor(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
		Foundations: Foundations{
			Colors: []ColorToken{
				{ID: "blue-500", Value: "#3B82F6"},
				{ID: "red-500", Value: "#EF4444"},
			},
		},
	}
	service := NewService(ds)

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.css")
	// Use a color very close to blue-500
	content := `.button { background: #3B82F5; }` // Off by 1
	if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	opts := &FixOptions{DryRun: true}
	result, err := service.FixFile(context.Background(), testFile, opts)
	if err != nil {
		t.Fatalf("FixFile failed: %v", err)
	}

	if len(result.Fixes) != 1 {
		t.Errorf("expected 1 fix for close color match, got %d", len(result.Fixes))
	}

	if !strings.Contains(result.FixedContent, "var(--color-blue-500)") {
		t.Errorf("expected closest color match, got: %s", result.FixedContent)
	}
}

func TestNormalizeHex(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"#FFF", "#ffffff"},
		{"#fff", "#ffffff"},
		{"#3B82F6", "#3b82f6"},
		{"#3b82f6", "#3b82f6"},
		{"invalid", ""},
		{"", ""},
	}

	for _, tc := range tests {
		result := normalizeHex(tc.input)
		if result != tc.expected {
			t.Errorf("normalizeHex(%q) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}

func TestToSentenceCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"IconTrash", "Icon trash"},
		{"icon-delete", "Icon delete"},
		{"close_button", "Close button"},
		{"save", "Save"},
	}

	for _, tc := range tests {
		result := toSentenceCase(tc.input)
		if result != tc.expected {
			t.Errorf("toSentenceCase(%q) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}
