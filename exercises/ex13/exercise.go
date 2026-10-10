package ex13

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
		ID:       13,
		Slug:     "13-stash",
		Name:     "13★ 하던 일 치워두기",
		Optional: true,
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
	words := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n"
	rules := "# 끝말잇기 규칙\n\n다음 단어의 첫 그자는 앞 단어의 끝 글자와 같아야 합니다.\n"
	if err := os.WriteFile(filepath.Join(workDir, "words.md"), []byte(words), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(workDir, "rules.md"), []byte(rules), 0644); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md", "rules.md"); err != nil {
		return err
	}
	if _, err := g.Run("commit", "-m", "단어장 초기화"); err != nil {
		return err
	}

	// new-word 브랜치 생성 후 words.md 수정 (커밋 안 함)
	if _, err := g.Run("switch", "-c", "new-word"); err != nil {
		return err
	}
	words2 := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n4. 범고래 (작업 중...)\n"
	return os.WriteFile(filepath.Join(workDir, "words.md"), []byte(words2), 0644)
}

func check(workDir string, g *git.Git) (exercise.Result, error) {
	// main에 오타 수정 커밋이 있는지
	// rules.md의 "그자" → "글자" 수정 여부 확인
	rulesContent, err := g.Run("show", "main:rules.md")
	if err != nil {
		return exercise.Result{Passed: false, Message: "❌ main의 rules.md를 읽을 수 없습니다."}, nil
	}
	if strings.Contains(rulesContent, "첫 그자") {
		return exercise.Result{
			Passed:  false,
			Message: "❌ main의 rules.md 오타가 아직 수정되지 않았습니다.\n   → main으로 이동해서 '첫 그자'를 '첫 글자'로 고치고 커밋하세요.",
		}, nil
	}

	// main에 수정 커밋이 있는지 (2개 이상)
	mainCount, err := g.Run("rev-list", "--count", "main")
	if err != nil || mainCount == "1" {
		return exercise.Result{
			Passed:  false,
			Message: "❌ main에 오타 수정 커밋이 없습니다.\n   → main 브랜치로 이동해서 rules.md를 수정하고 커밋하세요.",
		}, nil
	}

	// new-word 브랜치에 하던 작업이 남아 있는지
	// (stash pop 했거나 커밋했거나)
	currentBranch, _ := g.CurrentBranch()

	// new-word에서 words.md 상태 확인
	if currentBranch != "new-word" {
		// new-word 브랜치의 HEAD에 작업이 커밋됐는지 확인
		newWordWords, err := g.Run("show", "new-word:words.md")
		if err != nil || !strings.Contains(newWordWords, "범고래") {
			// stash가 비어 있지 않을 수도 있음 - stash에 있으면 OK
			stashList, _ := g.StashList()
			if len(stashList) > 0 {
				return exercise.Result{
					Passed:  false,
					Message: "❌ stash에 작업이 남아 있습니다.\n   → new-word 브랜치로 이동해서 'git stash pop' 을 실행하세요.",
				}, nil
			}
			return exercise.Result{
				Passed:  false,
				Message: "❌ new-word 브랜치에 하던 작업(범고래)이 없습니다.\n   → 'git stash pop' 으로 작업을 되살려 보세요.",
			}, nil
		}
		return exercise.Result{
			Passed:  false,
			Message: "❌ 현재 new-word 브랜치가 아닙니다.\n   → 'git switch new-word' 로 이동하세요.",
		}, nil
	}

	// new-word 브랜치에서 words.md 확인 (작업 트리 또는 커밋)
	diskContent, _ := os.ReadFile(filepath.Join(workDir, "words.md"))
	headContent, _ := g.FileContent("words.md")
	hasBumgorae := strings.Contains(string(diskContent), "범고래") || strings.Contains(headContent, "범고래")
	if !hasBumgorae {
		return exercise.Result{
			Passed:  false,
			Message: "❌ new-word에 하던 작업(범고래)이 없습니다.\n   → 'git stash pop' 으로 작업을 되살려 보세요.",
		}, nil
	}

	// stash 목록이 비어 있어야 함
	stashList, _ := g.StashList()
	if len(stashList) > 0 {
		return exercise.Result{
			Passed:  false,
			Message: "❌ stash 목록이 비어 있지 않습니다.\n   → 'git stash pop' 으로 stash를 되살리세요.",
		}, nil
	}

	return exercise.Result{
		Passed:  true,
		Message: "✅ stash로 작업을 치워두고, main에서 오타를 수정하고, 다시 되살렸습니다!",
	}, nil
}
