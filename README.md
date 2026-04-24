# mdreview

![Version](https://img.shields.io/badge/version-v0.1.0-blue)

`mdreview` is an MCP (Model Context Protocol) server that allows AI agents to render Markdown files to sanitized HTML and provide a temporary preview URL. It is specifically designed to work across devices using Tailscale.

## Current Version: v0.1.0
This initial release includes the core MCP server, HTML rendering with `goldmark`, XSS sanitization, and Tailscale IP detection.

## Features

- **Markdown Rendering:** Powered by `goldmark` for high-fidelity HTML output.
- **XSS Protection:** Strict sanitization using `bluemonday`.
- **Path Validation:** Prevents path traversal by restricting access to a configured workspace.
- **Tailscale Integration:** Automatically detects your Tailscale IP to return URLs accessible throughout your Tailnet.
- **In-Memory Store:** Previews are stored temporarily in memory with unique UUIDs.

## Installation

### 1. Build the Binary

Ensure you have Go installed (version 1.23+ recommended), then run:

```bash
go build -o mdreview
```

This will create a `mdreview` executable in the root directory.

### 2. Register with Gemini CLI

To make `mdreview` persistent and available in all projects, you can add it to your global configuration:

#### Option A: Global MCP Server (Recommended)
This registers the server directly in your `~/.gemini/settings.json`:

```bash
gemini mcp add --scope user mdreview ./mdreview -port 8080 -workspace .
```

#### Option B: Link as an Extension
If you want to use the `plugin.json` manifest (for agent integration/hooks), link this directory:

```bash
gemini extensions link .
```

### 3. Verify Installation

Restart your Gemini CLI session and check the connected servers:

- Run `/mcp list` to see `mdreview` and the `preview_markdown` tool.
- Run `gemini extensions list` (if using Option B) to verify the extension is active.

## Usage

Once connected, an agent can call the `preview_markdown` tool:

- **Input:** `path` (relative or absolute path to a `.md` file).
- **Output:** A URL (e.g., `http://100.x.y.z:8080/rev/uuid`) that can be opened in any browser on your network.

## Configuration

The server can be configured via `plugin.json` or CLI flags:

- `-port`: The port for the HTTP side-car (default: `8080`).
- `-workspace`: The root directory for file access validation (default: `.`).

## Security

This tool is built with a "security-first" approach:
- **Symlink resolution:** Prevents escaping the workspace via symbolic links.
- **Strict HTML Sanitization:** Blocks malicious scripts and event handlers.
- **Random UUIDs:** Preview paths are unguessable.
