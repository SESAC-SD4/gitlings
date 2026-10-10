package ex04

import (
	_ "embed"
	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/git"
	"os"
	"path/filepath"
)

//go:embed desc.txt
var desc string

//go:embed hints.txt
var hintsRaw string

func New() *exercise.Exercise {
	return &exercise.Exercise{
		ID:       4,
		Slug:     "04-discard",
		Name:     "04 망친 수정 버리기",
		Optional: false,
		Desc:     desc,
		Hints:    exercise.ParseHints(hintsRaw),
		Setup:    setup,
		Check:    check,
	}
}

const cleanWords = "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n"

func setup(workDir string, _ *git.Git) error {
	g, err := git.InitMain(workDir)
	if err != nil {
		return err
	}
	if err := g.SetupUserConfig(); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(workDir, "words.md"), []byte(cleanWords), 0644); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}
	if _, err := g.Run("commit", "-m", "단어장 초기화"); err != nil {
		return err
	}

	// 실수로 엉망이 된 words.md (add 안 함)
	messy := "이게 뭔가요\n????\nasdfasdf\n제발 살려주세요\n"
	return os.WriteFile(filepath.Join(workDir, "words.md"), []byte(messy), 0644)
}

func check(workDir string, g *git.Git) (exercise.Result, error) {
	// 커밋 수가 늘지 않아야 함
	count, _ := g.CommitCount()
	if count > 1 {
		return exercise.Result{
			Passed:  false,
			Message: "❌ 커밋이 새로 생겼습니다. 커밋을 추가하는 게 아니라 작업 트리를 되돌려야 합니다.\n   → git restore words.md 를 사용해 보세요.",
		}, nil
	}

	// 작업 트리가 깨끗해야 함
	if !g.IsClean() {
		return exercise.Result{
			Passed:  false,
			Message: "❌ words.md가 아직 변경된 상태입니다.\n   → 'git restore words.md' 를 실행해 보세요.",
		}, nil
	}

	return exercise.Result{
		Passed:  true,
		Message: "✅ 실수한 수정을 버리고 마지막 커밋 상태로 되돌렸습니다!",
	}, nil
}
