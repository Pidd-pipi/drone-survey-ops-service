package main

import "testing"

func TestListReflectsStatusUpdate(t *testing.T) {
	store := NewMissionStore()
	if _, err := store.UpdateStatus("mission-241", "landed"); err != nil {
		t.Fatal(err)
	}
	after := store.List()
	for _, m := range after {
		if m.ID == "mission-241" && m.Status != "landed" {
			t.Fatalf("列表未反映最新状态: %s", m.Status)
		}
	}
}

func TestMissionListSnapshotIndependent(t *testing.T) {
	store := NewMissionStore()
	first := store.List()
	first[0].Site = "被篡改"
	second := store.List()
	for _, m := range second {
		if m.Site == "被篡改" {
			t.Fatal("返回的列表与内部缓存共享底层数组，外部修改污染了存储")
		}
	}
}

func TestPreviousSnapshotNotOverwrittenByRefresh(t *testing.T) {
	store := NewMissionStore()
	first := store.List()
	if _, err := store.UpdateStatus("mission-241", "landed"); err != nil {
		t.Fatal(err)
	}
	_ = store.List()
	for _, m := range first {
		if m.ID == "mission-241" && m.Status == "landed" {
			t.Fatal("历史快照被后续刷新改写")
		}
	}
}

func TestListOrderStable(t *testing.T) {
	store := NewMissionStore()
	for i := 0; i < 20; i++ {
		id := "mission-bulk-" + string(rune('a'+i))
		store.missions[id] = Mission{ID: id, Site: "测试区", Pilot: "测试员", Images: 10, BatteryPct: 90, Status: "planned"}
	}
	items := store.List()
	if len(items) < 20 {
		t.Fatalf("任务数不足: %d", len(items))
	}
	for i := 1; i < len(items); i++ {
		if items[i-1].ID > items[i].ID {
			t.Fatalf("列表顺序不稳定: %s 在 %s 前面", items[i-1].ID, items[i].ID)
		}
	}
}

func TestServiceMissionsReflectsStatusChange(t *testing.T) {
	svc := NewSurveyService(NewMissionStore())
	_ = svc.Missions()
	if _, err := svc.ChangeStatus("mission-242", "landed"); err != nil {
		t.Fatal(err)
	}
	after := svc.Missions()
	for _, m := range after {
		if m.ID == "mission-242" && m.Status != "landed" {
			t.Fatalf("服务层列表未反映最新状态: %s", m.Status)
		}
	}
}
