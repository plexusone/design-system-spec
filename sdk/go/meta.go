// Package dss provides Go types for the Design System Specification.
package dss

// Meta contains system-level metadata for a design system.
type Meta struct {
	// Name is the design system name (e.g., "Material Design", "Carbon").
	Name string `json:"name"`

	// Version follows semantic versioning (e.g., "1.0.0").
	Version string `json:"version"`

	// Description provides a brief overview of the design system.
	Description string `json:"description,omitempty"`

	// MaturityLevel indicates system maturity (1=experimental, 5=stable).
	MaturityLevel int `json:"maturityLevel,omitempty" jsonschema:"minimum=1,maximum=5"`

	// Repository is the URL to the source repository.
	Repository string `json:"repository,omitempty" jsonschema:"format=uri"`

	// Documentation is the URL to the documentation site.
	Documentation string `json:"documentation,omitempty" jsonschema:"format=uri"`

	// Maintainers lists the people or teams responsible for the system.
	Maintainers []Maintainer `json:"maintainers,omitempty"`

	// License specifies the license under which the system is distributed.
	License string `json:"license,omitempty"`
}

// Maintainer represents a person or team responsible for the design system.
type Maintainer struct {
	// Name of the maintainer or team.
	Name string `json:"name"`

	// Email contact for the maintainer.
	Email string `json:"email,omitempty" jsonschema:"format=email"`

	// URL to the maintainer's profile or team page.
	URL string `json:"url,omitempty" jsonschema:"format=uri"`
}
