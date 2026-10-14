# Windows 확인 목록 (강사용)

학생은 **Git Bash**(Git for Windows에 포함된 mintty 터미널)를 사용합니다.  
명령 형식은 mac과 동일합니다 (`./gitlings.exe`).

---

## 0. 사전 확인

- Git for Windows가 설치되어 있는가? → 시작 메뉴에서 "Git Bash" 검색
- git 버전 확인 (Git Bash에서):
  ```sh
  git --version
  # git version 2.23 이상이어야 함
  ```

---

## 1. 실행 파일 구동

Git Bash에서:
```sh
cd gitlings-ko
./gitlings.exe list
```

**확인**: 문제 목록이 한글로 출력되는가?

**Defender 경고 처리**:
- 경고창 → "추가 정보" → "그래도 실행"
- 또는 PowerShell에서: `Unblock-File .\gitlings.exe`

---

## 2. 한글 출력 (Git Bash)

Git Bash(mintty)는 UTF-8이 기본이므로 별도 설정 불필요.

**확인**:
```sh
./gitlings.exe start 1
# 문제 설명이 한글로 올바르게 표시되는가?
```

PowerShell이나 cmd에서 직접 실행하면 한글이 깨질 수 있음.  
→ 학생에게 **Git Bash를 사용하도록** 안내.

---

## 3. 경로 · `~` 해석

Git Bash에서 `~`는 `/c/Users/<이름>` (Unix 스타일)으로 해석됩니다.  
gitlings는 `$HOME`을 먼저 읽고, 없을 때만 `os.UserHomeDir()`(`%USERPROFILE%`)을 씁니다.

**확인**:
```sh
./gitlings.exe start 1
# 출력 예: cd ~/gitlings-work/01-first-commit
# 실제 경로: C:\Users\<이름>\gitlings-work\01-first-commit
```

**이동 명령** (Git Bash):
```sh
cd ~/gitlings-work/01-first-commit
```

---

## 4. 분할 화면(TUI) 전체 흐름

터미널 두 개 열기.

**터미널 1 — 실행:**
```sh
cd gitlings-ko
./gitlings.exe
```

**터미널 2 — git 연습:**
```sh
cd ~/gitlings-work/01-first-commit
git init
git add words.md
git commit -m "첫 커밋"
```
터미널 1에서 **Enter** → ✅ 확인

- [ ] 켜자마자 `~/gitlings-work/` 안에 13개 연습 폴더(08·11번은 `.git-remote`도)가 생긴다. 처음 켤 때 몇 초 걸리는지 적는다
- [ ] 화면이 그려지고 `↑` `↓` Enter Tab Esc(힌트 닫기) Ctrl+R Ctrl+X Ctrl+C가 모두 먹는다
- [ ] **한글 입력 상태**에서도 Enter·Tab·Ctrl+R·Ctrl+X·화살표가 먹는다
- [ ] 한글 제목·설명이 깨지지 않고, 왼쪽 목록과 오른쪽 칸의 테두리가 어긋나지 않는다
- [ ] 시작할 때 이상한 문자(`]11;rgb:...` 같은 것)가 찍히거나 몇 초 멈추지 않는다(배경색 감지)
- [ ] 확인 질문에서 Enter는 취소되고, 영문 `y`로만 되돌려진다
- [ ] 다른 탭이 `~/gitlings-work/04-discard` **안에 있는 채로** Ctrl+R → `y`: 오류 없이 되돌려지고, 그 탭에서 `git status`가 바로 된다

---

## 5. ex08 / ex11 — 가짜 원격 저장소

08번, 11번 문제는 로컬 bare 저장소를 `origin`으로 사용합니다.  
Windows 경로가 포함되므로 Git Bash에서 확인 필요.

```sh
./gitlings.exe start 8
# 다른 터미널에서:
cd ~/gitlings-work/08-revert
git remote -v
# origin  /c/Users/.../gitlings-work/08-revert.git-remote (fetch/push) 형식인가?
git log --oneline
```

---

## 6. PowerShell에서의 동작 (보조 확인)

학생 주 환경은 Git Bash이지만 깨지지 않는지 확인:
```powershell
cd gitlings-ko
.\gitlings.exe list
# 한글이 어느 정도 표시되는가? (완벽하지 않아도 됨)
```

---

## 알려진 이슈 / 미확인 사항

- [ ] Git Bash에서 색상(굵기, 색)이 올바르게 표시되는가?
- [ ] Windows Defender가 자동 실행을 차단하는지
- [ ] `git log --oneline` 등에서 한글 커밋 메시지 표시 여부
