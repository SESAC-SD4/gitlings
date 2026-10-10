package ex10

import (
	_ "embed"
	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/git"
	"os"
	"path/filepath"
	"strings"
)

//go:embed desc.txt
var desc string

//go:embed hints.txt
var hintsRaw string

func New() *exercise.Exercise {
	return &exercise.Exercise{
		ID:       10,
		Slug:     "10-conflict",
		Name:     "10 충돌",
		Optional: false,
		Desc:     desc,
		Hints:    exercise.ParseHints(hintsRaw),
		Setup:    setup,
		Check:    check,
	}
}

func setup(workDir string, _ *git.Git) error {
	g, err := git.InitMain(workDir)
	if err != nil {
		return err
	}
	if err := g.SetupUserConfig(); err != nil {
		return err
	}

	// 공통 조상 커밋: 단어 2개
	words := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n"
	if err := os.WriteFile(filepath.Join(workDir, "words.md"), []byte(words), 0644); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}
	if _, err := g.Run("commit", "-m", "단어장 초기화"); err != nil {
		return err
	}

	// partner 브랜치: 규칙 위반 단어 "사과"
	if _, err := g.Run("switch", "-c", "partner"); err != nil {
		return err
	}
	wordsPartner := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 사과\n"
	if err := os.WriteFile(filepath.Join(workDir, "words.md"), []byte(wordsPartner), 0644); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}
	if _, err := g.Run("commit", "-m", "짝의 단어 '사과' 추가"); err != nil {
		return err
	}

	// main으로 돌아와 규칙 맞는 단어 "표범" 추가
	if _, err := g.Run("switch", "main"); err != nil {
		return err
	}
	wordsMain := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n"
	if err := os.WriteFile(filepath.Join(workDir, "words.md"), []byte(wordsMain), 0644); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md"); err != nil {
		return err
	}
	_, err = g.Run("commit", "-m", "단어 '표범' 추가")
	return err
}

func check(workDir string, g *git.Git) (exercise.Result, error) {
	// 충돌 표시가 없어야 함
	wordsContent, err := os.ReadFile(filepath.Join(workDir, "words.md"))
	if err == nil {
		content := string(wordsContent)
		if strings.Contains(content, "<<<<<<<") || strings.Contains(content, "=======") || strings.Contains(content, ">>>>>>>") {
			// 어느 줄인지 찾기
			for i, line := range exercise.Lines(content) {
				if strings.Contains(line, "<<<<<<<") || strings.Contains(line, "=======") || strings.Contains(line, ">>>>>>>") {
					return exercise.Result{
						Passed:  false,
						Message: "❌ words.md에 충돌 표시가 남아 있습니다 (" + lineName(i+1) + " 번 줄).\n   → 파일을 열어 <<<<<<< / ======= / >>>>>>> 표시를 모두 지우고 올바른 내용만 남기세요.",
					}, nil
				}
			}
		}
	}

	// 병합 커밋인지 확인 (부모 2개)
	parentCount, err := g.CommitParentCount()
	if err != nil || parentCount < 2 {
		return exercise.Result{
			Passed:  false,
			Message: "❌ 아직 병합 커밋이 없습니다.\n   → 'git merge partner' 후 충돌을 해결하고 커밋하세요.",
		}, nil
	}

	// 단어장 규칙 확인
	content, err := g.FileContent("words.md")
	if err != nil {
		return exercise.Result{Passed: false, Message: "❌ words.md를 읽을 수 없습니다."}, nil
	}
	words := exercise.ParseWordsMd(content)
	if err := exercise.ValidateWordChain(words); err != nil {
		return exercise.Result{
			Passed:  false,
			Message: "❌ 단어장 끝말잇기 규칙 위반: " + err.Error(),
		}, nil
	}

	// 3번 단어가 표범이어야 함
	for _, w := range words {
		if w.Number == 3 {
			if w.Word != "표범" {
				return exercise.Result{
					Passed:  false,
					Message: "❌ 3번 단어가 '표범'이 아닙니다 (현재: " + w.Word + ").\n   → 충돌 해결 시 끝말잇기 규칙에 맞는 '표범'을 선택하세요.",
				}, nil
			}
			break
		}
	}

	return exercise.Result{
		Passed:  true,
		Message: "✅ 충돌을 해결하고 올바른 단어를 선택했습니다!",
	}, nil
}

func lineName(n int) string {
	if n == 0 {
		return "0"
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
