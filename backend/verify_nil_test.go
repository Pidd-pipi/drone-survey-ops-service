package main

import "testing"

func TestValidateMissionStatusValid(t *testing.T) {
	if err := ValidateMissionStatus("landed"); err != nil {
		t.Fatalf("合法状态 landed 校验失败: %v", err)
	}
	if err := ValidateMissionStatus("planned"); err != nil {
		t.Fatalf("合法状态 planned 校验失败: %v", err)
	}
}

func TestStatusKnownCheck(t *testing.T) {
	if !(Mission{Status: "flying"}).StatusKnown() {
		t.Fatal("合法状态 flying 应被识别为已知")
	}
	if (Mission{Status: "unknown"}).StatusKnown() {
		t.Fatal("非法状态 unknown 不应被识别为已知")
	}
}

func TestStatusCheckerNilGuard(t *testing.T) {
	checker := defaultStatusChecker()
	if checker == nil {
		t.Fatal("默认状态校验器不应为 nil")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("状态校验器调用不应 panic: %v", r)
			}
		}()
		if err := ValidateMissionStatus("planned"); err != nil {
			t.Errorf("合法状态校验失败: %v", err)
		}
	}()
}
