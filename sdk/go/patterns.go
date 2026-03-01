package dss

// Pattern represents a multi-component UX solution for common scenarios.
type Pattern struct {
	// ID is a unique identifier (e.g., "form-validation", "data-table", "modal-flow").
	ID string `json:"id"`

	// Name is a display name for the pattern.
	Name string `json:"name"`

	// Description explains the pattern's purpose and when to use it.
	Description string `json:"description"`

	// Category groups related patterns (e.g., "forms", "navigation", "data-display").
	Category string `json:"category,omitempty"`

	// Components lists the component IDs that make up this pattern.
	Components []PatternComponent `json:"components"`

	// Layout describes how components are arranged.
	Layout *PatternLayout `json:"layout,omitempty"`

	// Behavior describes the interaction flow.
	Behavior string `json:"behavior,omitempty"`

	// Variations lists different configurations of the pattern.
	Variations []PatternVariation `json:"variations,omitempty"`

	// Accessibility defines a11y requirements for the pattern.
	Accessibility *PatternA11y `json:"accessibility,omitempty"`

	// LLM provides optimization fields for AI code generation.
	LLM *LLMContext `json:"llm,omitempty"`

	// DocumentationURL links to detailed documentation.
	DocumentationURL string `json:"documentationUrl,omitempty" jsonschema:"format=uri"`
}

// PatternComponent references a component used in a pattern.
type PatternComponent struct {
	// ComponentID references a component from the components list.
	ComponentID string `json:"componentId"`

	// Role describes the component's role in this pattern.
	Role string `json:"role,omitempty"`

	// Required indicates whether this component is essential.
	Required bool `json:"required,omitempty"`

	// Quantity specifies how many instances (e.g., "1", "1+", "0-1").
	Quantity string `json:"quantity,omitempty"`

	// Configuration provides default props for this usage.
	Configuration map[string]interface{} `json:"configuration,omitempty"`
}

// PatternLayout describes the spatial arrangement of components.
type PatternLayout struct {
	// Type is the layout strategy (flex, grid, stack, etc.).
	Type string `json:"type"`

	// Direction is the primary axis (row, column).
	Direction string `json:"direction,omitempty"`

	// Alignment describes cross-axis alignment.
	Alignment string `json:"alignment,omitempty"`

	// Spacing references a spacing token for gaps.
	Spacing string `json:"spacing,omitempty"`

	// Responsive describes breakpoint-specific changes.
	Responsive []ResponsiveLayout `json:"responsive,omitempty"`
}

// ResponsiveLayout describes layout changes at specific breakpoints.
type ResponsiveLayout struct {
	// Breakpoint references a breakpoint ID.
	Breakpoint string `json:"breakpoint"`

	// Changes describes what changes at this breakpoint.
	Changes map[string]string `json:"changes"`
}

// PatternVariation describes an alternative configuration of a pattern.
type PatternVariation struct {
	// ID is a unique identifier for the variation.
	ID string `json:"id"`

	// Name is a display name.
	Name string `json:"name"`

	// Description explains when to use this variation.
	Description string `json:"description,omitempty"`

	// ComponentChanges lists component differences from the base pattern.
	ComponentChanges []PatternComponent `json:"componentChanges,omitempty"`

	// LayoutChanges describes layout differences.
	LayoutChanges *PatternLayout `json:"layoutChanges,omitempty"`
}

// PatternA11y defines accessibility requirements for a pattern.
type PatternA11y struct {
	// LandmarkRole is the ARIA landmark for the pattern region.
	LandmarkRole string `json:"landmarkRole,omitempty"`

	// FocusOrder describes the expected tab order.
	FocusOrder []string `json:"focusOrder,omitempty"`

	// LiveRegions lists components that should announce changes.
	LiveRegions []string `json:"liveRegions,omitempty"`

	// Notes provides additional accessibility guidance.
	Notes string `json:"notes,omitempty"`
}
