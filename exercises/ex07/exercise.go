package ex07

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
		ID:       7,
		Slug:     "07-reset-soft",
		Name:     "07 방금 커밋 취소",
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

	// 초기 커밋
	words1 := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n"
	if err := os.WriteFile(filepath.Join(workDir, "words.md"), []byte(words1), 0644); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}
	if _, err := g.Run("commit", "-m", "단어장 초기화"); err != nil {
		return err
	}

	// 규칙 위반 단어 "사과" 커밋 (표범의 끝 글자 "범"으로 시작해야 함)
	words2 := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n4. 사과\n"
	if err := os.WriteFile(filepath.Join(workDir, "words.md"), []byte(words2), 0644); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}
	_, err = g.Run("commit", "-m", "단어 '사과' 추가 (규칙 위반)")
	return err
}

func check(workDir string, g *git.Git) (exercise.Result, error) {
	// HEAD가 setup의 HEAD~1이어야 함 → 커밋 수가 1이어야 함
	count, _ := g.CommitCount()
	if count > 1 {
		// staged에 무엇이 있는지 확인
		staged, _ := g.StagedFiles()
		statusOut, _ := g.Run("status", "--short")

		if count == 2 {
			// --hard를 쓴 경우
			if len(staged) == 0 {
				hasSagwa := false
				for _, line := range exercise.Lines(statusOut) {
					if strings.Contains(line, "words.md") {
						hasSagwa = true
					}
				}
				if !hasSagwa {
					return exercise.Result{
						Passed:  false,
						Message: "❌ --hard 옵션을 쓴 것 같습니다. 변경이 사라졌습니다.\n   → 'git reset --soft HEAD~1' 을 사용해야 합니다.\n   (현재 커밋을 취소하고 다시 시도해 보세요. gitlings reset으로 초기화 가능)",
					}, nil
				}
				// --mixed를 쓴 경우 (변경은 있지만 staged 아님)
				return exercise.Result{
					Passed:  false,
					Message: "❌ --mixed 옵션(기본값)을 쓴 것 같습니다. 변경이 무대 밖에 있습니다.\n   → 'git reset --soft HEAD~1' 을 사용해야 합니다.",
				}, nil
			}
		}

		return exercise.Result{
			Passed:  false,
			Message: "❌ 아직 커밋을 취소하지 않았습니다.\n   → 'git reset --soft HEAD~1' 을 실행해 보세요.",
		}, nil
	}

	// 인덱스에 words.md 변경이 있어야 함
	staged, _ := g.StagedFiles()
	if !exercise.Contains(staged, "words.md") {
		// --hard vs --mixed 구분: --hard는 작업 트리도 깨끗함
		if g.IsClean() {
			return exercise.Result{
				Passed:  false,
				Message: "❌ --hard 옵션을 썼습니다. 변경이 사라졌습니다. 무엇을 추가했었나요?\n   → 'git reset --soft HEAD~1' 을 사용해야 staged 상태로 남습니다.",
			}, nil
		}
		// --mixed: 변경은 있지만 staged 아님
		statusOut, _ := g.Run("status", "--short")
		for _, line := range exercise.Lines(statusOut) {
			if strings.Contains(line, "words.md") {
				return exercise.Result{
					Passed:  false,
					Message: "❌ --mixed(기본값)를 썼습니다. 변경은 남았지만 무대 밖(unstaged)에 있습니다.\n   → 'git reset --soft HEAD~1' 을 사용해야 staged 상태로 남습니다.",
				}, nil
			}
		}
		return exercise.Result{
			Passed:  false,
			Message: "❌ 인덱스(무대)에 words.md 변경이 없습니다.\n   → 'git reset --soft HEAD~1' 을 실행해 보세요.",
		}, nil
	}

	return exercise.Result{
		Passed:  true,
		Message: "✅ 커밋을 취소하고 변경을 무대(staged)에 남겼습니다!",
	}, nil
}
