package exercises

import (
	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/git"
)

func init() { register(8, "08-revert", "08 이미 올린 커밋", false, setupRevert, checkRevert) }

func setupRevert(dir string) error {
	g, err := git.InitMain(dir)
	if err != nil {
		return err
	}
	if err := g.SetupUserConfig(); err != nil {
		return err
	}
	// 초기 커밋, 그다음 잘못된 단어 커밋
	steps := []struct{ words, msg string }{
		{"# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n", "단어장 초기화"},
		{"# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n4. 사과\n", "단어 '사과' 추가 (규칙 위반)"},
	}
	for _, s := range steps {
		if err := exercise.WriteFiles(dir, map[string]string{"words.md": s.words}); err != nil {
			return err
		}
		if _, err := g.Run("add", "words.md"); err != nil {
			return err
		}
		if _, err := g.Run("commit", "-m", s.msg); err != nil {
			return err
		}
	}
	// 가짜 원격(bare 저장소)을 연습 폴더 옆에 만든다
	bareDir := dir + ".git-remote"
	if _, err := git.InitBareMain(bareDir); err != nil {
		return err
	}
	if _, err := g.Run("remote", "add", "origin", bareDir); err != nil {
		return err
	}
	_, err = g.Run("push", "-u", "origin", "main")
	return err
}

func checkRevert(dir string) exercise.Result {
	g := git.New(dir)
	// 강제 push를 막기 위해: 원래 잘못된 커밋이 이력에 남아 있어야 함
	badCommit, err := g.Run("log", "--oneline", "--all", "--grep=규칙 위반")
	if err != nil || badCommit == "" {
		// reset --hard 후 강제 push했을 가능성
		return exercise.Fail("원래 잘못된 커밋이 이력에 없습니다. 원격 이력을 지웠다면, 짝이 이미 받았다면 큰 문제가 됩니다! `git revert`로 되돌리세요 (reset + 강제 push가 아닌). Ctrl+R로 처음 상태로 되돌린 뒤 다시 시도하세요")
	}
	// revert 커밋이 있어야 함
	revertCommit, _ := g.Run("log", "--oneline", "--all", "--grep=Revert\\|revert\\|되돌")
	if revertCommit == "" {
		// revert 커밋을 찾는 다른 방법: 커밋 수
		if count, _ := g.CommitCount(); count < 3 {
			return exercise.Fail("아직 되돌린 커밋이 없습니다. 잘못된 커밋을 지우지 말고, 그 반대 변경을 새 커밋으로 만들어 보세요")
		}
	}
	// words.md가 잘못되기 전 상태여야 함
	content, err := g.FileContent("words.md")
	if err != nil {
		return exercise.Fail("`words.md`를 읽을 수 없습니다")
	}
	for _, w := range exercise.ParseWordsMd(content) {
		if w.Word == "사과" {
			return exercise.Fail("`words.md`에 '사과'가 아직 있습니다. `git revert`로 되돌리면 사과가 제거됩니다")
		}
	}
	if !g.LocalMatchesRemote("main", "origin") {
		return exercise.Fail("로컬 main이 origin/main과 다릅니다. 되돌린 커밋을 아직 `origin`에 올리지 않았습니다")
	}
	return exercise.Pass("이력을 보존하면서 잘못된 커밋을 되돌리고 push했습니다!")
}
