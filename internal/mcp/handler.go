package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"mdreview/internal/fs"
	"mdreview/internal/render"
	"mdreview/internal/server"

	"github.com/google/uuid"
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
	port          int
	store         *server.Store
}

func NewMCPServer(workspacePath string, port int, store *server.Store) *MCPServer {
	return &MCPServer{
		workspacePath: workspacePath,
		port:          port,
		store:         store,
	}
}

func (s *MCPServer) Run() {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			}
			break
		}

		var req JSONRPCRequest
		if err := json.Unmarshal(line, &req); err != nil {
			fmt.Fprintf(os.Stderr, "Error unmarshaling request: %v\n", err)
			continue
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
					"version": "0.1.0",
				},
			}
		case "tools/list":
			resp.Result = ListToolsResult{
				Tools: []Tool{
					{
						Name:        "preview_markdown",
						Description: "Render a markdown file to HTML and return a preview URL.",
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
			// Ignore other methods for now
			continue
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

	content, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	html, err := render.RenderMarkdown(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to render markdown: %w", err)
	}

	id := uuid.New().String()
	s.store.Set(id, html)

	ip := GetTailscaleIP()
	url := fmt.Sprintf("http://%s:%d/rev/%s", ip, s.port, id)

	return &CallToolResult{
		Content: []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		}{
			{Type: "text", Text: fmt.Sprintf("Preview available at: %s", url)},
		},
	}, nil
}
