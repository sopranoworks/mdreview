package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/sopranoworks/mdreview/internal/fs"
	"github.com/sopranoworks/mdreview/internal/version"
)

type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

type ListToolsResult struct {
	Tools []Tool `json:"tools"`
}

type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type CallToolResult struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}

type MCPServer struct {
	workspacePath string
	sidecarPort   int
	sidecarToken  string
}

func NewMCPServer(workspacePath string) *MCPServer {
	return &MCPServer{
		workspacePath: workspacePath,
	}
}

func (s *MCPServer) WithSidecar(port int, token string) *MCPServer {
	s.sidecarPort = port
	s.sidecarToken = token
	return s
}

func (s *MCPServer) Run() {
	decoder := json.NewDecoder(os.Stdin)
	for {
		var req JSONRPCRequest
		if err := decoder.Decode(&req); err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "Error decoding request: %v\n", err)
			}
			break
		}

		var resp JSONRPCResponse
		resp.JSONRPC = "2.0"
		resp.ID = req.ID

		switch req.Method {
		case "initialize":
			resp.Result = map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities":    map[string]interface{}{},
				"serverInfo": map[string]string{
					"name":    "mdreview",
					"version": version.Version,
				},
			}
		case "ping":
			resp.Result = map[string]interface{}{}
		case "tools/list":
			resp.Result = ListToolsResult{
				Tools: []Tool{
					{
						Name:        "preview_markdown",
						Description: "Render a markdown file to HTML and return a preview URL. You MUST call this tool immediately after creating or modifying any Markdown (.md) file to provide the user with a preview link for their review.",
						InputSchema: map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"path": map[string]interface{}{
									"type":        "string",
									"description": "The path to the markdown file relative to the workspace root.",
								},
							},
							"required": []string{"path"},
						},
					},
				},
			}
		case "tools/call":
			var params CallToolParams
			if err := json.Unmarshal(req.Params, &params); err != nil {
				fmt.Fprintf(os.Stderr, "Error unmarshaling tool params: %v\n", err)
				resp.Error = map[string]interface{}{
					"code":    -32602,
					"message": "Invalid params",
				}
			} else if params.Name == "preview_markdown" {
				path, ok := params.Arguments["path"].(string)
				if !ok {
					resp.Error = map[string]interface{}{
						"code":    -32602,
						"message": "Missing path argument",
					}
				} else {
					result, err := s.handlePreviewMarkdown(path)
					if err != nil {
						resp.Result = CallToolResult{
							Content: []struct {
								Type string `json:"type"`
								Text string `json:"text"`
							}{
								{Type: "text", Text: fmt.Sprintf("Error: %v", err)},
							},
						}
					} else {
						resp.Result = result
					}
				}
			} else {
				resp.Error = map[string]interface{}{
					"code":    -32601,
					"message": "Method not found",
				}
			}
		default:
			if req.ID == nil {
				continue // Ignore unhandled notifications
			}
			resp.Error = map[string]interface{}{
				"code":    -32601,
				"message": "Method not found",
			}
		}

		out, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling response: %v\n", err)
			continue
		}
		os.Stdout.Write(out)
		os.Stdout.Write([]byte("\n"))
	}
}

func (s *MCPServer) handlePreviewMarkdown(relPath string) (*CallToolResult, error) {
	absPath, err := fs.ValidatePath(s.workspacePath, relPath)
	if err != nil {
		return nil, err
	}

	if s.sidecarPort == 0 {
		port, token, err := DiscoverOrStartServer(0)
		if err != nil {
			return nil, fmt.Errorf("failed to discover or start sidecar server: %w", err)
		}
		s.sidecarPort = port
		s.sidecarToken = token
	}

	url, err := ProxyPreviewRequest(s.sidecarPort, s.sidecarToken, absPath)
	if err != nil {
		return nil, err
	}

	return &CallToolResult{
		Content: []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		}{
			{Type: "text", Text: fmt.Sprintf("Preview available at: %s", url)},
		},
	}, nil
}
