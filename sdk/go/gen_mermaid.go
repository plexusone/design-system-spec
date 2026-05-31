package dss

import (
	"fmt"
	"strings"
)

// MermaidOptions configures Mermaid diagram output.
type MermaidOptions struct {
	// DiagramType: "graph", "flowchart", "classDiagram"
	DiagramType string

	// Direction: "TD" (top-down), "LR" (left-right), "BT", "RL"
	Direction string

	// ShowTokens includes token relationships
	ShowTokens bool

	// ShowEvents includes event flow arrows
	ShowEvents bool

	// ShowSlots includes slot composition
	ShowSlots bool

	// GroupByCategory groups components by their category
	GroupByCategory bool

	// Theme: "default", "dark", "forest", "neutral"
	Theme string
}

// DefaultMermaidOptions returns sensible defaults.
func DefaultMermaidOptions() MermaidOptions {
	return MermaidOptions{
		DiagramType:     "graph",
		Direction:       "TD",
		ShowTokens:      false,
		ShowEvents:      true,
		ShowSlots:       true,
		GroupByCategory: true,
		Theme:           "default",
	}
}

// GenerateMermaid generates a Mermaid diagram showing component relationships.
func (ds *DesignSystem) GenerateMermaid(opts MermaidOptions) (string, error) {
	if len(ds.Components) == 0 {
		return "", fmt.Errorf("design system has no components defined")
	}

	var b strings.Builder

	// Theme configuration
	if opts.Theme != "" && opts.Theme != "default" {
		b.WriteString("%%{init: {'theme': '")
		b.WriteString(opts.Theme)
		b.WriteString("'}}%%\n")
	}

	// Diagram declaration
	b.WriteString(opts.DiagramType)
	b.WriteString(" ")
	b.WriteString(opts.Direction)
	b.WriteString("\n")

	// Build component ID to node mapping
	nodeIDs := make(map[string]string)
	for _, c := range ds.Components {
		nodeIDs[c.ID] = sanitizeMermaidID(c.ID)
	}

	// Group by category if enabled
	if opts.GroupByCategory {
		categories := make(map[string][]Component)
		for _, c := range ds.Components {
			cat := c.Category
			if cat == "" {
				cat = "other"
			}
			categories[cat] = append(categories[cat], c)
		}

		for cat, components := range categories {
			b.WriteString(fmt.Sprintf("    subgraph %s[%s]\n", sanitizeMermaidID(cat), formatCategoryName(cat)))
			for _, c := range components {
				writeComponentNode(&b, c, nodeIDs)
			}
			b.WriteString("    end\n")
		}
	} else {
		// Write all components without grouping
		for _, c := range ds.Components {
			writeComponentNode(&b, c, nodeIDs)
		}
	}

	b.WriteString("\n")

	// Write relationships
	for _, c := range ds.Components {
		nodeID := nodeIDs[c.ID]

		// Uses relationships (explicit dependencies)
		for _, usedID := range c.Uses {
			if targetNode, ok := nodeIDs[usedID]; ok {
				b.WriteString(fmt.Sprintf("    %s --> %s\n", nodeID, targetNode))
			}
		}

		// Slot relationships (composition)
		if opts.ShowSlots {
			for _, slot := range c.Slots {
				for _, allowedID := range slot.AllowedComponents {
					if targetNode, ok := nodeIDs[allowedID]; ok {
						b.WriteString(fmt.Sprintf("    %s -.->|%s| %s\n", nodeID, slot.Name, targetNode))
					}
				}
			}
		}

		// Event relationships (if ShowEvents and we can infer listeners)
		if opts.ShowEvents && len(c.Events) > 0 {
			// Add event annotations to the node
			for _, evt := range c.Events {
				b.WriteString(fmt.Sprintf("    %s -.->|emits: %s| %s\n", nodeID, evt.Name, nodeID))
			}
		}

		// Constraint relationships
		if c.Constraints != nil {
			if c.Constraints.RequiredParent != "" {
				if parentNode, ok := nodeIDs[c.Constraints.RequiredParent]; ok {
					b.WriteString(fmt.Sprintf("    %s -->|parent| %s\n", parentNode, nodeID))
				}
			}
		}
	}

	// Write pattern relationships
	if len(ds.Patterns) > 0 {
		b.WriteString("\n    %% Patterns\n")
		for _, p := range ds.Patterns {
			patternID := sanitizeMermaidID("pattern_" + p.ID)
			b.WriteString(fmt.Sprintf("    %s{{%s}}\n", patternID, p.Name))

			for _, pc := range p.Components {
				if targetNode, ok := nodeIDs[pc.ComponentID]; ok {
					label := pc.Role
					if label == "" {
						label = "uses"
					}
					b.WriteString(fmt.Sprintf("    %s -.->|%s| %s\n", patternID, label, targetNode))
				}
			}
		}
	}

	return b.String(), nil
}

// GenerateMermaidComponentDiagram generates a focused diagram for a single component.
func (ds *DesignSystem) GenerateMermaidComponentDiagram(componentID string, opts MermaidOptions) (string, error) {
	// Find the component
	var target *Component
	for i := range ds.Components {
		if ds.Components[i].ID == componentID {
			target = &ds.Components[i]
			break
		}
	}
	if target == nil {
		return "", fmt.Errorf("component %q not found", componentID)
	}

	var b strings.Builder

	b.WriteString(opts.DiagramType)
	b.WriteString(" ")
	b.WriteString(opts.Direction)
	b.WriteString("\n")

	nodeID := sanitizeMermaidID(target.ID)

	// Central component
	b.WriteString(fmt.Sprintf("    %s[<b>%s</b>]\n", nodeID, target.Name))

	// Dependencies (uses)
	if len(target.Uses) > 0 {
		b.WriteString("\n    %% Dependencies\n")
		for _, usedID := range target.Uses {
			usedNode := sanitizeMermaidID(usedID)
			b.WriteString(fmt.Sprintf("    %s[%s]\n", usedNode, usedID))
			b.WriteString(fmt.Sprintf("    %s --> %s\n", nodeID, usedNode))
		}
	}

	// Slots (what can be placed inside)
	if opts.ShowSlots && len(target.Slots) > 0 {
		b.WriteString("\n    %% Slots\n")
		for _, slot := range target.Slots {
			for _, allowedID := range slot.AllowedComponents {
				allowedNode := sanitizeMermaidID(allowedID)
				b.WriteString(fmt.Sprintf("    %s[%s]\n", allowedNode, allowedID))
				b.WriteString(fmt.Sprintf("    %s -.->|slot: %s| %s\n", nodeID, slot.Name, allowedNode))
			}
		}
	}

	// Events
	if opts.ShowEvents && len(target.Events) > 0 {
		b.WriteString("\n    %% Events\n")
		for _, evt := range target.Events {
			evtNode := sanitizeMermaidID("evt_" + evt.Name)
			b.WriteString(fmt.Sprintf("    %s((%s))\n", evtNode, evt.Name))
			b.WriteString(fmt.Sprintf("    %s -.->|emits| %s\n", nodeID, evtNode))
		}
	}

	// Props as a note
	if len(target.Props) > 0 {
		b.WriteString("\n    %% Props\n")
		propsNode := sanitizeMermaidID(target.ID + "_props")
		var propList []string
		for _, p := range target.Props {
			propList = append(propList, fmt.Sprintf("%s: %s", p.Name, p.Type))
		}
		// Mermaid note syntax
		b.WriteString(fmt.Sprintf("    %s[\"`Props:\n", propsNode))
		for _, prop := range propList[:min(5, len(propList))] { // Limit to 5 props
			b.WriteString(fmt.Sprintf("    - %s\n", prop))
		}
		if len(propList) > 5 {
			b.WriteString(fmt.Sprintf("    - ... +%d more\n", len(propList)-5))
		}
		b.WriteString("`\"]\n")
		b.WriteString(fmt.Sprintf("    %s --- %s\n", nodeID, propsNode))
	}

	return b.String(), nil
}

// GenerateMermaidTokenDiagram generates a diagram showing token usage by components.
func (ds *DesignSystem) GenerateMermaidTokenDiagram(opts MermaidOptions) (string, error) {
	var b strings.Builder

	b.WriteString("graph LR\n")

	// Collect all tokens used
	tokenUsers := make(map[string][]string)
	for _, c := range ds.Components {
		for _, tokenID := range c.TokensUsed {
			tokenUsers[tokenID] = append(tokenUsers[tokenID], c.ID)
		}
	}

	if len(tokenUsers) == 0 {
		return "", fmt.Errorf("no token usage found in components")
	}

	// Group tokens by category (infer from naming)
	b.WriteString("    subgraph Tokens\n")
	for tokenID := range tokenUsers {
		tokenNode := sanitizeMermaidID("token_" + tokenID)
		b.WriteString(fmt.Sprintf("        %s[/%s/]\n", tokenNode, tokenID))
	}
	b.WriteString("    end\n\n")

	b.WriteString("    subgraph Components\n")
	componentNodes := make(map[string]bool)
	for _, users := range tokenUsers {
		for _, compID := range users {
			if !componentNodes[compID] {
				componentNodes[compID] = true
				compNode := sanitizeMermaidID(compID)
				b.WriteString(fmt.Sprintf("        %s[%s]\n", compNode, compID))
			}
		}
	}
	b.WriteString("    end\n\n")

	// Draw relationships
	for tokenID, users := range tokenUsers {
		tokenNode := sanitizeMermaidID("token_" + tokenID)
		for _, compID := range users {
			compNode := sanitizeMermaidID(compID)
			b.WriteString(fmt.Sprintf("    %s -.-> %s\n", tokenNode, compNode))
		}
	}

	return b.String(), nil
}

// Helper functions

func writeComponentNode(b *strings.Builder, c Component, nodeIDs map[string]string) {
	nodeID := nodeIDs[c.ID]
	// Use different shapes based on component type/category
	shape := getComponentShape(c)
	b.WriteString(fmt.Sprintf("        %s%s%s%s\n", nodeID, shape.open, c.Name, shape.close))
}

type nodeShape struct {
	open  string
	close string
}

func getComponentShape(c Component) nodeShape {
	// Different shapes for different component types
	switch c.Category {
	case "inputs", "input", "form":
		return nodeShape{"[/", "/]"} // Parallelogram
	case "layout", "container":
		return nodeShape{"[[", "]]"} // Subroutine shape
	case "feedback", "notification":
		return nodeShape{"((", "))"} // Circle
	case "navigation", "nav":
		return nodeShape{">", "]"} // Asymmetric
	default:
		return nodeShape{"[", "]"} // Rectangle
	}
}

func sanitizeMermaidID(id string) string {
	// Replace characters that are invalid in Mermaid IDs
	replacer := strings.NewReplacer(
		"-", "_",
		" ", "_",
		".", "_",
		"/", "_",
	)
	return replacer.Replace(id)
}

func formatCategoryName(cat string) string {
	// Convert category ID to display name
	return strings.Title(strings.ReplaceAll(cat, "-", " "))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
