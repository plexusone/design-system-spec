package dss

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestServiceValidateFile(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
		Foundations: Foundations{
			Spacing: &SpacingScale{
				BaseUnit: "4px",
				Scale: []SpacingToken{
					{ID: "0", Value: "0px", PixelValue: 0},
					{ID: "4", Value: "16px", PixelValue: 16},
				},
			},
		},
	}
	service := NewService(ds)
	ctx := context.Background()

	// Create temp file with test content
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.tsx")

	t.Run("hardcoded colors", func(t *testing.T) {
		content := `
const Button = () => {
  return <button style={{ backgroundColor: "#FF0000" }}>Click</button>
}
`
		if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}

		result, err := service.ValidateFile(ctx, testFile, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Files != 1 {
			t.Errorf("expected 1 file, got %d", result.Files)
		}

		// Should have a warning for hardcoded color
		hasColorWarning := false
		for _, v := range result.Violations {
			if v.Rule == "no-hardcoded-colors" {
				hasColorWarning = true
				break
			}
		}
		if !hasColorWarning {
			t.Error("expected violation for hardcoded color")
		}
	})

	t.Run("hardcoded spacing", func(t *testing.T) {
		// CSS file content with hardcoded spacing
		content := `
.box {
  padding: 15px;
  margin: 17px;
}
`
		if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}

		result, err := service.ValidateFile(ctx, testFile, &ValidateOptions{
			Rules: []string{"use-spacing-scale"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		hasSpacingWarning := false
		for _, v := range result.Violations {
			if v.Rule == "use-spacing-scale" {
				hasSpacingWarning = true
				break
			}
		}
		if !hasSpacingWarning {
			t.Error("expected violation for non-standard spacing")
		}
	})

	t.Run("missing alt attribute", func(t *testing.T) {
		content := `
const Image = () => {
  return <img src="/photo.jpg" />
}
`
		if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}

		result, err := service.ValidateFile(ctx, testFile, &ValidateOptions{
			Rules: []string{"img-alt-required"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		hasAltError := false
		for _, v := range result.Violations {
			if v.Rule == "img-alt-required" && v.Severity == "error" {
				hasAltError = true
				break
			}
		}
		if !hasAltError {
			t.Error("expected error for missing alt attribute")
		}
	})

	t.Run("icon button without aria-label", func(t *testing.T) {
		content := `
const IconButton = () => {
  return <Button size="icon"><Icon /></Button>
}
`
		if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}

		result, err := service.ValidateFile(ctx, testFile, &ValidateOptions{
			Rules: []string{"button-accessible-name"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		hasButtonWarning := false
		for _, v := range result.Violations {
			if v.Rule == "button-accessible-name" {
				hasButtonWarning = true
				break
			}
		}
		if !hasButtonWarning {
			t.Error("expected warning for icon button without aria-label")
		}
	})

	t.Run("clean file", func(t *testing.T) {
		content := `
const Button = () => {
  return <button className="bg-primary-500 p-4">Click</button>
}
`
		if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}

		result, err := service.ValidateFile(ctx, testFile, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Summary.Errors > 0 {
			t.Errorf("expected no errors, got %d", result.Summary.Errors)
		}
	})
}

func TestServiceValidateDirectory(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
	}
	service := NewService(ds)
	ctx := context.Background()

	tmpDir := t.TempDir()

	// Create test files
	file1 := filepath.Join(tmpDir, "Button.tsx")
	file2 := filepath.Join(tmpDir, "Input.tsx")
	file3 := filepath.Join(tmpDir, "styles.css")

	content1 := `const Button = () => <button style={{ color: "#FF0000" }}>Click</button>`
	content2 := `const Input = () => <input type="text" />`
	content3 := `.btn { color: #00FF00; }`

	if err := os.WriteFile(file1, []byte(content1), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte(content2), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file3, []byte(content3), 0600); err != nil {
		t.Fatal(err)
	}

	t.Run("validate all files", func(t *testing.T) {
		result, err := service.ValidateDirectory(ctx, tmpDir, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Files != 3 {
			t.Errorf("expected 3 files, got %d", result.Files)
		}
	})

	t.Run("validate tsx only", func(t *testing.T) {
		result, err := service.ValidateDirectory(ctx, tmpDir, &ValidateOptions{
			Extensions: []string{".tsx"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Files != 2 {
			t.Errorf("expected 2 files, got %d", result.Files)
		}
	})
}

func TestValidationSummary(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
	}
	service := NewService(ds)
	ctx := context.Background()

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.tsx")

	// File with multiple issues
	content := `
const Component = () => {
  return (
    <div style={{ color: "#FF0000", padding: "15px" }}>
      <img src="/photo.jpg" />
    </div>
  )
}
`
	if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := service.ValidateFile(ctx, testFile, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have at least one error (missing alt) and warnings (color, spacing)
	if result.Summary.Errors < 1 {
		t.Errorf("expected at least 1 error, got %d", result.Summary.Errors)
	}
	if result.Summary.Warnings < 1 {
		t.Errorf("expected at least 1 warning, got %d", result.Summary.Warnings)
	}

	// Total should match
	total := result.Summary.Errors + result.Summary.Warnings + result.Summary.Infos
	if total != len(result.Violations) {
		t.Errorf("summary total (%d) doesn't match violations count (%d)", total, len(result.Violations))
	}
}

func TestValidateWithContext(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
	}
	service := NewService(ds)
	ctx := context.Background()

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.tsx")

	content := `const x = "#FF0000"`
	if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	result, err := service.ValidateFile(ctx, testFile, &ValidateOptions{
		IncludeContext: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that context is included
	for _, v := range result.Violations {
		if v.Rule == "no-hardcoded-colors" {
			if v.Context == "" {
				t.Error("expected context to be included")
			}
			break
		}
	}
}
