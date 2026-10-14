`git status`와 `git ls-files`로 지금 Git이 어떤 파일을 추적하고 있는지 먼저 보세요. `.env`가 목록에 있습니다.
---
`.gitignore`는 새 파일만 무시합니다. 이미 추적 중인 파일은 `git rm`으로 추적에서 빼야 하는데, 그냥 쓰면 디스크의 파일까지 지워집니다. 파일은 남기고 Git의 기록(캐시)에서만 빼는 옵션을 찾아보세요.
---
커밋한 뒤 `ls -a`에는 `.env`가 보이고, `git ls-files`에는 `.env`가 없고 `.gitignore`가 있어야 합니다. `git status`는 깨끗해야 합니다.
