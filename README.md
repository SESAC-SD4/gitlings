# gitlings-ko

터미널에서 **진짜 git**으로 배우는 한국어 Git 연습 문제 세트.  
끝말잇기 단어장을 Git으로 관리하면서 13가지 시나리오를 직접 풀어봅니다.

---

## 준비물

- **git 2.23 이상** (switch · restore 명령 지원)
  - mac: `brew install git` 또는 Xcode Command Line Tools
  - Windows: https://git-scm.com/ 에서 설치 — **Git Bash 포함**
- **Python·Go·Node 불필요** — 실행 파일이 이미 포함되어 있습니다.

---

## 준비: 터미널 두 개

| | 터미널 1 — 문제 화면 | 터미널 2 — git 명령 치는 곳 |
| --- | --- | --- |
| Windows | Windows Terminal에서 **Git Bash** 탭 | 같은 창에 **Git Bash** 탭 하나 더 (`Ctrl+Shift+T` 또는 탭 옆 `∨`에서 Git Bash) |
| mac | iTerm 또는 터미널 창 | 새 탭(`Cmd+T`) |

> Windows에서는 PowerShell 탭이 아니라 **Git Bash** 탭이어야 합니다. 프롬프트에 `MINGW64`가 보이면 맞습니다.  
> **Defender 경고가 뜨면**: "추가 정보" 클릭 후 "그래도 실행" 선택

## 시작하기

터미널 1에서:

```bash
cd ~
git clone <강사가 알려 준 주소> gitlings
cd gitlings
./gitlings.exe     # Windows
./gitlings         # mac
```

켜면 연습 폴더 13개가 `~/gitlings-work/` 안에 만들어지고, 왼쪽에 문제 목록, 오른쪽에 문제 설명이 보입니다. 이미 만들어진 연습 폴더는 다시 켜도 그대로 남습니다.

| 키 | 하는 일 |
| --- | --- |
| `↑` `↓` | 문제 고르기 |
| **Enter** | 채점 |
| **Tab** | 힌트 (누를 때마다 한 단계씩, 3단계까지) · **Esc**로 힌트 닫기 |
| **Ctrl+R** | 지금 문제만 처음 상태로 (✅ 표시도 지워짐) |
| **Ctrl+X** | 13문제 전체 초기화 |
| **Ctrl+C** | 끝내기 |
| 스페이스 / `b` | 문제 설명 스크롤 |

> 위 키들은 **한글 입력 상태에서도** 동작합니다. 되돌리기·전체 초기화는 확인 질문에 영문 **`y`**를 눌러야만 실행됩니다(Enter는 '아니오').

## 한 문제 푸는 흐름

1. 터미널 1: 문제를 고른다. 오른쪽 아래에 `cd ~/gitlings-work/01-first-commit` 같은 명령이 보입니다
2. 터미널 2: 그 명령으로 연습 폴더에 들어가서 git 명령으로 문제를 푼다
3. 터미널 1: **Enter** — 통과하면 ✅, 아니면 무엇이 아직인지 알려 줍니다

## 화면이 안 뜰 때 (명령으로 풀기)

```bash
./gitlings.exe list       # 목록 (mac은 ./gitlings)
./gitlings.exe start 1    # 1번 연습 폴더 만들기 + 문제 보기
./gitlings.exe check      # 채점
./gitlings.exe hint       # 힌트
./gitlings.exe reset      # 처음 상태로
```

---

## 문제 목록

| 번호 | 제목 | 배우는 것 |
|------|------|-----------|
| 01 | 첫 사진 | init · add · commit |
| 02 | 골라서 찍기 | 파일 선택 add |
| 03 | 무대에서 내리기 | restore --staged |
| 04 | 망친 수정 버리기 | restore |
| 05 | 메시지 오타 | commit --amend |
| 06 | 비밀이 섞였다 | .gitignore · rm --cached |
| 07 | 방금 커밋 취소 | reset --soft |
| 08 | 이미 올린 커밋 | revert · push |
| 09 | 갈래 만들고 합치기 | branch · switch · merge |
| 10 | 충돌 | 충돌 해결 |
| 11 | push가 거절됐다 | pull --no-rebase · push |
| 12 ★ | 날린 커밋 되살리기 | reflog |
| 13 ★ | 하던 일 치워두기 | stash |

★ = 선택 문제

---

## 작업 폴더

문제 폴더는 `~/gitlings-work/` 아래에 만들어집니다. 원격 저장소가 필요한 문제는 옆에 `<문제 폴더>.git-remote`도 만듭니다.  
진행 상황은 `~/gitlings-work/.progress.json`에 저장됩니다.

---

## 출처 · 라이선스

MIT License.  
분할 화면 구성은 [dockerlings](https://github.com/furkan/dockerlings)(MIT)를, 구조와 아이디어는 [rustlings](https://github.com/rust-lang/rustlings) (MIT),  
[git-katas](https://github.com/eficode-academy/git-katas) (MIT)를 참고했습니다. 화면과 명령 체계는 같은 수업의 `shellings-ko`와 같습니다.
