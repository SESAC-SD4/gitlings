## 상황

내가 단어를 추가해 커밋했습니다. 그런데 그사이 "짝"이 `rules.md`를 고쳐 `origin`에 먼저 push했습니다. 이제 `git push origin main`을 하면 거절됩니다.

## 할 일

짝의 커밋을 먼저 받아서 합친 뒤, `origin`에 push하세요.

## 참고

- mac: `git pull`만 실행하면 fatal 에러가 납니다 — 에러 메시지를 읽고 힌트를 활용하세요.
- Windows(Git Bash): `git pull`이 바로 merge를 시도합니다 — 편집기가 열리면 저장 후 닫으세요.

두 경우 모두 통과 조건은 같습니다.
