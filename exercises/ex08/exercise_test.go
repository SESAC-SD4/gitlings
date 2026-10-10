package ex08

import (
	"gitlings-ko/internal/git"
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
		t.Error("setup 직후에는 실패해야 함 (push 안 됨)")
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
	if _, err := g.Run("revert", "HEAD", "--no-edit"); err != nil {
		t.Fatalf("git revert 실패: %v", err)
	}
	if _, err := g.Run("push", "origin", "main"); err != nil {
		t.Fatalf("git push 실패: %v", err)
	}

	result, err := check(dir, g)
	if err != nil {
		t.Fatalf("check 오류: %v", err)
	}
	if !result.Passed {
		t.Errorf("모범 풀이 후 통과해야 함. 메시지: %s", result.Message)
	}
}

func TestCommonMistake_ResetAndForcePush(t *testing.T) {
	dir := t.TempDir()
	if err := setup(dir, nil); err != nil {
		t.Fatalf("setup 실패: %v", err)
	}
	g := git.New(dir)

	// 흔한 오답: reset --hard 후 강제 push
	if _, err := g.Run("reset", "--hard", "HEAD~1"); err != nil {
		t.Fatalf("git reset --hard 실패: %v", err)
	}
	if _, err := g.Run("push", "--force", "origin", "main"); err != nil {
		t.Fatalf("git push --force 실패: %v", err)
	}

	result, err := check(dir, g)
	if err != nil {
		t.Fatalf("check 오류: %v", err)
	}
	if result.Passed {
		t.Error("강제 push로 이력 삭제 시 실패해야 함")
	}
	t.Logf("오답 메시지: %s", result.Message)
}
