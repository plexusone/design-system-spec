import { LitElement, html, css } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import { sharedStyles } from '../styles.js';
import type { ColorToken } from '../types.js';

@customElement('dss-colors')
export class DssColors extends LitElement {
  static styles = [
    sharedStyles,
    css`
      .color-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
        gap: 1rem;
      }

      .color-card {
        background: var(--dss-bg-tertiary);
        border: 1px solid var(--dss-border);
        border-radius: var(--dss-radius-md);
        overflow: hidden;
        transition: transform 0.2s, box-shadow 0.2s;
        cursor: pointer;
      }

      .color-card:hover {
        transform: translateY(-2px);
        box-shadow: var(--dss-shadow);
      }

      .color-swatch {
        height: 80px;
        position: relative;
      }

      .color-swatch::after {
        content: attr(data-value);
        position: absolute;
        bottom: 4px;
        right: 4px;
        font-size: 0.65rem;
        font-family: var(--dss-font-mono);
        padding: 2px 4px;
        border-radius: 2px;
        background: rgba(0, 0, 0, 0.5);
        color: white;
      }

      .color-info {
        padding: 0.75rem;
      }

      .color-id {
        font-size: 0.8rem;
        font-weight: 600;
        margin: 0 0 0.25rem;
        word-break: break-word;
      }

      .color-semantic {
        font-size: 0.7rem;
        color: var(--dss-text-muted);
        margin: 0;
      }

      .color-modes {
        display: flex;
        gap: 0.5rem;
        margin-top: 0.5rem;
      }

      .color-mode {
        flex: 1;
        height: 20px;
        border-radius: 2px;
        border: 1px solid var(--dss-border);
      }

      .color-mode-label {
        font-size: 0.6rem;
        color: var(--dss-text-muted);
        text-align: center;
        margin-top: 2px;
      }

      .contrast-info {
        display: flex;
        gap: 0.25rem;
        margin-top: 0.5rem;
      }

      .contrast-badge {
        font-size: 0.6rem;
        padding: 1px 4px;
        border-radius: 2px;
        background: var(--dss-bg-secondary);
      }

      .contrast-badge--aa {
        background: var(--dss-success);
        color: white;
      }

      .contrast-badge--aaa {
        background: var(--dss-accent);
        color: white;
      }

      .group-header {
        font-size: 0.875rem;
        font-weight: 600;
        color: var(--dss-text-secondary);
        margin: 1.5rem 0 0.75rem;
        padding-bottom: 0.5rem;
        border-bottom: 1px solid var(--dss-border);
      }

      .group-header:first-child {
        margin-top: 0;
      }
    `,
  ];

  @property({ type: Array })
  colors: ColorToken[] = [];

  @property({ type: Boolean })
  editable = false;

  private _groupColors(colors: ColorToken[]): Map<string, ColorToken[]> {
    const groups = new Map<string, ColorToken[]>();

    for (const color of colors) {
      // Group by semantic or by prefix (e.g., "primary-500" -> "primary")
      let group = color.semantic || 'other';
      if (!color.semantic) {
        const match = color.id.match(/^([a-z]+)/i);
        if (match) {
          group = match[1];
        }
      }

      if (!groups.has(group)) {
        groups.set(group, []);
      }
      groups.get(group)!.push(color);
    }

    return groups;
  }

  private _getContrastColor(hex: string): string {
    // Simple luminance calculation for text contrast
    const r = parseInt(hex.slice(1, 3), 16);
    const g = parseInt(hex.slice(3, 5), 16);
    const b = parseInt(hex.slice(5, 7), 16);
    const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255;
    return luminance > 0.5 ? '#000000' : '#ffffff';
  }

  private _copyToClipboard(value: string) {
    navigator.clipboard.writeText(value);
    this.dispatchEvent(
      new CustomEvent('dss-copy', {
        detail: { value },
        bubbles: true,
        composed: true,
      })
    );
  }

  render() {
    if (!this.colors?.length) {
      return html`
        <div class="section">
          <div class="section-header">
            <h2 class="section-title">Colors</h2>
          </div>
          <div class="empty-state">No colors defined</div>
        </div>
      `;
    }

    const groups = this._groupColors(this.colors);

    return html`
      <div class="section">
        <div class="section-header">
          <h2 class="section-title">Colors</h2>
          <span class="badge">${this.colors.length} tokens</span>
        </div>

        ${Array.from(groups.entries()).map(
          ([group, colors]) => html`
            <div class="group-header">${group}</div>
            <div class="color-grid">
              ${colors.map((color) => this._renderColorCard(color))}
            </div>
          `
        )}
      </div>
    `;
  }

  private _renderColorCard(color: ColorToken) {
    const swatchColor = color.value || '#000000';

    return html`
      <div
        class="color-card"
        @click=${() => this._copyToClipboard(color.value)}
        title="Click to copy ${color.value}"
      >
        <div
          class="color-swatch"
          style="background: ${swatchColor}"
          data-value="${color.value}"
        ></div>
        <div class="color-info">
          <p class="color-id">${color.id}</p>
          ${color.semantic
            ? html`<p class="color-semantic">${color.semantic}</p>`
            : null}
          ${color.lightModeValue || color.darkModeValue
            ? html`
                <div class="color-modes">
                  ${color.lightModeValue
                    ? html`
                        <div>
                          <div
                            class="color-mode"
                            style="background: ${color.lightModeValue}"
                          ></div>
                          <div class="color-mode-label">Light</div>
                        </div>
                      `
                    : null}
                  ${color.darkModeValue
                    ? html`
                        <div>
                          <div
                            class="color-mode"
                            style="background: ${color.darkModeValue}"
                          ></div>
                          <div class="color-mode-label">Dark</div>
                        </div>
                      `
                    : null}
                </div>
              `
            : null}
          ${color.contrast?.wcagLevel
            ? html`
                <div class="contrast-info">
                  <span
                    class="contrast-badge contrast-badge--${color.contrast
                      .wcagLevel.toLowerCase()}"
                  >
                    WCAG ${color.contrast.wcagLevel}
                  </span>
                  ${color.contrast.ratio
                    ? html`<span class="contrast-badge"
                        >${color.contrast.ratio.toFixed(1)}:1</span
                      >`
                    : null}
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
    'dss-colors': DssColors;
  }
}
