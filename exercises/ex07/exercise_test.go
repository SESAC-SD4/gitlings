package ex07

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
		t.Error("setup 직후에는 실패해야 함 (규칙 위반 커밋 있음)")
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
	if _, err := g.Run("reset", "--soft", "HEAD~1"); err != nil {
		t.Fatalf("git reset --soft 실패: %v", err)
	}

	result, err := check(dir, g)
	if err != nil {
		t.Fatalf("check 오류: %v", err)
	}
	if !result.Passed {
		t.Errorf("모범 풀이 후 통과해야 함. 메시지: %s", result.Message)
	}
}

func TestCommonMistake_HardReset(t *testing.T) {
	dir := t.TempDir()
	if err := setup(dir, nil); err != nil {
		t.Fatalf("setup 실패: %v", err)
	}
	g := git.New(dir)

	// 흔한 오답: --hard 사용
	if _, err := g.Run("reset", "--hard", "HEAD~1"); err != nil {
		t.Fatalf("git reset --hard 실패: %v", err)
	}

	result, err := check(dir, g)
	if err != nil {
		t.Fatalf("check 오류: %v", err)
	}
	if result.Passed {
		t.Error("--hard 사용 시 실패해야 함")
	}
	t.Logf("오답 메시지: %s", result.Message)
}

func TestCommonMistake_MixedReset(t *testing.T) {
	dir := t.TempDir()
	if err := setup(dir, nil); err != nil {
		t.Fatalf("setup 실패: %v", err)
	}
	g := git.New(dir)

	// 흔한 오답: --mixed (기본값) 사용
	if _, err := g.Run("reset", "HEAD~1"); err != nil {
		t.Fatalf("git reset 실패: %v", err)
	}

	result, err := check(dir, g)
	if err != nil {
		t.Fatalf("check 오류: %v", err)
	}
	if result.Passed {
		t.Error("--mixed 사용 시 실패해야 함")
	}
	t.Logf("오답 메시지: %s", result.Message)
}
