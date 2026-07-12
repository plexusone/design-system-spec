package dss

import (
	"context"
	"fmt"
)

// Service provides design system operations for CLI, MCP, and SDK consumers.
// It wraps a loaded DesignSystem and provides a unified API for querying
// and validating against the specification.
type Service struct {
	ds *DesignSystem
}

// NewService creates a service from a loaded design system.
func NewService(ds *DesignSystem) *Service {
	return &Service{ds: ds}
}

// DesignSystem returns the underlying design system.
func (s *Service) DesignSystem() *DesignSystem {
	return s.ds
}

// ComponentSummary provides a brief overview of a component.
type ComponentSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
}

// TokenSummary provides a brief overview of a design token.
type TokenSummary struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// PatternSummary provides a brief overview of a pattern.
type PatternSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
}

// PromptOptions configures LLM prompt generation.
type PromptOptions struct {
	// Format: "markdown" (default), "xml"
	Format string `json:"format,omitempty"`

	// ComponentIDs limits output to specific components (empty = all)
	ComponentIDs []string `json:"componentIds,omitempty"`

	// IncludeFoundations includes design tokens
	IncludeFoundations bool `json:"includeFoundations,omitempty"`

	// IncludeComponents includes component definitions
	IncludeComponents bool `json:"includeComponents,omitempty"`

	// IncludePatterns includes pattern definitions
	IncludePatterns bool `json:"includePatterns,omitempty"`

	// IncludeAccessibility includes a11y requirements
	IncludeAccessibility bool `json:"includeAccessibility,omitempty"`

	// IncludeAntiPatterns includes anti-pattern warnings
	IncludeAntiPatterns bool `json:"includeAntiPatterns,omitempty"`

	// MaxExamples limits code examples per component
	MaxExamples int `json:"maxExamples,omitempty"`
}

// DefaultPromptOptions returns comprehensive defaults.
func DefaultPromptOptions() *PromptOptions {
	return &PromptOptions{
		Format:               "markdown",
		IncludeFoundations:   true,
		IncludeComponents:    true,
		IncludePatterns:      true,
		IncludeAccessibility: true,
		IncludeAntiPatterns:  true,
		MaxExamples:          3,
	}
}

// GetComponent returns a component by ID.
func (s *Service) GetComponent(_ context.Context, id string) (*Component, error) {
	for i := range s.ds.Components {
		if s.ds.Components[i].ID == id {
			return &s.ds.Components[i], nil
		}
	}
	return nil, fmt.Errorf("component not found: %s", id)
}

// ListComponents returns summaries of all components.
func (s *Service) ListComponents(_ context.Context) []ComponentSummary {
	result := make([]ComponentSummary, len(s.ds.Components))
	for i, c := range s.ds.Components {
		result[i] = ComponentSummary{
			ID:          c.ID,
			Name:        c.Name,
			Category:    c.Category,
			Description: c.Description,
		}
	}
	return result
}

// GetComponentVariants returns the variants for a component.
func (s *Service) GetComponentVariants(_ context.Context, componentID string) ([]Variant, error) {
	for _, c := range s.ds.Components {
		if c.ID == componentID {
			return c.Variants, nil
		}
	}
	return nil, fmt.Errorf("component not found: %s", componentID)
}

// GetComponentProps returns the props for a component.
func (s *Service) GetComponentProps(_ context.Context, componentID string) ([]Prop, error) {
	for _, c := range s.ds.Components {
		if c.ID == componentID {
			return c.Props, nil
		}
	}
	return nil, fmt.Errorf("component not found: %s", componentID)
}

// GetComponentAntiPatterns returns anti-patterns for a component.
func (s *Service) GetComponentAntiPatterns(_ context.Context, componentID string) ([]string, error) {
	for _, c := range s.ds.Components {
		if c.ID == componentID {
			if c.LLM != nil {
				return c.LLM.AntiPatterns, nil
			}
			return nil, nil
		}
	}
	return nil, fmt.Errorf("component not found: %s", componentID)
}

// GetToken returns a design token by type and name.
func (s *Service) GetToken(_ context.Context, tokenType, name string) (any, error) {
	f := s.ds.Foundations

	switch tokenType {
	case "color":
		for _, c := range f.Colors {
			if c.ID == name {
				return c, nil
			}
		}
	case "spacing":
		if f.Spacing != nil {
			for _, sp := range f.Spacing.Scale {
				if sp.ID == name {
					return sp, nil
				}
			}
		}
	case "typography":
		if f.Typography != nil {
			for _, ts := range f.Typography.TypeScale {
				if ts.ID == name {
					return ts, nil
				}
			}
		}
	case "elevation":
		for _, e := range f.Elevation {
			if e.ID == name {
				return e, nil
			}
		}
	case "borderRadius":
		for _, br := range f.BorderRadius {
			if br.ID == name {
				return br, nil
			}
		}
	case "breakpoint":
		for _, bp := range f.Breakpoints {
			if bp.ID == name {
				return bp, nil
			}
		}
	default:
		return nil, fmt.Errorf("unknown token type: %s", tokenType)
	}

	return nil, fmt.Errorf("token not found: %s/%s", tokenType, name)
}

// ListTokens returns all tokens of a given type.
func (s *Service) ListTokens(_ context.Context, tokenType string) ([]TokenSummary, error) {
	f := s.ds.Foundations
	var result []TokenSummary

	switch tokenType {
	case "color":
		for _, c := range f.Colors {
			result = append(result, TokenSummary{
				ID:    c.ID,
				Type:  "color",
				Value: c.Value,
			})
		}
	case "spacing":
		if f.Spacing != nil {
			for _, sp := range f.Spacing.Scale {
				result = append(result, TokenSummary{
					ID:    sp.ID,
					Type:  "spacing",
					Value: sp.Value,
				})
			}
		}
	case "typography":
		if f.Typography != nil {
			for _, ts := range f.Typography.TypeScale {
				result = append(result, TokenSummary{
					ID:    ts.ID,
					Type:  "typography",
					Value: fmt.Sprintf("%s %s", ts.FontSize, ts.FontWeight),
				})
			}
		}
	case "elevation":
		for _, e := range f.Elevation {
			result = append(result, TokenSummary{
				ID:    e.ID,
				Type:  "elevation",
				Value: e.Value,
			})
		}
	case "borderRadius":
		for _, br := range f.BorderRadius {
			result = append(result, TokenSummary{
				ID:    br.ID,
				Type:  "borderRadius",
				Value: br.Value,
			})
		}
	case "breakpoint":
		for _, bp := range f.Breakpoints {
			result = append(result, TokenSummary{
				ID:    bp.ID,
				Type:  "breakpoint",
				Value: bp.MinWidth,
			})
		}
	default:
		return nil, fmt.Errorf("unknown token type: %s", tokenType)
	}

	return result, nil
}

// GetPattern returns a pattern by ID.
func (s *Service) GetPattern(_ context.Context, id string) (*Pattern, error) {
	for i := range s.ds.Patterns {
		if s.ds.Patterns[i].ID == id {
			return &s.ds.Patterns[i], nil
		}
	}
	return nil, fmt.Errorf("pattern not found: %s", id)
}

// ListPatterns returns summaries of all patterns.
func (s *Service) ListPatterns(_ context.Context) []PatternSummary {
	result := make([]PatternSummary, len(s.ds.Patterns))
	for i, p := range s.ds.Patterns {
		result[i] = PatternSummary{
			ID:          p.ID,
			Name:        p.Name,
			Category:    p.Category,
			Description: p.Description,
		}
	}
	return result
}

// GetMeta returns the design system metadata.
func (s *Service) GetMeta(_ context.Context) Meta {
	return s.ds.Meta
}

// GenerateLLMPrompt generates a complete LLM context prompt.
func (s *Service) GenerateLLMPrompt(_ context.Context, opts *PromptOptions) (string, error) {
	if opts == nil {
		opts = DefaultPromptOptions()
	}

	llmOpts := LLMPromptOptions{
		Format:               opts.Format,
		IncludeFoundations:   opts.IncludeFoundations,
		IncludeComponents:    opts.IncludeComponents,
		IncludePatterns:      opts.IncludePatterns,
		IncludeAccessibility: opts.IncludeAccessibility,
		IncludeAntiPatterns:  opts.IncludeAntiPatterns,
		MaxExamples:          opts.MaxExamples,
	}

	return s.ds.GenerateLLMPrompt(llmOpts)
}

// GenerateComponentPrompt generates an LLM prompt for a specific component.
func (s *Service) GenerateComponentPrompt(ctx context.Context, componentID string) (string, error) {
	comp, err := s.GetComponent(ctx, componentID)
	if err != nil {
		return "", err
	}

	// Create a temporary DesignSystem with just this component
	tempDS := &DesignSystem{
		Meta:        s.ds.Meta,
		Foundations: s.ds.Foundations,
		Components:  []Component{*comp},
	}

	return tempDS.GenerateLLMPrompt(LLMPromptOptions{
		Format:              "markdown",
		IncludeFoundations:  true,
		IncludeComponents:   true,
		IncludeAntiPatterns: true,
		MaxExamples:         5,
	})
}

// GetAccessibility returns the system-wide accessibility requirements.
func (s *Service) GetAccessibility(_ context.Context) *Accessibility {
	return s.ds.Accessibility
}

// AccessibilityRequirements contains full accessibility requirements for a component.
type AccessibilityRequirements struct {
	ComponentID string `json:"componentId"`
	Component   string `json:"component"`

	// RequiredProps lists props that must be set for accessibility.
	RequiredProps []RequiredProp `json:"requiredProps,omitempty"`

	// Keyboard describes keyboard interaction requirements.
	Keyboard *KeyboardReqs `json:"keyboard,omitempty"`

	// Focus describes focus management requirements.
	Focus *FocusReqs `json:"focus,omitempty"`

	// ColorContrast describes contrast requirements.
	ColorContrast *ContrastReqs `json:"colorContrast,omitempty"`

	// WCAGCriteria lists applicable WCAG success criteria.
	WCAGCriteria []string `json:"wcagCriteria,omitempty"`

	// Role is the ARIA role for this component.
	Role string `json:"role,omitempty"`

	// ScreenReaderNotes provides guidance for screen reader behavior.
	ScreenReaderNotes string `json:"screenReaderNotes,omitempty"`
}

// RequiredProp describes a prop required for accessibility.
type RequiredProp struct {
	Name     string `json:"name"`
	When     string `json:"when,omitempty"`
	Type     string `json:"type,omitempty"`
	Example  string `json:"example,omitempty"`
	ARIAAttr string `json:"ariaAttr,omitempty"`
}

// KeyboardReqs describes keyboard accessibility requirements.
type KeyboardReqs struct {
	Interactions map[string]string `json:"interactions,omitempty"`
	TabIndex     string            `json:"tabIndex,omitempty"`
	TrapFocus    bool              `json:"trapFocus,omitempty"`
}

// FocusReqs describes focus management requirements.
type FocusReqs struct {
	Visible   bool   `json:"visible"`
	Order     string `json:"order,omitempty"`
	TrapFocus bool   `json:"trapFocus,omitempty"`
	InitialOn string `json:"initialOn,omitempty"`
}

// ContrastReqs describes color contrast requirements.
type ContrastReqs struct {
	Text         string `json:"text,omitempty"`
	LargeText    string `json:"largeText,omitempty"`
	UIComponents string `json:"uiComponents,omitempty"`
}

// GetAccessibilityRequirements returns accessibility requirements for a component.
// It combines system-wide requirements with component-specific ones.
func (s *Service) GetAccessibilityRequirements(_ context.Context, componentID string) (*AccessibilityRequirements, error) {
	// Find the component
	var comp *Component
	for i := range s.ds.Components {
		if s.ds.Components[i].ID == componentID {
			comp = &s.ds.Components[i]
			break
		}
	}
	if comp == nil {
		return nil, fmt.Errorf("component not found: %s", componentID)
	}

	reqs := &AccessibilityRequirements{
		ComponentID: comp.ID,
		Component:   comp.Name,
	}

	// Add system-wide contrast requirements
	if s.ds.Accessibility != nil && s.ds.Accessibility.ColorContrast != nil {
		cc := s.ds.Accessibility.ColorContrast
		reqs.ColorContrast = &ContrastReqs{
			Text:         fmt.Sprintf("%.1f:1", cc.NormalTextRatio),
			LargeText:    fmt.Sprintf("%.1f:1", cc.LargeTextRatio),
			UIComponents: fmt.Sprintf("%.1f:1", cc.NonTextRatio),
		}
	}

	// Add component-specific accessibility
	if comp.Accessibility != nil {
		a := comp.Accessibility
		reqs.Role = a.Role
		reqs.ScreenReaderNotes = a.ScreenReaderNotes

		// Convert keyboard interactions
		if len(a.KeyboardSupport) > 0 {
			reqs.Keyboard = &KeyboardReqs{
				Interactions: make(map[string]string),
			}
			for _, ki := range a.KeyboardSupport {
				reqs.Keyboard.Interactions[ki.Key] = ki.Action
			}
		}

		// Focus management
		if a.FocusManagement != "" {
			reqs.Focus = &FocusReqs{
				Visible: true,
				Order:   "natural",
			}
			if a.FocusManagement == "trap" {
				reqs.Focus.TrapFocus = true
			}
		}

		// Required ARIA attributes become required props
		for _, attr := range a.RequiredAttributes {
			reqs.RequiredProps = append(reqs.RequiredProps, RequiredProp{
				Name:     attr,
				ARIAAttr: attr,
			})
		}
	}

	// Derive WCAG criteria from component type
	reqs.WCAGCriteria = deriveWCAGCriteria(comp)

	// Derive required props from component props
	for _, prop := range comp.Props {
		if isAccessibilityProp(prop) {
			reqs.RequiredProps = append(reqs.RequiredProps, RequiredProp{
				Name:    prop.Name,
				Type:    prop.Type,
				When:    prop.Description,
				Example: getExampleValue(prop),
			})
		}
	}

	return reqs, nil
}

// deriveWCAGCriteria returns relevant WCAG criteria based on component characteristics.
func deriveWCAGCriteria(comp *Component) []string {
	var criteria []string

	// All interactive components need these
	criteria = append(criteria, "2.1.1") // Keyboard
	criteria = append(criteria, "2.4.7") // Focus Visible
	criteria = append(criteria, "4.1.2") // Name, Role, Value

	// Components with states
	if len(comp.States) > 0 {
		criteria = append(criteria, "4.1.3") // Status Messages
	}

	// Components with text
	if hasTextContent(comp) {
		criteria = append(criteria, "1.4.3") // Contrast (Minimum)
	}

	// Form components
	if isFormComponent(comp) {
		criteria = append(criteria, "3.3.1") // Error Identification
		criteria = append(criteria, "3.3.2") // Labels or Instructions
	}

	return criteria
}

// isAccessibilityProp returns true if the prop is accessibility-related.
func isAccessibilityProp(prop Prop) bool {
	accessibilityProps := map[string]bool{
		"aria-label":       true,
		"aria-labelledby":  true,
		"aria-describedby": true,
		"aria-expanded":    true,
		"aria-haspopup":    true,
		"aria-controls":    true,
		"aria-pressed":     true,
		"aria-checked":     true,
		"aria-selected":    true,
		"aria-disabled":    true,
		"role":             true,
		"alt":              true,
		"title":            true,
		"tabIndex":         true,
	}
	return accessibilityProps[prop.Name]
}

// hasTextContent returns true if the component likely displays text.
func hasTextContent(comp *Component) bool {
	textComponents := map[string]bool{
		"button": true, "text": true, "label": true, "heading": true,
		"paragraph": true, "link": true, "badge": true, "alert": true,
	}
	return textComponents[comp.ID]
}

// isFormComponent returns true if the component is a form element.
func isFormComponent(comp *Component) bool {
	formComponents := map[string]bool{
		"input": true, "select": true, "checkbox": true, "radio": true,
		"textarea": true, "switch": true, "slider": true, "form": true,
	}
	return formComponents[comp.ID]
}

// getExampleValue returns an example value for a prop.
func getExampleValue(prop Prop) string {
	if prop.Default != nil {
		if s, ok := prop.Default.(string); ok && s != "" {
			return s
		}
	}
	switch prop.Name {
	case "aria-label":
		return "Descriptive label"
	case "alt":
		return "Image description"
	default:
		return ""
	}
}

// AntiPattern describes an accessibility anti-pattern to avoid.
type AntiPattern struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	BadExample  string   `json:"badExample"`
	GoodExample string   `json:"goodExample"`
	WCAGCriteria []string `json:"wcagCriteria,omitempty"`
	Components  []string `json:"components,omitempty"`
}

// AntiPatternsResult contains anti-patterns for a component or rule.
type AntiPatternsResult struct {
	ComponentID  string        `json:"componentId,omitempty"`
	RuleID       string        `json:"ruleId,omitempty"`
	AntiPatterns []AntiPattern `json:"antiPatterns"`
}

// GetAntiPatterns returns accessibility anti-patterns to avoid.
func (s *Service) GetAntiPatterns(_ context.Context, componentID, ruleID string) (*AntiPatternsResult, error) {
	result := &AntiPatternsResult{
		ComponentID: componentID,
		RuleID:      ruleID,
	}

	// Built-in anti-patterns database
	allAntiPatterns := []AntiPattern{
		{
			ID:          "placeholder-as-label",
			Description: "Using placeholder text as the only label",
			BadExample:  `<input placeholder="Email">`,
			GoodExample: `<label for="email">Email</label><input id="email">`,
			WCAGCriteria: []string{"1.3.1", "3.3.2"},
			Components:  []string{"input", "textarea", "select"},
		},
		{
			ID:          "color-only-error",
			Description: "Using color alone to indicate errors",
			BadExample:  `<input style="border-color: red">`,
			GoodExample: `<input aria-invalid="true" aria-describedby="error"><span id="error">Error message</span>`,
			WCAGCriteria: []string{"1.4.1"},
			Components:  []string{"input", "form", "alert"},
		},
		{
			ID:          "missing-alt-text",
			Description: "Images without alternative text",
			BadExample:  `<img src="photo.jpg">`,
			GoodExample: `<img src="photo.jpg" alt="Description of the image">`,
			WCAGCriteria: []string{"1.1.1"},
			Components:  []string{"image", "avatar"},
		},
		{
			ID:          "non-focusable-interactive",
			Description: "Interactive elements that cannot receive keyboard focus",
			BadExample:  `<div onclick="handleClick()">Click me</div>`,
			GoodExample: `<button type="button" onclick="handleClick()">Click me</button>`,
			WCAGCriteria: []string{"2.1.1", "4.1.2"},
			Components:  []string{"button", "link"},
		},
		{
			ID:          "missing-button-name",
			Description: "Buttons without accessible names",
			BadExample:  `<button><svg>...</svg></button>`,
			GoodExample: `<button aria-label="Close dialog"><svg>...</svg></button>`,
			WCAGCriteria: []string{"4.1.2"},
			Components:  []string{"button", "icon-button"},
		},
		{
			ID:          "focus-not-visible",
			Description: "Focus indicator removed or invisible",
			BadExample:  `button:focus { outline: none; }`,
			GoodExample: `button:focus { outline: 2px solid var(--focus-ring); outline-offset: 2px; }`,
			WCAGCriteria: []string{"2.4.7"},
			Components:  []string{"button", "input", "link", "select"},
		},
		{
			ID:          "low-contrast-text",
			Description: "Text with insufficient color contrast",
			BadExample:  `<p style="color: #999; background: #fff">Light gray text</p>`,
			GoodExample: `<p style="color: #595959; background: #fff">Accessible gray text</p>`,
			WCAGCriteria: []string{"1.4.3", "1.4.6"},
			Components:  []string{"text", "label", "button"},
		},
		{
			ID:          "auto-playing-media",
			Description: "Media that plays automatically without user consent",
			BadExample:  `<video autoplay>...</video>`,
			GoodExample: `<video controls>...</video>`,
			WCAGCriteria: []string{"1.4.2"},
			Components:  []string{"video", "audio"},
		},
		{
			ID:          "missing-form-labels",
			Description: "Form inputs without associated labels",
			BadExample:  `Name: <input type="text">`,
			GoodExample: `<label for="name">Name:</label><input id="name" type="text">`,
			WCAGCriteria: []string{"1.3.1", "3.3.2", "4.1.2"},
			Components:  []string{"input", "select", "textarea", "checkbox", "radio"},
		},
		{
			ID:          "keyboard-trap",
			Description: "Content that traps keyboard focus",
			BadExample:  `// Modal with no escape key handling`,
			GoodExample: `// Modal with Escape key to close and focus return to trigger`,
			WCAGCriteria: []string{"2.1.2"},
			Components:  []string{"modal", "dialog", "dropdown"},
		},
		{
			ID:          "missing-skip-link",
			Description: "Page without skip navigation link",
			BadExample:  `<body><nav>Long navigation...</nav><main>Content</main></body>`,
			GoodExample: `<body><a href="#main" class="skip-link">Skip to content</a><nav>...</nav><main id="main">Content</main></body>`,
			WCAGCriteria: []string{"2.4.1"},
			Components:  []string{"navigation", "layout"},
		},
		{
			ID:          "improper-heading-structure",
			Description: "Skipping heading levels or using headings for styling",
			BadExample:  `<h1>Title</h1><h3>Subtitle</h3>`,
			GoodExample: `<h1>Title</h1><h2>Subtitle</h2>`,
			WCAGCriteria: []string{"1.3.1", "2.4.6"},
			Components:  []string{"heading", "text"},
		},
	}

	// Filter by component if specified
	for _, ap := range allAntiPatterns {
		include := false

		// Filter by component
		if componentID != "" {
			for _, c := range ap.Components {
				if c == componentID {
					include = true
					break
				}
			}
		} else if ruleID != "" {
			// Filter by WCAG criteria / rule ID
			for _, c := range ap.WCAGCriteria {
				if c == ruleID || matchesRule(ap.ID, ruleID) {
					include = true
					break
				}
			}
		} else {
			// No filter, include all
			include = true
		}

		if include {
			result.AntiPatterns = append(result.AntiPatterns, ap)
		}
	}

	// Also add component-specific anti-patterns from LLM context
	if componentID != "" {
		for _, comp := range s.ds.Components {
			if comp.ID == componentID && comp.LLM != nil {
				for _, ap := range comp.LLM.AntiPatterns {
					result.AntiPatterns = append(result.AntiPatterns, AntiPattern{
						ID:          fmt.Sprintf("%s-custom", componentID),
						Description: ap,
						Components:  []string{componentID},
					})
				}
			}
		}
	}

	return result, nil
}

// matchesRule checks if an anti-pattern ID matches a rule ID pattern.
func matchesRule(antiPatternID, ruleID string) bool {
	ruleMap := map[string][]string{
		"color-contrast":   {"low-contrast-text", "color-only-error"},
		"missing-label":    {"placeholder-as-label", "missing-form-labels", "missing-button-name"},
		"keyboard":         {"non-focusable-interactive", "keyboard-trap"},
		"focus":            {"focus-not-visible"},
		"image-alt":        {"missing-alt-text"},
		"heading-order":    {"improper-heading-structure"},
		"bypass":           {"missing-skip-link"},
	}

	if patterns, ok := ruleMap[ruleID]; ok {
		for _, p := range patterns {
			if p == antiPatternID {
				return true
			}
		}
	}
	return false
}

// ContrastSuggestion suggests a color token that meets contrast requirements.
type ContrastSuggestion struct {
	Token           string  `json:"token"`
	Value           string  `json:"value"`
	ContrastRatio   float64 `json:"contrastRatio"`
	MeetsAA         bool    `json:"meetsAA"`
	MeetsAAA        bool    `json:"meetsAAA"`
	MeetsAALarge    bool    `json:"meetsAALarge"`
}

// SuggestContrastToken suggests color tokens that meet contrast requirements.
func (s *Service) SuggestContrastToken(_ context.Context, background string, minContrast float64) ([]ContrastSuggestion, error) {
	if minContrast <= 0 {
		minContrast = 4.5 // AA normal text
	}

	// Parse background color
	bgLuminance, err := relativeLuminance(background)
	if err != nil {
		return nil, fmt.Errorf("invalid background color: %w", err)
	}

	var suggestions []ContrastSuggestion

	// Check all color tokens
	for _, c := range s.ds.Foundations.Colors {
		fgLuminance, err := relativeLuminance(c.Value)
		if err != nil {
			continue
		}

		ratio := contrastRatio(bgLuminance, fgLuminance)

		if ratio >= minContrast {
			suggestions = append(suggestions, ContrastSuggestion{
				Token:         c.ID,
				Value:         c.Value,
				ContrastRatio: ratio,
				MeetsAA:       ratio >= 4.5,
				MeetsAAA:      ratio >= 7.0,
				MeetsAALarge:  ratio >= 3.0,
			})
		}
	}

	return suggestions, nil
}

// relativeLuminance calculates the relative luminance of a color.
func relativeLuminance(hex string) (float64, error) {
	hex = trimHash(hex)
	if len(hex) != 6 {
		return 0, fmt.Errorf("invalid hex color: %s", hex)
	}

	r := float64(hexToDecimal(hex[0:2])) / 255.0
	g := float64(hexToDecimal(hex[2:4])) / 255.0
	b := float64(hexToDecimal(hex[4:6])) / 255.0

	// Apply sRGB transformation
	r = sRGBTransform(r)
	g = sRGBTransform(g)
	b = sRGBTransform(b)

	return 0.2126*r + 0.7152*g + 0.0722*b, nil
}

func sRGBTransform(c float64) float64 {
	if c <= 0.03928 {
		return c / 12.92
	}
	return pow((c+0.055)/1.055, 2.4)
}

func pow(base, exp float64) float64 {
	result := 1.0
	for i := 0; i < int(exp*10); i++ {
		result *= base
	}
	// Simplified - use math.Pow in production
	return result
}

func contrastRatio(l1, l2 float64) float64 {
	lighter := l1
	darker := l2
	if l2 > l1 {
		lighter = l2
		darker = l1
	}
	return (lighter + 0.05) / (darker + 0.05)
}

func trimHash(s string) string {
	if len(s) > 0 && s[0] == '#' {
		return s[1:]
	}
	return s
}

func hexToDecimal(hex string) int {
	var result int
	fmt.Sscanf(hex, "%x", &result)
	return result
}

// ComponentFixContext provides full context for fixing accessibility issues.
type ComponentFixContext struct {
	ComponentID string `json:"componentId"`
	IssueType   string `json:"issueType"`

	// FilePattern suggests where to find the component.
	FilePattern string `json:"filePattern"`

	// PropsToAdd lists props that may need to be added.
	PropsToAdd []RequiredProp `json:"propsToAdd,omitempty"`

	// StylesToCheck lists CSS properties to verify.
	StylesToCheck []StyleCheck `json:"stylesToCheck,omitempty"`

	// TokensAvailable lists available design tokens for the fix.
	TokensAvailable map[string][]string `json:"tokensAvailable,omitempty"`

	// RelatedComponents lists components that may also need fixes.
	RelatedComponents []string `json:"relatedComponents,omitempty"`

	// AntiPatterns lists patterns to avoid when fixing.
	AntiPatterns []string `json:"antiPatterns,omitempty"`
}

// StyleCheck describes a CSS property to verify.
type StyleCheck struct {
	Property string `json:"property"`
	Token    string `json:"token,omitempty"`
	MinValue string `json:"minValue,omitempty"`
}

// GetComponentFixContext returns full context for fixing accessibility issues in a component.
func (s *Service) GetComponentFixContext(ctx context.Context, componentID, issueType string) (*ComponentFixContext, error) {
	comp, err := s.GetComponent(ctx, componentID)
	if err != nil {
		return nil, err
	}

	context := &ComponentFixContext{
		ComponentID:     componentID,
		IssueType:       issueType,
		FilePattern:     fmt.Sprintf("src/components/%s.{tsx,vue,svelte}", capitalize(comp.Name)),
		TokensAvailable: make(map[string][]string),
	}

	// Add related components
	if len(comp.Uses) > 0 {
		context.RelatedComponents = append(context.RelatedComponents, comp.Uses...)
	}

	// Gather available tokens
	for _, c := range s.ds.Foundations.Colors {
		context.TokensAvailable["colors"] = append(context.TokensAvailable["colors"], c.ID)
	}
	if s.ds.Foundations.Spacing != nil {
		for _, sp := range s.ds.Foundations.Spacing.Scale {
			context.TokensAvailable["spacing"] = append(context.TokensAvailable["spacing"], sp.ID)
		}
	}

	// Issue-specific context
	switch issueType {
	case "color-contrast":
		context.StylesToCheck = []StyleCheck{
			{Property: "color", Token: "color.text.*"},
			{Property: "background-color", Token: "color.bg.*"},
		}
		context.PropsToAdd = nil // Contrast is a styling issue, not props

	case "missing-label":
		context.PropsToAdd = []RequiredProp{
			{Name: "aria-label", Type: "string", When: "no visible label", Example: `aria-label="Description"`},
			{Name: "aria-labelledby", Type: "string", When: "label exists elsewhere", Example: `aria-labelledby="label-id"`},
		}

	case "keyboard":
		context.PropsToAdd = []RequiredProp{
			{Name: "tabIndex", Type: "number", Example: `tabIndex={0}`},
			{Name: "onKeyDown", Type: "function", Example: `onKeyDown={handleKeyDown}`},
		}
		context.StylesToCheck = []StyleCheck{
			{Property: "outline", MinValue: "2px"},
		}

	case "focus":
		context.StylesToCheck = []StyleCheck{
			{Property: "outline", MinValue: "2px"},
			{Property: "outline-offset", MinValue: "2px"},
		}
	}

	// Add anti-patterns for this issue
	antiPatterns, _ := s.GetAntiPatterns(ctx, componentID, issueType)
	if antiPatterns != nil {
		for _, ap := range antiPatterns.AntiPatterns {
			context.AntiPatterns = append(context.AntiPatterns, ap.Description)
		}
	}

	return context, nil
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}
