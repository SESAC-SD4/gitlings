package ex06

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
		t.Error("setup 직후에는 실패해야 함 (.env 추적 중)")
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
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(".env\n"), 0644); err != nil {
		t.Fatalf("gitignore 생성 실패: %v", err)
	}
	if _, err := g.Run("rm", "--cached", ".env"); err != nil {
		t.Fatalf("git rm --cached 실패: %v", err)
	}
	if _, err := g.Run("add", ".gitignore"); err != nil {
		t.Fatalf("git add 실패: %v", err)
	}
	if _, err := g.Run("commit", "-m", ".env를 gitignore에 추가"); err != nil {
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

func TestCommonMistake_GitignoreOnly(t *testing.T) {
	dir := t.TempDir()
	if err := setup(dir, nil); err != nil {
		t.Fatalf("setup 실패: %v", err)
	}
	g := git.New(dir)

	// 흔한 오답: .gitignore만 추가하고 rm --cached 안 함
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(".env\n"), 0644); err != nil {
		t.Fatalf("gitignore 생성 실패: %v", err)
	}
	if _, err := g.Run("add", ".gitignore"); err != nil {
		t.Fatalf("git add 실패: %v", err)
	}
	if _, err := g.Run("commit", "-m", "gitignore만 추가"); err != nil {
		t.Fatalf("git commit 실패: %v", err)
	}

	result, err := check(dir, g)
	if err != nil {
		t.Fatalf("check 오류: %v", err)
	}
	if result.Passed {
		t.Error("rm --cached 없이 gitignore만 추가하면 실패해야 함")
	}
	t.Logf("오답 메시지: %s", result.Message)
}

func TestCommonMistake_DeleteEnv(t *testing.T) {
	dir := t.TempDir()
	if err := setup(dir, nil); err != nil {
		t.Fatalf("setup 실패: %v", err)
	}
	g := git.New(dir)

	// 흔한 오답: .env 파일 삭제
	if _, err := g.Run("rm", ".env"); err != nil {
		t.Fatalf("git rm 실패: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(".env\n"), 0644); err != nil {
		t.Fatalf("gitignore 생성 실패: %v", err)
	}
	if _, err := g.Run("add", ".gitignore"); err != nil {
		t.Fatalf("git add 실패: %v", err)
	}
	if _, err := g.Run("commit", "-m", ".env 삭제"); err != nil {
		t.Fatalf("git commit 실패: %v", err)
	}

	result, err := check(dir, g)
	if err != nil {
		t.Fatalf("check 오류: %v", err)
	}
	if result.Passed {
		t.Error(".env 파일 삭제 시 실패해야 함")
	}
	t.Logf("오답 메시지: %s", result.Message)
}
