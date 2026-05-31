package dss

import "testing"

func TestValidateContract(t *testing.T) {
	tests := []struct {
		name           string
		component      Component
		wantPassed     bool
		wantErrors     int
		wantWarnings   int
	}{
		{
			name: "no contract",
			component: Component{
				ID:   "button",
				Name: "Button",
			},
			wantPassed:   true,
			wantErrors:   0,
			wantWarnings: 0,
		},
		{
			name: "valid contract",
			component: Component{
				ID:   "button",
				Name: "Button",
				ThemingContract: &ThemingContract{
					Prefix: "--btn",
					Tokens: []ThemeToken{
						{
							ID:           "background",
							CSSProperty:  "--btn-background",
							Semantic:     "primary",
							DefaultLight: "#0066CC",
							DefaultDark:  "#3399FF",
						},
					},
				},
			},
			wantPassed:   true,
			wantErrors:   0,
			wantWarnings: 0,
		},
		{
			name: "missing prefix",
			component: Component{
				ID:   "button",
				Name: "Button",
				ThemingContract: &ThemingContract{
					Tokens: []ThemeToken{
						{
							ID:          "background",
							CSSProperty: "--btn-background",
						},
					},
				},
			},
			wantPassed:   false,
			wantErrors:   1,
			wantWarnings: 1, // missing defaults (combined into one warning)
		},
		{
			name: "prefix without --",
			component: Component{
				ID:   "button",
				Name: "Button",
				ThemingContract: &ThemingContract{
					Prefix: "btn",
					Tokens: []ThemeToken{
						{
							ID:          "background",
							CSSProperty: "--btn-background",
						},
					},
				},
			},
			wantPassed:   false,
			wantErrors:   1,
			wantWarnings: 2, // missing defaults
		},
		{
			name: "empty tokens",
			component: Component{
				ID:   "button",
				Name: "Button",
				ThemingContract: &ThemingContract{
					Prefix: "--btn",
					Tokens: []ThemeToken{},
				},
			},
			wantPassed:   true,
			wantErrors:   0,
			wantWarnings: 1, // no tokens warning
		},
		{
			name: "duplicate token id",
			component: Component{
				ID:   "button",
				Name: "Button",
				ThemingContract: &ThemingContract{
					Prefix: "--btn",
					Tokens: []ThemeToken{
						{ID: "bg", CSSProperty: "--btn-bg", DefaultLight: "#000", DefaultDark: "#fff"},
						{ID: "bg", CSSProperty: "--btn-bg-2", DefaultLight: "#000", DefaultDark: "#fff"},
					},
				},
			},
			wantPassed:   false,
			wantErrors:   1,
			wantWarnings: 0,
		},
		{
			name: "missing token id",
			component: Component{
				ID:   "button",
				Name: "Button",
				ThemingContract: &ThemingContract{
					Prefix: "--btn",
					Tokens: []ThemeToken{
						{CSSProperty: "--btn-bg", DefaultLight: "#000", DefaultDark: "#fff"},
					},
				},
			},
			wantPassed:   false,
			wantErrors:   1,
			wantWarnings: 0,
		},
		{
			name: "missing cssProperty",
			component: Component{
				ID:   "button",
				Name: "Button",
				ThemingContract: &ThemingContract{
					Prefix: "--btn",
					Tokens: []ThemeToken{
						{ID: "bg", DefaultLight: "#000", DefaultDark: "#fff"},
					},
				},
			},
			wantPassed:   false,
			wantErrors:   1,
			wantWarnings: 0,
		},
		{
			name: "css property doesn't match prefix",
			component: Component{
				ID:   "button",
				Name: "Button",
				ThemingContract: &ThemingContract{
					Prefix: "--btn",
					Tokens: []ThemeToken{
						{ID: "bg", CSSProperty: "--card-bg", DefaultLight: "#000", DefaultDark: "#fff"},
					},
				},
			},
			wantPassed:   true,
			wantErrors:   0,
			wantWarnings: 1,
		},
		{
			name: "invalid semantic",
			component: Component{
				ID:   "button",
				Name: "Button",
				ThemingContract: &ThemingContract{
					Prefix: "--btn",
					Tokens: []ThemeToken{
						{ID: "bg", CSSProperty: "--btn-bg", Semantic: "invalid-semantic", DefaultLight: "#000", DefaultDark: "#fff"},
					},
				},
			},
			wantPassed:   true,
			wantErrors:   0,
			wantWarnings: 1,
		},
		{
			name: "missing defaultLight",
			component: Component{
				ID:   "button",
				Name: "Button",
				ThemingContract: &ThemingContract{
					Prefix: "--btn",
					Tokens: []ThemeToken{
						{ID: "bg", CSSProperty: "--btn-bg", DefaultDark: "#fff"},
					},
				},
			},
			wantPassed:   true,
			wantErrors:   0,
			wantWarnings: 1,
		},
		{
			name: "missing both defaults",
			component: Component{
				ID:   "button",
				Name: "Button",
				ThemingContract: &ThemingContract{
					Prefix: "--btn",
					Tokens: []ThemeToken{
						{ID: "bg", CSSProperty: "--btn-bg"},
					},
				},
			},
			wantPassed:   true,
			wantErrors:   0,
			wantWarnings: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := ValidateContract(&tt.component)

			if report.Passed != tt.wantPassed {
				t.Errorf("Passed = %v, want %v", report.Passed, tt.wantPassed)
			}
			if len(report.Errors) != tt.wantErrors {
				t.Errorf("Errors = %d, want %d; errors: %+v", len(report.Errors), tt.wantErrors, report.Errors)
			}
			if len(report.Warnings) != tt.wantWarnings {
				t.Errorf("Warnings = %d, want %d; warnings: %+v", len(report.Warnings), tt.wantWarnings, report.Warnings)
			}
		})
	}
}

func TestValidateAllContracts(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
		Components: []Component{
			{ID: "button", Name: "Button"}, // No contract
			{
				ID:   "card",
				Name: "Card",
				ThemingContract: &ThemingContract{
					Prefix: "--card",
					Tokens: []ThemeToken{
						{ID: "bg", CSSProperty: "--card-bg", DefaultLight: "#fff", DefaultDark: "#000"},
					},
				},
			},
			{
				ID:   "input",
				Name: "Input",
				ThemingContract: &ThemingContract{
					Prefix: "--input",
					Tokens: []ThemeToken{
						{ID: "border", CSSProperty: "--input-border", DefaultLight: "#ccc", DefaultDark: "#333"},
					},
				},
			},
		},
	}

	reports := ValidateAllContracts(ds)

	// Should only have reports for components with contracts
	if len(reports) != 2 {
		t.Errorf("Reports count = %d, want 2", len(reports))
	}

	// Both should pass
	for _, r := range reports {
		if !r.Passed {
			t.Errorf("Component %s failed validation", r.ComponentID)
		}
	}
}

func TestTotalErrorsAndWarnings(t *testing.T) {
	reports := []*ContractReport{
		{ComponentID: "a", Passed: false, Errors: []ContractIssue{{}, {}}, Warnings: []ContractIssue{{}}},
		{ComponentID: "b", Passed: true, Errors: nil, Warnings: []ContractIssue{{}, {}, {}}},
		{ComponentID: "c", Passed: false, Errors: []ContractIssue{{}}, Warnings: nil},
	}

	if got := TotalErrors(reports); got != 3 {
		t.Errorf("TotalErrors = %d, want 3", got)
	}
	if got := TotalWarnings(reports); got != 4 {
		t.Errorf("TotalWarnings = %d, want 4", got)
	}
}
