import { LitElement, html, css } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import { sharedStyles } from '../styles.js';
import type { Typography, FontFamily, FontSize, TypeStyle } from '../types.js';

@customElement('dss-typography')
export class DssTypography extends LitElement {
  static styles = [
    sharedStyles,
    css`
      .font-family-card {
        background: var(--dss-bg-tertiary);
        border: 1px solid var(--dss-border);
        border-radius: var(--dss-radius-md);
        padding: 1rem;
        margin-bottom: 1rem;
      }

      .font-preview {
        font-size: 2rem;
        margin-bottom: 0.5rem;
        line-height: 1.2;
      }

      .font-name {
        font-size: 1rem;
        font-weight: 600;
        margin: 0 0 0.25rem;
      }

      .font-stack {
        font-size: 0.75rem;
        color: var(--dss-text-muted);
        font-family: var(--dss-font-mono);
        margin: 0;
        word-break: break-word;
      }

      .font-weights {
        display: flex;
        flex-wrap: wrap;
        gap: 0.5rem;
        margin-top: 0.75rem;
      }

      .font-weight-sample {
        padding: 0.25rem 0.5rem;
        background: var(--dss-bg-secondary);
        border-radius: var(--dss-radius-sm);
        font-size: 0.875rem;
      }

      .size-scale {
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
      }

      .size-row {
        display: flex;
        align-items: baseline;
        gap: 1rem;
        padding: 0.5rem;
        background: var(--dss-bg-tertiary);
        border-radius: var(--dss-radius-sm);
      }

      .size-label {
        min-width: 60px;
        font-size: 0.75rem;
        font-weight: 600;
        color: var(--dss-text-secondary);
      }

      .size-value {
        min-width: 80px;
        font-size: 0.7rem;
        font-family: var(--dss-font-mono);
        color: var(--dss-text-muted);
      }

      .size-sample {
        flex: 1;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .type-scale {
        display: flex;
        flex-direction: column;
        gap: 1rem;
      }

      .type-style-card {
        background: var(--dss-bg-tertiary);
        border: 1px solid var(--dss-border);
        border-radius: var(--dss-radius-md);
        padding: 1rem;
      }

      .type-style-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 0.5rem;
      }

      .type-style-name {
        font-size: 0.875rem;
        font-weight: 600;
        margin: 0;
      }

      .type-style-id {
        font-size: 0.7rem;
        font-family: var(--dss-font-mono);
        color: var(--dss-text-muted);
      }

      .type-style-sample {
        margin: 0.75rem 0;
        padding: 0.75rem;
        background: var(--dss-bg-secondary);
        border-radius: var(--dss-radius-sm);
      }

      .type-style-props {
        display: flex;
        flex-wrap: wrap;
        gap: 0.5rem;
        font-size: 0.7rem;
        font-family: var(--dss-font-mono);
        color: var(--dss-text-muted);
      }

      .type-prop {
        padding: 0.125rem 0.375rem;
        background: var(--dss-bg-secondary);
        border-radius: 2px;
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
  typography: Typography | null = null;

  @property({ type: Boolean })
  editable = false;

  render() {
    if (!this.typography) {
      return html`
        <div class="section">
          <div class="section-header">
            <h2 class="section-title">Typography</h2>
          </div>
          <div class="empty-state">No typography defined</div>
        </div>
      `;
    }

    const { fontFamilies, fontSizes, typeScale } = this.typography;

    return html`
      <div class="section">
        <div class="section-header">
          <h2 class="section-title">Typography</h2>
        </div>

        ${fontFamilies?.length
          ? html`
              <div class="subsection">
                <h3 class="subsection-title">Font Families</h3>
                ${fontFamilies.map((font) => this._renderFontFamily(font))}
              </div>
            `
          : null}
        ${fontSizes?.length
          ? html`
              <div class="subsection">
                <h3 class="subsection-title">Font Sizes</h3>
                <div class="size-scale">
                  ${fontSizes.map((size) => this._renderFontSize(size))}
                </div>
              </div>
            `
          : null}
        ${typeScale?.length
          ? html`
              <div class="subsection">
                <h3 class="subsection-title">Type Scale</h3>
                <div class="type-scale">
                  ${typeScale.map((style) => this._renderTypeStyle(style))}
                </div>
              </div>
            `
          : null}
      </div>
    `;
  }

  private _renderFontFamily(font: FontFamily) {
    const fontStack = font.stack || font.value || font.name;

    return html`
      <div class="font-family-card">
        <div class="font-preview" style="font-family: ${fontStack}">
          The quick brown fox jumps over the lazy dog
        </div>
        <p class="font-name">${font.name}</p>
        <p class="font-stack">${fontStack}</p>
        ${font.weights?.length
          ? html`
              <div class="font-weights">
                ${font.weights.map(
                  (w) => html`
                    <span
                      class="font-weight-sample"
                      style="font-family: ${fontStack}; font-weight: ${w}"
                    >
                      ${w}
                    </span>
                  `
                )}
              </div>
            `
          : null}
        ${font.usage
          ? html`<p class="color-semantic" style="margin-top: 0.5rem">
              ${font.usage}
            </p>`
          : null}
      </div>
    `;
  }

  private _renderFontSize(size: FontSize) {
    return html`
      <div class="size-row">
        <span class="size-label">${size.id}</span>
        <span class="size-value"
          >${size.value}${size.pixelValue ? ` (${size.pixelValue}px)` : ''}</span
        >
        <span class="size-sample" style="font-size: ${size.value}">
          The quick brown fox
        </span>
      </div>
    `;
  }

  private _renderTypeStyle(style: TypeStyle) {
    return html`
      <div class="type-style-card">
        <div class="type-style-header">
          <p class="type-style-name">${style.name}</p>
          <span class="type-style-id">${style.id}</span>
        </div>
        <div class="type-style-sample">
          <span style="font-size: 1.5rem">
            The quick brown fox jumps over the lazy dog
          </span>
        </div>
        <div class="type-style-props">
          <span class="type-prop">family: ${style.fontFamily}</span>
          <span class="type-prop">size: ${style.fontSize}</span>
          <span class="type-prop">weight: ${style.fontWeight}</span>
          <span class="type-prop">line-height: ${style.lineHeight}</span>
          ${style.letterSpacing
            ? html`<span class="type-prop"
                >letter-spacing: ${style.letterSpacing}</span
              >`
            : null}
        </div>
        ${style.usage
          ? html`<p class="color-semantic" style="margin-top: 0.5rem">
              ${style.usage}
            </p>`
          : null}
      </div>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'dss-typography': DssTypography;
  }
}
