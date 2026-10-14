`git status`를 쳐 보세요. "not a git repository"가 나오면 아직 이 폴더가 Git 저장소가 아니라는 뜻입니다.
---
`git init`은 지금 폴더를 Git 저장소로 만듭니다. 그다음 `git add`로 사진에 넣을 파일을 무대에 올리고, `git commit`으로 사진을 찍습니다. 커밋에는 `-m` 옵션으로 메시지를 붙입니다.
---
다 되면 `git log --oneline`에 커밋이 하나 보이고, `git status`에 "nothing to commit, working tree clean"이 나와야 합니다.
