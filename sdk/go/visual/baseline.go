package visual

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// BaselineManifest describes a baseline snapshot.
type BaselineManifest struct {
	Version   string            `json:"version"`
	CreatedAt time.Time         `json:"createdAt"`
	CreatedBy string            `json:"createdBy,omitempty"`
	TestSuite string            `json:"testSuite"`
	TestCount int               `json:"testCount"`
	Checksums map[string]string `json:"checksums"` // "testID/viewport" -> SHA256
}

// BaselineManager handles baseline storage and retrieval.
type BaselineManager struct {
	basePath string
}

// NewBaselineManager creates a manager for the given directory.
func NewBaselineManager(basePath string) *BaselineManager {
	return &BaselineManager{basePath: basePath}
}

// GetVersionPath returns the directory path for a specific version.
func (m *BaselineManager) GetVersionPath(version string) string {
	return filepath.Join(m.basePath, version)
}

// VersionExists checks if a baseline version exists.
func (m *BaselineManager) VersionExists(version string) bool {
	manifestPath := filepath.Join(m.GetVersionPath(version), "manifest.json")
	_, err := os.Stat(manifestPath)
	return err == nil
}

// GetBaseline returns the path to a baseline image.
func (m *BaselineManager) GetBaseline(version, testID, viewport string) (string, error) {
	filename := imageFilename(testID, viewport)
	path := filepath.Join(m.GetVersionPath(version), filename)

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("%w: %s/%s in version %s", ErrBaselineNotFound, testID, viewport, version)
		}
		return "", err
	}

	return path, nil
}

// SaveBaseline saves a baseline image.
func (m *BaselineManager) SaveBaseline(version, testID, viewport string, data []byte) error {
	versionPath := m.GetVersionPath(version)
	if err := os.MkdirAll(versionPath, 0755); err != nil {
		return fmt.Errorf("failed to create baseline directory: %w", err)
	}

	filename := imageFilename(testID, viewport)
	path := filepath.Join(versionPath, filename)

	//nolint:gosec // G306: Baseline images are intentionally world-readable
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write baseline: %w", err)
	}

	return nil
}

// DeleteBaseline deletes a baseline image.
func (m *BaselineManager) DeleteBaseline(version, testID, viewport string) error {
	filename := imageFilename(testID, viewport)
	path := filepath.Join(m.GetVersionPath(version), filename)
	return os.Remove(path)
}

// GenerateManifest creates a manifest for the version.
func (m *BaselineManager) GenerateManifest(version string, suite *VisualTestSuite) (*BaselineManifest, error) {
	versionPath := m.GetVersionPath(version)

	manifest := &BaselineManifest{
		Version:   version,
		CreatedAt: time.Now(),
		TestSuite: suite.Name,
		TestCount: 0,
		Checksums: make(map[string]string),
	}

	// Calculate checksums for all baseline images
	for _, test := range suite.Tests {
		if test.Skip {
			continue
		}

		for _, vp := range test.Viewports {
			filename := imageFilename(test.ID, vp.Name)
			path := filepath.Join(versionPath, filename)

			checksum, err := fileChecksum(path)
			if err != nil {
				continue // Skip missing files
			}

			key := fmt.Sprintf("%s/%s", test.ID, vp.Name)
			manifest.Checksums[key] = checksum
			manifest.TestCount++
		}
	}

	// Save manifest
	manifestPath := filepath.Join(versionPath, "manifest.json")
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal manifest: %w", err)
	}

	//nolint:gosec // G306: Manifest files are intentionally world-readable
	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write manifest: %w", err)
	}

	return manifest, nil
}

// LoadManifest loads the manifest for a version.
func (m *BaselineManager) LoadManifest(version string) (*BaselineManifest, error) {
	path := filepath.Join(m.GetVersionPath(version), "manifest.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: no manifest for version %s", ErrBaselineNotFound, version)
		}
		return nil, err
	}

	var manifest BaselineManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	return &manifest, nil
}

// ListVersions returns all available baseline versions.
func (m *BaselineManager) ListVersions() ([]string, error) {
	entries, err := os.ReadDir(m.basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var versions []string
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "latest" {
			// Check if it has a manifest
			manifestPath := filepath.Join(m.basePath, entry.Name(), "manifest.json")
			if _, err := os.Stat(manifestPath); err == nil {
				versions = append(versions, entry.Name())
			}
		}
	}

	// Sort versions (newest first)
	sort.Sort(sort.Reverse(sort.StringSlice(versions)))

	return versions, nil
}

// GetLatestVersion returns the most recent baseline version.
func (m *BaselineManager) GetLatestVersion() (string, error) {
	// Check for "latest" symlink
	latestPath := filepath.Join(m.basePath, "latest")
	if target, err := os.Readlink(latestPath); err == nil {
		return target, nil
	}

	// Otherwise, return the first version from the sorted list
	versions, err := m.ListVersions()
	if err != nil {
		return "", err
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("%w: no baseline versions found", ErrBaselineNotFound)
	}

	return versions[0], nil
}

// UpdateLatest updates the "latest" symlink to point to a version.
func (m *BaselineManager) UpdateLatest(version string) error {
	latestPath := filepath.Join(m.basePath, "latest")

	// Remove existing symlink if present
	os.Remove(latestPath)

	// Create new symlink (relative path)
	return os.Symlink(version, latestPath)
}

// VerifyBaseline verifies a baseline image's checksum.
func (m *BaselineManager) VerifyBaseline(version, testID, viewport string) error {
	manifest, err := m.LoadManifest(version)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s", testID, viewport)
	expectedChecksum, ok := manifest.Checksums[key]
	if !ok {
		return fmt.Errorf("%w: %s not in manifest", ErrBaselineNotFound, key)
	}

	path, err := m.GetBaseline(version, testID, viewport)
	if err != nil {
		return err
	}

	actualChecksum, err := fileChecksum(path)
	if err != nil {
		return err
	}

	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch for %s: expected %s, got %s", key, expectedChecksum, actualChecksum)
	}

	return nil
}

// PruneVersion removes a baseline version and all its images.
func (m *BaselineManager) PruneVersion(version string) error {
	versionPath := m.GetVersionPath(version)
	return os.RemoveAll(versionPath)
}

// imageFilename generates the filename for a baseline image.
func imageFilename(testID, viewport string) string {
	return fmt.Sprintf("%s-%s.png", testID, viewport)
}

// fileChecksum calculates the SHA256 checksum of a file.
func fileChecksum(path string) (string, error) {
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
