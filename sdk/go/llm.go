package dss

// LLMContext provides structured metadata for LLM reasoning and code generation.
// These fields help LLMs understand when and how to use design system elements.
type LLMContext struct {
	// Intent describes the primary purpose of this element.
	// Example: "Primary action button for form submissions and important CTAs"
	Intent string `json:"intent,omitempty"`

	// AllowedContexts lists situations where this element should be used.
	// Example: ["form-submit", "modal-confirm", "primary-cta"]
	AllowedContexts []string `json:"allowedContexts,omitempty"`

	// ForbiddenContexts lists situations where this element should NOT be used.
	// Example: ["destructive-action", "navigation", "inline-text"]
	ForbiddenContexts []string `json:"forbiddenContexts,omitempty"`

	// Dependencies lists other elements this one requires.
	// Example: ["form", "icon-library"]
	Dependencies []string `json:"dependencies,omitempty"`

	// ExampleUsage provides code snippets or descriptions of correct usage.
	// Example: ["<Button variant='primary'>Submit</Button>"]
	ExampleUsage []string `json:"exampleUsage,omitempty"`

	// AntiPatterns describes common mistakes to avoid.
	// Example: ["Don't use for destructive actions", "Avoid multiple primary buttons"]
	AntiPatterns []string `json:"antiPatterns,omitempty"`

	// SemanticMeaning describes the user-facing meaning of this element.
	// Example: "Signals the main action users should take"
	SemanticMeaning string `json:"semanticMeaning,omitempty"`

	// RelatedElements lists related elements for context.
	// Example: ["button-secondary", "button-outline", "link"]
	RelatedElements []string `json:"relatedElements,omitempty"`

	// PriorityScore indicates relative importance (1-100).
	// Higher scores mean this element should be preferred when multiple options exist.
	PriorityScore int `json:"priorityScore,omitempty" jsonschema:"minimum=1,maximum=100"`

	// Tags are freeform labels for categorization and search.
	Tags []string `json:"tags,omitempty"`
}

// WithLLMContext is an interface for types that have LLM optimization metadata.
type WithLLMContext interface {
	// GetLLMContext returns the LLM context for this element.
	GetLLMContext() *LLMContext
}
