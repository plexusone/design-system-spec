# Case Study: PlexusOne Design System

How PlexusOne uses DSS to maintain visual consistency across 1 website, 33 MkDocs documentation sites, and Marp presentation decks.

## The Challenge

PlexusOne is a suite of AI agent tools and libraries spanning multiple products:

- **OmniLLM** - Multi-provider LLM abstraction
- **OmniVault** - Secrets management for AI agents
- **OmniVoice** - Voice AI integration
- **AgentKit** - Agent development framework
- **AssistantKit** - Assistant UI components
- Plus 30+ additional libraries and tools

Each product has its own GitHub repository with documentation. Over time, visual inconsistency crept in:

| Property | Count | Problem |
|----------|-------|---------|
| Main website | 1 | Custom CSS, no system |
| MkDocs sites | 33 | Mix of indigo, teal, readthedocs themes |
| Marp presentations | ~20 | Inconsistent colors, fonts |

### Symptoms of Design Drift

```yaml
# go-elevenlabs/mkdocs.yml
theme:
  palette:
    primary: indigo
    accent: indigo

# plexusone-apps/mkdocs.yml
theme:
  palette:
    primary: teal
    accent: cyan

# opik-go/mkdocs.yml
theme: readthedocs  # Different theme entirely
```

The brand was fragmenting. Cyan, purple, and pink were the intended accent colors, but many sites used Material's default indigo. Dark mode was inconsistent. Fonts varied.

## The Solution: Design System as Code

Instead of manually updating 33+ repos, we created a **single source of truth** using DSS.

### One Spec, Multiple Outputs

```
                    ┌─────────────────────────────────────────────────────────┐
                    │            plexusone-specs/design-system/               │
                    │                                                         │
                    │   colors.json  typography.json  spacing.json  ...       │
                    │                                                         │
                    │              THE SINGLE SOURCE OF TRUTH                 │
                    └─────────────────────────┬───────────────────────────────┘
                                              │
                              dss generate    │
                                              │
              ┌───────────────────────────────┼───────────────────────────────┐
              │                               │                               │
              ▼                               ▼                               ▼
    ┌─────────────────┐             ┌─────────────────┐             ┌─────────────────┐
    │   web.css       │             │   mkdocs.css    │             │   slides.css    │
    │   Tailwind v4   │             │   MkDocs Mat.   │             │   Marp Theme    │
    └────────┬────────┘             └────────┬────────┘             └────────┬────────┘
             │                               │                               │
             ▼                               ▼                               ▼
    ┌─────────────────┐             ┌─────────────────┐             ┌─────────────────┐
    │                 │             │  33 MkDocs      │             │  ~20 Marp       │
    │    Website      │             │  Documentation  │             │  Presentation   │
    │                 │             │  Sites          │             │  Decks          │
    └─────────────────┘             └─────────────────┘             └─────────────────┘
```

**The key insight**: Define design tokens once, generate context-appropriate CSS for each medium. When the brand evolves, update the spec and regenerate—all 54+ properties update together.

```
plexusone-specs/
└── design-system/
    ├── meta.json              # Brand metadata
    ├── principles.json        # Design philosophy
    ├── foundations/
    │   ├── colors.json        # Cyan/purple/pink palette
    │   ├── typography.json    # Inter + JetBrains Mono
    │   └── spacing.json       # Spacing scale
    ├── contexts/
    │   ├── web.json           # Website-specific adjustments
    │   ├── slides.json        # Brighter colors for projection
    │   └── docs.json          # MkDocs Material config
    └── components/
        ├── button.json
        └── card.json
```

### Design Principles Encoded

```json
[
  {
    "id": "dark-first",
    "name": "Dark-First",
    "description": "Dark backgrounds are the default. All UI is designed for dark mode first."
  },
  {
    "id": "gradient-accents",
    "name": "Gradient Accents",
    "description": "The cyan-purple-pink gradient is the signature visual element."
  }
]
```

### Core Color Palette

| Token | Value | Usage |
|-------|-------|-------|
| `background` | `#0a0e1a` | Primary dark background |
| `cyan` | `#06b6d4` | Primary accent, CTAs |
| `purple` | `#8b5cf6` | Secondary accent |
| `pink` | `#ec4899` | Tertiary accent |
| `foreground` | `#f1f5f9` | Primary text |

### Context-Aware Generation

The same design tokens generate different outputs per context:

**Web (Tailwind v4)**
```css
@theme {
  --color-cyan: #06b6d4;
  --color-purple: #8b5cf6;
  --font-sans: 'Inter', system-ui, sans-serif;
}
```

**Slides (Marp)** - Brighter for projection
```css
:root {
  --color-cyan: #00d4ff;    /* +15% brightness */
  --color-purple: #b388ff;
  --color-pink: #ff80ab;
}
```

**Docs (MkDocs Material)**
```css
:root {
  --md-primary-fg-color: #06b6d4;
  --md-accent-fg-color: #8b5cf6;
  --md-default-bg-color: #0a0e1a;
}
```

## Implementation

### Step 1: Define the Spec

Created `plexusone-specs/design-system/` with JSON files defining colors, typography, spacing, and components.

### Step 2: Generate Outputs

```bash
# Generate web CSS (Tailwind v4)
dss generate --css output/web.css --css-format tailwind4

# Generate MkDocs theme CSS
dss generate --css output/mkdocs.css --css-format mkdocs-material

# Generate Marp theme CSS
dss generate --css output/slides.css --css-format marp --context slides
```

### Step 3: Distribute via Package

For MkDocs sites, we created a PyPi package:

```bash
pip install plexusone-mkdocs-theme
```

Each of the 33 repos updates their `mkdocs.yml`:

```yaml
theme:
  name: material
  palette:
    scheme: slate
    primary: custom
    accent: custom

extra_css:
  - https://cdn.plexusone.dev/mkdocs/plexusone.css
  # Or if using the package:
  # - plexusone_mkdocs_theme/plexusone.css

plugins:
  - plexusone-theme  # Auto-configures Material theme
```

### Step 4: LLM Context for AI Development

DSS also generates context for AI assistants:

```bash
dss generate --llm DESIGN_CONTEXT.md
```

This creates a prompt that AI coding assistants can use to build compliant UI:

```markdown
# PlexusOne Design System

## Colors
- Primary accent: cyan (#06b6d4) - Use for CTAs, links
- Secondary: purple (#8b5cf6) - Use for focus rings
...

## Components
### Button
- Primary variant: cyan with glow effect
- Don't: Multiple primary buttons in same section
```

## Results

### Before DSS

- 33 sites with 5+ different color schemes
- No documented design principles
- AI assistants had no context for building UI
- Each site updated manually
- No accessibility guidelines
- No content/voice documentation

### After DSS

- Single source of truth in `plexusone-specs/design-system/`
- Consistent dark theme with cyan/purple/pink accents
- AI assistants receive design context automatically
- Updates propagate via package or CDN
- **98.1% DSS coverage (Grade A)** across all 9 categories
- WCAG AA accessibility built into the spec
- Voice, tone, and terminology guidelines documented

### Properties Now Unified

All 54+ properties now share the same design tokens:

| Type | Count | Examples |
|------|-------|----------|
| **Website** | 1 | plexusone.dev |
| **MkDocs Sites** | 33 | omnillm, omnivault, agentkit, assistantkit, go-elevenlabs, opik-go, mcp-go, go-bedrock, go-openrouter, chroma-go... |
| **Marp Decks** | ~20 | Product demos, conference talks, internal training |
| **Total** | **54+** | All using same cyan/purple/pink palette, Inter + JetBrains Mono fonts, dark-first theme |

A single change to `colors.json` propagates to all 54+ properties on next build.

### Maintenance Workflow

When the design system changes:

1. Update JSON files in `plexusone-specs/design-system/`
2. Run `dss generate` to create new outputs
3. Publish new version of `plexusone-mkdocs-theme`
4. Sites pick up changes on next build (or immediately via CDN)

## Key Learnings

### What Worked

1. **JSON as source of truth** - Easy to version, diff, and review
2. **Context-aware generation** - Same tokens, different outputs per medium
3. **Package distribution** - One update propagates to all sites
4. **LLM context generation** - AI can build compliant UI without human intervention

### Challenges

1. **Migration effort** - 33 sites needed mkdocs.yml updates (scripted)
2. **Edge cases** - Some sites had custom extensions needing manual review
3. **Font loading** - Ensuring consistent font loading across all contexts

### Metrics

| Metric | Before | After |
|--------|--------|-------|
| Properties unified | 0 | 54+ (1 website, 33 docs, 20 decks) |
| Time to update all properties | ~2 hours | ~5 minutes |
| Color scheme variations | 5+ | 1 |
| Design documentation | Scattered | Centralized |
| AI compliance rate | Unknown | Measurable via validation |
| DSS coverage score | 0% | 98.1% (Grade A) |
| Documented categories | 0/9 | 9/9 |
| Accessibility compliance | Undefined | WCAG AA |

## Artifacts

The PlexusOne design system generates:

| Output | Format | Purpose |
|--------|--------|---------|
| `web.css` | Tailwind v4 @theme | Main website |
| `mkdocs.css` | CSS variables | Documentation sites |
| `slides.css` | Marp theme | Presentations |
| `types.ts` | TypeScript | Component prop types |
| `DESIGN_CONTEXT.md` | Markdown | LLM prompt context |

The DSS spec itself contains:

| Category | Files | Content |
|----------|-------|---------|
| Meta | `meta.json` | Name, version, description, maintainers |
| Principles | `principles.json` | 4 design principles with examples |
| Foundations | 11 files | Colors, typography, spacing, elevation, motion, grid, breakpoints, border-radius, border-width, opacity, z-index |
| Components | 2 files | Button, Card (with variants, states, props, slots) |
| Patterns | 2 files | Form Layout, Data Table |
| Templates | 2 files | Dashboard, Settings |
| Content | `content.json` | Voice, tone, terminology, writing style |
| Accessibility | `accessibility.json` | WCAG level, contrast, keyboard, screen reader |
| Governance | `governance.json` | Versioning, contribution, deprecation policies |

## DSS Spec Coverage Journey

The PlexusOne design system evolved from a minimal spec to comprehensive coverage through iterative improvement.

### Initial State: Grade F (29.4%)

The initial implementation covered only basic foundations:

```
╔══════════════════════════════════════════════════════════════════════════╗
║ 🟢 Meta             83.3%                                                ║
║ 🟡 Principles       66.7%  (no examples)                                 ║
║ 🟡 Foundations      41.7%  (colors, typography, spacing only)            ║
║ 🟡 Components       42.9%  (no states, no slots)                         ║
║ 🔴 Patterns          0.0%                                                ║
║ 🔴 Templates         0.0%                                                ║
║ 🔴 Content           0.0%                                                ║
║ 🔴 Accessibility     0.0%                                                ║
║ 🔴 Governance        0.0%                                                ║
╠══════════════════════════════════════════════════════════════════════════╣
║ 🔴 Grade: F  Score: 29.4%  (1 complete, 3 partial, 5 missing)            ║
╚══════════════════════════════════════════════════════════════════════════╝
```

### Final State: Grade A (98.1%)

After comprehensive implementation:

```
╔══════════════════════════════════════════════════════════════════════════╗
║ 🟢 Meta             83.3%  (name, version, description)                  ║
║ 🟢 Principles      100.0%  (4 principles with examples)                  ║
║ 🟢 Foundations     100.0%  (colors, typography, spacing, elevation,      ║
║                            motion, grid, breakpoints, border-radius,     ║
║                            border-width, opacity, z-index)               ║
║ 🟢 Components      100.0%  (variants, states, props, slots)              ║
║ 🟢 Patterns        100.0%  (form-layout, data-table)                     ║
║ 🟢 Templates       100.0%  (dashboard, settings)                         ║
║ 🟢 Content         100.0%  (voice, tone, terminology, writing style)     ║
║ 🟢 Accessibility   100.0%  (WCAG level, contrast, keyboard, screen reader)║
║ 🟢 Governance      100.0%  (versioning, contribution, deprecation)       ║
╠══════════════════════════════════════════════════════════════════════════╣
║ 🟢 Grade: A  Score: 98.1%  (9 complete, 0 partial, 0 missing)            ║
╚══════════════════════════════════════════════════════════════════════════╝
```

### What We Added

| Category | Files Added | Purpose |
|----------|-------------|---------|
| **Foundations** | `elevation.json` | Shadow tokens (sm, md, lg, xl) and glow effects |
| | `motion.json` | Animation durations and easing curves |
| | `breakpoints.json` | Responsive breakpoints (sm, md, lg, xl, 2xl) |
| | `border-radius.json` | Corner radius scale |
| | `border-width.json` | Border thickness tokens |
| | `opacity.json` | Opacity scale for overlays |
| | `z-index.json` | Stacking order tokens |
| | `grid.json` | 12-column grid configuration |
| **Accessibility** | `accessibility.json` | WCAG AA compliance, contrast ratios, keyboard navigation, screen reader guidelines |
| **Content** | `content.json` | Voice attributes, tone guidelines per context, terminology (preferred/avoided), writing style |
| **Governance** | `governance.json` | Semver versioning policy, contribution workflow, deprecation rules, ownership |
| **Patterns** | `form-layout.json` | Standard form pattern with validation |
| | `data-table.json` | Sortable, filterable table pattern |
| **Templates** | `dashboard.json` | App shell with sidebar, header, content |
| | `settings.json` | Settings page layout |
| **Components** | Updated `button.json` | Added states (hover, active, focus, disabled, loading) and slots (leftIcon, rightIcon) |
| | Updated `card.json` | Added states (hover, focus, selected) and slots (header, footer, media) |
| **Principles** | Updated `principles.json` | Added examples (good/bad) and rationale for each principle |

### Coverage Validation Command

The `dss coverage` command enables continuous monitoring:

```bash
$ dss coverage -d ./design-system

╔══════════════════════════════════════════════════════════════════════════╗
║                         DSS COVERAGE REPORT                              ║
╠══════════════════════════════════════════════════════════════════════════╣
║ Design System: PlexusOne Design System v1.0.0                            ║
╠══════════════════════════════════════════════════════════════════════════╣
║ 🟢 Grade: A  Score: 98.1%  (9 complete, 0 partial, 0 missing)            ║
╚══════════════════════════════════════════════════════════════════════════╝
```

### Benefits of Comprehensive Coverage

| Benefit | Description |
|---------|-------------|
| **LLM Code Generation** | AI assistants have complete context for generating compliant UI code |
| **Design Decisions Documented** | Principles with examples explain *why*, not just *what* |
| **Accessibility Built-in** | WCAG requirements are part of the spec, not an afterthought |
| **Content Consistency** | Voice and terminology guidelines prevent brand drift |
| **Maintainability** | Governance policies define how changes are made |
| **Pattern Library** | Common UI patterns are documented with composition rules |

### Coverage CI Integration

Add to CI pipeline to prevent coverage regression:

```yaml
# .github/workflows/design-system.yml
- name: Check DSS Coverage
  run: |
    dss coverage --json -d ./design-system > coverage.json
    SCORE=$(jq '.summary.overallScore' coverage.json)
    if (( $(echo "$SCORE < 90" | bc -l) )); then
      echo "DSS coverage dropped below 90%: $SCORE%"
      exit 1
    fi
```

## When DSS Makes Sense

DSS is not the right solution for everyone. Here's an honest assessment of when it's worth the investment.

### Good Fit Criteria

DSS is worthwhile when you have:

| Criteria | Why It Matters |
|----------|----------------|
| **10+ properties** | Manual synchronization becomes unsustainable |
| **Multiple output formats** | Same tokens → Tailwind, MkDocs, Marp, etc. |
| **AI-assisted development** | LLMs benefit from structured design context |
| **Code-first team** | Comfortable with JSON, CLI tools, CI/CD |
| **Documentation-as-code philosophy** | Want design decisions versioned and reviewable |

PlexusOne checks all five boxes, making DSS a strong fit.

### When Simpler Alternatives Are Better

| Situation | Better Alternative |
|-----------|-------------------|
| Single website | CSS custom properties in one file |
| 2-3 MkDocs sites | Shared `extra.css` in a gist or repo |
| No AI development | Figma + manual CSS |
| Simple brand (just colors/fonts) | Style Dictionary or a 50-line build script |
| Figma-first workflow | Figma Tokens plugin → Style Dictionary |

**The 80/20 rule**: A shell script that copies a CSS file to multiple repos solves 80% of the consistency problem with 20% of the effort.

### Honest Trade-offs

| Benefit | Cost |
|---------|------|
| Single source of truth | Another system to maintain |
| Comprehensive coverage (98.1%) | Significant upfront effort (11 foundation files, patterns, templates, governance, content) |
| LLM-ready context | Requires structured JSON authoring |
| Automated validation | Learning curve for `dss` CLI |
| Future-proof | May be over-engineering for simple needs |

### Questions to Ask Before Adopting

1. **Do you have 10+ properties?** If not, manual sync may be fine.
2. **Do you need multiple CSS formats?** If everything is React/Tailwind, a shared npm package suffices.
3. **Are AI assistants part of your dev workflow?** If not, the LLM context benefit doesn't apply.
4. **Does your team prefer code or visual tools?** Figma-first teams may resist JSON-based design specs.
5. **Do you need governance documentation?** Small teams communicating directly may not need formal deprecation policies.

### Why It Worked for PlexusOne

PlexusOne is a tooling company building AI agent infrastructure. The team:

- Already builds CLI tools and libraries
- Has 54+ properties across repos
- Uses AI assistants for code generation
- Values automation over manual processes
- Wants design decisions documented alongside code

For this profile, DSS provides clear ROI. For a small team with a single website, it would be over-engineering.

## Conclusion

DSS transformed PlexusOne's fragmented design landscape into a cohesive system. The key insight: **treat design as code**. Version it, generate outputs for each context, and distribute via packages.

The coverage journey from 29.4% to 98.1% demonstrates that a design system spec is not a one-time effort—it evolves as the system matures. The `dss coverage` command provides visibility into completeness and helps teams prioritize what to add next.

**Bottom line**: DSS is worth it if you have the scale (10+ properties), need multiple output formats, and want AI assistants to understand your design system. If you just need shared colors across a few sites, a simple CSS file is the better choice.
