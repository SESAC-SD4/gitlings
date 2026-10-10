package ex02

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
		ID:       2,
		Slug:     "02-pick-files",
		Name:     "02 골라서 찍기",
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

	words := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n"
	if err := os.WriteFile(filepath.Join(workDir, "words.md"), []byte(words), 0644); err != nil {
		return err
	}

	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}
	if _, err := g.Run("commit", "-m", "단어장 초기화"); err != nil {
		return err
	}

	// 3번째 단어 추가 (unstaged)
	words2 := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n"
	if err := os.WriteFile(filepath.Join(workDir, "words.md"), []byte(words2), 0644); err != nil {
		return err
	}
	// notes.md 추가 (unstaged)
	notes := "오늘 공부한 것: git add, git commit\n아직 잘 모르겠음...\n"
	return os.WriteFile(filepath.Join(workDir, "notes.md"), []byte(notes), 0644)
}

func check(workDir string, g *git.Git) (exercise.Result, error) {
	count, err := g.CommitCount()
	if err != nil || count < 2 {
		return exercise.Result{
			Passed:  false,
			Message: "❌ 새 커밋이 없습니다.\n   → words.md만 add하고 commit해 보세요.",
		}, nil
	}

	out, _ := g.Run("diff", "--name-only", "HEAD~1", "HEAD")
	changed := exercise.Lines(out)

	if exercise.Contains(changed, "notes.md") {
		return exercise.Result{
			Passed:  false,
			Message: "❌ 낙서(notes.md)까지 사진에 찍혔습니다.\n   → 커밋을 취소하고 words.md만 add해서 다시 커밋해 보세요.",
		}, nil
	}

	if !exercise.Contains(changed, "words.md") {
		return exercise.Result{
			Passed:  false,
			Message: "❌ 새 커밋에 words.md 변경이 없습니다.\n   → git add words.md 후 commit을 실행해 보세요.",
		}, nil
	}

	// notes.md가 working tree에 있는지 (unstaged 또는 untracked)
	statusOut, _ := g.Run("status", "--short")
	notesPresent := false
	for _, line := range exercise.Lines(statusOut) {
		if strings.Contains(line, "notes.md") {
			notesPresent = true
		}
	}
	if !notesPresent {
		return exercise.Result{
			Passed:  false,
			Message: "❌ notes.md가 수정된 상태로 남아 있지 않습니다.\n   → notes.md는 add 없이 그냥 두세요.",
		}, nil
	}

	content, err := g.FileContent("words.md")
	if err != nil {
		return exercise.Result{Passed: false, Message: "❌ words.md를 읽을 수 없습니다."}, nil
	}
	words := exercise.ParseWordsMd(content)
	if err := exercise.ValidateWordChain(words); err != nil {
		return exercise.Result{
			Passed:  false,
			Message: "❌ 단어장 끝말잇기 규칙 위반: " + err.Error(),
		}, nil
	}

	return exercise.Result{
		Passed:  true,
		Message: "✅ words.md만 커밋하고 notes.md는 수정 상태로 남겼습니다!",
	}, nil
}
