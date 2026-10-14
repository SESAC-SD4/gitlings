package exercises

import "testing"

func TestGitignore(t *testing.T) {
	ex := ByID(6)
	t.Run("setup 직후에는 실패", func(t *testing.T) {
		mustFail(t, ex, prepare(t, ex), "아직 커밋에 포함")
	})
	t.Run("gitignore·rm --cached·commit이면 통과", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, `echo ".env" > .gitignore && git rm --cached .env && git add .gitignore && git commit -m ".env를 gitignore에 추가"`)
		mustPass(t, ex, dir)
	})
	t.Run("gitignore만 추가하면 추적 중 안내", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, `echo ".env" > .gitignore && git add .gitignore && git commit -m "gitignore만 추가"`)
		mustFail(t, ex, dir, "아직 추적 중")
	})
	t.Run(".env를 지우면 안내", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, `git rm .env && echo ".env" > .gitignore && git add .gitignore && git commit -m ".env 삭제"`)
		mustFail(t, ex, dir, "삭제됐습니다")
	})
}
