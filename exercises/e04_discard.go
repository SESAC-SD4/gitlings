package exercises

import (
	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/git"
)

func init() {
	register(4, "04-discard", "04 망친 수정 버리기", false, setupDiscard, checkDiscard)
}

const discardCleanWords = "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n"

func setupDiscard(dir string) error {
	g, err := git.InitMain(dir)
	if err != nil {
		return err
	}
	if err := g.SetupUserConfig(); err != nil {
		return err
	}
	if err := exercise.WriteFiles(dir, map[string]string{"words.md": discardCleanWords}); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}
	if _, err := g.Run("commit", "-m", "단어장 초기화"); err != nil {
		return err
	}
	// 실수로 엉망이 된 words.md (add 안 함)
	return exercise.WriteFiles(dir, map[string]string{"words.md": "이게 뭔가요\n????\nasdfasdf\n제발 살려주세요\n"})
}

func checkDiscard(dir string) exercise.Result {
	g := git.New(dir)
	if count, _ := g.CommitCount(); count > 1 {
		return exercise.Fail("커밋이 새로 생겼습니다. 커밋을 추가하는 게 아니라 작업 트리의 `words.md`를 마지막 커밋 상태로 되돌려야 합니다")
	}
	if !g.IsClean() {
		return exercise.Fail("`words.md`가 아직 변경된 상태입니다. `git diff`로 무엇이 바뀌었는지 보고, 그 변경을 버려 보세요")
	}
	return exercise.Pass("실수한 수정을 버리고 마지막 커밋 상태로 되돌렸습니다!")
}
