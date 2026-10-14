package exercises

import (
	"strings"

	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/git"
)

func init() { register(13, "13-stash", "13★ 하던 일 치워두기", true, setupStash, checkStash) }

func setupStash(dir string) error {
	g, err := git.InitMain(dir)
	if err != nil {
		return err
	}
	if err := g.SetupUserConfig(); err != nil {
		return err
	}
	if err := exercise.WriteFiles(dir, map[string]string{
		"words.md": "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n",
		"rules.md": "# 끝말잇기 규칙\n\n다음 단어의 첫 그자는 앞 단어의 끝 글자와 같아야 합니다.\n",
	}); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md", "rules.md"); err != nil {
		return err
	}
	if _, err := g.Run("commit", "-m", "단어장 초기화"); err != nil {
		return err
	}
	// new-word 브랜치에서 words.md 수정 중 (커밋 안 함)
	if _, err := g.Run("switch", "-c", "new-word"); err != nil {
		return err
	}
	return exercise.WriteFiles(dir, map[string]string{"words.md": "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n4. 범고래 (작업 중...)\n"})
}

func checkStash(dir string) exercise.Result {
	g := git.New(dir)
	// main의 rules.md에서 "그자" → "글자" 수정 여부
	rules, err := g.Run("show", "main:rules.md")
	if err != nil {
		return exercise.Fail("main의 `rules.md`를 읽을 수 없습니다")
	}
	if strings.Contains(rules, "첫 그자") {
		return exercise.Fail("main의 `rules.md` 오타가 아직 수정되지 않았습니다. main으로 이동해서 '첫 그자'를 '첫 글자'로 고치고 커밋하세요")
	}
	if n, err := g.Run("rev-list", "--count", "main"); err != nil || n == "1" {
		return exercise.Fail("main에 오타 수정 커밋이 없습니다. main 브랜치로 이동해서 `rules.md`를 수정하고 커밋하세요")
	}

	// new-word 브랜치에 하던 작업이 남아 있는지 (stash에서 되살렸거나 커밋했거나)
	if branch, _ := g.CurrentBranch(); branch != "new-word" {
		words, err := g.Run("show", "new-word:words.md")
		if err != nil || !strings.Contains(words, "범고래") {
			if stashes, _ := g.StashList(); len(stashes) > 0 {
				return exercise.Fail("stash에 작업이 남아 있습니다. new-word 브랜치로 돌아가서 치워둔 작업을 되살리세요")
			}
			return exercise.Fail("new-word 브랜치에 하던 작업(범고래)이 없습니다. 치워둔 작업을 되살려 보세요")
		}
		return exercise.Fail("현재 new-word 브랜치가 아닙니다. new-word 브랜치로 돌아가세요")
	}

	disk, _ := exercise.Read(dir, "words.md")
	head, _ := g.FileContent("words.md")
	if !strings.Contains(disk, "범고래") && !strings.Contains(head, "범고래") {
		return exercise.Fail("new-word에 하던 작업(범고래)이 없습니다. stash 목록에 치워둔 작업이 있는지 보고 되살려 보세요")
	}
	if stashes, _ := g.StashList(); len(stashes) > 0 {
		return exercise.Fail("stash 목록이 비어 있지 않습니다. 치워둔 작업을 꺼내면서 stash 목록에서도 지워지게 되살리세요")
	}
	return exercise.Pass("stash로 작업을 치워두고, main에서 오타를 수정하고, 다시 되살렸습니다!")
}
