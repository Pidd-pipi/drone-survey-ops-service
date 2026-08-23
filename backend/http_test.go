package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPWorkflow(t *testing.T) {
	ts := httptest.NewServer(NewRouter(NewSurveyService(NewMissionStore())))
	defer ts.Close()
	for _, path := range []string{"/health", "/api/missions", "/", "/app.js"} {
		resp, err := http.Get(ts.URL + path)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("GET %s status=%v err=%v", path, resp.StatusCode, err)
		}
		resp.Body.Close()
	}
	for _, tc := range []struct {
		body string
		want int
	}{{`{"id":"mission-241","status":"landed"}`, 200}, {`{"id":"mission-241","status":"unknown"}`, 400}, {`{"id":"mission-x","status":"landed"}`, 404}} {
		resp, err := http.Post(ts.URL+"/api/missions/status", "application/json", bytes.NewBufferString(tc.body))
		if err != nil || resp.StatusCode != tc.want {
			t.Fatalf("POST status=%v want=%v err=%v", resp.StatusCode, tc.want, err)
		}
		resp.Body.Close()
	}
}
