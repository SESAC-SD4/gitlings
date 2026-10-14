package exercises

import (
	"fmt"
	"strings"

	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/git"
)

func init() { register(10, "10-conflict", "10 충돌", false, setupConflict, checkConflict) }

func setupConflict(dir string) error {
	g, err := git.InitMain(dir)
	if err != nil {
		return err
	}
	if err := g.SetupUserConfig(); err != nil {
		return err
	}
	// 공통 조상 커밋: 단어 2개 → partner 브랜치: 규칙 위반 "사과" → main: 규칙 맞는 "표범"
	steps := []struct{ branch, words, msg string }{
		{"", "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n", "단어장 초기화"},
		{"-c partner", "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 사과\n", "짝의 단어 '사과' 추가"},
		{"main", "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n", "단어 '표범' 추가"},
	}
	for _, s := range steps {
		if s.branch != "" {
			if _, err := g.Run(append([]string{"switch"}, strings.Fields(s.branch)...)...); err != nil {
				return err
			}
		}
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

func conflictMarker(line string) bool {
	return strings.Contains(line, "<<<<<<<") || strings.Contains(line, "=======") || strings.Contains(line, ">>>>>>>")
}

func checkConflict(dir string) exercise.Result {
	g := git.New(dir)
	// 충돌 표시가 없어야 함
	if content, ok := exercise.Read(dir, "words.md"); ok {
		for i, line := range exercise.Lines(content) {
			if conflictMarker(line) {
				return exercise.Fail(fmt.Sprintf("`words.md`에 충돌 표시가 남아 있습니다(%d번 줄). 파일을 열어 `<<<<<<<` / `=======` / `>>>>>>>` 표시를 모두 지우고 올바른 내용만 남기세요", i+1))
			}
		}
	}
	// 병합 커밋인지 확인 (부모 2개)
	if n, err := g.CommitParentCount(); err != nil || n < 2 {
		return exercise.Fail("아직 병합 커밋이 없습니다. `partner` 브랜치를 `main`에 병합하고, 충돌을 해결한 뒤 커밋까지 해야 합니다")
	}
	content, err := g.FileContent("words.md")
	if err != nil {
		return exercise.Fail("`words.md`를 읽을 수 없습니다")
	}
	words := exercise.ParseWordsMd(content)
	if err := exercise.ValidateWordChain(words); err != nil {
		return exercise.Fail("단어장 끝말잇기 규칙 위반: " + err.Error())
	}
	for _, w := range words {
		if w.Number == 3 {
			if w.Word != "표범" {
				return exercise.Fail("3번 단어가 '표범'이 아닙니다(현재: " + w.Word + "). 충돌 해결 시 끝말잇기 규칙에 맞는 '표범'을 선택하세요")
			}
			break
		}
	}
	return exercise.Pass("충돌을 해결하고 올바른 단어를 선택했습니다!")
}
