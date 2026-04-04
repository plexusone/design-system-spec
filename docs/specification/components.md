# Components

Components are the UI building blocks of your design system. Each component spec defines its variants, states, props, and LLM context.

## File Location

```
my-design-system/
└── components/
    ├── button.json
    ├── card.json
    └── input.json
```

## Schema

```json
{
  "id": "Button",
  "name": "Button",
  "description": "Interactive button for triggering actions",
  "category": "actions",
  "variants": [...],
  "states": [...],
  "props": [...],
  "slots": [...],
  "anatomy": [...],
  "constraints": {...},
  "llm": {...}
}
```

## Fields

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique identifier (used in code generation) |
| `name` | string | Human-readable name |
| `description` | string | What this component does |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `category` | string | Component category (actions, inputs, layout, etc.) |
| `variants` | array | Visual/behavioral variations |
| `states` | array | Interactive states (hover, focus, disabled) |
| `props` | array | Configurable properties |
| `slots` | array | Content insertion points |
| `anatomy` | array | Internal structure parts |
| `constraints` | object | Size and layout constraints |
| `llm` | object | LLM optimization context |

## Variants

Variants define the visual or behavioral variations of a component.

```json
{
  "variants": [
    {
      "id": "default",
      "name": "Default",
      "description": "Primary action button",
      "isDefault": true
    },
    {
      "id": "secondary",
      "name": "Secondary",
      "description": "Less prominent actions"
    },
    {
      "id": "destructive",
      "name": "Destructive",
      "description": "Dangerous actions like delete"
    },
    {
      "id": "ghost",
      "name": "Ghost",
      "description": "Minimal button for inline actions"
    }
  ]
}
```

## States

States define interactive states beyond variants.

```json
{
  "states": [
    {
      "id": "default",
      "name": "Default",
      "description": "Normal resting state"
    },
    {
      "id": "hover",
      "name": "Hover",
      "description": "Mouse over state"
    },
    {
      "id": "focus",
      "name": "Focus",
      "description": "Keyboard focus state"
    },
    {
      "id": "disabled",
      "name": "Disabled",
      "description": "Non-interactive state"
    }
  ]
}
```

## Props

Props define the configurable properties of a component.

```json
{
  "props": [
    {
      "name": "variant",
      "type": "enum",
      "description": "Visual style variant",
      "required": false
    },
    {
      "name": "size",
      "type": "enum",
      "description": "Button size",
      "defaultValue": "default",
      "required": false
    },
    {
      "name": "disabled",
      "type": "boolean",
      "description": "Whether the button is disabled",
      "defaultValue": false
    },
    {
      "name": "asChild",
      "type": "boolean",
      "description": "Render as child element for composition"
    }
  ]
}
```

### Prop Types

| Type | Description |
|------|-------------|
| `string` | Text value |
| `number` | Numeric value |
| `boolean` | True/false |
| `enum` | One of predefined values (from variants or sizes) |
| `node` | React node content |
| `function` | Callback function |

## Slots

Slots define content insertion points.

```json
{
  "slots": [
    {
      "name": "icon",
      "description": "Icon displayed before text",
      "required": false
    },
    {
      "name": "children",
      "description": "Button label content",
      "required": true
    }
  ]
}
```

## Constraints

Constraints define size and layout limitations.

```json
{
  "constraints": {
    "minWidth": "64px",
    "maxWidth": "100%",
    "minHeight": "32px"
  }
}
```

## LLM Context

The `llm` field optimizes the component for AI code generation.

```json
{
  "llm": {
    "intent": "Trigger user-initiated actions like form submission",
    "allowedContexts": [
      "forms",
      "dialogs",
      "toolbars",
      "cards",
      "page-actions"
    ],
    "forbiddenContexts": [
      "inline-text",
      "navigation-menu-items"
    ],
    "exampleUsage": [
      "<Button>Save Changes</Button>",
      "<Button variant=\"secondary\">Cancel</Button>",
      "<Button variant=\"destructive\">Delete</Button>"
    ],
    "antiPatterns": [
      "Multiple primary buttons in the same view",
      "Using button for navigation (use Link instead)",
      "Ghost buttons for important primary actions",
      "Disabled buttons without explanation tooltip"
    ],
    "semanticMeaning": "Signals an actionable element that triggers behavior",
    "relatedElements": ["link", "icon-button"],
    "priorityScore": 90
  }
}
```

### LLM Fields

| Field | Type | Description |
|-------|------|-------------|
| `intent` | string | What this component is for |
| `allowedContexts` | array | Where this component should be used |
| `forbiddenContexts` | array | Where this component should NOT be used |
| `exampleUsage` | array | Code examples for the LLM |
| `antiPatterns` | array | Common mistakes to avoid |
| `semanticMeaning` | string | What this component communicates |
| `relatedElements` | array | Similar or complementary components |
| `priorityScore` | number | Relative importance (0-100) |

## Complete Example

```json
{
  "id": "Button",
  "name": "Button",
  "description": "Interactive button for triggering actions. Core component for all user interactions.",
  "category": "actions",
  "variants": [
    {
      "id": "default",
      "name": "Default",
      "description": "Primary action button with solid background",
      "isDefault": true
    },
    {
      "id": "secondary",
      "name": "Secondary",
      "description": "Less prominent actions, secondary CTAs"
    },
    {
      "id": "outline",
      "name": "Outline",
      "description": "Bordered button for tertiary actions"
    },
    {
      "id": "ghost",
      "name": "Ghost",
      "description": "Minimal button for inline actions"
    },
    {
      "id": "destructive",
      "name": "Destructive",
      "description": "Dangerous actions like delete"
    }
  ],
  "states": [
    { "id": "default", "name": "Default" },
    { "id": "hover", "name": "Hover" },
    { "id": "focus", "name": "Focus" },
    { "id": "active", "name": "Active" },
    { "id": "disabled", "name": "Disabled" }
  ],
  "props": [
    {
      "name": "variant",
      "type": "enum",
      "description": "Visual style variant"
    },
    {
      "name": "size",
      "type": "enum",
      "description": "Button size"
    },
    {
      "name": "disabled",
      "type": "boolean",
      "description": "Whether the button is disabled"
    }
  ],
  "accessibility": {
    "role": "button",
    "keyboardInteractions": [
      { "key": "Enter", "action": "Activate button" },
      { "key": "Space", "action": "Activate button" }
    ]
  },
  "llm": {
    "intent": "Trigger user-initiated actions",
    "allowedContexts": ["forms", "dialogs", "toolbars"],
    "forbiddenContexts": ["inline-text"],
    "antiPatterns": [
      "Multiple primary buttons in the same view",
      "Using button for navigation"
    ],
    "exampleUsage": [
      "<Button>Save</Button>",
      "<Button variant=\"destructive\">Delete</Button>"
    ]
  }
}
```
