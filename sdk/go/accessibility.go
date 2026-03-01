package dss

// Accessibility defines system-wide accessibility requirements and guidelines.
type Accessibility struct {
	// WCAGLevel is the target WCAG conformance level (A, AA, AAA).
	WCAGLevel string `json:"wcagLevel" jsonschema:"enum=A,enum=AA,enum=AAA"`

	// WCAGVersion is the WCAG version (e.g., "2.1", "2.2").
	WCAGVersion string `json:"wcagVersion,omitempty"`

	// ColorContrast defines color contrast requirements.
	ColorContrast *ColorContrastRequirements `json:"colorContrast,omitempty"`

	// Keyboard defines keyboard accessibility requirements.
	Keyboard *KeyboardRequirements `json:"keyboard,omitempty"`

	// ScreenReader defines screen reader support requirements.
	ScreenReader *ScreenReaderRequirements `json:"screenReader,omitempty"`

	// Motion defines motion and animation requirements.
	Motion *MotionRequirements `json:"motion,omitempty"`

	// TouchTarget defines touch target size requirements.
	TouchTarget *TouchTargetRequirements `json:"touchTarget,omitempty"`

	// Guidelines provides additional accessibility guidelines.
	Guidelines []AccessibilityGuideline `json:"guidelines,omitempty"`

	// TestingRequirements describes required accessibility testing.
	TestingRequirements []TestingRequirement `json:"testingRequirements,omitempty"`
}

// ColorContrastRequirements defines contrast ratio requirements.
type ColorContrastRequirements struct {
	// NormalTextRatio is the minimum ratio for normal text.
	NormalTextRatio float64 `json:"normalTextRatio"`

	// LargeTextRatio is the minimum ratio for large text.
	LargeTextRatio float64 `json:"largeTextRatio"`

	// NonTextRatio is the minimum ratio for non-text elements.
	NonTextRatio float64 `json:"nonTextRatio,omitempty"`

	// LargeTextThreshold defines what counts as large text (in pixels).
	LargeTextThreshold int `json:"largeTextThreshold,omitempty"`
}

// KeyboardRequirements defines keyboard accessibility standards.
type KeyboardRequirements struct {
	// FocusVisible requires visible focus indicators.
	FocusVisible bool `json:"focusVisible"`

	// FocusOrder requires logical focus order.
	FocusOrder bool `json:"focusOrder"`

	// NoKeyboardTrap prevents keyboard traps.
	NoKeyboardTrap bool `json:"noKeyboardTrap"`

	// SkipLinks requires skip navigation links.
	SkipLinks bool `json:"skipLinks,omitempty"`

	// FocusIndicatorMinSize minimum size for focus indicators (in CSS pixels).
	FocusIndicatorMinSize int `json:"focusIndicatorMinSize,omitempty"`

	// ShortcutConflicts lists keyboard shortcuts to avoid.
	ShortcutConflicts []string `json:"shortcutConflicts,omitempty"`
}

// ScreenReaderRequirements defines screen reader support standards.
type ScreenReaderRequirements struct {
	// SemanticHTML requires semantic HTML elements.
	SemanticHTML bool `json:"semanticHtml"`

	// ARIALabels requires ARIA labels for interactive elements.
	ARIALabels bool `json:"ariaLabels"`

	// LiveRegions requires live regions for dynamic content.
	LiveRegions bool `json:"liveRegions,omitempty"`

	// HeadingStructure requires proper heading hierarchy.
	HeadingStructure bool `json:"headingStructure,omitempty"`

	// SupportedReaders lists screen readers to test with.
	SupportedReaders []string `json:"supportedReaders,omitempty"`
}

// MotionRequirements defines motion and animation standards.
type MotionRequirements struct {
	// ReducedMotion requires respecting prefers-reduced-motion.
	ReducedMotion bool `json:"reducedMotion"`

	// NoAutoPlay prevents auto-playing animations longer than 5 seconds.
	NoAutoPlay bool `json:"noAutoPlay,omitempty"`

	// PauseControl requires pause/stop controls for animations.
	PauseControl bool `json:"pauseControl,omitempty"`

	// FlashingThreshold maximum flashes per second (must be <=3).
	FlashingThreshold int `json:"flashingThreshold,omitempty" jsonschema:"maximum=3"`
}

// TouchTargetRequirements defines touch target size standards.
type TouchTargetRequirements struct {
	// MinSize is the minimum touch target size (e.g., "44px").
	MinSize string `json:"minSize"`

	// MinSpacing is the minimum spacing between targets.
	MinSpacing string `json:"minSpacing,omitempty"`
}

// AccessibilityGuideline provides additional accessibility guidance.
type AccessibilityGuideline struct {
	// ID is a unique identifier.
	ID string `json:"id"`

	// Title is a short title for the guideline.
	Title string `json:"title"`

	// Description explains the guideline.
	Description string `json:"description"`

	// WCAGCriteria lists related WCAG success criteria.
	WCAGCriteria []string `json:"wcagCriteria,omitempty"`

	// Examples demonstrate the guideline.
	Examples []string `json:"examples,omitempty"`
}

// TestingRequirement describes an accessibility testing requirement.
type TestingRequirement struct {
	// Type is the testing type (automated, manual, assistive-tech).
	Type string `json:"type"`

	// Tools lists required testing tools.
	Tools []string `json:"tools,omitempty"`

	// Frequency describes how often testing should occur.
	Frequency string `json:"frequency,omitempty"`

	// Coverage describes what should be tested.
	Coverage string `json:"coverage,omitempty"`
}
