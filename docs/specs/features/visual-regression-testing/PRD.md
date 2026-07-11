# Visual Regression Testing - Product Requirements Document

**Version:** 1.0.0
**Status:** Draft
**Author:** PlexusOne
**Last Updated:** 2024-07-11

## Executive Summary

Visual regression testing enables automated detection of unintended visual changes in UI components. This feature integrates w3pilot's screenshot capabilities with DSS to provide a complete validation loop for design system compliance—covering not just token/code compliance but actual rendered output.

## Problem Statement

### Current Gap

DSS currently validates:

- ✅ Token usage (colors, spacing)
- ✅ Accessibility attributes (alt, aria-label)
- ✅ Component variant compliance
- ✅ Anti-pattern detection

DSS cannot validate:

- ❌ Visual correctness ("does it look right?")
- ❌ Rendering consistency across browsers
- ❌ Layout and spacing accuracy
- ❌ Responsive behavior
- ❌ Visual regressions between versions

### Impact

Without visual validation:

1. **Manual QA burden** - Humans must visually inspect every component change
2. **Regression risk** - Visual bugs slip through code review
3. **Inconsistent releases** - Different reviewers apply different standards
4. **Slow iteration** - Fear of breaking visuals slows development

### Agentic Maturity Impact

Per the [Agentic Maturity Assessment](../../../agentic-maturity.md):

| Phase | Current | With Visual Testing |
|-------|---------|---------------------|
| Validation | L3 | L4 |
| Overall | L2.8 | L3.2 |

Visual testing addresses the #1 gap identified in the maturity assessment.

## Goals

### Primary Goals

1. **Automated visual regression detection** - Catch visual changes before release
2. **Baseline management** - Version-controlled reference images per release
3. **CI/CD integration** - Block PRs with unreviewed visual changes
4. **Agent-accessible validation** - MCP tools for AI agent workflows

### Secondary Goals

1. **Cross-browser testing** - Validate Chrome, Firefox, Safari rendering
2. **Responsive testing** - Multiple viewport sizes per component
3. **Component isolation** - Test components in isolation from page context
4. **Performance** - Complete test suite in <5 minutes

### Non-Goals (v1.0)

1. **Animation testing** - Static snapshots only
2. **Interaction testing** - Visual state, not behavior
3. **Full page testing** - Component-level only
4. **Figma integration** - Code-to-code comparison only

## User Personas

### 1. Design System Maintainer

**Needs:**

- Generate baselines for new releases
- Review visual diffs before approving PRs
- Update baselines when intentional changes are made

**Workflow:**

```
1. Merge component changes
2. CI runs visual tests
3. Review diff report
4. Approve or request changes
5. Update baseline if intentional
```

### 2. AI Coding Agent (Claude, Copilot)

**Needs:**

- Validate generated code produces correct visuals
- Identify which components changed visually
- Auto-fix or report visual regressions

**Workflow:**

```
1. Generate component code
2. Call visual_test MCP tool
3. If regression detected, analyze diff
4. Attempt fix or escalate to human
```

### 3. Frontend Developer

**Needs:**

- Run visual tests locally before pushing
- Debug visual differences
- Understand what changed and why

**Workflow:**

```
1. Make component changes
2. Run dss visual-test locally
3. Review diff images
4. Fix unintended changes
5. Push with confidence
```

## Requirements

### Functional Requirements

#### FR-1: Test Definition

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-1.1 | Define visual tests in JSON/YAML format | P0 |
| FR-1.2 | Support component-level test targeting | P0 |
| FR-1.3 | Support multiple viewports per test | P0 |
| FR-1.4 | Support custom thresholds per test | P1 |
| FR-1.5 | Auto-generate tests from DSS component spec | P2 |

#### FR-2: Baseline Management

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-2.1 | Generate baseline screenshots for release | P0 |
| FR-2.2 | Store baselines by version tag | P0 |
| FR-2.3 | Support baseline update workflow | P0 |
| FR-2.4 | Validate baseline completeness | P1 |
| FR-2.5 | Baseline pruning for deleted tests | P2 |

#### FR-3: Test Execution

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-3.1 | Capture screenshots via w3pilot | P0 |
| FR-3.2 | Compare against baseline with diff generation | P0 |
| FR-3.3 | Support parallel test execution | P0 |
| FR-3.4 | Support single-test debugging mode | P1 |
| FR-3.5 | Retry flaky tests with stabilization | P1 |

#### FR-4: Reporting

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-4.1 | Generate JSON report with pass/fail status | P0 |
| FR-4.2 | Generate diff images for failures | P0 |
| FR-4.3 | Generate HTML report for human review | P1 |
| FR-4.4 | Integrate with compliance report | P1 |
| FR-4.5 | GitHub PR comment with diff summary | P2 |

#### FR-5: MCP Integration

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-5.1 | `visual_test` tool for running tests | P0 |
| FR-5.2 | `visual_baseline_generate` tool | P0 |
| FR-5.3 | `visual_baseline_update` tool | P1 |
| FR-5.4 | `visual_test_single` for debugging | P1 |
| FR-5.5 | Include visual status in compliance report | P1 |

### Non-Functional Requirements

#### NFR-1: Performance

| ID | Requirement | Target |
|----|-------------|--------|
| NFR-1.1 | Single test execution | <2 seconds |
| NFR-1.2 | Full suite (100 tests) | <5 minutes |
| NFR-1.3 | Image comparison | <100ms per pair |
| NFR-1.4 | Baseline generation | <10 minutes |

#### NFR-2: Accuracy

| ID | Requirement | Target |
|----|-------------|--------|
| NFR-2.1 | False positive rate | <1% |
| NFR-2.2 | False negative rate | <0.1% |
| NFR-2.3 | Pixel threshold default | 0.1% (configurable) |
| NFR-2.4 | Anti-aliasing tolerance | Built-in |

#### NFR-3: Storage

| ID | Requirement | Target |
|----|-------------|--------|
| NFR-3.1 | Screenshot size | <200 KB average |
| NFR-3.2 | Baseline per version | <50 MB |
| NFR-3.3 | Diff image size | <300 KB |

## Success Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Visual regressions caught | >95% | Bugs found in CI vs production |
| False positive rate | <5% | Tests requiring manual override |
| Test suite run time | <5 min | CI timing |
| Developer adoption | >80% | Teams using visual tests |
| Baseline coverage | 100% | Components with baselines |

## User Stories

### Epic: Visual Test Definition

**US-1:** As a design system maintainer, I want to define visual tests in a declarative format so that tests are version-controlled and reviewable.

**Acceptance Criteria:**

- Tests defined in `visual-tests/` directory
- JSON/YAML format with schema validation
- Component, variant, viewport, and threshold fields
- Schema documented and validated by `dss lint-spec`

### Epic: Baseline Management

**US-2:** As a release manager, I want to generate baselines for a tagged release so that future changes are compared against a known-good state.

**Acceptance Criteria:**

- `dss visual-baseline generate --version v1.0.0`
- Baselines stored in `baselines/{version}/`
- Manifest includes test metadata and generation timestamp
- Baselines can be regenerated idempotently

**US-3:** As a developer, I want to update baselines when I make intentional visual changes so that CI passes for approved changes.

**Acceptance Criteria:**

- `dss visual-baseline update --test button-primary`
- Requires explicit confirmation (not automatic)
- Audit log of baseline updates
- Works for single test or full suite

### Epic: Test Execution

**US-4:** As a CI pipeline, I want to run visual tests on every PR so that regressions are caught before merge.

**Acceptance Criteria:**

- `dss visual-test --baseline v1.0.0`
- Exit code 1 if any test fails threshold
- JSON report for programmatic consumption
- Diff images as CI artifacts

**US-5:** As a developer, I want to run a single visual test locally so that I can debug a specific failure.

**Acceptance Criteria:**

- `dss visual-test --test button-primary --interactive`
- Opens diff viewer or generates local diff image
- Shows pixel difference percentage
- Can compare against any baseline version

### Epic: MCP Integration

**US-6:** As an AI coding agent, I want to validate my generated code visually so that I can catch rendering issues before human review.

**Acceptance Criteria:**

- `visual_test` MCP tool available
- Returns pass/fail with diff percentage
- Can request diff image for analysis
- Integrated into compliance report

## Appendix

### A. Test Definition Schema

```json
{
  "$schema": "https://dss.dev/schema/visual-test.json",
  "id": "button-primary-desktop",
  "component": "button",
  "variant": "primary",
  "url": "http://localhost:6006/iframe.html?id=button--primary",
  "selector": "#storybook-root > *",
  "viewport": {
    "width": 1280,
    "height": 800
  },
  "threshold": 0.001,
  "stabilization": {
    "waitForSelector": "[data-ready]",
    "waitMs": 100
  }
}
```

### B. Baseline Directory Structure

```
baselines/
├── v1.0.0/
│   ├── manifest.json
│   ├── button-primary-desktop.png
│   ├── button-primary-mobile.png
│   ├── button-secondary-desktop.png
│   └── ...
├── v1.1.0/
│   └── ...
└── latest -> v1.1.0
```

### C. Report Schema

```json
{
  "timestamp": "2024-07-11T10:30:00Z",
  "baselineVersion": "v1.0.0",
  "summary": {
    "total": 100,
    "passed": 98,
    "failed": 2,
    "skipped": 0
  },
  "results": [
    {
      "testId": "button-primary-desktop",
      "passed": false,
      "diffPercent": 0.15,
      "threshold": 0.001,
      "baselinePath": "baselines/v1.0.0/button-primary-desktop.png",
      "actualPath": "results/button-primary-desktop.png",
      "diffPath": "results/button-primary-desktop.diff.png"
    }
  ]
}
```

### D. Related Documents

- [Technical Requirements Document](TRD.md)
- [Implementation Plan](PLAN.md)
- [Roadmap](ROADMAP.md)
- [Agentic Maturity Assessment](../../../agentic-maturity.md)
