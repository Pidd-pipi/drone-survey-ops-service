package main

import (
	"context"
	"sync"
	"testing"
)

func seedRecord(id string) OpsRecord {
	return OpsRecord{ID: id, Subject: "外业任务 " + id, Owner: "lin", Priority: OpsPriorityNormal, Status: OpsStatusActive, Labels: map[string]string{"site": "北坡", "evidence": "pic-" + id}}
}

func TestOpsStoreConcurrentReadWriteNoRace(t *testing.T) {
	store := newOpsStore([]OpsRecord{seedRecord("r-000")})
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			<-start
			for j := 0; j < 400; j++ {
				_, _ = store.List(context.Background())
			}
		}(i)
	}
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			<-start
			for j := 0; j < 400; j++ {
				_ = store.Count()
			}
		}(i)
	}
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			<-start
			for j := 0; j < 400; j++ {
				_ = store.Put(context.Background(), seedRecord("w-"+string(rune('a'+n))+"-"+string(rune('0'+j%10))))
			}
		}(i)
	}
	close(start)
	wg.Wait()
}

func TestOpsStoreConcurrentGetUpdateNoRace(t *testing.T) {
	store := newOpsStore([]OpsRecord{seedRecord("r-000")})
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 400; j++ {
				_, _ = store.Get(context.Background(), "r-000")
			}
		}()
	}
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 400; j++ {
				_ = store.Update(context.Background(), seedRecord("r-000"), 0)
			}
		}()
	}
	close(start)
	wg.Wait()
}

func TestOpsStoreGetReturnsIndependentLabels(t *testing.T) {
	store := newOpsStore(nil)
	rec := seedRecord("r-001")
	if err := store.Put(context.Background(), rec); err != nil {
		t.Fatal(err)
	}
	got, err := store.Get(context.Background(), "r-001")
	if err != nil {
		t.Fatal(err)
	}
	got.Labels["site"] = "被篡改"
	again, err := store.Get(context.Background(), "r-001")
	if err != nil {
		t.Fatal(err)
	}
	if again.Labels["site"] != "北坡" {
		t.Fatalf("读到的记录被外部改动污染: site=%q", again.Labels["site"])
	}
}

func TestOpsStorePutDoesNotAliasCallerLabels(t *testing.T) {
	store := newOpsStore(nil)
	rec := seedRecord("r-002")
	if err := store.Put(context.Background(), rec); err != nil {
		t.Fatal(err)
	}
	rec.Labels["evidence"] = "被篡改"
	got, err := store.Get(context.Background(), "r-002")
	if err != nil {
		t.Fatal(err)
	}
	if got.Labels["evidence"] != "pic-r-002" {
		t.Fatalf("存储内容受调用方后续修改影响: evidence=%q", got.Labels["evidence"])
	}
}

func TestOpsStoreListReturnsIndependentLabels(t *testing.T) {
	store := newOpsStore(nil)
	if err := store.Put(context.Background(), seedRecord("r-003")); err != nil {
		t.Fatal(err)
	}
	items, err := store.List(context.Background())
	if err != nil || len(items) != 1 {
		t.Fatalf("List 失败: items=%d err=%v", len(items), err)
	}
	items[0].Labels["site"] = "被篡改"
	again, err := store.List(context.Background())
	if err != nil || len(again) != 1 {
		t.Fatalf("二次 List 失败: items=%d err=%v", len(again), err)
	}
	if again[0].Labels["site"] != "北坡" {
		t.Fatalf("List 返回内容受外部改动污染: site=%q", again[0].Labels["site"])
	}
}
