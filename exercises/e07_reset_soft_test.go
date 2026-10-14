package exercises

import "testing"

func TestResetSoft(t *testing.T) {
	ex := ByID(7)
	t.Run("setup 직후에는 실패", func(t *testing.T) {
		mustFail(t, ex, prepare(t, ex), "아직 커밋을 취소하지 않았습니다")
	})
	t.Run("reset --soft면 통과", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, "git reset --soft HEAD~1")
		mustPass(t, ex, dir)
	})
	t.Run("--hard면 안내", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, "git reset --hard HEAD~1")
		mustFail(t, ex, dir, "--hard")
	})
	t.Run("--mixed면 안내", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, "git reset HEAD~1")
		mustFail(t, ex, dir, "--mixed")
	})
}
