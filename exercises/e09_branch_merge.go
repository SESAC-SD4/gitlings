package exercises

import (
	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/git"
)

func init() {
	register(9, "09-branch-merge", "09 갈래 만들고 합치기", false, setupBranchMerge, checkBranchMerge)
}

func setupBranchMerge(dir string) error {
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
	_, err = g.Run("commit", "-m", "단어장 초기화")
	return err
}

func checkBranchMerge(dir string) exercise.Result {
	g := git.New(dir)
	if !g.BranchExists("new-word") {
		return exercise.Fail("`new-word` 브랜치가 없습니다. `git branch`로 지금 있는 브랜치를 확인해 보세요")
	}
	if cur, _ := g.CurrentBranch(); cur != "main" {
		return exercise.Fail("현재 main 브랜치가 아닙니다. 합칠 때는 받는 쪽인 main으로 돌아와서 merge해야 합니다")
	}
	// new-word 브랜치의 커밋이 main에 포함됐는지
	newWordHead, err := g.Run("rev-parse", "new-word")
	if err != nil {
		return exercise.Fail("`new-word` 브랜치를 읽을 수 없습니다")
	}
	if !g.CommitContains("main", newWordHead) {
		return exercise.Fail("`new-word` 브랜치의 커밋이 아직 main에 병합되지 않았습니다")
	}
	content, err := g.FileContent("words.md")
	if err != nil {
		return exercise.Fail("`words.md`를 읽을 수 없습니다")
	}
	words := exercise.ParseWordsMd(content)
	if len(words) < 3 {
		return exercise.Fail("단어가 2개뿐입니다. `new-word` 브랜치에서 단어를 추가해야 합니다")
	}
	if err := exercise.ValidateWordChain(words); err != nil {
		return exercise.Fail("단어장 끝말잇기 규칙 위반: " + err.Error() + ". 규칙에 맞는 단어를 추가해 보세요 (차표의 끝 글자 '표'로 시작하는 단어)")
	}
	return exercise.Pass("new-word 브랜치에서 단어를 추가하고 main에 병합했습니다!")
}
