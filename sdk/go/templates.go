package dss

// Template represents a page-level layout structure.
type Template struct {
	// ID is a unique identifier (e.g., "dashboard", "settings", "auth").
	ID string `json:"id"`

	// Name is a display name for the template.
	Name string `json:"name"`

	// Description explains the template's purpose and when to use it.
	Description string `json:"description"`

	// Category groups related templates (e.g., "admin", "marketing", "app").
	Category string `json:"category,omitempty"`

	// Regions defines the major layout regions.
	Regions []TemplateRegion `json:"regions"`

	// Grid describes the overall grid configuration.
	Grid *TemplateGrid `json:"grid,omitempty"`

	// Patterns lists pattern IDs commonly used in this template.
	Patterns []string `json:"patterns,omitempty"`

	// Responsive describes breakpoint-specific adaptations.
	Responsive []ResponsiveTemplate `json:"responsive,omitempty"`

	// LLM provides optimization fields for AI code generation.
	LLM *LLMContext `json:"llm,omitempty"`

	// DocumentationURL links to detailed documentation.
	DocumentationURL string `json:"documentationUrl,omitempty" jsonschema:"format=uri"`

	// PreviewURL links to a visual preview.
	PreviewURL string `json:"previewUrl,omitempty" jsonschema:"format=uri"`
}

// TemplateRegion defines a major section of a page template.
type TemplateRegion struct {
	// ID is a unique identifier for the region.
	ID string `json:"id"`

	// Name is a display name (e.g., "Header", "Sidebar", "Main Content").
	Name string `json:"name"`

	// Role describes the semantic role of this region.
	Role string `json:"role,omitempty"`

	// Required indicates whether this region must be present.
	Required bool `json:"required,omitempty"`

	// GridArea specifies the CSS grid area name.
	GridArea string `json:"gridArea,omitempty"`

	// MinHeight specifies minimum height.
	MinHeight string `json:"minHeight,omitempty"`

	// MaxHeight specifies maximum height.
	MaxHeight string `json:"maxHeight,omitempty"`

	// Sticky indicates whether the region should stick on scroll.
	Sticky bool `json:"sticky,omitempty"`

	// AllowedPatterns lists pattern IDs that can be placed here.
	AllowedPatterns []string `json:"allowedPatterns,omitempty"`

	// AllowedComponents lists component IDs that can be placed here.
	AllowedComponents []string `json:"allowedComponents,omitempty"`
}

// TemplateGrid defines the grid configuration for a template.
type TemplateGrid struct {
	// Areas is the CSS grid-template-areas value.
	Areas string `json:"areas,omitempty"`

	// Columns is the CSS grid-template-columns value.
	Columns string `json:"columns,omitempty"`

	// Rows is the CSS grid-template-rows value.
	Rows string `json:"rows,omitempty"`

	// Gap references a spacing token for grid gaps.
	Gap string `json:"gap,omitempty"`
}

// ResponsiveTemplate describes template changes at specific breakpoints.
type ResponsiveTemplate struct {
	// Breakpoint references a breakpoint ID.
	Breakpoint string `json:"breakpoint"`

	// GridChanges describes grid changes at this breakpoint.
	GridChanges *TemplateGrid `json:"gridChanges,omitempty"`

	// HiddenRegions lists region IDs to hide at this breakpoint.
	HiddenRegions []string `json:"hiddenRegions,omitempty"`

	// RegionChanges lists region modifications at this breakpoint.
	RegionChanges []TemplateRegionChange `json:"regionChanges,omitempty"`
}

// TemplateRegionChange describes a region change at a breakpoint.
type TemplateRegionChange struct {
	// RegionID references the region being changed.
	RegionID string `json:"regionId"`

	// Changes maps property names to new values.
	Changes map[string]string `json:"changes"`
}
