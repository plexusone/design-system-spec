# Meta

The Meta layer defines system-level metadata for your design system.

## File Location

```
my-design-system/
└── meta.json
```

## Schema

```json
{
  "name": "My Design System",
  "version": "1.0.0",
  "description": "Description of the design system",
  "maintainers": [...]
}
```

## Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Name of the design system |
| `version` | string | Yes | Semantic version (e.g., "1.0.0") |
| `description` | string | No | Brief description |
| `maintainers` | array | No | List of maintainers |

## Maintainers

```json
{
  "maintainers": [
    {
      "name": "Design Team",
      "email": "design@example.com",
      "role": "owner"
    }
  ]
}
```

## Example

```json
{
  "name": "Acme Design System",
  "version": "2.1.0",
  "description": "The official design system for Acme products",
  "maintainers": [
    {
      "name": "Design Systems Team",
      "email": "ds-team@acme.com",
      "role": "owner"
    }
  ]
}
```
