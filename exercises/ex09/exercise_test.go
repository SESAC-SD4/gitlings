package ex09

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
		t.Error("setup 직후에는 실패해야 함 (브랜치 없음)")
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
	if _, err := g.Run("switch", "-c", "new-word"); err != nil {
		t.Fatalf("git switch -c 실패: %v", err)
	}
	words := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n"
	if err := os.WriteFile(filepath.Join(dir, "words.md"), []byte(words), 0644); err != nil {
		t.Fatalf("파일 수정 실패: %v", err)
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		t.Fatalf("git add 실패: %v", err)
	}
	if _, err := g.Run("commit", "-m", "단어 '표범' 추가"); err != nil {
		t.Fatalf("git commit 실패: %v", err)
	}
	if _, err := g.Run("switch", "main"); err != nil {
		t.Fatalf("git switch main 실패: %v", err)
	}
	if _, err := g.Run("merge", "new-word"); err != nil {
		t.Fatalf("git merge 실패: %v", err)
	}

	result, err := check(dir, g)
	if err != nil {
		t.Fatalf("check 오류: %v", err)
	}
	if !result.Passed {
		t.Errorf("모범 풀이 후 통과해야 함. 메시지: %s", result.Message)
	}
}

func TestCommonMistake_InvalidWord(t *testing.T) {
	dir := t.TempDir()
	if err := setup(dir, nil); err != nil {
		t.Fatalf("setup 실패: %v", err)
	}
	g := git.New(dir)

	// 흔한 오답: 규칙에 맞지 않는 단어 추가
	if _, err := g.Run("switch", "-c", "new-word"); err != nil {
		t.Fatalf("git switch -c 실패: %v", err)
	}
	words := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 사과\n" // 규칙 위반
	if err := os.WriteFile(filepath.Join(dir, "words.md"), []byte(words), 0644); err != nil {
		t.Fatalf("파일 수정 실패: %v", err)
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		t.Fatalf("git add 실패: %v", err)
	}
	if _, err := g.Run("commit", "-m", "규칙 위반 단어 추가"); err != nil {
		t.Fatalf("git commit 실패: %v", err)
	}
	if _, err := g.Run("switch", "main"); err != nil {
		t.Fatalf("git switch main 실패: %v", err)
	}
	if _, err := g.Run("merge", "new-word"); err != nil {
		t.Fatalf("git merge 실패: %v", err)
	}

	result, err := check(dir, g)
	if err != nil {
		t.Fatalf("check 오류: %v", err)
	}
	if result.Passed {
		t.Error("규칙 위반 단어 시 실패해야 함")
	}
	t.Logf("오답 메시지: %s", result.Message)
}
