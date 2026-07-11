package dss

// Validators defines external validation tools for requirements delegation.
// DSS spec defines requirements (e.g., WCAG AA); validators perform the actual checks.
type Validators struct {
	// Accessibility references accessibility validation tools.
	// Requirements are defined in the Accessibility section; this specifies how to validate.
	Accessibility *AccessibilityValidator `json:"accessibility,omitempty"`

	// APIStyle references API style guide validators.
	APIStyle *APIStyleValidator `json:"apiStyle,omitempty"`

	// Custom references custom validation tools for other domains.
	Custom []CustomValidator `json:"custom,omitempty"`
}

// AccessibilityValidator references an accessibility validation tool.
// Works in conjunction with the Accessibility section which defines requirements.
type AccessibilityValidator struct {
	// Tool is the validator tool identifier.
	// Examples: "agent-a11y", "axe-core", "pa11y", "lighthouse"
	Tool string `json:"tool"`

	// Type specifies how to invoke the tool.
	Type ValidatorType `json:"type"`

	// Command is the command or path to run (for mcp/cli types).
	// For MCP: the server command (e.g., "agent-a11y mcp serve")
	// For CLI: the executable (e.g., "npx axe")
	Command string `json:"command,omitempty"`

	// Standards lists which accessibility standards to validate.
	// If empty, inherits from Accessibility.WCAGLevel and WCAGVersion.
	// Examples: ["WCAG2.1-AA", "WCAG2.2-AAA", "Section508"]
	Standards []string `json:"standards,omitempty"`

	// Checks specifies which checks to run.
	// Tool-specific; if empty, runs all checks.
	Checks []string `json:"checks,omitempty"`

	// Required indicates whether validation must pass for compliance.
	Required bool `json:"required,omitempty"`
}

// APIStyleValidator references an API style guide validator.
type APIStyleValidator struct {
	// Tool is the validator tool identifier.
	// Examples: "spectral", "vacuum", "redocly"
	Tool string `json:"tool"`

	// Type specifies how to invoke the tool.
	Type ValidatorType `json:"type"`

	// Command is the command or path to run.
	Command string `json:"command,omitempty"`

	// Ruleset is the path or URL to the style guide ruleset.
	// For Spectral: path to .spectral.yaml or ruleset URL
	// For api-style-spec: path to style-spec.yaml
	Ruleset string `json:"ruleset,omitempty"`

	// Rules specifies which rules to enable/disable.
	// Tool-specific configuration.
	Rules map[string]RuleSeverity `json:"rules,omitempty"`

	// Required indicates whether validation must pass for compliance.
	Required bool `json:"required,omitempty"`
}

// CustomValidator references a custom validation tool for other domains.
type CustomValidator struct {
	// ID is a unique identifier for this validator.
	ID string `json:"id"`

	// Name is a human-readable name.
	Name string `json:"name"`

	// Domain describes what this validator checks.
	// Examples: "security", "performance", "i18n", "seo"
	Domain string `json:"domain"`

	// Tool is the validator tool identifier.
	Tool string `json:"tool"`

	// Type specifies how to invoke the tool.
	Type ValidatorType `json:"type"`

	// Command is the command or path to run.
	Command string `json:"command,omitempty"`

	// Config provides tool-specific configuration.
	Config map[string]any `json:"config,omitempty"`

	// Required indicates whether validation must pass for compliance.
	Required bool `json:"required,omitempty"`
}

// ValidatorType specifies how to invoke a validator.
type ValidatorType string

const (
	// ValidatorTypeMCP invokes the validator as an MCP server.
	// The tool runs as a subprocess exposing MCP tools.
	ValidatorTypeMCP ValidatorType = "mcp"

	// ValidatorTypeCLI invokes the validator as a command-line tool.
	// Output is parsed from stdout/stderr.
	ValidatorTypeCLI ValidatorType = "cli"

	// ValidatorTypeNPM invokes the validator via npm/npx.
	// Useful for JavaScript-based tools like axe-core.
	ValidatorTypeNPM ValidatorType = "npm"

	// ValidatorTypeAPI invokes the validator via HTTP API.
	// Command field contains the base URL.
	ValidatorTypeAPI ValidatorType = "api"

	// ValidatorTypeLibrary invokes the validator as a Go library.
	// Command field contains the import path.
	ValidatorTypeLibrary ValidatorType = "library"
)

// RuleSeverity specifies the severity for a validation rule.
type RuleSeverity string

const (
	RuleSeverityError   RuleSeverity = "error"
	RuleSeverityWarning RuleSeverity = "warning"
	RuleSeverityInfo    RuleSeverity = "info"
	RuleSeverityOff     RuleSeverity = "off"
)

// ValidatorResult represents the result of running an external validator.
type ValidatorResult struct {
	// ValidatorID identifies which validator produced this result.
	ValidatorID string `json:"validatorId"`

	// Tool is the tool name.
	Tool string `json:"tool"`

	// Domain is the validation domain (accessibility, api, custom).
	Domain string `json:"domain"`

	// Passed indicates overall pass/fail status.
	Passed bool `json:"passed"`

	// Issues contains validation findings.
	Issues []ValidatorIssue `json:"issues,omitempty"`

	// Summary provides aggregate statistics.
	Summary ValidatorSummary `json:"summary"`

	// Metadata contains tool-specific result data.
	Metadata map[string]any `json:"metadata,omitempty"`
}

// ValidatorIssue represents a single validation finding.
type ValidatorIssue struct {
	// Rule is the rule ID that was violated.
	Rule string `json:"rule"`

	// Severity is the issue severity.
	Severity string `json:"severity"` // error, warning, info

	// Message describes the issue.
	Message string `json:"message"`

	// Path is the location in code/spec where the issue was found.
	Path string `json:"path,omitempty"`

	// Line is the line number if applicable.
	Line int `json:"line,omitempty"`

	// Element is the DOM element or selector if applicable (for a11y).
	Element string `json:"element,omitempty"`

	// Help provides guidance on fixing the issue.
	Help string `json:"help,omitempty"`

	// HelpURL links to documentation.
	HelpURL string `json:"helpUrl,omitempty"`
}

// ValidatorSummary provides aggregate statistics.
type ValidatorSummary struct {
	// Errors is the count of error-level issues.
	Errors int `json:"errors"`

	// Warnings is the count of warning-level issues.
	Warnings int `json:"warnings"`

	// Infos is the count of info-level issues.
	Infos int `json:"infos"`

	// Passed is the count of passed checks.
	Passed int `json:"passed,omitempty"`

	// Total is the total number of checks run.
	Total int `json:"total,omitempty"`
}
