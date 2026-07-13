# Rendering Components - Roadmap

## Version History

| Version | Date | Status | Description |
|---------|------|--------|-------------|
| 0.1.0 | TBD | Planned | Live demos with Material Web |
| 0.2.0 | TBD | Planned | React/Shadcn generator |
| 0.3.0 | TBD | Planned | iOS (SwiftUI) generator |
| 0.4.0 | TBD | Planned | Android (Compose) generator |
| 0.5.0 | TBD | Planned | Flutter generator |
| 1.0.0 | TBD | Planned | Full release with all targets |

## v0.1.0 - Live Demo Infrastructure

**Goal:** Render Material Web components in DSS HTML documentation

### Features
- [ ] Material Web CDN integration
- [ ] Live component rendering in docs
- [ ] Variant toggle controls
- [ ] State manipulation (disabled, loading)
- [ ] Light/dark theme toggle

### Technical Tasks
- [ ] Update `layout.html` with import map
- [ ] Create `demo.js` controller
- [ ] Update `component.html` template
- [ ] Embed static assets in Go binary

### Dependencies
- material-web via esm.run CDN
- Google Fonts (Roboto, Material Symbols)

---

## v0.2.0 - React/Shadcn Generator

**Goal:** Generate React components from DSS specs

### Features
- [ ] TSX component generation
- [ ] TypeScript interface generation
- [ ] Shadcn/ui variant with Tailwind
- [ ] Accessibility attributes
- [ ] Event handler mapping

### CLI Commands
```bash
dss generate --react ./output/
dss generate --react --shadcn ./output/
```

### Output Structure
```
output/
├── components/
│   ├── Button.tsx
│   ├── Checkbox.tsx
│   └── ...
└── types/
    └── index.ts
```

---

## v0.3.0 - iOS SwiftUI Generator

**Goal:** Generate SwiftUI Views from DSS specs

### Features
- [ ] Swift View generation
- [ ] ViewModifier for variants
- [ ] @State/@Binding properties
- [ ] Accessibility modifiers
- [ ] Environment-based theming

### CLI Commands
```bash
dss generate --swift ./output/ios/
```

### Output Structure
```
output/ios/
├── Components/
│   ├── MDButton.swift
│   ├── MDCheckbox.swift
│   └── ...
└── Theme/
    └── MaterialTheme.swift
```

---

## v0.4.0 - Android Compose Generator

**Goal:** Generate Jetpack Compose functions from DSS specs

### Features
- [ ] @Composable function generation
- [ ] Material 3 theming
- [ ] Modifier parameters
- [ ] Semantics for accessibility
- [ ] Preview composables

### CLI Commands
```bash
dss generate --compose ./output/android/
```

### Output Structure
```
output/android/
├── components/
│   ├── MDButton.kt
│   ├── MDCheckbox.kt
│   └── ...
└── theme/
    └── MaterialTheme.kt
```

---

## v0.5.0 - Flutter Generator

**Goal:** Generate Flutter Widgets from DSS specs

### Features
- [ ] Dart Widget generation
- [ ] StatelessWidget/StatefulWidget
- [ ] Material 3 Flutter theming
- [ ] Semantics widget for a11y
- [ ] Theme extensions

### CLI Commands
```bash
dss generate --flutter ./output/flutter/
```

### Output Structure
```
output/flutter/
├── lib/
│   ├── components/
│   │   ├── md_button.dart
│   │   ├── md_checkbox.dart
│   │   └── ...
│   └── theme/
│       └── material_theme.dart
└── pubspec.yaml
```

---

## v1.0.0 - Full Release

**Goal:** Production-ready code generation for all targets

### Features
- [ ] All generators stable
- [ ] Comprehensive test coverage
- [ ] Documentation complete
- [ ] CI/CD integration examples

### CLI Commands
```bash
# Generate all targets
dss generate --all ./output/

# Generate with options
dss generate --react --include-comments --accessibility ./output/
```

### Quality Gates
- [ ] 90%+ prop coverage
- [ ] Generated code compiles on all platforms
- [ ] Accessibility audit passes
- [ ] Performance benchmarks met

---

## Future Considerations (Post-1.0)

### Additional Targets
- Vue SFC generator
- Svelte component generator
- Web Components (Lit) generator
- React Native generator

### Advanced Features
- Visual regression testing integration
- Storybook story generation
- Component playground hosting
- AI-assisted code refinement

### Platform Extensions
- Figma plugin integration
- VS Code extension
- GitHub Actions for CI
- npm/pub/pod package publishing

## Milestones

```
┌─────────┬─────────┬─────────┬─────────┬─────────┬─────────┐
│  v0.1   │  v0.2   │  v0.3   │  v0.4   │  v0.5   │  v1.0   │
│  Demos  │  React  │  Swift  │ Compose │ Flutter │ Release │
├─────────┼─────────┼─────────┼─────────┼─────────┼─────────┤
│ Week 1  │ Week 2  │ Week 3  │ Week 3  │ Week 4  │ Week 5  │
└─────────┴─────────┴─────────┴─────────┴─────────┴─────────┘
```

## Success Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Prop Coverage | 80%+ | Props generated / Props in spec |
| Compile Success | 100% | Generated code compiles |
| Demo Load Time | < 2s | Time to interactive |
| A11y Coverage | 100% | Components with ARIA |
| Test Coverage | 80%+ | Lines covered |
