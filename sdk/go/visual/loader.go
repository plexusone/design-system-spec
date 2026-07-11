package visual

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Loader loads visual test definitions from files.
type Loader struct {
	basePath string
}

// NewLoader creates a loader for the given base path.
func NewLoader(basePath string) *Loader {
	return &Loader{basePath: basePath}
}

// Load loads the visual test suite from the base path.
// It looks for visual-tests.yaml or visual-tests.json.
func (l *Loader) Load() (*VisualTestSuite, error) {
	// Try YAML first, then JSON
	yamlPath := filepath.Join(l.basePath, "visual-tests.yaml")
	jsonPath := filepath.Join(l.basePath, "visual-tests.json")

	var data []byte
	var err error
	var isYAML bool

	if _, err = os.Stat(yamlPath); err == nil {
		data, err = os.ReadFile(yamlPath)
		isYAML = true
	} else if _, err = os.Stat(jsonPath); err == nil {
		data, err = os.ReadFile(jsonPath)
		isYAML = false
	} else {
		return nil, fmt.Errorf("no visual-tests.yaml or visual-tests.json found in %s", l.basePath)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read test file: %w", err)
	}

	var suite VisualTestSuite
	if isYAML {
		if err := yaml.Unmarshal(data, &suite); err != nil {
			return nil, fmt.Errorf("failed to parse YAML test file: %w", err)
		}
	} else {
		if err := json.Unmarshal(data, &suite); err != nil {
			return nil, fmt.Errorf("failed to parse JSON test file: %w", err)
		}
	}

	// Validate and apply defaults
	if err := l.validate(&suite); err != nil {
		return nil, err
	}
	l.applyDefaults(&suite)

	return &suite, nil
}

// LoadFromFile loads a visual test suite from a specific file path.
func (l *Loader) LoadFromFile(path string) (*VisualTestSuite, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read test file: %w", err)
	}

	var suite VisualTestSuite
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &suite); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(data, &suite); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}

	if err := l.validate(&suite); err != nil {
		return nil, err
	}
	l.applyDefaults(&suite)

	return &suite, nil
}

// validate checks that the suite has required fields.
func (l *Loader) validate(suite *VisualTestSuite) error {
	if suite.Version == "" {
		return fmt.Errorf("suite version is required")
	}
	if suite.Name == "" {
		return fmt.Errorf("suite name is required")
	}
	if len(suite.Tests) == 0 {
		return fmt.Errorf("at least one test is required")
	}

	for i, test := range suite.Tests {
		if test.ID == "" {
			return fmt.Errorf("test[%d]: id is required", i)
		}
		if test.Component == "" {
			return fmt.Errorf("test[%d] (%s): component is required", i, test.ID)
		}
		if test.URL == "" {
			return fmt.Errorf("test[%d] (%s): url is required", i, test.ID)
		}
	}

	return nil
}

// applyDefaults merges suite defaults into individual tests.
func (l *Loader) applyDefaults(suite *VisualTestSuite) {
	for i := range suite.Tests {
		test := &suite.Tests[i]

		// Apply default viewports
		if len(test.Viewports) == 0 {
			if len(suite.Defaults.Viewports) > 0 {
				test.Viewports = suite.Defaults.Viewports
			} else {
				// Use built-in defaults if nothing specified
				test.Viewports = DefaultViewports()
			}
		}

		// Apply default threshold
		if test.Threshold == 0 {
			if suite.Defaults.Threshold > 0 {
				test.Threshold = suite.Defaults.Threshold
			} else {
				test.Threshold = DefaultThreshold()
			}
		}

		// Apply default stabilization
		if test.Stabilization == nil {
			if suite.Defaults.Stabilization != nil {
				// Copy to avoid shared mutation
				stab := *suite.Defaults.Stabilization
				test.Stabilization = &stab
			} else {
				test.Stabilization = DefaultStabilization()
			}
		}

		// Expand relative URL with base URL
		if suite.BaseURL != "" && !isAbsoluteURL(test.URL) {
			test.URL = suite.BaseURL + test.URL
		}

		// Set name from ID if not provided
		if test.Name == "" {
			test.Name = test.ID
		}
	}
}

// isAbsoluteURL checks if a URL is absolute.
func isAbsoluteURL(url string) bool {
	return strings.HasPrefix(url, "http://") ||
		strings.HasPrefix(url, "https://") ||
		strings.HasPrefix(url, "//")
}

// FilterTests returns tests matching the given IDs.
func FilterTests(tests []VisualTest, ids []string) []VisualTest {
	if len(ids) == 0 {
		return tests
	}

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

// FilterByComponent returns tests for a specific component.
func FilterByComponent(tests []VisualTest, component string) []VisualTest {
	var filtered []VisualTest
	for _, t := range tests {
		if t.Component == component {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

// ExpandTests expands tests with multiple viewports into individual test cases.
// Returns a list of (test, viewport) pairs.
func ExpandTests(tests []VisualTest) []struct {
	Test     VisualTest
	Viewport Viewport
} {
	var expanded []struct {
		Test     VisualTest
		Viewport Viewport
	}

	for _, test := range tests {
		for _, vp := range test.Viewports {
			expanded = append(expanded, struct {
				Test     VisualTest
				Viewport Viewport
			}{
				Test:     test,
				Viewport: vp,
			})
		}
	}

	return expanded
}
