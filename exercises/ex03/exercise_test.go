package ex03

import (
	"gitlings-ko/internal/git"
	"os"
	"path/filepath"
	"testing"
)

func TestSetupThenFail(t *testing.T) {
	dir := t.TempDir()
	if err := setup(dir, nil); err != nil {
		t.Fatalf("setup 실패: %v", err)
	}
	g := git.New(dir)
	result, err := check(dir, g)
	if err != nil {
		t.Fatalf("check 오류: %v", err)
	}
	if result.Passed {
		t.Error("setup 직후에는 실패해야 함")
	}
	t.Logf("실패 메시지: %s", result.Message)
}

func TestSolveThenPass(t *testing.T) {
	dir := t.TempDir()
	if err := setup(dir, nil); err != nil {
		t.Fatalf("setup 실패: %v", err)
	}
	g := git.New(dir)

	// 모범 풀이
	if _, err := g.Run("restore", "--staged", "notes.md"); err != nil {
		t.Fatalf("git restore --staged 실패: %v", err)
	}
	if _, err := g.Run("commit", "-m", "단어 '범고래' 추가"); err != nil {
		t.Fatalf("git commit 실패: %v", err)
	}

	result, err := check(dir, g)
	if err != nil {
		t.Fatalf("check 오류: %v", err)
	}
	if !result.Passed {
		t.Errorf("모범 풀이 후 통과해야 함. 메시지: %s", result.Message)
	}
}

func TestCommonMistake_DeleteNotes(t *testing.T) {
	dir := t.TempDir()
	if err := setup(dir, nil); err != nil {
		t.Fatalf("setup 실패: %v", err)
	}
	g := git.New(dir)

	// 흔한 오답: notes.md 삭제
	os.Remove(filepath.Join(dir, "notes.md"))
	if _, err := g.Run("add", "."); err != nil {
		t.Fatalf("git add 실패: %v", err)
	}
	if _, err := g.Run("commit", "-m", "notes.md 삭제함"); err != nil {
		t.Fatalf("git commit 실패: %v", err)
	}

	result, err := check(dir, g)
	if err != nil {
		t.Fatalf("check 오류: %v", err)
	}
	if result.Passed {
		t.Error("notes.md 삭제 시 실패해야 함")
	}
	t.Logf("오답 메시지: %s", result.Message)
}
