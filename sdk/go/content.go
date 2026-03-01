package dss

// Content defines voice, tone, and content guidelines.
type Content struct {
	// Voice describes the brand's personality and communication style.
	Voice *VoiceGuidelines `json:"voice,omitempty"`

	// Tone provides guidance on adapting voice to different contexts.
	Tone []ToneGuideline `json:"tone,omitempty"`

	// Terminology defines preferred and avoided terms.
	Terminology *Terminology `json:"terminology,omitempty"`

	// Microcopy provides templates for common UI text.
	Microcopy []MicrocopyGuideline `json:"microcopy,omitempty"`

	// WritingStyle provides general writing guidelines.
	WritingStyle *WritingStyle `json:"writingStyle,omitempty"`
}

// VoiceGuidelines describes the brand's communication personality.
type VoiceGuidelines struct {
	// Description provides an overview of the voice.
	Description string `json:"description"`

	// Attributes lists key voice characteristics.
	Attributes []VoiceAttribute `json:"attributes"`

	// Examples demonstrate the voice in action.
	Examples []VoiceExample `json:"examples,omitempty"`
}

// VoiceAttribute describes a characteristic of the brand voice.
type VoiceAttribute struct {
	// Name is the attribute name (e.g., "Friendly", "Professional").
	Name string `json:"name"`

	// Description explains the attribute.
	Description string `json:"description"`

	// DoExamples show correct usage.
	DoExamples []string `json:"doExamples,omitempty"`

	// DontExamples show incorrect usage.
	DontExamples []string `json:"dontExamples,omitempty"`
}

// VoiceExample demonstrates the brand voice.
type VoiceExample struct {
	// Context describes when this example applies.
	Context string `json:"context"`

	// Example is the sample text.
	Example string `json:"example"`
}

// ToneGuideline provides guidance for specific contexts.
type ToneGuideline struct {
	// Context describes the situation (e.g., "error messages", "success states").
	Context string `json:"context"`

	// Description explains how to adapt tone.
	Description string `json:"description"`

	// Attributes lists tone characteristics for this context.
	Attributes []string `json:"attributes,omitempty"`

	// Examples demonstrate appropriate tone.
	Examples []string `json:"examples,omitempty"`
}

// Terminology defines preferred and avoided terms.
type Terminology struct {
	// PreferredTerms lists terms to use.
	PreferredTerms []TermDefinition `json:"preferredTerms,omitempty"`

	// AvoidedTerms lists terms to avoid.
	AvoidedTerms []TermDefinition `json:"avoidedTerms,omitempty"`

	// Glossary provides definitions for domain-specific terms.
	Glossary []TermDefinition `json:"glossary,omitempty"`
}

// TermDefinition defines a term and its usage.
type TermDefinition struct {
	// Term is the word or phrase.
	Term string `json:"term"`

	// Definition explains the term.
	Definition string `json:"definition,omitempty"`

	// Usage explains when and how to use (or avoid) the term.
	Usage string `json:"usage,omitempty"`

	// Alternatives lists alternative terms to use instead (for avoided terms).
	Alternatives []string `json:"alternatives,omitempty"`
}

// MicrocopyGuideline provides templates for common UI text patterns.
type MicrocopyGuideline struct {
	// ID is a unique identifier.
	ID string `json:"id"`

	// Context describes when this microcopy applies.
	Context string `json:"context"`

	// Template is the text template (may include placeholders).
	Template string `json:"template"`

	// Variables lists placeholder variables in the template.
	Variables []MicrocopyVariable `json:"variables,omitempty"`

	// Examples show the template filled in.
	Examples []string `json:"examples,omitempty"`

	// Notes provides additional guidance.
	Notes string `json:"notes,omitempty"`
}

// MicrocopyVariable describes a placeholder in a microcopy template.
type MicrocopyVariable struct {
	// Name is the variable name.
	Name string `json:"name"`

	// Description explains what value to use.
	Description string `json:"description"`

	// Examples show sample values.
	Examples []string `json:"examples,omitempty"`
}

// WritingStyle provides general writing guidelines.
type WritingStyle struct {
	// Capitalization describes capitalization rules.
	Capitalization string `json:"capitalization,omitempty"`

	// Punctuation describes punctuation guidelines.
	Punctuation string `json:"punctuation,omitempty"`

	// Contractions indicates whether contractions are allowed.
	Contractions string `json:"contractions,omitempty"`

	// Numbers describes how to format numbers.
	Numbers string `json:"numbers,omitempty"`

	// DateTimeFormat describes date/time formatting.
	DateTimeFormat string `json:"dateTimeFormat,omitempty"`

	// ActiveVoice indicates preference for active voice.
	ActiveVoice bool `json:"activeVoice,omitempty"`

	// MaxSentenceLength is the recommended max sentence length.
	MaxSentenceLength int `json:"maxSentenceLength,omitempty"`
}
