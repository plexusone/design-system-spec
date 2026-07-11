package visual

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Executor runs visual tests.
type Executor struct {
	testsDir     string
	baselinesDir string
	outputDir    string
	parallel     int
	loader       *Loader
	baseline     *BaselineManager
	comparator   Comparator
}

// ExecutorOptions configures the executor.
type ExecutorOptions struct {
	TestsDir     string // Directory containing visual-tests.yaml
	BaselinesDir string // Directory containing baseline versions
	OutputDir    string // Directory for test results
	Parallel     int    // Number of parallel workers
}

// NewExecutor creates a test executor.
func NewExecutor(opts ExecutorOptions) (*Executor, error) {
	if opts.Parallel <= 0 {
		opts.Parallel = 4
	}

	comparator, err := NewComparator()
	if err != nil {
		return nil, fmt.Errorf("failed to create comparator: %w", err)
	}

	return &Executor{
		testsDir:     opts.TestsDir,
		baselinesDir: opts.BaselinesDir,
		outputDir:    opts.OutputDir,
		parallel:     opts.Parallel,
		loader:       NewLoader(opts.TestsDir),
		baseline:     NewBaselineManager(opts.BaselinesDir),
		comparator:   comparator,
	}, nil
}

// RunOptions configures a test run.
type RunOptions struct {
	BaselineVersion string   // Version to compare against (or "latest")
	TestIDs         []string // Specific tests to run (empty = all)
	Viewports       []string // Specific viewports (empty = all)
	Threshold       float64  // Override threshold (0 = use test default)
	UpdateBaseline  bool     // Update baseline instead of comparing
}

// Run executes visual tests and returns a report.
func (e *Executor) Run(ctx context.Context, opts RunOptions) (*VisualTestReport, error) {
	startTime := time.Now()

	// Load test suite
	suite, err := e.loader.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load test suite: %w", err)
	}

	// Resolve baseline version
	baselineVersion := opts.BaselineVersion
	if baselineVersion == "" || baselineVersion == "latest" {
		baselineVersion, err = e.baseline.GetLatestVersion()
		if err != nil && !opts.UpdateBaseline {
			return nil, fmt.Errorf("no baseline available: %w", err)
		}
	}

	// Filter tests
	tests := suite.Tests
	if len(opts.TestIDs) > 0 {
		tests = FilterTests(tests, opts.TestIDs)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(e.outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Run tests
	var results []VisualTestResult
	if opts.UpdateBaseline {
		results = e.runBaseline(ctx, tests, baselineVersion)
	} else {
		results = e.runParallel(ctx, tests, baselineVersion, opts)
	}

	// Build report
	report := &VisualTestReport{
		Timestamp:       startTime,
		BaselineVersion: baselineVersion,
		Duration:        time.Since(startTime),
		Results:         results,
	}

	// Calculate summary
	for _, r := range results {
		report.Summary.Total++
		switch r.Status {
		case TestStatusPassed:
			report.Summary.Passed++
		case TestStatusFailed:
			report.Summary.Failed++
		case TestStatusSkipped:
			report.Summary.Skipped++
		case TestStatusError:
			report.Summary.Errors++
		}
	}

	return report, nil
}

// runParallel runs tests in parallel with worker pool.
func (e *Executor) runParallel(ctx context.Context, tests []VisualTest, baselineVersion string, opts RunOptions) []VisualTestResult {
	// Expand tests into individual jobs
	type job struct {
		test     VisualTest
		viewport Viewport
	}

	var jobs []job
	for _, test := range tests {
		if test.Skip {
			continue
		}
		for _, vp := range test.Viewports {
			// Filter by viewport if specified
			if len(opts.Viewports) > 0 && !contains(opts.Viewports, vp.Name) {
				continue
			}
			jobs = append(jobs, job{test: test, viewport: vp})
		}
	}

	// Add skipped tests to results
	var results []VisualTestResult
	var mu sync.Mutex

	for _, test := range tests {
		if test.Skip {
			for _, vp := range test.Viewports {
				results = append(results, VisualTestResult{
					TestID:   test.ID,
					Viewport: vp.Name,
					Status:   TestStatusSkipped,
					Error:    test.SkipReason,
				})
			}
		}
	}

	if len(jobs) == 0 {
		return results
	}

	// Create job channel
	jobCh := make(chan job, len(jobs))
	for _, j := range jobs {
		jobCh <- j
	}
	close(jobCh)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < e.parallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Each worker gets its own w3pilot client
			client, err := NewW3PilotClient(ctx)
			if err != nil {
				// Log error and skip this worker
				return
			}
			defer client.Close()

			// Launch browser
			if err := client.LaunchBrowser(ctx, true); err != nil {
				return
			}
			defer func() { _ = client.CloseBrowser(ctx) }()

			for j := range jobCh {
				threshold := j.test.Threshold
				if opts.Threshold > 0 {
					threshold = opts.Threshold
				}

				result := e.runSingleTest(ctx, client, j.test, j.viewport, baselineVersion, threshold)

				mu.Lock()
				results = append(results, result)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return results
}

// runSingleTest executes a single test case.
func (e *Executor) runSingleTest(ctx context.Context, client *W3PilotClient, test VisualTest, vp Viewport, baselineVersion string, threshold float64) VisualTestResult {
	startTime := time.Now()
	result := VisualTestResult{
		TestID:    test.ID,
		Viewport:  vp.Name,
		Threshold: threshold,
	}

	// 1. Capture screenshot
	screenshot, err := client.CaptureScreenshot(ctx, CaptureOptions{
		URL:           test.URL,
		Selector:      test.Selector,
		Viewport:      vp,
		Stabilization: test.Stabilization,
	})
	if err != nil {
		result.Status = TestStatusError
		result.Error = err.Error()
		result.Duration = time.Since(startTime)
		return result
	}

	// 2. Save actual screenshot
	actualPath := filepath.Join(e.outputDir, imageFilename(test.ID, vp.Name))
	//nolint:gosec // G306: Screenshot files are intentionally world-readable
	if err := os.WriteFile(actualPath, screenshot, 0644); err != nil {
		result.Status = TestStatusError
		result.Error = fmt.Sprintf("failed to save screenshot: %v", err)
		result.Duration = time.Since(startTime)
		return result
	}
	result.ActualPath = actualPath

	// 3. Get baseline
	baselinePath, err := e.baseline.GetBaseline(baselineVersion, test.ID, vp.Name)
	if err != nil {
		result.Status = TestStatusError
		result.Error = err.Error()
		result.Duration = time.Since(startTime)
		return result
	}
	result.BaselinePath = baselinePath

	// 4. Compare
	diffPath := filepath.Join(e.outputDir, fmt.Sprintf("%s-%s.diff.png", test.ID, vp.Name))
	compareResult, err := e.comparator.Compare(ctx, baselinePath, actualPath, diffPath)
	if err != nil {
		result.Status = TestStatusError
		result.Error = fmt.Sprintf("comparison failed: %v", err)
		result.Duration = time.Since(startTime)
		return result
	}

	result.DiffPercent = compareResult.DiffPercent
	result.Duration = time.Since(startTime)

	// 5. Determine pass/fail
	if compareResult.DiffPercent <= threshold {
		result.Status = TestStatusPassed
		// Clean up diff for passing tests
		os.Remove(diffPath)
	} else {
		result.Status = TestStatusFailed
		result.DiffPath = diffPath
	}

	return result
}

// runBaseline generates baselines instead of comparing.
func (e *Executor) runBaseline(ctx context.Context, tests []VisualTest, version string) []VisualTestResult {
	var results []VisualTestResult

	// Create single client for baseline generation
	client, err := NewW3PilotClient(ctx)
	if err != nil {
		for _, test := range tests {
			for _, vp := range test.Viewports {
				results = append(results, VisualTestResult{
					TestID:   test.ID,
					Viewport: vp.Name,
					Status:   TestStatusError,
					Error:    err.Error(),
				})
			}
		}
		return results
	}
	defer client.Close()

	if err := client.LaunchBrowser(ctx, true); err != nil {
		for _, test := range tests {
			for _, vp := range test.Viewports {
				results = append(results, VisualTestResult{
					TestID:   test.ID,
					Viewport: vp.Name,
					Status:   TestStatusError,
					Error:    fmt.Sprintf("failed to launch browser: %v", err),
				})
			}
		}
		return results
	}
	defer func() { _ = client.CloseBrowser(ctx) }()

	for _, test := range tests {
		if test.Skip {
			for _, vp := range test.Viewports {
				results = append(results, VisualTestResult{
					TestID:   test.ID,
					Viewport: vp.Name,
					Status:   TestStatusSkipped,
					Error:    test.SkipReason,
				})
			}
			continue
		}

		for _, vp := range test.Viewports {
			startTime := time.Now()
			result := VisualTestResult{
				TestID:    test.ID,
				Viewport:  vp.Name,
				Threshold: test.Threshold,
			}

			// Capture screenshot
			screenshot, err := client.CaptureScreenshot(ctx, CaptureOptions{
				URL:           test.URL,
				Selector:      test.Selector,
				Viewport:      vp,
				Stabilization: test.Stabilization,
			})
			if err != nil {
				result.Status = TestStatusError
				result.Error = err.Error()
				result.Duration = time.Since(startTime)
				results = append(results, result)
				continue
			}

			// Save as baseline
			if err := e.baseline.SaveBaseline(version, test.ID, vp.Name, screenshot); err != nil {
				result.Status = TestStatusError
				result.Error = fmt.Sprintf("failed to save baseline: %v", err)
			} else {
				result.Status = TestStatusPassed
				result.BaselinePath = filepath.Join(e.baseline.GetVersionPath(version), imageFilename(test.ID, vp.Name))
			}

			result.Duration = time.Since(startTime)
			results = append(results, result)
		}
	}

	return results
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
