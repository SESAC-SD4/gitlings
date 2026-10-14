package exercises

import "testing"

func TestConflict(t *testing.T) {
	ex := ByID(10)
	t.Run("setup 직후에는 실패", func(t *testing.T) {
		mustFail(t, ex, prepare(t, ex), "병합 커밋이 없습니다")
	})
	t.Run("표범을 남기고 커밋하면 통과", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, `git merge partner || true
printf '# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n3. 표범\n' > words.md
git add words.md && git commit -m "충돌 해결: 표범 선택"`)
		mustPass(t, ex, dir)
	})
	t.Run("충돌 표시를 남긴 채 커밋하면 실패", func(t *testing.T) {
		dir := prepare(t, ex)
		sh(t, dir, `git merge partner || true
printf '# 끝말잇기 단어장\n\n1. 기차\n2. 차표\n<<<<<<< HEAD\n3. 표범\n=======\n3. 사과\n>>>>>>> partner\n' > words.md
git add words.md && git commit -m "충돌 표시 남긴 채 커밋"`)
		mustFail(t, ex, dir, "충돌 표시가 남아 있습니다")
	})
}
