package ex01

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
		ID:       1,
		Slug:     "01-first-commit",
		Name:     "01 첫 사진",
		Optional: false,
		Desc:     desc,
		Hints:    exercise.ParseHints(hintsRaw),
		Setup:    setup,
		Check:    check,
	}
}

func setup(workDir string, _ *git.Git) error {
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return err
	}
	words := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n"
	return os.WriteFile(filepath.Join(workDir, "words.md"), []byte(words), 0644)
}

func check(workDir string, g *git.Git) (exercise.Result, error) {
	if !g.HasGitDir() {
		return exercise.Result{
			Passed:  false,
			Message: "❌ 아직 Git 저장소가 아닙니다.\n   → 'git init' 을 먼저 실행해 보세요.",
		}, nil
	}

	count, err := g.CommitCount()
	if err != nil || count == 0 {
		// staged만 있는지 확인
		staged, _ := g.StagedFiles()
		if len(staged) > 0 {
			return exercise.Result{
				Passed:  false,
				Message: "❌ 파일이 무대에 올라가 있지만 아직 커밋하지 않았습니다.\n   → 'git commit -m \"메시지\"' 를 실행해 보세요.",
			}, nil
		}
		return exercise.Result{
			Passed:  false,
			Message: "❌ 사진을 아직 안 찍었습니다 (커밋이 없습니다).\n   → 'git add words.md' 다음 'git commit -m \"메시지\"' 를 실행해 보세요.",
		}, nil
	}

	if !g.FileAtHEAD("words.md") {
		return exercise.Result{
			Passed:  false,
			Message: "❌ words.md가 커밋에 포함되지 않았습니다.\n   → 'git add words.md' 후 커밋해 보세요.",
		}, nil
	}

	if !g.IsClean() {
		staged, _ := g.StagedFiles()
		if len(staged) > 0 {
			return exercise.Result{
				Passed:  false,
				Message: "❌ 파일이 무대에 올라가 있지만 커밋하지 않았습니다.\n   → 'git commit -m \"메시지\"' 를 실행해 보세요.",
			}, nil
		}
		return exercise.Result{
			Passed:  false,
			Message: "❌ 작업 트리에 변경이 남아 있습니다.\n   → 'git status' 로 상태를 확인해 보세요.",
		}, nil
	}

	// user.name 설정 여부 확인 (경고만)
	out, _ := g.Run("log", "-1", "--format=%an")
	warning := ""
	if strings.TrimSpace(out) == "" {
		warning = "\n\n⚠️  커밋 작성자 이름이 없습니다. 'git config user.name \"이름\"' 으로 설정해 보세요."
	}

	return exercise.Result{
		Passed:  true,
		Message: "✅ 첫 번째 커밋 완료! 단어장의 첫 사진을 찍었습니다." + warning,
	}, nil
}
