package main

import (
	"bufio"
	"fmt"
	"gitlings-ko/exercises"
	"gitlings-ko/internal/console"
	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/git"
	"gitlings-ko/internal/progress"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	console.Init()

	if err := git.CheckVersion(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	args := os.Args[1:]
	if len(args) == 0 {
		runList()
		return
	}

	switch args[0] {
	case "list":
		runList()
	case "start":
		if len(args) < 2 {
			fmt.Println("사용법: gitlings start <번호>")
			fmt.Println("예시: gitlings start 1")
			os.Exit(1)
		}
		n, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("번호를 인식할 수 없습니다: %s\n", args[1])
			os.Exit(1)
		}
		runStart(n)
	case "check":
		runCheck()
	case "hint":
		runHint()
	case "reset":
		runReset()
	case "help", "--help", "-h":
		runHelp()
	default:
		fmt.Printf("알 수 없는 명령: %s\n", args[0])
		fmt.Println("사용 가능한 명령: list, start, check, hint, reset, help")
		os.Exit(1)
	}
}

func baseDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "gitlings-work")
}

func runList() {
	p, _ := progress.Load()
	all := exercises.All()

	fmt.Println()
	fmt.Println("  gitlings-ko — Git 연습 문제 목록")
	fmt.Println()
	fmt.Println("  번호  제목                         상태")
	fmt.Println("  ─────────────────────────────────────────")
	for _, ex := range all {
		status := "⬜"
		if p.IsPassed(ex.ID) {
			status = "✅"
		}
		optional := ""
		if ex.Optional {
			optional = " ★"
		}
		fmt.Printf("  %s   %s%s\n", status, ex.Name, optional)
	}
	fmt.Println()
	fmt.Println("  시작: gitlings start <번호>   예) gitlings start 1")
	fmt.Println("  채점: gitlings check")
	fmt.Println("  힌트: gitlings hint")
	fmt.Println()
}

func runStart(id int) {
	ex := exercises.ByID(id)
	if ex == nil {
		fmt.Printf("❌ %d번 문제가 없습니다. gitlings list 로 목록을 확인하세요.\n", id)
		os.Exit(1)
	}

	workDir := filepath.Join(baseDir(), ex.Slug)

	// 이미 있으면 확인
	if _, err := os.Stat(workDir); err == nil {
		fmt.Printf("⚠️  %s 폴더가 이미 있습니다.\n", workDir)
		fmt.Print("   덮어쓸까요? (y/N) ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(input)) != "y" {
			fmt.Println("   취소했습니다.")
			return
		}
		if err := os.RemoveAll(workDir); err != nil {
			fmt.Printf("❌ 폴더 삭제 실패: %v\n", err)
			os.Exit(1)
		}
		// bare 원격도 정리
		bareDir := workDir + ".git-remote"
		os.RemoveAll(bareDir)
		partnerDir := workDir + ".partner"
		os.RemoveAll(partnerDir)
	}

	g := git.New(workDir)
	fmt.Printf("\n🔧 %s 문제를 준비하는 중...\n", ex.Name)
	if err := ex.Setup(workDir, g); err != nil {
		fmt.Printf("❌ setup 오류: %v\n", err)
		os.Exit(1)
	}

	// 마커: .git/info/exclude에 .gitlings_id 추가
	// (.git이 있는 경우만)
	if _, err := os.Stat(filepath.Join(workDir, ".git")); err == nil {
		idFile := filepath.Join(workDir, ".gitlings_id")
		os.WriteFile(idFile, []byte(strconv.Itoa(id)), 0644)
		g.AddExclude(".gitlings_id")
	} else {
		// .git 없는 문제(ex01)는 파일만 만들어 둠
		idFile := filepath.Join(workDir, ".gitlings_id")
		os.WriteFile(idFile, []byte(strconv.Itoa(id)), 0644)
	}

	// 진행 상황 저장
	p, _ := progress.Load()
	p.LastStarted = id
	p.Save()

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  %s\n", ex.Name)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println(ex.Desc)
	fmt.Println()
	fmt.Printf("📁 작업 폴더로 이동하세요:\n")
	fmt.Printf("   cd %s\n", workDir)
	fmt.Println()
	fmt.Println("채점: gitlings check   힌트: gitlings hint")
	fmt.Println()
}

func currentExercise() (*exercise.Exercise, string) {
	cwd, _ := os.Getwd()

	// 현재 디렉토리가 ~/gitlings-work/<slug>인지 확인
	base := baseDir()
	if strings.HasPrefix(cwd, base+string(os.PathSeparator)) {
		rel := cwd[len(base)+1:]
		parts := strings.SplitN(rel, string(os.PathSeparator), 2)
		slug := parts[0]
		for _, ex := range exercises.All() {
			if ex.Slug == slug {
				return ex, cwd
			}
		}
	}

	// .gitlings_id 파일로 찾기
	idFile := filepath.Join(cwd, ".gitlings_id")
	if data, err := os.ReadFile(idFile); err == nil {
		if id, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil {
			if ex := exercises.ByID(id); ex != nil {
				return ex, cwd
			}
		}
	}

	// last_started fallback
	p, _ := progress.Load()
	if p.LastStarted > 0 {
		ex := exercises.ByID(p.LastStarted)
		if ex != nil {
			workDir := filepath.Join(base, ex.Slug)
			return ex, workDir
		}
	}

	return nil, ""
}

func runCheck() {
	ex, workDir := currentExercise()
	if ex == nil {
		fmt.Println("❌ 현재 어느 문제를 풀고 있는지 알 수 없습니다.")
		fmt.Println("   → gitlings start <번호> 로 문제를 시작하세요.")
		os.Exit(1)
	}

	g := git.New(workDir)
	// .git이 생겼으면 .gitlings_id를 exclude에 등록 (git status에 안 보이게)
	if _, err := os.Stat(filepath.Join(workDir, ".git")); err == nil {
		g.AddExclude(".gitlings_id")
	}
	result, err := ex.Check(workDir, g)
	if err != nil {
		fmt.Printf("❌ 채점 중 오류: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println(result.Message)
	fmt.Println()

	if result.Passed {
		p, _ := progress.Load()
		p.MarkPassed(ex.ID)
		p.Save()

		// 다음 문제 안내
		next := exercises.ByID(ex.ID + 1)
		if next != nil {
			fmt.Printf("   다음 문제: gitlings start %d\n", next.ID)
		} else {
			fmt.Println("   🎉 모든 문제를 완료했습니다!")
		}
		fmt.Println()
	}
}

func runHint() {
	ex, _ := currentExercise()
	if ex == nil {
		fmt.Println("❌ 현재 어느 문제를 풀고 있는지 알 수 없습니다.")
		fmt.Println("   → gitlings start <번호> 로 문제를 시작하세요.")
		os.Exit(1)
	}

	p, _ := progress.Load()
	// hint_level은 진행 상황에 저장하지 않고 단순히 순서대로 보여줌
	// 매 호출마다 다음 힌트를 보여주기 위해 별도 카운터 사용
	// 간단히: .gitlings_hint_level 파일 사용
	hintFile := filepath.Join(baseDir(), fmt.Sprintf(".hint_%d", ex.ID))
	level := 0
	if data, err := os.ReadFile(hintFile); err == nil {
		level, _ = strconv.Atoi(strings.TrimSpace(string(data)))
	}
	_ = p

	hints := ex.Hints
	if len(hints) == 0 {
		fmt.Println("이 문제에는 힌트가 없습니다.")
		return
	}

	if level >= len(hints) {
		fmt.Println("💡 모든 힌트를 다 보셨습니다!")
		fmt.Println()
		fmt.Printf("마지막 힌트:\n%s\n", hints[len(hints)-1])
		return
	}

	fmt.Printf("\n💡 힌트 %d/%d:\n\n%s\n\n", level+1, len(hints), hints[level])

	// 다음 레벨 저장
	os.WriteFile(hintFile, []byte(strconv.Itoa(level+1)), 0644)
	if level+1 < len(hints) {
		fmt.Printf("   더 자세한 힌트: gitlings hint (아직 %d단계 남음)\n\n", len(hints)-level-1)
	}
}

func runReset() {
	ex, workDir := currentExercise()
	if ex == nil {
		fmt.Println("❌ 현재 어느 문제를 풀고 있는지 알 수 없습니다.")
		fmt.Println("   → gitlings start <번호> 로 문제를 시작하세요.")
		os.Exit(1)
	}

	fmt.Printf("⚠️  %s 문제를 처음 상태로 되돌립니다.\n", ex.Name)
	fmt.Print("   계속할까요? (y/N) ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(input)) != "y" {
		fmt.Println("   취소했습니다.")
		return
	}

	// 작업 폴더 삭제 후 재생성
	os.RemoveAll(workDir)
	bareDir := workDir + ".git-remote"
	os.RemoveAll(bareDir)
	partnerDir := workDir + ".partner"
	os.RemoveAll(partnerDir)

	// 힌트 레벨 초기화
	hintFile := filepath.Join(baseDir(), fmt.Sprintf(".hint_%d", ex.ID))
	os.Remove(hintFile)

	g := git.New(workDir)
	if err := ex.Setup(workDir, g); err != nil {
		fmt.Printf("❌ setup 오류: %v\n", err)
		os.Exit(1)
	}

	// 마커 재생성
	if _, err := os.Stat(filepath.Join(workDir, ".git")); err == nil {
		idFile := filepath.Join(workDir, ".gitlings_id")
		os.WriteFile(idFile, []byte(strconv.Itoa(ex.ID)), 0644)
		g.AddExclude(".gitlings_id")
	} else {
		idFile := filepath.Join(workDir, ".gitlings_id")
		os.WriteFile(idFile, []byte(strconv.Itoa(ex.ID)), 0644)
	}

	fmt.Printf("\n✅ %s 문제가 초기화됐습니다.\n", ex.Name)
	fmt.Printf("   cd %s\n\n", workDir)
}

func runHelp() {
	fmt.Print(`
gitlings-ko — 터미널에서 진짜 git으로 배우는 한국어 연습 문제

사용법:
  gitlings              문제 목록 보기
  gitlings list         문제 목록 보기
  gitlings start <N>    N번 문제 시작 (폴더 만들고 설명 출력)
  gitlings check        현재 문제 채점
  gitlings hint         힌트 보기 (3단계, 한 번에 하나씩)
  gitlings reset        현재 문제 초기화
  gitlings help         이 도움말

작업 폴더: ~/gitlings-work/
진행 기록: ~/gitlings-work/.progress.json

✅ = 통과   ⬜ = 미완료   ★ = 선택 문제
`)
}
