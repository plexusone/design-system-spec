package visual

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Service provides the public API for visual testing.
type Service struct {
	testsDir     string
	baselinesDir string
	outputDir    string
	parallel     int
}

// ServiceOptions configures the visual testing service.
type ServiceOptions struct {
	TestsDir     string // Directory containing visual-tests.yaml
	BaselinesDir string // Directory containing baseline versions
	OutputDir    string // Directory for test results
	Parallel     int    // Number of parallel workers (default: 4)
}

// NewService creates a visual testing service.
func NewService(opts ServiceOptions) *Service {
	if opts.Parallel <= 0 {
		opts.Parallel = 4
	}
	if opts.OutputDir == "" {
		opts.OutputDir = filepath.Join(opts.TestsDir, "results")
	}

	return &Service{
		testsDir:     opts.TestsDir,
		baselinesDir: opts.BaselinesDir,
		outputDir:    opts.OutputDir,
		parallel:     opts.Parallel,
	}
}

// LoadTestSuite loads the visual test suite.
func (s *Service) LoadTestSuite() (*VisualTestSuite, error) {
	loader := NewLoader(s.testsDir)
	return loader.Load()
}

// RunTests executes visual tests and returns a report.
func (s *Service) RunTests(ctx context.Context, opts *TestOptions) (*VisualTestReport, error) {
	if opts == nil {
		opts = &TestOptions{}
	}

	executor, err := NewExecutor(ExecutorOptions{
		TestsDir:     s.testsDir,
		BaselinesDir: s.baselinesDir,
		OutputDir:    s.outputDir,
		Parallel:     s.parallel,
	})
	if err != nil {
		return nil, err
	}

	return executor.Run(ctx, RunOptions{
		BaselineVersion: opts.BaselineVersion,
		TestIDs:         opts.TestIDs,
		Viewports:       opts.Viewports,
		Threshold:       opts.Threshold,
		UpdateBaseline:  false,
	})
}

// TestOptions configures a test run.
type TestOptions struct {
	BaselineVersion string   // Version to compare against (default: "latest")
	TestIDs         []string // Specific tests to run (empty = all)
	Viewports       []string // Specific viewports (empty = all)
	Threshold       float64  // Override threshold (0 = use test default)
}

// GenerateBaseline generates baselines for a new version.
func (s *Service) GenerateBaseline(ctx context.Context, version string) (*BaselineResult, error) {
	if version == "" {
		return nil, fmt.Errorf("version is required")
	}

	executor, err := NewExecutor(ExecutorOptions{
		TestsDir:     s.testsDir,
		BaselinesDir: s.baselinesDir,
		OutputDir:    s.outputDir,
		Parallel:     1, // Single worker for baseline generation
	})
	if err != nil {
		return nil, err
	}

	// Run in baseline mode
	report, err := executor.Run(ctx, RunOptions{
		BaselineVersion: version,
		UpdateBaseline:  true,
	})
	if err != nil {
		return nil, err
	}

	// Load test suite for manifest
	suite, err := s.LoadTestSuite()
	if err != nil {
		return nil, err
	}

	// Generate manifest
	baseline := NewBaselineManager(s.baselinesDir)
	manifest, err := baseline.GenerateManifest(version, suite)
	if err != nil {
		return nil, fmt.Errorf("failed to generate manifest: %w", err)
	}

	// Update latest symlink
	if err := baseline.UpdateLatest(version); err != nil {
		// Non-fatal
	}

	return &BaselineResult{
		Version:   version,
		TestCount: manifest.TestCount,
		Path:      baseline.GetVersionPath(version),
		Errors:    countErrors(report),
	}, nil
}

// BaselineResult contains the result of baseline generation.
type BaselineResult struct {
	Version   string `json:"version"`
	TestCount int    `json:"testCount"`
	Path      string `json:"path"`
	Errors    int    `json:"errors"`
}

// UpdateBaseline updates specific test baselines.
func (s *Service) UpdateBaseline(ctx context.Context, version string, testIDs []string) (*BaselineResult, error) {
	if version == "" {
		return nil, fmt.Errorf("version is required")
	}
	if len(testIDs) == 0 {
		return nil, fmt.Errorf("at least one test ID is required")
	}

	executor, err := NewExecutor(ExecutorOptions{
		TestsDir:     s.testsDir,
		BaselinesDir: s.baselinesDir,
		OutputDir:    s.outputDir,
		Parallel:     1,
	})
	if err != nil {
		return nil, err
	}

	report, err := executor.Run(ctx, RunOptions{
		BaselineVersion: version,
		TestIDs:         testIDs,
		UpdateBaseline:  true,
	})
	if err != nil {
		return nil, err
	}

	// Regenerate manifest
	suite, err := s.LoadTestSuite()
	if err != nil {
		return nil, err
	}

	baseline := NewBaselineManager(s.baselinesDir)
	manifest, err := baseline.GenerateManifest(version, suite)
	if err != nil {
		return nil, err
	}

	return &BaselineResult{
		Version:   version,
		TestCount: manifest.TestCount,
		Path:      baseline.GetVersionPath(version),
		Errors:    countErrors(report),
	}, nil
}

// ListBaselineVersions returns all available baseline versions.
func (s *Service) ListBaselineVersions() ([]string, error) {
	baseline := NewBaselineManager(s.baselinesDir)
	return baseline.ListVersions()
}

// GetBaselineManifest returns the manifest for a baseline version.
func (s *Service) GetBaselineManifest(version string) (*BaselineManifest, error) {
	baseline := NewBaselineManager(s.baselinesDir)
	return baseline.LoadManifest(version)
}

// RunSingleTest runs a single test for debugging.
func (s *Service) RunSingleTest(ctx context.Context, testID string, opts *SingleTestOptions) (*VisualTestResult, error) {
	if opts == nil {
		opts = &SingleTestOptions{}
	}

	// Load test suite to find the test
	suite, err := s.LoadTestSuite()
	if err != nil {
		return nil, err
	}

	var test *VisualTest
	for i := range suite.Tests {
		if suite.Tests[i].ID == testID {
			test = &suite.Tests[i]
			break
		}
	}

	if test == nil {
		return nil, fmt.Errorf("%w: %s", ErrTestNotFound, testID)
	}

	// Determine viewport
	var viewport Viewport
	if opts.Viewport != "" {
		for _, vp := range test.Viewports {
			if vp.Name == opts.Viewport {
				viewport = vp
				break
			}
		}
		if viewport.Name == "" {
			return nil, fmt.Errorf("viewport %s not found for test %s", opts.Viewport, testID)
		}
	} else if len(test.Viewports) > 0 {
		viewport = test.Viewports[0]
	} else {
		viewport = DefaultViewports()[0]
	}

	// Create client
	client, err := NewW3PilotClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	if err := client.LaunchBrowser(ctx, !opts.Headful); err != nil {
		return nil, err
	}
	defer func() { _ = client.CloseBrowser(ctx) }()

	// Create executor for single test
	executor, err := NewExecutor(ExecutorOptions{
		TestsDir:     s.testsDir,
		BaselinesDir: s.baselinesDir,
		OutputDir:    s.outputDir,
		Parallel:     1,
	})
	if err != nil {
		return nil, err
	}

	threshold := test.Threshold
	if opts.Threshold > 0 {
		threshold = opts.Threshold
	}

	baselineVersion := opts.BaselineVersion
	if baselineVersion == "" {
		baseline := NewBaselineManager(s.baselinesDir)
		baselineVersion, _ = baseline.GetLatestVersion()
	}

	result := executor.runSingleTest(ctx, client, *test, viewport, baselineVersion, threshold)
	return &result, nil
}

// SingleTestOptions configures a single test run.
type SingleTestOptions struct {
	BaselineVersion string  // Version to compare against
	Viewport        string  // Specific viewport to test
	Threshold       float64 // Override threshold
	Headful         bool    // Run browser in headful mode (visible)
}

// SaveReport saves a test report to a file.
func (s *Service) SaveReport(report *VisualTestReport, path string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	//nolint:gosec // G306: Report files are intentionally world-readable
	return os.WriteFile(path, data, 0644)
}

// CleanResults removes old test results.
func (s *Service) CleanResults() error {
	return os.RemoveAll(s.outputDir)
}

func countErrors(report *VisualTestReport) int {
	count := 0
	for _, r := range report.Results {
		if r.Status == TestStatusError {
			count++
		}
	}
	return count
}
