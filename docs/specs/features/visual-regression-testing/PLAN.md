# Visual Regression Testing - Implementation Plan

**Version:** 1.0.0
**Status:** Draft
**Author:** PlexusOne
**Last Updated:** 2024-07-11

## Overview

This document outlines the step-by-step implementation plan for adding visual regression testing to DSS.

## Implementation Phases

```
Phase 1: Core Types & Loader     (Week 1)
    │
    ▼
Phase 2: W3Pilot Integration     (Week 1-2)
    │
    ▼
Phase 3: Image Comparison        (Week 2)
    │
    ▼
Phase 4: Baseline Management     (Week 2-3)
    │
    ▼
Phase 5: Test Executor           (Week 3)
    │
    ▼
Phase 6: CLI Commands            (Week 3-4)
    │
    ▼
Phase 7: MCP Tools               (Week 4)
    │
    ▼
Phase 8: Compliance Integration  (Week 4)
    │
    ▼
Phase 9: Documentation & Tests   (Week 4-5)
```

## Phase 1: Core Types & Loader

**Goal:** Define data models and test definition loading.

### Tasks

| Task | File | Description |
|------|------|-------------|
| 1.1 | `sdk/go/visual/types.go` | Define core types (VisualTest, Result, Report) |
| 1.2 | `sdk/go/visual/loader.go` | YAML/JSON test definition loader |
| 1.3 | `sdk/go/visual/loader_test.go` | Loader unit tests |
| 1.4 | `schema/visual-test-suite.schema.json` | JSON Schema for validation |

### Implementation Details

**1.1 types.go**

```go
package visual

import "time"

// VisualTestSuite represents a collection of visual tests
type VisualTestSuite struct {
    Version     string             `json:"version" yaml:"version"`
    Name        string             `json:"name" yaml:"name"`
    Description string             `json:"description,omitempty" yaml:"description,omitempty"`
    BaseURL     string             `json:"baseUrl" yaml:"baseUrl"`
    Defaults    VisualTestDefaults `json:"defaults,omitempty" yaml:"defaults,omitempty"`
    Tests       []VisualTest       `json:"tests" yaml:"tests"`
}

// VisualTestDefaults provides default values for all tests
type VisualTestDefaults struct {
    Viewports     []Viewport     `json:"viewports,omitempty" yaml:"viewports,omitempty"`
    Threshold     float64        `json:"threshold,omitempty" yaml:"threshold,omitempty"`
    Stabilization *Stabilization `json:"stabilization,omitempty" yaml:"stabilization,omitempty"`
}

// VisualTest defines a single visual test case
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

// Viewport defines browser viewport dimensions
type Viewport struct {
    Name   string `json:"name" yaml:"name"`
    Width  int    `json:"width" yaml:"width"`
    Height int    `json:"height" yaml:"height"`
}

// Stabilization defines wait conditions before screenshot
type Stabilization struct {
    WaitForSelector   string `json:"waitForSelector,omitempty" yaml:"waitForSelector,omitempty"`
    WaitForTimeout    int    `json:"waitForTimeout,omitempty" yaml:"waitForTimeout,omitempty"`
    WaitMs            int    `json:"waitMs,omitempty" yaml:"waitMs,omitempty"`
    DisableAnimations bool   `json:"disableAnimations,omitempty" yaml:"disableAnimations,omitempty"`
}

// VisualTestReport contains all test results
type VisualTestReport struct {
    Timestamp       time.Time          `json:"timestamp"`
    BaselineVersion string             `json:"baselineVersion"`
    Duration        time.Duration      `json:"duration"`
    Summary         VisualTestSummary  `json:"summary"`
    Results         []VisualTestResult `json:"results"`
    Errors          []string           `json:"errors,omitempty"`
}

// VisualTestSummary provides aggregate stats
type VisualTestSummary struct {
    Total   int `json:"total"`
    Passed  int `json:"passed"`
    Failed  int `json:"failed"`
    Skipped int `json:"skipped"`
    Errors  int `json:"errors"`
}

// VisualTestResult represents a single test outcome
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

// TestStatus represents the outcome of a test
type TestStatus string

const (
    TestStatusPassed  TestStatus = "passed"
    TestStatusFailed  TestStatus = "failed"
    TestStatusSkipped TestStatus = "skipped"
    TestStatusError   TestStatus = "error"
)
```

**1.2 loader.go**

```go
package visual

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"
)

// Loader loads visual test definitions
type Loader struct {
    basePath string
}

// NewLoader creates a loader for the given path
func NewLoader(basePath string) *Loader {
    return &Loader{basePath: basePath}
}

// Load loads all test suites from the base path
func (l *Loader) Load() (*VisualTestSuite, error) {
    // Look for visual-tests.yaml or visual-tests.json
    yamlPath := filepath.Join(l.basePath, "visual-tests.yaml")
    jsonPath := filepath.Join(l.basePath, "visual-tests.json")

    var data []byte
    var err error

    if _, err = os.Stat(yamlPath); err == nil {
        data, err = os.ReadFile(yamlPath)
    } else if _, err = os.Stat(jsonPath); err == nil {
        data, err = os.ReadFile(jsonPath)
    } else {
        return nil, fmt.Errorf("no visual-tests.yaml or visual-tests.json found in %s", l.basePath)
    }

    if err != nil {
        return nil, fmt.Errorf("failed to read test file: %w", err)
    }

    var suite VisualTestSuite
    if err := yaml.Unmarshal(data, &suite); err != nil {
        return nil, fmt.Errorf("failed to parse test file: %w", err)
    }

    // Apply defaults
    l.applyDefaults(&suite)

    return &suite, nil
}

// applyDefaults merges suite defaults into individual tests
func (l *Loader) applyDefaults(suite *VisualTestSuite) {
    for i := range suite.Tests {
        test := &suite.Tests[i]

        // Apply default viewports
        if len(test.Viewports) == 0 && len(suite.Defaults.Viewports) > 0 {
            test.Viewports = suite.Defaults.Viewports
        }

        // Apply default threshold
        if test.Threshold == 0 && suite.Defaults.Threshold > 0 {
            test.Threshold = suite.Defaults.Threshold
        }

        // Apply default stabilization
        if test.Stabilization == nil && suite.Defaults.Stabilization != nil {
            test.Stabilization = suite.Defaults.Stabilization
        }

        // Expand URL with base URL
        if suite.BaseURL != "" && !isAbsoluteURL(test.URL) {
            test.URL = suite.BaseURL + test.URL
        }
    }
}

func isAbsoluteURL(url string) bool {
    return len(url) > 4 && (url[:4] == "http" || url[:2] == "//")
}
```

### Exit Criteria

- [ ] All types defined with JSON/YAML tags
- [ ] Loader parses YAML and JSON formats
- [ ] Defaults are applied correctly
- [ ] Unit tests pass with 90% coverage

---

## Phase 2: W3Pilot Integration

**Goal:** Create client for w3pilot screenshot capture via MCP.

### Tasks

| Task | File | Description |
|------|------|-------------|
| 2.1 | `sdk/go/visual/capture.go` | W3Pilot MCP client |
| 2.2 | `sdk/go/visual/capture_test.go` | Integration tests |

### Implementation Details

**2.1 capture.go**

```go
package visual

import (
    "context"
    "encoding/base64"
    "fmt"
    "os/exec"
    "time"

    "github.com/modelcontextprotocol/go-sdk/mcp"
)

// CaptureOptions configures screenshot capture
type CaptureOptions struct {
    URL           string
    Selector      string
    Viewport      Viewport
    Stabilization *Stabilization
}

// W3PilotClient wraps MCP client for w3pilot
type W3PilotClient struct {
    client  *mcp.Client
    session *mcp.Session
    cmd     *exec.Cmd
}

// NewW3PilotClient starts w3pilot and connects via MCP
func NewW3PilotClient(ctx context.Context) (*W3PilotClient, error) {
    // Start w3pilot MCP server as subprocess
    cmd := exec.CommandContext(ctx, "w3pilot", "mcp", "serve")

    client := mcp.NewClient(&mcp.ClientInfo{
        Name:    "dss-visual",
        Version: "1.0.0",
    })

    session, err := client.ConnectCommand(ctx, cmd)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to w3pilot: %w", err)
    }

    return &W3PilotClient{
        client:  client,
        session: session,
        cmd:     cmd,
    }, nil
}

// LaunchBrowser starts a browser session
func (c *W3PilotClient) LaunchBrowser(ctx context.Context, headless bool) error {
    _, err := c.session.CallTool(ctx, "browser_launch", map[string]any{
        "headless": headless,
    })
    return err
}

// CaptureScreenshot navigates to URL and captures screenshot
func (c *W3PilotClient) CaptureScreenshot(ctx context.Context, opts CaptureOptions) ([]byte, error) {
    // 1. Navigate
    if _, err := c.session.CallTool(ctx, "page_navigate", map[string]any{
        "url": opts.URL,
    }); err != nil {
        return nil, fmt.Errorf("navigation failed: %w", err)
    }

    // 2. Set viewport
    if _, err := c.session.CallTool(ctx, "page_set_viewport", map[string]any{
        "width":  opts.Viewport.Width,
        "height": opts.Viewport.Height,
    }); err != nil {
        return nil, fmt.Errorf("set viewport failed: %w", err)
    }

    // 3. Stabilization
    if opts.Stabilization != nil {
        if err := c.stabilize(ctx, opts.Stabilization); err != nil {
            return nil, fmt.Errorf("stabilization failed: %w", err)
        }
    }

    // 4. Capture
    var result map[string]any
    var err error

    if opts.Selector != "" {
        result, err = c.session.CallTool(ctx, "element_screenshot", map[string]any{
            "selector": opts.Selector,
        })
    } else {
        result, err = c.session.CallTool(ctx, "page_screenshot", map[string]any{
            "format": "base64",
        })
    }

    if err != nil {
        return nil, fmt.Errorf("screenshot failed: %w", err)
    }

    // 5. Decode
    data, ok := result["data"].(string)
    if !ok {
        return nil, fmt.Errorf("invalid screenshot response")
    }

    return base64.StdEncoding.DecodeString(data)
}

func (c *W3PilotClient) stabilize(ctx context.Context, s *Stabilization) error {
    if s.WaitForSelector != "" {
        timeout := s.WaitForTimeout
        if timeout == 0 {
            timeout = 5000
        }
        if _, err := c.session.CallTool(ctx, "page_wait_for_selector", map[string]any{
            "selector":   s.WaitForSelector,
            "timeout_ms": timeout,
        }); err != nil {
            return err
        }
    }

    if s.WaitMs > 0 {
        time.Sleep(time.Duration(s.WaitMs) * time.Millisecond)
    }

    if s.DisableAnimations {
        // Inject CSS to disable animations
        if _, err := c.session.CallTool(ctx, "page_evaluate", map[string]any{
            "expression": `
                const style = document.createElement('style');
                style.textContent = '*, *::before, *::after { animation: none !important; transition: none !important; }';
                document.head.appendChild(style);
            `,
        }); err != nil {
            return err
        }
    }

    return nil
}

// Close terminates w3pilot subprocess
func (c *W3PilotClient) Close() error {
    if c.session != nil {
        c.session.Close()
    }
    if c.cmd != nil && c.cmd.Process != nil {
        return c.cmd.Process.Kill()
    }
    return nil
}
```

### Exit Criteria

- [ ] W3Pilot client connects via MCP
- [ ] Screenshots capture correctly
- [ ] Viewport settings applied
- [ ] Stabilization works (wait for selector, disable animations)
- [ ] Graceful subprocess cleanup

---

## Phase 3: Image Comparison

**Goal:** Compare images and generate diff visualizations.

### Tasks

| Task | File | Description |
|------|------|-------------|
| 3.1 | `sdk/go/visual/compare.go` | ImageMagick-based comparison |
| 3.2 | `sdk/go/visual/compare_go.go` | Pure Go fallback |
| 3.3 | `sdk/go/visual/compare_test.go` | Comparison tests |

### Implementation Details

**3.1 compare.go (ImageMagick)**

```go
package visual

import (
    "bytes"
    "context"
    "fmt"
    "os/exec"
    "strconv"
    "strings"
)

// Comparator compares images and generates diffs
type Comparator interface {
    Compare(ctx context.Context, baseline, actual, diffOut string) (*CompareResult, error)
}

// CompareResult contains comparison metrics
type CompareResult struct {
    DiffPixels  int64   `json:"diffPixels"`
    TotalPixels int64   `json:"totalPixels"`
    DiffPercent float64 `json:"diffPercent"`
    DiffPath    string  `json:"diffPath,omitempty"`
}

// ImageMagickComparator uses ImageMagick for comparison
type ImageMagickComparator struct {
    comparePath string
    fuzz        string // Anti-aliasing tolerance
}

// NewImageMagickComparator creates a comparator using ImageMagick
func NewImageMagickComparator() (*ImageMagickComparator, error) {
    path, err := exec.LookPath("compare")
    if err != nil {
        return nil, fmt.Errorf("ImageMagick not found: %w", err)
    }
    return &ImageMagickComparator{
        comparePath: path,
        fuzz:        "2%",
    }, nil
}

// Compare compares two images and generates a diff
func (c *ImageMagickComparator) Compare(ctx context.Context, baseline, actual, diffOut string) (*CompareResult, error) {
    // Run ImageMagick compare
    cmd := exec.CommandContext(ctx,
        c.comparePath,
        "-metric", "AE",
        "-fuzz", c.fuzz,
        baseline,
        actual,
        diffOut,
    )

    var stderr bytes.Buffer
    cmd.Stderr = &stderr

    // compare returns exit code 1 if images differ, which is not an error
    cmd.Run()

    // Parse pixel count from stderr
    output := strings.TrimSpace(stderr.String())
    diffPixels, err := strconv.ParseInt(output, 10, 64)
    if err != nil {
        return nil, fmt.Errorf("failed to parse diff output: %w", err)
    }

    // Get image dimensions for percentage
    totalPixels, err := c.getPixelCount(ctx, baseline)
    if err != nil {
        return nil, err
    }

    diffPercent := float64(diffPixels) / float64(totalPixels)

    return &CompareResult{
        DiffPixels:  diffPixels,
        TotalPixels: totalPixels,
        DiffPercent: diffPercent,
        DiffPath:    diffOut,
    }, nil
}

func (c *ImageMagickComparator) getPixelCount(ctx context.Context, image string) (int64, error) {
    cmd := exec.CommandContext(ctx, "identify", "-format", "%w %h", image)
    output, err := cmd.Output()
    if err != nil {
        return 0, err
    }

    parts := strings.Fields(string(output))
    if len(parts) != 2 {
        return 0, fmt.Errorf("unexpected identify output")
    }

    width, _ := strconv.ParseInt(parts[0], 10, 64)
    height, _ := strconv.ParseInt(parts[1], 10, 64)

    return width * height, nil
}
```

### Exit Criteria

- [ ] ImageMagick comparison works
- [ ] Pure Go fallback available
- [ ] Diff percentage calculated correctly
- [ ] Diff images generated with highlights

---

## Phase 4: Baseline Management

**Goal:** Store and manage baseline screenshots per version.

### Tasks

| Task | File | Description |
|------|------|-------------|
| 4.1 | `sdk/go/visual/baseline.go` | Baseline storage and retrieval |
| 4.2 | `sdk/go/visual/baseline_test.go` | Baseline management tests |

### Implementation Details

**4.1 baseline.go**

```go
package visual

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

// BaselineManifest describes a baseline snapshot
type BaselineManifest struct {
    Version   string            `json:"version"`
    CreatedAt time.Time         `json:"createdAt"`
    CreatedBy string            `json:"createdBy,omitempty"`
    TestSuite string            `json:"testSuite"`
    TestCount int               `json:"testCount"`
    Checksums map[string]string `json:"checksums"`
}

// BaselineManager handles baseline storage
type BaselineManager struct {
    basePath string
}

// NewBaselineManager creates a manager for the given directory
func NewBaselineManager(basePath string) *BaselineManager {
    return &BaselineManager{basePath: basePath}
}

// GetVersionPath returns the path for a specific version
func (m *BaselineManager) GetVersionPath(version string) string {
    return filepath.Join(m.basePath, version)
}

// GetBaseline returns the baseline image for a test
func (m *BaselineManager) GetBaseline(version, testID, viewport string) (string, error) {
    filename := m.imageFilename(testID, viewport)
    path := filepath.Join(m.GetVersionPath(version), filename)

    if _, err := os.Stat(path); err != nil {
        return "", fmt.Errorf("baseline not found: %s/%s", testID, viewport)
    }

    return path, nil
}

// SaveBaseline saves a baseline image
func (m *BaselineManager) SaveBaseline(version, testID, viewport string, data []byte) error {
    versionPath := m.GetVersionPath(version)
    if err := os.MkdirAll(versionPath, 0755); err != nil {
        return err
    }

    filename := m.imageFilename(testID, viewport)
    path := filepath.Join(versionPath, filename)

    return os.WriteFile(path, data, 0644)
}

// GenerateManifest creates a manifest for the version
func (m *BaselineManager) GenerateManifest(version string, suite *VisualTestSuite) (*BaselineManifest, error) {
    versionPath := m.GetVersionPath(version)

    manifest := &BaselineManifest{
        Version:   version,
        CreatedAt: time.Now(),
        TestSuite: suite.Name,
        TestCount: len(suite.Tests),
        Checksums: make(map[string]string),
    }

    // Calculate checksums for all baseline images
    for _, test := range suite.Tests {
        for _, vp := range test.Viewports {
            filename := m.imageFilename(test.ID, vp.Name)
            path := filepath.Join(versionPath, filename)

            checksum, err := m.fileChecksum(path)
            if err != nil {
                continue // Skip missing files
            }

            key := fmt.Sprintf("%s/%s", test.ID, vp.Name)
            manifest.Checksums[key] = checksum
        }
    }

    // Save manifest
    manifestPath := filepath.Join(versionPath, "manifest.json")
    data, _ := json.MarshalIndent(manifest, "", "  ")
    if err := os.WriteFile(manifestPath, data, 0644); err != nil {
        return nil, err
    }

    return manifest, nil
}

// LoadManifest loads the manifest for a version
func (m *BaselineManager) LoadManifest(version string) (*BaselineManifest, error) {
    path := filepath.Join(m.GetVersionPath(version), "manifest.json")
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var manifest BaselineManifest
    if err := json.Unmarshal(data, &manifest); err != nil {
        return nil, err
    }

    return &manifest, nil
}

// ListVersions returns all available baseline versions
func (m *BaselineManager) ListVersions() ([]string, error) {
    entries, err := os.ReadDir(m.basePath)
    if err != nil {
        return nil, err
    }

    var versions []string
    for _, entry := range entries {
        if entry.IsDir() && entry.Name() != "latest" {
            versions = append(versions, entry.Name())
        }
    }

    return versions, nil
}

// UpdateLatest updates the "latest" symlink
func (m *BaselineManager) UpdateLatest(version string) error {
    latestPath := filepath.Join(m.basePath, "latest")
    os.Remove(latestPath) // Remove existing symlink
    return os.Symlink(version, latestPath)
}

func (m *BaselineManager) imageFilename(testID, viewport string) string {
    return fmt.Sprintf("%s-%s.png", testID, viewport)
}

func (m *BaselineManager) fileChecksum(path string) (string, error) {
    f, err := os.Open(path)
    if err != nil {
        return "", err
    }
    defer f.Close()

    h := sha256.New()
    if _, err := io.Copy(h, f); err != nil {
        return "", err
    }

    return hex.EncodeToString(h.Sum(nil)), nil
}
```

### Exit Criteria

- [ ] Baselines stored by version
- [ ] Manifest tracks checksums
- [ ] Latest symlink works
- [ ] Version listing works

---

## Phase 5: Test Executor

**Goal:** Orchestrate test execution with parallelization.

### Tasks

| Task | File | Description |
|------|------|-------------|
| 5.1 | `sdk/go/visual/executor.go` | Test execution engine |
| 5.2 | `sdk/go/visual/service.go` | Service layer (public API) |
| 5.3 | `sdk/go/visual/executor_test.go` | Executor tests |

### Implementation Details

**5.1 executor.go**

```go
package visual

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"
)

// Executor runs visual tests
type Executor struct {
    loader     *Loader
    baseline   *BaselineManager
    comparator Comparator
    outputDir  string
    parallel   int
}

// NewExecutor creates a test executor
func NewExecutor(opts ExecutorOptions) (*Executor, error) {
    comparator, err := NewImageMagickComparator()
    if err != nil {
        return nil, err
    }

    return &Executor{
        loader:     NewLoader(opts.TestsDir),
        baseline:   NewBaselineManager(opts.BaselinesDir),
        comparator: comparator,
        outputDir:  opts.OutputDir,
        parallel:   opts.Parallel,
    }, nil
}

// ExecutorOptions configures the executor
type ExecutorOptions struct {
    TestsDir     string
    BaselinesDir string
    OutputDir    string
    Parallel     int
}

// Run executes all tests
func (e *Executor) Run(ctx context.Context, baselineVersion string, testIDs []string) (*VisualTestReport, error) {
    startTime := time.Now()

    // Load test suite
    suite, err := e.loader.Load()
    if err != nil {
        return nil, err
    }

    // Filter tests if specific IDs provided
    tests := suite.Tests
    if len(testIDs) > 0 {
        tests = e.filterTests(tests, testIDs)
    }

    // Ensure output directory exists
    if err := os.MkdirAll(e.outputDir, 0755); err != nil {
        return nil, err
    }

    // Run tests in parallel
    results := e.runParallel(ctx, tests, baselineVersion)

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

func (e *Executor) runParallel(ctx context.Context, tests []VisualTest, baselineVersion string) []VisualTestResult {
    // Create job channel
    type job struct {
        test     VisualTest
        viewport Viewport
    }

    jobs := make(chan job, len(tests)*3) // Assume max 3 viewports
    results := make(chan VisualTestResult, len(tests)*3)

    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < e.parallel; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()

            // Each worker gets its own w3pilot client
            client, err := NewW3PilotClient(ctx)
            if err != nil {
                return
            }
            defer client.Close()

            // Launch browser once per worker
            client.LaunchBrowser(ctx, true)

            for j := range jobs {
                result := e.runSingleTest(ctx, client, j.test, j.viewport, baselineVersion)
                results <- result
            }
        }()
    }

    // Send jobs
    for _, test := range tests {
        if test.Skip {
            for _, vp := range test.Viewports {
                results <- VisualTestResult{
                    TestID:   test.ID,
                    Viewport: vp.Name,
                    Status:   TestStatusSkipped,
                }
            }
            continue
        }

        for _, vp := range test.Viewports {
            jobs <- job{test: test, viewport: vp}
        }
    }
    close(jobs)

    // Wait and collect
    wg.Wait()
    close(results)

    var allResults []VisualTestResult
    for r := range results {
        allResults = append(allResults, r)
    }

    return allResults
}

func (e *Executor) runSingleTest(ctx context.Context, client *W3PilotClient, test VisualTest, vp Viewport, baselineVersion string) VisualTestResult {
    startTime := time.Now()
    result := VisualTestResult{
        TestID:    test.ID,
        Viewport:  vp.Name,
        Threshold: test.Threshold,
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
    actualPath := filepath.Join(e.outputDir, fmt.Sprintf("%s-%s.png", test.ID, vp.Name))
    if err := os.WriteFile(actualPath, screenshot, 0644); err != nil {
        result.Status = TestStatusError
        result.Error = err.Error()
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
        result.Error = err.Error()
        result.Duration = time.Since(startTime)
        return result
    }

    result.DiffPercent = compareResult.DiffPercent
    result.DiffPath = diffPath
    result.Duration = time.Since(startTime)

    // 5. Determine pass/fail
    if compareResult.DiffPercent <= test.Threshold {
        result.Status = TestStatusPassed
        os.Remove(diffPath) // Clean up diff for passing tests
        result.DiffPath = ""
    } else {
        result.Status = TestStatusFailed
    }

    return result
}

func (e *Executor) filterTests(tests []VisualTest, ids []string) []VisualTest {
    idSet := make(map[string]bool)
    for _, id := range ids {
        idSet[id] = true
    }

    var filtered []VisualTest
    for _, t := range tests {
        if idSet[t.ID] {
            filtered = append(filtered, t)
        }
    }

    return filtered
}
```

### Exit Criteria

- [ ] Parallel execution works
- [ ] Worker pool manages w3pilot instances
- [ ] Results aggregated correctly
- [ ] Graceful error handling

---

## Phase 6: CLI Commands

**Goal:** Add `dss visual-test` and `dss visual-baseline` commands.

### Tasks

| Task | File | Description |
|------|------|-------------|
| 6.1 | `cmd/dss/cmd/visual_test.go` | Test runner command |
| 6.2 | `cmd/dss/cmd/visual_baseline.go` | Baseline management commands |
| 6.3 | `cmd/dss/cmd/visual_test_test.go` | CLI tests |

### Exit Criteria

- [ ] `dss visual-test` runs tests
- [ ] `dss visual-baseline generate` creates baselines
- [ ] `dss visual-baseline update` updates specific tests
- [ ] `dss visual-baseline list` shows versions
- [ ] JSON output format works

---

## Phase 7: MCP Tools

**Goal:** Expose visual testing via MCP for AI agents.

### Tasks

| Task | File | Description |
|------|------|-------------|
| 7.1 | `skills/designsystem/tools_visual.go` | Visual test MCP tools |

### MCP Tools

| Tool | Description |
|------|-------------|
| `visual_test` | Run visual tests against baseline |
| `visual_baseline_generate` | Generate baselines for version |
| `visual_baseline_update` | Update specific test baselines |
| `visual_test_single` | Run single test (debugging) |

### Exit Criteria

- [ ] All MCP tools implemented
- [ ] Tools integrated with dss-mcp server
- [ ] JSON responses match CLI output

---

## Phase 8: Compliance Integration

**Goal:** Include visual testing in compliance reports.

### Tasks

| Task | File | Description |
|------|------|-------------|
| 8.1 | `sdk/go/compliance.go` | Add visual category to compliance |
| 8.2 | `skills/designsystem/tools_compliance.go` | Update compliance tools |

### Exit Criteria

- [ ] Visual test results in compliance report
- [ ] Visual failures block release gate (configurable)
- [ ] Agentic maturity updated to L4

---

## Phase 9: Documentation & Tests

**Goal:** Complete documentation and test coverage.

### Tasks

| Task | File | Description |
|------|------|-------------|
| 9.1 | `docs/visual-testing.md` | User documentation |
| 9.2 | `docs/mcp-server.md` | Update with visual tools |
| 9.3 | `sdk/go/visual/*_test.go` | Unit tests (90% coverage) |
| 9.4 | `examples/visual-tests/` | Example test suite |

### Exit Criteria

- [ ] User documentation complete
- [ ] MCP tools documented
- [ ] 90% test coverage
- [ ] Example tests work

---

## Testing Strategy

### Unit Tests

- Loader parsing
- Baseline management
- Image comparison (mock images)
- Result aggregation

### Integration Tests

- W3Pilot communication (real subprocess)
- Screenshot capture (real browser)
- End-to-end test run

### CI Tests

- Run on PRs
- Use baseline from main branch
- Upload diff artifacts

---

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| W3Pilot unavailable | Clear error message, skip visual tests |
| ImageMagick not installed | Pure Go fallback |
| Flaky tests | Retry logic, stabilization options |
| Large baselines | Compression, pruning old versions |
| Slow execution | Parallel workers, test filtering |

---

## Dependencies

### External

| Dependency | Required | Fallback |
|------------|----------|----------|
| w3pilot | Yes | Error |
| ImageMagick | Yes | go-imagediff |
| Chrome/Chromium | Yes | Firefox |

### Go Packages

```go
require (
    github.com/modelcontextprotocol/go-sdk v1.6.0
    gopkg.in/yaml.v3 v3.0.1
    github.com/spf13/cobra v1.8.0
)
```

---

## Success Criteria

| Metric | Target |
|--------|--------|
| Test coverage | 90% |
| CLI commands | 4 |
| MCP tools | 4 |
| Example tests | 10+ |
| Documentation pages | 2 |
| Agentic maturity | L4 (from L3) |
