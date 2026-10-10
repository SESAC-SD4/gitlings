# Windows 확인 목록 (강사용)

gitlings-ko를 Windows에서 처음 실행하기 전 점검 항목.

---

## 1. 실행 파일 구동

```powershell
cd gitlings-ko
.\gitlings.exe list
```

**확인**: 문제 목록이 한글로 출력되는가?

**Defender 경고 처리**:
- "Windows가 PC를 보호했습니다" 경고창 → "추가 정보" → "그래도 실행"
- 또는 PowerShell에서: `Unblock-File .\gitlings.exe`

---

## 2. 콘솔 한글 (UTF-8)

`gitlings.exe`는 시작 시 자동으로 `chcp 65001`을 실행합니다.

**확인**:
- PowerShell 5.1에서 한글 문제 설명이 깨지지 않는가?
- Windows Terminal에서는 기본 UTF-8이라 보통 문제없음
- 구형 cmd.exe / PowerShell 5.1에서 확인 필요

**수동 설정 (필요 시)**:
```powershell
[Console]::OutputEncoding = [Text.Encoding]::UTF8
chcp 65001
```

---

## 3. 경로 · `~` 해석

gitlings는 작업 폴더를 `~/gitlings-work/` 아래에 만듭니다.  
Go의 `os.UserHomeDir()`를 사용하므로 `%USERPROFILE%`을 읽습니다.

**확인**:
```powershell
.\gitlings.exe start 1
# 출력된 경로가 C:\Users\<이름>\gitlings-work\01-first-commit 형식인가?
```

**이동 명령**:
```powershell
cd $HOME\gitlings-work\01-first-commit   # PowerShell 권장
# 또는
cd ~\gitlings-work\01-first-commit       # PowerShell 7+
```

---

## 4. git 버전 확인

```powershell
git --version
```

**확인**: 2.23.0 이상인가?  
미만이면 https://git-scm.com/ 에서 최신 버전 설치.

---

## 5. ex08 / ex11 — 가짜 원격 저장소

08번, 11번 문제는 로컬 bare 저장소를 `origin`으로 사용합니다.  
Windows 경로(`C:\Users\...`)가 포함되므로 확인 필요.

**확인**:
```powershell
.\gitlings.exe start 8
cd $HOME\gitlings-work\08-revert
git remote -v
# origin  C:\Users\...\gitlings-work\08-revert.git-remote (fetch/push) 형식인가?
git log --oneline
.\gitlings.exe check
```

---

## 6. 실행 권한

Windows에서는 별도 실행 권한 설정 불필요.  
파일명이 `.exe`이므로 PowerShell에서 `.\gitlings.exe`로 실행.

---

## 7. 전체 흐름 손으로 확인

```powershell
.\gitlings.exe list
.\gitlings.exe start 1
cd $HOME\gitlings-work\01-first-commit
git init
git config user.name "홍길동"
git config user.email "hong@example.com"
git add words.md
git commit -m "첫 커밋"
.\gitlings.exe check    # 통과 확인
.\gitlings.exe hint     # 힌트 출력 확인
```

---

## 알려진 이슈 / 미확인 사항

- [ ] PowerShell 5.1에서 한글 깨짐 여부
- [ ] Windows Defender가 자동 실행을 차단하는지
- [ ] `cd ~\gitlings-work\...` 축약 경로 동작 여부 (PowerShell 5.1)
- [ ] `git log --oneline` 등에서 한글 커밋 메시지 표시 여부
