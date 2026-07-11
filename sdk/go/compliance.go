package dss

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// ComplianceReport represents a complete validation report for release.
type ComplianceReport struct {
	// Timestamp when the report was generated
	Timestamp time.Time `json:"timestamp"`

	// DesignSystem identifies the spec being validated against
	DesignSystem DesignSystemRef `json:"designSystem"`

	// Target describes what was validated
	Target ValidationTarget `json:"target"`

	// Status is the overall compliance status
	Status ComplianceStatus `json:"status"`

	// Score is the overall compliance score (0-100)
	Score int `json:"score"`

	// Categories breaks down compliance by category
	Categories []CategoryCompliance `json:"categories"`

	// Summary provides aggregate statistics
	Summary ComplianceSummary `json:"summary"`

	// Issues lists all validation findings
	Issues []ComplianceIssue `json:"issues,omitempty"`

	// Metadata contains additional context
	Metadata map[string]any `json:"metadata,omitempty"`
}

// DesignSystemRef identifies the design system spec.
type DesignSystemRef struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ValidationTarget describes what was validated.
type ValidationTarget struct {
	// Type is "file", "directory", or "project"
	Type string `json:"type"`

	// Path is the file or directory path
	Path string `json:"path"`

	// FileCount is the number of files validated
	FileCount int `json:"fileCount"`

	// Commit is the git commit hash if available
	Commit string `json:"commit,omitempty"`

	// Branch is the git branch if available
	Branch string `json:"branch,omitempty"`
}

// ComplianceStatus represents overall compliance state.
type ComplianceStatus string

const (
	// ComplianceStatusPass means all required checks passed
	ComplianceStatusPass ComplianceStatus = "pass"

	// ComplianceStatusFail means required checks failed
	ComplianceStatusFail ComplianceStatus = "fail"

	// ComplianceStatusWarn means no failures but warnings exist
	ComplianceStatusWarn ComplianceStatus = "warn"
)

// CategoryCompliance tracks compliance for a specific category.
type CategoryCompliance struct {
	// Category name (e.g., "colors", "spacing", "accessibility")
	Category string `json:"category"`

	// Status for this category
	Status ComplianceStatus `json:"status"`

	// Score for this category (0-100)
	Score int `json:"score"`

	// Checked is the number of checks run
	Checked int `json:"checked"`

	// Passed is the number of checks that passed
	Passed int `json:"passed"`

	// Failed is the number of checks that failed
	Failed int `json:"failed"`

	// Blocking indicates if failures in this category block release
	Blocking bool `json:"blocking"`
}

// ComplianceSummary provides aggregate statistics.
type ComplianceSummary struct {
	// TotalChecks is the total number of checks run
	TotalChecks int `json:"totalChecks"`

	// Passed is the number of passed checks
	Passed int `json:"passed"`

	// Errors is the count of error-level issues
	Errors int `json:"errors"`

	// Warnings is the count of warning-level issues
	Warnings int `json:"warnings"`

	// Infos is the count of info-level issues
	Infos int `json:"infos"`

	// BlockingIssues is the count of issues that block release
	BlockingIssues int `json:"blockingIssues"`

	// FixableIssues is the count of auto-fixable issues
	FixableIssues int `json:"fixableIssues"`
}

// ComplianceIssue represents a single compliance finding.
type ComplianceIssue struct {
	// Category this issue belongs to
	Category string `json:"category"`

	// Rule is the rule ID
	Rule string `json:"rule"`

	// Severity is error, warning, or info
	Severity string `json:"severity"`

	// Message describes the issue
	Message string `json:"message"`

	// File is the file path
	File string `json:"file,omitempty"`

	// Line is the line number
	Line int `json:"line,omitempty"`

	// Blocking indicates if this issue blocks release
	Blocking bool `json:"blocking"`

	// Fixable indicates if this can be auto-fixed
	Fixable bool `json:"fixable"`

	// FixSuggestion provides guidance on fixing
	FixSuggestion string `json:"fixSuggestion,omitempty"`
}

// ComplianceOptions configures compliance report generation.
type ComplianceOptions struct {
	// IncludeIssues includes detailed issue list in report
	IncludeIssues bool `json:"includeIssues,omitempty"`

	// BlockingCategories lists categories that block release
	// Default: ["colors", "accessibility"]
	BlockingCategories []string `json:"blockingCategories,omitempty"`

	// MinScore is the minimum acceptable score
	MinScore int `json:"minScore,omitempty"`

	// GitInfo includes git commit/branch info if available
	GitInfo bool `json:"gitInfo,omitempty"`
}

// DefaultComplianceOptions returns sensible defaults.
func DefaultComplianceOptions() *ComplianceOptions {
	return &ComplianceOptions{
		IncludeIssues:      true,
		BlockingCategories: []string{"colors", "accessibility"},
		MinScore:           80,
		GitInfo:            true,
	}
}

// GenerateComplianceReport creates a compliance report for a directory.
func (s *Service) GenerateComplianceReport(ctx context.Context, dir string, opts *ComplianceOptions) (*ComplianceReport, error) {
	if opts == nil {
		opts = DefaultComplianceOptions()
	}

	// Run validation
	validateOpts := &ValidateOptions{
		IncludeContext: true,
	}
	validationResult, err := s.ValidateDirectory(ctx, dir, validateOpts)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Build category compliance
	categoryMap := make(map[string]*CategoryCompliance)
	blockingSet := make(map[string]bool)
	for _, cat := range opts.BlockingCategories {
		blockingSet[cat] = true
	}

	// Initialize categories
	categories := []string{"colors", "spacing", "accessibility", "variants", "patterns"}
	for _, cat := range categories {
		categoryMap[cat] = &CategoryCompliance{
			Category: cat,
			Status:   ComplianceStatusPass,
			Score:    100,
			Blocking: blockingSet[cat],
		}
	}

	// Process violations
	var issues []ComplianceIssue
	for _, v := range validationResult.Violations {
		cat := ruleToCategory(v.Rule)
		if catComp, ok := categoryMap[cat]; ok {
			catComp.Checked++
			if v.Severity == "error" {
				catComp.Failed++
				catComp.Status = ComplianceStatusFail
			} else if v.Severity == "warning" && catComp.Status == ComplianceStatusPass {
				catComp.Status = ComplianceStatusWarn
			}
		}

		issue := ComplianceIssue{
			Category:      cat,
			Rule:          v.Rule,
			Severity:      v.Severity,
			Message:       v.Message,
			File:          v.File,
			Line:          v.Line,
			Blocking:      blockingSet[cat] && v.Severity == "error",
			Fixable:       isFixableRule(v.Rule),
			FixSuggestion: getFixSuggestion(v.Rule),
		}
		issues = append(issues, issue)
	}

	// Calculate category scores
	for _, catComp := range categoryMap {
		if catComp.Checked > 0 {
			catComp.Passed = catComp.Checked - catComp.Failed
			catComp.Score = (catComp.Passed * 100) / catComp.Checked
		}
	}

	// Build category list
	var categoryList []CategoryCompliance
	for _, cat := range categories {
		categoryList = append(categoryList, *categoryMap[cat])
	}

	// Calculate summary
	summary := ComplianceSummary{
		TotalChecks: validationResult.Files,
		Passed:      validationResult.Summary.Passed,
		Errors:      validationResult.Summary.Errors,
		Warnings:    validationResult.Summary.Warnings,
		Infos:       validationResult.Summary.Infos,
	}

	for _, issue := range issues {
		if issue.Blocking {
			summary.BlockingIssues++
		}
		if issue.Fixable {
			summary.FixableIssues++
		}
	}

	// Calculate overall score
	score := calculateComplianceScore(categoryList, summary)

	// Determine overall status
	status := ComplianceStatusPass
	if summary.BlockingIssues > 0 {
		status = ComplianceStatusFail
	} else if summary.Errors > 0 || score < opts.MinScore {
		status = ComplianceStatusFail
	} else if summary.Warnings > 0 {
		status = ComplianceStatusWarn
	}

	report := &ComplianceReport{
		Timestamp: time.Now().UTC(),
		DesignSystem: DesignSystemRef{
			Name:    s.ds.Meta.Name,
			Version: s.ds.Meta.Version,
		},
		Target: ValidationTarget{
			Type:      "directory",
			Path:      dir,
			FileCount: validationResult.Files,
		},
		Status:     status,
		Score:      score,
		Categories: categoryList,
		Summary:    summary,
	}

	if opts.IncludeIssues {
		report.Issues = issues
	}

	return report, nil
}

// ComplianceCertificate represents proof of compliance.
type ComplianceCertificate struct {
	// ID is a unique certificate identifier
	ID string `json:"id"`

	// Timestamp when certificate was issued
	Timestamp time.Time `json:"timestamp"`

	// ExpiresAt when the certificate expires (optional)
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`

	// DesignSystem the compliance is against
	DesignSystem DesignSystemRef `json:"designSystem"`

	// Target what was certified
	Target ValidationTarget `json:"target"`

	// Status compliance status
	Status ComplianceStatus `json:"status"`

	// Score compliance score
	Score int `json:"score"`

	// Hash cryptographic hash of the report
	Hash string `json:"hash"`

	// Categories summary per category
	Categories []CategoryCertification `json:"categories"`
}

// CategoryCertification is a simplified category status.
type CategoryCertification struct {
	Category string           `json:"category"`
	Status   ComplianceStatus `json:"status"`
	Score    int              `json:"score"`
}

// GenerateCertificate creates a compliance certificate from a report.
func (s *Service) GenerateCertificate(report *ComplianceReport) *ComplianceCertificate {
	// Generate unique ID
	idData := fmt.Sprintf("%s-%s-%s-%d",
		report.DesignSystem.Name,
		report.DesignSystem.Version,
		report.Target.Path,
		report.Timestamp.UnixNano(),
	)
	hash := sha256.Sum256([]byte(idData))
	id := hex.EncodeToString(hash[:8])

	// Build category certifications
	var categories []CategoryCertification
	for _, cat := range report.Categories {
		categories = append(categories, CategoryCertification{
			Category: cat.Category,
			Status:   cat.Status,
			Score:    cat.Score,
		})
	}

	// Generate report hash
	reportData := fmt.Sprintf("%v", report)
	reportHash := sha256.Sum256([]byte(reportData))

	return &ComplianceCertificate{
		ID:           id,
		Timestamp:    report.Timestamp,
		DesignSystem: report.DesignSystem,
		Target:       report.Target,
		Status:       report.Status,
		Score:        report.Score,
		Hash:         hex.EncodeToString(reportHash[:]),
		Categories:   categories,
	}
}

// ReleaseGate represents the release decision.
type ReleaseGate struct {
	// Approved indicates if release is approved
	Approved bool `json:"approved"`

	// Reason explains the decision
	Reason string `json:"reason"`

	// BlockingIssues lists issues preventing release
	BlockingIssues []string `json:"blockingIssues,omitempty"`

	// Warnings lists non-blocking concerns
	Warnings []string `json:"warnings,omitempty"`

	// Score is the compliance score
	Score int `json:"score"`

	// Certificate is the compliance certificate if approved
	Certificate *ComplianceCertificate `json:"certificate,omitempty"`
}

// ReleaseGateOptions configures release gate behavior.
type ReleaseGateOptions struct {
	// MinScore is the minimum acceptable score
	MinScore int `json:"minScore,omitempty"`

	// RequireZeroErrors requires no error-level issues
	RequireZeroErrors bool `json:"requireZeroErrors,omitempty"`

	// BlockingCategories lists categories that must pass
	BlockingCategories []string `json:"blockingCategories,omitempty"`

	// AllowWarnings allows release with warnings
	AllowWarnings bool `json:"allowWarnings,omitempty"`
}

// DefaultReleaseGateOptions returns strict release gate defaults.
func DefaultReleaseGateOptions() *ReleaseGateOptions {
	return &ReleaseGateOptions{
		MinScore:           80,
		RequireZeroErrors:  true,
		BlockingCategories: []string{"colors", "accessibility"},
		AllowWarnings:      true,
	}
}

// CheckReleaseGate evaluates if a release should be approved.
func (s *Service) CheckReleaseGate(ctx context.Context, dir string, opts *ReleaseGateOptions) (*ReleaseGate, error) {
	if opts == nil {
		opts = DefaultReleaseGateOptions()
	}

	// Generate compliance report
	complianceOpts := &ComplianceOptions{
		IncludeIssues:      true,
		BlockingCategories: opts.BlockingCategories,
		MinScore:           opts.MinScore,
	}

	report, err := s.GenerateComplianceReport(ctx, dir, complianceOpts)
	if err != nil {
		return nil, err
	}

	gate := &ReleaseGate{
		Score: report.Score,
	}

	// Check score
	if report.Score < opts.MinScore {
		gate.Approved = false
		gate.Reason = fmt.Sprintf("Score %d is below minimum %d", report.Score, opts.MinScore)
		gate.BlockingIssues = append(gate.BlockingIssues,
			fmt.Sprintf("Compliance score %d < %d required", report.Score, opts.MinScore))
		return gate, nil
	}

	// Check blocking issues
	for _, issue := range report.Issues {
		if issue.Blocking {
			gate.BlockingIssues = append(gate.BlockingIssues,
				fmt.Sprintf("[%s] %s: %s", issue.Category, issue.Rule, issue.Message))
		} else if issue.Severity == "warning" {
			gate.Warnings = append(gate.Warnings,
				fmt.Sprintf("[%s] %s: %s", issue.Category, issue.Rule, issue.Message))
		}
	}

	if len(gate.BlockingIssues) > 0 {
		gate.Approved = false
		gate.Reason = fmt.Sprintf("%d blocking issues found", len(gate.BlockingIssues))
		return gate, nil
	}

	// Check zero errors requirement
	if opts.RequireZeroErrors && report.Summary.Errors > 0 {
		gate.Approved = false
		gate.Reason = fmt.Sprintf("%d errors found, zero required", report.Summary.Errors)
		return gate, nil
	}

	// Check warnings
	if !opts.AllowWarnings && len(gate.Warnings) > 0 {
		gate.Approved = false
		gate.Reason = fmt.Sprintf("%d warnings found, none allowed", len(gate.Warnings))
		return gate, nil
	}

	// Approved!
	gate.Approved = true
	if len(gate.Warnings) > 0 {
		gate.Reason = fmt.Sprintf("Approved with %d warnings", len(gate.Warnings))
	} else {
		gate.Reason = "All checks passed"
	}

	// Generate certificate
	gate.Certificate = s.GenerateCertificate(report)

	return gate, nil
}

// Helper functions

func ruleToCategory(rule string) string {
	switch rule {
	case "no-hardcoded-colors":
		return "colors"
	case "use-spacing-scale":
		return "spacing"
	case "img-alt-required", "button-accessible-name":
		return "accessibility"
	case "valid-variant":
		return "variants"
	case "no-multiple-primary-buttons", "no-nested-cards":
		return "patterns"
	default:
		return "other"
	}
}

func isFixableRule(rule string) bool {
	fixable := map[string]bool{
		"no-hardcoded-colors":    true,
		"use-spacing-scale":      true,
		"img-alt-required":       true,
		"button-accessible-name": true,
	}
	return fixable[rule]
}

func getFixSuggestion(rule string) string {
	suggestions := map[string]string{
		"no-hardcoded-colors":         "Use design token CSS variable instead of hardcoded color",
		"use-spacing-scale":           "Use design token spacing value instead of pixel value",
		"img-alt-required":            "Add descriptive alt attribute to image",
		"button-accessible-name":      "Add aria-label to icon-only button",
		"valid-variant":               "Use one of the allowed variants from the spec",
		"no-multiple-primary-buttons": "Use only one primary button per view",
		"no-nested-cards":             "Avoid nesting Card components",
	}
	return suggestions[rule]
}

func calculateComplianceScore(categories []CategoryCompliance, summary ComplianceSummary) int {
	if len(categories) == 0 {
		return 100
	}

	// Weight categories
	weights := map[string]int{
		"colors":        25,
		"spacing":       20,
		"accessibility": 30,
		"variants":      15,
		"patterns":      10,
	}

	totalWeight := 0
	weightedScore := 0

	for _, cat := range categories {
		weight := weights[cat.Category]
		if weight == 0 {
			weight = 10
		}
		totalWeight += weight
		weightedScore += cat.Score * weight
	}

	if totalWeight == 0 {
		return 100
	}

	// Start with weighted category score
	score := weightedScore / totalWeight

	// Deduct for errors
	score -= summary.Errors * 5
	score -= summary.Warnings * 1

	// Clamp to 0-100
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}
