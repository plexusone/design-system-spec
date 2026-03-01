package dss

import (
	"github.com/invopop/jsonschema"
)

// GenerateSchema generates a JSON Schema for the DesignSystem type.
func GenerateSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{
		DoNotReference:             false,
		ExpandedStruct:             false,
		AllowAdditionalProperties:  false,
		RequiredFromJSONSchemaTags: true,
	}
	return r.Reflect(&DesignSystem{})
}

// GenerateFoundationsSchema generates a JSON Schema for the Foundations type.
func GenerateFoundationsSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{
		DoNotReference:             false,
		AllowAdditionalProperties:  false,
		RequiredFromJSONSchemaTags: true,
	}
	return r.Reflect(&Foundations{})
}

// GenerateComponentSchema generates a JSON Schema for the Component type.
func GenerateComponentSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{
		DoNotReference:             false,
		AllowAdditionalProperties:  false,
		RequiredFromJSONSchemaTags: true,
	}
	return r.Reflect(&Component{})
}

// GenerateMetaSchema generates a JSON Schema for the Meta type.
func GenerateMetaSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{
		DoNotReference:             false,
		AllowAdditionalProperties:  false,
		RequiredFromJSONSchemaTags: true,
	}
	return r.Reflect(&Meta{})
}

// GeneratePatternSchema generates a JSON Schema for the Pattern type.
func GeneratePatternSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{
		DoNotReference:             false,
		AllowAdditionalProperties:  false,
		RequiredFromJSONSchemaTags: true,
	}
	return r.Reflect(&Pattern{})
}

// GenerateTemplateSchema generates a JSON Schema for the Template type.
func GenerateTemplateSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{
		DoNotReference:             false,
		AllowAdditionalProperties:  false,
		RequiredFromJSONSchemaTags: true,
	}
	return r.Reflect(&Template{})
}

// GenerateAccessibilitySchema generates a JSON Schema for the Accessibility type.
func GenerateAccessibilitySchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{
		DoNotReference:             false,
		AllowAdditionalProperties:  false,
		RequiredFromJSONSchemaTags: true,
	}
	return r.Reflect(&Accessibility{})
}

// GenerateGovernanceSchema generates a JSON Schema for the Governance type.
func GenerateGovernanceSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{
		DoNotReference:             false,
		AllowAdditionalProperties:  false,
		RequiredFromJSONSchemaTags: true,
	}
	return r.Reflect(&Governance{})
}
