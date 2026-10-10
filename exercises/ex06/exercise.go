package ex06

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
		ID:       6,
		Slug:     "06-gitignore",
		Name:     "06 비밀이 섞였다",
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

	words := "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n"
	if err := os.WriteFile(filepath.Join(workDir, "words.md"), []byte(words), 0644); err != nil {
		return err
	}
	env := "ADMIN_PASSWORD=super_secret_1234\n"
	if err := os.WriteFile(filepath.Join(workDir, ".env"), []byte(env), 0644); err != nil {
		return err
	}

	if _, err := g.Run("add", "words.md", ".env"); err != nil {
		return err
	}
	_, err = g.Run("commit", "-m", "단어장 및 환경설정 추가")
	return err
}

func check(workDir string, g *git.Git) (exercise.Result, error) {
	// .env가 HEAD에 없어야 함
	if g.FileAtHEAD(".env") {
		// .gitignore만 있고 rm --cached 안 한 경우
		if g.FileAtHEAD(".gitignore") && g.IsTracked(".env") {
			return exercise.Result{
				Passed:  false,
				Message: "❌ .gitignore에 .env를 넣었지만 아직 추적 중입니다.\n   → 'git rm --cached .env' 로 추적에서 제거하고 커밋하세요.",
			}, nil
		}
		return exercise.Result{
			Passed:  false,
			Message: "❌ .env가 아직 커밋에 포함되어 있습니다.\n   → 'git rm --cached .env' 실행 후 커밋하세요.",
		}, nil
	}

	// .env 파일이 디스크에 있어야 함
	if _, err := os.Stat(filepath.Join(workDir, ".env")); os.IsNotExist(err) {
		return exercise.Result{
			Passed:  false,
			Message: "❌ .env 파일이 삭제됐습니다. 파일은 디스크에 남겨두어야 합니다.\n   → git restore .env 또는 echo로 파일을 다시 만들고 git rm --cached .env를 사용하세요.",
		}, nil
	}

	// .gitignore가 HEAD에 있어야 함
	if !g.FileAtHEAD(".gitignore") {
		return exercise.Result{
			Passed:  false,
			Message: "❌ .gitignore 파일이 커밋에 없습니다.\n   → .gitignore를 만들고 커밋하세요.",
		}, nil
	}

	// .gitignore가 .env를 무시하는지 확인
	if !g.IsIgnored(".env") {
		return exercise.Result{
			Passed:  false,
			Message: "❌ .gitignore가 .env를 무시하지 않습니다.\n   → .gitignore 파일에 '.env' 줄을 추가하세요.",
		}, nil
	}

	return exercise.Result{
		Passed:  true,
		Message: "✅ .env를 추적에서 제거하고 .gitignore에 등록했습니다!\n\n⚠️  이전 커밋에는 비밀번호가 여전히 남아 있습니다. 진짜 상황이라면 비밀번호를 바꿔야 합니다.",
	}, nil
}
