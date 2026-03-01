package dss

// Principle represents a design philosophy or brand value that guides decisions.
type Principle struct {
	// ID is a unique identifier for the principle.
	ID string `json:"id"`

	// Name is a short, memorable title for the principle.
	Name string `json:"name"`

	// Description explains the principle in detail.
	Description string `json:"description"`

	// Rationale explains why this principle matters.
	Rationale string `json:"rationale,omitempty"`

	// Examples demonstrate the principle in practice.
	Examples []PrincipleExample `json:"examples,omitempty"`

	// Order determines display ordering (lower = first).
	Order int `json:"order,omitempty"`
}

// PrincipleExample illustrates a principle with a concrete case.
type PrincipleExample struct {
	// Title summarizes the example.
	Title string `json:"title"`

	// Description provides detail about how the principle applies.
	Description string `json:"description"`

	// ImageURL points to a visual demonstration.
	ImageURL string `json:"imageUrl,omitempty" jsonschema:"format=uri"`

	// Good indicates whether this is a positive or negative example.
	Good bool `json:"good"`
}
