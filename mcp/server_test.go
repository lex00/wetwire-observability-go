package mcp

import (
	"encoding/json"
	"testing"
)

func TestNewServer(t *testing.T) {
	s := NewServer("wetwire-obs", "1.0.0")
	if s.name != "wetwire-obs" {
		t.Errorf("name = %q, want wetwire-obs", s.name)
	}
	if s.version != "1.0.0" {
		t.Errorf("version = %q, want 1.0.0", s.version)
	}
}

func TestServer_RegisterTool(t *testing.T) {
	s := NewServer("test", "1.0.0")
	tool := Tool{
		Name:        "test-tool",
		Description: "A test tool",
	}
	s.RegisterTool(tool, func(params json.RawMessage) (any, error) {
		return map[string]string{"result": "ok"}, nil
	})

	if len(s.tools) != 1 {
		t.Errorf("len(tools) = %d, want 1", len(s.tools))
	}
	if s.tools[0].Name != "test-tool" {
		t.Errorf("tools[0].Name = %q", s.tools[0].Name)
	}
}

func TestServer_HandleInitialize(t *testing.T) {
	s := NewServer("wetwire-obs", "1.0.0")

	req := Request{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`1`),
		Method:  "initialize",
		Params: json.RawMessage(`{
			"protocolVersion": "2024-11-05",
			"capabilities": {},
			"clientInfo": {"name": "test-client", "version": "1.0.0"}
		}`),
	}

	resp := s.HandleRequest(req)
	if resp.Error != nil {
		t.Fatalf("HandleRequest error: %v", resp.Error)
	}

	result, ok := resp.Result.(InitializeResult)
	if !ok {
		t.Fatalf("result type = %T, want InitializeResult", resp.Result)
	}

	if result.ServerInfo.Name != "wetwire-obs" {
		t.Errorf("ServerInfo.Name = %q", result.ServerInfo.Name)
	}
}

func TestServer_HandleToolsList(t *testing.T) {
	s := NewServer("test", "1.0.0")
	s.RegisterTool(Tool{
		Name:        "build",
		Description: "Build configurations",
	}, func(params json.RawMessage) (any, error) {
		return nil, nil
	})

	req := Request{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`2`),
		Method:  "tools/list",
	}

	resp := s.HandleRequest(req)
	if resp.Error != nil {
		t.Fatalf("HandleRequest error: %v", resp.Error)
	}

	result, ok := resp.Result.(ToolsListResult)
	if !ok {
		t.Fatalf("result type = %T, want ToolsListResult", resp.Result)
	}

	if len(result.Tools) != 1 {
		t.Errorf("len(tools) = %d, want 1", len(result.Tools))
	}
}

func TestServer_HandleToolsCall(t *testing.T) {
	s := NewServer("test", "1.0.0")
	s.RegisterTool(Tool{
		Name:        "echo",
		Description: "Echo input",
	}, func(params json.RawMessage) (any, error) {
		var p struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		return map[string]string{"echo": p.Message}, nil
	})

	req := Request{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`3`),
		Method:  "tools/call",
		Params: json.RawMessage(`{
			"name": "echo",
			"arguments": {"message": "hello"}
		}`),
	}

	resp := s.HandleRequest(req)
	if resp.Error != nil {
		t.Fatalf("HandleRequest error: %v", resp.Error)
	}

	result, ok := resp.Result.(*ToolCallResult)
	if !ok {
		t.Fatalf("result type = %T, want *ToolCallResult", resp.Result)
	}

	if len(result.Content) == 0 {
		t.Error("expected content in result")
	}
}

func TestServer_HandleUnknownMethod(t *testing.T) {
	s := NewServer("test", "1.0.0")

	req := Request{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`4`),
		Method:  "unknown/method",
	}

	resp := s.HandleRequest(req)
	if resp.Error == nil {
		t.Error("expected error for unknown method")
	}
	if resp.Error.Code != -32601 {
		t.Errorf("error code = %d, want -32601", resp.Error.Code)
	}
}

func TestServer_HandleToolNotFound(t *testing.T) {
	s := NewServer("test", "1.0.0")

	req := Request{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`5`),
		Method:  "tools/call",
		Params: json.RawMessage(`{
			"name": "nonexistent",
			"arguments": {}
		}`),
	}

	resp := s.HandleRequest(req)
	if resp.Error == nil {
		t.Error("expected error for nonexistent tool")
	}
}
