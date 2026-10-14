`git log --oneline`으로 마지막 커밋이 "사과"를 추가한 커밋인지 확인하세요.
---
`git reset`은 브랜치를 앞 커밋으로 되돌려 커밋을 취소합니다. 옵션에 따라 취소한 변경이 남는 곳이 달라집니다: 무대(staged) 위에 그대로 두기, 무대 아래(unstaged)로 내리기(기본값), 다 치우기. 무대에 남기는 옵션을 고르고, 몇 커밋 전으로 갈지도 정해야 합니다.
---
`git log --oneline`에서 "사과" 커밋이 사라지고, `git status`에서 `words.md`가 초록색 "Changes to be committed"에 있어야 합니다.
