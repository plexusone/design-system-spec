import { LitElement, html, css } from 'lit';
import { customElement, property, state } from 'lit/decorators.js';
import { sharedStyles } from '../styles.js';
import type { ThemeBindings, Component, ColorToken } from '../types.js';

@customElement('dss-contracts')
export class DssContracts extends LitElement {
  static styles = [
    sharedStyles,
    css`
      .bindings-list {
        display: flex;
        flex-direction: column;
        gap: 1rem;
      }

      .binding-card {
        background: var(--dss-bg-tertiary);
        border: 1px solid var(--dss-border);
        border-radius: var(--dss-radius-md);
        overflow: hidden;
      }

      .binding-header {
        padding: 1rem;
        border-bottom: 1px solid var(--dss-border);
        display: flex;
        justify-content: space-between;
        align-items: center;
      }

      .binding-title {
        font-size: 1rem;
        font-weight: 600;
        margin: 0;
      }

      .binding-meta {
        display: flex;
        gap: 0.5rem;
      }

      .binding-badge {
        font-size: 0.7rem;
        padding: 0.125rem 0.5rem;
        border-radius: var(--dss-radius-sm);
        background: var(--dss-bg-secondary);
        color: var(--dss-text-secondary);
      }

      .binding-badge--strategy {
        background: var(--dss-accent);
        color: var(--dss-bg-primary);
      }

      .binding-badge--mode {
        background: var(--dss-bg-hover);
      }

      .mappings-table {
        width: 100%;
        border-collapse: collapse;
        font-size: 0.8rem;
      }

      .mappings-table th,
      .mappings-table td {
        text-align: left;
        padding: 0.75rem 1rem;
        border-bottom: 1px solid var(--dss-border);
      }

      .mappings-table th {
        font-weight: 600;
        color: var(--dss-text-muted);
        text-transform: uppercase;
        font-size: 0.65rem;
        letter-spacing: 0.05em;
        background: var(--dss-bg-secondary);
      }

      .mappings-table tbody tr:last-child td {
        border-bottom: none;
      }

      .mappings-table tbody tr:hover {
        background: var(--dss-bg-hover);
      }

      .mapping-from {
        font-family: var(--dss-font-mono);
        color: var(--dss-accent);
      }

      .mapping-to {
        font-family: var(--dss-font-mono);
        color: var(--dss-text-primary);
      }

      .mapping-arrow {
        color: var(--dss-text-muted);
        padding: 0 0.5rem;
      }

      .mapping-preview {
        display: flex;
        align-items: center;
        gap: 0.5rem;
      }

      .color-preview {
        width: 20px;
        height: 20px;
        border-radius: 4px;
        border: 1px solid var(--dss-border);
      }

      .mapping-transform {
        font-size: 0.7rem;
        font-family: var(--dss-font-mono);
        color: var(--dss-warning);
        background: var(--dss-bg-secondary);
        padding: 0.125rem 0.375rem;
        border-radius: 2px;
      }

      .empty-mappings {
        padding: 2rem;
        text-align: center;
        color: var(--dss-text-muted);
        font-size: 0.875rem;
      }

      .spec-url {
        font-size: 0.7rem;
        font-family: var(--dss-font-mono);
        color: var(--dss-text-muted);
        margin-top: 0.25rem;
      }

      .generate-section {
        padding: 1rem;
        border-top: 1px solid var(--dss-border);
        background: var(--dss-bg-secondary);
      }

      .generate-title {
        font-size: 0.75rem;
        font-weight: 600;
        color: var(--dss-text-muted);
        margin: 0 0 0.75rem;
        text-transform: uppercase;
        letter-spacing: 0.05em;
      }

      .generate-output {
        background: var(--dss-bg-primary);
        border: 1px solid var(--dss-border);
        border-radius: var(--dss-radius-sm);
        padding: 0.75rem;
        font-family: var(--dss-font-mono);
        font-size: 0.7rem;
        white-space: pre-wrap;
        overflow-x: auto;
        max-height: 200px;
        overflow-y: auto;
      }

      .no-bindings {
        text-align: center;
        padding: 3rem;
      }

      .no-bindings-title {
        font-size: 1rem;
        font-weight: 600;
        margin: 0 0 0.5rem;
        color: var(--dss-text-secondary);
      }

      .no-bindings-text {
        font-size: 0.875rem;
        color: var(--dss-text-muted);
        margin: 0;
      }
    `,
  ];

  @property({ type: Array })
  themeBindings: ThemeBindings[] = [];

  @property({ type: Array })
  components: Component[] = [];

  @property({ type: Array })
  colors: ColorToken[] = [];

  @property({ type: Boolean })
  editable = false;

  @state()
  private _expandedBindings: Set<number> = new Set();

  private _toggleExpanded(index: number) {
    if (this._expandedBindings.has(index)) {
      this._expandedBindings.delete(index);
    } else {
      this._expandedBindings.add(index);
    }
    this.requestUpdate();
  }

  private _getColorValue(tokenId: string): string | null {
    const color = this.colors.find((c) => c.id === tokenId);
    return color?.value || null;
  }

  private _getComponentContract(componentId: string) {
    const comp = this.components.find((c) => c.id === componentId);
    return comp?.themingContract || null;
  }

  private _generateCSS(binding: ThemeBindings): string {
    const lines: string[] = [];
    const contract = this._getComponentContract(binding.component);

    lines.push(`/* Theme bindings for ${binding.component} */`);

    let selector = ':root';
    if (binding.themeMode === 'dark') {
      selector = ":root[data-theme='dark'], .dark";
    } else if (binding.themeMode === 'light') {
      selector = ":root[data-theme='light'], .light";
    }

    lines.push(`${selector} {`);

    for (const mapping of binding.mappings) {
      const colorValue = this._getColorValue(mapping.from);
      const token = contract?.tokens.find((t) => t.id === mapping.to);
      const cssProperty = token?.cssProperty || `--${binding.component}-${mapping.to}`;

      let value = colorValue || `var(--${mapping.from})`;
      if (mapping.transform) {
        value = `${mapping.transform}(${value})`;
      }

      lines.push(`  ${cssProperty}: ${value};`);
    }

    lines.push('}');

    return lines.join('\n');
  }

  render() {
    if (!this.themeBindings?.length) {
      return html`
        <div class="section">
          <div class="section-header">
            <h2 class="section-title">Theme Bindings</h2>
          </div>
          <div class="no-bindings">
            <p class="no-bindings-title">No theme bindings defined</p>
            <p class="no-bindings-text">
              Theme bindings map your design tokens to component theming
              contracts.
            </p>
          </div>
        </div>
      `;
    }

    return html`
      <div class="section">
        <div class="section-header">
          <h2 class="section-title">Theme Bindings</h2>
          <span class="badge">${this.themeBindings.length} bindings</span>
        </div>

        <div class="bindings-list">
          ${this.themeBindings.map((binding, index) =>
            this._renderBinding(binding, index)
          )}
        </div>
      </div>
    `;
  }

  private _renderBinding(binding: ThemeBindings, index: number) {
    const isExpanded = this._expandedBindings.has(index);
    const contract = this._getComponentContract(binding.component);

    return html`
      <div class="binding-card">
        <div class="binding-header">
          <div>
            <h3 class="binding-title">${binding.component}</h3>
            ${binding.specUrl
              ? html`<p class="spec-url">${binding.specUrl}</p>`
              : null}
          </div>
          <div class="binding-meta">
            ${binding.strategy
              ? html`<span class="binding-badge binding-badge--strategy"
                  >${binding.strategy}</span
                >`
              : null}
            ${binding.themeMode
              ? html`<span class="binding-badge binding-badge--mode"
                  >${binding.themeMode}</span
                >`
              : null}
            <span class="binding-badge"
              >${binding.mappings.length} mappings</span
            >
          </div>
        </div>

        ${binding.mappings.length
          ? html`
              <table class="mappings-table">
                <thead>
                  <tr>
                    <th>From (Your Token)</th>
                    <th></th>
                    <th>To (Component Token)</th>
                    <th>Preview</th>
                  </tr>
                </thead>
                <tbody>
                  ${binding.mappings.map((m) => {
                    const colorValue = this._getColorValue(m.from);
                    const targetToken = contract?.tokens.find(
                      (t) => t.id === m.to
                    );

                    return html`
                      <tr>
                        <td class="mapping-from">${m.from}</td>
                        <td class="mapping-arrow">→</td>
                        <td>
                          <span class="mapping-to">${m.to}</span>
                          ${m.transform
                            ? html`<span class="mapping-transform"
                                >${m.transform}</span
                              >`
                            : null}
                        </td>
                        <td>
                          <div class="mapping-preview">
                            ${colorValue
                              ? html`<span
                                  class="color-preview"
                                  style="background: ${colorValue}"
                                ></span>`
                              : null}
                            <span style="font-size: 0.7rem; color: var(--dss-text-muted)">
                              ${colorValue || 'Token not found'}
                            </span>
                          </div>
                        </td>
                      </tr>
                    `;
                  })}
                </tbody>
              </table>
            `
          : html`
              <div class="empty-mappings">
                No explicit mappings. Using ${binding.strategy || 'default'}
                strategy.
              </div>
            `}

        <div class="generate-section">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem">
            <h4 class="generate-title" style="margin: 0">Generated CSS</h4>
            <button
              class="btn"
              style="padding: 0.25rem 0.5rem; font-size: 0.7rem"
              @click=${() => this._toggleExpanded(index)}
            >
              ${isExpanded ? 'Hide' : 'Show'}
            </button>
          </div>
          ${isExpanded
            ? html`
                <pre class="generate-output">${this._generateCSS(binding)}</pre>
              `
            : null}
        </div>
      </div>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'dss-contracts': DssContracts;
  }
}
