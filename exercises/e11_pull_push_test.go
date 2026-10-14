package exercises

import "testing"

func TestPullPush(t *testing.T) {
	ex := ByID(11)
	t.Run("setup 직후에는 실패", func(t *testing.T) {
		mustFail(t, ex, prepare(t, ex), "아직 push하지 않았거나")
	})
	t.Run("pull 후 push하면 통과", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, "git pull --no-rebase --no-edit && git push origin main")
		mustPass(t, ex, dir)
	})
}
