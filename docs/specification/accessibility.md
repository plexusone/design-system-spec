# Accessibility

The Accessibility layer defines WCAG compliance requirements and guidelines.

## File Location

```
my-design-system/
└── accessibility.json
```

## Schema

```json
{
  "wcagLevel": "AA",
  "colorContrast": {...},
  "keyboard": {...},
  "guidelines": [...]
}
```

## Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `wcagLevel` | string | Yes | Target WCAG level: "A", "AA", or "AAA" |
| `colorContrast` | object | No | Contrast requirements |
| `keyboard` | object | No | Keyboard navigation requirements |
| `guidelines` | array | No | Additional guidelines |

## Color Contrast

```json
{
  "colorContrast": {
    "normalTextRatio": 4.5,
    "largeTextRatio": 3.0,
    "uiComponentRatio": 3.0
  }
}
```

## Keyboard Requirements

```json
{
  "keyboard": {
    "focusVisible": true,
    "skipLinks": true,
    "noKeyboardTraps": true
  }
}
```

## Guidelines

```json
{
  "guidelines": [
    {
      "id": "focus-visibility",
      "name": "Focus Visibility",
      "description": "All interactive elements must have a visible focus indicator that meets contrast requirements"
    },
    {
      "id": "alt-text",
      "name": "Alternative Text",
      "description": "All meaningful images must have descriptive alt text"
    },
    {
      "id": "color-independence",
      "name": "Color Independence",
      "description": "Information must not be conveyed by color alone"
    }
  ]
}
```

## Complete Example

```json
{
  "wcagLevel": "AA",
  "colorContrast": {
    "normalTextRatio": 4.5,
    "largeTextRatio": 3.0,
    "uiComponentRatio": 3.0
  },
  "keyboard": {
    "focusVisible": true,
    "skipLinks": true,
    "noKeyboardTraps": true
  },
  "guidelines": [
    {
      "id": "focus-visibility",
      "name": "Focus Visibility",
      "description": "All interactive elements must have a visible focus indicator"
    },
    {
      "id": "alt-text",
      "name": "Alternative Text",
      "description": "All meaningful images must have descriptive alt text"
    },
    {
      "id": "color-independence",
      "name": "Color Independence",
      "description": "Information must not be conveyed by color alone"
    }
  ]
}
```

## Validation

The `dss validate` command checks for common accessibility issues:

- Images missing `alt` attributes
- Icon-only buttons without `aria-label`
- Color-only information indicators

See [CLI Reference](../cli.md) for validation details.
