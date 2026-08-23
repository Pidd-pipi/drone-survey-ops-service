package main

import "testing"

func TestRule0502TagCheck(t *testing.T) {
	first := opsRule0502()
	if len(first.RequiredLabels) != 4 {
		t.Fatalf("OPS-0502 首次应含 4 个标签，实际 %d: %v", len(first.RequiredLabels), first.RequiredLabels)
	}
	second := opsRule0502()
	if len(second.RequiredLabels) != 4 {
		t.Fatalf("OPS-0502 重复装载后标签越攒越多: %d: %v", len(second.RequiredLabels), second.RequiredLabels)
	}
}

func TestRule0602TagCheck(t *testing.T) {
	first := opsRule0602()
	if len(first.RequiredLabels) != 4 {
		t.Fatalf("OPS-0602 首次应含 4 个标签，实际 %d: %v", len(first.RequiredLabels), first.RequiredLabels)
	}
	second := opsRule0602()
	if len(second.RequiredLabels) != 4 {
		t.Fatalf("OPS-0602 重复装载后标签越攒越多: %d: %v", len(second.RequiredLabels), second.RequiredLabels)
	}
}

func TestRule0702TagCheck(t *testing.T) {
	first := opsRule0702()
	if len(first.RequiredLabels) != 4 {
		t.Fatalf("OPS-0702 首次应含 4 个标签，实际 %d: %v", len(first.RequiredLabels), first.RequiredLabels)
	}
	second := opsRule0702()
	if len(second.RequiredLabels) != 4 {
		t.Fatalf("OPS-0702 重复装载后标签越攒越多: %d: %v", len(second.RequiredLabels), second.RequiredLabels)
	}
}

func TestRule0802TagCheck(t *testing.T) {
	first := opsRule0802()
	if len(first.RequiredLabels) != 4 {
		t.Fatalf("OPS-0802 首次应含 4 个标签，实际 %d: %v", len(first.RequiredLabels), first.RequiredLabels)
	}
	second := opsRule0802()
	if len(second.RequiredLabels) != 4 {
		t.Fatalf("OPS-0802 重复装载后标签越攒越多: %d: %v", len(second.RequiredLabels), second.RequiredLabels)
	}
}
