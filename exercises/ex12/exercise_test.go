package ex12

import (
	"gitlings-ko/internal/git"
	"os"
	"path/filepath"
	"strings"
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
		t.Error("setup 직후에는 실패해야 함 (커밋 날린 상태)")
	}
	t.Logf("실패 메시지: %s", result.Message)
}

func TestSolveThenPass(t *testing.T) {
	dir := t.TempDir()
	if err := setup(dir, nil); err != nil {
		t.Fatalf("setup 실패: %v", err)
	}
	g := git.New(dir)

	// 모범 풀이: reflog에서 해시 찾아 reset
	targetHash, err := os.ReadFile(filepath.Join(dir, ".git", targetCommitFile))
	if err != nil {
		t.Fatalf("target hash 읽기 실패: %v", err)
	}
	hash := strings.TrimSpace(string(targetHash))

	if _, err := g.Run("reset", "--hard", hash); err != nil {
		t.Fatalf("git reset --hard 실패: %v", err)
	}

	result, err := check(dir, g)
	if err != nil {
		t.Fatalf("check 오류: %v", err)
	}
	if !result.Passed {
		t.Errorf("모범 풀이 후 통과해야 함. 메시지: %s", result.Message)
	}
}
