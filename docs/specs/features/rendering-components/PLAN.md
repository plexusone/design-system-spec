# Rendering Components - Implementation Plan

## Phase 1: Live Demo Infrastructure

### 1.1 Update HTML Templates with Material Web CDN

**Files to modify:**
- `sdk/go/templates/html/layout.html`
- `sdk/go/templates/html/component.html`

**Tasks:**
1. Add Material Web import map to layout.html
2. Add Material Symbols font
3. Create demo container structure in component.html
4. Add variant toggle buttons

### 1.2 Create Demo Renderer JavaScript

**New files:**
- `sdk/go/templates/html/static/demo.js`

**Tasks:**
1. Create ComponentDemo class
2. Map DSS component IDs to Material Web tags
3. Implement variant switching
4. Add state controls (disabled, loading, etc.)

### 1.3 Update gen_html.go

**Tasks:**
1. Embed demo.js in static files
2. Add LiveDemos option to HTMLOptions
3. Generate demo initialization code per component

## Phase 2: React/Shadcn Generator

### 2.1 Create Generator Interface

**New file:** `sdk/go/generator.go`

```go
type CodeGenerator interface {
    Generate(component Component) (string, error)
    GenerateTypes(component Component) (string, error)
    FileExtension() string
    TargetName() string
}
```

### 2.2 Implement React Generator

**New file:** `sdk/go/gen_react.go`

**Tasks:**
1. Define ReactGenerator struct
2. Create TSX template
3. Map DSS props to React props interface
4. Map DSS events to React event handlers
5. Generate variant handling (switch/if statements)
6. Add accessibility attributes

### 2.3 Implement Shadcn Variant

**New file:** `sdk/go/gen_react_shadcn.go`

**Tasks:**
1. Extend ReactGenerator with Shadcn patterns
2. Use Radix primitive imports
3. Generate Tailwind class utilities
4. Create variants with cva (class-variance-authority)

## Phase 3: iOS Generators

### 3.1 SwiftUI Generator

**New file:** `sdk/go/gen_swift.go`

**Tasks:**
1. Define SwiftUIGenerator struct
2. Create Swift View template
3. Map DSS props to Swift properties
4. Generate ViewModifiers for variants
5. Add accessibility modifiers
6. Handle @State/@Binding patterns

### 3.2 UIKit Generator (Optional)

**New file:** `sdk/go/gen_uikit.go`

**Tasks:**
1. Generate UIView subclass
2. Create programmatic layout code
3. Map props to UIKit properties

## Phase 4: Android Generators

### 4.1 Jetpack Compose Generator

**New file:** `sdk/go/gen_compose.go`

**Tasks:**
1. Define ComposeGenerator struct
2. Create Kotlin @Composable template
3. Map DSS props to Compose parameters
4. Generate Material 3 theming
5. Add Semantics for accessibility
6. Handle Modifier chains

### 4.2 XML Views Generator (Optional)

**New file:** `sdk/go/gen_android_xml.go`

**Tasks:**
1. Generate XML layout files
2. Generate attrs.xml for custom attributes
3. Generate styles.xml for variants

## Phase 5: Flutter Generator

### 5.1 Implement Flutter Generator

**New file:** `sdk/go/gen_flutter.go`

**Tasks:**
1. Define FlutterGenerator struct
2. Create Dart Widget template
3. Map DSS props to constructor parameters
4. Generate StatelessWidget or StatefulWidget
5. Use Material 3 Flutter widgets
6. Add Semantics widget for accessibility

## Phase 6: CLI Integration

### 6.1 Extend Generate Command

**File:** `cmd/dss/cmd/generate.go`

**Tasks:**
1. Add --react, --swift, --compose, --flutter flags
2. Add --all flag for all targets
3. Output files to specified directories
4. Add --types-only for interface generation

### 6.2 Extend Render Command

**File:** `cmd/dss/cmd/render.go`

**Tasks:**
1. Add --live-demos flag
2. Add --demo-cdn flag for CDN URL override
3. Add --code-examples to include generated code

## Phase 7: Documentation & Testing

### 7.1 Documentation

**Tasks:**
1. Update CLI help text
2. Create usage examples in README
3. Document generator options
4. Add troubleshooting guide

### 7.2 Testing

**Tasks:**
1. Unit tests for each generator
2. Integration tests (code compiles)
3. Snapshot tests for output stability
4. Demo rendering tests (Playwright)

## Implementation Order

```
Week 1: Phase 1 (Live Demos)
├── 1.1 HTML templates
├── 1.2 Demo JavaScript
└── 1.3 gen_html.go updates

Week 2: Phase 2 (React)
├── 2.1 Generator interface
├── 2.2 React generator
└── 2.3 Shadcn variant

Week 3: Phase 3-4 (iOS/Android)
├── 3.1 SwiftUI generator
└── 4.1 Compose generator

Week 4: Phase 5-6 (Flutter/CLI)
├── 5.1 Flutter generator
└── 6.1-6.2 CLI integration

Week 5: Phase 7 (Docs/Tests)
├── 7.1 Documentation
└── 7.2 Testing
```

## Verification Checklist

### Phase 1 Complete
- [ ] Material Web components render in HTML docs
- [ ] Variant toggles work
- [ ] State controls (disabled) work
- [ ] Light/dark theme works

### Phase 2 Complete
- [ ] `dss generate --react` produces TSX files
- [ ] Generated code compiles with TypeScript
- [ ] Props interface matches DSS spec
- [ ] Accessibility attributes present

### Phase 3 Complete
- [ ] `dss generate --swift` produces Swift files
- [ ] Generated code compiles with Xcode
- [ ] SwiftUI previews work

### Phase 4 Complete
- [ ] `dss generate --compose` produces Kotlin files
- [ ] Generated code compiles with Gradle
- [ ] @Preview composables work

### Phase 5 Complete
- [ ] `dss generate --flutter` produces Dart files
- [ ] Generated code passes `dart analyze`
- [ ] Widgets render in Flutter

### Phase 6 Complete
- [ ] All CLI flags work
- [ ] --all generates all targets
- [ ] Output directories created correctly

### Phase 7 Complete
- [ ] All tests pass
- [ ] Documentation complete
- [ ] README updated
