package exercises

import "testing"

func TestFirstCommit(t *testing.T) {
	ex := ByID(1)
	t.Run("setup 직후에는 실패", func(t *testing.T) {
		mustFail(t, ex, prepare(t, ex), "Git 저장소가 아닙니다")
	})
	t.Run("init·add·commit이면 통과", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, `git init -b main && git add words.md && git commit -m "첫 번째 커밋"`)
		mustPass(t, ex, dir)
	})
	t.Run("init만 하면 커밋 안내", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, "git init -b main")
		mustFail(t, ex, dir, "커밋이 없습니다")
	})
	t.Run("add만 하면 commit 안내", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, "git init -b main && git add words.md")
		mustFail(t, ex, dir, "아직 커밋하지 않았습니다")
	})
}
