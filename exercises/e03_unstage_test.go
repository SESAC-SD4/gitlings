package exercises

import "testing"

func TestUnstage(t *testing.T) {
	ex := ByID(3)
	t.Run("setup 직후에는 실패", func(t *testing.T) {
		mustFail(t, ex, prepare(t, ex), "아직 커밋하지 않았습니다")
	})
	t.Run("restore --staged 후 commit이면 통과", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, `git restore --staged notes.md && git commit -m "단어 '범고래' 추가"`)
		mustPass(t, ex, dir)
	})
	t.Run("notes.md를 지우면 안내", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, `rm notes.md && git add . && git commit -m "notes.md 삭제함"`)
		mustFail(t, ex, dir, "삭제됐습니다")
	})
}
