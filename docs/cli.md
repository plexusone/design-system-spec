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
```

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

## Usage in Makefiles

Example `Makefile` for a project using DSS:

```makefile
DSS ?= dss
WEB_SRC = ./web/src

.PHONY: generate validate

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
