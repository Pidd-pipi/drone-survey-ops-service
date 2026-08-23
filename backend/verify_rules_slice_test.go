package main

import "testing"

func TestRulePool03LengthStable(t *testing.T) {
	first := opsRules03()
	if len(first) != 8 {
		t.Fatalf("分组 03 首次应有 8 条规则，实际 %d", len(first))
	}
	_ = opsRules04()
	second := opsRules03()
	if len(second) != 8 {
		t.Fatalf("分组 03 重复装载后长度增长: %d", len(second))
	}
}

func TestRulePool04LengthStable(t *testing.T) {
	first := opsRules04()
	if len(first) != 8 {
		t.Fatalf("分组 04 首次应有 8 条规则，实际 %d", len(first))
	}
	_ = opsRules03()
	second := opsRules04()
	if len(second) != 8 {
		t.Fatalf("分组 04 重复装载后长度增长: %d", len(second))
	}
}

func TestRule0302LabelsStable(t *testing.T) {
	first := opsRule0302()
	if len(first.RequiredLabels) != 4 {
		t.Fatalf("OPS-0302 首次应含 4 个标签，实际 %d: %v", len(first.RequiredLabels), first.RequiredLabels)
	}
	second := opsRule0302()
	if len(second.RequiredLabels) != 4 {
		t.Fatalf("OPS-0302 重复装载后标签越攒越多: %d: %v", len(second.RequiredLabels), second.RequiredLabels)
	}
}

func TestRule0406LabelsStable(t *testing.T) {
	first := opsRule0406()
	if len(first.RequiredLabels) != 4 {
		t.Fatalf("OPS-0406 首次应含 4 个标签，实际 %d: %v", len(first.RequiredLabels), first.RequiredLabels)
	}
	second := opsRule0406()
	if len(second.RequiredLabels) != 4 {
		t.Fatalf("OPS-0406 重复装载后标签越攒越多: %d: %v", len(second.RequiredLabels), second.RequiredLabels)
	}
}

func TestOpsRulesAllCodesDistinct(t *testing.T) {
	rules := opsRules()
	seen := map[string]bool{}
	for _, r := range rules {
		if seen[r.Code] {
			t.Fatalf("规则码重复: %s", r.Code)
		}
		seen[r.Code] = true
	}
	if len(seen) != 112 {
		t.Fatalf("规则码总数应为 112，实际 %d", len(seen))
	}
}
