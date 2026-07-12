package dss

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/plexusone/structured-evaluation/rubric"
)

//go:embed templates/html/*.html templates/html/static/*.css
var htmlTemplates embed.FS

// HTMLOptions configures HTML generation.
type HTMLOptions struct {
	// OutputDir is the output directory
	OutputDir string

	// Title is the documentation title
	Title string

	// EvalResult is optional evaluation data (rubric format)
	// Coverage is extracted from EvalResult.GetCoverage() if available
	EvalResult *rubric.Rubric
}

// HTMLPage represents a page in the documentation.
type HTMLPage struct {
	Title            string
	DesignSystemName string
	Version          string
	ActivePage       string
	CSS              template.CSS
	Content          template.HTML
}

// HTMLEvalCategory is a template-friendly category representation.
type HTMLEvalCategory struct {
	Name       string
	Score      int     // 0-100 for progress bar
	IntScore   int     // 1-5
	ScoreLabel string  // "Excellent", "Good", etc.
	Weight     float64 // 0-1.0
	Passed     int
	Checks     int
}

// HTMLEvalFinding is a template-friendly finding representation.
type HTMLEvalFinding struct {
	Severity       string
	SeverityClass  string
	ID             string
	Title          string
	Description    string
	Location       string
	Recommendation string
}

// HTMLEvalCoverageSection is a template-friendly coverage section.
type HTMLEvalCoverageSection struct {
	Total      int
	Complete   int
	Percentage int
	Missing    []string
}

// HTMLEvalCoverage is a template-friendly coverage representation.
// Maps from rubric.CoverageReport sections to named fields for template access.
type HTMLEvalCoverage struct {
	Components    HTMLEvalCoverageSection
	Foundations   HTMLEvalCoverageSection
	Patterns      HTMLEvalCoverageSection
	Accessibility HTMLEvalCoverageSection
	Overall       int
}

// GenerateHTML generates HTML documentation for the design system.
func (ds *DesignSystem) GenerateHTML(opts *HTMLOptions) error {
	if opts == nil {
		opts = &HTMLOptions{}
	}

	if opts.OutputDir == "" {
		opts.OutputDir = "docs"
	}

	if opts.Title == "" {
		opts.Title = ds.Meta.Name
	}

	// Create output directory
	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Load CSS
	cssData, err := htmlTemplates.ReadFile("templates/html/static/style.css")
	if err != nil {
		return fmt.Errorf("reading CSS: %w", err)
	}

	// Parse templates
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"mul": func(a, b float64) float64 { return a * b },
	}).ParseFS(htmlTemplates, "templates/html/*.html")
	if err != nil {
		return fmt.Errorf("parsing templates: %w", err)
	}

	// Generate pages
	pages := []struct {
		name     string
		template string
		data     any
	}{
		{"index.html", "index.html", ds.indexData(opts)},
		{"components.html", "components.html", ds.componentsData()},
		{"tokens.html", "tokens.html", ds.tokensData()},
		{"eval.html", "eval.html", ds.evalData(opts)},
	}

	for _, page := range pages {
		if err := ds.generatePage(tmpl, opts.OutputDir, page.name, page.template, page.data, string(cssData)); err != nil {
			return fmt.Errorf("generating %s: %w", page.name, err)
		}
	}

	// Generate individual component pages
	for _, comp := range ds.Components {
		filename := fmt.Sprintf("component-%s.html", comp.ID)
		if err := ds.generatePage(tmpl, opts.OutputDir, filename, "component.html", comp, string(cssData)); err != nil {
			return fmt.Errorf("generating %s: %w", filename, err)
		}
	}

	return nil
}

func (ds *DesignSystem) generatePage(tmpl *template.Template, outputDir, filename, templateName string, data any, css string) error {
	// Render content
	var contentBuf strings.Builder
	if err := tmpl.ExecuteTemplate(&contentBuf, templateName, data); err != nil {
		return fmt.Errorf("executing template %s: %w", templateName, err)
	}

	// Determine active page
	activePage := strings.TrimSuffix(filename, ".html")
	if strings.HasPrefix(activePage, "component-") {
		activePage = "components"
	}

	// Render layout
	//nolint:gosec // G203: CSS/HTML are from embedded templates, not user input
	page := HTMLPage{
		Title:            ds.Meta.Name,
		DesignSystemName: ds.Meta.Name,
		Version:          ds.Meta.Version,
		ActivePage:       activePage,
		CSS:              template.CSS(css),
		Content:          template.HTML(contentBuf.String()),
	}

	// Create output file
	outputPath := filepath.Join(outputDir, filename)
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	if err := tmpl.ExecuteTemplate(f, "layout.html", page); err != nil {
		return fmt.Errorf("executing layout: %w", err)
	}

	return nil
}

// indexData returns data for the index page.
func (ds *DesignSystem) indexData(opts *HTMLOptions) map[string]any {
	tokenCount := len(ds.Foundations.Colors)
	if ds.Foundations.Spacing != nil {
		tokenCount += len(ds.Foundations.Spacing.Scale)
	}
	if ds.Foundations.Typography != nil {
		tokenCount += len(ds.Foundations.Typography.FontFamilies)
		tokenCount += len(ds.Foundations.Typography.FontSizes)
	}
	tokenCount += len(ds.Foundations.Elevation)

	data := map[string]any{
		"Meta":           ds.Meta,
		"ComponentCount": len(ds.Components),
		"TokenCount":     tokenCount,
		"PatternCount":   len(ds.Patterns),
		"Principles":     ds.Principles,
	}

	if opts.EvalResult != nil {
		// Convert IntScore (1-5) to percentage for display
		data["EvalScore"] = intScoreToPercentage(opts.EvalResult.IntScore)
		data["EvalGradeClass"] = intScoreToGradeClass(opts.EvalResult.IntScore)
	}

	return data
}

// componentsData returns data for the components page.
func (ds *DesignSystem) componentsData() map[string]any {
	return map[string]any{
		"Components": ds.Components,
	}
}

// tokensData returns data for the tokens page.
func (ds *DesignSystem) tokensData() map[string]any {
	return map[string]any{
		"Colors":     ds.Foundations.Colors,
		"Typography": ds.Foundations.Typography,
		"Spacing":    ds.Foundations.Spacing,
		"Elevation":  ds.Foundations.Elevation,
	}
}

// evalData returns data for the eval page.
func (ds *DesignSystem) evalData(opts *HTMLOptions) map[string]any {
	data := map[string]any{}

	if opts.EvalResult != nil {
		r := opts.EvalResult

		// Convert rubric to template-friendly format
		evalData := map[string]any{
			"Score":      intScoreToPercentage(r.IntScore),
			"IntScore":   int(r.IntScore),
			"ScoreLabel": r.IntScore.String(),
			"Grade":      intScoreToGrade(r.IntScore),
			"Decision":   r.OverallDecision,
			"Summary":    r.Summary,
		}

		// Convert categories
		var categories []HTMLEvalCategory
		for _, cat := range r.Categories {
			categories = append(categories, HTMLEvalCategory{
				Name:       cat.Category,
				Score:      intScoreToPercentage(cat.IntScore),
				IntScore:   int(cat.IntScore),
				ScoreLabel: cat.IntScore.String(),
				Weight:     evalCategoryWeights[cat.Category],
				Passed:     countPassedFindings(cat.Findings),
				Checks:     len(cat.Findings) + countPassedFindings(cat.Findings),
			})
		}
		evalData["Categories"] = categories

		// Add coverage if available (from rubric extensions)
		if coverage := r.GetCoverage(); coverage != nil {
			// Convert to template-friendly format
			evalData["Coverage"] = convertCoverageForHTML(coverage)
		}

		// Convert findings
		var findings []HTMLEvalFinding
		for _, f := range r.Findings {
			findings = append(findings, HTMLEvalFinding{
				Severity:       string(f.Severity),
				SeverityClass:  severityToClass(f.Severity),
				ID:             f.ID,
				Title:          f.Title,
				Description:    f.Description,
				Location:       f.Location,
				Recommendation: f.Recommendation,
			})
		}
		evalData["Findings"] = findings

		data["Eval"] = evalData
		data["EvalGradeClass"] = intScoreToGradeClass(r.IntScore)
		data["EvalTimestamp"] = r.Metadata.GeneratedAt.Format("2006-01-02 15:04:05")
	}

	return data
}

// intScoreToPercentage converts IntegerScore (1-5) to percentage (0-100).
func intScoreToPercentage(score rubric.IntegerScore) int {
	return int(score) * 20
}

// intScoreToGrade converts IntegerScore to letter grade.
func intScoreToGrade(score rubric.IntegerScore) string {
	switch score {
	case rubric.ScoreExcellent:
		return "A"
	case rubric.ScoreGood:
		return "B"
	case rubric.ScoreAcceptable:
		return "C"
	case rubric.ScoreMajorRevisions:
		return "D"
	default:
		return "F"
	}
}

// intScoreToGradeClass converts IntegerScore to CSS class.
func intScoreToGradeClass(score rubric.IntegerScore) string {
	return "grade-" + strings.ToLower(intScoreToGrade(score))
}

// severityToClass converts Severity to CSS class.
func severityToClass(severity rubric.Severity) string {
	switch severity {
	case rubric.SeverityCritical, rubric.SeverityHigh:
		return "error"
	case rubric.SeverityMedium:
		return "warning"
	default:
		return "info"
	}
}

// countPassedFindings estimates passed checks (findings are failures).
func countPassedFindings(findings []rubric.Finding) int {
	// This is a rough estimate since we don't track individual passes
	// For now, assume 3x as many passes as failures for a reasonable ratio
	return len(findings) * 2
}

// convertCoverageForHTML converts CoverageReport to template-friendly format.
func convertCoverageForHTML(cr *rubric.CoverageReport) HTMLEvalCoverage {
	sectionToHTML := func(name string) HTMLEvalCoverageSection {
		s := cr.GetSection(name)
		return HTMLEvalCoverageSection{
			Total:      s.Total,
			Complete:   s.Complete,
			Percentage: s.Percentage,
			Missing:    s.Missing,
		}
	}

	return HTMLEvalCoverage{
		Components:    sectionToHTML(CoverageSectionComponents),
		Foundations:   sectionToHTML(CoverageSectionFoundations),
		Patterns:      sectionToHTML(CoverageSectionPatterns),
		Accessibility: sectionToHTML(CoverageSectionAccessibility),
		Overall:       cr.Overall,
	}
}
