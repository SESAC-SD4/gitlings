# gitlings-ko 사양 (v1, 2026-10-11)

> 이 문서는 제작 세션에 넘기는 사양입니다. 결정 배경은 강사 저장소(private) `lecture/00_curriculum/02_기획/교안_설계_2026-09.md`의 **D26**.
> 모르는 것은 추측하지 말고 `## 열린 질문`에 적어 강사에게 묻습니다.

## 1. 무엇을 만드나

**터미널에서 진짜 git으로 풀고 자동 채점되는 한국어 Git 연습 문제 세트.** rustlings·dockerlings 방식.

- 대상: 한국 부트캠프 비전공 성인 30명. Windows(PowerShell)·mac(zsh) 혼재
- 쓰는 날: 10/14(수) Git 이틀차 마지막 정리 시간에 13문제를 한꺼번에 푼다
- 학생이 이미 배운 것: 터미널 기초, `init·add·commit·status·log·diff`, `restore`·`restore --staged`·`commit --amend`, `.gitignore`·`rm --cached`, `remote·push`(gh 인증), 그리고 10/14 당일 `HEAD~`·`reset`·`revert`·`branch·switch·merge`·충돌·PR
- 학생이 모르는 것: Python·Go 등 프로그래밍. **도구 내부를 몰라도 쓸 수 있어야 한다**

## 2. 반드시 지킬 조건

1. **`git clone` 후 바로 실행** — 런타임 설치 없음
   - Go(1.25)로 빌드한 실행 파일을 저장소에 커밋한다
   - mac: arm64+amd64를 `lipo`로 합친 유니버설 바이너리 `gitlings` (실행 권한 비트 유지 확인)
   - Windows: `gitlings.exe` (amd64)
   - 학생 명령: mac `./gitlings list` / Windows `.\gitlings.exe list`
2. **진짜 git** — 시스템에 설치된 `git`을 호출한다. 시작할 때 `git --version`이 **2.23 이상**(switch·restore)인지 확인하고, 아니면 한국어로 안내
3. **네트워크·GitHub 없음** — 원격이 필요한 문제는 로컬 bare 저장소를 가짜 원격(`origin`)으로 만든다
4. **한국어** — 모든 안내·문제·힌트·채점 결과. 단 **폴더·파일 이름은 ASCII**(Windows 인코딩 사고 방지). 파일 *내용*은 한국어
5. **Windows 콘솔 한글** — Windows에서는 시작 시 콘솔 출력 코드페이지를 UTF-8(65001)로 맞춘다(PowerShell 5.1에서 깨지지 않게). 실제 확인은 강사가 Windows에서 한다
6. **연습 저장소는 gitlings 폴더 밖에 만든다** — 기본 `~/gitlings-work/NN-이름/`. 이유: 1번 문제(아직 저장소 아님)에서 학생이 `git status`를 치면 상위의 gitlings-ko 저장소를 잡아버린다
7. **학생의 전역 git 설정에 의존하지 않는다** — 문제를 만들 때 필요한 커밋 작성자는 `-c user.name/-c user.email`로 넣는다. 학생에게 `user.name`이 없을 수도 있다(그 경우 안내 메시지)
8. 기본 브랜치는 `main`. setup은 `init.defaultBranch` 설정과 무관하게 `main`을 만든다

## 3. 명령

```
gitlings            = gitlings list
gitlings list       문제 목록 · 통과 여부(✅/⬜) · ★는 선택
gitlings start N    N번 문제의 상황을 ~/gitlings-work/NN-이름/ 에 만든다. 이미 있으면 덮어쓸지 묻는다
                    → 문제 설명(무엇이 어떤 상태이고, 무엇을 하면 되는지)을 출력하고
                       "cd ~/gitlings-work/NN-이름" 안내
gitlings check      현재 폴더(또는 가장 최근 start한 문제)를 채점. 통과면 ✅ + 다음 문제 안내,
                    실패면 "무엇이 아직 아닌지"를 한국어로. 흔한 오답은 따로 알아보고 짚는다
gitlings hint       힌트 한 단계씩 (문제당 3단계: 방향 → 명령 이름 → 거의 정답)
gitlings reset      현재 문제를 처음 상태로 다시 만든다
gitlings help
```

- 진행 상황은 `~/gitlings-work/.progress.json`에 저장
- `check`는 **현재 폴더가 어느 문제인지** 알아야 한다 → 문제 폴더에 표식 파일을 두되, **git이 추적하지 않게** `.git/info/exclude`에 넣는다(학생의 `git status`에 보이지 않도록)

## 4. 소재 — 끝말잇기 단어장

13문제가 하나의 이야기로 이어진다: **친구들과 끝말잇기 단어장을 Git으로 관리한다.**

- `words.md` — 단어장. 한 줄에 `N. 단어` 형식
  ```
  # 끝말잇기 단어장
  1. 기차
  2. 차표
  ```
- 규칙: 다음 단어의 첫 글자 = 앞 단어의 끝 글자, 두 글자 이상. **채점에서 규칙 검사 함수 하나를 공용으로 쓴다**
- `rules.md` — 규칙 설명 (11·13번에서 "짝이 고친 다른 파일")
- `notes.md` — 낙서·메모 (2·3번의 "커밋하면 안 되는 파일")
- `.env` — `ADMIN_PASSWORD=...` (6번)
- 10/14 PR 실습(끝말잇기)과 같은 소재라서, 혼자 gitlings로 연습한 것을 짝과 GitHub에서 다시 하게 된다

## 5. 문제 목록 (정규 11 + 선택 2)

각 문제: **상황(setup) · 해야 할 일 · 통과 조건(check) · 잡아낼 흔한 오답 · 힌트 3단계**.
통과 조건은 **명령이 아니라 저장소 상태**로 본다(같은 결과를 내는 다른 방법도 통과).

### 10/13 범위 — 혼자 기록하기

**01 첫 사진** (`01-first-commit`)
- 상황: 저장소가 아닌 폴더. `words.md`(기차·차표)만 있다
- 할 일: 저장소로 만들고 `words.md`를 첫 커밋
- 통과: `.git` 있음 · 커밋 ≥1 · `words.md`가 HEAD에 있음 · 작업 트리 깨끗
- 흔한 오답: init만 하고 커밋 안 함 → "사진을 아직 안 찍었다" / add 없이 commit → git 메시지를 읽게 안내

**02 골라서 찍기** (`02-pick-files`)
- 상황: 커밋 1개. `words.md`에 3번째 단어(표범) 추가됨, `notes.md`에 낙서 추가됨 — 둘 다 수정 상태
- 할 일: `words.md`만 커밋, `notes.md`는 수정 상태로 남긴다
- 통과: 새 커밋이 `words.md`만 바꿈 · `notes.md`는 수정됨(unstaged) · 단어장 규칙 유효
- 흔한 오답: `git add .`로 둘 다 커밋 → "낙서까지 사진에 찍혔다" + 되돌리는 법은 hint

**03 프레임에서 빼기** (`03-unstage`)
- 상황: `words.md` 수정 + 새 파일 `notes.md` 둘 다 add 된 상태(staged)
- 할 일: `notes.md`는 프레임에서 빼고 `words.md`만 커밋. `notes.md` 파일은 지우지 않는다
- 통과: 마지막 커밋에 `notes.md` 없음 · `notes.md`가 디스크에 있고 untracked · 커밋에 `words.md` 변경 있음
- 흔한 오답: `notes.md`를 삭제 → "파일까지 지웠다"

**04 망친 수정 버리기** (`04-discard`)
- 상황: 커밋된 `words.md`를 실수로 엉망으로 고침(저장만 함, add 안 함)
- 할 일: 마지막 커밋 상태로 되돌린다
- 통과: `words.md` == HEAD 내용 · 커밋 수 그대로 · 작업 트리 깨끗

**05 메시지 오타** (`05-amend`)
- 상황: 마지막 커밋 메시지가 `단어 추가함ㅁㅁㄴㅇ`
- 할 일: 메시지만 `단어 '표범' 추가`로 고친다 (내용은 그대로)
- 통과: HEAD 메시지 == `단어 '표범' 추가`(앞뒤 공백 무시) · 커밋 수가 **늘지 않음** · HEAD 트리가 setup 때와 같음
- 흔한 오답: 새 커밋을 하나 더 만듦 → "커밋이 하나 늘었다. 고치는 게 아니라 새로 찍었다"

**06 비밀이 섞였다** (`06-gitignore`)
- 상황: `.env`(ADMIN_PASSWORD=...)가 이미 커밋되어 추적 중
- 할 일: `.env`를 `.gitignore`에 넣고, 추적에서 빼서 커밋. **파일은 디스크에 남긴다**
- 통과: HEAD 트리에 `.env` 없음 · `.gitignore`가 HEAD에 있고 `.env`가 무시됨(`git check-ignore`) · `.env` 파일 존재
- 흔한 오답: `.env`를 지움 → 실패 / `.gitignore`만 쓰고 rm --cached 안 함 → "아직 추적 중이다"
- 문제 설명 끝에 한 줄: "이전 커밋에는 비밀번호가 여전히 남아 있습니다. 진짜라면 비밀번호를 바꿔야 합니다"

### 10/14 범위 — 이력 다루기와 함께 합치기

**07 방금 커밋 취소** (`07-reset-soft`)
- 상황: 마지막 커밋이 규칙에 안 맞는 단어(`사과` — "표"로 시작해야 함)를 넣었다
- 할 일: 커밋만 취소하고 변경은 프레임(staged)에 남긴다
- 통과: HEAD == setup의 HEAD~1 · 인덱스에 `사과` 추가 변경이 있음
- 흔한 오답: `--hard` → "변경이 사라졌다. 무엇을 썼나?" / `--mixed`(기본) → "변경은 남았지만 프레임 밖이다"

**08 이미 올린 커밋** (`08-revert`)
- 상황: 가짜 원격 `origin`에 잘못된 단어 커밋까지 push된 상태
- 할 일: 이력을 지우지 않고 되돌린 뒤 push
- 통과: 원래 잘못된 커밋이 이력에 남아 있음 · 그 뒤에 되돌린 커밋 · `words.md`가 잘못되기 전과 같음 · 로컬 main == origin/main
- 흔한 오답: reset 후 강제 push → "원격 이력을 지웠다. 짝이 이미 받았다면?" (실패 처리)

**09 갈래 만들고 합치기** (`09-branch-merge`)
- 상황: main에 단어 2개
- 할 일: `new-word` 브랜치를 만들어 규칙에 맞는 단어를 커밋하고, main에 병합
- 통과: `new-word` 브랜치 존재 · 그 브랜치의 커밋이 main에 포함 · 단어장 규칙 유효

**10 충돌** (`10-conflict`)
- 상황: main과 `partner` 브랜치가 3번 줄에 각각 다른 단어 — main은 `표범`(규칙 맞음), partner는 `사과`(규칙 틀림)
- 할 일: `partner`를 main에 병합 → 충돌 → **끝말이 이어지는 단어를 남겨** 해결하고 커밋
- 통과: 충돌 표시(`<<<<<<<`·`=======`·`>>>>>>>`) 없음 · 병합 커밋(부모 2개) · 단어장 규칙 유효 · 3번 단어가 `표범`
- 흔한 오답: 표시를 남긴 채 커밋 → 줄 번호까지 짚어 준다

**11 push가 거절됐다** (`11-pull-push`)
- 상황: 내가 단어를 추가해 커밋. 그사이 "짝"이 `rules.md`를 고쳐 origin에 먼저 push → 내 push가 거절된다
- 할 일: 짝의 커밋을 받아 합친 뒤 push
- 통과: 로컬 main == origin/main · 내 커밋과 짝의 커밋 둘 다 포함 (merge든 rebase든 통과)
- ⚠️ **확인된 사실 (git 2.50, 2026-10-11)**: pull 방식 설정이 없으면 갈라진 이력에서 `git pull`이 `fatal: Need to specify how to reconcile divergent branches.`로 실패한다. setup에서 이 설정을 **넣지 않는다** — 학생이 이 메시지를 만나는 것이 수업 내용이다. 힌트 2단계에서 `git pull --no-rebase`를 알려 준다

### 선택 ★

**12★ 날린 커밋 되살리기** (`12-reflog`)
- 상황: 커밋 3개 중 2개를 `reset --hard HEAD~2`로 날린 상태(setup이 실행)
- 할 일: 날린 커밋을 되살린다
- 통과: main이 setup 때 날리기 전의 커밋을 가리킴

**13★ 하던 일 치워두기** (`13-stash`)
- 상황: `new-word` 브랜치에서 `words.md` 수정 중(커밋 안 함). main의 `rules.md`에 급한 오타가 있다
- 할 일: 하던 일을 치워두고 main으로 가서 오타를 고쳐 커밋 → 다시 돌아와 하던 일을 되살린다
- 통과: main에 오타 수정 커밋 · `new-word`의 작업 트리에 하던 수정이 있음(또는 그 수정이 커밋됨) · stash 목록이 비어 있음

## 6. 구조 · 품질

- 문제 정의(설명·힌트·흔한 오답 메시지)는 문제별 파일로 분리하고 `go:embed`로 넣는다 → 문장만 고칠 때 찾기 쉽게
- **자동 테스트**: 각 문제마다 (a) setup 직후 check가 **실패**하고 (b) 모범 풀이 명령을 실행하면 **통과**하고 (c) 흔한 오답을 재현하면 해당 메시지가 나오는지 Go 테스트로 검증
- 빌드 스크립트: 3종 빌드 + mac 유니버설 합치기를 한 번에 (`make` 또는 `build.sh`)
- 라이선스: MIT. Githug(MIT)·git-katas(MIT)에서 구조·아이디어를 참고했다면 README에 출처 표기. **라이선스 파일이 없는 저장소(Git-Gud·gitsim 등)의 코드는 가져오지 않는다**
- README.md(학생용, 한국어): clone → 실행 → `list` → `start 1` → `check`까지 Windows·mac 따로 화면 예시. 학생이 처음 보는 문서이므로 짧고 친절하게
- 강사용 모범 풀이는 이 저장소에 두지 않는다(공개 예정). 테스트 코드의 풀이 명령은 괜찮다

## 7. 하지 않는 것

- GitHub·네트워크 사용 · AI 힌트 · 웹 UI · 진행 상황 서버 전송
- `rebase`·`cherry-pick`·`bisect` 문제 (Learn Git Branching 심화로 넘김)

## 8. 완료 기준

- [ ] mac(이 Mac)에서 13문제 모두 테스트 통과 · 실제 터미널에서 `start → 풀이 → check` 손으로 1회 이상
- [ ] `gitlings` · `gitlings.exe` 빌드되어 커밋됨 · 새 폴더에 clone해서 바로 실행되는지 확인
- [ ] Windows 실행은 **강사가 확인** — 확인할 목록을 `WINDOWS_CHECK.md`에 적어 둔다 (콘솔 한글, 경로, `~` 해석, 실행 권한·Defender 반응)
- [ ] GitHub 원격 저장소 생성·공개는 **하지 않는다** (강사 승인 후 별도)

## 열린 질문 (제작 중 생기면 여기에 추가)

- 문제 수가 늘 수 있게 번호 체계를 어떻게 둘지

---

## 구현 중 결정사항 (강사 확인 후 반영)

1. **번호 체계**: 13개 이후 추가 시 14, 15... 순번을 이어감. 슬러그(`14-cherry-pick`)가 실질 식별자. 중간 삽입 시 기존 슬러그를 유지하고 번호만 조정. → 특별한 이슈 없으면 그대로 진행.

2. **표식 파일 방식**: SPEC §3에서 `.gitlings_id` 파일을 `.git/info/exclude`에 등록하도록 명시함. 구현에서는 이 방식을 채택. 단 ex01처럼 `.git`이 없는 문제에서는 `check` 실행 시 `.git`이 생겼을 때 자동으로 exclude에 추가하도록 처리함.

3. **check의 workDir 감지 우선순위**: (1) 현재 경로가 `~/gitlings-work/<slug>/` 패턴에 맞으면 슬러그로 식별 → (2) `.gitlings_id` 파일 → (3) `last_started` fallback.

4. **ex07 setup이 --hard 상태**: setup 자체가 규칙 위반 커밋을 만든 상태이므로, `check`를 실행하면 "커밋이 있고 staged도 없음 → --hard처럼 보임"이라는 메시지가 뜸. 학생이 setup 직후 check를 실행하면 "--hard 옵션을 쓴 것 같습니다"라는 안내가 나오는데, 이는 혼란스러울 수 있음. → 문제 설명에서 먼저 `gitlings check`를 실행하지 말고 `git reset --soft HEAD~1`부터 시도하도록 안내 추가를 고려.

5. **ex11 git pull 동작**: git 2.50에서 확인된 것처럼 갈라진 이력에서 `git pull` 시 `fatal: Need to specify how to reconcile divergent branches.` 에러가 발생함. 모범 풀이는 `git pull --no-rebase`. rebase 방식도 허용 (로컬 main == origin/main 조건만 확인).
