package visual

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Comparator compares images and generates diffs.
type Comparator interface {
	// Compare compares two images and generates a diff image.
	// Returns the comparison result including diff percentage.
	Compare(ctx context.Context, baseline, actual, diffOut string) (*CompareResult, error)
}

// CompareResult contains comparison metrics.
type CompareResult struct {
	DiffPixels  int64   `json:"diffPixels"`
	TotalPixels int64   `json:"totalPixels"`
	DiffPercent float64 `json:"diffPercent"`
	DiffPath    string  `json:"diffPath,omitempty"`
}

// ImageMagickComparator uses ImageMagick for comparison.
type ImageMagickComparator struct {
	comparePath  string
	identifyPath string
	fuzz         string // Anti-aliasing tolerance (e.g., "2%")
}

// NewImageMagickComparator creates a comparator using ImageMagick.
func NewImageMagickComparator() (*ImageMagickComparator, error) {
	comparePath, err := exec.LookPath("compare")
	if err != nil {
		return nil, fmt.Errorf("%w: compare command not found", ErrImageMagickMissing)
	}

	identifyPath, err := exec.LookPath("identify")
	if err != nil {
		return nil, fmt.Errorf("%w: identify command not found", ErrImageMagickMissing)
	}

	return &ImageMagickComparator{
		comparePath:  comparePath,
		identifyPath: identifyPath,
		fuzz:         "2%",
	}, nil
}

// SetFuzz sets the anti-aliasing tolerance.
func (c *ImageMagickComparator) SetFuzz(fuzz string) {
	c.fuzz = fuzz
}

// Compare compares two images using ImageMagick.
func (c *ImageMagickComparator) Compare(ctx context.Context, baseline, actual, diffOut string) (*CompareResult, error) {
	// Get total pixels from baseline image
	totalPixels, err := c.getPixelCount(ctx, baseline)
	if err != nil {
		return nil, fmt.Errorf("failed to get baseline dimensions: %w", err)
	}

	// Run ImageMagick compare
	// -metric AE: Absolute Error (pixel count that differ)
	// -fuzz: Tolerance for anti-aliasing
	//nolint:gosec // G204: comparePath is from exec.LookPath, file paths are user-provided as expected
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

	// compare returns exit code 1 if images differ (not an error)
	// Exit code 2+ indicates actual error
	err = cmd.Run()
	exitCode := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	}

	if exitCode > 1 {
		return nil, fmt.Errorf("compare failed: %s", stderr.String())
	}

	// Parse pixel count from stderr
	output := strings.TrimSpace(stderr.String())
	diffPixels, err := strconv.ParseInt(output, 10, 64)
	if err != nil {
		// Sometimes ImageMagick outputs with scientific notation
		var diffFloat float64
		if _, parseErr := fmt.Sscanf(output, "%e", &diffFloat); parseErr == nil {
			diffPixels = int64(diffFloat)
		} else {
			return nil, fmt.Errorf("failed to parse diff output %q: %w", output, err)
		}
	}

	diffPercent := float64(diffPixels) / float64(totalPixels)

	return &CompareResult{
		DiffPixels:  diffPixels,
		TotalPixels: totalPixels,
		DiffPercent: diffPercent,
		DiffPath:    diffOut,
	}, nil
}

// getPixelCount returns the total pixel count of an image.
func (c *ImageMagickComparator) getPixelCount(ctx context.Context, imagePath string) (int64, error) {
	//nolint:gosec // G204: identifyPath is from exec.LookPath, imagePath is user-provided as expected
	cmd := exec.CommandContext(ctx, c.identifyPath, "-format", "%w %h", imagePath)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	parts := strings.Fields(string(output))
	if len(parts) != 2 {
		return 0, fmt.Errorf("unexpected identify output: %s", output)
	}

	width, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, err
	}

	height, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, err
	}

	return width * height, nil
}

// GoImageComparator is a pure Go image comparator (no ImageMagick required).
type GoImageComparator struct {
	threshold float64 // Per-channel threshold (0-255)
}

// NewGoImageComparator creates a pure Go comparator.
func NewGoImageComparator() *GoImageComparator {
	return &GoImageComparator{
		threshold: 5, // Allow 5 units of difference per channel
	}
}

// SetThreshold sets the per-channel difference threshold.
func (c *GoImageComparator) SetThreshold(threshold float64) {
	c.threshold = threshold
}

// Compare compares two PNG images using pure Go.
func (c *GoImageComparator) Compare(ctx context.Context, baseline, actual, diffOut string) (*CompareResult, error) {
	// Load images
	baselineImg, err := loadPNG(baseline)
	if err != nil {
		return nil, fmt.Errorf("failed to load baseline: %w", err)
	}

	actualImg, err := loadPNG(actual)
	if err != nil {
		return nil, fmt.Errorf("failed to load actual: %w", err)
	}

	// Check dimensions match
	baseBounds := baselineImg.Bounds()
	actualBounds := actualImg.Bounds()

	if baseBounds.Dx() != actualBounds.Dx() || baseBounds.Dy() != actualBounds.Dy() {
		return nil, fmt.Errorf("image dimensions differ: %dx%d vs %dx%d",
			baseBounds.Dx(), baseBounds.Dy(), actualBounds.Dx(), actualBounds.Dy())
	}

	width := baseBounds.Dx()
	height := baseBounds.Dy()
	totalPixels := int64(width * height)

	// Create diff image
	diffImg := image.NewRGBA(baseBounds)
	var diffPixels int64

	// Compare pixel by pixel
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			br, bg, bb, ba := baselineImg.At(x+baseBounds.Min.X, y+baseBounds.Min.Y).RGBA()
			ar, ag, ab, aa := actualImg.At(x+actualBounds.Min.X, y+actualBounds.Min.Y).RGBA()

			// Convert from 16-bit to 8-bit
			br, bg, bb, ba = br>>8, bg>>8, bb>>8, ba>>8
			ar, ag, ab, aa = ar>>8, ag>>8, ab>>8, aa>>8

			// Check if pixels differ beyond threshold
			if absDiff(br, ar) > uint32(c.threshold) ||
				absDiff(bg, ag) > uint32(c.threshold) ||
				absDiff(bb, ab) > uint32(c.threshold) ||
				absDiff(ba, aa) > uint32(c.threshold) {
				diffPixels++
				// Mark diff in red
				diffImg.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
			} else {
				// Copy baseline pixel with reduced opacity
				//nolint:gosec // G115: Values are safe after >>8 shift (0-255 range)
				diffImg.Set(x, y, color.RGBA{R: uint8(br), G: uint8(bg), B: uint8(bb), A: 128})
			}
		}
	}

	// Save diff image
	if diffOut != "" {
		if err := savePNG(diffOut, diffImg); err != nil {
			return nil, fmt.Errorf("failed to save diff: %w", err)
		}
	}

	diffPercent := float64(diffPixels) / float64(totalPixels)

	return &CompareResult{
		DiffPixels:  diffPixels,
		TotalPixels: totalPixels,
		DiffPercent: diffPercent,
		DiffPath:    diffOut,
	}, nil
}

func loadPNG(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return png.Decode(f)
}

func savePNG(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func absDiff(a, b uint32) uint32 {
	if a > b {
		return a - b
	}
	return b - a
}

// NewComparator creates the best available comparator.
// It prefers ImageMagick if available, falling back to pure Go.
func NewComparator() (Comparator, error) {
	// Try ImageMagick first
	im, err := NewImageMagickComparator()
	if err == nil {
		return im, nil
	}

	// Fall back to Go implementation
	return NewGoImageComparator(), nil
}
