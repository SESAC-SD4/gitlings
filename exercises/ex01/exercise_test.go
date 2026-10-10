package ex01

import (
	"gitlings-ko/internal/git"
	"os"
	"path/filepath"
	"testing"
)

func TestSetupThenFail(t *testing.T) {
	dir := t.TempDir()
	g := git.New(dir)

	if err := setup(dir, g); err != nil {
		t.Fatalf("setup 실패: %v", err)
	}

	// setup 직후: 저장소가 아니고 words.md만 있음
	if _, err := os.Stat(filepath.Join(dir, "words.md")); err != nil {
		t.Fatal("words.md가 없어야 하는 곳에 없음")
	}

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
	g := git.New(dir)

	if err := setup(dir, g); err != nil {
		t.Fatalf("setup 실패: %v", err)
	}

	// 모범 풀이
	if _, err := g.Run("init"); err != nil {
		t.Fatalf("git init 실패: %v", err)
	}
	if _, err := g.Run("symbolic-ref", "HEAD", "refs/heads/main"); err != nil {
		t.Fatalf("branch 설정 실패: %v", err)
	}
	if err := g.SetupUserConfig(); err != nil {
		t.Fatalf("user config 실패: %v", err)
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		t.Fatalf("git add 실패: %v", err)
	}
	if _, err := g.Run("commit", "-m", "첫 번째 커밋"); err != nil {
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

func TestCommonMistake_InitOnly(t *testing.T) {
	dir := t.TempDir()
	g := git.New(dir)

	if err := setup(dir, g); err != nil {
		t.Fatalf("setup 실패: %v", err)
	}

	// 흔한 오답: init만 하고 커밋 안 함
	if _, err := g.Run("init"); err != nil {
		t.Fatalf("git init 실패: %v", err)
	}
	if _, err := g.Run("symbolic-ref", "HEAD", "refs/heads/main"); err != nil {
		t.Fatalf("branch 설정 실패: %v", err)
	}

	result, err := check(dir, g)
	if err != nil {
		t.Fatalf("check 오류: %v", err)
	}
	if result.Passed {
		t.Error("커밋 없으면 실패해야 함")
	}
	if result.Message == "" {
		t.Error("오답 메시지가 있어야 함")
	}
	t.Logf("오답 메시지: %s", result.Message)
}

func TestCommonMistake_AddButNoCommit(t *testing.T) {
	dir := t.TempDir()
	g := git.New(dir)

	if err := setup(dir, g); err != nil {
		t.Fatalf("setup 실패: %v", err)
	}

	// 흔한 오답: add 했지만 commit 안 함
	if _, err := g.Run("init"); err != nil {
		t.Fatalf("git init 실패: %v", err)
	}
	if _, err := g.Run("symbolic-ref", "HEAD", "refs/heads/main"); err != nil {
		t.Fatalf("branch 설정 실패: %v", err)
	}
	if err := g.SetupUserConfig(); err != nil {
		t.Fatalf("user config 실패: %v", err)
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		t.Fatalf("git add 실패: %v", err)
	}

	result, err := check(dir, g)
	if err != nil {
		t.Fatalf("check 오류: %v", err)
	}
	if result.Passed {
		t.Error("커밋 없으면 실패해야 함")
	}
	t.Logf("staged만 있을 때 메시지: %s", result.Message)
}
