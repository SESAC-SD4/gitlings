package exercises

import "testing"

func TestReflog(t *testing.T) {
	ex := ByID(12)
	t.Run("setup 직후에는 실패", func(t *testing.T) {
		mustFail(t, ex, prepare(t, ex), "아직 날린 커밋이 복구되지 않았습니다")
	})
	t.Run("reflog의 해시로 reset하면 통과", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, `git reset --hard "$(git rev-parse 'HEAD@{1}')"`)
		mustPass(t, ex, dir)
	})
}
