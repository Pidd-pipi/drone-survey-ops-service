package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDomainHeaderPresent(t *testing.T) {
	handler := opsEnterpriseMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	srv := httptest.NewServer(handler)
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/x")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("状态码错误: %d", resp.StatusCode)
	}
	if got := resp.Header.Get("X-Operations-Domain"); got != "drone-survey-ops-service" {
		t.Fatalf("缺少 X-Operations-Domain 头: %q", got)
	}
	if got := resp.Header.Get("X-Operations-Request"); got == "" {
		t.Fatal("缺少 X-Operations-Request 头")
	}
}

func TestJSONContentTypeSet(t *testing.T) {
	rec := httptest.NewRecorder()
	opsJSON(rec, http.StatusOK, map[string]string{"ok": "1"})
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("opsJSON 应设置 JSON Content-Type，实际 %q", ct)
	}
}

func TestAllowedMethodsMultiple(t *testing.T) {
	if !opsAllowed("GET", "GET", "POST") {
		t.Fatal("GET 应被允许")
	}
	if !opsAllowed("POST", "GET", "POST") {
		t.Fatal("POST 应被允许")
	}
	if opsAllowed("DELETE", "GET", "POST") {
		t.Fatal("DELETE 不应被允许")
	}
}

func TestHealthServiceName(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(healthHandler))
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["service"] != "drone-survey-ops" {
		t.Fatalf("健康检查服务名错误: %q", body["service"])
	}
}
