package dss

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// FixResult contains the results of a fix operation.
type FixResult struct {
	// File is the path to the file
	File string `json:"file"`

	// OriginalContent is the original file content
	OriginalContent string `json:"originalContent,omitempty"`

	// FixedContent is the fixed file content
	FixedContent string `json:"fixedContent"`

	// Fixes contains all applied fixes
	Fixes []Fix `json:"fixes"`

	// Summary provides aggregate counts
	Summary FixSummary `json:"summary"`
}

// Fix represents a single fix applied to the code.
type Fix struct {
	// Line is the line number (1-indexed)
	Line int `json:"line"`

	// Rule is the rule that triggered this fix
	Rule string `json:"rule"`

	// Original is the original code
	Original string `json:"original"`

	// Replacement is the fixed code
	Replacement string `json:"replacement"`

	// Description explains what was fixed
	Description string `json:"description"`
}

// FixSummary provides aggregate counts.
type FixSummary struct {
	// TotalFixes is the number of fixes applied
	TotalFixes int `json:"totalFixes"`

	// ColorFixes is the number of color fixes
	ColorFixes int `json:"colorFixes"`

	// SpacingFixes is the number of spacing fixes
	SpacingFixes int `json:"spacingFixes"`

	// AccessibilityFixes is the number of accessibility fixes
	AccessibilityFixes int `json:"accessibilityFixes"`
}

// FixOptions configures fix behavior.
type FixOptions struct {
	// Rules limits fixes to specific rules (empty = all)
	Rules []string `json:"rules,omitempty"`

	// DryRun returns fixes without applying them
	DryRun bool `json:"dryRun,omitempty"`

	// IncludeOriginal includes original content in result
	IncludeOriginal bool `json:"includeOriginal,omitempty"`
}

// DefaultFixOptions returns sensible defaults.
func DefaultFixOptions() *FixOptions {
	return &FixOptions{
		DryRun:          false,
		IncludeOriginal: false,
	}
}

// FixFile fixes violations in a single file.
func (s *Service) FixFile(ctx context.Context, path string, opts *FixOptions) (*FixResult, error) {
	if opts == nil {
		opts = DefaultFixOptions()
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	result := &FixResult{
		File:  path,
		Fixes: []Fix{},
	}

	if opts.IncludeOriginal {
		result.OriginalContent = string(content)
	}

	code := string(content)
	f := &fileFixer{
		service: s,
		path:    path,
		code:    code,
		opts:    opts,
		result:  result,
	}

	fixedCode := f.fix(ctx)
	result.FixedContent = fixedCode

	// Write if not dry run and there were changes
	if !opts.DryRun && len(result.Fixes) > 0 {
		if err := os.WriteFile(path, []byte(fixedCode), 0600); err != nil {
			return nil, fmt.Errorf("writing file: %w", err)
		}
	}

	return result, nil
}

// SuggestFixes returns suggested fixes without applying them.
func (s *Service) SuggestFixes(ctx context.Context, path string, opts *FixOptions) (*FixResult, error) {
	if opts == nil {
		opts = DefaultFixOptions()
	}
	opts.DryRun = true
	opts.IncludeOriginal = true
	return s.FixFile(ctx, path, opts)
}

// fileFixer handles fixing of a single file.
type fileFixer struct {
	service *Service
	path    string
	code    string
	opts    *FixOptions
	result  *FixResult
}

func (f *fileFixer) fix(_ context.Context) string {
	rules := f.opts.Rules
	allRules := len(rules) == 0

	code := f.code

	if allRules || f.hasRule(rules, "no-hardcoded-colors") {
		code = f.fixHardcodedColors(code)
	}

	if allRules || f.hasRule(rules, "use-spacing-scale") {
		code = f.fixHardcodedSpacing(code)
	}

	if allRules || f.hasRule(rules, "img-alt-required") {
		code = f.fixMissingAlt(code)
	}

	if allRules || f.hasRule(rules, "button-accessible-name") {
		code = f.fixMissingAriaLabel(code)
	}

	f.updateSummary()
	return code
}

func (f *fileFixer) hasRule(rules []string, wanted ...string) bool {
	for _, r := range rules {
		for _, w := range wanted {
			if r == w {
				return true
			}
		}
	}
	return false
}

func (f *fileFixer) lineAt(code string, offset int) int {
	return strings.Count(code[:offset], "\n") + 1
}

func (f *fileFixer) fixHardcodedColors(code string) string {
	// Build color mapping from design system
	colorMap := f.buildColorMap()

	// Fix hex colors (6-char or 3-char hex codes)
	hexPattern := regexp.MustCompile(`#([0-9a-fA-F]{6}|[0-9a-fA-F]{3})\b`)

	code = hexPattern.ReplaceAllStringFunc(code, func(match string) string {
		// Skip if in comment
		idx := strings.LastIndex(code, match)
		if idx > 0 {
			before := ""
			if idx > 50 {
				before = code[idx-50 : idx]
			} else {
				before = code[:idx]
			}
			if strings.Contains(before, "//") || strings.Contains(before, "/*") {
				return match
			}
		}

		// Normalize to 6-char hex
		normalized := normalizeHex(match)

		// Look for matching token
		if varName, ok := colorMap[normalized]; ok {
			f.result.Fixes = append(f.result.Fixes, Fix{
				Line:        f.lineAt(code, strings.Index(code, match)),
				Rule:        "no-hardcoded-colors",
				Original:    match,
				Replacement: fmt.Sprintf("var(--%s)", varName),
				Description: fmt.Sprintf("Replaced hardcoded color with design token"),
			})
			return fmt.Sprintf("var(--%s)", varName)
		}

		// Find closest color match
		if varName := f.findClosestColor(normalized); varName != "" {
			f.result.Fixes = append(f.result.Fixes, Fix{
				Line:        f.lineAt(code, strings.Index(code, match)),
				Rule:        "no-hardcoded-colors",
				Original:    match,
				Replacement: fmt.Sprintf("var(--%s)", varName),
				Description: fmt.Sprintf("Replaced with closest matching design token"),
			})
			return fmt.Sprintf("var(--%s)", varName)
		}

		return match
	})

	return code
}

func (f *fileFixer) buildColorMap() map[string]string {
	colorMap := make(map[string]string)

	for _, color := range f.service.ds.Foundations.Colors {
		normalized := normalizeHex(color.Value)
		if normalized != "" {
			varName := "color-" + color.ID
			colorMap[normalized] = varName
		}
	}

	return colorMap
}

func (f *fileFixer) findClosestColor(hex string) string {
	if len(f.service.ds.Foundations.Colors) == 0 {
		return ""
	}

	r1, g1, b1 := hexToRGB(hex)

	var closest string
	minDist := 999999.0

	for _, color := range f.service.ds.Foundations.Colors {
		normalized := normalizeHex(color.Value)
		if normalized == "" {
			continue
		}

		r2, g2, b2 := hexToRGB(normalized)
		dist := colorDistance(r1, g1, b1, r2, g2, b2)

		if dist < minDist {
			minDist = dist
			closest = "color-" + color.ID
		}
	}

	// Only suggest if reasonably close (threshold of 50)
	if minDist < 50 {
		return closest
	}

	return ""
}

func (f *fileFixer) fixHardcodedSpacing(code string) string {
	// Build spacing mapping
	spacingMap := f.buildSpacingMap()

	// Match patterns like: padding: 15px, margin: 10px, gap: 20px
	pxPattern := regexp.MustCompile(`(padding|margin|gap|top|right|bottom|left|width|height):\s*(\d+)px`)

	code = pxPattern.ReplaceAllStringFunc(code, func(match string) string {
		parts := pxPattern.FindStringSubmatch(match)
		if len(parts) < 3 {
			return match
		}

		prop := parts[1]
		value := parts[2]

		// Look for exact match first
		if varName, ok := spacingMap[value]; ok {
			f.result.Fixes = append(f.result.Fixes, Fix{
				Line:        f.lineAt(code, strings.Index(code, match)),
				Rule:        "use-spacing-scale",
				Original:    match,
				Replacement: fmt.Sprintf("%s: var(--%s)", prop, varName),
				Description: fmt.Sprintf("Replaced hardcoded spacing with design token"),
			})
			return fmt.Sprintf("%s: var(--%s)", prop, varName)
		}

		// Find closest spacing value
		if varName := f.findClosestSpacing(value); varName != "" {
			f.result.Fixes = append(f.result.Fixes, Fix{
				Line:        f.lineAt(code, strings.Index(code, match)),
				Rule:        "use-spacing-scale",
				Original:    match,
				Replacement: fmt.Sprintf("%s: var(--%s)", prop, varName),
				Description: fmt.Sprintf("Replaced with closest spacing token"),
			})
			return fmt.Sprintf("%s: var(--%s)", prop, varName)
		}

		return match
	})

	return code
}

func (f *fileFixer) buildSpacingMap() map[string]string {
	spacingMap := make(map[string]string)

	if f.service.ds.Foundations.Spacing == nil {
		// Default spacing scale
		defaults := map[string]string{
			"0": "spacing-0", "4": "spacing-1", "8": "spacing-2",
			"12": "spacing-3", "16": "spacing-4", "20": "spacing-5",
			"24": "spacing-6", "32": "spacing-8", "40": "spacing-10",
			"48": "spacing-12", "64": "spacing-16",
		}
		return defaults
	}

	for _, sp := range f.service.ds.Foundations.Spacing.Scale {
		var pxValue string
		if sp.PixelValue > 0 {
			pxValue = fmt.Sprintf("%d", sp.PixelValue)
		} else {
			// Try to extract from value like "16px" or "1rem"
			pxValue = strings.TrimSuffix(sp.Value, "px")
		}
		if pxValue != "" {
			spacingMap[pxValue] = "spacing-" + sp.ID
		}
	}

	return spacingMap
}

func (f *fileFixer) findClosestSpacing(value string) string {
	px, err := strconv.Atoi(value)
	if err != nil {
		return ""
	}

	// Standard spacing scale
	scale := []int{0, 4, 8, 12, 16, 20, 24, 32, 40, 48, 64, 80, 96}
	scaleNames := []string{
		"spacing-0", "spacing-1", "spacing-2", "spacing-3", "spacing-4",
		"spacing-5", "spacing-6", "spacing-8", "spacing-10", "spacing-12",
		"spacing-16", "spacing-20", "spacing-24",
	}

	// Add from design system
	if f.service.ds.Foundations.Spacing != nil {
		for _, sp := range f.service.ds.Foundations.Spacing.Scale {
			if sp.PixelValue > 0 {
				scale = append(scale, sp.PixelValue)
				scaleNames = append(scaleNames, "spacing-"+sp.ID)
			}
		}
	}

	// Find closest
	minDiff := 999
	closestIdx := -1

	for i, s := range scale {
		diff := abs(px - s)
		if diff < minDiff {
			minDiff = diff
			closestIdx = i
		}
	}

	// Only suggest if within 4px
	if closestIdx >= 0 && minDiff <= 4 && closestIdx < len(scaleNames) {
		return scaleNames[closestIdx]
	}

	return ""
}

func (f *fileFixer) fixMissingAlt(code string) string {
	// Match img tags without alt attribute
	imgPattern := regexp.MustCompile(`<img([^>]*)>`)

	code = imgPattern.ReplaceAllStringFunc(code, func(match string) string {
		if strings.Contains(match, "alt=") || strings.Contains(match, "alt =") {
			return match
		}

		// Try to extract src to generate meaningful alt text
		srcPattern := regexp.MustCompile(`src=["']([^"']+)["']`)
		srcMatch := srcPattern.FindStringSubmatch(match)

		altText := ""
		if len(srcMatch) > 1 {
			// Extract filename without extension as alt text base
			src := srcMatch[1]
			parts := strings.Split(src, "/")
			filename := parts[len(parts)-1]
			filename = strings.TrimSuffix(filename, ".png")
			filename = strings.TrimSuffix(filename, ".jpg")
			filename = strings.TrimSuffix(filename, ".jpeg")
			filename = strings.TrimSuffix(filename, ".svg")
			filename = strings.TrimSuffix(filename, ".gif")
			filename = strings.TrimSuffix(filename, ".webp")
			altText = strings.ReplaceAll(filename, "-", " ")
			altText = strings.ReplaceAll(altText, "_", " ")
		}

		if altText == "" {
			altText = "TODO: Add descriptive alt text"
		}

		// Insert alt attribute before closing >
		replacement := strings.TrimSuffix(match, ">")
		replacement = strings.TrimSuffix(replacement, " ")
		replacement = replacement + fmt.Sprintf(` alt="%s">`, altText)

		f.result.Fixes = append(f.result.Fixes, Fix{
			Line:        f.lineAt(code, strings.Index(code, match)),
			Rule:        "img-alt-required",
			Original:    match,
			Replacement: replacement,
			Description: "Added alt attribute to image",
		})

		return replacement
	})

	return code
}

func (f *fileFixer) fixMissingAriaLabel(code string) string {
	// Match icon-only buttons without aria-label
	buttonPattern := regexp.MustCompile(`<[Bb]utton([^>]*size=["']icon["'][^>]*)>`)

	code = buttonPattern.ReplaceAllStringFunc(code, func(match string) string {
		if strings.Contains(match, "aria-label") {
			return match
		}

		// Try to infer label from icon name or other context
		iconPattern := regexp.MustCompile(`(?:icon|Icon)=["']{?([^"'}>]+)`)
		iconMatch := iconPattern.FindStringSubmatch(match)

		ariaLabel := "TODO: Add button description"
		if len(iconMatch) > 1 {
			// Convert icon name to readable label
			iconName := iconMatch[1]
			iconName = strings.TrimPrefix(iconName, "Icon")
			iconName = strings.TrimPrefix(iconName, "icon")
			// Convert camelCase or kebab-case to sentence
			ariaLabel = toSentenceCase(iconName)
		}

		// Insert aria-label before closing >
		replacement := strings.TrimSuffix(match, ">")
		replacement = strings.TrimSuffix(replacement, " ")
		replacement = replacement + fmt.Sprintf(` aria-label="%s">`, ariaLabel)

		f.result.Fixes = append(f.result.Fixes, Fix{
			Line:        f.lineAt(code, strings.Index(code, match)),
			Rule:        "button-accessible-name",
			Original:    match,
			Replacement: replacement,
			Description: "Added aria-label to icon button",
		})

		return replacement
	})

	return code
}

func (f *fileFixer) updateSummary() {
	for _, fix := range f.result.Fixes {
		f.result.Summary.TotalFixes++
		switch fix.Rule {
		case "no-hardcoded-colors":
			f.result.Summary.ColorFixes++
		case "use-spacing-scale":
			f.result.Summary.SpacingFixes++
		case "img-alt-required", "button-accessible-name":
			f.result.Summary.AccessibilityFixes++
		}
	}
}

// Helper functions

func normalizeHex(color string) string {
	color = strings.TrimPrefix(color, "#")
	color = strings.ToLower(color)

	// Handle 3-char hex
	if len(color) == 3 {
		color = string(color[0]) + string(color[0]) +
			string(color[1]) + string(color[1]) +
			string(color[2]) + string(color[2])
	}

	// Only return valid 6-char hex
	if len(color) == 6 {
		return "#" + color
	}

	// Try to extract hex from hsl/rgb
	if strings.HasPrefix(strings.ToLower(color), "hsl") || strings.HasPrefix(strings.ToLower(color), "rgb") {
		// Skip complex color functions for now
		return ""
	}

	return ""
}

func hexToRGB(hex string) (int, int, int) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 0, 0, 0
	}

	r, _ := strconv.ParseInt(hex[0:2], 16, 64)
	g, _ := strconv.ParseInt(hex[2:4], 16, 64)
	b, _ := strconv.ParseInt(hex[4:6], 16, 64)

	return int(r), int(g), int(b)
}

func colorDistance(r1, g1, b1, r2, g2, b2 int) float64 {
	// Simple Euclidean distance
	dr := float64(r1 - r2)
	dg := float64(g1 - g2)
	db := float64(b1 - b2)
	return dr*dr + dg*dg + db*db
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func toSentenceCase(s string) string {
	// Handle kebab-case
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")

	// Handle camelCase
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune(' ')
		}
		result.WriteRune(r)
	}

	s = result.String()
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)

	// Capitalize first letter
	if len(s) > 0 {
		s = strings.ToUpper(string(s[0])) + s[1:]
	}

	return s
}

// FixDirectory fixes all files in a directory.
func (s *Service) FixDirectory(ctx context.Context, dir string, opts *FixOptions) ([]*FixResult, error) {
	if opts == nil {
		opts = DefaultFixOptions()
	}

	var results []*FixResult

	extensions := map[string]bool{
		".tsx": true, ".jsx": true, ".ts": true, ".js": true, ".css": true,
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading directory: %w", err)
	}

	// Sort for consistent ordering
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Check if file has a supported extension
		for ext := range extensions {
			if strings.HasSuffix(strings.ToLower(entry.Name()), ext) {
				path := dir + "/" + entry.Name()
				result, err := s.FixFile(ctx, path, opts)
				if err != nil {
					continue
				}
				if len(result.Fixes) > 0 {
					results = append(results, result)
				}
				break
			}
		}
	}

	return results, nil
}
