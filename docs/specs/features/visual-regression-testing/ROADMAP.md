# Visual Regression Testing - Roadmap

**Version:** 1.0.0
**Status:** Draft
**Author:** PlexusOne
**Last Updated:** 2024-07-11

## Overview

This roadmap outlines the milestones for implementing visual regression testing in DSS.

## Timeline Summary

```
                    2024
    Week 1      Week 2      Week 3      Week 4      Week 5
    ┌───────────┬───────────┬───────────┬───────────┬───────────┐
    │           │           │           │           │           │
    │  M1: Core │  M2: Cap  │  M3: Exec │  M4: CLI  │  M5: Docs │
    │  Types    │  & Comp   │  & Base   │  & MCP    │  & Polish │
    │           │           │           │           │           │
    │  ████     │  ████     │  ████     │  ████     │  ████     │
    │           │           │           │           │           │
    └───────────┴───────────┴───────────┴───────────┴───────────┘
         │           │           │           │           │
         ▼           ▼           ▼           ▼           ▼
       v0.5.0-    v0.5.0-    v0.5.0-    v0.5.0-    v0.5.0
       alpha.1    alpha.2    alpha.3    rc.1

```

## Milestone 1: Core Types & Loader

**Target:** Week 1
**Tag:** `v0.5.0-alpha.1`

### Deliverables

| Deliverable | Status | Owner |
|-------------|--------|-------|
| `visual/types.go` - Core data types | ⬜ Pending | |
| `visual/loader.go` - Test definition loader | ⬜ Pending | |
| `visual/loader_test.go` - Unit tests | ⬜ Pending | |
| `schema/visual-test-suite.schema.json` | ⬜ Pending | |
| Example test definition | ⬜ Pending | |

### Acceptance Criteria

- [ ] VisualTest, VisualTestSuite, Result types defined
- [ ] YAML and JSON test definitions load correctly
- [ ] Defaults applied to individual tests
- [ ] JSON Schema validates test definitions
- [ ] Unit tests pass with 90% coverage

### Definition of Done

```bash
# All tests pass
go test -v ./sdk/go/visual/...

# Example loads successfully
go run ./cmd/dss visual-test --help
```

---

## Milestone 2: Capture & Comparison

**Target:** Week 2
**Tag:** `v0.5.0-alpha.2`

### Deliverables

| Deliverable | Status | Owner |
|-------------|--------|-------|
| `visual/capture.go` - W3Pilot MCP client | ⬜ Pending | |
| `visual/compare.go` - ImageMagick comparison | ⬜ Pending | |
| `visual/compare_go.go` - Pure Go fallback | ⬜ Pending | |
| Integration tests with w3pilot | ⬜ Pending | |

### Acceptance Criteria

- [ ] W3Pilot subprocess starts and connects via MCP
- [ ] Screenshots captured for page and elements
- [ ] Viewport configuration works
- [ ] Stabilization (wait for selector, wait ms) works
- [ ] ImageMagick comparison produces diff percentage
- [ ] Diff images generated with visual highlights
- [ ] Pure Go fallback works when ImageMagick unavailable

### Definition of Done

```bash
# Capture test (requires w3pilot)
go test -v ./sdk/go/visual/ -run TestCapture

# Comparison test
go test -v ./sdk/go/visual/ -run TestCompare
```

---

## Milestone 3: Executor & Baseline

**Target:** Week 3
**Tag:** `v0.5.0-alpha.3`

### Deliverables

| Deliverable | Status | Owner |
|-------------|--------|-------|
| `visual/baseline.go` - Baseline management | ⬜ Pending | |
| `visual/executor.go` - Test executor | ⬜ Pending | |
| `visual/service.go` - Public service API | ⬜ Pending | |
| Parallel execution support | ⬜ Pending | |

### Acceptance Criteria

- [ ] Baselines stored by version with manifest
- [ ] Baseline checksums calculated and verified
- [ ] Parallel test execution with worker pool
- [ ] Each worker manages its own w3pilot instance
- [ ] Results aggregated into report structure
- [ ] Error handling for individual test failures

### Definition of Done

```bash
# Generate baseline (requires running storybook)
go run ./cmd/dss visual-baseline generate --version v1.0.0

# Run tests
go run ./cmd/dss visual-test --baseline v1.0.0
```

---

## Milestone 4: CLI & MCP

**Target:** Week 4
**Tag:** `v0.5.0-rc.1`

### Deliverables

| Deliverable | Status | Owner |
|-------------|--------|-------|
| `cmd/dss/cmd/visual_test.go` - Test command | ⬜ Pending | |
| `cmd/dss/cmd/visual_baseline.go` - Baseline commands | ⬜ Pending | |
| `skills/designsystem/tools_visual.go` - MCP tools | ⬜ Pending | |
| Compliance integration | ⬜ Pending | |

### CLI Commands

| Command | Status |
|---------|--------|
| `dss visual-test` | ⬜ Pending |
| `dss visual-test --baseline <version>` | ⬜ Pending |
| `dss visual-test --test <id>` | ⬜ Pending |
| `dss visual-test --json` | ⬜ Pending |
| `dss visual-baseline generate --version <v>` | ⬜ Pending |
| `dss visual-baseline update --test <id>` | ⬜ Pending |
| `dss visual-baseline list` | ⬜ Pending |

### MCP Tools

| Tool | Status |
|------|--------|
| `visual_test` | ⬜ Pending |
| `visual_baseline_generate` | ⬜ Pending |
| `visual_baseline_update` | ⬜ Pending |
| `visual_test_single` | ⬜ Pending |

### Acceptance Criteria

- [ ] All CLI commands functional
- [ ] JSON output format for CI integration
- [ ] All MCP tools exposed via dss-mcp
- [ ] Visual tests included in compliance report
- [ ] Release gate can block on visual failures

### Definition of Done

```bash
# CLI works
dss visual-test --baseline v1.0.0 --json

# MCP tools work
echo '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"visual_test","arguments":{}}}' | dss-mcp --spec ./examples/
```

---

## Milestone 5: Documentation & Polish

**Target:** Week 5
**Tag:** `v0.5.0`

### Deliverables

| Deliverable | Status | Owner |
|-------------|--------|-------|
| `docs/visual-testing.md` - User guide | ⬜ Pending | |
| `docs/mcp-server.md` - Update with visual tools | ⬜ Pending | |
| Example visual test suite | ⬜ Pending | |
| CI/CD example workflow | ⬜ Pending | |
| Release notes | ⬜ Pending | |

### Documentation

| Page | Status |
|------|--------|
| Visual Testing Guide | ⬜ Pending |
| CLI Reference (visual commands) | ⬜ Pending |
| MCP Tools (visual tools) | ⬜ Pending |
| CI/CD Integration | ⬜ Pending |
| Troubleshooting | ⬜ Pending |

### Acceptance Criteria

- [ ] User documentation covers all features
- [ ] Example test suite with 10+ component tests
- [ ] GitHub Actions workflow example
- [ ] 90% test coverage
- [ ] Agentic maturity updated to L4

### Definition of Done

```bash
# All tests pass
go test -v ./...

# Documentation builds
cd docs && mkdocs build

# Coverage meets target
go test -cover ./sdk/go/visual/... | grep "coverage: 90"
```

---

## Future Milestones (v0.6.0+)

### M6: Cross-Browser Testing

**Target:** v0.6.0

| Feature | Description |
|---------|-------------|
| Firefox support | w3pilot --browser firefox |
| Safari support | w3pilot --browser safari (webkit) |
| Browser matrix | Run tests across all browsers |
| Browser-specific baselines | Separate baselines per browser |

### M7: Advanced Features

**Target:** v0.7.0

| Feature | Description |
|---------|-------------|
| Animation testing | Capture key animation frames |
| Interaction states | Hover, focus, active states |
| Dark mode testing | Theme variant baselines |
| Responsive breakpoints | Auto-test at breakpoints |

### M8: CI/CD Integration

**Target:** v0.8.0

| Feature | Description |
|---------|-------------|
| GitHub Action | `dss-visual-action` |
| PR comments | Auto-comment with diff images |
| Artifact upload | Diff images as PR artifacts |
| Approval workflow | Approve visual changes in PR |

### M9: Visual AI

**Target:** v0.9.0

| Feature | Description |
|---------|-------------|
| Perceptual hashing | Ignore minor pixel differences |
| Layout detection | Detect layout shifts |
| Component recognition | Auto-detect component boundaries |
| AI-powered diff explanation | "Button border changed from 1px to 2px" |

---

## Version History

| Version | Date | Description |
|---------|------|-------------|
| v0.5.0-alpha.1 | TBD | Core types and loader |
| v0.5.0-alpha.2 | TBD | Capture and comparison |
| v0.5.0-alpha.3 | TBD | Executor and baselines |
| v0.5.0-rc.1 | TBD | CLI and MCP tools |
| v0.5.0 | TBD | GA release |

---

## Dependencies

### Required for v0.5.0

| Dependency | Version | Status |
|------------|---------|--------|
| w3pilot | >=0.7.0 | ✅ Available |
| ImageMagick | >=7.0 | ✅ Available |
| Chrome/Chromium | >=120 | ✅ Available |

### Planned Dependencies

| Dependency | Version | Milestone |
|------------|---------|-----------|
| go-imagediff | latest | M2 (fallback) |
| Firefox | >=120 | M6 |
| WebKit | latest | M6 |

---

## Risks & Mitigations

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| w3pilot API changes | High | Low | Pin version, add integration tests |
| ImageMagick unavailable in CI | Medium | Medium | Pure Go fallback |
| Flaky visual tests | Medium | High | Stabilization, retry logic, thresholds |
| Large baseline storage | Low | Medium | Compression, pruning, git-lfs |
| Slow test execution | Medium | Medium | Parallelization, test filtering |

---

## Success Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Test suite execution | <5 min | CI timing |
| False positive rate | <5% | Manual review |
| Coverage improvement | L3 → L4 | Agentic maturity |
| Adoption | 100% components | Baseline coverage |

---

## Communication Plan

### Milestone Announcements

| Milestone | Channel |
|-----------|---------|
| Alpha releases | GitHub Releases (pre-release) |
| RC releases | GitHub Releases + Discord |
| GA release | GitHub Releases + Blog post |

### Documentation Updates

| Content | Location |
|---------|----------|
| Release notes | `CHANGELOG.md`, GitHub Release |
| User guide | `docs/visual-testing.md` |
| API reference | GoDoc |
| MCP tools | `docs/mcp-server.md` |

---

## Appendix: Test Suite Example

```yaml
# visual-tests/visual-tests.yaml
version: "1.0"
name: "Design System Visual Tests"
baseUrl: "http://localhost:6006/iframe.html"

defaults:
  viewports:
    - name: desktop
      width: 1280
      height: 800
    - name: mobile
      width: 375
      height: 667
  threshold: 0.001
  stabilization:
    waitMs: 100
    disableAnimations: true

tests:
  - id: button-primary
    component: button
    variant: primary
    url: "?id=button--primary"
    selector: "#storybook-root > *"

  - id: button-secondary
    component: button
    variant: secondary
    url: "?id=button--secondary"
    selector: "#storybook-root > *"

  - id: button-destructive
    component: button
    variant: destructive
    url: "?id=button--destructive"
    selector: "#storybook-root > *"

  - id: card-default
    component: card
    url: "?id=card--default"
    selector: "#storybook-root > *"
    viewports:
      - name: desktop
        width: 1280
        height: 800

  - id: input-default
    component: input
    url: "?id=input--default"
    selector: "#storybook-root > *"
    threshold: 0.005  # Higher tolerance for text rendering

  - id: modal-open
    component: modal
    url: "?id=modal--open"
    stabilization:
      waitForSelector: "[data-state='open']"
      waitMs: 300  # Wait for animation
```

---

## Related Documents

- [PRD.md](PRD.md) - Product requirements
- [TRD.md](TRD.md) - Technical requirements
- [PLAN.md](PLAN.md) - Implementation plan
