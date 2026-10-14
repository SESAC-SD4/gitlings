package exercises

import (
	"strings"

	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/git"
)

func init() {
	register(1, "01-first-commit", "01 첫 사진", false, setupFirstCommit, checkFirstCommit)
}

func setupFirstCommit(dir string) error {
	return exercise.WriteFiles(dir, map[string]string{"words.md": "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n"})
}

func checkFirstCommit(dir string) exercise.Result {
	g := git.New(dir)
	if !g.HasGitDir() {
		return exercise.Fail("아직 Git 저장소가 아닙니다. 이 폴더를 먼저 Git 저장소로 만들어야 합니다")
	}
	count, err := g.CommitCount()
	if err != nil || count == 0 {
		if staged, _ := g.StagedFiles(); len(staged) > 0 {
			return exercise.Fail("파일이 무대에 올라가 있지만 아직 커밋하지 않았습니다. `git commit`으로 사진을 찍어 보세요")
		}
		return exercise.Fail("사진을 아직 안 찍었습니다(커밋이 없습니다). 파일을 무대에 올린 뒤 커밋해 보세요")
	}
	if !g.FileAtHEAD("words.md") {
		return exercise.Fail("`words.md`가 커밋에 포함되지 않았습니다. `words.md`를 무대에 올린 뒤 커밋해 보세요")
	}
	if !g.IsClean() {
		if staged, _ := g.StagedFiles(); len(staged) > 0 {
			return exercise.Fail("파일이 무대에 올라가 있지만 커밋하지 않았습니다. `git commit`으로 사진을 찍어 보세요")
		}
		return exercise.Fail("작업 트리에 변경이 남아 있습니다. `git status`로 상태를 확인해 보세요")
	}
	msg := "첫 번째 커밋 완료! 단어장의 첫 사진을 찍었습니다"
	// user.name이 없으면 경고만 한다.
	if out, _ := g.Run("log", "-1", "--format=%an"); strings.TrimSpace(out) == "" {
		msg += "\n\n⚠️ 커밋 작성자 이름이 없습니다. `git config`로 `user.name`을 설정해 보세요"
	}
	return exercise.Pass(msg)
}
