# mdreview

![Version](https://img.shields.io/badge/version-v0.2.0-green)

`mdreview` is an MCP (Model Context Protocol) server that allows AI agents to render Markdown files to sanitized HTML and provide a temporary preview URL. It is specifically designed to work across devices using Tailscale.

## Current Version: v0.2.0 (Split Sidecar Edition)
This major update introduces a persistent architecture and significant security hardening.

### What's New in v0.2.0:
- **Split Sidecar Architecture**: Separation into `mdreview-mcp-srv` (persistent background server) and `mdreview-mcp-cli` (transient MCP handler).
- **Session Persistence**: Previews now stay alive even after the AI agent session ends.
- **Security Hardening**:
  - **Token-based IPC**: Communication between CLI and Server is secured with random per-session tokens.
  - **Salted Preview IDs**: URLs are now unpredictable, generated using a combination of the IPC token and file path.
  - **Least Privilege**: State files are restricted to `0600` permissions and are user-specific.
- **Improved Installation**: Direct `go install` of specific binaries.

## Features

- **Markdown Rendering:** Powered by `goldmark` for high-fidelity HTML output.
- **XSS Protection:** Strict sanitization using `bluemonday`.
- **Path Validation:** Prevents path traversal by restricting access to a configured workspace.
- **Tailscale Integration:** Automatically detects your Tailscale IP to return URLs accessible throughout your Tailnet.
- **Sidecar Architecture:** Uses a persistent background server (`mdreview-mcp-srv`) to manage previews independently of the MCP session.
- **Auto-Expiration:** Previews are stored temporarily in memory and automatically expire after 24 hours.

## Installation

`mdreview` supports **zero-setup execution**. If you have Go installed (version 1.23+), you can install the binaries with a single command:

```bash
go install github.com/sopranoworks/mdreview/cmd/...@latest
```

This will install `mdreview-mcp-cli` and `mdreview-mcp-srv` to your `$GOPATH/bin`.

### 1. Register with Gemini CLI

#### Option A: Install as an Extension (Recommended)
This is the most complete method. It installs the `preview_markdown` tool and the **AI Agent personality** from `agents/mdreview.md`. The agent is specifically tuned to be proactive and will automatically show you previews after every edit.

```bash
# 1. Install binaries
go install github.com/sopranoworks/mdreview/cmd/...@latest

# 2. Install extension
gemini extensions install https://github.com/sopranoworks/mdreview
```

#### Option B: Standalone MCP Server
This method adds the tool to your toolbox. The tool's description includes embedded instructions advising the agent to use it immediately after edits.

```bash
gemini mcp add --scope user mdreview mdreview-mcp-cli -port 8080 -workspace .
```

### 2. Register with Claude Code

#### Option A: Recommended Registration
Ensure you have installed the binaries via `go install`, then add the tool:

```bash
claude mcp add mdreview mdreview-mcp-cli -- -port 8080 -workspace .
```

**Proactive Previews:** The tool itself contains instructions for Claude to use it after every Markdown edit. For the full experience, you can copy the specific instructions from [agents/mdreview.md](agents/mdreview.md) into your Claude custom instructions.

#### Option B: Direct Execution (No Pre-install)
```bash
claude mcp add mdreview go -- run github.com/sopranoworks/mdreview/cmd/mdreview-mcp-cli@latest -port 8080 -workspace .
```

### 3. (Optional) Manual Build
If you'd like to build the binaries manually from the source:

```bash
go build -o mdreview-mcp-cli ./cmd/mdreview-mcp-cli
go build -o mdreview-mcp-srv ./cmd/mdreview-mcp-srv
```

### 4. Verify Installation


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
