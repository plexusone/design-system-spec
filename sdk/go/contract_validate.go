package dss

import (
	"fmt"
	"strings"
)

// ContractReport contains the results of validating a theming contract.
type ContractReport struct {
	// ComponentID is the ID of the validated component.
	ComponentID string `json:"componentId"`

	// Passed indicates whether validation passed without errors.
	Passed bool `json:"passed"`

	// Errors are critical issues that must be fixed.
	Errors []ContractIssue `json:"errors,omitempty"`

	// Warnings are non-critical issues that should be addressed.
	Warnings []ContractIssue `json:"warnings,omitempty"`
}

// ContractIssue represents a validation error or warning.
type ContractIssue struct {
	// Field is the JSON path to the problematic field.
	Field string `json:"field"`

	// Message describes the issue.
	Message string `json:"message"`

	// TokenID is the affected token ID, if applicable.
	TokenID string `json:"tokenId,omitempty"`
}

// ValidateContract validates a component's theming contract for completeness and correctness.
func ValidateContract(component *Component) *ContractReport {
	report := &ContractReport{
		ComponentID: component.ID,
		Passed:      true,
	}

	contract := component.ThemingContract
	if contract == nil {
		// No contract to validate
		return report
	}

	// Validate prefix
	if contract.Prefix == "" {
		report.addError("themingContract.prefix", "", "prefix is required")
	} else if !strings.HasPrefix(contract.Prefix, "--") {
		report.addError("themingContract.prefix", "", "prefix must start with '--'")
	}

	// Validate tokens
	if len(contract.Tokens) == 0 {
		report.addWarning("themingContract.tokens", "", "no tokens defined in contract")
	}

	tokenIDs := make(map[string]bool)
	for i, token := range contract.Tokens {
		field := fmt.Sprintf("themingContract.tokens[%d]", i)

		// Check for duplicate IDs
		if token.ID == "" {
			report.addError(field+".id", "", "token id is required")
		} else if tokenIDs[token.ID] {
			report.addError(field+".id", token.ID, "duplicate token id")
		} else {
			tokenIDs[token.ID] = true
		}

		// Validate CSS property
		if token.CSSProperty == "" {
			report.addError(field+".cssProperty", token.ID, "cssProperty is required")
		} else {
			// Check CSS property follows prefix convention
			if contract.Prefix != "" && !strings.HasPrefix(token.CSSProperty, contract.Prefix) {
				report.addWarning(field+".cssProperty", token.ID,
					fmt.Sprintf("cssProperty should start with prefix '%s'", contract.Prefix))
			}
		}

		// Validate semantic value
		if !IsValidSemantic(token.Semantic) {
			report.addWarning(field+".semantic", token.ID,
				fmt.Sprintf("unknown semantic value '%s'; expected one of: %s",
					token.Semantic, strings.Join(ValidSemantics, ", ")))
		}

		// Check for default values
		if token.DefaultLight == "" && token.DefaultDark == "" {
			report.addWarning(field, token.ID, "no default values (defaultLight/defaultDark) provided")
		} else if token.DefaultLight == "" {
			report.addWarning(field+".defaultLight", token.ID, "defaultLight not provided")
		} else if token.DefaultDark == "" {
			report.addWarning(field+".defaultDark", token.ID, "defaultDark not provided")
		}
	}

	return report
}

// ValidateAllContracts validates all components' theming contracts in a design system.
func ValidateAllContracts(ds *DesignSystem) []*ContractReport {
	var reports []*ContractReport

	for i := range ds.Components {
		component := &ds.Components[i]
		if component.ThemingContract != nil {
			reports = append(reports, ValidateContract(component))
		}
	}

	return reports
}

// addError adds an error to the report and marks it as failed.
func (r *ContractReport) addError(field, tokenID, message string) {
	r.Passed = false
	r.Errors = append(r.Errors, ContractIssue{
		Field:   field,
		TokenID: tokenID,
		Message: message,
	})
}

// addWarning adds a warning to the report (does not fail validation).
func (r *ContractReport) addWarning(field, tokenID, message string) {
	r.Warnings = append(r.Warnings, ContractIssue{
		Field:   field,
		TokenID: tokenID,
		Message: message,
	})
}

// TotalErrors returns the total number of errors across all reports.
func TotalErrors(reports []*ContractReport) int {
	total := 0
	for _, r := range reports {
		total += len(r.Errors)
	}
	return total
}

// TotalWarnings returns the total number of warnings across all reports.
func TotalWarnings(reports []*ContractReport) int {
	total := 0
	for _, r := range reports {
		total += len(r.Warnings)
	}
	return total
}
