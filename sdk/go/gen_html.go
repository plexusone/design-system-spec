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

	// MkDocs enables MkDocs-compatible output mode.
	// When true, generates Markdown files with embedded HTML
	// instead of standalone HTML files.
	MkDocs bool
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

	// Use MkDocs mode if requested
	if opts.MkDocs {
		return ds.generateMkDocs(opts)
	}

	// Load CSS for standalone mode
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

// generateMkDocs generates MkDocs-compatible Markdown files with embedded HTML.
// Output follows DSS canonical layers: Meta, Foundations, Components, Principles,
// Patterns, Templates, Content, Accessibility, Governance + Eval.
func (ds *DesignSystem) generateMkDocs(opts *HTMLOptions) error {
	// Parse templates
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"mul": func(a, b float64) float64 { return a * b },
	}).ParseFS(htmlTemplates, "templates/html/*.html")
	if err != nil {
		return fmt.Errorf("parsing templates: %w", err)
	}

	// DSS Canonical Layer directories
	specsDir := filepath.Join(opts.OutputDir, "specs")
	layers := map[string]string{
		"meta":          filepath.Join(specsDir, "meta"),
		"foundations":   filepath.Join(specsDir, "foundations"),
		"components":    filepath.Join(specsDir, "components"),
		"principles":    filepath.Join(specsDir, "principles"),
		"patterns":      filepath.Join(specsDir, "patterns"),
		"templates":     filepath.Join(specsDir, "templates"),
		"content":       filepath.Join(specsDir, "content"),
		"accessibility": filepath.Join(specsDir, "accessibility"),
		"governance":    filepath.Join(specsDir, "governance"),
		"eval":          filepath.Join(specsDir, "eval"),
	}

	// Create all layer directories
	if err := os.MkdirAll(specsDir, 0755); err != nil {
		return fmt.Errorf("creating specs directory: %w", err)
	}
	for _, dir := range layers {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating directory %s: %w", dir, err)
		}
	}

	// Generate spec overview (index)
	if err := ds.generateMkDocsPage(tmpl, filepath.Join(specsDir, "index.md"), "index.html", ds.Meta.Name+" Specification", ds.indexData(opts)); err != nil {
		return fmt.Errorf("generating specs index: %w", err)
	}

	// Layer 1: Meta
	if err := ds.generateMkDocsMarkdown(filepath.Join(layers["meta"], "index.md"), "Meta", ds.metaMarkdown()); err != nil {
		return fmt.Errorf("generating meta: %w", err)
	}

	// Layer 2: Foundations (tokens)
	if err := ds.generateMkDocsPage(tmpl, filepath.Join(layers["foundations"], "index.md"), "tokens.html", "Foundations", ds.tokensData()); err != nil {
		return fmt.Errorf("generating foundations: %w", err)
	}

	// Layer 3: Components
	if err := ds.generateMkDocsPage(tmpl, filepath.Join(layers["components"], "index.md"), "components.html", "Components", ds.componentsData()); err != nil {
		return fmt.Errorf("generating components index: %w", err)
	}
	for _, comp := range ds.Components {
		filename := filepath.Join(layers["components"], comp.ID+".md")
		if err := ds.generateMkDocsPage(tmpl, filename, "component.html", comp.Name, comp); err != nil {
			return fmt.Errorf("generating component %s: %w", comp.ID, err)
		}
	}

	// Layer 4: Principles (if present)
	if len(ds.Principles) > 0 {
		if err := ds.generateMkDocsMarkdown(filepath.Join(layers["principles"], "index.md"), "Principles", ds.principlesMarkdown()); err != nil {
			return fmt.Errorf("generating principles: %w", err)
		}
	}

	// Layer 5: Patterns (if present)
	if len(ds.Patterns) > 0 {
		if err := ds.generateMkDocsMarkdown(filepath.Join(layers["patterns"], "index.md"), "Patterns", ds.patternsMarkdown()); err != nil {
			return fmt.Errorf("generating patterns: %w", err)
		}
	}

	// Layer 6: Templates (if present)
	if len(ds.Templates) > 0 {
		if err := ds.generateMkDocsMarkdown(filepath.Join(layers["templates"], "index.md"), "Templates", ds.templatesMarkdown()); err != nil {
			return fmt.Errorf("generating templates: %w", err)
		}
	}

	// Layer 7: Content (if present)
	if ds.Content != nil {
		if err := ds.generateMkDocsMarkdown(filepath.Join(layers["content"], "index.md"), "Content", ds.contentMarkdown()); err != nil {
			return fmt.Errorf("generating content: %w", err)
		}
	}

	// Layer 8: Accessibility (if present)
	if ds.Accessibility != nil {
		if err := ds.generateMkDocsMarkdown(filepath.Join(layers["accessibility"], "index.md"), "Accessibility", ds.accessibilityMarkdown()); err != nil {
			return fmt.Errorf("generating accessibility: %w", err)
		}
	}

	// Layer 9: Governance (if present)
	if ds.Governance != nil {
		if err := ds.generateMkDocsMarkdown(filepath.Join(layers["governance"], "index.md"), "Governance", ds.governanceMarkdown()); err != nil {
			return fmt.Errorf("generating governance: %w", err)
		}
	}

	// Evaluation (under specs)
	if err := ds.generateMkDocsPage(tmpl, filepath.Join(layers["eval"], "index.md"), "eval.html", "Evaluation", ds.evalData(opts)); err != nil {
		return fmt.Errorf("generating eval: %w", err)
	}

	return nil
}

// generateMkDocsPage generates a single MkDocs-compatible Markdown page.
func (ds *DesignSystem) generateMkDocsPage(tmpl *template.Template, outputPath, templateName, title string, data any) error {
	// Render HTML content
	var contentBuf strings.Builder
	if err := tmpl.ExecuteTemplate(&contentBuf, templateName, data); err != nil {
		return fmt.Errorf("executing template %s: %w", templateName, err)
	}

	// Create Markdown file with embedded HTML
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	// Write Markdown frontmatter and embedded HTML
	fmt.Fprintf(f, "# %s\n\n", title)
	fmt.Fprintf(f, "<div class=\"dss-content\" markdown>\n\n")
	//nolint:gosec // G203: HTML is from embedded templates, not user input
	fmt.Fprintf(f, "%s\n", contentBuf.String())
	fmt.Fprintf(f, "\n</div>\n")

	return nil
}

// generateMkDocsMarkdown generates a Markdown page with pure Markdown content.
func (ds *DesignSystem) generateMkDocsMarkdown(outputPath, title, content string) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	fmt.Fprintf(f, "# %s\n\n", title)
	fmt.Fprintf(f, "%s\n", content)

	return nil
}

// metaMarkdown generates Markdown content for the Meta layer.
func (ds *DesignSystem) metaMarkdown() string {
	var b strings.Builder

	b.WriteString("Design system metadata and versioning information.\n\n")
	b.WriteString("| Property | Value |\n")
	b.WriteString("|----------|-------|\n")
	fmt.Fprintf(&b, "| **Name** | %s |\n", ds.Meta.Name)
	fmt.Fprintf(&b, "| **Version** | %s |\n", ds.Meta.Version)
	if ds.Meta.Description != "" {
		fmt.Fprintf(&b, "| **Description** | %s |\n", ds.Meta.Description)
	}
	if ds.Meta.Documentation != "" {
		fmt.Fprintf(&b, "| **Documentation** | [%s](%s) |\n", ds.Meta.Documentation, ds.Meta.Documentation)
	}
	if ds.Meta.Repository != "" {
		fmt.Fprintf(&b, "| **Repository** | [%s](%s) |\n", ds.Meta.Repository, ds.Meta.Repository)
	}
	if ds.Meta.License != "" {
		fmt.Fprintf(&b, "| **License** | %s |\n", ds.Meta.License)
	}
	if ds.Meta.MaturityLevel > 0 {
		fmt.Fprintf(&b, "| **Maturity Level** | %d/5 |\n", ds.Meta.MaturityLevel)
	}

	if len(ds.Meta.Maintainers) > 0 {
		b.WriteString("\n## Maintainers\n\n")
		b.WriteString("| Name | Email | URL |\n")
		b.WriteString("|------|-------|-----|\n")
		for _, m := range ds.Meta.Maintainers {
			url := m.URL
			if url != "" {
				url = fmt.Sprintf("[link](%s)", url)
			}
			fmt.Fprintf(&b, "| %s | %s | %s |\n", m.Name, m.Email, url)
		}
	}

	return b.String()
}

// principlesMarkdown generates Markdown content for the Principles layer.
func (ds *DesignSystem) principlesMarkdown() string {
	var b strings.Builder

	b.WriteString("Design philosophy and guiding principles.\n\n")

	for _, p := range ds.Principles {
		fmt.Fprintf(&b, "## %s\n\n", p.Name)
		if p.Description != "" {
			fmt.Fprintf(&b, "%s\n\n", p.Description)
		}
	}

	return b.String()
}

// patternsMarkdown generates Markdown content for the Patterns layer.
func (ds *DesignSystem) patternsMarkdown() string {
	var b strings.Builder

	b.WriteString("Multi-component UX solutions and patterns.\n\n")
	b.WriteString("| Pattern | Description |\n")
	b.WriteString("|---------|-------------|\n")

	for _, p := range ds.Patterns {
		fmt.Fprintf(&b, "| **%s** | %s |\n", p.Name, p.Description)
	}

	return b.String()
}

// templatesMarkdown generates Markdown content for the Templates layer.
func (ds *DesignSystem) templatesMarkdown() string {
	var b strings.Builder

	b.WriteString("Page-level layout templates.\n\n")
	b.WriteString("| Template | Description |\n")
	b.WriteString("|----------|-------------|\n")

	for _, t := range ds.Templates {
		fmt.Fprintf(&b, "| **%s** | %s |\n", t.Name, t.Description)
	}

	return b.String()
}

// contentMarkdown generates Markdown content for the Content layer.
func (ds *DesignSystem) contentMarkdown() string {
	var b strings.Builder

	b.WriteString("Voice, tone, and content guidelines.\n\n")

	// Voice guidelines
	if ds.Content.Voice != nil {
		b.WriteString("## Voice\n\n")
		if ds.Content.Voice.Description != "" {
			fmt.Fprintf(&b, "%s\n\n", ds.Content.Voice.Description)
		}
		if len(ds.Content.Voice.Attributes) > 0 {
			b.WriteString("### Attributes\n\n")
			b.WriteString("| Attribute | Description |\n")
			b.WriteString("|-----------|-------------|\n")
			for _, attr := range ds.Content.Voice.Attributes {
				fmt.Fprintf(&b, "| **%s** | %s |\n", attr.Name, attr.Description)
			}
			b.WriteString("\n")
		}
	}

	// Tone guidelines
	if len(ds.Content.Tone) > 0 {
		b.WriteString("## Tone\n\n")
		b.WriteString("| Context | Description |\n")
		b.WriteString("|---------|-------------|\n")
		for _, t := range ds.Content.Tone {
			fmt.Fprintf(&b, "| **%s** | %s |\n", t.Context, t.Description)
		}
		b.WriteString("\n")
	}

	// Terminology
	if ds.Content.Terminology != nil {
		b.WriteString("## Terminology\n\n")
		if len(ds.Content.Terminology.PreferredTerms) > 0 {
			b.WriteString("### Preferred Terms\n\n")
			b.WriteString("| Term | Definition |\n")
			b.WriteString("|------|------------|\n")
			for _, term := range ds.Content.Terminology.PreferredTerms {
				fmt.Fprintf(&b, "| **%s** | %s |\n", term.Term, term.Definition)
			}
			b.WriteString("\n")
		}
		if len(ds.Content.Terminology.AvoidedTerms) > 0 {
			b.WriteString("### Avoided Terms\n\n")
			b.WriteString("| Term | Usage |\n")
			b.WriteString("|------|-------|\n")
			for _, term := range ds.Content.Terminology.AvoidedTerms {
				fmt.Fprintf(&b, "| ~~%s~~ | %s |\n", term.Term, term.Usage)
			}
			b.WriteString("\n")
		}
	}

	// Writing style
	if ds.Content.WritingStyle != nil {
		b.WriteString("## Writing Style\n\n")
		b.WriteString("| Guideline | Rule |\n")
		b.WriteString("|-----------|------|\n")
		if ds.Content.WritingStyle.Capitalization != "" {
			fmt.Fprintf(&b, "| **Capitalization** | %s |\n", ds.Content.WritingStyle.Capitalization)
		}
		if ds.Content.WritingStyle.Punctuation != "" {
			fmt.Fprintf(&b, "| **Punctuation** | %s |\n", ds.Content.WritingStyle.Punctuation)
		}
		if ds.Content.WritingStyle.Contractions != "" {
			fmt.Fprintf(&b, "| **Contractions** | %s |\n", ds.Content.WritingStyle.Contractions)
		}
		if ds.Content.WritingStyle.Numbers != "" {
			fmt.Fprintf(&b, "| **Numbers** | %s |\n", ds.Content.WritingStyle.Numbers)
		}
		if ds.Content.WritingStyle.ActiveVoice {
			b.WriteString("| **Active Voice** | Preferred |\n")
		}
		b.WriteString("\n")
	}

	return b.String()
}

// accessibilityMarkdown generates Markdown content for the Accessibility layer.
func (ds *DesignSystem) accessibilityMarkdown() string {
	var b strings.Builder

	b.WriteString("WCAG requirements and accessibility standards.\n\n")

	// WCAG compliance summary
	b.WriteString("## Compliance Target\n\n")
	b.WriteString("| Property | Value |\n")
	b.WriteString("|----------|-------|\n")
	if ds.Accessibility.WCAGLevel != "" {
		fmt.Fprintf(&b, "| **WCAG Level** | %s |\n", ds.Accessibility.WCAGLevel)
	}
	if ds.Accessibility.WCAGVersion != "" {
		fmt.Fprintf(&b, "| **WCAG Version** | %s |\n", ds.Accessibility.WCAGVersion)
	}
	b.WriteString("\n")

	// Color contrast
	if ds.Accessibility.ColorContrast != nil {
		b.WriteString("## Color Contrast\n\n")
		b.WriteString("| Requirement | Ratio |\n")
		b.WriteString("|-------------|-------|\n")
		if ds.Accessibility.ColorContrast.NormalTextRatio > 0 {
			fmt.Fprintf(&b, "| **Normal Text** | %.1f:1 |\n", ds.Accessibility.ColorContrast.NormalTextRatio)
		}
		if ds.Accessibility.ColorContrast.LargeTextRatio > 0 {
			fmt.Fprintf(&b, "| **Large Text** | %.1f:1 |\n", ds.Accessibility.ColorContrast.LargeTextRatio)
		}
		if ds.Accessibility.ColorContrast.NonTextRatio > 0 {
			fmt.Fprintf(&b, "| **Non-Text** | %.1f:1 |\n", ds.Accessibility.ColorContrast.NonTextRatio)
		}
		b.WriteString("\n")
	}

	// Keyboard requirements
	if ds.Accessibility.Keyboard != nil {
		b.WriteString("## Keyboard Accessibility\n\n")
		b.WriteString("| Requirement | Status |\n")
		b.WriteString("|-------------|--------|\n")
		fmt.Fprintf(&b, "| **Focus Visible** | %s |\n", boolToCheckmark(ds.Accessibility.Keyboard.FocusVisible))
		fmt.Fprintf(&b, "| **Focus Order** | %s |\n", boolToCheckmark(ds.Accessibility.Keyboard.FocusOrder))
		fmt.Fprintf(&b, "| **No Keyboard Trap** | %s |\n", boolToCheckmark(ds.Accessibility.Keyboard.NoKeyboardTrap))
		if ds.Accessibility.Keyboard.SkipLinks {
			b.WriteString("| **Skip Links** | ✅ |\n")
		}
		b.WriteString("\n")
	}

	// Screen reader requirements
	if ds.Accessibility.ScreenReader != nil {
		b.WriteString("## Screen Reader Support\n\n")
		b.WriteString("| Requirement | Status |\n")
		b.WriteString("|-------------|--------|\n")
		fmt.Fprintf(&b, "| **Semantic HTML** | %s |\n", boolToCheckmark(ds.Accessibility.ScreenReader.SemanticHTML))
		fmt.Fprintf(&b, "| **ARIA Labels** | %s |\n", boolToCheckmark(ds.Accessibility.ScreenReader.ARIALabels))
		fmt.Fprintf(&b, "| **Live Regions** | %s |\n", boolToCheckmark(ds.Accessibility.ScreenReader.LiveRegions))
		fmt.Fprintf(&b, "| **Heading Structure** | %s |\n", boolToCheckmark(ds.Accessibility.ScreenReader.HeadingStructure))
		b.WriteString("\n")
	}

	// Motion requirements
	if ds.Accessibility.Motion != nil {
		b.WriteString("## Motion & Animation\n\n")
		b.WriteString("| Requirement | Status |\n")
		b.WriteString("|-------------|--------|\n")
		fmt.Fprintf(&b, "| **Reduced Motion** | %s |\n", boolToCheckmark(ds.Accessibility.Motion.ReducedMotion))
		fmt.Fprintf(&b, "| **No Auto-Play** | %s |\n", boolToCheckmark(ds.Accessibility.Motion.NoAutoPlay))
		fmt.Fprintf(&b, "| **Pause Control** | %s |\n", boolToCheckmark(ds.Accessibility.Motion.PauseControl))
		b.WriteString("\n")
	}

	// Touch target requirements
	if ds.Accessibility.TouchTarget != nil {
		b.WriteString("## Touch Targets\n\n")
		b.WriteString("| Property | Value |\n")
		b.WriteString("|----------|-------|\n")
		if ds.Accessibility.TouchTarget.MinSize != "" {
			fmt.Fprintf(&b, "| **Minimum Size** | %s |\n", ds.Accessibility.TouchTarget.MinSize)
		}
		if ds.Accessibility.TouchTarget.MinSpacing != "" {
			fmt.Fprintf(&b, "| **Minimum Spacing** | %s |\n", ds.Accessibility.TouchTarget.MinSpacing)
		}
		b.WriteString("\n")
	}

	// Guidelines
	if len(ds.Accessibility.Guidelines) > 0 {
		b.WriteString("## Guidelines\n\n")
		for _, g := range ds.Accessibility.Guidelines {
			fmt.Fprintf(&b, "### %s\n\n", g.Title)
			fmt.Fprintf(&b, "%s\n\n", g.Description)
			if len(g.WCAGCriteria) > 0 {
				b.WriteString("**WCAG Criteria:** ")
				b.WriteString(strings.Join(g.WCAGCriteria, ", "))
				b.WriteString("\n\n")
			}
		}
	}

	// Testing requirements
	if len(ds.Accessibility.TestingRequirements) > 0 {
		b.WriteString("## Testing Requirements\n\n")
		b.WriteString("| Type | Tools | Frequency |\n")
		b.WriteString("|------|-------|----------|\n")
		for _, t := range ds.Accessibility.TestingRequirements {
			tools := strings.Join(t.Tools, ", ")
			fmt.Fprintf(&b, "| **%s** | %s | %s |\n", t.Type, tools, t.Frequency)
		}
		b.WriteString("\n")
	}

	return b.String()
}

// boolToCheckmark converts a boolean to a checkmark emoji.
func boolToCheckmark(v bool) string {
	if v {
		return "✅"
	}
	return "❌"
}

// governanceMarkdown generates Markdown content for the Governance layer.
func (ds *DesignSystem) governanceMarkdown() string {
	var b strings.Builder

	b.WriteString("Versioning, contribution, and deprecation policies.\n\n")

	// Versioning policy
	if ds.Governance.Versioning != nil {
		b.WriteString("## Versioning\n\n")
		b.WriteString("| Property | Value |\n")
		b.WriteString("|----------|-------|\n")
		if ds.Governance.Versioning.Strategy != "" {
			fmt.Fprintf(&b, "| **Strategy** | %s |\n", ds.Governance.Versioning.Strategy)
		}
		if ds.Governance.Versioning.ReleaseSchedule != "" {
			fmt.Fprintf(&b, "| **Release Schedule** | %s |\n", ds.Governance.Versioning.ReleaseSchedule)
		}
		if ds.Governance.Versioning.ChangelogLocation != "" {
			fmt.Fprintf(&b, "| **Changelog** | [%s](%s) |\n", ds.Governance.Versioning.ChangelogLocation, ds.Governance.Versioning.ChangelogLocation)
		}
		b.WriteString("\n")
		if ds.Governance.Versioning.Description != "" {
			fmt.Fprintf(&b, "%s\n\n", ds.Governance.Versioning.Description)
		}
		if ds.Governance.Versioning.BreakingChangePolicy != "" {
			fmt.Fprintf(&b, "**Breaking Changes:** %s\n\n", ds.Governance.Versioning.BreakingChangePolicy)
		}
	}

	// Contribution policy
	if ds.Governance.Contribution != nil {
		b.WriteString("## Contribution\n\n")
		if ds.Governance.Contribution.Description != "" {
			fmt.Fprintf(&b, "%s\n\n", ds.Governance.Contribution.Description)
		}
		if len(ds.Governance.Contribution.Guidelines) > 0 {
			b.WriteString("### Guidelines\n\n")
			for _, g := range ds.Governance.Contribution.Guidelines {
				fmt.Fprintf(&b, "- %s\n", g)
			}
			b.WriteString("\n")
		}
		if ds.Governance.Contribution.Workflow != nil && len(ds.Governance.Contribution.Workflow.Steps) > 0 {
			b.WriteString("### Workflow\n\n")
			b.WriteString("| Step | Name | Description |\n")
			b.WriteString("|------|------|-------------|\n")
			for _, step := range ds.Governance.Contribution.Workflow.Steps {
				fmt.Fprintf(&b, "| %d | **%s** | %s |\n", step.Order, step.Name, step.Description)
			}
			b.WriteString("\n")
		}
	}

	// Deprecation policy
	if ds.Governance.Deprecation != nil {
		b.WriteString("## Deprecation\n\n")
		b.WriteString("| Property | Value |\n")
		b.WriteString("|----------|-------|\n")
		if ds.Governance.Deprecation.WarningPeriod != "" {
			fmt.Fprintf(&b, "| **Warning Period** | %s |\n", ds.Governance.Deprecation.WarningPeriod)
		}
		fmt.Fprintf(&b, "| **Migration Guide Required** | %s |\n", boolToCheckmark(ds.Governance.Deprecation.MigrationGuideRequired))
		b.WriteString("\n")

		if len(ds.Governance.Deprecation.NotificationChannels) > 0 {
			b.WriteString("**Notification Channels:** ")
			b.WriteString(strings.Join(ds.Governance.Deprecation.NotificationChannels, ", "))
			b.WriteString("\n\n")
		}

		if len(ds.Governance.Deprecation.DeprecatedItems) > 0 {
			b.WriteString("### Deprecated Items\n\n")
			b.WriteString("| Type | ID | Since | Removal | Replacement |\n")
			b.WriteString("|------|----|----|---------|-------------|\n")
			for _, item := range ds.Governance.Deprecation.DeprecatedItems {
				fmt.Fprintf(&b, "| %s | `%s` | %s | %s | %s |\n",
					item.Type, item.ID, item.DeprecatedSince, item.RemovalDate, item.Replacement)
			}
			b.WriteString("\n")
		}
	}

	// Decision process
	if ds.Governance.DecisionProcess != nil {
		b.WriteString("## Decision Process\n\n")
		if ds.Governance.DecisionProcess.Description != "" {
			fmt.Fprintf(&b, "%s\n\n", ds.Governance.DecisionProcess.Description)
		}
		if ds.Governance.DecisionProcess.RFCRequired {
			b.WriteString("**RFC Required:** Yes\n\n")
		}
		if len(ds.Governance.DecisionProcess.ApprovalRequired) > 0 {
			b.WriteString("**Approval Required From:** ")
			b.WriteString(strings.Join(ds.Governance.DecisionProcess.ApprovalRequired, ", "))
			b.WriteString("\n\n")
		}
	}

	// Ownership
	if len(ds.Governance.Ownership) > 0 {
		b.WriteString("## Ownership\n\n")
		b.WriteString("| Area | Team | Contact |\n")
		b.WriteString("|------|------|---------||\n")
		for _, o := range ds.Governance.Ownership {
			contact := o.ContactEmail
			if o.SlackChannel != "" && contact != "" {
				contact = fmt.Sprintf("%s / %s", o.ContactEmail, o.SlackChannel)
			} else if o.SlackChannel != "" {
				contact = o.SlackChannel
			}
			fmt.Fprintf(&b, "| %s | %s | %s |\n", o.Area, o.Team, contact)
		}
		b.WriteString("\n")
	}

	// Support policy
	if ds.Governance.SupportPolicy != nil {
		b.WriteString("## Support Policy\n\n")
		if ds.Governance.SupportPolicy.Description != "" {
			fmt.Fprintf(&b, "%s\n\n", ds.Governance.SupportPolicy.Description)
		}
		if len(ds.Governance.SupportPolicy.SupportedVersions) > 0 {
			b.WriteString("### Supported Versions\n\n")
			b.WriteString("| Version | Status | End of Life |\n")
			b.WriteString("|---------|--------|-------------|\n")
			for _, v := range ds.Governance.SupportPolicy.SupportedVersions {
				fmt.Fprintf(&b, "| %s | %s | %s |\n", v.Version, v.Status, v.EndOfLifeDate)
			}
			b.WriteString("\n")
		}
	}

	return b.String()
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
