package exercises

import (
	"os"
	"path/filepath"

	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/git"
)

func init() {
	register(12, "12-reflog", "12★ 날린 커밋 되살리기", true, setupReflog, checkReflog)
}

// setup 때 저장해 두는 목표 커밋 해시 (check에서 비교용, .git 안에 숨겨 둔다)
const reflogTargetFile = ".gitlings_target_commit"

func setupReflog(dir string) error {
	g, err := git.InitMain(dir)
	if err != nil {
		return err
	}
	if err := g.SetupUserConfig(); err != nil {
		return err
	}
	commits := []struct{ words, msg string }{
		{"# 끝말잇기 단어장\n\n1. 기차\n", "단어 1개"},
		{"# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n", "단어 2개"},
		{"# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n", "단어 3개"},
	}
	for _, c := range commits {
		if err := exercise.WriteFiles(dir, map[string]string{"words.md": c.words}); err != nil {
			return err
		}
		if _, err := g.Run("add", "words.md"); err != nil {
			return err
		}
		if _, err := g.Run("commit", "-m", c.msg); err != nil {
			return err
		}
	}
	target, err := g.HeadCommit()
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, ".git", reflogTargetFile), []byte(target), 0o644); err != nil {
		return err
	}
	// 2개 날리기
	_, err = g.Run("reset", "--hard", "HEAD~2")
	return err
}

func checkReflog(dir string) exercise.Result {
	g := git.New(dir)
	target, err := os.ReadFile(filepath.Join(dir, ".git", reflogTargetFile))
	if err != nil {
		return exercise.Fail("setup 정보를 읽을 수 없습니다. 이 문제를 처음 상태로 되돌려 보세요(Ctrl+R)")
	}
	head, err := g.HeadCommit()
	if err != nil {
		return exercise.Fail("HEAD를 읽을 수 없습니다")
	}
	if head != string(target) {
		if count, _ := g.CommitCount(); count < 3 {
			return exercise.Fail("아직 날린 커밋이 복구되지 않았습니다. `git reflog`로 날리기 전 커밋의 해시를 찾아, 그 커밋으로 되돌아가 보세요")
		}
		return exercise.Fail("올바른 커밋으로 복구되지 않았습니다. `git reflog`로 `HEAD~2` 이전의 커밋 해시를 찾아보세요")
	}
	return exercise.Pass("날린 커밋을 reflog로 되살렸습니다!")
}
