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

`mdreview` supports **zero-setup execution**. If you have Go installed (version 1.23+), you can simply install and it will "just work" everywhere—Windows, macOS, and Linux.

### 1. Register with Gemini CLI

To make `mdreview` persistent and available in all projects:

#### Option A: Global MCP Server (Recommended)
This registers the server directly in your `~/.gemini/settings.json`:

```bash
gemini mcp add --scope user mdreview go run bootstrap.go -port 8080 -workspace .
```

#### Option B: Link as an Extension
Link this directory to enable the plugin manifest and hooks:

```bash
gemini extensions link .
```

#### Option C: Claude Code (Plugin)
Claude Code supports direct plugin installation. Installation is **Auto-Build** out of the box:

```bash
# Install and it just works
claude plugin add github:sopranoworks/mdreview
```

### 2. (Optional) Manual Build
The server will automatically build itself on the first run. However, if you'd like to build the binary manually:

```bash
go build -o mdreview
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

The server can be configured via flags or environment variables:

- `-port`: The **preferred** port for the HTTP side-car (default: `8080`).
- `-workspace`: The root directory for file access validation (default: `.`).

**Automatic Port Selection:**
If the requested port (e.g., 8080) is already in use by another application, `mdreview` will automatically search for and bind to an available port. It will log the actual port used to `stderr`, and the `preview_markdown` tool will always return URLs with the correct, active port. You do not need to manually change the port unless you have a specific requirement.

## Security

This tool is built with a "security-first" approach:
- **Symlink resolution:** Prevents escaping the workspace via symbolic links.
- **Strict HTML Sanitization:** Blocks malicious scripts and event handlers.
- **Random UUIDs:** Preview paths are unguessable.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

Copyright (c) 2026 Sopranoworks, Osamu Takahashi.
