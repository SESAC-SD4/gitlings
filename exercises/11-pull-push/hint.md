push가 거절될 때 나온 메시지를 읽어 보세요. origin에 내 컴퓨터에는 없는 커밋(짝의 커밋)이 있다는 뜻입니다. `git log --oneline --graph --all`로 내 커밋과 origin의 커밋을 비교해 보세요.
---
`git pull`은 origin의 커밋을 받아 지금 브랜치에 합칩니다.
- mac에서 `fatal: Need to specify how to reconcile divergent branches.`가 나면, 그 아래 Git이 보여 주는 설정 세 가지를 읽어 보세요. 그중 **병합(merge) 방식**을 고르는 설정을 한 번 해 두고 다시 pull하면 됩니다.
- Windows(Git Bash)에서 편집기(VSCode 등)가 열리면 그냥 저장 후 닫으면 됩니다.
---
`git log --oneline --graph`에서 내 커밋(범고래)과 짝의 커밋(두음 법칙)이 합쳐져 있어야 합니다. 그다음 다시 push하면 거절되지 않고, `git status`에 origin과 같다는 안내가 보여야 합니다.
