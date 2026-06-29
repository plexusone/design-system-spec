# CLI Reference

The `dss` command-line tool generates code artifacts and validates implementations against design system specifications.

## Installation

```bash
go install github.com/plexusone/design-system-spec/cmd/dss@latest
```

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--dir` | `-d` | Directory containing the design system spec (default: `.`) |
| `--help` | `-h` | Help for any command |
| `--version` | `-v` | Version information |

## Commands

### dss info

Display information about a design system specification.

```bash
dss info [flags]
```

**Examples:**

```bash
# Show info for current directory
dss info

# Show info for a specific directory
dss info -d ./my-design-system
```

**Output includes:**

- Design system name and version
- Description
- Principles count
- Foundation token counts (colors, typography, spacing)
- Component list with variant counts
- Accessibility target level
- Validation status

---

### dss generate

Generate code artifacts from a design system specification.

```bash
dss generate [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--css` | Output path for CSS file |
| `--types` | Output path for TypeScript types |
| `--llm` | Output path for LLM prompt/context |
| `--css-format` | CSS format: `tailwind4` (default), `css-vars`, `scss` |

**Examples:**

```bash
# Preview all outputs to stdout
dss generate

# Generate CSS only
dss generate --css ./src/index.css

# Generate all artifacts
dss generate \
  --css ./src/index.css \
  --types ./src/lib/design-system-types.ts \
  --llm ./DESIGN_CONTEXT.md

# Use different spec directory
dss generate -d ./design-system --css ./web/src/styles.css

# Generate SCSS variables instead of Tailwind
dss generate --css-format scss --css ./src/variables.scss

# Generate standard CSS custom properties
dss generate --css-format css-vars --css ./src/vars.css

# Generate NPM package with all targets
dss generate --package ./dist --targets all

# Generate NPM package with specific targets
dss generate --package ./dist --targets css,tailwind,shadcn
```

**NPM Package Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--package` | `-p` | Output directory for NPM package |
| `--scope` | `-s` | NPM scope (e.g., `@myorg`) |
| `--name` | `-n` | Package name (default: `design-tokens`) |
| `--targets` | `-t` | Comma-separated targets (default: `css,tailwind`) |
| `--dry-run` | | Preview without writing files |

**Available Targets:**

| Target | Description |
|--------|-------------|
| `css` | CSS custom properties |
| `tailwind` | Tailwind CSS v4 preset |
| `shadcn` | ShadCN/UI theme variables |
| `mkdocs-material` | MkDocs Material theme |
| `scss` | SCSS variables |
| `json` | Raw JSON tokens |
| `w3c` | W3C Design Tokens format |
| `all` | All of the above |

**CSS Formats:**

| Format | Description |
|--------|-------------|
| `tailwind4` | Tailwind CSS v4 `@theme` block with CSS custom properties |
| `css-vars` | Standard CSS `:root` with custom properties |
| `scss` | SCSS variables (`$color-primary`, etc.) |

---

### dss validate

Validate that component implementations comply with the design system specification.

```bash
dss validate <components-dir> [flags]
```

**Arguments:**

| Argument | Description |
|----------|-------------|
| `components-dir` | Directory containing React/TypeScript component files |

**Flags:**

| Flag | Description |
|------|-------------|
| `--json` | Output results as JSON (useful for CI) |

**Examples:**

```bash
# Validate components
dss validate ./src/components

# Validate with different spec directory
dss validate -d ./design-system ./web/src/components

# JSON output for CI integration
dss validate --json ./src/components
```

**Validation Checks:**

| Rule | Severity | Description |
|------|----------|-------------|
| `no-hardcoded-colors` | Warning | Hex/RGB/HSL colors should use CSS variables |
| `use-spacing-scale` | Warning | Pixel values should use spacing scale |
| `img-alt-required` | Error | Images must have alt attributes |
| `button-accessible-name` | Warning | Icon-only buttons need aria-label |
| `valid-variant` | Error | Variant values must match component spec |
| `single-primary-button` | Warning | Only one primary button per view |
| `no-nested-cards` | Warning | Cards should not be nested |

**Exit Codes:**

| Code | Meaning |
|------|---------|
| 0 | Validation passed (may have warnings) |
| 1 | Validation failed with errors |

---

### dss bind

Generate theme bindings from design system token mappings.

```bash
dss bind [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Output file (default: stdout) |
| `--format` | `-f` | Output format: `css`, `typescript`, `scss` (default: `css`) |
| `--strategy` | | Mapping strategy: `explicit`, `semantic`, `inherit` (default: `explicit`) |

**Examples:**

```bash
# Generate CSS bindings to stdout
dss bind

# Generate CSS to file
dss bind --output ./theme.css

# Generate TypeScript constants
dss bind --format typescript --output ./theme.ts

# Use semantic auto-mapping
dss bind --strategy semantic --output ./theme.css

# Generate SCSS variables
dss bind --format scss --output ./theme.scss
```

**Mapping Strategies:**

| Strategy | Description |
|----------|-------------|
| `explicit` | Only use defined `from`/`to` mappings, skip unmapped tokens |
| `semantic` | Auto-map by semantic field, fall back to component defaults |
| `inherit` | Use component defaults for all unmapped tokens |

---

### dss contract

Display and validate component theming contracts.

```bash
dss contract <subcommand> [flags]
```

**Subcommands:**

#### dss contract show

Display a component's theming contract.

```bash
dss contract show <component-id> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--json` | Output as JSON instead of table |

**Examples:**

```bash
# Show button's theming contract
dss contract show button

# JSON output
dss contract show button --json
```

**Output:**

```
Theming Contract: Button
Prefix: --btn

ID          CSS PROPERTY      SEMANTIC  LIGHT DEFAULT  DARK DEFAULT
--          ------------      --------  -------------  ------------
background  --btn-background  primary   #0066CC        #3399FF
text        --btn-text        text      #FFFFFF        #0A0E1A
```

#### dss contract validate

Validate all theming contracts in the design system.

```bash
dss contract validate [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--json` | Output as JSON instead of table |

**Examples:**

```bash
# Validate all contracts
dss contract validate

# JSON output for CI
dss contract validate --json
```

**Validation Checks:**

| Check | Severity | Description |
|-------|----------|-------------|
| Prefix required | Error | Contract must have a prefix starting with `--` |
| Token ID required | Error | Each token must have a unique ID |
| CSS property required | Error | Each token must have a cssProperty |
| Duplicate token ID | Error | Token IDs must be unique within a contract |
| CSS property prefix | Warning | cssProperty should start with the contract prefix |
| Invalid semantic | Warning | Semantic value should be from the allowed set |
| Missing defaults | Warning | Tokens should have defaultLight and defaultDark |

**Output:**

```
✓ button: OK
⚠ card: 2 warnings
    WARN  [bg] themingContract.tokens[0].defaultLight: defaultLight not provided
✗ input: 1 error, 0 warnings
    ERROR themingContract.prefix: prefix is required

Validation passed: 3 contracts validated, 2 warnings
```

---

## Usage in Makefiles

Example `Makefile` for a project using DSS:

```makefile
DSS ?= dss
WEB_SRC = ./web/src

.PHONY: generate validate package

generate:
	@$(DSS) generate \
		--css $(WEB_SRC)/index.css \
		--types $(WEB_SRC)/lib/design-system-types.ts \
		--llm ./DESIGN_CONTEXT.md

validate:
	@$(DSS) validate $(WEB_SRC)/components

# CI validation with JSON output
validate-ci:
	@$(DSS) validate --json $(WEB_SRC)/components

# Check if generated files are current
check:
	@$(DSS) generate --css /tmp/check.css
	@diff -q $(WEB_SRC)/index.css /tmp/check.css

# Generate NPM package
package:
	@$(DSS) generate --package ./dist --targets all
```

---

## Using Generated NPM Packages

### In a Tailwind Project

```javascript
// tailwind.config.js
import preset from '@myorg/design-tokens/tailwind'

export default {
  presets: [preset],
  content: ['./src/**/*.{js,ts,jsx,tsx}'],
}
```

### In a ShadCN Project

```css
/* Import theme variables */
@import '@myorg/design-tokens/shadcn/theme.css';
```

### In MkDocs

```yaml
# mkdocs.yml
extra_css:
  - node_modules/@myorg/design-tokens/mkdocs/extra.css
```

### Direct CSS Import

```css
@import '@myorg/design-tokens/css/tokens.css';

.my-button {
  background: var(--color-primary);
  color: var(--color-text);
}
```

### JavaScript/TypeScript

```javascript
import { colors, spacing } from '@myorg/design-tokens'

console.log(colors.primary) // #...
```

---

## CI/CD Integration

### GitHub Actions

```yaml
name: Design System Validation

on: [push, pull_request]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Install dss
        run: go install github.com/plexusone/design-system-spec/cmd/dss@latest

      - name: Validate components
        run: dss validate -d ./design-system ./src/components

      - name: Check generated files are current
        run: |
          dss generate -d ./design-system --css /tmp/index.css
          diff ./src/index.css /tmp/index.css
```

### Pre-commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit

# Validate design system compliance
dss validate ./src/components

if [ $? -ne 0 ]; then
  echo "Design system validation failed!"
  exit 1
fi
```
