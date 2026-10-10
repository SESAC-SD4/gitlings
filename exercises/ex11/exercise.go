package ex11

import (
	_ "embed"
	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/git"
	"os"
	"path/filepath"
)

//go:embed desc.txt
var desc string

//go:embed hints.txt
var hintsRaw string

func New() *exercise.Exercise {
	return &exercise.Exercise{
		ID:       11,
		Slug:     "11-pull-push",
		Name:     "11 push가 거절됐다",
		Optional: false,
		Desc:     desc,
		Hints:    exercise.ParseHints(hintsRaw),
		Setup:    setup,
		Check:    check,
	}
}

func setup(workDir string, _ *git.Git) error {
	// 로컬 저장소
	g, err := git.InitMain(workDir)
	if err != nil {
		return err
	}
	if err := g.SetupUserConfig(); err != nil {
		return err
	}

	// 공통 조상: words.md + rules.md
	words := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n"
	rules := "# 끝말잇기 규칙\n\n다음 단어의 첫 글자는 앞 단어의 끝 글자와 같아야 합니다.\n단어는 두 글자 이상이어야 합니다.\n"
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

	// 가짜 원격 bare 저장소
	bareDir := workDir + ".git-remote"
	if err := os.MkdirAll(bareDir, 0755); err != nil {
		return err
	}
	gBare := git.New(bareDir)
	if _, err := gBare.Run("init", "--bare"); err != nil {
		return err
	}
	if _, err := g.Run("remote", "add", "origin", bareDir); err != nil {
		return err
	}
	if _, err := g.Run("push", "-u", "origin", "main"); err != nil {
		return err
	}

	// "짝"이 rules.md를 고쳐 origin에 push (로컬에는 없는 상태 시뮬레이션)
	// 짝의 저장소를 임시로 clone해서 push
	partnerDir := workDir + ".partner"
	// clone은 기존 디렉토리에서 실행해야 함 (workDir 사용)
	if _, err := g.Run("clone", bareDir, partnerDir); err != nil {
		return err
	}
	gPartner := git.New(partnerDir)
	if err := gPartner.SetupUserConfig(); err != nil {
		return err
	}
	rulesFixed := "# 끝말잇기 규칙\n\n다음 단어의 첫 글자는 앞 단어의 끝 글자와 같아야 합니다.\n단어는 두 글자 이상이어야 합니다.\n\n(짝이 추가: 두음 법칙은 적용하지 않습니다.)\n"
	if err := os.WriteFile(filepath.Join(partnerDir, "rules.md"), []byte(rulesFixed), 0644); err != nil {
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
	words2 := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n4. 범고래\n"
	if err := os.WriteFile(filepath.Join(workDir, "words.md"), []byte(words2), 0644); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}
	_, err = g.Run("commit", "-m", "단어 '범고래' 추가")
	return err
}

func check(workDir string, g *git.Git) (exercise.Result, error) {
	// 로컬 main == origin/main
	if !g.LocalMatchesRemote("main", "origin") {
		return exercise.Result{
			Passed:  false,
			Message: "❌ 아직 push하지 않았거나 origin의 내용을 받지 않았습니다.\n   → 'git pull --no-rebase' 후 'git push origin main' 을 실행해 보세요.",
		}, nil
	}

	// 내 커밋이 이력에 있어야 함
	myCommit, _ := g.Run("log", "--oneline", "--all", "--grep=범고래")
	if myCommit == "" {
		return exercise.Result{
			Passed:  false,
			Message: "❌ 내 커밋(범고래)이 이력에 없습니다.",
		}, nil
	}

	// 짝의 커밋도 이력에 있어야 함
	partnerCommit, _ := g.Run("log", "--oneline", "--all", "--grep=두음")
	if partnerCommit == "" {
		return exercise.Result{
			Passed:  false,
			Message: "❌ 짝의 커밋(rules.md 수정)이 이력에 없습니다.\n   → 'git pull --no-rebase' 로 짝의 커밋을 받아 오세요.",
		}, nil
	}

	return exercise.Result{
		Passed:  true,
		Message: "✅ 짝의 커밋을 받아 합치고 push했습니다!",
	}, nil
}
