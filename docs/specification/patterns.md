# Patterns

Patterns define multi-component solutions for common UX scenarios.

## File Location

```
my-design-system/
└── patterns/
    ├── form.json
    └── data-table.json
```

## Schema

```json
{
  "id": "login-form",
  "name": "Login Form",
  "description": "Standard login form pattern",
  "components": [...],
  "layout": {...},
  "llm": {...}
}
```

## Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique pattern identifier |
| `name` | string | Yes | Human-readable name |
| `description` | string | Yes | What this pattern solves |
| `components` | array | Yes | Components used in this pattern |
| `layout` | object | No | Layout configuration |
| `llm` | object | No | LLM optimization context |

## Example

```json
{
  "id": "login-form",
  "name": "Login Form",
  "description": "Standard login form with email, password, and submit button",
  "components": [
    {
      "componentId": "Input",
      "role": "email-input",
      "props": {
        "type": "email",
        "placeholder": "Email address"
      }
    },
    {
      "componentId": "Input",
      "role": "password-input",
      "props": {
        "type": "password",
        "placeholder": "Password"
      }
    },
    {
      "componentId": "Button",
      "role": "submit",
      "props": {
        "variant": "default"
      }
    }
  ],
  "layout": {
    "type": "stack",
    "direction": "vertical",
    "gap": "4"
  },
  "llm": {
    "intent": "Authenticate users with email and password",
    "allowedContexts": ["auth-page", "modal"],
    "antiPatterns": [
      "Don't auto-submit on input change",
      "Always show password visibility toggle"
    ]
  }
}
```
