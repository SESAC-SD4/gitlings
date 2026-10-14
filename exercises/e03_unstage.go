package exercises

import (
	"strings"

	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/git"
)

func init() {
	register(3, "03-unstage", "03 무대에서 내리기", false, setupUnstage, checkUnstage)
}

func setupUnstage(dir string) error {
	g, err := git.InitMain(dir)
	if err != nil {
		return err
	}
	if err := g.SetupUserConfig(); err != nil {
		return err
	}
	if err := exercise.WriteFiles(dir, map[string]string{"words.md": "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n"}); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}
	if _, err := g.Run("commit", "-m", "단어장 초기화"); err != nil {
		return err
	}
	// 4번째 단어와 낙서를 모두 stage
	if err := exercise.WriteFiles(dir, map[string]string{
		"words.md": "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n4. 범고래\n",
		"notes.md": "오늘은 unstage 연습!\n",
	}); err != nil {
		return err
	}
	_, err = g.Run("add", "words.md", "notes.md")
	return err
}

func checkUnstage(dir string) exercise.Result {
	g := git.New(dir)
	if g.FileAtHEAD("notes.md") {
		return exercise.Fail("`notes.md`가 커밋에 포함됐습니다. 커밋하기 전에 `notes.md`를 무대에서 내려야 합니다")
	}
	if !exercise.Exists(dir, "notes.md") {
		return exercise.Fail("`notes.md` 파일이 삭제됐습니다. 파일은 지우지 말고 무대에서만 내리세요")
	}
	if !g.FileAtHEAD("words.md") {
		return exercise.Fail("`words.md`가 커밋에 없습니다. `words.md`를 커밋해야 합니다")
	}
	if count, _ := g.CommitCount(); count < 2 {
		return exercise.Fail("아직 커밋하지 않았습니다. `notes.md`를 무대에서 내린 후 `words.md`를 커밋해 보세요")
	}
	out, _ := g.Run("diff", "--name-only", "HEAD~1", "HEAD")
	if exercise.Contains(exercise.Lines(out), "notes.md") {
		return exercise.Fail("마지막 커밋에 `notes.md`가 포함됐습니다. 커밋을 취소하고 다시 시도해 보세요")
	}
	statusOut, _ := g.Run("status", "--short")
	notesStatus := ""
	for _, line := range exercise.Lines(statusOut) {
		if strings.Contains(line, "notes.md") {
			notesStatus = line
		}
	}
	if notesStatus == "" {
		return exercise.Fail("`notes.md`가 작업 트리에 없습니다")
	}
	return exercise.Pass("`notes.md`를 무대에서 내리고 `words.md`만 커밋했습니다!")
}
