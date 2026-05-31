import { css } from 'lit';

// Shared styles for DSS UI components
export const sharedStyles = css`
  :host {
    --dss-font-sans: 'Inter', system-ui, -apple-system, sans-serif;
    --dss-font-mono: 'JetBrains Mono', 'Fira Code', monospace;

    /* Dark theme (default) */
    --dss-bg-primary: #0a0e1a;
    --dss-bg-secondary: #1e293b;
    --dss-bg-tertiary: #334155;
    --dss-bg-hover: #475569;

    --dss-text-primary: #f1f5f9;
    --dss-text-secondary: #94a3b8;
    --dss-text-muted: #64748b;

    --dss-border: #334155;
    --dss-border-hover: #475569;

    --dss-accent: #06b6d4;
    --dss-accent-hover: #22d3ee;
    --dss-success: #22c55e;
    --dss-warning: #f59e0b;
    --dss-danger: #ef4444;

    --dss-radius-sm: 4px;
    --dss-radius-md: 8px;
    --dss-radius-lg: 12px;

    --dss-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.3);

    font-family: var(--dss-font-sans);
    color: var(--dss-text-primary);
    line-height: 1.5;
  }

  :host([theme='light']) {
    --dss-bg-primary: #ffffff;
    --dss-bg-secondary: #f8fafc;
    --dss-bg-tertiary: #f1f5f9;
    --dss-bg-hover: #e2e8f0;

    --dss-text-primary: #0f172a;
    --dss-text-secondary: #475569;
    --dss-text-muted: #94a3b8;

    --dss-border: #e2e8f0;
    --dss-border-hover: #cbd5e1;

    --dss-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1);
  }

  * {
    box-sizing: border-box;
  }

  .section {
    background: var(--dss-bg-secondary);
    border: 1px solid var(--dss-border);
    border-radius: var(--dss-radius-lg);
    padding: 1.5rem;
    margin-bottom: 1.5rem;
  }

  .section-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 1rem;
    padding-bottom: 0.75rem;
    border-bottom: 1px solid var(--dss-border);
  }

  .section-title {
    font-size: 1.25rem;
    font-weight: 600;
    margin: 0;
    color: var(--dss-text-primary);
  }

  .section-subtitle {
    font-size: 0.875rem;
    color: var(--dss-text-secondary);
    margin: 0;
  }

  .badge {
    display: inline-flex;
    align-items: center;
    padding: 0.25rem 0.5rem;
    font-size: 0.75rem;
    font-weight: 500;
    border-radius: var(--dss-radius-sm);
    background: var(--dss-bg-tertiary);
    color: var(--dss-text-secondary);
  }

  .badge--accent {
    background: var(--dss-accent);
    color: var(--dss-bg-primary);
  }

  .grid {
    display: grid;
    gap: 1rem;
  }

  .grid--2 {
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  }

  .grid--3 {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  }

  .grid--4 {
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  }

  .card {
    background: var(--dss-bg-tertiary);
    border: 1px solid var(--dss-border);
    border-radius: var(--dss-radius-md);
    padding: 1rem;
    transition: border-color 0.2s, box-shadow 0.2s;
  }

  .card:hover {
    border-color: var(--dss-border-hover);
    box-shadow: var(--dss-shadow);
  }

  .card-title {
    font-size: 0.875rem;
    font-weight: 600;
    margin: 0 0 0.25rem;
  }

  .card-subtitle {
    font-size: 0.75rem;
    color: var(--dss-text-muted);
    margin: 0;
    font-family: var(--dss-font-mono);
  }

  .mono {
    font-family: var(--dss-font-mono);
  }

  .btn {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    font-size: 0.875rem;
    font-weight: 500;
    border: 1px solid var(--dss-border);
    border-radius: var(--dss-radius-md);
    background: var(--dss-bg-tertiary);
    color: var(--dss-text-primary);
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn:hover {
    background: var(--dss-bg-hover);
    border-color: var(--dss-border-hover);
  }

  .btn--primary {
    background: var(--dss-accent);
    border-color: var(--dss-accent);
    color: var(--dss-bg-primary);
  }

  .btn--primary:hover {
    background: var(--dss-accent-hover);
    border-color: var(--dss-accent-hover);
  }

  .empty-state {
    text-align: center;
    padding: 2rem;
    color: var(--dss-text-muted);
  }

  .tag-list {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-top: 0.5rem;
  }

  .tag {
    display: inline-block;
    padding: 0.125rem 0.5rem;
    font-size: 0.7rem;
    border-radius: var(--dss-radius-sm);
    background: var(--dss-bg-secondary);
    color: var(--dss-text-secondary);
  }
`;
