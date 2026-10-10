package ex03

import (
	_ "embed"
	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/git"
	"os"
	"path/filepath"
	"strings"
)

//go:embed desc.txt
var desc string

//go:embed hints.txt
var hintsRaw string

func New() *exercise.Exercise {
	return &exercise.Exercise{
		ID:       3,
		Slug:     "03-unstage",
		Name:     "03 프레임에서 빼기",
		Optional: false,
		Desc:     desc,
		Hints:    exercise.ParseHints(hintsRaw),
		Setup:    setup,
		Check:    check,
	}
}

func setup(workDir string, _ *git.Git) error {
	g, err := git.InitMain(workDir)
	if err != nil {
		return err
	}
	if err := g.SetupUserConfig(); err != nil {
		return err
	}

	// 초기 상태 만들기
	words := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n"
	if err := os.WriteFile(filepath.Join(workDir, "words.md"), []byte(words), 0644); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}
	if _, err := g.Run("commit", "-m", "단어장 초기화"); err != nil {
		return err
	}

	// 4번째 단어 추가 및 stage
	words2 := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n4. 범고래\n"
	if err := os.WriteFile(filepath.Join(workDir, "words.md"), []byte(words2), 0644); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}

	// notes.md 추가 및 stage
	notes := "오늘은 unstage 연습!\n"
	if err := os.WriteFile(filepath.Join(workDir, "notes.md"), []byte(notes), 0644); err != nil {
		return err
	}
	_, err = g.Run("add", "notes.md")
	return err
}

func check(workDir string, g *git.Git) (exercise.Result, error) {
	// notes.md가 HEAD 커밋에 없어야 함
	if g.FileAtHEAD("notes.md") {
		return exercise.Result{
			Passed:  false,
			Message: "❌ notes.md가 커밋에 포함됐습니다.\n   → 'git restore --staged notes.md' 로 무대에서 내리세요.",
		}, nil
	}

	// notes.md가 디스크에 있어야 함 (삭제하면 안 됨)
	if _, err := os.Stat(filepath.Join(workDir, "notes.md")); os.IsNotExist(err) {
		return exercise.Result{
			Passed:  false,
			Message: "❌ notes.md 파일이 삭제됐습니다.\n   → 파일은 지우지 말고 무대에서만 내리세요.\n   힌트: git restore --staged notes.md",
		}, nil
	}

	// words.md가 HEAD에 있어야 함
	if !g.FileAtHEAD("words.md") {
		return exercise.Result{
			Passed:  false,
			Message: "❌ words.md가 커밋에 없습니다.\n   → words.md를 커밋해야 합니다.",
		}, nil
	}

	// 커밋이 2개 이상인지 (새 커밋이 생겼는지)
	count, _ := g.CommitCount()
	if count < 2 {
		return exercise.Result{
			Passed:  false,
			Message: "❌ 아직 커밋하지 않았습니다.\n   → notes.md를 무대에서 내린 후 words.md를 커밋해 보세요.",
		}, nil
	}

	// 마지막 커밋에 notes.md가 없는지 재확인
	out, _ := g.Run("diff", "--name-only", "HEAD~1", "HEAD")
	changed := exercise.Lines(out)
	if exercise.Contains(changed, "notes.md") {
		return exercise.Result{
			Passed:  false,
			Message: "❌ 마지막 커밋에 notes.md가 포함됐습니다.\n   → 커밋을 취소하고 다시 시도해 보세요.",
		}, nil
	}

	// notes.md가 untracked 상태인지 확인
	statusOut, _ := g.Run("status", "--short")
	notesStatus := ""
	for _, line := range exercise.Lines(statusOut) {
		if strings.Contains(line, "notes.md") {
			notesStatus = line
		}
	}
	if notesStatus == "" {
		return exercise.Result{
			Passed:  false,
			Message: "❌ notes.md가 작업 트리에 없습니다.",
		}, nil
	}

	return exercise.Result{
		Passed:  true,
		Message: "✅ notes.md를 무대에서 내리고 words.md만 커밋했습니다!",
	}, nil
}
