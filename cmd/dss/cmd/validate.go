package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	dss "github.com/plexusone/design-system-spec/sdk/go"
)

var (
	jsonOutput bool
)

// ComplianceReport holds validation results
type ComplianceReport struct {
	Passed   []string          `json:"passed"`
	Warnings []ComplianceIssue `json:"warnings"`
	Errors   []ComplianceIssue `json:"errors"`
}

// ComplianceIssue represents a single compliance violation
type ComplianceIssue struct {
	Component string `json:"component,omitempty"`
	File      string `json:"file"`
	Line      int    `json:"line,omitempty"`
	Rule      string `json:"rule"`
	Message   string `json:"message"`
	Severity  string `json:"severity"`
}

var validateCmd = &cobra.Command{
	Use:   "validate [components-dir]",
	Short: "Validate component implementations against spec",
	Long: `Validate that React/TypeScript component implementations comply with the design system spec.

Checks for:
  - Hardcoded colors (should use CSS variables)
  - Non-standard spacing values
  - Accessibility issues (missing alt, aria-label)
  - Anti-patterns defined in the spec
  - Variant values matching component specs

Examples:
  dss validate ./src/components
  dss validate -d ./design-system ./web/src/components
  dss validate --json ./src/components`,
	Args: cobra.ExactArgs(1),
	RunE: runValidate,
}

func init() {
	validateCmd.Flags().BoolVar(&jsonOutput, "json", false, "output results as JSON")

	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	componentsDir := args[0]
	dir := getSpecDir()

	// Load design system
	ds, err := dss.LoadDesignSystem(dir)
	if err != nil {
		return fmt.Errorf("loading design system: %w", err)
	}

	if !jsonOutput {
		fmt.Printf("Validating against: %s v%s\n", ds.Meta.Name, ds.Meta.Version)
		fmt.Printf("Components directory: %s\n\n", componentsDir)
	}

	report := &ComplianceReport{
		Passed:   []string{},
		Warnings: []ComplianceIssue{},
		Errors:   []ComplianceIssue{},
	}

	// Build component lookup
	componentSpecs := make(map[string]dss.Component)
	for _, c := range ds.Components {
		componentSpecs[c.ID] = c
	}

	// Scan component files
	err = filepath.Walk(componentsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".tsx") {
			return nil
		}

		validateComponentFile(path, componentSpecs, ds, report)
		return nil
	})

	if err != nil {
		return fmt.Errorf("scanning components: %w", err)
	}

	// Output report
	if jsonOutput {
		jsonData, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(jsonData))
	} else {
		printReport(report)
	}

	// Exit with error if there are errors
	if len(report.Errors) > 0 {
		return fmt.Errorf("validation found %d errors", len(report.Errors))
	}

	return nil
}

func validateComponentFile(path string, specs map[string]dss.Component, ds *dss.DesignSystem, report *ComplianceReport) {
	content, err := os.ReadFile(path)
	if err != nil {
		report.Errors = append(report.Errors, ComplianceIssue{
			File:     path,
			Rule:     "file-read",
			Message:  fmt.Sprintf("Could not read file: %v", err),
			Severity: "error",
		})
		return
	}

	code := string(content)
	filename := filepath.Base(path)
	componentName := strings.TrimSuffix(filename, ".tsx")

	// Check for hardcoded colors (should use CSS variables)
	checkHardcodedColors(path, code, report)

	// Check for hardcoded spacing (should use spacing scale)
	checkHardcodedSpacing(path, code, report)

	// Check for accessibility issues
	checkAccessibility(path, code, report)

	// Check against specific component spec if it exists
	if spec, ok := specs[componentName]; ok {
		validateAgainstSpec(path, code, spec, report)
	}

	// Check for anti-pattern violations
	checkAntiPatterns(path, code, ds, report)
}

func checkHardcodedColors(path string, code string, report *ComplianceReport) {
	// Look for hardcoded hex colors
	hexPattern := regexp.MustCompile(`#[0-9a-fA-F]{3,8}`)
	matches := hexPattern.FindAllStringIndex(code, -1)

	for _, match := range matches {
		line := strings.Count(code[:match[0]], "\n") + 1
		colorValue := code[match[0]:match[1]]

		// Ignore common exceptions (like in comments or SVG paths)
		context := ""
		if match[0] > 20 {
			context = code[match[0]-20 : match[0]]
		}
		if strings.Contains(context, "//") || strings.Contains(context, "/*") {
			continue
		}

		report.Warnings = append(report.Warnings, ComplianceIssue{
			File:     path,
			Line:     line,
			Rule:     "no-hardcoded-colors",
			Message:  fmt.Sprintf("Hardcoded color '%s' - use CSS variable from design system", colorValue),
			Severity: "warning",
		})
	}

	// Look for hardcoded rgb/hsl
	rgbPattern := regexp.MustCompile(`(?:rgb|hsl)a?\([^)]+\)`)
	rgbMatches := rgbPattern.FindAllStringIndex(code, -1)

	for _, match := range rgbMatches {
		line := strings.Count(code[:match[0]], "\n") + 1
		colorValue := code[match[0]:match[1]]

		// Skip if it's in a CSS variable definition context
		context := ""
		if match[0] > 50 {
			context = code[match[0]-50 : match[0]]
		}
		if strings.Contains(context, "--color-") || strings.Contains(context, "var(") {
			continue
		}

		report.Warnings = append(report.Warnings, ComplianceIssue{
			File:     path,
			Line:     line,
			Rule:     "no-hardcoded-colors",
			Message:  fmt.Sprintf("Hardcoded color '%s' - use CSS variable", colorValue),
			Severity: "warning",
		})
	}
}

func checkHardcodedSpacing(path string, code string, report *ComplianceReport) {
	// Look for pixel values that should use spacing scale
	pxPattern := regexp.MustCompile(`:\s*(\d+)px`)
	matches := pxPattern.FindAllStringSubmatch(code, -1)

	validSpacing := map[string]bool{
		"0": true, "1": true, "2": true, "4": true, "8": true,
		"12": true, "16": true, "20": true, "24": true, "32": true,
		"40": true, "48": true, "64": true, "80": true, "96": true,
	}

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		value := match[1]

		// Skip if it's a valid spacing value
		if validSpacing[value] {
			continue
		}

		// Skip small values that might be borders
		if value == "1" || value == "2" || value == "3" {
			continue
		}

		idx := strings.Index(code, match[0])
		line := strings.Count(code[:idx], "\n") + 1

		report.Warnings = append(report.Warnings, ComplianceIssue{
			File:     path,
			Line:     line,
			Rule:     "use-spacing-scale",
			Message:  fmt.Sprintf("Value '%spx' not in spacing scale - use design system spacing", value),
			Severity: "warning",
		})
	}
}

func checkAccessibility(path string, code string, report *ComplianceReport) {
	// Check for images without alt
	imgPattern := regexp.MustCompile(`<img[^>]*>`)
	imgMatches := imgPattern.FindAllString(code, -1)

	for _, img := range imgMatches {
		if !strings.Contains(img, "alt=") && !strings.Contains(img, "alt =") {
			idx := strings.Index(code, img)
			line := strings.Count(code[:idx], "\n") + 1

			report.Errors = append(report.Errors, ComplianceIssue{
				File:     path,
				Line:     line,
				Rule:     "img-alt-required",
				Message:  "Image missing alt attribute",
				Severity: "error",
			})
		}
	}

	// Check for buttons without accessible name
	buttonPattern := regexp.MustCompile(`<[Bb]utton[^>]*>`)
	buttonMatches := buttonPattern.FindAllString(code, -1)

	for _, btn := range buttonMatches {
		// Icon-only buttons need aria-label
		if strings.Contains(btn, "size=\"icon\"") || strings.Contains(btn, "size='icon'") {
			if !strings.Contains(btn, "aria-label") {
				idx := strings.Index(code, btn)
				line := strings.Count(code[:idx], "\n") + 1

				report.Warnings = append(report.Warnings, ComplianceIssue{
					File:     path,
					Line:     line,
					Rule:     "button-accessible-name",
					Message:  "Icon-only button should have aria-label",
					Severity: "warning",
				})
			}
		}
	}
}

func validateAgainstSpec(path string, code string, spec dss.Component, report *ComplianceReport) {
	// Check that variant values match spec
	if len(spec.Variants) > 0 {
		validVariants := make(map[string]bool)
		for _, v := range spec.Variants {
			validVariants[v.ID] = true
		}

		// Look for variant prop usage
		variantPattern := regexp.MustCompile(`variant=["']([^"']+)["']`)
		matches := variantPattern.FindAllStringSubmatch(code, -1)

		for _, match := range matches {
			if len(match) < 2 {
				continue
			}
			variant := match[1]
			if !validVariants[variant] {
				idx := strings.Index(code, match[0])
				line := strings.Count(code[:idx], "\n") + 1

				report.Errors = append(report.Errors, ComplianceIssue{
					Component: spec.Name,
					File:      path,
					Line:      line,
					Rule:      "valid-variant",
					Message:   fmt.Sprintf("Unknown variant '%s' - valid variants: %v", variant, keys(validVariants)),
					Severity:  "error",
				})
			}
		}
	}

	report.Passed = append(report.Passed, fmt.Sprintf("%s: validated against spec", spec.Name))
}

func checkAntiPatterns(path string, code string, ds *dss.DesignSystem, report *ComplianceReport) {
	// Check for multiple primary buttons (anti-pattern)
	primaryBtnPattern := regexp.MustCompile(`<Button[^>]*variant=["'](?:default|primary)["']`)
	primaryMatches := primaryBtnPattern.FindAllString(code, -1)

	if len(primaryMatches) > 1 {
		report.Warnings = append(report.Warnings, ComplianceIssue{
			File:     path,
			Rule:     "single-primary-button",
			Message:  fmt.Sprintf("Found %d primary buttons - design system recommends one primary button per view", len(primaryMatches)),
			Severity: "warning",
		})
	}

	// Check for nested cards
	cardPattern := regexp.MustCompile(`<Card[^>]*>[^<]*<Card`)
	if cardPattern.MatchString(code) {
		report.Warnings = append(report.Warnings, ComplianceIssue{
			File:     path,
			Rule:     "no-nested-cards",
			Message:  "Nested cards detected - design system recommends avoiding card nesting",
			Severity: "warning",
		})
	}
}

func keys(m map[string]bool) []string {
	result := make([]string, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

func printReport(report *ComplianceReport) {
	fmt.Print("=== Compliance Report ===\n\n")

	if len(report.Passed) > 0 {
		fmt.Println("✓ Passed:")
		for _, p := range report.Passed {
			fmt.Printf("  %s\n", p)
		}
		fmt.Println()
	}

	if len(report.Warnings) > 0 {
		fmt.Printf("⚠ Warnings (%d):\n", len(report.Warnings))
		for _, w := range report.Warnings {
			loc := w.File
			if w.Line > 0 {
				loc = fmt.Sprintf("%s:%d", w.File, w.Line)
			}
			fmt.Printf("  [%s] %s\n    %s\n", w.Rule, loc, w.Message)
		}
		fmt.Println()
	}

	if len(report.Errors) > 0 {
		fmt.Printf("✗ Errors (%d):\n", len(report.Errors))
		for _, e := range report.Errors {
			loc := e.File
			if e.Line > 0 {
				loc = fmt.Sprintf("%s:%d", e.File, e.Line)
			}
			fmt.Printf("  [%s] %s\n    %s\n", e.Rule, loc, e.Message)
		}
		fmt.Println()
	}

	// Summary
	total := len(report.Passed) + len(report.Warnings) + len(report.Errors)
	fmt.Printf("Summary: %d checks (%d passed, %d warnings, %d errors)\n",
		total, len(report.Passed), len(report.Warnings), len(report.Errors))

	if len(report.Errors) == 0 && len(report.Warnings) == 0 {
		fmt.Println("\n✓ All components comply with design system!")
	} else if len(report.Errors) == 0 {
		fmt.Println("\n⚠ Components have warnings but no blocking errors")
	} else {
		fmt.Println("\n✗ Components have compliance errors that should be fixed")
	}
}
