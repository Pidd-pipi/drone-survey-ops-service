package main

import "testing"

func TestPausedToActiveAllowed(t *testing.T) {
	sm := newOpsStateMachine()
	if err := sm.Move(OpsStatusPaused, OpsStatusActive, "resume"); err != nil {
		t.Fatalf("暂停恢复执行应被允许: %v", err)
	}
}

func TestClosedToActiveRejected(t *testing.T) {
	sm := newOpsStateMachine()
	if err := sm.Move(OpsStatusClosed, OpsStatusActive, "reopen"); err == nil {
		t.Fatal("已关闭记录不应允许转回执行")
	}
}

func TestPausedValidity(t *testing.T) {
	if !opsStatusValid(OpsStatusPaused) {
		t.Fatal("paused 应是合法状态")
	}
}

func TestStateCanMoveSameStatus(t *testing.T) {
	sm := newOpsStateMachine()
	if !sm.CanMove(OpsStatusQueued, OpsStatusQueued) {
		t.Fatal("同状态应视为可保持")
	}
}

func TestMatchKeepsPausedRecords(t *testing.T) {
	rec := OpsRecord{ID: "r-001", Subject: "外业任务", Status: OpsStatusPaused, Priority: OpsPriorityNormal, Owner: "lin"}
	if !opsMatch(rec, OpsQuery{}) {
		t.Fatal("无过滤条件的查询不应丢掉 paused 记录")
	}
	if !opsMatch(rec, OpsQuery{Status: OpsStatusPaused}) {
		t.Fatal("按 paused 过滤应命中 paused 记录")
	}
}

func TestPausedNotTerminal(t *testing.T) {
	if opsStatusTerminal(OpsStatusPaused) {
		t.Fatal("paused 不应是终态")
	}
	if !opsStatusTerminal(OpsStatusClosed) {
		t.Fatal("closed 应是终态")
	}
}
