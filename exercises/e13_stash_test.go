package exercises

import "testing"

func TestStash(t *testing.T) {
	ex := ByID(13)
	t.Run("setup 직후에는 실패", func(t *testing.T) {
		mustFail(t, ex, prepare(t, ex), "오타가 아직 수정되지 않았습니다")
	})
	t.Run("stash → main 수정 → stash pop이면 통과", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, `git stash && git switch main
printf '# 끝말잇기 규칙\n\n다음 단어의 첫 글자는 앞 단어의 끝 글자와 같아야 합니다.\n' > rules.md
git add rules.md && git commit -m "rules.md 오타 수정"
git switch new-word && git stash pop`)
		mustPass(t, ex, dir)
	})
}
