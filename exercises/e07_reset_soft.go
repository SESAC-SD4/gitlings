package exercises

import (
	"strings"

	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/git"
)

func init() {
	register(7, "07-reset-soft", "07 방금 커밋 취소", false, setupResetSoft, checkResetSoft)
}

func setupResetSoft(dir string) error {
	g, err := git.InitMain(dir)
	if err != nil {
		return err
	}
	if err := g.SetupUserConfig(); err != nil {
		return err
	}
	// 초기 커밋, 그다음 규칙 위반 단어 "사과" 커밋 (표범의 끝 글자 "범"으로 시작해야 함)
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
	return nil
}

func checkResetSoft(dir string) exercise.Result {
	g := git.New(dir)
	// "사과" 커밋을 취소했으면 커밋 수가 1이어야 함
	count, _ := g.CommitCount()
	// 커밋이 setup 그대로 2개면 아직 reset하지 않은 것이다. --hard·--mixed 구분은 커밋이 1개일 때만 한다.
	if count > 1 {
		return exercise.Fail("아직 커밋을 취소하지 않았습니다. `git log --oneline`에 \"사과\" 커밋이 남아 있습니다")
	}

	// 인덱스에 words.md 변경이 있어야 함
	staged, _ := g.StagedFiles()
	if !exercise.Contains(staged, "words.md") {
		// --hard vs --mixed 구분: --hard는 작업 트리도 깨끗함
		if g.IsClean() {
			return exercise.Fail("--hard 옵션을 썼습니다. 변경이 사라졌습니다. 변경을 무대(staged)에 남기는 옵션을 써야 합니다. Ctrl+R로 처음 상태로 되돌린 뒤 다시 해 보세요")
		}
		statusOut, _ := g.Run("status", "--short")
		for _, line := range exercise.Lines(statusOut) {
			if strings.Contains(line, "words.md") {
				return exercise.Fail("--mixed(기본값)를 썼습니다. 변경은 남았지만 무대 아래(unstaged)에 있습니다. 변경을 무대(staged)에 남기는 옵션을 써야 합니다. Ctrl+R로 처음 상태로 되돌린 뒤 다시 해 보세요")
			}
		}
		return exercise.Fail("무대(staged)에 `words.md` 변경이 없습니다. 취소한 변경이 무대(staged)에 남아 있어야 합니다")
	}
	return exercise.Pass("커밋을 취소하고 변경을 무대(staged)에 남겼습니다!")
}
