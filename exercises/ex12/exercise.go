package ex12

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
		ID:       12,
		Slug:     "12-reflog",
		Name:     "12★ 날린 커밋 되살리기",
		Optional: true,
		Desc:     desc,
		Hints:    exercise.ParseHints(hintsRaw),
		Setup:    setup,
		Check:    check,
	}
}

// setup 시 저장해두는 커밋 해시 (check에서 비교용)
const targetCommitFile = ".gitlings_target_commit"

func setup(workDir string, _ *git.Git) error {
	g, err := git.InitMain(workDir)
	if err != nil {
		return err
	}
	if err := g.SetupUserConfig(); err != nil {
		return err
	}

	// 커밋 1
	words1 := "# 끝말잇기 단어장\n\n1. 기차\n"
	if err := os.WriteFile(filepath.Join(workDir, "words.md"), []byte(words1), 0644); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}
	if _, err := g.Run("commit", "-m", "단어 1개"); err != nil {
		return err
	}

	// 커밋 2
	words2 := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n"
	if err := os.WriteFile(filepath.Join(workDir, "words.md"), []byte(words2), 0644); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}
	if _, err := g.Run("commit", "-m", "단어 2개"); err != nil {
		return err
	}

	// 커밋 3 (목표)
	words3 := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n"
	if err := os.WriteFile(filepath.Join(workDir, "words.md"), []byte(words3), 0644); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}
	if _, err := g.Run("commit", "-m", "단어 3개"); err != nil {
		return err
	}

	// 목표 해시 저장 (check에서 비교용, .git 안에 숨겨둠)
	targetHash, err := g.HeadCommit()
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(workDir, ".git", targetCommitFile), []byte(targetHash), 0644); err != nil {
		return err
	}

	// 2개 날리기
	_, err = g.Run("reset", "--hard", "HEAD~2")
	return err
}

func check(workDir string, g *git.Git) (exercise.Result, error) {
	// 목표 해시 읽기
	targetHash, err := os.ReadFile(filepath.Join(workDir, ".git", targetCommitFile))
	if err != nil {
		return exercise.Result{Passed: false, Message: "❌ setup 정보를 읽을 수 없습니다. gitlings reset으로 초기화해 보세요."}, nil
	}
	target := string(targetHash)

	// HEAD가 목표 해시를 가리키는지 확인
	head, err := g.HeadCommit()
	if err != nil {
		return exercise.Result{Passed: false, Message: "❌ HEAD를 읽을 수 없습니다."}, nil
	}

	if head != target {
		count, _ := g.CommitCount()
		if count < 3 {
			return exercise.Result{
				Passed:  false,
				Message: "❌ 아직 날린 커밋이 복구되지 않았습니다.\n   → 'git reflog' 로 이전 커밋 해시를 찾아보세요.\n   그 다음: git reset --hard [해시]",
			}, nil
		}
		return exercise.Result{
			Passed:  false,
			Message: "❌ 올바른 커밋으로 복구되지 않았습니다.\n   → 'git reflog' 로 'HEAD~2' 이전의 커밋 해시를 찾아보세요.",
		}, nil
	}

	return exercise.Result{
		Passed:  true,
		Message: "✅ 날린 커밋을 reflog로 되살렸습니다!",
	}, nil
}
