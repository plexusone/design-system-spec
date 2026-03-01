package dss

// Component represents a reusable UI element in the design system.
type Component struct {
	// ID is a unique identifier (e.g., "button", "input", "modal").
	ID string `json:"id"`

	// Name is a display name for the component.
	Name string `json:"name"`

	// Description provides an overview of the component's purpose.
	Description string `json:"description,omitempty"`

	// Category groups related components (e.g., "inputs", "navigation", "feedback").
	Category string `json:"category,omitempty"`

	// Variants define different visual or behavioral versions.
	Variants []Variant `json:"variants,omitempty"`

	// States define interactive states (hover, focus, disabled, etc.).
	States []State `json:"states,omitempty"`

	// Props define configurable properties/attributes.
	Props []Prop `json:"props,omitempty"`

	// Slots define content insertion points for composition.
	Slots []Slot `json:"slots,omitempty"`

	// TokensUsed references foundation token IDs used by this component.
	TokensUsed []string `json:"tokensUsed,omitempty"`

	// Constraints define usage limitations and requirements.
	Constraints *Constraints `json:"constraints,omitempty"`

	// Accessibility defines a11y requirements for this component.
	Accessibility *ComponentA11y `json:"accessibility,omitempty"`

	// LLM provides optimization fields for AI code generation.
	LLM *LLMContext `json:"llm,omitempty"`

	// DocumentationURL links to detailed documentation.
	DocumentationURL string `json:"documentationUrl,omitempty" jsonschema:"format=uri"`

	// FigmaURL links to the Figma component.
	FigmaURL string `json:"figmaUrl,omitempty" jsonschema:"format=uri"`
}

// Variant represents a visual or behavioral variant of a component.
type Variant struct {
	// ID is a unique identifier for the variant.
	ID string `json:"id"`

	// Name is a display name (e.g., "Primary", "Secondary", "Outline").
	Name string `json:"name"`

	// Description explains when to use this variant.
	Description string `json:"description,omitempty"`

	// IsDefault indicates this is the default variant.
	IsDefault bool `json:"isDefault,omitempty"`

	// TokenOverrides lists tokens that differ from the base component.
	TokenOverrides []TokenOverride `json:"tokenOverrides,omitempty"`
}

// TokenOverride specifies a token value that differs for a variant.
type TokenOverride struct {
	// TokenID references the token being overridden.
	TokenID string `json:"tokenId"`

	// Value is the override value.
	Value string `json:"value"`
}

// State represents an interactive state of a component.
type State struct {
	// ID is a unique identifier (hover, active, focus, disabled, error, loading).
	ID string `json:"id"`

	// Name is a display name for the state.
	Name string `json:"name,omitempty"`

	// Description explains when this state applies.
	Description string `json:"description,omitempty"`

	// TokenOverrides lists tokens that change in this state.
	TokenOverrides []TokenOverride `json:"tokenOverrides,omitempty"`
}

// Prop represents a configurable property of a component.
type Prop struct {
	// Name is the prop name as used in code.
	Name string `json:"name"`

	// Type is the data type (string, number, boolean, enum, etc.).
	Type string `json:"type"`

	// Description explains what the prop controls.
	Description string `json:"description,omitempty"`

	// Required indicates whether the prop must be provided.
	Required bool `json:"required,omitempty"`

	// Default is the default value if not specified.
	Default interface{} `json:"default,omitempty"`

	// EnumValues lists allowed values for enum types.
	EnumValues []string `json:"enumValues,omitempty"`
}

// Slot represents a content insertion point for component composition.
type Slot struct {
	// Name is the slot name.
	Name string `json:"name"`

	// Description explains what content belongs in this slot.
	Description string `json:"description,omitempty"`

	// Required indicates whether content must be provided.
	Required bool `json:"required,omitempty"`

	// AllowedComponents lists component IDs that can be inserted.
	AllowedComponents []string `json:"allowedComponents,omitempty"`
}

// Constraints define usage limitations for a component.
type Constraints struct {
	// MinWidth is the minimum recommended width.
	MinWidth string `json:"minWidth,omitempty"`

	// MaxWidth is the maximum recommended width.
	MaxWidth string `json:"maxWidth,omitempty"`

	// MaxChildren limits the number of child elements.
	MaxChildren int `json:"maxChildren,omitempty"`

	// RequiredParent specifies a required parent component.
	RequiredParent string `json:"requiredParent,omitempty"`

	// ForbiddenChildren lists components that cannot be nested inside.
	ForbiddenChildren []string `json:"forbiddenChildren,omitempty"`

	// ResponsiveBehavior describes how the component adapts.
	ResponsiveBehavior string `json:"responsiveBehavior,omitempty"`
}

// ComponentA11y defines accessibility requirements for a component.
type ComponentA11y struct {
	// Role is the ARIA role for the component.
	Role string `json:"role,omitempty"`

	// RequiredAttributes lists required ARIA attributes.
	RequiredAttributes []string `json:"requiredAttributes,omitempty"`

	// KeyboardSupport describes keyboard interaction requirements.
	KeyboardSupport []KeyboardInteraction `json:"keyboardSupport,omitempty"`

	// ScreenReaderNotes provides guidance for screen reader users.
	ScreenReaderNotes string `json:"screenReaderNotes,omitempty"`

	// FocusManagement describes focus behavior.
	FocusManagement string `json:"focusManagement,omitempty"`
}

// KeyboardInteraction describes a keyboard shortcut or interaction.
type KeyboardInteraction struct {
	// Key is the key or key combination.
	Key string `json:"key"`

	// Action describes what happens when pressed.
	Action string `json:"action"`
}
