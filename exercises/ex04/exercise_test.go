package ex04

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
		t.Error("setup 직후에는 실패해야 함 (words.md 엉망)")
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
	if _, err := g.Run("restore", "words.md"); err != nil {
		t.Fatalf("git restore 실패: %v", err)
	}

	result, err := check(dir, g)
	if err != nil {
		t.Fatalf("check 오류: %v", err)
	}
	if !result.Passed {
		t.Errorf("모범 풀이 후 통과해야 함. 메시지: %s", result.Message)
	}
}

func TestCommonMistake_NewCommit(t *testing.T) {
	dir := t.TempDir()
	if err := setup(dir, nil); err != nil {
		t.Fatalf("setup 실패: %v", err)
	}
	g := git.New(dir)

	// 흔한 오답: 파일을 수동으로 고쳐서 새 커밋
	// cleanWords와 다른 내용으로 덮어씀 (단어 추가)
	modified := cleanWords + "4. 래퍼\n"
	if err := os.WriteFile(filepath.Join(dir, "words.md"), []byte(modified), 0644); err != nil {
		t.Fatalf("파일 수정 실패: %v", err)
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		t.Fatalf("git add 실패: %v", err)
	}
	if _, err := g.Run("commit", "-m", "새 커밋으로 수정"); err != nil {
		t.Skipf("커밋 불가 (환경 문제): %v", err)
	}

	result, err := check(dir, g)
	if err != nil {
		t.Fatalf("check 오류: %v", err)
	}
	if result.Passed {
		t.Error("새 커밋으로 해결하면 실패해야 함")
	}
	t.Logf("오답 메시지: %s", result.Message)
}
