package exercises

import "testing"

func TestRevert(t *testing.T) {
	ex := ByID(8)
	t.Run("setup 직후에는 실패", func(t *testing.T) {
		mustFail(t, ex, prepare(t, ex), "되돌린 커밋이 없습니다")
	})
	t.Run("revert·push면 통과", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, "git revert HEAD --no-edit && git push origin main")
		mustPass(t, ex, dir)
	})
	t.Run("reset --hard 후 강제 push면 안내", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, "git reset --hard HEAD~1 && git push --force origin main")
		mustFail(t, ex, dir, "이력에 없습니다")
	})
}
