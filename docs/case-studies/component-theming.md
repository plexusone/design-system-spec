# Case Study: Component Theming Contracts

How a component library (markdown-editor) and an application (PlexusOne) use DSS to coordinate theming through a formal contract.

## The Challenge

When integrating third-party components, theming becomes complex:

| Challenge | Without DSS |
|-----------|-------------|
| **Discovery** | Read component source to find CSS variables |
| **Mapping** | Manually map app tokens to component tokens |
| **Maintenance** | No validation when either spec changes |
| **Documentation** | Scattered across READMEs and code comments |

### The Scenario

- **markdown-editor**: A Lit Web Component with its own design tokens
- **PlexusOne**: An application with its own design system
- **Goal**: PlexusOne wants to theme markdown-editor using its colors

## The Solution: Theming Contracts

DSS enables a formal contract between component and consumer:

```
┌─────────────────────────────────────────────────────────────────────┐
│                    COMPONENT (markdown-editor)                       │
│                                                                     │
│   design-system.json                                                │
│   ├── foundations/colors        ← Internal token definitions        │
│   └── themingContract           ← PUBLIC: What consumers can theme  │
│         prefix: "--mde-theme"                                       │
│         tokens: [bg-primary, accent, ...]                           │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  │ Contract
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    APPLICATION (PlexusOne)                           │
│                                                                     │
│   design-system.json                                                │
│   ├── foundations/colors        ← App token definitions             │
│   └── themeBindings             ← Maps app tokens to components     │
│         @grokify/markdown-editor:                                   │
│           plexus-dark → bg-primary                                  │
│           plexus-cyan → accent                                      │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  │ dss generate
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    GENERATED CSS                                     │
│                                                                     │
│   markdown-editor {                                                 │
│     --mde-theme-bg-primary: var(--plexus-dark);                     │
│     --mde-theme-accent: var(--plexus-cyan);                         │
│   }                                                                 │
└─────────────────────────────────────────────────────────────────────┘
```

## Schema: Theming Contract (Component Side)

Components declare what can be themed via `themingContract`:

```json
{
  "meta": {
    "name": "Markdown Editor",
    "version": "1.0.0"
  },
  "themingContract": {
    "prefix": "--mde-theme",
    "description": "CSS custom properties for external theming",
    "tokens": [
      {
        "id": "bg-primary",
        "cssProperty": "--mde-theme-bg-primary",
        "semantic": "background",
        "description": "Main background color",
        "defaultLight": "#ffffff",
        "defaultDark": "#0f172a"
      },
      {
        "id": "bg-secondary",
        "cssProperty": "--mde-theme-bg-secondary",
        "semantic": "background",
        "description": "Toolbar and secondary surfaces",
        "defaultLight": "#f8fafc",
        "defaultDark": "#1e293b"
      },
      {
        "id": "accent",
        "cssProperty": "--mde-theme-accent",
        "semantic": "primary",
        "description": "Primary accent for interactive elements",
        "defaultLight": "#3b82f6",
        "defaultDark": "#60a5fa"
      },
      {
        "id": "text-primary",
        "cssProperty": "--mde-theme-text-primary",
        "semantic": "foreground",
        "description": "Primary text color",
        "defaultLight": "#1e293b",
        "defaultDark": "#f1f5f9"
      }
    ]
  }
}
```

### Token Semantics

The `semantic` field enables automatic mapping:

| Semantic | Meaning | Examples |
|----------|---------|----------|
| `background` | Surface colors | `bg-primary`, `bg-secondary` |
| `foreground` | Text colors | `text-primary`, `text-muted` |
| `primary` | Brand/accent colors | `accent`, `link` |
| `border` | Edge colors | `border`, `divider` |
| `success` | Positive states | `success`, `valid` |
| `danger` | Negative states | `error`, `destructive` |

## Schema: Theme Bindings (Application Side)

Applications declare how to map their tokens to components:

```json
{
  "meta": {
    "name": "PlexusOne",
    "version": "1.0.0"
  },
  "foundations": {
    "colors": [
      { "id": "plexus-dark", "value": "#0a0e1a", "semantic": "background" },
      { "id": "plexus-slate", "value": "#1e293b", "semantic": "background" },
      { "id": "plexus-cyan", "value": "#06b6d4", "semantic": "primary" },
      { "id": "foreground", "value": "#f1f5f9", "semantic": "foreground" }
    ]
  },
  "themeBindings": [
    {
      "component": "@grokify/markdown-editor",
      "specUrl": "https://github.com/grokify/markdown-editor/blob/main/design-system.json",
      "mappings": [
        { "from": "plexus-dark", "to": "bg-primary" },
        { "from": "plexus-slate", "to": "bg-secondary" },
        { "from": "plexus-cyan", "to": "accent" },
        { "from": "foreground", "to": "text-primary" }
      ],
      "themeMode": "dark"
    }
  ]
}
```

### Mapping Strategies

| Strategy | Description | When to Use |
|----------|-------------|-------------|
| **Explicit** | Manual `from`/`to` mapping | Full control, custom branding |
| **Semantic** | Auto-match by `semantic` field | Quick integration, standard colors |
| **Hybrid** | Semantic defaults + explicit overrides | Balance of automation and control |

## Automatic Semantic Mapping

When semantics align, mappings can be generated automatically:

```bash
# Generate bindings by matching semantic fields
dss bind @grokify/markdown-editor --strategy=semantic
```

The tool:

1. Fetches component's `themingContract`
2. Matches app tokens by `semantic` field
3. Generates `themeBindings` section
4. Warns about unmatched tokens

```
✓ Matched: bg-primary ← plexus-dark (semantic: background)
✓ Matched: accent ← plexus-cyan (semantic: primary)
⚠ Unmatched: bg-tertiary (no app token with semantic: background)
  Suggestion: Add mapping or use plexus-slate
```

## Generated Output

### CSS Output

```css
/* Generated by: dss generate --css --bindings */

/* markdown-editor theme bindings */
markdown-editor {
  /* Dark theme (default for PlexusOne) */
  --mde-theme-bg-primary: var(--plexus-dark);
  --mde-theme-bg-secondary: var(--plexus-slate);
  --mde-theme-accent: var(--plexus-cyan);
  --mde-theme-text-primary: var(--plexus-foreground);
}

/* Light theme override */
[data-theme="light"] markdown-editor {
  --mde-theme-bg-primary: #ffffff;
  --mde-theme-bg-secondary: #f8fafc;
  --mde-theme-accent: #3b82f6;
  --mde-theme-text-primary: #1e293b;
}
```

### TypeScript Output

```typescript
// Generated by: dss generate --typescript --bindings

export const markdownEditorTheme = {
  dark: {
    '--mde-theme-bg-primary': 'var(--plexus-dark)',
    '--mde-theme-bg-secondary': 'var(--plexus-slate)',
    '--mde-theme-accent': 'var(--plexus-cyan)',
    '--mde-theme-text-primary': 'var(--plexus-foreground)',
  },
  light: {
    '--mde-theme-bg-primary': '#ffffff',
    '--mde-theme-bg-secondary': '#f8fafc',
    '--mde-theme-accent': '#3b82f6',
    '--mde-theme-text-primary': '#1e293b',
  },
} as const;
```

## Implementation: markdown-editor

### 1. Add Theming Contract to DSS

```json
// design-system.json
{
  "meta": { "name": "Markdown Editor", "version": "1.0.0" },
  "themingContract": {
    "prefix": "--mde-theme",
    "tokens": [...]
  },
  "foundations": {
    "colors": [...]
  }
}
```

### 2. Component Uses Contract Tokens

```typescript
// shared.styles.ts
export const colors = css`
  :host {
    /* Use theming API with fallbacks to internal defaults */
    --mde-bg-primary: var(--mde-theme-bg-primary, #ffffff);
    --mde-accent: var(--mde-theme-accent, #3b82f6);
  }

  :host([theme='dark']) {
    --mde-bg-primary: var(--mde-theme-bg-primary, #0f172a);
    --mde-accent: var(--mde-theme-accent, #60a5fa);
  }
`;
```

### 3. Document the Contract

The theming contract becomes part of the component's public API:

- Published in `design-system.json`
- Documented in MkDocs (theming.md)
- Validated by CI (schema check)

## Implementation: PlexusOne

### 1. Add Theme Bindings to DSS

```json
// design-system.json
{
  "meta": { "name": "PlexusOne", "version": "1.0.0" },
  "themeBindings": [
    {
      "component": "@grokify/markdown-editor",
      "mappings": [...]
    }
  ]
}
```

### 2. Generate and Include CSS

```bash
# Generate theme bindings CSS
dss generate --css output/component-themes.css --bindings
```

```html
<!-- Include in app -->
<link rel="stylesheet" href="/css/component-themes.css">
```

### 3. Or Inline in Component Integration

```html
<!-- tools/markdown-editor/index.html -->
<style>
  /* PlexusOne theme → markdown-editor */
  markdown-editor {
    --mde-theme-bg-primary: var(--plexus-dark);
    --mde-theme-accent: var(--plexus-cyan);
  }
</style>
```

## Validation

### Contract Validation (Component CI)

```yaml
# .github/workflows/dss.yml
- name: Validate Theming Contract
  run: |
    dss validate design-system.json
    dss contract check --ensure-defaults
```

Validates:

- All contract tokens have `defaultLight` and `defaultDark`
- CSS properties follow naming convention
- Semantics are from allowed set

### Binding Validation (App CI)

```yaml
# .github/workflows/dss.yml
- name: Validate Theme Bindings
  run: |
    dss validate design-system.json
    dss bindings check --warn-unbound
```

Validates:

- Referenced components exist
- All `to` tokens exist in component contract
- All `from` tokens exist in app foundations
- Warns about unmapped contract tokens

## Benefits

| Benefit | Description |
|---------|-------------|
| **Discoverability** | `themingContract` is the single source of truth for what can be themed |
| **Type Safety** | Schema validation catches mismatches at build time |
| **Auto-Generation** | CSS mappings generated from specs, not hand-written |
| **Semantic Matching** | Tokens with same semantics can auto-bind |
| **Documentation** | Contract is self-documenting, rendered in MkDocs |
| **LLM Context** | AI assistants understand the contract for code generation |

## Real-World Results

### Before DSS Contracts

```css
/* Hand-written, discovered by reading component source */
markdown-editor {
  --mde-bg-primary: #0a0e1a;  /* Copied from PlexusOne */
  --mde-accent: #06b6d4;       /* Hope this is right... */
  /* What else can I theme? */
}
```

### After DSS Contracts

```bash
# Discover themeable tokens
$ dss contract show @grokify/markdown-editor

Theming Contract: @grokify/markdown-editor
Prefix: --mde-theme

Tokens:
  bg-primary     (background)  Main background color
  bg-secondary   (background)  Toolbar and secondary surfaces
  accent         (primary)     Primary accent for interactive elements
  text-primary   (foreground)  Primary text color
  ...

# Generate bindings
$ dss bind @grokify/markdown-editor --output themes/markdown-editor.css
✓ Generated 12 token mappings
```

## Future: Theming Contract Registry

A central registry could enable:

```bash
# Search for themeable components
$ dss search --has-contract

@grokify/markdown-editor  v1.0.0  Markdown editor with PDF export
@shadcn/ui                v2.0.0  Beautifully designed components
@radix-ui/themes          v3.0.0  Radix Themes

# Quick integration
$ dss bind @shadcn/ui --strategy=semantic --output themes/
```

## Conclusion

DSS theming contracts formalize the relationship between component libraries and consuming applications:

1. **Components publish** what can be themed (`themingContract`)
2. **Apps declare** how to map tokens (`themeBindings`)
3. **Tools generate** CSS from specs (`dss generate --bindings`)
4. **CI validates** contracts don't break (`dss bindings check`)

This eliminates guesswork, enables automation, and ensures theming stays in sync as both specs evolve.
