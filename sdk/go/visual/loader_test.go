package visual

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoaderLoadYAML(t *testing.T) {
	// Create temp directory with test file
	dir := t.TempDir()

	yamlContent := `
version: "1.0"
name: "Test Suite"
baseUrl: "http://localhost:6006"
defaults:
  viewports:
    - name: desktop
      width: 1280
      height: 800
  threshold: 0.002
  stabilization:
    waitMs: 100
tests:
  - id: button-primary
    component: button
    variant: primary
    url: "/button--primary"
  - id: button-secondary
    component: button
    variant: secondary
    url: "/button--secondary"
    threshold: 0.005
`
	//nolint:gosec // G306: Test fixture file
	if err := os.WriteFile(filepath.Join(dir, "visual-tests.yaml"), []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	loader := NewLoader(dir)
	suite, err := loader.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify suite
	if suite.Version != "1.0" {
		t.Errorf("Expected version 1.0, got %s", suite.Version)
	}
	if suite.Name != "Test Suite" {
		t.Errorf("Expected name 'Test Suite', got %s", suite.Name)
	}
	if len(suite.Tests) != 2 {
		t.Errorf("Expected 2 tests, got %d", len(suite.Tests))
	}

	// Verify defaults applied to first test
	test1 := suite.Tests[0]
	if test1.ID != "button-primary" {
		t.Errorf("Expected ID button-primary, got %s", test1.ID)
	}
	if len(test1.Viewports) != 1 {
		t.Errorf("Expected 1 viewport, got %d", len(test1.Viewports))
	}
	if test1.Threshold != 0.002 {
		t.Errorf("Expected threshold 0.002, got %f", test1.Threshold)
	}
	if test1.URL != "http://localhost:6006/button--primary" {
		t.Errorf("Expected expanded URL, got %s", test1.URL)
	}

	// Verify custom threshold on second test
	test2 := suite.Tests[1]
	if test2.Threshold != 0.005 {
		t.Errorf("Expected threshold 0.005, got %f", test2.Threshold)
	}
}

func TestLoaderLoadJSON(t *testing.T) {
	dir := t.TempDir()

	jsonContent := `{
  "version": "1.0",
  "name": "JSON Test Suite",
  "tests": [
    {
      "id": "card-default",
      "component": "card",
      "url": "http://localhost:6006/card"
    }
  ]
}`
	//nolint:gosec // G306: Test fixture file
	if err := os.WriteFile(filepath.Join(dir, "visual-tests.json"), []byte(jsonContent), 0644); err != nil {
		t.Fatal(err)
	}

	loader := NewLoader(dir)
	suite, err := loader.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if suite.Name != "JSON Test Suite" {
		t.Errorf("Expected name 'JSON Test Suite', got %s", suite.Name)
	}

	// Verify built-in defaults applied
	test := suite.Tests[0]
	if len(test.Viewports) != 3 { // DefaultViewports() returns 3
		t.Errorf("Expected 3 default viewports, got %d", len(test.Viewports))
	}
	if test.Threshold != DefaultThreshold() {
		t.Errorf("Expected default threshold, got %f", test.Threshold)
	}
	if test.Stabilization == nil {
		t.Error("Expected default stabilization")
	}
}

func TestLoaderValidation(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr string
	}{
		{
			name: "missing version",
			content: `name: Test
tests:
  - id: x
    component: y
    url: z`,
			wantErr: "version is required",
		},
		{
			name: "missing name",
			content: `version: "1.0"
tests:
  - id: x
    component: y
    url: z`,
			wantErr: "name is required",
		},
		{
			name: "no tests",
			content: `version: "1.0"
name: Test
tests: []`,
			wantErr: "at least one test is required",
		},
		{
			name: "missing test id",
			content: `version: "1.0"
name: Test
tests:
  - component: x
    url: y`,
			wantErr: "id is required",
		},
		{
			name: "missing component",
			content: `version: "1.0"
name: Test
tests:
  - id: x
    url: y`,
			wantErr: "component is required",
		},
		{
			name: "missing url",
			content: `version: "1.0"
name: Test
tests:
  - id: x
    component: y`,
			wantErr: "url is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			//nolint:gosec // G306: Test fixture file
			if err := os.WriteFile(filepath.Join(dir, "visual-tests.yaml"), []byte(tc.content), 0644); err != nil {
				t.Fatal(err)
			}

			loader := NewLoader(dir)
			_, err := loader.Load()
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if tc.wantErr != "" && !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("Expected error containing %q, got %q", tc.wantErr, err.Error())
			}
		})
	}
}

func TestLoaderNoFile(t *testing.T) {
	dir := t.TempDir()
	loader := NewLoader(dir)
	_, err := loader.Load()
	if err == nil {
		t.Fatal("Expected error for missing file")
	}
}

func TestFilterTests(t *testing.T) {
	tests := []VisualTest{
		{ID: "a", Component: "button"},
		{ID: "b", Component: "card"},
		{ID: "c", Component: "button"},
	}

	// Filter by IDs
	filtered := FilterTests(tests, []string{"a", "c"})
	if len(filtered) != 2 {
		t.Errorf("Expected 2 filtered tests, got %d", len(filtered))
	}

	// Empty filter returns all
	all := FilterTests(tests, nil)
	if len(all) != 3 {
		t.Errorf("Expected 3 tests with empty filter, got %d", len(all))
	}
}

func TestFilterByComponent(t *testing.T) {
	tests := []VisualTest{
		{ID: "a", Component: "button"},
		{ID: "b", Component: "card"},
		{ID: "c", Component: "button"},
	}

	buttons := FilterByComponent(tests, "button")
	if len(buttons) != 2 {
		t.Errorf("Expected 2 button tests, got %d", len(buttons))
	}
}

func TestExpandTests(t *testing.T) {
	tests := []VisualTest{
		{
			ID:        "a",
			Component: "button",
			Viewports: []Viewport{
				{Name: "desktop", Width: 1280, Height: 800},
				{Name: "mobile", Width: 375, Height: 667},
			},
		},
		{
			ID:        "b",
			Component: "card",
			Viewports: []Viewport{
				{Name: "desktop", Width: 1280, Height: 800},
			},
		},
	}

	expanded := ExpandTests(tests)
	if len(expanded) != 3 { // 2 viewports for a + 1 for b
		t.Errorf("Expected 3 expanded tests, got %d", len(expanded))
	}

	// Verify first expansion
	if expanded[0].Test.ID != "a" || expanded[0].Viewport.Name != "desktop" {
		t.Errorf("Unexpected first expansion: %+v", expanded[0])
	}
}

func TestIsAbsoluteURL(t *testing.T) {
	tests := []struct {
		url  string
		want bool
	}{
		{"http://localhost:6006", true},
		{"https://example.com", true},
		{"//cdn.example.com", true},
		{"/path/to/page", false},
		{"relative/path", false},
		{"?query=param", false},
	}

	for _, tc := range tests {
		got := isAbsoluteURL(tc.url)
		if got != tc.want {
			t.Errorf("isAbsoluteURL(%q) = %v, want %v", tc.url, got, tc.want)
		}
	}
}
