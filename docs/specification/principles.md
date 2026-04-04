# Principles

Principles define the design philosophy and guidelines that inform all design decisions.

## File Location

```
my-design-system/
└── principles.json
```

## Schema

```json
{
  "principles": [
    {
      "name": "Clarity Over Complexity",
      "description": "Keep interfaces simple and understandable",
      "examples": [...]
    }
  ]
}
```

## Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Principle name |
| `description` | string | Yes | What the principle means |
| `examples` | array | No | Examples of the principle in action |

## Example

```json
{
  "principles": [
    {
      "name": "Clarity Over Complexity",
      "description": "Keep interfaces simple and understandable. Avoid jargon, use clear labels, and present information hierarchically.",
      "examples": [
        {
          "good": "Use 'Save' instead of 'Persist Changes'",
          "bad": "Using technical jargon in user-facing text"
        }
      ]
    },
    {
      "name": "Consistency",
      "description": "Similar elements should look and behave similarly across the system.",
      "examples": [
        {
          "good": "All primary buttons use the same color",
          "bad": "Different button styles for the same action type"
        }
      ]
    },
    {
      "name": "Accessibility First",
      "description": "Design for all users, including those with disabilities.",
      "examples": [
        {
          "good": "All interactive elements are keyboard accessible",
          "bad": "Relying only on hover states for important information"
        }
      ]
    }
  ]
}
```

## Usage in LLM Context

Principles are included in the generated LLM prompt to help AI assistants understand and apply your design philosophy when generating code.
