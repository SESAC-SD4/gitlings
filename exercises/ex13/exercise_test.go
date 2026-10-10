package ex13

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
		t.Error("setup 직후에는 실패해야 함 (오타 수정 안 됨)")
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
	// 1. stash
	if _, err := g.Run("stash"); err != nil {
		t.Fatalf("git stash 실패: %v", err)
	}
	// 2. main으로 이동, 오타 수정
	if _, err := g.Run("switch", "main"); err != nil {
		t.Fatalf("git switch main 실패: %v", err)
	}
	rulesFixed := "# 끝말잇기 규칙\n\n다음 단어의 첫 글자는 앞 단어의 끝 글자와 같아야 합니다.\n"
	if err := os.WriteFile(filepath.Join(dir, "rules.md"), []byte(rulesFixed), 0644); err != nil {
		t.Fatalf("파일 수정 실패: %v", err)
	}
	if _, err := g.Run("add", "rules.md"); err != nil {
		t.Fatalf("git add 실패: %v", err)
	}
	if _, err := g.Run("commit", "-m", "rules.md 오타 수정"); err != nil {
		t.Fatalf("git commit 실패: %v", err)
	}
	// 3. new-word로 돌아와 stash pop
	if _, err := g.Run("switch", "new-word"); err != nil {
		t.Fatalf("git switch new-word 실패: %v", err)
	}
	if _, err := g.Run("stash", "pop"); err != nil {
		t.Fatalf("git stash pop 실패: %v", err)
	}

	result, err := check(dir, g)
	if err != nil {
		t.Fatalf("check 오류: %v", err)
	}
	if !result.Passed {
		t.Errorf("모범 풀이 후 통과해야 함. 메시지: %s", result.Message)
	}
}
