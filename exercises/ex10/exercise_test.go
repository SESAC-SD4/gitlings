package ex10

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
		t.Error("setup 직후에는 실패해야 함 (병합 안 됨)")
	}
	t.Logf("실패 메시지: %s", result.Message)
}

func TestSolveThenPass(t *testing.T) {
	dir := t.TempDir()
	if err := setup(dir, nil); err != nil {
		t.Fatalf("setup 실패: %v", err)
	}
	g := git.New(dir)

	// 모범 풀이: merge → 충돌 → 표범 선택
	g.Run("merge", "partner") // 충돌 발생 (에러 무시)

	// 충돌 해결: 표범 선택
	words := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n"
	if err := os.WriteFile(filepath.Join(dir, "words.md"), []byte(words), 0644); err != nil {
		t.Fatalf("파일 수정 실패: %v", err)
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		t.Fatalf("git add 실패: %v", err)
	}
	if _, err := g.Run("commit", "-m", "충돌 해결: 표범 선택"); err != nil {
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

func TestCommonMistake_LeaveConflictMarkers(t *testing.T) {
	dir := t.TempDir()
	if err := setup(dir, nil); err != nil {
		t.Fatalf("setup 실패: %v", err)
	}
	g := git.New(dir)

	g.Run("merge", "partner") // 충돌 발생

	// 흔한 오답: 충돌 표시를 남긴 채 커밋
	wordsContent, _ := os.ReadFile(filepath.Join(dir, "words.md"))
	_ = wordsContent
	conflictContent := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n<<<<<<< HEAD\n3. 표범\n=======\n3. 사과\n>>>>>>> partner\n"
	if err := os.WriteFile(filepath.Join(dir, "words.md"), []byte(conflictContent), 0644); err != nil {
		t.Fatalf("파일 수정 실패: %v", err)
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		t.Fatalf("git add 실패: %v", err)
	}
	if _, err := g.Run("commit", "-m", "충돌 표시 남긴 채 커밋"); err != nil {
		t.Fatalf("git commit 실패: %v", err)
	}

	result, err := check(dir, g)
	if err != nil {
		t.Fatalf("check 오류: %v", err)
	}
	if result.Passed {
		t.Error("충돌 표시 남긴 채 커밋 시 실패해야 함")
	}
	t.Logf("오답 메시지: %s", result.Message)
}
