package dss

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ValidationResult contains all validation findings for a file or directory.
type ValidationResult struct {
	// Files is the number of files checked
	Files int `json:"files"`

	// Violations contains all findings
	Violations []Violation `json:"violations"`

	// Summary provides aggregate counts
	Summary ValidationSummary `json:"summary"`
}

// Violation represents a single compliance issue.
type Violation struct {
	// File is the path to the file
	File string `json:"file"`

	// Line is the line number (1-indexed, 0 if unknown)
	Line int `json:"line,omitempty"`

	// Column is the column number (1-indexed, 0 if unknown)
	Column int `json:"column,omitempty"`

	// Rule is the rule ID (e.g., "no-hardcoded-colors")
	Rule string `json:"rule"`

	// Message describes the issue
	Message string `json:"message"`

	// Severity is "error", "warning", or "info"
	Severity string `json:"severity"`

	// Component is the component ID if applicable
	Component string `json:"component,omitempty"`

	// Context is a code snippet for context
	Context string `json:"context,omitempty"`
}

// ValidationSummary provides aggregate counts.
type ValidationSummary struct {
	Errors   int `json:"errors"`
	Warnings int `json:"warnings"`
	Infos    int `json:"infos"`
	Passed   int `json:"passed"`
}

// ValidateOptions configures validation behavior.
type ValidateOptions struct {
	// Rules limits validation to specific rules (empty = all)
	Rules []string `json:"rules,omitempty"`

	// Extensions limits file types to check (default: .tsx, .jsx, .ts, .js, .css)
	Extensions []string `json:"extensions,omitempty"`

	// IncludeContext includes code snippets in violations
	IncludeContext bool `json:"includeContext,omitempty"`
}

// DefaultValidateOptions returns sensible defaults.
func DefaultValidateOptions() *ValidateOptions {
	return &ValidateOptions{
		Extensions:     []string{".tsx", ".jsx", ".ts", ".js", ".css"},
		IncludeContext: false,
	}
}

// ValidateFile checks a single file against the design system.
func (s *Service) ValidateFile(ctx context.Context, path string, opts *ValidateOptions) (*ValidationResult, error) {
	if opts == nil {
		opts = DefaultValidateOptions()
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	result := &ValidationResult{
		Files:      1,
		Violations: []Violation{},
	}

	code := string(content)
	v := &fileValidator{
		service: s,
		path:    path,
		code:    code,
		opts:    opts,
		result:  result,
	}

	v.validate(ctx)
	v.updateSummary()

	return result, nil
}

// ValidateDirectory checks all files in a directory against the design system.
func (s *Service) ValidateDirectory(ctx context.Context, dir string, opts *ValidateOptions) (*ValidationResult, error) {
	if opts == nil {
		opts = DefaultValidateOptions()
	}

	result := &ValidationResult{
		Files:      0,
		Violations: []Violation{},
	}

	extSet := make(map[string]bool)
	for _, ext := range opts.Extensions {
		extSet[ext] = true
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if !extSet[ext] {
			return nil
		}

		fileResult, err := s.ValidateFile(ctx, path, opts)
		if err != nil {
			result.Violations = append(result.Violations, Violation{
				File:     path,
				Rule:     "file-read",
				Message:  fmt.Sprintf("Could not read file: %v", err),
				Severity: "error",
			})
			return nil
		}

		result.Files++
		result.Violations = append(result.Violations, fileResult.Violations...)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking directory: %w", err)
	}

	// Update summary
	for _, v := range result.Violations {
		switch v.Severity {
		case "error":
			result.Summary.Errors++
		case "warning":
			result.Summary.Warnings++
		case "info":
			result.Summary.Infos++
		}
	}

	return result, nil
}

// fileValidator handles validation of a single file.
type fileValidator struct {
	service *Service
	path    string
	code    string
	opts    *ValidateOptions
	result  *ValidationResult
}

func (v *fileValidator) validate(_ context.Context) {
	rules := v.opts.Rules
	allRules := len(rules) == 0

	if allRules || v.hasRule(rules, "no-hardcoded-colors") {
		v.checkHardcodedColors()
	}

	if allRules || v.hasRule(rules, "use-spacing-scale") {
		v.checkHardcodedSpacing()
	}

	if allRules || v.hasRule(rules, "img-alt-required", "button-accessible-name") {
		v.checkAccessibility()
	}

	if allRules || v.hasRule(rules, "valid-variant") {
		v.checkVariants()
	}

	if allRules || v.hasRule(rules, "single-primary-button", "no-nested-cards") {
		v.checkAntiPatterns()
	}
}

func (v *fileValidator) hasRule(rules []string, wanted ...string) bool {
	for _, r := range rules {
		for _, w := range wanted {
			if r == w {
				return true
			}
		}
	}
	return false
}

func (v *fileValidator) addViolation(rule, message, severity string, line int) {
	violation := Violation{
		File:     v.path,
		Line:     line,
		Rule:     rule,
		Message:  message,
		Severity: severity,
	}

	if v.opts.IncludeContext && line > 0 {
		violation.Context = v.getLineContext(line)
	}

	v.result.Violations = append(v.result.Violations, violation)
}

func (v *fileValidator) getLineContext(line int) string {
	lines := strings.Split(v.code, "\n")
	if line <= 0 || line > len(lines) {
		return ""
	}
	return strings.TrimSpace(lines[line-1])
}

func (v *fileValidator) lineAt(offset int) int {
	return strings.Count(v.code[:offset], "\n") + 1
}

func (v *fileValidator) checkHardcodedColors() {
	// Check for hex colors
	hexPattern := regexp.MustCompile(`#[0-9a-fA-F]{3,8}`)
	matches := hexPattern.FindAllStringIndex(v.code, -1)

	for _, match := range matches {
		// Skip if in comment
		context := ""
		if match[0] > 20 {
			context = v.code[match[0]-20 : match[0]]
		}
		if strings.Contains(context, "//") || strings.Contains(context, "/*") {
			continue
		}

		colorValue := v.code[match[0]:match[1]]
		line := v.lineAt(match[0])
		v.addViolation(
			"no-hardcoded-colors",
			fmt.Sprintf("Hardcoded color '%s' - use CSS variable from design system", colorValue),
			"warning",
			line,
		)
	}

	// Check for rgb/hsl colors
	rgbPattern := regexp.MustCompile(`(?:rgb|hsl)a?\([^)]+\)`)
	rgbMatches := rgbPattern.FindAllStringIndex(v.code, -1)

	for _, match := range rgbMatches {
		// Skip if in CSS variable context
		context := ""
		if match[0] > 50 {
			context = v.code[match[0]-50 : match[0]]
		}
		if strings.Contains(context, "--color-") || strings.Contains(context, "var(") {
			continue
		}

		colorValue := v.code[match[0]:match[1]]
		line := v.lineAt(match[0])
		v.addViolation(
			"no-hardcoded-colors",
			fmt.Sprintf("Hardcoded color '%s' - use CSS variable", colorValue),
			"warning",
			line,
		)
	}
}

func (v *fileValidator) checkHardcodedSpacing() {
	pxPattern := regexp.MustCompile(`:\s*(\d+)px`)
	matches := pxPattern.FindAllStringSubmatch(v.code, -1)

	// Build valid spacing values from design system
	validSpacing := map[string]bool{
		"0": true, "1": true, "2": true, "4": true, "8": true,
		"12": true, "16": true, "20": true, "24": true, "32": true,
		"40": true, "48": true, "64": true, "80": true, "96": true,
	}

	// Add spacing from design system
	if v.service.ds.Foundations.Spacing != nil {
		for _, sp := range v.service.ds.Foundations.Spacing.Scale {
			// Extract numeric value from spacing token
			var numVal string
			_, _ = fmt.Sscanf(sp.Value, "%s", &numVal)
			numVal = strings.TrimSuffix(numVal, "px")
			numVal = strings.TrimSuffix(numVal, "rem")
			if sp.PixelValue > 0 {
				numVal = fmt.Sprintf("%d", sp.PixelValue)
			}
			validSpacing[numVal] = true
		}
	}

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		value := match[1]

		// Skip valid spacing values
		if validSpacing[value] {
			continue
		}

		// Skip small values that might be borders
		if value == "1" || value == "2" || value == "3" {
			continue
		}

		idx := strings.Index(v.code, match[0])
		line := v.lineAt(idx)
		v.addViolation(
			"use-spacing-scale",
			fmt.Sprintf("Value '%spx' not in spacing scale - use design system spacing", value),
			"warning",
			line,
		)
	}
}

func (v *fileValidator) checkAccessibility() {
	// Check for images without alt
	imgPattern := regexp.MustCompile(`<img[^>]*>`)
	imgMatches := imgPattern.FindAllString(v.code, -1)

	for _, img := range imgMatches {
		if !strings.Contains(img, "alt=") && !strings.Contains(img, "alt =") {
			idx := strings.Index(v.code, img)
			line := v.lineAt(idx)
			v.addViolation(
				"img-alt-required",
				"Image missing alt attribute",
				"error",
				line,
			)
		}
	}

	// Check for icon-only buttons without aria-label
	buttonPattern := regexp.MustCompile(`<[Bb]utton[^>]*>`)
	buttonMatches := buttonPattern.FindAllString(v.code, -1)

	for _, btn := range buttonMatches {
		if strings.Contains(btn, "size=\"icon\"") || strings.Contains(btn, "size='icon'") {
			if !strings.Contains(btn, "aria-label") {
				idx := strings.Index(v.code, btn)
				line := v.lineAt(idx)
				v.addViolation(
					"button-accessible-name",
					"Icon-only button should have aria-label",
					"warning",
					line,
				)
			}
		}
	}
}

func (v *fileValidator) checkVariants() {
	// Build component specs map
	specs := make(map[string]Component)
	for _, c := range v.service.ds.Components {
		specs[c.ID] = c
	}

	// Get component name from filename
	filename := filepath.Base(v.path)
	componentName := strings.TrimSuffix(filename, filepath.Ext(filename))

	spec, ok := specs[componentName]
	if !ok {
		return
	}

	if len(spec.Variants) == 0 {
		return
	}

	validVariants := make(map[string]bool)
	for _, variant := range spec.Variants {
		validVariants[variant.ID] = true
	}

	variantPattern := regexp.MustCompile(`variant=["']([^"']+)["']`)
	matches := variantPattern.FindAllStringSubmatch(v.code, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		variant := match[1]
		if !validVariants[variant] {
			idx := strings.Index(v.code, match[0])
			line := v.lineAt(idx)

			variantList := make([]string, 0, len(validVariants))
			for k := range validVariants {
				variantList = append(variantList, k)
			}

			v.result.Violations = append(v.result.Violations, Violation{
				File:      v.path,
				Line:      line,
				Rule:      "valid-variant",
				Message:   fmt.Sprintf("Unknown variant '%s' - valid variants: %v", variant, variantList),
				Severity:  "error",
				Component: spec.Name,
			})
		}
	}
}

func (v *fileValidator) checkAntiPatterns() {
	// Check for multiple primary buttons
	primaryBtnPattern := regexp.MustCompile(`<Button[^>]*variant=["'](?:default|primary)["']`)
	primaryMatches := primaryBtnPattern.FindAllString(v.code, -1)

	if len(primaryMatches) > 1 {
		v.addViolation(
			"single-primary-button",
			fmt.Sprintf("Found %d primary buttons - design system recommends one primary button per view", len(primaryMatches)),
			"warning",
			0,
		)
	}

	// Check for nested cards
	cardPattern := regexp.MustCompile(`<Card[^>]*>[^<]*<Card`)
	if cardPattern.MatchString(v.code) {
		v.addViolation(
			"no-nested-cards",
			"Nested cards detected - design system recommends avoiding card nesting",
			"warning",
			0,
		)
	}
}

func (v *fileValidator) updateSummary() {
	for _, violation := range v.result.Violations {
		switch violation.Severity {
		case "error":
			v.result.Summary.Errors++
		case "warning":
			v.result.Summary.Warnings++
		case "info":
			v.result.Summary.Infos++
		}
	}
}
