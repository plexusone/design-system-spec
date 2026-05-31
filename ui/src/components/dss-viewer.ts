import { LitElement, html, css, nothing } from 'lit';
import { customElement, property, state } from 'lit/decorators.js';
import { sharedStyles } from '../styles.js';
import type { DesignSystem } from '../types.js';

// Import sub-components
import './dss-colors.js';
import './dss-typography.js';
import './dss-spacing.js';
import './dss-components.js';
import './dss-contracts.js';
import './dss-editor.js';

type ViewMode = 'viewer' | 'editor';
type Section = 'overview' | 'colors' | 'typography' | 'spacing' | 'components' | 'contracts';

@customElement('dss-viewer')
export class DssViewer extends LitElement {
  static styles = [
    sharedStyles,
    css`
      :host {
        display: block;
        min-height: 100vh;
        background: var(--dss-bg-primary);
      }

      .viewer-container {
        max-width: 1400px;
        margin: 0 auto;
        padding: 1.5rem;
      }

      .header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 2rem;
        padding-bottom: 1.5rem;
        border-bottom: 1px solid var(--dss-border);
      }

      .header-info {
        flex: 1;
      }

      .header-title {
        font-size: 2rem;
        font-weight: 700;
        margin: 0 0 0.5rem;
      }

      .header-version {
        display: inline-block;
        padding: 0.25rem 0.5rem;
        font-size: 0.75rem;
        font-weight: 500;
        background: var(--dss-accent);
        color: var(--dss-bg-primary);
        border-radius: var(--dss-radius-sm);
        margin-left: 0.75rem;
      }

      .header-description {
        font-size: 1rem;
        color: var(--dss-text-secondary);
        margin: 0;
        max-width: 600px;
      }

      .header-actions {
        display: flex;
        gap: 0.5rem;
        align-items: center;
      }

      .mode-toggle {
        display: flex;
        border: 1px solid var(--dss-border);
        border-radius: var(--dss-radius-md);
        overflow: hidden;
      }

      .mode-btn {
        padding: 0.5rem 1rem;
        font-size: 0.875rem;
        font-weight: 500;
        background: var(--dss-bg-secondary);
        color: var(--dss-text-secondary);
        border: none;
        cursor: pointer;
        transition: all 0.2s;
      }

      .mode-btn:hover {
        background: var(--dss-bg-tertiary);
      }

      .mode-btn--active {
        background: var(--dss-accent);
        color: var(--dss-bg-primary);
      }

      .nav {
        display: flex;
        gap: 0.25rem;
        margin-bottom: 1.5rem;
        padding: 0.25rem;
        background: var(--dss-bg-secondary);
        border-radius: var(--dss-radius-md);
        overflow-x: auto;
      }

      .nav-btn {
        padding: 0.5rem 1rem;
        font-size: 0.875rem;
        font-weight: 500;
        background: transparent;
        color: var(--dss-text-secondary);
        border: none;
        border-radius: var(--dss-radius-sm);
        cursor: pointer;
        white-space: nowrap;
        transition: all 0.2s;
      }

      .nav-btn:hover {
        background: var(--dss-bg-tertiary);
        color: var(--dss-text-primary);
      }

      .nav-btn--active {
        background: var(--dss-bg-primary);
        color: var(--dss-text-primary);
        box-shadow: var(--dss-shadow);
      }

      .nav-count {
        font-size: 0.7rem;
        margin-left: 0.375rem;
        padding: 0.125rem 0.375rem;
        background: var(--dss-bg-tertiary);
        border-radius: var(--dss-radius-sm);
      }

      .content {
        animation: fadeIn 0.2s ease;
      }

      @keyframes fadeIn {
        from {
          opacity: 0;
          transform: translateY(4px);
        }
        to {
          opacity: 1;
          transform: translateY(0);
        }
      }

      .overview-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
        gap: 1rem;
        margin-bottom: 2rem;
      }

      .overview-card {
        background: var(--dss-bg-secondary);
        border: 1px solid var(--dss-border);
        border-radius: var(--dss-radius-md);
        padding: 1.25rem;
        cursor: pointer;
        transition: all 0.2s;
      }

      .overview-card:hover {
        border-color: var(--dss-accent);
        transform: translateY(-2px);
        box-shadow: var(--dss-shadow);
      }

      .overview-card-count {
        font-size: 2rem;
        font-weight: 700;
        color: var(--dss-accent);
        margin: 0;
      }

      .overview-card-label {
        font-size: 0.875rem;
        color: var(--dss-text-secondary);
        margin: 0.25rem 0 0;
      }

      .empty-state {
        text-align: center;
        padding: 4rem 2rem;
      }

      .empty-state-title {
        font-size: 1.5rem;
        font-weight: 600;
        margin: 0 0 0.5rem;
        color: var(--dss-text-secondary);
      }

      .empty-state-text {
        font-size: 1rem;
        color: var(--dss-text-muted);
        margin: 0 0 1.5rem;
      }

      .toast {
        position: fixed;
        bottom: 1.5rem;
        right: 1.5rem;
        padding: 0.75rem 1.25rem;
        background: var(--dss-bg-tertiary);
        border: 1px solid var(--dss-border);
        border-radius: var(--dss-radius-md);
        font-size: 0.875rem;
        color: var(--dss-text-primary);
        box-shadow: var(--dss-shadow);
        animation: slideIn 0.3s ease;
        z-index: 1000;
      }

      @keyframes slideIn {
        from {
          opacity: 0;
          transform: translateY(1rem);
        }
        to {
          opacity: 1;
          transform: translateY(0);
        }
      }

      .principles-section {
        margin-bottom: 2rem;
      }

      .principles-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
        gap: 1rem;
      }

      .principle-card {
        background: var(--dss-bg-secondary);
        border: 1px solid var(--dss-border);
        border-radius: var(--dss-radius-md);
        padding: 1rem;
      }

      .principle-name {
        font-size: 1rem;
        font-weight: 600;
        margin: 0 0 0.5rem;
      }

      .principle-description {
        font-size: 0.875rem;
        color: var(--dss-text-secondary);
        margin: 0;
      }

      .file-input {
        display: none;
      }

      .upload-area {
        border: 2px dashed var(--dss-border);
        border-radius: var(--dss-radius-lg);
        padding: 3rem 2rem;
        text-align: center;
        cursor: pointer;
        transition: all 0.2s;
        margin-top: 1rem;
      }

      .upload-area:hover {
        border-color: var(--dss-accent);
        background: var(--dss-bg-secondary);
      }

      .upload-icon {
        font-size: 3rem;
        margin-bottom: 1rem;
      }

      .upload-text {
        font-size: 1rem;
        color: var(--dss-text-secondary);
        margin: 0;
      }

      .upload-hint {
        font-size: 0.875rem;
        color: var(--dss-text-muted);
        margin: 0.5rem 0 0;
      }

      .empty-actions {
        display: flex;
        gap: 1rem;
        justify-content: center;
        flex-wrap: wrap;
      }
    `,
  ];

  @property({ type: Object })
  spec: DesignSystem | null = null;

  @property({ type: String })
  src = '';

  @property({ type: String, reflect: true })
  theme: 'dark' | 'light' = 'dark';

  @property({ type: String })
  mode: ViewMode = 'viewer';

  @state()
  private _activeSection: Section = 'overview';

  @state()
  private _loading = false;

  @state()
  private _error: string | null = null;

  @state()
  private _toast: string | null = null;

  connectedCallback() {
    super.connectedCallback();

    // Listen for copy events from children
    this.addEventListener('dss-copy', ((e: CustomEvent) => {
      this._showToast(e.detail.value || 'Copied!');
    }) as EventListener);

    // Listen for spec changes from editor
    this.addEventListener('dss-spec-change', ((e: CustomEvent) => {
      this.spec = e.detail.spec;
      this._showToast('Changes applied');
    }) as EventListener);

    // Load from src if provided
    if (this.src && !this.spec) {
      this._loadFromUrl(this.src);
    }
  }

  private async _loadFromUrl(url: string) {
    this._loading = true;
    this._error = null;

    try {
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error(`Failed to load: ${response.status}`);
      }
      this.spec = await response.json();
    } catch (err) {
      this._error = (err as Error).message;
    } finally {
      this._loading = false;
    }
  }

  private _showToast(message: string) {
    this._toast = message;
    setTimeout(() => {
      this._toast = null;
    }, 2000);
  }

  private _getCounts() {
    if (!this.spec) return {};

    return {
      colors: this.spec.foundations?.colors?.length || 0,
      typography: this.spec.foundations?.typography ? 1 : 0,
      spacing: (this.spec.foundations?.spacing?.scale?.length || 0) +
        (this.spec.foundations?.borderRadius?.length || 0) +
        (this.spec.foundations?.elevation?.length || 0),
      components: this.spec.components?.length || 0,
      contracts: this.spec.themeBindings?.length || 0,
    };
  }

  private _triggerFileUpload() {
    const input = this.shadowRoot?.querySelector('.file-input') as HTMLInputElement;
    input?.click();
  }

  private _handleFileUpload(e: Event) {
    const input = e.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;

    this._loading = true;
    this._error = null;

    const reader = new FileReader();
    reader.onload = (event) => {
      try {
        const text = event.target?.result as string;
        this.spec = JSON.parse(text) as DesignSystem;
        this._showToast(`Loaded: ${this.spec.meta?.name || file.name}`);
      } catch (err) {
        this._error = `Invalid JSON: ${(err as Error).message}`;
      } finally {
        this._loading = false;
      }
    };
    reader.onerror = () => {
      this._error = 'Failed to read file';
      this._loading = false;
    };
    reader.readAsText(file);

    // Reset input so the same file can be selected again
    input.value = '';
  }

  private _handleDrop(e: DragEvent) {
    e.preventDefault();
    const file = e.dataTransfer?.files?.[0];
    if (!file || !file.name.endsWith('.json')) {
      this._showToast('Please drop a JSON file');
      return;
    }

    // Create a fake event to reuse the handler
    const fakeInput = { files: [file], value: '' } as unknown as HTMLInputElement;
    this._handleFileUpload({ target: fakeInput } as unknown as Event);
  }

  private _handleDragOver(e: DragEvent) {
    e.preventDefault();
  }

  render() {
    if (this._loading) {
      return html`
        <div class="viewer-container">
          <div class="empty-state">
            <p class="empty-state-title">Loading...</p>
          </div>
        </div>
      `;
    }

    if (this._error) {
      return html`
        <div class="viewer-container">
          <div class="empty-state">
            <p class="empty-state-title">Error Loading Spec</p>
            <p class="empty-state-text">${this._error}</p>
            <button class="btn btn--primary" @click=${() => this._loadFromUrl(this.src)}>
              Retry
            </button>
          </div>
        </div>
      `;
    }

    if (!this.spec && this.mode === 'viewer') {
      return html`
        <div class="viewer-container">
          <input
            type="file"
            class="file-input"
            accept=".json"
            @change=${this._handleFileUpload}
          />
          <div class="empty-state">
            <p class="empty-state-title">No Design System Loaded</p>
            <p class="empty-state-text">
              Upload a JSON file, provide a URL, or create a new spec in Editor mode.
            </p>
            <div class="empty-actions">
              <button class="btn btn--primary" @click=${this._triggerFileUpload}>
                Upload JSON File
              </button>
              <button class="btn" @click=${() => (this.mode = 'editor')}>
                Open Editor
              </button>
            </div>
            <div
              class="upload-area"
              @click=${this._triggerFileUpload}
              @drop=${this._handleDrop}
              @dragover=${this._handleDragOver}
            >
              <div class="upload-icon">📁</div>
              <p class="upload-text">Drop your design-system.json here</p>
              <p class="upload-hint">or click to browse</p>
            </div>
          </div>
        </div>
      `;
    }

    return html`
      <div class="viewer-container">
        ${this._renderHeader()}
        ${this.mode === 'editor'
          ? this._renderEditor()
          : this._renderViewer()}
        ${this._toast ? html`<div class="toast">${this._toast}</div>` : nothing}
      </div>
    `;
  }

  private _renderHeader() {
    return html`
      <input
        type="file"
        class="file-input"
        accept=".json"
        @change=${this._handleFileUpload}
      />
      <header class="header">
        <div class="header-info">
          <h1 class="header-title">
            ${this.spec?.meta?.name || 'Design System'}
            ${this.spec?.meta?.version
              ? html`<span class="header-version">v${this.spec.meta.version}</span>`
              : nothing}
          </h1>
          ${this.spec?.meta?.description
            ? html`<p class="header-description">${this.spec.meta.description}</p>`
            : nothing}
        </div>
        <div class="header-actions">
          <button class="btn" @click=${this._triggerFileUpload}>
            Load File
          </button>
          <div class="mode-toggle">
            <button
              class="mode-btn ${this.mode === 'viewer' ? 'mode-btn--active' : ''}"
              @click=${() => (this.mode = 'viewer')}
            >
              Viewer
            </button>
            <button
              class="mode-btn ${this.mode === 'editor' ? 'mode-btn--active' : ''}"
              @click=${() => (this.mode = 'editor')}
            >
              Editor
            </button>
          </div>
        </div>
      </header>
    `;
  }

  private _renderViewer() {
    const counts = this._getCounts();

    return html`
      <nav class="nav">
        <button
          class="nav-btn ${this._activeSection === 'overview' ? 'nav-btn--active' : ''}"
          @click=${() => (this._activeSection = 'overview')}
        >
          Overview
        </button>
        <button
          class="nav-btn ${this._activeSection === 'colors' ? 'nav-btn--active' : ''}"
          @click=${() => (this._activeSection = 'colors')}
        >
          Colors
          ${counts.colors ? html`<span class="nav-count">${counts.colors}</span>` : nothing}
        </button>
        <button
          class="nav-btn ${this._activeSection === 'typography' ? 'nav-btn--active' : ''}"
          @click=${() => (this._activeSection = 'typography')}
        >
          Typography
        </button>
        <button
          class="nav-btn ${this._activeSection === 'spacing' ? 'nav-btn--active' : ''}"
          @click=${() => (this._activeSection = 'spacing')}
        >
          Spacing & Effects
          ${counts.spacing ? html`<span class="nav-count">${counts.spacing}</span>` : nothing}
        </button>
        <button
          class="nav-btn ${this._activeSection === 'components' ? 'nav-btn--active' : ''}"
          @click=${() => (this._activeSection = 'components')}
        >
          Components
          ${counts.components ? html`<span class="nav-count">${counts.components}</span>` : nothing}
        </button>
        <button
          class="nav-btn ${this._activeSection === 'contracts' ? 'nav-btn--active' : ''}"
          @click=${() => (this._activeSection = 'contracts')}
        >
          Theme Bindings
          ${counts.contracts ? html`<span class="nav-count">${counts.contracts}</span>` : nothing}
        </button>
      </nav>

      <div class="content">
        ${this._renderActiveSection()}
      </div>
    `;
  }

  private _renderActiveSection() {
    switch (this._activeSection) {
      case 'overview':
        return this._renderOverview();
      case 'colors':
        return html`<dss-colors .colors=${this.spec?.foundations?.colors || []}></dss-colors>`;
      case 'typography':
        return html`<dss-typography .typography=${this.spec?.foundations?.typography}></dss-typography>`;
      case 'spacing':
        return html`
          <dss-spacing
            .spacing=${this.spec?.foundations?.spacing}
            .borderRadius=${this.spec?.foundations?.borderRadius || []}
            .elevation=${this.spec?.foundations?.elevation || []}
          ></dss-spacing>
        `;
      case 'components':
        return html`<dss-components .components=${this.spec?.components || []}></dss-components>`;
      case 'contracts':
        return html`
          <dss-contracts
            .themeBindings=${this.spec?.themeBindings || []}
            .components=${this.spec?.components || []}
            .colors=${this.spec?.foundations?.colors || []}
          ></dss-contracts>
        `;
      default:
        return nothing;
    }
  }

  private _renderOverview() {
    const counts = this._getCounts();

    return html`
      <div class="overview-grid">
        <div class="overview-card" @click=${() => (this._activeSection = 'colors')}>
          <p class="overview-card-count">${counts.colors}</p>
          <p class="overview-card-label">Color Tokens</p>
        </div>
        <div class="overview-card" @click=${() => (this._activeSection = 'typography')}>
          <p class="overview-card-count">
            ${this.spec?.foundations?.typography?.fontFamilies?.length || 0}
          </p>
          <p class="overview-card-label">Font Families</p>
        </div>
        <div class="overview-card" @click=${() => (this._activeSection = 'spacing')}>
          <p class="overview-card-count">
            ${this.spec?.foundations?.spacing?.scale?.length || 0}
          </p>
          <p class="overview-card-label">Spacing Values</p>
        </div>
        <div class="overview-card" @click=${() => (this._activeSection = 'components')}>
          <p class="overview-card-count">${counts.components}</p>
          <p class="overview-card-label">Components</p>
        </div>
        <div class="overview-card" @click=${() => (this._activeSection = 'contracts')}>
          <p class="overview-card-count">${counts.contracts}</p>
          <p class="overview-card-label">Theme Bindings</p>
        </div>
      </div>

      ${this.spec?.principles?.length
        ? html`
            <div class="principles-section section">
              <div class="section-header">
                <h2 class="section-title">Design Principles</h2>
              </div>
              <div class="principles-grid">
                ${this.spec.principles.map(
                  (p) => html`
                    <div class="principle-card">
                      <h3 class="principle-name">${p.name}</h3>
                      ${p.description
                        ? html`<p class="principle-description">${p.description}</p>`
                        : nothing}
                    </div>
                  `
                )}
              </div>
            </div>
          `
        : nothing}
    `;
  }

  private _renderEditor() {
    return html`
      <dss-editor .spec=${this.spec}></dss-editor>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'dss-viewer': DssViewer;
  }
}
