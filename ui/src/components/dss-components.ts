import { LitElement, html, css } from 'lit';
import { customElement, property, state } from 'lit/decorators.js';
import { sharedStyles } from '../styles.js';
import type { Component } from '../types.js';

@customElement('dss-components')
export class DssComponents extends LitElement {
  static styles = [
    sharedStyles,
    css`
      .components-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
        gap: 1rem;
      }

      .component-card {
        background: var(--dss-bg-tertiary);
        border: 1px solid var(--dss-border);
        border-radius: var(--dss-radius-md);
        overflow: hidden;
      }

      .component-header {
        padding: 1rem;
        border-bottom: 1px solid var(--dss-border);
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
      }

      .component-title {
        font-size: 1rem;
        font-weight: 600;
        margin: 0;
      }

      .component-id {
        font-size: 0.7rem;
        font-family: var(--dss-font-mono);
        color: var(--dss-text-muted);
        margin: 0.25rem 0 0;
      }

      .component-category {
        font-size: 0.7rem;
        padding: 0.125rem 0.5rem;
        background: var(--dss-accent);
        color: var(--dss-bg-primary);
        border-radius: var(--dss-radius-sm);
      }

      .component-body {
        padding: 1rem;
      }

      .component-description {
        font-size: 0.875rem;
        color: var(--dss-text-secondary);
        margin: 0 0 1rem;
        line-height: 1.5;
      }

      .component-section {
        margin-bottom: 1rem;
      }

      .component-section:last-child {
        margin-bottom: 0;
      }

      .component-section-title {
        font-size: 0.75rem;
        font-weight: 600;
        color: var(--dss-text-muted);
        text-transform: uppercase;
        letter-spacing: 0.05em;
        margin: 0 0 0.5rem;
      }

      .variants-list {
        display: flex;
        flex-wrap: wrap;
        gap: 0.5rem;
      }

      .variant-tag {
        padding: 0.25rem 0.5rem;
        font-size: 0.75rem;
        background: var(--dss-bg-secondary);
        border-radius: var(--dss-radius-sm);
        border: 1px solid var(--dss-border);
      }

      .variant-tag--default {
        border-color: var(--dss-accent);
        color: var(--dss-accent);
      }

      .states-list {
        display: flex;
        flex-wrap: wrap;
        gap: 0.375rem;
      }

      .state-tag {
        padding: 0.125rem 0.375rem;
        font-size: 0.7rem;
        background: var(--dss-bg-secondary);
        border-radius: 2px;
        color: var(--dss-text-secondary);
      }

      .props-table {
        width: 100%;
        font-size: 0.75rem;
        border-collapse: collapse;
      }

      .props-table th,
      .props-table td {
        text-align: left;
        padding: 0.375rem 0.5rem;
        border-bottom: 1px solid var(--dss-border);
      }

      .props-table th {
        font-weight: 600;
        color: var(--dss-text-muted);
        text-transform: uppercase;
        font-size: 0.65rem;
        letter-spacing: 0.05em;
      }

      .props-table td {
        color: var(--dss-text-secondary);
      }

      .props-table .prop-name {
        font-family: var(--dss-font-mono);
        color: var(--dss-text-primary);
      }

      .props-table .prop-type {
        font-family: var(--dss-font-mono);
        color: var(--dss-accent);
        font-size: 0.7rem;
      }

      .props-table .prop-required {
        color: var(--dss-warning);
      }

      .expand-btn {
        background: none;
        border: none;
        color: var(--dss-accent);
        font-size: 0.75rem;
        cursor: pointer;
        padding: 0.25rem 0;
      }

      .expand-btn:hover {
        text-decoration: underline;
      }

      .a11y-info {
        font-size: 0.75rem;
        color: var(--dss-text-secondary);
      }

      .a11y-role {
        font-family: var(--dss-font-mono);
        background: var(--dss-bg-secondary);
        padding: 0.125rem 0.375rem;
        border-radius: 2px;
      }

      .keyboard-list {
        margin: 0.5rem 0 0;
        padding: 0;
        list-style: none;
      }

      .keyboard-item {
        display: flex;
        gap: 0.5rem;
        font-size: 0.7rem;
        margin-bottom: 0.25rem;
      }

      .keyboard-key {
        font-family: var(--dss-font-mono);
        background: var(--dss-bg-secondary);
        padding: 0.125rem 0.375rem;
        border-radius: 2px;
        min-width: 50px;
      }

      .theming-tokens {
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
      }

      .theming-token {
        display: flex;
        justify-content: space-between;
        font-size: 0.7rem;
        padding: 0.25rem 0.5rem;
        background: var(--dss-bg-secondary);
        border-radius: 2px;
      }

      .theming-token-name {
        font-family: var(--dss-font-mono);
      }

      .theming-token-semantic {
        color: var(--dss-text-muted);
      }

      .filter-bar {
        display: flex;
        gap: 0.5rem;
        margin-bottom: 1rem;
        flex-wrap: wrap;
      }

      .filter-btn {
        padding: 0.375rem 0.75rem;
        font-size: 0.75rem;
        border: 1px solid var(--dss-border);
        border-radius: var(--dss-radius-sm);
        background: var(--dss-bg-tertiary);
        color: var(--dss-text-secondary);
        cursor: pointer;
        transition: all 0.2s;
      }

      .filter-btn:hover {
        border-color: var(--dss-border-hover);
      }

      .filter-btn--active {
        background: var(--dss-accent);
        border-color: var(--dss-accent);
        color: var(--dss-bg-primary);
      }
    `,
  ];

  @property({ type: Array })
  components: Component[] = [];

  @property({ type: Boolean })
  editable = false;

  @state()
  private _expandedComponents: Set<string> = new Set();

  @state()
  private _categoryFilter: string | null = null;

  private get _categories(): string[] {
    const categories = new Set<string>();
    for (const comp of this.components) {
      if (comp.category) {
        categories.add(comp.category);
      }
    }
    return Array.from(categories).sort();
  }

  private get _filteredComponents(): Component[] {
    if (!this._categoryFilter) return this.components;
    return this.components.filter((c) => c.category === this._categoryFilter);
  }

  private _toggleExpanded(id: string) {
    if (this._expandedComponents.has(id)) {
      this._expandedComponents.delete(id);
    } else {
      this._expandedComponents.add(id);
    }
    this.requestUpdate();
  }

  render() {
    if (!this.components?.length) {
      return html`
        <div class="section">
          <div class="section-header">
            <h2 class="section-title">Components</h2>
          </div>
          <div class="empty-state">No components defined</div>
        </div>
      `;
    }

    return html`
      <div class="section">
        <div class="section-header">
          <h2 class="section-title">Components</h2>
          <span class="badge">${this.components.length} components</span>
        </div>

        ${this._categories.length > 1
          ? html`
              <div class="filter-bar">
                <button
                  class="filter-btn ${!this._categoryFilter
                    ? 'filter-btn--active'
                    : ''}"
                  @click=${() => (this._categoryFilter = null)}
                >
                  All
                </button>
                ${this._categories.map(
                  (cat) => html`
                    <button
                      class="filter-btn ${this._categoryFilter === cat
                        ? 'filter-btn--active'
                        : ''}"
                      @click=${() => (this._categoryFilter = cat)}
                    >
                      ${cat}
                    </button>
                  `
                )}
              </div>
            `
          : null}

        <div class="components-grid">
          ${this._filteredComponents.map((comp) =>
            this._renderComponent(comp)
          )}
        </div>
      </div>
    `;
  }

  private _renderComponent(comp: Component) {
    const isExpanded = this._expandedComponents.has(comp.id);

    return html`
      <div class="component-card">
        <div class="component-header">
          <div>
            <h3 class="component-title">${comp.name}</h3>
            <p class="component-id">${comp.id}</p>
          </div>
          ${comp.category
            ? html`<span class="component-category">${comp.category}</span>`
            : null}
        </div>

        <div class="component-body">
          ${comp.description
            ? html`<p class="component-description">${comp.description}</p>`
            : null}

          <!-- Variants -->
          ${comp.variants?.length
            ? html`
                <div class="component-section">
                  <h4 class="component-section-title">Variants</h4>
                  <div class="variants-list">
                    ${comp.variants.map(
                      (v) => html`
                        <span
                          class="variant-tag ${v.isDefault
                            ? 'variant-tag--default'
                            : ''}"
                          title="${v.description || ''}"
                        >
                          ${v.name}${v.isDefault ? ' (default)' : ''}
                        </span>
                      `
                    )}
                  </div>
                </div>
              `
            : null}

          <!-- States -->
          ${comp.states?.length
            ? html`
                <div class="component-section">
                  <h4 class="component-section-title">States</h4>
                  <div class="states-list">
                    ${comp.states.map(
                      (s) => html`
                        <span class="state-tag" title="${s.description || ''}"
                          >${s.id}</span
                        >
                      `
                    )}
                  </div>
                </div>
              `
            : null}

          <!-- Props (collapsible) -->
          ${comp.props?.length
            ? html`
                <div class="component-section">
                  <h4 class="component-section-title">
                    Props
                    <button
                      class="expand-btn"
                      @click=${() => this._toggleExpanded(comp.id)}
                    >
                      ${isExpanded ? 'Hide' : 'Show'} ${comp.props.length} props
                    </button>
                  </h4>
                  ${isExpanded
                    ? html`
                        <table class="props-table">
                          <thead>
                            <tr>
                              <th>Name</th>
                              <th>Type</th>
                              <th>Default</th>
                            </tr>
                          </thead>
                          <tbody>
                            ${comp.props.map(
                              (p) => html`
                                <tr>
                                  <td>
                                    <span class="prop-name">${p.name}</span>
                                    ${p.required
                                      ? html`<span class="prop-required"
                                          >*</span
                                        >`
                                      : null}
                                  </td>
                                  <td class="prop-type">${p.type}</td>
                                  <td>${p.default ?? '-'}</td>
                                </tr>
                              `
                            )}
                          </tbody>
                        </table>
                      `
                    : null}
                </div>
              `
            : null}

          <!-- Theming Contract -->
          ${comp.themingContract?.tokens?.length
            ? html`
                <div class="component-section">
                  <h4 class="component-section-title">
                    Theming Contract
                    <span class="badge" style="margin-left: 0.5rem"
                      >${comp.themingContract.prefix}</span
                    >
                  </h4>
                  <div class="theming-tokens">
                    ${comp.themingContract.tokens.slice(0, 4).map(
                      (t) => html`
                        <div class="theming-token">
                          <span class="theming-token-name"
                            >${t.cssProperty}</span
                          >
                          ${t.semantic
                            ? html`<span class="theming-token-semantic"
                                >${t.semantic}</span
                              >`
                            : null}
                        </div>
                      `
                    )}
                    ${comp.themingContract.tokens.length > 4
                      ? html`<span
                          style="font-size: 0.7rem; color: var(--dss-text-muted)"
                          >+${comp.themingContract.tokens.length - 4} more
                          tokens</span
                        >`
                      : null}
                  </div>
                </div>
              `
            : null}

          <!-- Accessibility -->
          ${comp.accessibility
            ? html`
                <div class="component-section">
                  <h4 class="component-section-title">Accessibility</h4>
                  <div class="a11y-info">
                    ${comp.accessibility.role
                      ? html`Role:
                          <span class="a11y-role"
                            >${comp.accessibility.role}</span
                          >`
                      : null}
                    ${comp.accessibility.keyboardSupport?.length
                      ? html`
                          <ul class="keyboard-list">
                            ${comp.accessibility.keyboardSupport
                              .slice(0, 3)
                              .map(
                                (k) => html`
                                  <li class="keyboard-item">
                                    <span class="keyboard-key">${k.key}</span>
                                    <span>${k.action}</span>
                                  </li>
                                `
                              )}
                          </ul>
                        `
                      : null}
                  </div>
                </div>
              `
            : null}
        </div>
      </div>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'dss-components': DssComponents;
  }
}
