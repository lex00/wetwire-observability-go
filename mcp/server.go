// Package mcp provides a Model Context Protocol (MCP) server implementation.
package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Server is an MCP server that handles JSON-RPC requests over stdio.
type Server struct {
	name     string
	version  string
	tools    []Tool
	handlers map[string]ToolHandler
}

// Tool represents an MCP tool definition.
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	InputSchema InputSchema `json:"inputSchema,omitempty"`
}

// InputSchema defines the JSON Schema for tool inputs.
type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
}

// Property defines a JSON Schema property.
type Property struct {
	Type        string   `json:"type"`
	Description string   `json:"description,omitempty"`
	Enum        []string `json:"enum,omitempty"`
	Default     any      `json:"default,omitempty"`
}

// ToolHandler is a function that handles a tool call.
type ToolHandler func(params json.RawMessage) (any, error)

// Request is a JSON-RPC request.
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response is a JSON-RPC response.
type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  any             `json:"result,omitempty"`
	Error   *ResponseError  `json:"error,omitempty"`
}

// ResponseError is a JSON-RPC error.
type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// MCP Protocol types

// InitializeParams contains the parameters for initialize request.
type InitializeParams struct {
	ProtocolVersion string          `json:"protocolVersion"`
	Capabilities    ClientCaps      `json:"capabilities"`
	ClientInfo      Implementation  `json:"clientInfo"`
}

// ClientCaps contains client capabilities.
type ClientCaps struct {
	Roots    *RootsCap    `json:"roots,omitempty"`
	Sampling *SamplingCap `json:"sampling,omitempty"`
}

// RootsCap indicates roots capability.
type RootsCap struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// SamplingCap indicates sampling capability.
type SamplingCap struct{}

// Implementation identifies a client or server.
type Implementation struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// InitializeResult is the response to initialize.
type InitializeResult struct {
	ProtocolVersion string         `json:"protocolVersion"`
	Capabilities    ServerCaps     `json:"capabilities"`
	ServerInfo      Implementation `json:"serverInfo"`
}

// ServerCaps contains server capabilities.
type ServerCaps struct {
	Tools *ToolsCap `json:"tools,omitempty"`
}

// ToolsCap indicates tools capability.
type ToolsCap struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// ToolsListResult is the response to tools/list.
type ToolsListResult struct {
	Tools []Tool `json:"tools"`
}

// ToolsCallParams contains the parameters for tools/call.
type ToolsCallParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// ToolCallResult is the response to tools/call.
type ToolCallResult struct {
	Content []ContentBlock `json:"content"`
	IsError bool           `json:"isError,omitempty"`
}

// ContentBlock is a content block in a tool result.
type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// JSON-RPC error codes
const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)

// NewServer creates a new MCP server.
func NewServer(name, version string) *Server {
	return &Server{
		name:     name,
		version:  version,
		handlers: make(map[string]ToolHandler),
	}
}

// RegisterTool registers a tool with its handler.
func (s *Server) RegisterTool(tool Tool, handler ToolHandler) {
	s.tools = append(s.tools, tool)
	s.handlers[tool.Name] = handler
}

// HandleRequest processes a single JSON-RPC request.
func (s *Server) HandleRequest(req Request) Response {
	resp := Response{
		JSONRPC: "2.0",
		ID:      req.ID,
	}

	switch req.Method {
	case "initialize":
		resp.Result = s.handleInitialize(req.Params)
	case "initialized":
		// Notification, no response needed for notifications
		// But since we return Response, just set empty result
		resp.Result = struct{}{}
	case "tools/list":
		resp.Result = s.handleToolsList()
	case "tools/call":
		result, err := s.handleToolsCall(req.Params)
		if err != nil {
			resp.Error = &ResponseError{
				Code:    InternalError,
				Message: err.Error(),
			}
		} else {
			resp.Result = result
		}
	default:
		resp.Error = &ResponseError{
			Code:    MethodNotFound,
			Message: fmt.Sprintf("unknown method: %s", req.Method),
		}
	}

	return resp
}

func (s *Server) handleInitialize(params json.RawMessage) InitializeResult {
	return InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: ServerCaps{
			Tools: &ToolsCap{
				ListChanged: false,
			},
		},
		ServerInfo: Implementation{
			Name:    s.name,
			Version: s.version,
		},
	}
}

func (s *Server) handleToolsList() ToolsListResult {
	return ToolsListResult{
		Tools: s.tools,
	}
}

func (s *Server) handleToolsCall(params json.RawMessage) (*ToolCallResult, error) {
	var callParams ToolsCallParams
	if err := json.Unmarshal(params, &callParams); err != nil {
		return nil, fmt.Errorf("invalid params: %w", err)
	}

	handler, ok := s.handlers[callParams.Name]
	if !ok {
		return nil, fmt.Errorf("unknown tool: %s", callParams.Name)
	}

	result, err := handler(callParams.Arguments)
	if err != nil {
		return &ToolCallResult{
			Content: []ContentBlock{
				{Type: "text", Text: fmt.Sprintf("Error: %s", err.Error())},
			},
			IsError: true,
		}, nil
	}

	// Format result as text
	var text string
	switch v := result.(type) {
	case string:
		text = v
	case nil:
		text = "Success"
	default:
		data, _ := json.MarshalIndent(result, "", "  ")
		text = string(data)
	}

	return &ToolCallResult{
		Content: []ContentBlock{
			{Type: "text", Text: text},
		},
	}, nil
}

// Run starts the MCP server, reading from stdin and writing to stdout.
func (s *Server) Run() error {
	reader := bufio.NewReader(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			resp := Response{
				JSONRPC: "2.0",
				Error: &ResponseError{
					Code:    ParseError,
					Message: "parse error",
				},
			}
			encoder.Encode(resp)
			continue
		}

		resp := s.HandleRequest(req)

		// Don't send response for notifications (requests without ID)
		if req.ID == nil {
			continue
		}

		if err := encoder.Encode(resp); err != nil {
			return err
		}
	}
}
