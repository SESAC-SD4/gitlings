package exercises

import (
	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/git"
)

func init() {
	register(6, "06-gitignore", "06 비밀이 섞였다", false, setupGitignore, checkGitignore)
}

func setupGitignore(dir string) error {
	g, err := git.InitMain(dir)
	if err != nil {
		return err
	}
	if err := g.SetupUserConfig(); err != nil {
		return err
	}
	if err := exercise.WriteFiles(dir, map[string]string{
		"words.md": "# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n",
		".env":     "ADMIN_PASSWORD=super_secret_1234\n",
	}); err != nil {
		return err
	}
	if _, err := g.Run("add", "words.md", ".env"); err != nil {
		return err
	}
	_, err = g.Run("commit", "-m", "단어장 및 환경설정 추가")
	return err
}

func checkGitignore(dir string) exercise.Result {
	g := git.New(dir)
	// .env가 HEAD에 없어야 함
	if g.FileAtHEAD(".env") {
		// .gitignore는 커밋했지만 .env를 추적에서 빼지 않은 경우
		if g.FileAtHEAD(".gitignore") && g.IsTracked(".env") {
			return exercise.Fail("`.gitignore`에 `.env`를 넣었지만 아직 추적 중입니다. `.gitignore`는 이미 추적 중인 파일을 빼 주지 않습니다. `.env`를 추적에서 빼고 커밋하세요")
		}
		return exercise.Fail("`.env`가 아직 커밋에 포함되어 있습니다. 파일은 남기고 Git 추적에서만 뺀 뒤 커밋하세요")
	}
	// .env 파일이 디스크에 있어야 함
	if !exercise.IsFile(dir, ".env") {
		return exercise.Fail("`.env` 파일이 삭제됐습니다. 파일은 디스크에 남겨 두고 Git 추적에서만 빼야 합니다. Ctrl+R로 처음 상태로 되돌린 뒤 다시 해 보세요")
	}
	if !g.FileAtHEAD(".gitignore") {
		return exercise.Fail("`.gitignore` 파일이 커밋에 없습니다. `.gitignore`를 만들고 커밋하세요")
	}
	if !g.IsIgnored(".env") {
		return exercise.Fail("`.gitignore`가 `.env`를 무시하지 않습니다. `.gitignore` 파일에 `.env` 줄을 추가하세요")
	}
	return exercise.Pass(".env를 추적에서 제거하고 .gitignore에 등록했습니다!\n\n⚠️ 이전 커밋에는 비밀번호가 여전히 남아 있습니다. 진짜 상황이라면 비밀번호를 바꿔야 합니다")
}
