package ex02

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

	// 모범 풀이: words.md만 add하고 commit
	if _, err := g.Run("add", "words.md"); err != nil {
		t.Fatalf("git add 실패: %v", err)
	}
	if _, err := g.Run("commit", "-m", "단어 '표범' 추가"); err != nil {
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

func TestCommonMistake_AddAll(t *testing.T) {
	dir := t.TempDir()
	if err := setup(dir, nil); err != nil {
		t.Fatalf("setup 실패: %v", err)
	}
	g := git.New(dir)

	// 흔한 오답: git add . 로 둘 다 커밋
	if _, err := g.Run("add", "."); err != nil {
		t.Fatalf("git add 실패: %v", err)
	}
	// notes.md를 제외하기 위해 .gitlings_id도 무시
	os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(".gitlings_id\n"), 0644)
	g.Run("add", ".")
	if _, err := g.Run("commit", "-m", "잘못된 커밋"); err != nil {
		t.Fatalf("git commit 실패: %v", err)
	}

	result, err := check(dir, g)
	if err != nil {
		t.Fatalf("check 오류: %v", err)
	}
	if result.Passed {
		t.Error("notes.md까지 커밋하면 실패해야 함")
	}
	t.Logf("오답 메시지: %s", result.Message)
}
