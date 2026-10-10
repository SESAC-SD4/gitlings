# gitlings-ko 구현 계획 (PLAN v1)

> SPEC.md v1(2026-10-11)을 읽고 작성한 구현 계획.

---

## 이해 요약

**목표**: Go로 빌드한 CLI 바이너리 `gitlings` / `gitlings.exe`. 학생이 clone 후 바로 실행, 한국어로 진짜 git을 연습. 끝말잇기 단어장 소재 13문제.

**핵심 제약**
- 런타임 설치 없음 (바이너리 커밋)
- 시스템 git 2.23+ 필수 (시작 시 확인)
- 로컬 bare 저장소로 가짜 origin (네트워크 없음)
- 작업 폴더 `~/gitlings-work/` (gitlings-ko 저장소 밖)
- 전역 git 설정 불의존 (`-c user.name/-c user.email` 주입)
- Windows PowerShell 5.1 UTF-8 (chcp 65001)

---

## 패키지 구조

```
gitlings-ko/
├── cmd/gitlings/
│   └── main.go               # 진입점, 서브커맨드 라우팅
├── internal/
│   ├── console/
│   │   └── console.go        # Windows UTF-8 초기화
│   ├── git/
│   │   └── git.go            # git 명령 실행기, 버전 확인
│   ├── progress/
│   │   └── progress.go       # ~/gitlings-work/.progress.json
│   └── exercise/
│       ├── exercise.go       # Exercise 구조체 및 Result 타입
│       └── wordchain.go      # 끝말잇기 규칙 검증 공용 함수
├── exercises/
│   ├── registry.go           # 전체 문제 목록 등록
│   ├── shared/
│   │   └── words_initial.md  # 초기 단어장 (embed)
│   ├── ex01/
│   │   ├── exercise.go       # Setup + Check Go 함수
│   │   ├── desc.txt          # 문제 설명 (go:embed)
│   │   └── hints.txt         # 힌트 3단계 (go:embed, 빈 줄 구분)
│   ├── ex02/ ... ex13/
├── gitlings                  # mac 유니버설 바이너리 (커밋됨)
├── gitlings.exe              # Windows amd64 바이너리 (커밋됨)
├── Makefile                  # make build, make test, make clean
├── build.sh                  # 3종 빌드 + lipo 유니버설 합치기
├── README.md                 # 학생용 한국어, Windows/mac 예시
└── WINDOWS_CHECK.md          # 강사용 Windows 확인 목록
```

---

## 문제 정의 형식

각 `exercises/exNN/exercise.go`는 아래 패턴을 따른다.

```go
package exNN

import (
    _ "embed"
    "gitlings-ko/internal/exercise"
    "gitlings-ko/internal/git"
)

//go:embed desc.txt
var desc string

//go:embed hints.txt
var hintsRaw string

func New() *exercise.Exercise {
    return &exercise.Exercise{
        ID:       NN,
        Slug:     "NN-slug",
        Name:     "NN 제목",
        Optional: false,
        Desc:     desc,
        Hints:    exercise.ParseHints(hintsRaw),
        Setup:    setup,
        Check:    check,
    }
}

func setup(workDir string, g *git.Git) error { ... }
func check(workDir string, g *git.Git) (exercise.Result, error) { ... }
```

### `desc.txt` 형식
```
상황: ...
할 일: ...
(선택) 참고: ...
```

### `hints.txt` 형식 — 빈 줄로 3단계 구분
```
(1단계: 방향)

(2단계: 명령 이름)

(3단계: 거의 정답)
```

---

## 진행 저장

- 파일: `~/gitlings-work/.progress.json`
- 내용: `{"passed": [1,3], "last_started": 3}`
- 문제 식별자: 문제 폴더 `.git/info/exclude`에도 등록된 `.gitlings_id` 파일 (내용: 문제 번호)
  - `check`는 현재 디렉토리에서 이 파일을 찾아 문제 번호를 확인
  - 현재 폴더가 문제 폴더가 아니면 `last_started` 번호로 fallback

---

## 테스트 방식

각 `exercises/exNN/exercise_test.go`:

```go
func TestSetupThenFail(t *testing.T)    // setup 직후 → check 실패
func TestSolveThenPass(t *testing.T)    // 모범 풀이 후 → check 통과
func TestCommonMistake(t *testing.T)    // 흔한 오답 → 특정 메시지 포함 확인
```

테스트마다 `t.TempDir()`를 workDir로 사용하여 격리.

---

## 구현 순서

1. `PLAN.md` 커밋
2. Go 모듈 init + 공용 뼈대(git, progress, console, exercise) + 문제 없이 `list/start/check/hint/reset` 동작
3. ex01 구현 + 테스트 통과
4. ex02~ex11 각각 구현 + 테스트 통과 후 커밋
5. ex12★, ex13★ 구현 + 테스트
6. build.sh / Makefile + 바이너리 빌드 후 커밋
7. README.md, WINDOWS_CHECK.md 작성 후 커밋

---

## 열린 질문

> SPEC.md `## 열린 질문`에 추가한 항목

1. **번호 체계 확장**: 13개 이후 문제 추가 시 번호를 이어감 (14, 15, ...). 슬러그(예: `14-cherry-pick`)가 실질 식별자라 중간 삽입 시 번호 재정렬 없이 슬러그만 수정. Optional 여부는 `Optional: true` 필드로 관리.
2. **check의 workDir 감지 우선순위**: (1) 현재 폴더의 `.gitlings_id` → (2) `last_started` fallback → (3) 오류 안내. SPEC §3에서 "현재 폴더 또는 가장 최근 start한 문제"라고 명시되어 있어 이 순서가 적합하다고 판단.
3. **모범 풀이 공개 범위**: 테스트 코드에 풀이 명령이 포함되므로 테스트는 `_test.go`에만 있고 공개 저장소 방침은 강사 결정에 맡김. SPEC §6에서 허용됨.
