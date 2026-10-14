package exercises

import (
	"strings"

	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/git"
)

func init() {
	register(2, "02-pick-files", "02 골라서 찍기", false, setupPickFiles, checkPickFiles)
}

func setupPickFiles(dir string) error {
	g, err := git.InitMain(dir)
	if err != nil {
		return err
	}
	if err := g.SetupUserConfig(); err != nil {
		return err
	}
	if err := exercise.WriteFiles(dir, map[string]string{"words.md": "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n"}); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}
	if _, err := g.Run("commit", "-m", "단어장 초기화"); err != nil {
		return err
	}
	// 3번째 단어 추가와 낙서 (둘 다 unstaged)
	return exercise.WriteFiles(dir, map[string]string{
		"words.md": "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n",
		"notes.md": "오늘 공부한 것: git add, git commit\n아직 잘 모르겠음...\n",
	})
}

func checkPickFiles(dir string) exercise.Result {
	g := git.New(dir)
	count, err := g.CommitCount()
	if err != nil || count < 2 {
		return exercise.Fail("새 커밋이 없습니다. `words.md`만 add하고 commit해 보세요")
	}
	out, _ := g.Run("diff", "--name-only", "HEAD~1", "HEAD")
	changed := exercise.Lines(out)
	if exercise.Contains(changed, "notes.md") {
		return exercise.Fail("낙서(`notes.md`)까지 사진에 찍혔습니다. 커밋을 취소하고 `words.md`만 add해서 다시 커밋해 보세요")
	}
	if !exercise.Contains(changed, "words.md") {
		return exercise.Fail("새 커밋에 `words.md` 변경이 없습니다. `words.md`를 무대에 올린 뒤 커밋해 보세요")
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
		return exercise.Fail("`notes.md`가 수정된 상태로 남아 있지 않습니다. `notes.md`는 add 없이 그냥 두세요")
	}
	content, err := g.FileContent("words.md")
	if err != nil {
		return exercise.Fail("`words.md`를 읽을 수 없습니다")
	}
	if err := exercise.ValidateWordChain(exercise.ParseWordsMd(content)); err != nil {
		return exercise.Fail("단어장 끝말잇기 규칙 위반: " + err.Error())
	}
	return exercise.Pass("`words.md`만 커밋하고 `notes.md`는 수정 상태로 남겼습니다!")
}
