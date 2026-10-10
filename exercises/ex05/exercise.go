package ex05

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
		ID:       5,
		Slug:     "05-amend",
		Name:     "05 메시지 오타",
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

	words1 := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n"
	if err := os.WriteFile(filepath.Join(workDir, "words.md"), []byte(words1), 0644); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}
	if _, err := g.Run("commit", "-m", "단어장 초기화"); err != nil {
		return err
	}

	// 오타 있는 커밋 메시지로 커밋
	words2 := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n"
	if err := os.WriteFile(filepath.Join(workDir, "words.md"), []byte(words2), 0644); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}
	_, err = g.Run("commit", "-m", "단어 추가함ㅁㅁㄴㅇ")
	return err
}

func check(workDir string, g *git.Git) (exercise.Result, error) {
	// 커밋 수가 늘지 않아야 함
	count, _ := g.CommitCount()
	if count > 2 {
		return exercise.Result{
			Passed:  false,
			Message: "❌ 커밋이 하나 늘었습니다. 새로 찍은 게 아니라 고쳐야 합니다.\n   → git commit --amend -m \"단어 '표범' 추가\" 를 사용해 보세요.",
		}, nil
	}

	// HEAD 메시지 확인
	msg, err := g.HeadMessage()
	if err != nil {
		return exercise.Result{Passed: false, Message: "❌ 커밋 메시지를 읽을 수 없습니다."}, nil
	}

	if strings.TrimSpace(msg) != "단어 '표범' 추가" {
		if strings.Contains(msg, "ㅁㅁㄴㅇ") {
			return exercise.Result{
				Passed:  false,
				Message: "❌ 아직 오타가 남아 있습니다.\n   → git commit --amend -m \"단어 '표범' 추가\" 를 실행해 보세요.",
			}, nil
		}
		return exercise.Result{
			Passed:  false,
			Message: "❌ 커밋 메시지가 다릅니다.\n   현재: " + strings.TrimSpace(msg) + "\n   원하는: 단어 '표범' 추가",
		}, nil
	}

	return exercise.Result{
		Passed:  true,
		Message: "✅ 커밋 메시지를 올바르게 수정했습니다!",
	}, nil
}
