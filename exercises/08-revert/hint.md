`git log --oneline`으로 잘못된 커밋이 맨 위에 있는지 보세요. 이 커밋은 이미 `origin`에도 올라가 있습니다.
---
`git revert`는 지정한 커밋을 지우지 않고, 그 반대 변경을 담은 새 커밋을 만듭니다. 되돌릴 커밋을 해시나 `HEAD`로 지정합니다. 커밋 메시지 편집기가 열리면 저장하고 닫으면 됩니다. 그다음 `git push origin main`으로 올립니다.
---
`git log --oneline`에 "사과" 커밋과 그것을 되돌린 "Revert" 커밋이 모두 보이고, `words.md`에는 "사과"가 없어야 합니다. `git status`가 origin과 같다고(up to date) 나와야 합니다.
