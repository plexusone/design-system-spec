# Rendering Components - Technical Requirements

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     DSS Component Spec                           │
│                    (JSON specification)                          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Code Generators (sdk/go/)                      │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌────────────┐│
│  │ gen_react   │ │ gen_swift   │ │ gen_compose │ │ gen_flutter││
│  │ .go         │ │ .go         │ │ .go         │ │ .go        ││
│  └─────────────┘ └─────────────┘ └─────────────┘ └────────────┘│
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Output Targets                              │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌────────────┐│
│  │ .tsx        │ │ .swift      │ │ .kt         │ │ .dart      ││
│  │ React/Shadcn│ │ SwiftUI     │ │ Compose     │ │ Flutter    ││
│  └─────────────┘ └─────────────┘ └─────────────┘ └────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

## Code Generator Interface

```go
// sdk/go/generator.go

// CodeGenerator generates platform-specific code from DSS specs
type CodeGenerator interface {
    // Generate produces code for a component
    Generate(component Component) (string, error)

    // GenerateTypes produces type definitions
    GenerateTypes(component Component) (string, error)

    // FileExtension returns the output file extension
    FileExtension() string

    // TargetName returns the generator target name
    TargetName() string
}

// GeneratorOptions configures code generation
type GeneratorOptions struct {
    // IncludeComments adds documentation comments
    IncludeComments bool

    // IncludeAccessibility adds a11y attributes
    IncludeAccessibility bool

    // Style formatting style (e.g., "tailwind", "css-modules")
    Style string

    // TypesOnly generates only type definitions
    TypesOnly bool
}
```

## Generator Implementations

### 1. React Generator (`gen_react.go`)

```go
type ReactGenerator struct {
    Options GeneratorOptions
    // UseShadcn uses Radix primitives + Tailwind
    UseShadcn bool
}

func (g *ReactGenerator) Generate(c Component) (string, error) {
    // Generate TSX component with:
    // - Props interface from component.Props
    // - Event handlers from component.Events
    // - Slots as children/render props
    // - Accessibility attributes
    // - Variant handling via props
}
```

**Output Example:**
```tsx
import * as React from "react";
import { cn } from "@/lib/utils";

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "filled" | "outlined" | "text" | "elevated" | "tonal";
  disabled?: boolean;
  loading?: boolean;
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = "filled", disabled, loading, children, ...props }, ref) => {
    return (
      <button
        ref={ref}
        className={cn(buttonVariants({ variant }), className)}
        disabled={disabled || loading}
        aria-busy={loading}
        {...props}
      >
        {loading && <Spinner />}
        {children}
      </button>
    );
  }
);
```

### 2. SwiftUI Generator (`gen_swift.go`)

```go
type SwiftUIGenerator struct {
    Options GeneratorOptions
}

func (g *SwiftUIGenerator) Generate(c Component) (string, error) {
    // Generate SwiftUI View with:
    // - @State/@Binding properties
    // - ViewModifiers for variants
    // - Accessibility modifiers
    // - Environment values for theming
}
```

**Output Example:**
```swift
import SwiftUI

public struct MDButton: View {
    public enum Variant {
        case filled, outlined, text, elevated, tonal
    }

    let label: String
    var variant: Variant = .filled
    var disabled: Bool = false
    var action: () -> Void

    public var body: some View {
        Button(action: action) {
            Text(label)
                .padding(.horizontal, 24)
                .padding(.vertical, 10)
        }
        .buttonStyle(MDButtonStyle(variant: variant))
        .disabled(disabled)
        .accessibilityLabel(label)
    }
}
```

### 3. Jetpack Compose Generator (`gen_compose.go`)

```go
type ComposeGenerator struct {
    Options GeneratorOptions
}

func (g *ComposeGenerator) Generate(c Component) (string, error) {
    // Generate Kotlin @Composable with:
    // - Parameters from props
    // - Modifier chain for styling
    // - Material 3 theming
    // - Semantics for accessibility
}
```

**Output Example:**
```kotlin
@Composable
fun MDButton(
    label: String,
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
    variant: ButtonVariant = ButtonVariant.Filled,
    enabled: Boolean = true,
) {
    when (variant) {
        ButtonVariant.Filled -> FilledButton(onClick, modifier, enabled) {
            Text(label)
        }
        ButtonVariant.Outlined -> OutlinedButton(onClick, modifier, enabled) {
            Text(label)
        }
        // ...
    }
}
```

### 4. Flutter Generator (`gen_flutter.go`)

```go
type FlutterGenerator struct {
    Options GeneratorOptions
}

func (g *FlutterGenerator) Generate(c Component) (string, error) {
    // Generate Dart Widget with:
    // - StatelessWidget or StatefulWidget
    // - Constructor parameters from props
    // - build() method with Material widgets
    // - Semantics widget for accessibility
}
```

**Output Example:**
```dart
class MDButton extends StatelessWidget {
  final String label;
  final VoidCallback? onPressed;
  final ButtonVariant variant;
  final bool disabled;

  const MDButton({
    super.key,
    required this.label,
    this.onPressed,
    this.variant = ButtonVariant.filled,
    this.disabled = false,
  });

  @override
  Widget build(BuildContext context) {
    return switch (variant) {
      ButtonVariant.filled => FilledButton(
        onPressed: disabled ? null : onPressed,
        child: Text(label),
      ),
      ButtonVariant.outlined => OutlinedButton(
        onPressed: disabled ? null : onPressed,
        child: Text(label),
      ),
      // ...
    };
  }
}
```

## Live Demo Integration

### HTML Template Updates

Update `sdk/go/templates/html/component.html`:

```html
{{define "component"}}
<article class="component-page">
  <h1>{{.Name}}</h1>

  <!-- Live Demo Section -->
  <section class="component-demo">
    <h2>Live Demo</h2>
    <div class="demo-container" id="demo-{{.ID}}">
      <!-- Material Web component renders here -->
    </div>
    <div class="demo-controls">
      {{range .Variants}}
      <button data-variant="{{.ID}}">{{.Name}}</button>
      {{end}}
    </div>
  </section>

  <!-- Code Examples -->
  <section class="component-code">
    <div class="code-tabs">
      <button data-tab="web">Web Component</button>
      <button data-tab="react">React</button>
      <button data-tab="swift">SwiftUI</button>
      <button data-tab="compose">Compose</button>
      <button data-tab="flutter">Flutter</button>
    </div>
    <div class="code-panels">
      {{/* Generated code for each target */}}
    </div>
  </section>
</article>
{{end}}
```

### Material Web CDN Integration

Add to `sdk/go/templates/html/layout.html`:

```html
<head>
  <!-- Material Web via CDN -->
  <script type="importmap">
  {
    "imports": {
      "@material/web/": "https://esm.run/@material/web/"
    }
  }
  </script>
  <script type="module">
    import '@material/web/all.js';
    import {styles as typescaleStyles} from '@material/web/typography/md-typescale-styles.js';
    document.adoptedStyleSheets.push(typescaleStyles.styleSheet);
  </script>
  <link href="https://fonts.googleapis.com/css2?family=Roboto:wght@400;500;700&display=swap" rel="stylesheet">
  <link href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined" rel="stylesheet">
</head>
```

### Component Demo Renderer

```javascript
// sdk/go/templates/html/static/demo.js

class ComponentDemo {
  constructor(componentId, variants) {
    this.container = document.getElementById(`demo-${componentId}`);
    this.variants = variants;
    this.currentVariant = variants[0]?.id || 'default';
  }

  render() {
    const html = this.getComponentHTML(this.currentVariant);
    this.container.innerHTML = html;
  }

  getComponentHTML(variant) {
    // Maps DSS component ID to Material Web tag
    const tagMap = {
      'button': {
        'filled': '<md-filled-button>Button</md-filled-button>',
        'outlined': '<md-outlined-button>Button</md-outlined-button>',
        'text': '<md-text-button>Button</md-text-button>',
        'elevated': '<md-elevated-button>Button</md-elevated-button>',
        'tonal': '<md-filled-tonal-button>Button</md-filled-tonal-button>',
      },
      // ... other components
    };
    return tagMap[this.componentId]?.[variant] || '';
  }

  setVariant(variantId) {
    this.currentVariant = variantId;
    this.render();
  }
}
```

## CLI Commands

### Generate Command Extensions

```bash
# Generate React components
dss generate -d ./specs/v3 --react ./output/react/

# Generate SwiftUI views
dss generate -d ./specs/v3 --swift ./output/ios/

# Generate Jetpack Compose
dss generate -d ./specs/v3 --compose ./output/android/

# Generate Flutter widgets
dss generate -d ./specs/v3 --flutter ./output/flutter/

# Generate all targets
dss generate -d ./specs/v3 --all ./output/
```

### Render Command Extensions

```bash
# Render HTML with live demos
dss render -d ./specs/v3 --output ./docs --live-demos

# Render with specific demo CDN
dss render -d ./specs/v3 --output ./docs --demo-cdn "https://esm.run/@material/web/"
```

## File Structure

```
sdk/go/
├── generator.go          # Interface definition
├── gen_react.go          # React/TSX generator
├── gen_react_shadcn.go   # Shadcn variant
├── gen_swift.go          # SwiftUI generator
├── gen_compose.go        # Jetpack Compose generator
├── gen_flutter.go        # Flutter generator
├── gen_vue.go            # Vue SFC generator
├── gen_webcomponent.go   # Lit element generator
└── templates/
    ├── react/
    │   └── component.tsx.tmpl
    ├── swift/
    │   └── view.swift.tmpl
    ├── compose/
    │   └── composable.kt.tmpl
    ├── flutter/
    │   └── widget.dart.tmpl
    └── html/
        ├── layout.html       # Updated with CDN
        ├── component.html    # Updated with demos
        └── static/
            └── demo.js       # Demo controller
```

## Dependencies

### Go Dependencies
- `text/template` - Code generation templates
- Existing DSS SDK types

### CDN Dependencies (Runtime)
- `@material/web` via esm.run
- Google Fonts (Roboto, Material Symbols)
- Optional: playground-elements for code editing

## Testing Strategy

1. **Unit Tests**: Each generator produces valid syntax
2. **Integration Tests**: Generated code compiles
3. **Snapshot Tests**: Output matches expected templates
4. **Demo Tests**: Live demos render correctly (Playwright)
