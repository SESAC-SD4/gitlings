package progress

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSaveRoundTrip(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	p, err := Load()
	if err != nil || len(p.Passed) != 0 || p.LastStarted != 0 {
		t.Fatalf("파일이 없으면 빈 진행 상황: %+v %v", p, err)
	}
	p.MarkPassed(3)
	p.MarkPassed(3)
	p.LastStarted = 4
	if p.NextHint(4, 3) != 1 || p.NextHint(4, 3) != 2 {
		t.Fatal("힌트 단계가 1, 2로 올라가야 합니다")
	}
	if err := p.Save(); err != nil {
		t.Fatal(err)
	}
	q, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if !q.IsPassed(3) || len(q.Passed) != 1 || q.LastStarted != 4 || q.HintLevel[4] != 2 {
		t.Fatalf("저장한 그대로 읽혀야 합니다: %+v", q)
	}
	if _, err := os.Stat(filepath.Join(home, "gitlings-work", ".progress.json")); err != nil {
		t.Fatal("~/gitlings-work/.progress.json에 저장해야 합니다")
	}
}

func TestNextHintCapsAndResets(t *testing.T) {
	p := &Progress{}
	for i := 0; i < 5; i++ {
		p.NextHint(1, 3)
	}
	if p.HintLevel[1] != 3 {
		t.Fatalf("3단계에서 멈춰야 합니다: %d", p.HintLevel[1])
	}
	p.ResetHint(1)
	if p.NextHint(1, 3) != 1 {
		t.Fatal("ResetHint 뒤에는 1단계부터")
	}
}

func TestUnpass(t *testing.T) {
	p := &Progress{Passed: []int{1, 2, 3}}
	p.Unpass(2)
	if p.IsPassed(2) || !p.IsPassed(1) || !p.IsPassed(3) {
		t.Fatalf("2번만 통과 기록에서 빠져야 합니다: %v", p.Passed)
	}
	p.Unpass(9)
	if len(p.Passed) != 2 {
		t.Fatal("없는 번호를 빼도 다른 기록은 그대로여야 합니다")
	}
}

func TestLoadBrokenFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := filepath.Join(home, "gitlings-work")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".progress.json"), []byte("{망가짐"), 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := Load()
	if err != nil || p == nil || len(p.Passed) != 0 {
		t.Fatalf("망가진 파일은 빈 진행 상황으로 시작: %+v %v", p, err)
	}
}
