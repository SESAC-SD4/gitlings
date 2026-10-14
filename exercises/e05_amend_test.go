package exercises

import "testing"

func TestAmend(t *testing.T) {
	ex := ByID(5)
	t.Run("setup 직후에는 실패", func(t *testing.T) {
		mustFail(t, ex, prepare(t, ex), "오타가 남아 있습니다")
	})
	t.Run("commit --amend면 통과", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, `git commit --amend -m "단어 '표범' 추가"`)
		mustPass(t, ex, dir)
	})
	t.Run("새 커밋을 추가하면 안내", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, `git commit --allow-empty -m "단어 '표범' 추가"`)
		mustFail(t, ex, dir, "커밋이 하나 늘었습니다")
	})
}
