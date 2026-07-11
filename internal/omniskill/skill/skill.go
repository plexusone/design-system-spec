// Package skill provides interfaces for defining MCP tools and skills.
// This is a local stub that will be replaced with github.com/plexusone/omniskill/skill
// when that package becomes available.
package skill

import "context"

// Skill represents a collection of related tools that can be exposed via MCP.
type Skill interface {
	// Name returns the skill identifier.
	Name() string

	// Description returns a human-readable description.
	Description() string

	// Init performs any initialization needed by the skill.
	Init(ctx context.Context) error

	// Close cleans up any resources.
	Close() error

	// Tools returns all tools provided by this skill.
	Tools() []Tool
}

// Tool represents a single MCP tool.
type Tool interface {
	// Name returns the tool identifier.
	Name() string

	// Description returns a human-readable description.
	Description() string

	// Parameters returns the parameter schema for this tool.
	Parameters() Parameters

	// Execute runs the tool with the given parameters.
	Execute(ctx context.Context, params map[string]any) (any, error)
}

// Parameters is a map of parameter names to their definitions.
type Parameters map[string]Parameter

// Parameter defines a tool parameter.
type Parameter struct {
	// Type is the JSON Schema type (string, number, boolean, array, object).
	Type string

	// Description explains the parameter's purpose.
	Description string

	// Required indicates whether the parameter must be provided.
	Required bool

	// Default is the default value if not provided.
	Default any

	// Enum lists allowed values for enum types.
	Enum []any
}

// toolImpl implements the Tool interface.
type toolImpl struct {
	name        string
	description string
	parameters  Parameters
	handler     func(ctx context.Context, params map[string]any) (any, error)
}

// NewTool creates a new tool with the given name, description, parameters, and handler.
func NewTool(
	name string,
	description string,
	parameters Parameters,
	handler func(ctx context.Context, params map[string]any) (any, error),
) Tool {
	return &toolImpl{
		name:        name,
		description: description,
		parameters:  parameters,
		handler:     handler,
	}
}

func (t *toolImpl) Name() string {
	return t.name
}

func (t *toolImpl) Description() string {
	return t.description
}

func (t *toolImpl) Parameters() Parameters {
	return t.parameters
}

func (t *toolImpl) Execute(ctx context.Context, params map[string]any) (any, error) {
	return t.handler(ctx, params)
}
