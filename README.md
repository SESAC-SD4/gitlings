# gitlings-ko

터미널에서 **진짜 git**으로 배우는 한국어 Git 연습 문제 세트.  
끝말잇기 단어장을 Git으로 관리하면서 13가지 시나리오를 직접 풀어봅니다.

---

## 준비물

- **git 2.23 이상** (switch · restore 명령 지원)
  - mac: `brew install git` 또는 Xcode Command Line Tools
  - Windows: https://git-scm.com/ 에서 설치
- **Python·Go·Node 불필요** — 실행 파일이 이미 포함되어 있습니다.

---

## 빠른 시작

### macOS (zsh)

```sh
git clone https://github.com/...   # 강사가 URL 안내
cd gitlings-ko
./gitlings list
./gitlings start 1
# 안내에 따라 cd ~/gitlings-work/01-first-commit
./gitlings check
```

### Windows (PowerShell)

```powershell
git clone https://github.com/...   # 강사가 URL 안내
cd gitlings-ko
.\gitlings.exe list
.\gitlings.exe start 1
# 안내에 따라 cd ~/gitlings-work/01-first-commit
.\gitlings.exe check
```

> **Defender 경고가 뜨면**: "추가 정보" 클릭 후 "그래도 실행" 선택  
> (서명되지 않은 파일이라 경고가 나옵니다. 강사 확인 파일입니다.)

---

## 명령 정리

| 명령 | 설명 |
|------|------|
| `gitlings list` | 문제 목록 · 통과 여부 |
| `gitlings start <번호>` | N번 문제 시작 |
| `gitlings check` | 현재 문제 채점 |
| `gitlings hint` | 힌트 (3단계, 한 번에 하나씩) |
| `gitlings reset` | 현재 문제 초기화 |
| `gitlings help` | 도움말 |

---

## 문제 목록

| 번호 | 제목 | 배우는 것 |
|------|------|-----------|
| 01 | 첫 사진 | init · add · commit |
| 02 | 골라서 찍기 | 파일 선택 add |
| 03 | 프레임에서 빼기 | restore --staged |
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

문제 폴더는 `~/gitlings-work/` 아래에 만들어집니다.  
진행 상황은 `~/gitlings-work/.progress.json`에 저장됩니다.

---

## 출처 · 라이선스

MIT License.  
구조와 아이디어는 [rustlings](https://github.com/rust-lang/rustlings) (MIT),  
[git-katas](https://github.com/eficode-academy/git-katas) (MIT)를 참고했습니다.
