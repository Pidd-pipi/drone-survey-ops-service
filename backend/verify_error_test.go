package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func seedErrorRecord(id string) OpsRecord {
	return OpsRecord{ID: id, Subject: "外业任务 " + id, Owner: "lin", Priority: OpsPriorityNormal, Status: OpsStatusQueued, Labels: map[string]string{"site": "北坡"}}
}

func TestWrapOpsPreservesOpsErrorType(t *testing.T) {
	err := wrapOps("create", "store.put", ErrOpsConflict)
	var typed *OpsError
	if !errors.As(err, &typed) {
		t.Fatalf("wrapOps 应保留 *OpsError 类型，实际 %T", err)
	}
	if typed.Code != "create" {
		t.Fatalf("OpsError.Code 错误: %s", typed.Code)
	}
	if !errors.Is(err, ErrOpsConflict) {
		t.Fatal("错误链应保留 ErrOpsConflict 哨兵")
	}
}

func TestErrorCodeConflictClassified(t *testing.T) {
	err := fmt.Errorf("%w: revision mismatch", ErrOpsConflict)
	if code := opsCode(err); code != "conflict" {
		t.Fatalf("冲突错误应归类为 conflict，实际 %s", code)
	}
}

func TestErrorCodeNotFoundClassified(t *testing.T) {
	if code := opsCode(ErrOpsNotFound); code != "not_found" {
		t.Fatalf("不存在错误应归类为 not_found，实际 %s", code)
	}
}

func TestNotFoundSentinelClassified(t *testing.T) {
	if !opsIsNotFound(ErrOpsNotFound) {
		t.Fatal("opsIsNotFound 应识别 ErrOpsNotFound")
	}
}

func TestTransitionSentinelClassified(t *testing.T) {
	if !opsIsTransition(ErrOpsTransition) {
		t.Fatal("opsIsTransition 应识别 ErrOpsTransition")
	}
}

func TestCreateConflictPreservesChain(t *testing.T) {
	svc := newOpsService(nil)
	rec := seedErrorRecord("r-001")
	if _, err := svc.Create(context.Background(), rec); err != nil {
		t.Fatal(err)
	}
	_, err := svc.Create(context.Background(), rec)
	if err == nil {
		t.Fatal("重复创建应报错")
	}
	var typed *OpsError
	if !errors.As(err, &typed) {
		t.Fatalf("服务层错误应保留 *OpsError 类型，实际 %T", err)
	}
	if !errors.Is(err, ErrOpsConflict) {
		t.Fatal("服务层错误链应保留 ErrOpsConflict 哨兵")
	}
}

func TestHTTPInvalidStatusReturnsBadRequest(t *testing.T) {
	ts := httptest.NewServer(NewRouter(NewSurveyService(NewMissionStore())))
	defer ts.Close()
	resp, err := http.Post(ts.URL+"/api/missions/status", "application/json", bytes.NewBufferString(`{"id":"mission-241","status":"unknown"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("非法状态应返回 400，实际 %d", resp.StatusCode)
	}
}

func TestHTTPUnknownMissionReturnsNotFound(t *testing.T) {
	ts := httptest.NewServer(NewRouter(NewSurveyService(NewMissionStore())))
	defer ts.Close()
	resp, err := http.Post(ts.URL+"/api/missions/status", "application/json", bytes.NewBufferString(`{"id":"mission-x","status":"landed"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("未知任务应返回 404，实际 %d", resp.StatusCode)
	}
}
