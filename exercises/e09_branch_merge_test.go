package exercises

import "testing"

func TestBranchMerge(t *testing.T) {
	ex := ByID(9)
	t.Run("setup 직후에는 실패", func(t *testing.T) {
		mustFail(t, ex, prepare(t, ex), "브랜치가 없습니다")
	})
	t.Run("브랜치에서 단어 추가 후 병합하면 통과", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, `git switch -c new-word && echo "3. 표범" >> words.md && git add words.md && git commit -m "단어 '표범' 추가" && git switch main && git merge new-word`)
		mustPass(t, ex, dir)
	})
	t.Run("규칙에 맞지 않는 단어면 안내", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, `git switch -c new-word && echo "3. 사과" >> words.md && git add words.md && git commit -m "규칙 위반 단어 추가" && git switch main && git merge new-word`)
		mustFail(t, ex, dir, "규칙 위반")
	})
}
