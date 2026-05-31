import { LitElement, html, css } from 'lit';
import { customElement, property, state } from 'lit/decorators.js';
import { sharedStyles } from '../styles.js';
import type { DesignSystem } from '../types.js';

@customElement('dss-editor')
export class DssEditor extends LitElement {
  static styles = [
    sharedStyles,
    css`
      :host {
        display: block;
      }

      .editor-container {
        display: flex;
        flex-direction: column;
        height: 100%;
        min-height: 400px;
      }

      .editor-toolbar {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0.75rem 1rem;
        background: var(--dss-bg-secondary);
        border: 1px solid var(--dss-border);
        border-bottom: none;
        border-radius: var(--dss-radius-md) var(--dss-radius-md) 0 0;
      }

      .toolbar-left {
        display: flex;
        align-items: center;
        gap: 0.75rem;
      }

      .toolbar-right {
        display: flex;
        align-items: center;
        gap: 0.5rem;
      }

      .editor-status {
        font-size: 0.75rem;
        color: var(--dss-text-muted);
      }

      .editor-status--modified {
        color: var(--dss-warning);
      }

      .editor-status--error {
        color: var(--dss-danger);
      }

      .editor-status--valid {
        color: var(--dss-success);
      }

      .textarea-wrapper {
        flex: 1;
        position: relative;
        border: 1px solid var(--dss-border);
        border-radius: 0 0 var(--dss-radius-md) var(--dss-radius-md);
        overflow: hidden;
      }

      .editor-textarea {
        width: 100%;
        height: 100%;
        min-height: 300px;
        padding: 1rem;
        font-family: var(--dss-font-mono);
        font-size: 0.8rem;
        line-height: 1.6;
        background: var(--dss-bg-primary);
        color: var(--dss-text-primary);
        border: none;
        resize: none;
        outline: none;
      }

      .editor-textarea::placeholder {
        color: var(--dss-text-muted);
      }

      .error-panel {
        padding: 0.75rem 1rem;
        background: rgba(239, 68, 68, 0.1);
        border-top: 1px solid var(--dss-danger);
        font-size: 0.75rem;
        color: var(--dss-danger);
        font-family: var(--dss-font-mono);
      }

      .download-modal {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0, 0, 0, 0.7);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 1000;
      }

      .modal-content {
        background: var(--dss-bg-secondary);
        border: 1px solid var(--dss-border);
        border-radius: var(--dss-radius-lg);
        padding: 1.5rem;
        max-width: 500px;
        width: 90%;
      }

      .modal-title {
        font-size: 1.25rem;
        font-weight: 600;
        margin: 0 0 1rem;
      }

      .modal-text {
        font-size: 0.875rem;
        color: var(--dss-text-secondary);
        margin: 0 0 1.5rem;
      }

      .modal-actions {
        display: flex;
        justify-content: flex-end;
        gap: 0.5rem;
      }

      .file-input {
        display: none;
      }

      .tabs {
        display: flex;
        gap: 0;
        border-bottom: 1px solid var(--dss-border);
        background: var(--dss-bg-secondary);
      }

      .tab {
        padding: 0.75rem 1.5rem;
        font-size: 0.875rem;
        font-weight: 500;
        color: var(--dss-text-secondary);
        background: transparent;
        border: none;
        border-bottom: 2px solid transparent;
        cursor: pointer;
        transition: all 0.2s;
      }

      .tab:hover {
        color: var(--dss-text-primary);
        background: var(--dss-bg-tertiary);
      }

      .tab--active {
        color: var(--dss-accent);
        border-bottom-color: var(--dss-accent);
      }
    `,
  ];

  @property({ type: Object })
  spec: DesignSystem | null = null;

  @property({ type: String })
  initialJson = '';

  @state()
  private _jsonText = '';

  @state()
  private _isModified = false;

  @state()
  private _parseError: string | null = null;

  @state()
  private _showDownloadModal = false;

  @state()
  private _activeTab: 'edit' | 'formatted' = 'edit';

  connectedCallback() {
    super.connectedCallback();
    if (this.spec) {
      this._jsonText = JSON.stringify(this.spec, null, 2);
    } else if (this.initialJson) {
      this._jsonText = this.initialJson;
    }
  }

  updated(changedProperties: Map<string, unknown>) {
    if (changedProperties.has('spec') && this.spec && !this._isModified) {
      this._jsonText = JSON.stringify(this.spec, null, 2);
    }
  }

  private _handleInput(e: Event) {
    const textarea = e.target as HTMLTextAreaElement;
    this._jsonText = textarea.value;
    this._isModified = true;

    // Validate JSON
    try {
      JSON.parse(this._jsonText);
      this._parseError = null;
    } catch (err) {
      this._parseError = (err as Error).message;
    }
  }

  private _formatJson() {
    try {
      const parsed = JSON.parse(this._jsonText);
      this._jsonText = JSON.stringify(parsed, null, 2);
      this._parseError = null;
    } catch (err) {
      this._parseError = (err as Error).message;
    }
  }

  private _minifyJson() {
    try {
      const parsed = JSON.parse(this._jsonText);
      this._jsonText = JSON.stringify(parsed);
      this._parseError = null;
    } catch (err) {
      this._parseError = (err as Error).message;
    }
  }

  private _copyToClipboard() {
    navigator.clipboard.writeText(this._jsonText);
    this.dispatchEvent(
      new CustomEvent('dss-copy', {
        detail: { value: 'JSON copied to clipboard' },
        bubbles: true,
        composed: true,
      })
    );
  }

  private _downloadJson() {
    try {
      // Validate before download
      const parsed = JSON.parse(this._jsonText);
      const formatted = JSON.stringify(parsed, null, 2);

      const blob = new Blob([formatted], { type: 'application/json' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `design-system-${Date.now()}.json`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);

      this._showDownloadModal = false;
    } catch (err) {
      this._parseError = (err as Error).message;
    }
  }

  private _applyChanges() {
    try {
      const parsed = JSON.parse(this._jsonText) as DesignSystem;
      this._parseError = null;
      this._isModified = false;

      this.dispatchEvent(
        new CustomEvent('dss-spec-change', {
          detail: { spec: parsed },
          bubbles: true,
          composed: true,
        })
      );
    } catch (err) {
      this._parseError = (err as Error).message;
    }
  }

  private _handleFileUpload(e: Event) {
    const input = e.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (event) => {
      const text = event.target?.result as string;
      try {
        const parsed = JSON.parse(text);
        this._jsonText = JSON.stringify(parsed, null, 2);
        this._parseError = null;
        this._isModified = true;
      } catch (err) {
        this._parseError = `Invalid JSON file: ${(err as Error).message}`;
      }
    };
    reader.readAsText(file);
  }

  private _triggerFileUpload() {
    const input = this.shadowRoot?.querySelector('.file-input') as HTMLInputElement;
    input?.click();
  }

  render() {
    const statusClass = this._parseError
      ? 'editor-status--error'
      : this._isModified
      ? 'editor-status--modified'
      : 'editor-status--valid';

    const statusText = this._parseError
      ? 'Invalid JSON'
      : this._isModified
      ? 'Modified'
      : 'Valid';

    return html`
      <div class="editor-container">
        <div class="editor-toolbar">
          <div class="toolbar-left">
            <span class="editor-status ${statusClass}">${statusText}</span>
          </div>
          <div class="toolbar-right">
            <input
              type="file"
              class="file-input"
              accept=".json"
              @change=${this._handleFileUpload}
            />
            <button class="btn" @click=${this._triggerFileUpload}>
              Load File
            </button>
            <button class="btn" @click=${this._formatJson}>Format</button>
            <button class="btn" @click=${this._minifyJson}>Minify</button>
            <button class="btn" @click=${this._copyToClipboard}>Copy</button>
            <button
              class="btn btn--primary"
              @click=${() => (this._showDownloadModal = true)}
              ?disabled=${!!this._parseError}
            >
              Download
            </button>
            ${this._isModified && !this._parseError
              ? html`
                  <button class="btn btn--primary" @click=${this._applyChanges}>
                    Apply
                  </button>
                `
              : null}
          </div>
        </div>

        <div class="textarea-wrapper">
          <textarea
            class="editor-textarea"
            .value=${this._jsonText}
            @input=${this._handleInput}
            placeholder="Paste or edit your design system JSON here..."
            spellcheck="false"
          ></textarea>
        </div>

        ${this._parseError
          ? html`<div class="error-panel">${this._parseError}</div>`
          : null}
      </div>

      ${this._showDownloadModal
        ? html`
            <div
              class="download-modal"
              @click=${(e: Event) => {
                if (e.target === e.currentTarget) {
                  this._showDownloadModal = false;
                }
              }}
            >
              <div class="modal-content">
                <h3 class="modal-title">Download Design System</h3>
                <p class="modal-text">
                  Your design system JSON will be downloaded as a file. You can
                  use this file with the DSS CLI or import it later.
                </p>
                <div class="modal-actions">
                  <button
                    class="btn"
                    @click=${() => (this._showDownloadModal = false)}
                  >
                    Cancel
                  </button>
                  <button class="btn btn--primary" @click=${this._downloadJson}>
                    Download JSON
                  </button>
                </div>
              </div>
            </div>
          `
        : null}
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'dss-editor': DssEditor;
  }
}
