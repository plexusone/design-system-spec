import { LitElement, html, css } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import { sharedStyles } from '../styles.js';
import type { SpacingScale, BorderRadiusToken, ElevationToken } from '../types.js';

@customElement('dss-spacing')
export class DssSpacing extends LitElement {
  static styles = [
    sharedStyles,
    css`
      .spacing-scale {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
      }

      .spacing-row {
        display: flex;
        align-items: center;
        gap: 1rem;
        padding: 0.5rem;
        background: var(--dss-bg-tertiary);
        border-radius: var(--dss-radius-sm);
      }

      .spacing-label {
        min-width: 40px;
        font-size: 0.75rem;
        font-weight: 600;
        color: var(--dss-text-secondary);
      }

      .spacing-value {
        min-width: 80px;
        font-size: 0.7rem;
        font-family: var(--dss-font-mono);
        color: var(--dss-text-muted);
      }

      .spacing-visual {
        flex: 1;
        display: flex;
        align-items: center;
      }

      .spacing-bar {
        height: 16px;
        background: var(--dss-accent);
        border-radius: 2px;
        min-width: 2px;
      }

      .base-unit {
        font-size: 0.875rem;
        color: var(--dss-text-secondary);
        margin-bottom: 1rem;
        padding: 0.5rem;
        background: var(--dss-bg-tertiary);
        border-radius: var(--dss-radius-sm);
        display: inline-block;
      }

      .radius-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
        gap: 1rem;
      }

      .radius-card {
        text-align: center;
        padding: 1rem;
        background: var(--dss-bg-tertiary);
        border-radius: var(--dss-radius-md);
      }

      .radius-preview {
        width: 60px;
        height: 60px;
        background: var(--dss-accent);
        margin: 0 auto 0.5rem;
      }

      .radius-label {
        font-size: 0.75rem;
        font-weight: 600;
        margin: 0;
      }

      .radius-value {
        font-size: 0.7rem;
        font-family: var(--dss-font-mono);
        color: var(--dss-text-muted);
        margin: 0;
      }

      .elevation-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
        gap: 1.5rem;
      }

      .elevation-card {
        padding: 1.5rem;
        background: var(--dss-bg-tertiary);
        border-radius: var(--dss-radius-md);
        text-align: center;
      }

      .elevation-label {
        font-size: 0.875rem;
        font-weight: 600;
        margin: 0 0 0.25rem;
      }

      .elevation-value {
        font-size: 0.65rem;
        font-family: var(--dss-font-mono);
        color: var(--dss-text-muted);
        margin: 0;
        word-break: break-word;
      }

      .subsection {
        margin-bottom: 1.5rem;
      }

      .subsection-title {
        font-size: 1rem;
        font-weight: 600;
        margin: 0 0 1rem;
        color: var(--dss-text-secondary);
      }
    `,
  ];

  @property({ type: Object })
  spacing: SpacingScale | null = null;

  @property({ type: Array })
  borderRadius: BorderRadiusToken[] = [];

  @property({ type: Array })
  elevation: ElevationToken[] = [];

  @property({ type: Boolean })
  editable = false;

  render() {
    const hasSpacing = this.spacing?.scale?.length;
    const hasRadius = this.borderRadius?.length;
    const hasElevation = this.elevation?.length;

    if (!hasSpacing && !hasRadius && !hasElevation) {
      return html`
        <div class="section">
          <div class="section-header">
            <h2 class="section-title">Spacing & Effects</h2>
          </div>
          <div class="empty-state">No spacing or effects defined</div>
        </div>
      `;
    }

    return html`
      <div class="section">
        <div class="section-header">
          <h2 class="section-title">Spacing & Effects</h2>
        </div>

        ${hasSpacing ? this._renderSpacing() : null}
        ${hasRadius ? this._renderBorderRadius() : null}
        ${hasElevation ? this._renderElevation() : null}
      </div>
    `;
  }

  private _renderSpacing() {
    return html`
      <div class="subsection">
        <h3 class="subsection-title">Spacing Scale</h3>
        ${this.spacing?.baseUnit
          ? html`<div class="base-unit">
              Base unit: <strong>${this.spacing.baseUnit}</strong>
            </div>`
          : null}
        <div class="spacing-scale">
          ${this.spacing?.scale?.map((token) => {
            // Parse pixel value for visual width
            const pixelValue =
              token.pixelValue || parseInt(token.value) || 0;
            const maxWidth = 200;
            const width = Math.min(pixelValue * 2, maxWidth);

            return html`
              <div class="spacing-row">
                <span class="spacing-label">${token.id}</span>
                <span class="spacing-value"
                  >${token.value}${token.pixelValue
                    ? ` (${token.pixelValue}px)`
                    : ''}</span
                >
                <div class="spacing-visual">
                  <div
                    class="spacing-bar"
                    style="width: ${width}px"
                  ></div>
                </div>
              </div>
            `;
          })}
        </div>
      </div>
    `;
  }

  private _renderBorderRadius() {
    return html`
      <div class="subsection">
        <h3 class="subsection-title">Border Radius</h3>
        <div class="radius-grid">
          ${this.borderRadius.map(
            (token) => html`
              <div class="radius-card">
                <div
                  class="radius-preview"
                  style="border-radius: ${token.value}"
                ></div>
                <p class="radius-label">${token.id}</p>
                <p class="radius-value">${token.value}</p>
              </div>
            `
          )}
        </div>
      </div>
    `;
  }

  private _renderElevation() {
    return html`
      <div class="subsection">
        <h3 class="subsection-title">Elevation / Shadows</h3>
        <div class="elevation-grid">
          ${this.elevation.map(
            (token) => html`
              <div class="elevation-card" style="box-shadow: ${token.value}">
                <p class="elevation-label">${token.id}</p>
                <p class="elevation-value">${token.value}</p>
                ${token.usage
                  ? html`<p
                      class="color-semantic"
                      style="margin-top: 0.5rem; font-size: 0.7rem"
                    >
                      ${token.usage}
                    </p>`
                  : null}
              </div>
            `
          )}
        </div>
      </div>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'dss-spacing': DssSpacing;
  }
}
