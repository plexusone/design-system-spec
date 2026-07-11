package visual

import (
	"errors"
	"fmt"
)

// Standard errors for visual testing.
var (
	// ErrBaselineNotFound indicates no baseline exists for the test.
	ErrBaselineNotFound = errors.New("baseline not found")

	// ErrTestNotFound indicates the requested test does not exist.
	ErrTestNotFound = errors.New("test not found")

	// ErrW3PilotUnavailable indicates w3pilot is not available.
	ErrW3PilotUnavailable = errors.New("w3pilot not available")

	// ErrImageMagickMissing indicates ImageMagick is not installed.
	ErrImageMagickMissing = errors.New("imagemagick not installed")

	// ErrScreenshotFailed indicates screenshot capture failed.
	ErrScreenshotFailed = errors.New("screenshot capture failed")

	// ErrComparisonFailed indicates image comparison failed.
	ErrComparisonFailed = errors.New("image comparison failed")

	// ErrThresholdExceeded indicates the diff exceeded the threshold.
	ErrThresholdExceeded = errors.New("diff threshold exceeded")

	// ErrNavigationFailed indicates page navigation failed.
	ErrNavigationFailed = errors.New("page navigation failed")

	// ErrStabilizationFailed indicates stabilization conditions were not met.
	ErrStabilizationFailed = errors.New("stabilization failed")
)

// VisualTestError wraps errors with test context.
type VisualTestError struct {
	TestID   string
	Viewport string
	Op       string
	Err      error
}

// Error returns the error message.
func (e *VisualTestError) Error() string {
	if e.Viewport != "" {
		return fmt.Sprintf("%s [%s/%s]: %v", e.Op, e.TestID, e.Viewport, e.Err)
	}
	return fmt.Sprintf("%s [%s]: %v", e.Op, e.TestID, e.Err)
}

// Unwrap returns the underlying error.
func (e *VisualTestError) Unwrap() error {
	return e.Err
}

// NewTestError creates a new VisualTestError.
func NewTestError(testID, viewport, op string, err error) *VisualTestError {
	return &VisualTestError{
		TestID:   testID,
		Viewport: viewport,
		Op:       op,
		Err:      err,
	}
}

// CaptureError wraps a capture error with context.
func CaptureError(testID, viewport string, err error) error {
	return NewTestError(testID, viewport, "capture", err)
}

// CompareError wraps a comparison error with context.
func CompareError(testID, viewport string, err error) error {
	return NewTestError(testID, viewport, "compare", err)
}

// BaselineError wraps a baseline error with context.
func BaselineError(testID, viewport string, err error) error {
	return NewTestError(testID, viewport, "baseline", err)
}
