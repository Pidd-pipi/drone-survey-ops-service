package main

import (
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestAuditAddReturnsEvent(t *testing.T) {
	audit := newOpsAudit()
	event := audit.Add("r-001", "created", "lin")
	if event.ID == "" {
		t.Fatal("Add 返回的事件 ID 为空")
	}
	if event.RecordID != "r-001" || event.Type != "created" || event.Actor != "lin" {
		t.Fatalf("Add 返回的事件字段不对: %+v", event)
	}
}

func TestAuditClearReleasesBacking(t *testing.T) {
	audit := newOpsAudit()
	for i := 0; i < 50; i++ {
		audit.Add("r-001", "created", "lin")
	}
	audit.Clear()
	if len(audit.events) != 0 {
		t.Fatalf("Clear 后事件数应为 0，实际 %d", len(audit.events))
	}
	if cap(audit.events) != 0 {
		t.Fatalf("Clear 后底层数组未释放，cap=%d", cap(audit.events))
	}
}

func TestRequestIDsUniqueUnderConcurrency(t *testing.T) {
	handler := requestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	const workers = 8
	const rounds = 50
	start := make(chan struct{})
	var wg sync.WaitGroup
	ids := make(chan string, workers*rounds)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < rounds; j++ {
				rec := httptest.NewRecorder()
				handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))
				ids <- rec.Header().Get("X-Request-ID")
			}
		}()
	}
	close(start)
	wg.Wait()
	close(ids)
	seen := map[string]bool{}
	for id := range ids {
		if id == "" {
			t.Fatal("请求 ID 为空")
		}
		if seen[id] {
			t.Fatalf("请求 ID 重复: %s", id)
		}
		seen[id] = true
	}
}

func runShutdownCycle(t *testing.T) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	if err := ln.Close(); err != nil {
		t.Fatal(err)
	}
	srv := newEnterpriseServer(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	done := make(chan struct{})
	go func() {
		_ = serveHTTP(srv)
		close(done)
	}()
	deadline := time.Now().Add(2 * time.Second)
	for {
		conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("server 未开始监听")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("serveHTTP 未在信号后返回")
	}
}

func TestServeHTTPNoGoroutineLeakAfterShutdown(t *testing.T) {
	runShutdownCycle(t)
	time.Sleep(200 * time.Millisecond)
	before := runtime.NumGoroutine()
	for i := 0; i < 3; i++ {
		runShutdownCycle(t)
	}
	time.Sleep(300 * time.Millisecond)
	after := runtime.NumGoroutine()
	if after > before+1 {
		t.Fatalf("关停后 goroutine 数增长: before=%d after=%d", before, after)
	}
}
