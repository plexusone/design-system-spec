package visual

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"time"
)

// CaptureOptions configures screenshot capture.
type CaptureOptions struct {
	URL           string
	Selector      string
	Viewport      Viewport
	Stabilization *Stabilization
}

// W3PilotClient wraps communication with w3pilot via MCP.
type W3PilotClient struct {
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	decoder *json.Decoder
	reqID   int
}

// MCPRequest represents an MCP JSON-RPC request.
type MCPRequest struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      int            `json:"id"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params,omitempty"`
}

// MCPResponse represents an MCP JSON-RPC response.
type MCPResponse struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      int            `json:"id"`
	Result  map[string]any `json:"result,omitempty"`
	Error   *MCPError      `json:"error,omitempty"`
}

// MCPError represents an MCP error.
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewW3PilotClient starts w3pilot as a subprocess and connects via MCP.
func NewW3PilotClient(ctx context.Context) (*W3PilotClient, error) {
	// Check if w3pilot is available
	w3pilotPath, err := exec.LookPath("w3pilot")
	if err != nil {
		return nil, fmt.Errorf("%w: w3pilot not found in PATH", ErrW3PilotUnavailable)
	}

	// Start w3pilot MCP server
	cmd := exec.CommandContext(ctx, w3pilotPath, "mcp", "serve")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start w3pilot: %w", err)
	}

	client := &W3PilotClient{
		cmd:     cmd,
		stdin:   stdin,
		stdout:  stdout,
		decoder: json.NewDecoder(stdout),
		reqID:   0,
	}

	// Initialize MCP connection
	if err := client.initialize(ctx); err != nil {
		client.Close()
		return nil, err
	}

	return client, nil
}

// initialize performs MCP handshake.
func (c *W3PilotClient) initialize(ctx context.Context) error {
	// Send initialize request
	_, err := c.call(ctx, "initialize", map[string]any{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]any{},
		"clientInfo": map[string]any{
			"name":    "dss-visual",
			"version": "1.0.0",
		},
	})
	if err != nil {
		return fmt.Errorf("MCP initialize failed: %w", err)
	}

	// Send initialized notification
	return c.notify(ctx, "notifications/initialized", nil)
}

// call sends an MCP request and waits for response.
func (c *W3PilotClient) call(ctx context.Context, method string, params map[string]any) (map[string]any, error) {
	_ = ctx // Reserved for future cancellation support
	c.reqID++

	req := MCPRequest{
		JSONRPC: "2.0",
		ID:      c.reqID,
		Method:  method,
		Params:  params,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Write request
	if _, err := c.stdin.Write(append(data, '\n')); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Read response
	var resp MCPResponse
	if err := c.decoder.Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("MCP error %d: %s", resp.Error.Code, resp.Error.Message)
	}

	return resp.Result, nil
}

// notify sends an MCP notification (no response expected).
func (c *W3PilotClient) notify(ctx context.Context, method string, params map[string]any) error {
	_ = ctx // Reserved for future cancellation support
	req := map[string]any{
		"jsonrpc": "2.0",
		"method":  method,
	}
	if params != nil {
		req["params"] = params
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	_, err = c.stdin.Write(append(data, '\n'))
	return err
}

// callTool calls an MCP tool.
func (c *W3PilotClient) callTool(ctx context.Context, name string, args map[string]any) (map[string]any, error) {
	result, err := c.call(ctx, "tools/call", map[string]any{
		"name":      name,
		"arguments": args,
	})
	if err != nil {
		return nil, err
	}

	// Extract content from result
	if content, ok := result["content"].([]any); ok && len(content) > 0 {
		if item, ok := content[0].(map[string]any); ok {
			if text, ok := item["text"].(string); ok {
				// Parse JSON text
				var parsed map[string]any
				if err := json.Unmarshal([]byte(text), &parsed); err == nil {
					return parsed, nil
				}
				// Return as-is if not JSON
				return map[string]any{"text": text}, nil
			}
		}
	}

	return result, nil
}

// LaunchBrowser starts a browser session.
func (c *W3PilotClient) LaunchBrowser(ctx context.Context, headless bool) error {
	_, err := c.callTool(ctx, "browser_launch", map[string]any{
		"headless": headless,
	})
	return err
}

// CloseBrowser closes the browser session.
func (c *W3PilotClient) CloseBrowser(ctx context.Context) error {
	_, err := c.callTool(ctx, "browser_close", nil)
	return err
}

// Navigate navigates to a URL.
func (c *W3PilotClient) Navigate(ctx context.Context, url string) error {
	_, err := c.callTool(ctx, "page_navigate", map[string]any{
		"url": url,
	})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrNavigationFailed, err)
	}
	return nil
}

// SetViewport sets the viewport size.
func (c *W3PilotClient) SetViewport(ctx context.Context, width, height int) error {
	_, err := c.callTool(ctx, "page_set_viewport", map[string]any{
		"width":  width,
		"height": height,
	})
	return err
}

// WaitForSelector waits for a selector to appear.
func (c *W3PilotClient) WaitForSelector(ctx context.Context, selector string, timeoutMs int) error {
	_, err := c.callTool(ctx, "page_wait_for_selector", map[string]any{
		"selector":   selector,
		"timeout_ms": timeoutMs,
	})
	if err != nil {
		return fmt.Errorf("%w: selector %s: %v", ErrStabilizationFailed, selector, err)
	}
	return nil
}

// Evaluate runs JavaScript in the page.
func (c *W3PilotClient) Evaluate(ctx context.Context, expression string) error {
	_, err := c.callTool(ctx, "page_evaluate", map[string]any{
		"expression": expression,
	})
	return err
}

// PageScreenshot captures a full page screenshot.
func (c *W3PilotClient) PageScreenshot(ctx context.Context) ([]byte, error) {
	result, err := c.callTool(ctx, "page_screenshot", map[string]any{
		"format": "base64",
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrScreenshotFailed, err)
	}

	data, ok := result["data"].(string)
	if !ok {
		return nil, fmt.Errorf("%w: invalid screenshot response", ErrScreenshotFailed)
	}

	return base64.StdEncoding.DecodeString(data)
}

// ElementScreenshot captures an element screenshot.
func (c *W3PilotClient) ElementScreenshot(ctx context.Context, selector string) ([]byte, error) {
	result, err := c.callTool(ctx, "element_screenshot", map[string]any{
		"selector": selector,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrScreenshotFailed, err)
	}

	data, ok := result["data"].(string)
	if !ok {
		return nil, fmt.Errorf("%w: invalid screenshot response", ErrScreenshotFailed)
	}

	return base64.StdEncoding.DecodeString(data)
}

// CaptureScreenshot navigates to URL and captures a screenshot.
func (c *W3PilotClient) CaptureScreenshot(ctx context.Context, opts CaptureOptions) ([]byte, error) {
	// 1. Navigate
	if err := c.Navigate(ctx, opts.URL); err != nil {
		return nil, err
	}

	// 2. Set viewport
	if err := c.SetViewport(ctx, opts.Viewport.Width, opts.Viewport.Height); err != nil {
		return nil, fmt.Errorf("failed to set viewport: %w", err)
	}

	// 3. Stabilization
	if opts.Stabilization != nil {
		if err := c.stabilize(ctx, opts.Stabilization); err != nil {
			return nil, err
		}
	}

	// 4. Capture screenshot
	if opts.Selector != "" {
		return c.ElementScreenshot(ctx, opts.Selector)
	}
	return c.PageScreenshot(ctx)
}

// stabilize applies stabilization settings.
func (c *W3PilotClient) stabilize(ctx context.Context, s *Stabilization) error {
	// Wait for selector
	if s.WaitForSelector != "" {
		timeout := s.WaitForTimeout
		if timeout == 0 {
			timeout = 5000
		}
		if err := c.WaitForSelector(ctx, s.WaitForSelector, timeout); err != nil {
			return err
		}
	}

	// Fixed wait
	if s.WaitMs > 0 {
		time.Sleep(time.Duration(s.WaitMs) * time.Millisecond)
	}

	// Disable animations
	if s.DisableAnimations {
		err := c.Evaluate(ctx, `
			const style = document.createElement('style');
			style.textContent = '*, *::before, *::after { animation: none !important; transition: none !important; }';
			document.head.appendChild(style);
		`)
		if err != nil {
			// Non-fatal, continue
		}
	}

	return nil
}

// Close terminates the w3pilot subprocess.
func (c *W3PilotClient) Close() error {
	if c.stdin != nil {
		c.stdin.Close()
	}
	if c.stdout != nil {
		c.stdout.Close()
	}
	if c.cmd != nil && c.cmd.Process != nil {
		return c.cmd.Process.Kill()
	}
	return nil
}
