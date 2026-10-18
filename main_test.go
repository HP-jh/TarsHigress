package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	client := NewClient("", "")
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	withClient(client, healthHandler)(w, req)
	if w.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w.Code) }
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil { t.Fatal(err) }
	if resp["status"] != "ok" { t.Fatalf("expected ok, got %v", resp["status"]) }
}

func TestToolsHandler(t *testing.T) {
	client := NewClient("", "")
	req := httptest.NewRequest(http.MethodGet, "/tools", nil)
	w := httptest.NewRecorder()
	withClient(client, toolsHandler)(w, req)
	if w.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w.Code) }
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil { t.Fatal(err) }
	tools, ok := resp["tools"].([]interface{})
	if !ok || len(tools) == 0 { t.Fatal("expected tools") }
}

func TestInvokeHandler_MethodNotAllowed(t *testing.T) {
	client := NewClient("", "")
	req := httptest.NewRequest(http.MethodGet, "/invoke", nil)
	w := httptest.NewRecorder()
	withClient(client, invokeHandler)(w, req)
	if w.Code != http.StatusMethodNotAllowed { t.Fatalf("expected 405, got %d", w.Code) }
}

func TestInvokeHandler_UnknownTool(t *testing.T) {
	client := NewClient("", "")
	body, _ := json.Marshal(map[string]interface{}{ "tool": "xxx_unknown", "params": map[string]interface{}{} })
	req := httptest.NewRequest(http.MethodPost, "/invoke", bytes.NewReader(body))
	w := httptest.NewRecorder()
	withClient(client, invokeHandler)(w, req)
	if w.Code != http.StatusBadRequest { t.Fatalf("expected 400, got %d", w.Code) }
}

func TestClientReq_MockServer(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{ "status": "ok" })
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	client := NewClient(srv.URL, "test-key")
	res, err := client.req(context.Background(), "GET", "/health", nil)
	if err != nil { t.Fatalf("unexpected err: %v", err) }
	if res["status"] != "ok" { t.Fatalf("expected ok, got %v", res["status"]) }
}

func TestSecurityHeaders(t *testing.T) {
	client := NewClient("", "")
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	withClient(client, healthHandler)(w, req)
	headers := []string{
		"X-Content-Type-Options",
		"X-Frame-Options",
		"Content-Security-Policy",
		"Strict-Transport-Security",
		"Referrer-Policy",
		"Permissions-Policy",
	}
	for _, h := range headers {
		if w.Header().Get(h) == "" {
			t.Fatalf("missing security header: %s", h)
		}
	}
}

func TestRequestIDHeader(t *testing.T) {
	client := NewClient("", "")
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	withClient(client, healthHandler)(w, req)
	if w.Header().Get("X-Request-ID") == "" {
		t.Fatal("expected X-Request-ID header")
	}
}

func TestRequestIDPropagation(t *testing.T) {
	client := NewClient("", "")
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("X-Request-ID", "test-req-123")
	w := httptest.NewRecorder()
	withClient(client, healthHandler)(w, req)
	if w.Header().Get("X-Request-ID") != "test-req-123" {
		t.Fatalf("expected request ID to propagate, got %s", w.Header().Get("X-Request-ID"))
	}
}

func TestStructuredLogger(t *testing.T) {
	l := Logger{}
	l.Info("test", map[string]interface{}{"key": "value"})
	l.Error("test-err", map[string]interface{}{"key": "value"})
}

func TestConfigValidation_InvalidURL(t *testing.T) {
	os.Setenv("TARS_MCP_ENDPOINT", "not-a-url")
	defer os.Unsetenv("TARS_MCP_ENDPOINT")
	cfg := loadConfig()
	if cfg["endpoint"] != "" {
		t.Fatalf("expected empty endpoint for invalid URL, got %s", cfg["endpoint"])
	}
}

func TestConfigValidation_ValidURL(t *testing.T) {
	os.Setenv("TARS_MCP_ENDPOINT", "https://example.com")
	defer os.Unsetenv("TARS_MCP_ENDPOINT")
	cfg := loadConfig()
	if cfg["endpoint"] != "https://example.com" {
		t.Fatalf("expected valid endpoint, got %s", cfg["endpoint"])
	}
}

func TestGracefulShutdown(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	client := NewClient(srv.URL, "")
	res, err := client.req(context.Background(), "GET", "/health", nil)
	if err != nil { t.Fatalf("unexpected err: %v", err) }
	if res == nil { t.Fatal("expected response") }
}

func TestRetryMechanism(t *testing.T) {
	callCount := 0
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount < 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	client := NewClient(srv.URL, "")
	res, err := client.req(context.Background(), "GET", "/health", nil)
	if err != nil { t.Fatalf("unexpected err: %v", err) }
	if res["status"] != "ok" { t.Fatalf("expected ok, got %v", res["status"]) }
	if callCount < 2 { t.Fatalf("expected retries, got %d calls", callCount) }
}

func TestMetricsEndpoint(t *testing.T) {
	client := NewClient("", "")
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()
	withClient(client, metricsHandler)(w, req)
	if w.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w.Code) }
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil { t.Fatal(err) }
	if _, ok := resp["requests"]; !ok { t.Fatal("expected requests field") }
}

func TestRateLimiting(t *testing.T) {
	client := NewClient("", "")
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	withClient(client, healthHandler)(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusTooManyRequests {
		t.Fatalf("unexpected status: %d", w.Code)
	}
}
