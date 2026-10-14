package exercises

import (
	"os"

	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/git"
)

func init() {
	register(11, "11-pull-push", "11 push가 거절됐다", false, setupPullPush, checkPullPush)
}

const (
	pullPushRules      = "# 끝말잇기 규칙\n\n다음 단어의 첫 글자는 앞 단어의 끝 글자와 같아야 합니다.\n단어는 두 글자 이상이어야 합니다.\n"
	pullPushRulesFixed = pullPushRules + "\n(짝이 추가: 두음 법칙은 적용하지 않습니다.)\n"
)

func setupPullPush(dir string) error {
	g, err := git.InitMain(dir)
	if err != nil {
		return err
	}
	if err := g.SetupUserConfig(); err != nil {
		return err
	}
	// 공통 조상: words.md + rules.md
	if err := exercise.WriteFiles(dir, map[string]string{
		"words.md": "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n",
		"rules.md": pullPushRules,
	}); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md", "rules.md"); err != nil {
		return err
	}
	if _, err := g.Run("commit", "-m", "단어장 초기화"); err != nil {
		return err
	}

	// 가짜 원격 bare 저장소
	bareDir := dir + ".git-remote"
	if err := os.MkdirAll(bareDir, 0o755); err != nil {
		return err
	}
	if _, err := git.New(bareDir).Run("init", "--bare"); err != nil {
		return err
	}
	if _, err := g.Run("remote", "add", "origin", bareDir); err != nil {
		return err
	}
	if _, err := g.Run("push", "-u", "origin", "main"); err != nil {
		return err
	}

	// "짝"이 rules.md를 고쳐 origin에 push (짝의 저장소를 옆 폴더에 clone해서 push)
	partnerDir := dir + ".partner"
	if _, err := g.Run("clone", bareDir, partnerDir); err != nil {
		return err
	}
	gPartner := git.New(partnerDir)
	if err := gPartner.SetupUserConfig(); err != nil {
		return err
	}
	if err := exercise.WriteFiles(partnerDir, map[string]string{"rules.md": pullPushRulesFixed}); err != nil {
		return err
	}
	if _, err := gPartner.Run("add", "rules.md"); err != nil {
		return err
	}
	if _, err := gPartner.Run("commit", "-m", "rules.md: 두음 법칙 안내 추가"); err != nil {
		return err
	}
	if _, err := gPartner.Run("push", "origin", "main"); err != nil {
		return err
	}

	// 내 로컬에도 단어 추가 커밋 (push 전)
	if err := exercise.WriteFiles(dir, map[string]string{"words.md": "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n4. 범고래\n"}); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}
	_, err = g.Run("commit", "-m", "단어 '범고래' 추가")
	return err
}

func checkPullPush(dir string) exercise.Result {
	g := git.New(dir)
	if !g.LocalMatchesRemote("main", "origin") {
		return exercise.Fail("아직 push하지 않았거나 origin의 내용을 받지 않았습니다. 짝의 커밋을 먼저 받아 합친 뒤 다시 push해 보세요")
	}
	if mine, _ := g.Run("log", "--oneline", "--all", "--grep=범고래"); mine == "" {
		return exercise.Fail("내 커밋(범고래)이 이력에 없습니다")
	}
	if partner, _ := g.Run("log", "--oneline", "--all", "--grep=두음"); partner == "" {
		return exercise.Fail("짝의 커밋(`rules.md` 수정)이 이력에 없습니다. origin에서 짝의 커밋을 받아 오세요")
	}
	return exercise.Pass("짝의 커밋을 받아 합치고 push했습니다!")
}
