package dss

import (
	"context"
	"fmt"
)

// FixLoopResult represents the outcome of a fix-validate loop.
type FixLoopResult struct {
	// Converged indicates if the loop reached a stable state
	Converged bool `json:"converged"`

	// Iterations is the number of fix-validate cycles run
	Iterations int `json:"iterations"`

	// InitialViolations is the count before fixing
	InitialViolations int `json:"initialViolations"`

	// FinalViolations is the count after all fixes
	FinalViolations int `json:"finalViolations"`

	// FixedCount is the total number of issues fixed
	FixedCount int `json:"fixedCount"`

	// UnfixableCount is issues that couldn't be auto-fixed
	UnfixableCount int `json:"unfixableCount"`

	// Iterations detail for each cycle
	IterationDetails []IterationDetail `json:"iterationDetails"`

	// RemainingIssues lists issues that couldn't be fixed
	RemainingIssues []Violation `json:"remainingIssues,omitempty"`

	// Status is the final compliance status
	Status ComplianceStatus `json:"status"`
}

// IterationDetail captures a single fix-validate cycle.
type IterationDetail struct {
	// Iteration number (1-indexed)
	Iteration int `json:"iteration"`

	// ViolationsBefore at start of iteration
	ViolationsBefore int `json:"violationsBefore"`

	// FixesApplied in this iteration
	FixesApplied int `json:"fixesApplied"`

	// ViolationsAfter at end of iteration
	ViolationsAfter int `json:"violationsAfter"`

	// NewViolations introduced by fixes (regression)
	NewViolations int `json:"newViolations"`

	// FilesModified in this iteration
	FilesModified []string `json:"filesModified"`
}

// FixLoopOptions configures the fix-validate loop.
type FixLoopOptions struct {
	// MaxIterations is the maximum number of fix-validate cycles
	// Default: 3
	MaxIterations int `json:"maxIterations,omitempty"`

	// DryRun previews fixes without applying
	DryRun bool `json:"dryRun,omitempty"`

	// StopOnRegression stops if fixes introduce new violations
	StopOnRegression bool `json:"stopOnRegression,omitempty"`

	// Rules limits which rules to fix (empty = all fixable rules)
	Rules []string `json:"rules,omitempty"`

	// Files limits which files to process (empty = all)
	Files []string `json:"files,omitempty"`

	// Extensions limits file extensions to process
	Extensions []string `json:"extensions,omitempty"`
}

// DefaultFixLoopOptions returns sensible defaults.
func DefaultFixLoopOptions() *FixLoopOptions {
	return &FixLoopOptions{
		MaxIterations:    3,
		DryRun:           false,
		StopOnRegression: true,
		Extensions:       []string{".tsx", ".jsx", ".ts", ".js", ".css"},
	}
}

// RunFixLoop executes the fix-validate loop until convergence or max iterations.
func (s *Service) RunFixLoop(ctx context.Context, dir string, opts *FixLoopOptions) (*FixLoopResult, error) {
	if opts == nil {
		opts = DefaultFixLoopOptions()
	}

	result := &FixLoopResult{
		IterationDetails: []IterationDetail{},
	}

	// Initial validation
	validateOpts := &ValidateOptions{
		Extensions: opts.Extensions,
	}
	initialValidation, err := s.ValidateDirectory(ctx, dir, validateOpts)
	if err != nil {
		return nil, fmt.Errorf("initial validation failed: %w", err)
	}

	result.InitialViolations = len(initialValidation.Violations)

	// Track violations across iterations
	currentViolations := initialValidation.Violations
	previousViolationCount := len(currentViolations)

	// Run fix-validate loop
	for i := 1; i <= opts.MaxIterations; i++ {
		result.Iterations = i

		detail := IterationDetail{
			Iteration:        i,
			ViolationsBefore: len(currentViolations),
		}

		// No violations = converged
		if len(currentViolations) == 0 {
			result.Converged = true
			break
		}

		// Apply fixes
		fixOpts := &FixOptions{
			DryRun: opts.DryRun,
			Rules:  opts.Rules,
		}
		fixResults, err := s.FixDirectory(ctx, dir, fixOpts)
		if err != nil {
			return nil, fmt.Errorf("fix iteration %d failed: %w", i, err)
		}

		// Count fixes applied
		for _, fr := range fixResults {
			detail.FixesApplied += len(fr.Fixes)
			if len(fr.Fixes) > 0 {
				detail.FilesModified = append(detail.FilesModified, fr.File)
			}
		}
		result.FixedCount += detail.FixesApplied

		// If dry run, we can't re-validate
		if opts.DryRun {
			result.IterationDetails = append(result.IterationDetails, detail)
			break
		}

		// Re-validate
		revalidation, err := s.ValidateDirectory(ctx, dir, validateOpts)
		if err != nil {
			return nil, fmt.Errorf("re-validation iteration %d failed: %w", i, err)
		}

		currentViolations = revalidation.Violations
		detail.ViolationsAfter = len(currentViolations)

		// Check for regressions
		detail.NewViolations = countNewViolations(initialValidation.Violations, currentViolations)
		if opts.StopOnRegression && detail.NewViolations > 0 {
			result.IterationDetails = append(result.IterationDetails, detail)
			result.Status = ComplianceStatusFail
			result.RemainingIssues = currentViolations
			return result, fmt.Errorf("fixes introduced %d new violations", detail.NewViolations)
		}

		result.IterationDetails = append(result.IterationDetails, detail)

		// Check convergence
		if len(currentViolations) == 0 {
			result.Converged = true
			break
		}

		// Check if making progress
		if len(currentViolations) >= previousViolationCount {
			// No progress, remaining issues are unfixable
			break
		}

		previousViolationCount = len(currentViolations)
	}

	result.FinalViolations = len(currentViolations)
	result.RemainingIssues = currentViolations
	result.UnfixableCount = len(currentViolations)

	// Determine final status
	if result.FinalViolations == 0 {
		result.Status = ComplianceStatusPass
	} else {
		hasErrors := false
		for _, v := range currentViolations {
			if v.Severity == "error" {
				hasErrors = true
				break
			}
		}
		if hasErrors {
			result.Status = ComplianceStatusFail
		} else {
			result.Status = ComplianceStatusWarn
		}
	}

	return result, nil
}

// FixAndVerify runs a single fix cycle and verifies the fix resolved the issue.
func (s *Service) FixAndVerify(ctx context.Context, file string, opts *FixOptions) (*FixVerifyResult, error) {
	if opts == nil {
		opts = DefaultFixOptions()
	}

	result := &FixVerifyResult{
		File: file,
	}

	// Validate before
	validateOpts := &ValidateOptions{}
	beforeResult, err := s.ValidateFile(ctx, file, validateOpts)
	if err != nil {
		return nil, fmt.Errorf("pre-fix validation failed: %w", err)
	}
	result.ViolationsBefore = len(beforeResult.Violations)

	// Apply fix
	fixResult, err := s.FixFile(ctx, file, opts)
	if err != nil {
		return nil, fmt.Errorf("fix failed: %w", err)
	}
	result.FixesApplied = len(fixResult.Fixes)

	// If dry run, return preview
	if opts.DryRun {
		result.Verified = false
		result.Reason = "dry run mode - fixes not applied"
		return result, nil
	}

	// If no fixes applied
	if result.FixesApplied == 0 {
		result.Verified = true
		result.Reason = "no fixes needed"
		return result, nil
	}

	// Validate after
	afterResult, err := s.ValidateFile(ctx, file, validateOpts)
	if err != nil {
		return nil, fmt.Errorf("post-fix validation failed: %w", err)
	}
	result.ViolationsAfter = len(afterResult.Violations)

	// Calculate improvement
	result.Improvement = result.ViolationsBefore - result.ViolationsAfter

	// Check for regressions
	result.Regressions = countNewViolations(beforeResult.Violations, afterResult.Violations)

	// Determine verification status
	if result.ViolationsAfter == 0 {
		result.Verified = true
		result.Reason = "all violations fixed"
	} else if result.Improvement > 0 && result.Regressions == 0 {
		result.Verified = true
		result.Reason = fmt.Sprintf("improved: %d violations fixed, %d remaining", result.Improvement, result.ViolationsAfter)
	} else if result.Regressions > 0 {
		result.Verified = false
		result.Reason = fmt.Sprintf("regression: %d new violations introduced", result.Regressions)
	} else {
		result.Verified = false
		result.Reason = "no improvement"
	}

	return result, nil
}

// FixVerifyResult captures the outcome of a fix-and-verify operation.
type FixVerifyResult struct {
	// File that was processed
	File string `json:"file"`

	// ViolationsBefore the fix
	ViolationsBefore int `json:"violationsBefore"`

	// FixesApplied count
	FixesApplied int `json:"fixesApplied"`

	// ViolationsAfter the fix
	ViolationsAfter int `json:"violationsAfter"`

	// Improvement is violations reduced
	Improvement int `json:"improvement"`

	// Regressions is new violations introduced
	Regressions int `json:"regressions"`

	// Verified indicates if fix was successful
	Verified bool `json:"verified"`

	// Reason explains the verification outcome
	Reason string `json:"reason"`
}

// countNewViolations counts violations in after that weren't in before.
func countNewViolations(before, after []Violation) int {
	// Build set of before violations (file:line:rule)
	beforeSet := make(map[string]bool)
	for _, v := range before {
		key := fmt.Sprintf("%s:%d:%s", v.File, v.Line, v.Rule)
		beforeSet[key] = true
	}

	// Count new violations
	newCount := 0
	for _, v := range after {
		key := fmt.Sprintf("%s:%d:%s", v.File, v.Line, v.Rule)
		if !beforeSet[key] {
			newCount++
		}
	}

	return newCount
}
