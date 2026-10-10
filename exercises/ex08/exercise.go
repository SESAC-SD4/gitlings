package ex08

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
		ID:       8,
		Slug:     "08-revert",
		Name:     "08 이미 올린 커밋",
		Optional: false,
		Desc:     desc,
		Hints:    exercise.ParseHints(hintsRaw),
		Setup:    setup,
		Check:    check,
	}
}

func setup(workDir string, _ *git.Git) error {
	// 로컬 저장소 생성
	g, err := git.InitMain(workDir)
	if err != nil {
		return err
	}
	if err := g.SetupUserConfig(); err != nil {
		return err
	}

	// 초기 커밋
	words1 := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n"
	if err := os.WriteFile(filepath.Join(workDir, "words.md"), []byte(words1), 0644); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}
	if _, err := g.Run("commit", "-m", "단어장 초기화"); err != nil {
		return err
	}

	// 잘못된 단어 커밋
	words2 := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n4. 사과\n"
	if err := os.WriteFile(filepath.Join(workDir, "words.md"), []byte(words2), 0644); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}
	if _, err := g.Run("commit", "-m", "단어 '사과' 추가 (규칙 위반)"); err != nil {
		return err
	}

	// 가짜 원격(bare 저장소) 만들기
	bareDir := workDir + ".git-remote"
	if err := os.MkdirAll(bareDir, 0755); err != nil {
		return err
	}
	gBare := git.New(bareDir)
	if _, err := gBare.Run("init", "--bare"); err != nil {
		return err
	}

	// origin으로 등록 후 push
	if _, err := g.Run("remote", "add", "origin", bareDir); err != nil {
		return err
	}
	_, err = g.Run("push", "-u", "origin", "main")
	return err
}

func check(workDir string, g *git.Git) (exercise.Result, error) {
	// 강제 push를 막기 위해: 원격이 setup 때의 커밋을 포함하는지 확인
	// 원래 잘못된 커밋이 이력에 남아 있어야 함
	badCommit, err := g.Run("log", "--oneline", "--all", "--grep=규칙 위반")
	if err != nil || badCommit == "" {
		// reset --hard 후 강제 push했을 가능성
		return exercise.Result{
			Passed:  false,
			Message: "❌ 원래 잘못된 커밋이 이력에 없습니다.\n   원격 이력을 지웠다면, 짝이 이미 받았다면 큰 문제가 됩니다!\n   → git revert로 되돌리세요 (reset + 강제 push가 아닌).\n   gitlings reset으로 초기화 후 다시 시도하세요.",
		}, nil
	}

	// revert 커밋이 있어야 함
	revertCommit, _ := g.Run("log", "--oneline", "--all", "--grep=Revert\\|revert\\|되돌")
	if revertCommit == "" {
		// revert 커밋을 찾는 다른 방법: 마지막 커밋과 그 전 커밋의 트리가 같은지
		count, _ := g.CommitCount()
		if count < 3 {
			return exercise.Result{
				Passed:  false,
				Message: "❌ 아직 되돌린 커밋이 없습니다.\n   → 'git revert HEAD --no-edit' 를 실행해 보세요.",
			}, nil
		}
	}

	// words.md가 잘못되기 전 상태여야 함
	content, err := g.FileContent("words.md")
	if err != nil {
		return exercise.Result{Passed: false, Message: "❌ words.md를 읽을 수 없습니다."}, nil
	}
	words := exercise.ParseWordsMd(content)
	// 사과가 없어야 함
	for _, w := range words {
		if w.Word == "사과" {
			return exercise.Result{
				Passed:  false,
				Message: "❌ words.md에 '사과'가 아직 있습니다.\n   → git revert로 되돌리면 사과가 제거됩니다.",
			}, nil
		}
	}

	// 로컬 main == origin/main
	if !g.LocalMatchesRemote("main", "origin") {
		return exercise.Result{
			Passed:  false,
			Message: "❌ 로컬 main이 origin/main과 다릅니다.\n   → 'git push origin main' 을 실행해 보세요.",
		}, nil
	}

	return exercise.Result{
		Passed:  true,
		Message: "✅ 이력을 보존하면서 잘못된 커밋을 되돌리고 push했습니다!",
	}, nil
}
