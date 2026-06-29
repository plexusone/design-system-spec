# NPM Package Generation Roadmap

This document specifies the NPM package generation feature for design-system-spec.

## Overview

The `dss generate --package` command generates a complete, publishable NPM package containing design tokens and framework-specific presets from a design system specification.

## Motivation

Currently, design systems like PlexusOne, ProductBuildersHQ, and AIStandardsIO define tokens in JSON specs but must manually create NPM packages for distribution. This feature automates that process, enabling:

1. **Consistent token distribution** across projects
2. **Framework-specific presets** (Tailwind, ShadCN, MkDocs Material)
3. **TypeScript type safety** for token consumers
4. **Version synchronization** between spec and published package

## CLI Interface

### Basic Usage

```bash
# Generate NPM package to ./dist
dss generate --package ./dist

# Generate with specific targets
dss generate --package ./dist --targets tailwind,shadcn,mkdocs-material

# Generate with custom scope
dss generate --package ./dist --scope @myorg
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--package` | `-p` | - | Output directory for NPM package |
| `--targets` | `-t` | `css,tailwind` | Comma-separated list of targets |
| `--scope` | `-s` | From meta | NPM scope (e.g., `@plexusone`) |
| `--name` | `-n` | From meta | Package name (default: `design-tokens`) |
| `--dry-run` | | `false` | Preview without writing files |

### Available Targets

| Target | Description | Output Files |
|--------|-------------|--------------|
| `css` | CSS custom properties | `css/tokens.css` |
| `tailwind` | Tailwind CSS v4 preset | `tailwind/preset.js`, `tailwind/theme.css` |
| `shadcn` | ShadCN/UI theme variables | `shadcn/theme.css`, `shadcn/colors.json` |
| `mkdocs-material` | MkDocs Material theme | `mkdocs/extra.css`, `mkdocs/palette.yml` |
| `scss` | SCSS variables | `scss/_variables.scss` |
| `json` | Raw JSON tokens | `tokens.json` |
| `w3c` | W3C Design Tokens format | `tokens.w3c.json` |

## Output Structure

```
dist/
├── package.json              # Generated from meta
├── README.md                 # Auto-generated documentation
├── index.js                  # CommonJS entry
├── index.mjs                 # ESM entry
├── index.d.ts                # TypeScript declarations
│
├── css/
│   └── tokens.css            # CSS custom properties
│
├── tailwind/
│   ├── preset.js             # Tailwind v4 preset (theme.extend)
│   ├── preset.d.ts           # TypeScript types for preset
│   └── theme.css             # @theme block for Tailwind v4
│
├── shadcn/
│   ├── theme.css             # ShadCN theme variables
│   ├── colors.json           # Color palette for shadcn init
│   └── components.json       # ShadCN configuration
│
├── mkdocs/
│   ├── extra.css             # MkDocs Material theme overrides
│   └── palette.yml           # Color scheme configuration
│
├── scss/
│   └── _variables.scss       # SCSS variables
│
├── tokens.json               # Raw token data
└── tokens.w3c.json           # W3C Design Tokens format
```

## Generated package.json

```json
{
  "name": "@plexusone/design-tokens",
  "version": "0.1.0",
  "description": "Design tokens for PlexusOne",
  "main": "index.js",
  "module": "index.mjs",
  "types": "index.d.ts",
  "exports": {
    ".": {
      "import": "./index.mjs",
      "require": "./index.js",
      "types": "./index.d.ts"
    },
    "./css": "./css/tokens.css",
    "./tailwind": {
      "import": "./tailwind/preset.js",
      "types": "./tailwind/preset.d.ts"
    },
    "./shadcn": "./shadcn/theme.css",
    "./mkdocs": "./mkdocs/extra.css",
    "./scss": "./scss/_variables.scss"
  },
  "files": ["css", "tailwind", "shadcn", "mkdocs", "scss", "*.js", "*.mjs", "*.d.ts", "*.json"],
  "keywords": ["design-tokens", "design-system", "css", "tailwind"],
  "repository": {
    "type": "git",
    "url": "https://github.com/plexusone/plexusone-design-system"
  },
  "license": "MIT",
  "generatedBy": "design-system-spec"
}
```

## Framework-Specific Outputs

### Tailwind CSS v4

**tailwind/preset.js:**
```javascript
/** @type {import('tailwindcss').Config} */
export default {
  theme: {
    extend: {
      colors: {
        primary: 'var(--color-primary)',
        'primary-light': 'var(--color-primary-light)',
        // ...
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'ui-monospace', 'monospace'],
      },
      spacing: {
        // Uses CSS variables for Tailwind v4
      },
    },
  },
}
```

**tailwind/theme.css:**
```css
@import "tailwindcss";

@theme {
  --color-primary: #06b6d4;
  --color-primary-light: #22d3ee;
  /* ... */
}
```

### ShadCN/UI

**shadcn/theme.css:**
```css
@layer base {
  :root {
    --background: 0 0% 100%;
    --foreground: 222.2 84% 4.9%;
    --primary: 187 94% 43%;
    --primary-foreground: 210 40% 98%;
    /* HSL format for ShadCN */
  }

  .dark {
    --background: 222.2 84% 4.9%;
    --foreground: 210 40% 98%;
    /* Dark mode overrides */
  }
}
```

**shadcn/colors.json:**
```json
{
  "primary": {
    "DEFAULT": "hsl(187, 94%, 43%)",
    "foreground": "hsl(210, 40%, 98%)"
  }
}
```

### MkDocs Material

**mkdocs/extra.css:**
```css
[data-md-color-scheme="slate"] {
  --md-primary-fg-color: #06b6d4;
  --md-accent-fg-color: #8b5cf6;
  /* ... */
}
```

**mkdocs/palette.yml:**
```yaml
palette:
  - scheme: slate
    primary: custom
    accent: custom
    toggle:
      icon: material/brightness-4
      name: Switch to light mode
  - scheme: default
    primary: custom
    accent: custom
    toggle:
      icon: material/brightness-7
      name: Switch to dark mode
```

## Implementation Plan

### Phase 1: Core Package Generator

1. Add `PackageGeneratorOptions` struct to `sdk/go/`
2. Implement `GeneratePackage()` method on `DesignSystem`
3. Generate `package.json` from `Meta`
4. Generate `index.js`, `index.mjs`, `index.d.ts` exports
5. Add `--package` flag to `dss generate` command

### Phase 2: Framework Targets

1. **css** - Reuse existing `GenerateCSS()` with `css-vars` format
2. **tailwind** - New `GenerateTailwindPreset()` method
3. **shadcn** - New `GenerateShadCNTheme()` method
4. **mkdocs-material** - Reuse existing `generateMkDocsMaterialCSS()`
5. **scss** - Reuse existing `GenerateCSS()` with `scss` format

### Phase 3: TypeScript Types

1. Generate `index.d.ts` with token type definitions
2. Generate `tailwind/preset.d.ts` for Tailwind preset types
3. Export token constants as typed objects

### Phase 4: Documentation

1. Auto-generate `README.md` with usage examples
2. Document available exports and imports
3. Include version and generation timestamp

## Testing

```bash
# Generate package from minimal-system example
dss generate -d ./examples/minimal-system --package ./test-output

# Verify package structure
ls -la ./test-output/

# Validate package.json
node -e "require('./test-output/package.json')"

# Test Tailwind preset
npx tailwindcss -i ./test-output/tailwind/theme.css -o /dev/null
```

## Usage Examples

### In a Tailwind Project

```javascript
// tailwind.config.js
import preset from '@plexusone/design-tokens/tailwind'

export default {
  presets: [preset],
  content: ['./src/**/*.{js,ts,jsx,tsx}'],
}
```

### In a ShadCN Project

```bash
# Copy theme to your project
cp node_modules/@plexusone/design-tokens/shadcn/theme.css ./src/styles/
```

### In MkDocs

```yaml
# mkdocs.yml
extra_css:
  - node_modules/@plexusone/design-tokens/mkdocs/extra.css
```

### Direct CSS Import

```css
@import '@plexusone/design-tokens/css/tokens.css';

.my-button {
  background: var(--color-primary);
  color: var(--color-text);
}
```

## Success Criteria

1. Generated package passes `npm publish --dry-run`
2. Tailwind preset works with `npx tailwindcss`
3. ShadCN theme integrates without errors
4. MkDocs Material theme renders correctly
5. TypeScript types provide autocomplete for all tokens
6. Package size is reasonable (< 50KB uncompressed)

## Related

- [TASKS.md](../../TASKS.md) - Project task tracking
- [cli.md](../cli.md) - CLI reference
- [W3C Design Tokens](https://design-tokens.github.io/community-group/format/) - Standard format
