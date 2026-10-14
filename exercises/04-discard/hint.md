`git diff`로 마지막 커밋과 비교해 `words.md`가 어떻게 망가졌는지 보세요. `git status`가 변경을 버리는 방법도 알려 줍니다.
---
`git restore`는 작업 트리의 파일을 마지막 커밋 상태로 되돌립니다. 되돌린 변경은 다시 살릴 수 없으니 버려도 되는 변경인지 먼저 확인하세요.
---
다 되면 `git status`에 "nothing to commit, working tree clean"이 나오고, `git log --oneline`의 커밋 수는 그대로여야 합니다.
