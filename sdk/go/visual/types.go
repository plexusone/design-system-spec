// Package visual provides visual regression testing capabilities for DSS.
package visual

import "time"

// VisualTestSuite represents a collection of visual tests.
type VisualTestSuite struct {
	Version     string             `json:"version" yaml:"version"`
	Name        string             `json:"name" yaml:"name"`
	Description string             `json:"description,omitempty" yaml:"description,omitempty"`
	BaseURL     string             `json:"baseUrl,omitempty" yaml:"baseUrl,omitempty"`
	Defaults    VisualTestDefaults `json:"defaults,omitempty" yaml:"defaults,omitempty"`
	Tests       []VisualTest       `json:"tests" yaml:"tests"`
}

// VisualTestDefaults provides default values for all tests in a suite.
type VisualTestDefaults struct {
	Viewports     []Viewport     `json:"viewports,omitempty" yaml:"viewports,omitempty"`
	Threshold     float64        `json:"threshold,omitempty" yaml:"threshold,omitempty"`
	Stabilization *Stabilization `json:"stabilization,omitempty" yaml:"stabilization,omitempty"`
}

// VisualTest defines a single visual test case.
type VisualTest struct {
	ID            string         `json:"id" yaml:"id"`
	Name          string         `json:"name,omitempty" yaml:"name,omitempty"`
	Component     string         `json:"component" yaml:"component"`
	Variant       string         `json:"variant,omitempty" yaml:"variant,omitempty"`
	URL           string         `json:"url" yaml:"url"`
	Selector      string         `json:"selector,omitempty" yaml:"selector,omitempty"`
	Viewports     []Viewport     `json:"viewports,omitempty" yaml:"viewports,omitempty"`
	Threshold     float64        `json:"threshold,omitempty" yaml:"threshold,omitempty"`
	Stabilization *Stabilization `json:"stabilization,omitempty" yaml:"stabilization,omitempty"`
	Skip          bool           `json:"skip,omitempty" yaml:"skip,omitempty"`
	SkipReason    string         `json:"skipReason,omitempty" yaml:"skipReason,omitempty"`
}

// Viewport defines browser viewport dimensions.
type Viewport struct {
	Name   string `json:"name" yaml:"name"`
	Width  int    `json:"width" yaml:"width"`
	Height int    `json:"height" yaml:"height"`
}

// Stabilization defines wait conditions before capturing a screenshot.
type Stabilization struct {
	WaitForSelector   string `json:"waitForSelector,omitempty" yaml:"waitForSelector,omitempty"`
	WaitForTimeout    int    `json:"waitForTimeout,omitempty" yaml:"waitForTimeout,omitempty"`
	WaitMs            int    `json:"waitMs,omitempty" yaml:"waitMs,omitempty"`
	DisableAnimations bool   `json:"disableAnimations,omitempty" yaml:"disableAnimations,omitempty"`
}

// VisualTestReport contains all test results from a test run.
type VisualTestReport struct {
	Timestamp       time.Time          `json:"timestamp"`
	BaselineVersion string             `json:"baselineVersion"`
	Duration        time.Duration      `json:"duration"`
	Summary         VisualTestSummary  `json:"summary"`
	Results         []VisualTestResult `json:"results"`
	Errors          []string           `json:"errors,omitempty"`
}

// VisualTestSummary provides aggregate statistics for a test run.
type VisualTestSummary struct {
	Total   int `json:"total"`
	Passed  int `json:"passed"`
	Failed  int `json:"failed"`
	Skipped int `json:"skipped"`
	Errors  int `json:"errors"`
}

// VisualTestResult represents the outcome of a single visual test.
type VisualTestResult struct {
	TestID       string        `json:"testId"`
	Viewport     string        `json:"viewport"`
	Status       TestStatus    `json:"status"`
	DiffPercent  float64       `json:"diffPercent,omitempty"`
	Threshold    float64       `json:"threshold"`
	Duration     time.Duration `json:"duration"`
	BaselinePath string        `json:"baselinePath,omitempty"`
	ActualPath   string        `json:"actualPath,omitempty"`
	DiffPath     string        `json:"diffPath,omitempty"`
	Error        string        `json:"error,omitempty"`
}

// TestStatus represents the outcome of a test.
type TestStatus string

const (
	// TestStatusPassed indicates the test passed (diff within threshold).
	TestStatusPassed TestStatus = "passed"
	// TestStatusFailed indicates the test failed (diff exceeded threshold).
	TestStatusFailed TestStatus = "failed"
	// TestStatusSkipped indicates the test was skipped.
	TestStatusSkipped TestStatus = "skipped"
	// TestStatusError indicates an error occurred during the test.
	TestStatusError TestStatus = "error"
)

// DefaultViewports returns the standard viewports for testing.
func DefaultViewports() []Viewport {
	return []Viewport{
		{Name: "desktop", Width: 1280, Height: 800},
		{Name: "tablet", Width: 768, Height: 1024},
		{Name: "mobile", Width: 375, Height: 667},
	}
}

// DefaultThreshold returns the default diff threshold (0.1%).
func DefaultThreshold() float64 {
	return 0.001
}

// DefaultStabilization returns default stabilization settings.
func DefaultStabilization() *Stabilization {
	return &Stabilization{
		WaitMs:            100,
		DisableAnimations: true,
	}
}
