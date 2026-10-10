package ex09

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
		ID:       9,
		Slug:     "09-branch-merge",
		Name:     "09 갈래 만들고 합치기",
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
	_, err = g.Run("commit", "-m", "단어장 초기화")
	return err
}

func check(workDir string, g *git.Git) (exercise.Result, error) {
	// new-word 브랜치 존재 확인
	if !g.BranchExists("new-word") {
		return exercise.Result{
			Passed:  false,
			Message: "❌ 'new-word' 브랜치가 없습니다.\n   → 'git switch -c new-word' 로 브랜치를 만들어 보세요.",
		}, nil
	}

	// main에서 확인
	currentBranch, _ := g.CurrentBranch()
	if currentBranch != "main" {
		return exercise.Result{
			Passed:  false,
			Message: "❌ 현재 main 브랜치가 아닙니다.\n   → 'git switch main' 으로 이동 후 merge하세요.",
		}, nil
	}

	// new-word 브랜치의 커밋이 main에 포함됐는지
	newWordHead, err := g.Run("rev-parse", "new-word")
	if err != nil {
		return exercise.Result{Passed: false, Message: "❌ new-word 브랜치를 읽을 수 없습니다."}, nil
	}
	if !g.CommitContains("main", newWordHead) {
		return exercise.Result{
			Passed:  false,
			Message: "❌ new-word 브랜치의 커밋이 main에 병합되지 않았습니다.\n   → 'git merge new-word' 를 실행해 보세요.",
		}, nil
	}

	// 단어 수가 3 이상이어야 함
	content, err := g.FileContent("words.md")
	if err != nil {
		return exercise.Result{Passed: false, Message: "❌ words.md를 읽을 수 없습니다."}, nil
	}
	words := exercise.ParseWordsMd(content)
	if len(words) < 3 {
		return exercise.Result{
			Passed:  false,
			Message: "❌ 단어가 2개뿐입니다. new-word 브랜치에서 단어를 추가해야 합니다.",
		}, nil
	}

	// 끝말잇기 규칙 확인
	if err := exercise.ValidateWordChain(words); err != nil {
		return exercise.Result{
			Passed:  false,
			Message: "❌ 단어장 끝말잇기 규칙 위반: " + err.Error() + "\n   → 규칙에 맞는 단어를 추가해 보세요.\n   (차표의 끝 글자 '표'로 시작하는 단어)",
		}, nil
	}

	return exercise.Result{
		Passed:  true,
		Message: "✅ new-word 브랜치에서 단어를 추가하고 main에 병합했습니다!",
	}, nil
}
