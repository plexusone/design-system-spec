# Visual Regression Testing - Technical Requirements Document

**Version:** 1.0.0
**Status:** Draft
**Author:** PlexusOne
**Last Updated:** 2024-07-11

## Overview

This document defines the technical architecture and implementation details for visual regression testing in DSS, using w3pilot as the screenshot capture engine.

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              DSS Visual Testing                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐                  │
│  │ CLI Commands │    │  MCP Tools   │    │   Go SDK     │                  │
│  │              │    │              │    │              │                  │
│  │ visual-test  │    │ visual_test  │    │ RunVisualTests│                 │
│  │ visual-base  │    │ visual_base* │    │ GenerateBase │                  │
│  └──────┬───────┘    └──────┬───────┘    └──────┬───────┘                  │
│         │                   │                   │                           │
│         └───────────────────┴───────────────────┘                           │
│                             │                                                │
│                             ▼                                                │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │                     Visual Test Service                               │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │  │
│  │  │ Test Loader │  │  Executor   │  │  Comparator │  │  Reporter   │ │  │
│  │  │             │  │             │  │             │  │             │ │  │
│  │  │ Parse YAML  │  │ Run w3pilot │  │ ImageMagick │  │ JSON/HTML   │ │  │
│  │  │ Validate    │  │ Screenshots │  │ Diff Gen    │  │ Reports     │ │  │
│  │  └─────────────┘  └──────┬──────┘  └─────────────┘  └─────────────┘ │  │
│  └───────────────────────────┼──────────────────────────────────────────┘  │
│                              │                                              │
└──────────────────────────────┼──────────────────────────────────────────────┘
                               │
                               ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│                            w3pilot (Subprocess)                               │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐              │
│  │ Browser Control │  │ Page Navigation │  │ Screenshot Cap  │              │
│  │                 │  │                 │  │                 │              │
│  │ Chrome/Firefox  │  │ Set Viewport    │  │ Full Page       │              │
│  │ Headless        │  │ Wait for Ready  │  │ Element         │              │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘              │
└──────────────────────────────────────────────────────────────────────────────┘
```

### Component Architecture

```
sdk/go/
├── visual/                      # Visual testing package
│   ├── types.go                 # Core types (Test, Result, Report)
│   ├── loader.go                # Test definition loader
│   ├── executor.go              # Test execution orchestration
│   ├── capture.go               # w3pilot integration
│   ├── compare.go               # Image comparison
│   ├── baseline.go              # Baseline management
│   ├── report.go                # Report generation
│   └── service.go               # Service layer (public API)
│
├── visual_test.go               # Package tests
│
skills/designsystem/
├── tools_visual.go              # MCP tools for visual testing
│
cmd/dss/cmd/
├── visual_test.go               # CLI: dss visual-test
├── visual_baseline.go           # CLI: dss visual-baseline
```

## Data Models

### Test Definition

```go
// visual/types.go

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
    WaitForSelector string `json:"waitForSelector,omitempty" yaml:"waitForSelector,omitempty"`
    WaitForTimeout  int    `json:"waitForTimeout,omitempty" yaml:"waitForTimeout,omitempty"`
    WaitMs          int    `json:"waitMs,omitempty" yaml:"waitMs,omitempty"`
    DisableAnimations bool `json:"disableAnimations,omitempty" yaml:"disableAnimations,omitempty"`
}
```

### Baseline

```go
// visual/baseline.go

// BaselineManifest describes a baseline snapshot
type BaselineManifest struct {
    Version     string            `json:"version"`
    CreatedAt   time.Time         `json:"createdAt"`
    CreatedBy   string            `json:"createdBy,omitempty"`
    TestSuite   string            `json:"testSuite"`
    TestCount   int               `json:"testCount"`
    Checksums   map[string]string `json:"checksums"` // testID -> SHA256
    Metadata    map[string]string `json:"metadata,omitempty"`
}

// BaselineEntry represents a single baseline image
type BaselineEntry struct {
    TestID    string    `json:"testId"`
    Viewport  string    `json:"viewport"`
    Path      string    `json:"path"`
    Checksum  string    `json:"checksum"`
    CreatedAt time.Time `json:"createdAt"`
    Size      int64     `json:"size"`
}
```

### Test Results

```go
// visual/types.go

// VisualTestReport contains all test results
type VisualTestReport struct {
    Timestamp       time.Time           `json:"timestamp"`
    BaselineVersion string              `json:"baselineVersion"`
    Duration        time.Duration       `json:"duration"`
    Summary         VisualTestSummary   `json:"summary"`
    Results         []VisualTestResult  `json:"results"`
    Errors          []string            `json:"errors,omitempty"`
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
    Status       TestStatus    `json:"status"` // passed, failed, skipped, error
    DiffPercent  float64       `json:"diffPercent,omitempty"`
    Threshold    float64       `json:"threshold"`
    Duration     time.Duration `json:"duration"`
    BaselinePath string        `json:"baselinePath,omitempty"`
    ActualPath   string        `json:"actualPath,omitempty"`
    DiffPath     string        `json:"diffPath,omitempty"`
    Error        string        `json:"error,omitempty"`
}

type TestStatus string

const (
    TestStatusPassed  TestStatus = "passed"
    TestStatusFailed  TestStatus = "failed"
    TestStatusSkipped TestStatus = "skipped"
    TestStatusError   TestStatus = "error"
)
```

## W3Pilot Integration

### Communication Protocol

DSS communicates with w3pilot via MCP (Model Context Protocol) over stdio:

```go
// visual/capture.go

// W3PilotClient wraps MCP client for w3pilot communication
type W3PilotClient struct {
    session *mcp.Session
    cmd     *exec.Cmd
}

// NewW3PilotClient starts w3pilot subprocess and connects via MCP
func NewW3PilotClient(ctx context.Context) (*W3PilotClient, error) {
    cmd := exec.Command("w3pilot", "mcp", "serve")

    // Connect MCP client to subprocess stdio
    client := mcp.NewClient("dss-visual", "1.0.0")
    session, err := client.ConnectCommand(ctx, cmd, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to w3pilot: %w", err)
    }

    return &W3PilotClient{
        session: session,
        cmd:     cmd,
    }, nil
}

// CaptureScreenshot takes a screenshot of the specified element or page
func (c *W3PilotClient) CaptureScreenshot(ctx context.Context, opts CaptureOptions) ([]byte, error) {
    // 1. Navigate to URL
    _, err := c.session.CallTool(ctx, "page_navigate", map[string]any{
        "url": opts.URL,
    })
    if err != nil {
        return nil, fmt.Errorf("navigation failed: %w", err)
    }

    // 2. Set viewport
    _, err = c.session.CallTool(ctx, "page_set_viewport", map[string]any{
        "width":  opts.Viewport.Width,
        "height": opts.Viewport.Height,
    })
    if err != nil {
        return nil, fmt.Errorf("viewport set failed: %w", err)
    }

    // 3. Wait for stabilization
    if opts.Stabilization != nil {
        if opts.Stabilization.WaitForSelector != "" {
            _, err = c.session.CallTool(ctx, "page_wait_for_selector", map[string]any{
                "selector":   opts.Stabilization.WaitForSelector,
                "timeout_ms": opts.Stabilization.WaitForTimeout,
            })
            if err != nil {
                return nil, fmt.Errorf("wait for selector failed: %w", err)
            }
        }
        if opts.Stabilization.WaitMs > 0 {
            time.Sleep(time.Duration(opts.Stabilization.WaitMs) * time.Millisecond)
        }
    }

    // 4. Capture screenshot
    var result map[string]any
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

    // 5. Decode base64 response
    data, ok := result["data"].(string)
    if !ok {
        return nil, fmt.Errorf("invalid screenshot response")
    }

    return base64.StdEncoding.DecodeString(data)
}

// Close terminates the w3pilot subprocess
func (c *W3PilotClient) Close() error {
    c.session.Close()
    return c.cmd.Process.Kill()
}
```

### W3Pilot MCP Tools Used

| Tool | Purpose | Parameters |
|------|---------|------------|
| `browser_launch` | Start browser session | `headless`, `browser` |
| `page_navigate` | Load test URL | `url` |
| `page_set_viewport` | Set viewport size | `width`, `height` |
| `page_wait_for_selector` | Wait for element | `selector`, `timeout_ms` |
| `page_screenshot` | Full page capture | `format` |
| `element_screenshot` | Element capture | `selector` |
| `browser_close` | Close browser | - |

## Image Comparison

### Comparison Algorithm

Use ImageMagick `compare` for pixel-level diffing:

```go
// visual/compare.go

// ImageComparator handles image comparison operations
type ImageComparator struct {
    magickPath string // Path to ImageMagick compare binary
}

// CompareImages compares two images and generates a diff
func (c *ImageComparator) CompareImages(ctx context.Context, baseline, actual, diffOut string) (*CompareResult, error) {
    // Use ImageMagick compare with metric
    cmd := exec.CommandContext(ctx,
        c.magickPath,
        "-metric", "AE",           // Absolute Error (pixel count)
        "-fuzz", "2%",             // Anti-aliasing tolerance
        baseline,
        actual,
        diffOut,
    )

    // ImageMagick outputs metric to stderr
    var stderr bytes.Buffer
    cmd.Stderr = &stderr

    err := cmd.Run()

    // Parse pixel difference from stderr
    diffPixels, _ := strconv.ParseInt(strings.TrimSpace(stderr.String()), 10, 64)

    // Calculate total pixels
    baselineInfo, _ := c.getImageInfo(baseline)
    totalPixels := baselineInfo.Width * baselineInfo.Height

    diffPercent := float64(diffPixels) / float64(totalPixels)

    return &CompareResult{
        DiffPixels:  diffPixels,
        TotalPixels: totalPixels,
        DiffPercent: diffPercent,
        DiffPath:    diffOut,
    }, nil
}

// CompareResult contains comparison metrics
type CompareResult struct {
    DiffPixels  int64   `json:"diffPixels"`
    TotalPixels int64   `json:"totalPixels"`
    DiffPercent float64 `json:"diffPercent"`
    DiffPath    string  `json:"diffPath"`
}
```

### Comparison Options

| Option | Default | Description |
|--------|---------|-------------|
| `threshold` | 0.001 (0.1%) | Max allowed diff percentage |
| `fuzz` | 2% | Anti-aliasing tolerance |
| `metric` | AE | Absolute Error (pixel count) |
| `highlight` | red | Diff highlight color |

### Alternative: Pure Go Comparison

For environments without ImageMagick:

```go
// visual/compare_go.go

import "github.com/olegfedoseev/image-diff"

func (c *GoImageComparator) CompareImages(baseline, actual string) (*CompareResult, error) {
    img1, _ := loadImage(baseline)
    img2, _ := loadImage(actual)

    diff, percent, _ := imagediff.Diff(img1, img2, &imagediff.Options{
        Threshold: 0.1,
    })

    return &CompareResult{
        DiffPercent: percent,
        DiffImage:   diff,
    }, nil
}
```

## CLI Commands

### dss visual-test

```go
// cmd/dss/cmd/visual_test.go

var visualTestCmd = &cobra.Command{
    Use:   "visual-test",
    Short: "Run visual regression tests",
    Long:  `Run visual regression tests comparing current rendering against baselines.`,
    RunE:  runVisualTest,
}

func init() {
    visualTestCmd.Flags().StringP("baseline", "b", "latest", "Baseline version to compare against")
    visualTestCmd.Flags().StringP("tests", "t", "./visual-tests", "Test definition directory")
    visualTestCmd.Flags().StringP("output", "o", "./visual-results", "Output directory for results")
    visualTestCmd.Flags().StringSlice("test", nil, "Run specific test(s) by ID")
    visualTestCmd.Flags().StringSlice("viewport", nil, "Run specific viewport(s)")
    visualTestCmd.Flags().Float64("threshold", 0.001, "Default diff threshold")
    visualTestCmd.Flags().Int("parallel", 4, "Parallel test execution")
    visualTestCmd.Flags().Bool("update-baseline", false, "Update baseline with current results")
    visualTestCmd.Flags().Bool("json", false, "Output JSON report")
    visualTestCmd.Flags().Bool("interactive", false, "Interactive mode (open diff viewer)")

    rootCmd.AddCommand(visualTestCmd)
}
```

### dss visual-baseline

```go
// cmd/dss/cmd/visual_baseline.go

var visualBaselineCmd = &cobra.Command{
    Use:   "visual-baseline",
    Short: "Manage visual test baselines",
}

var visualBaselineGenerateCmd = &cobra.Command{
    Use:   "generate",
    Short: "Generate baseline screenshots",
    RunE:  runVisualBaselineGenerate,
}

var visualBaselineUpdateCmd = &cobra.Command{
    Use:   "update",
    Short: "Update baseline for specific tests",
    RunE:  runVisualBaselineUpdate,
}

var visualBaselineListCmd = &cobra.Command{
    Use:   "list",
    Short: "List available baselines",
    RunE:  runVisualBaselineList,
}

func init() {
    visualBaselineGenerateCmd.Flags().StringP("version", "v", "", "Version tag for baseline (required)")
    visualBaselineGenerateCmd.Flags().StringP("tests", "t", "./visual-tests", "Test definition directory")
    visualBaselineGenerateCmd.Flags().StringP("output", "o", "./baselines", "Baseline output directory")
    visualBaselineGenerateCmd.Flags().Int("parallel", 4, "Parallel execution")
    visualBaselineGenerateCmd.MarkFlagRequired("version")

    visualBaselineUpdateCmd.Flags().StringP("version", "v", "", "Baseline version to update")
    visualBaselineUpdateCmd.Flags().StringSlice("test", nil, "Test IDs to update")
    visualBaselineUpdateCmd.Flags().Bool("all", false, "Update all tests")

    visualBaselineCmd.AddCommand(visualBaselineGenerateCmd)
    visualBaselineCmd.AddCommand(visualBaselineUpdateCmd)
    visualBaselineCmd.AddCommand(visualBaselineListCmd)

    rootCmd.AddCommand(visualBaselineCmd)
}
```

## MCP Tools

### tools_visual.go

```go
// skills/designsystem/tools_visual.go

func (s *Skill) visualTestTool() skill.Tool {
    return skill.NewTool(
        "visual_test",
        "Run visual regression tests against baseline",
        map[string]skill.Parameter{
            "baseline_version": {
                Type:        "string",
                Description: "Baseline version to compare against (default: latest)",
                Required:    false,
            },
            "test_ids": {
                Type:        "array",
                Description: "Specific test IDs to run (default: all)",
                Required:    false,
            },
            "threshold": {
                Type:        "number",
                Description: "Diff threshold (default: 0.001)",
                Required:    false,
            },
        },
        func(ctx context.Context, params map[string]any) (any, error) {
            opts := &visual.TestOptions{
                BaselineVersion: "latest",
                Threshold:       0.001,
            }

            if v, ok := params["baseline_version"].(string); ok {
                opts.BaselineVersion = v
            }
            if ids, ok := params["test_ids"].([]any); ok {
                opts.TestIDs = toStringSlice(ids)
            }
            if t, ok := params["threshold"].(float64); ok {
                opts.Threshold = t
            }

            report, err := s.visualService.RunTests(ctx, opts)
            if err != nil {
                return nil, err
            }

            return report, nil
        },
    )
}

func (s *Skill) visualBaselineGenerateTool() skill.Tool {
    return skill.NewTool(
        "visual_baseline_generate",
        "Generate baseline screenshots for a version",
        map[string]skill.Parameter{
            "version": {
                Type:        "string",
                Description: "Version tag for the baseline",
                Required:    true,
            },
        },
        func(ctx context.Context, params map[string]any) (any, error) {
            version := params["version"].(string)

            result, err := s.visualService.GenerateBaseline(ctx, version)
            if err != nil {
                return nil, err
            }

            return result, nil
        },
    )
}
```

## Test Definition Schema

### JSON Schema

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "https://dss.dev/schema/visual-test-suite.json",
  "title": "Visual Test Suite",
  "type": "object",
  "required": ["version", "name", "tests"],
  "properties": {
    "version": {
      "type": "string",
      "pattern": "^\\d+\\.\\d+$"
    },
    "name": {
      "type": "string"
    },
    "description": {
      "type": "string"
    },
    "baseUrl": {
      "type": "string",
      "format": "uri"
    },
    "defaults": {
      "$ref": "#/definitions/defaults"
    },
    "tests": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/test"
      }
    }
  },
  "definitions": {
    "viewport": {
      "type": "object",
      "required": ["name", "width", "height"],
      "properties": {
        "name": { "type": "string" },
        "width": { "type": "integer", "minimum": 320 },
        "height": { "type": "integer", "minimum": 200 }
      }
    },
    "stabilization": {
      "type": "object",
      "properties": {
        "waitForSelector": { "type": "string" },
        "waitForTimeout": { "type": "integer" },
        "waitMs": { "type": "integer" },
        "disableAnimations": { "type": "boolean" }
      }
    },
    "defaults": {
      "type": "object",
      "properties": {
        "viewports": {
          "type": "array",
          "items": { "$ref": "#/definitions/viewport" }
        },
        "threshold": { "type": "number" },
        "stabilization": { "$ref": "#/definitions/stabilization" }
      }
    },
    "test": {
      "type": "object",
      "required": ["id", "component", "url"],
      "properties": {
        "id": { "type": "string" },
        "name": { "type": "string" },
        "component": { "type": "string" },
        "variant": { "type": "string" },
        "url": { "type": "string" },
        "selector": { "type": "string" },
        "viewports": {
          "type": "array",
          "items": { "$ref": "#/definitions/viewport" }
        },
        "threshold": { "type": "number" },
        "stabilization": { "$ref": "#/definitions/stabilization" },
        "skip": { "type": "boolean" },
        "skipReason": { "type": "string" }
      }
    }
  }
}
```

## Error Handling

### Error Types

```go
// visual/errors.go

var (
    ErrBaselineNotFound    = errors.New("baseline not found")
    ErrTestNotFound        = errors.New("test not found")
    ErrW3PilotUnavailable  = errors.New("w3pilot not available")
    ErrImageMagickMissing  = errors.New("imagemagick not installed")
    ErrScreenshotFailed    = errors.New("screenshot capture failed")
    ErrComparisonFailed    = errors.New("image comparison failed")
    ErrThresholdExceeded   = errors.New("diff threshold exceeded")
)

// VisualTestError wraps errors with context
type VisualTestError struct {
    TestID   string
    Viewport string
    Op       string
    Err      error
}

func (e *VisualTestError) Error() string {
    return fmt.Sprintf("%s [%s/%s]: %v", e.Op, e.TestID, e.Viewport, e.Err)
}
```

### Recovery Strategies

| Error | Recovery |
|-------|----------|
| W3Pilot crash | Restart subprocess, retry test |
| Navigation timeout | Retry with increased timeout |
| Screenshot timeout | Retry once, mark as error |
| Missing baseline | Generate baseline or skip test |
| Comparison failure | Log error, include in report |

## Performance Considerations

### Parallelization

```go
// visual/executor.go

func (e *Executor) RunParallel(ctx context.Context, tests []VisualTest, workers int) (*VisualTestReport, error) {
    // Create worker pool
    jobs := make(chan VisualTest, len(tests))
    results := make(chan VisualTestResult, len(tests))

    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            // Each worker has its own w3pilot client
            client, _ := NewW3PilotClient(ctx)
            defer client.Close()

            for test := range jobs {
                result := e.runSingleTest(ctx, client, test)
                results <- result
            }
        }()
    }

    // Send jobs
    for _, test := range tests {
        jobs <- test
    }
    close(jobs)

    // Collect results
    wg.Wait()
    close(results)

    return e.aggregateResults(results), nil
}
```

### Resource Limits

| Resource | Limit | Rationale |
|----------|-------|-----------|
| Workers | 4 (default) | Browser memory usage |
| Screenshot timeout | 10s | Prevent hanging |
| Navigation timeout | 30s | Slow page loads |
| Total suite timeout | 10m | CI timeout safety |
| Image size | 4096x4096 | Memory/storage |

## Dependencies

### Required

| Dependency | Version | Purpose |
|------------|---------|---------|
| w3pilot | >=0.7.0 | Screenshot capture |
| ImageMagick | >=7.0 | Image comparison |

### Go Dependencies

```go
require (
    github.com/plexusone/w3pilot v0.7.0  // Not a Go import, used as subprocess
    github.com/modelcontextprotocol/go-sdk v1.6.0
    gopkg.in/yaml.v3 v3.0.1
)
```

### Optional

| Dependency | Purpose |
|------------|---------|
| Chrome/Chromium | Default browser |
| Firefox | Cross-browser testing |
| go-imagediff | Pure Go comparison (no ImageMagick) |

## Security Considerations

### Subprocess Security

- W3pilot runs as subprocess with limited permissions
- No shell injection in parameters (use exec.Command, not shell)
- Sandbox browser with `--no-sandbox` only in containers

### File Access

- Baselines stored in project directory (git-managed)
- Results stored in temp/output directory
- No write access outside designated directories

### Network

- Test URLs must be localhost or configured allowlist
- No external network access during tests (optional)

## Related Documents

- [PRD.md](PRD.md) - Product requirements
- [PLAN.md](PLAN.md) - Implementation plan
- [ROADMAP.md](ROADMAP.md) - Milestone roadmap
