package main

import (
	"context"
	"testing"
	"time"
)

func seedContextRecord(id string) OpsRecord {
	return OpsRecord{ID: id, Subject: "外业任务 " + id, Owner: "lin", Priority: OpsPriorityNormal, Status: OpsStatusQueued, Labels: map[string]string{"site": "北坡"}}
}

func TestContextKeepsParentDeadline(t *testing.T) {
	parent, cancel := context.WithDeadline(context.Background(), time.Now().Add(200*time.Millisecond))
	defer cancel()
	derived, derivedCancel := opsContext(parent, 3*time.Second)
	defer derivedCancel()
	deadline, ok := derived.Deadline()
	if !ok {
		t.Fatal("派生上下文丢失 deadline")
	}
	parentDeadline, _ := parent.Deadline()
	if deadline.After(parentDeadline) {
		t.Fatalf("派生 deadline 晚于父 deadline: derived=%v parent=%v", deadline, parentDeadline)
	}
	parentCtx, parentCancel := context.WithCancel(context.Background())
	derived2, derivedCancel2 := opsContext(parentCtx, 3*time.Second)
	defer derivedCancel2()
	parentCancel()
	select {
	case <-derived2.Done():
	case <-time.After(300 * time.Millisecond):
		t.Fatal("父 ctx 取消后派生上下文未随之取消")
	}
}

func TestDelayHonorsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := time.Now()
	err := opsDelay(ctx, 300*time.Millisecond)
	if err == nil {
		t.Fatalf("已取消的 ctx 下 opsDelay 应返回错误，实际 nil（耗时 %v）", time.Since(start))
	}
	if time.Since(start) > 200*time.Millisecond {
		t.Fatalf("取消后 opsDelay 未立即返回，耗时 %v", time.Since(start))
	}
}

func TestTransitionHonorsParentContext(t *testing.T) {
	svc := newOpsService([]OpsRecord{seedContextRecord("r-001")})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := svc.Transition(ctx, "r-001", 0, OpsStatusActive, "tester")
	if err == nil {
		t.Fatal("父 ctx 已取消时 Transition 应返回错误")
	}
	rec, _ := svc.Get(context.Background(), "r-001")
	if rec.Status != OpsStatusQueued {
		t.Fatalf("取消后记录状态不应改变: %s", rec.Status)
	}
}

func TestCreateHonorsCancelledContext(t *testing.T) {
	svc := newOpsService(nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := svc.Create(ctx, seedContextRecord("r-002"))
	if err == nil {
		t.Fatal("父 ctx 已取消时 Create 应返回错误")
	}
	if svc.Count() != 0 {
		t.Fatal("取消后不应写入记录")
	}
}

func TestSearchHonorsCancelledContext(t *testing.T) {
	svc := newOpsService([]OpsRecord{seedContextRecord("r-003")})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := svc.Search(ctx, OpsQuery{})
	if err == nil {
		t.Fatal("父 ctx 已取消时 Search 应返回错误")
	}
}
