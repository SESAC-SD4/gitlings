package exercises

import (
	"strings"

	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/git"
)

func init() {
	register(5, "05-amend", "05 메시지 오타", false, setupAmend, checkAmend)
}

const amendWant = "단어 '표범' 추가"

func setupAmend(dir string) error {
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
	// 오타 있는 커밋 메시지로 커밋
	if err := exercise.WriteFiles(dir, map[string]string{"words.md": "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n"}); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}
	_, err = g.Run("commit", "-m", "단어 추가함ㅁㅁㄴㅇ")
	return err
}

func checkAmend(dir string) exercise.Result {
	g := git.New(dir)
	if count, _ := g.CommitCount(); count > 2 {
		return exercise.Fail("커밋이 하나 늘었습니다. 새로 찍은 게 아니라 마지막 커밋을 고쳐야 합니다. Ctrl+R로 처음 상태로 되돌려 다시 해 보세요")
	}
	msg, err := g.HeadMessage()
	if err != nil {
		return exercise.Fail("커밋 메시지를 읽을 수 없습니다")
	}
	msg = strings.TrimSpace(msg)
	if msg != amendWant {
		if strings.Contains(msg, "ㅁㅁㄴㅇ") {
			return exercise.Fail("아직 오타가 남아 있습니다. `git log --oneline`으로 마지막 커밋 메시지를 확인하고, 그 커밋을 고쳐 보세요")
		}
		return exercise.Fail("커밋 메시지가 다릅니다.\n현재: " + msg + "\n원하는: " + amendWant)
	}
	return exercise.Pass("커밋 메시지를 올바르게 수정했습니다!")
}
