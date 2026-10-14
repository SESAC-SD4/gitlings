package exercises

import "testing"

func TestPickFiles(t *testing.T) {
	ex := ByID(2)
	t.Run("setup 직후에는 실패", func(t *testing.T) {
		mustFail(t, ex, prepare(t, ex), "새 커밋이 없습니다")
	})
	t.Run("words.md만 add·commit이면 통과", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, `git add words.md && git commit -m "단어 '표범' 추가"`)
		mustPass(t, ex, dir)
	})
	t.Run("git add .로 둘 다 커밋하면 안내", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, `git add . && git commit -m "잘못된 커밋"`)
		mustFail(t, ex, dir, "낙서")
	})
}
