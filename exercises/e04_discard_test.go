package exercises

import "testing"

func TestDiscard(t *testing.T) {
	ex := ByID(4)
	t.Run("setup 직후에는 실패", func(t *testing.T) {
		mustFail(t, ex, prepare(t, ex), "아직 변경된 상태")
	})
	t.Run("git restore면 통과", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, "git restore words.md")
		mustPass(t, ex, dir)
	})
	t.Run("고쳐서 새 커밋하면 안내", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, `git show HEAD:words.md > words.md && echo "4. 래퍼" >> words.md && git add words.md && git commit -m "새 커밋으로 수정"`)
		mustFail(t, ex, dir, "커밋이 새로 생겼습니다")
	})
}
