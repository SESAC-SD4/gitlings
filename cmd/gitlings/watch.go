package main

import (
	"fmt"
	"gitlings-ko/exercises"
	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/git"
	"gitlings-ko/internal/progress"
	"golang.org/x/term"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

const (
	ansiClear  = "\033[2J\033[H"
	ansiBold   = "\033[1m"
	ansiReset  = "\033[0m"
	ansiGreen  = "\033[32m"
	ansiRed    = "\033[31m"
	ansiYellow = "\033[33m"
	ansiCyan   = "\033[36m"
	ansiGray   = "\033[90m"
)

func runWatch(id int) {
	var ex *exercise.Exercise
	if id > 0 {
		ex = exercises.ByID(id)
		if ex == nil {
			fmt.Printf("❌ %d번 문제가 없습니다.\n", id)
			os.Exit(1)
		}
	} else {
		ex, _ = currentExercise()
	}

	if ex == nil {
		// 아직 시작한 문제 없음 — 목록 보여주고 번호 입력
		runList()
		fmt.Print("  시작할 문제 번호를 입력하세요: ")
		var n int
		fmt.Scanf("%d", &n)
		ex = exercises.ByID(n)
		if ex == nil {
			fmt.Printf("❌ %d번 문제가 없습니다.\n", n)
			os.Exit(1)
		}
	}

	watchExercise(ex)
}

// watchExercise는 ex 문제를 watch mode로 실행한다.
// 재귀가 아니라 루프로 문제 전환을 처리한다.
func watchExercise(startEx *exercise.Exercise) {
	ex := startEx
	for {
		next := watchOne(ex)
		if next == nil {
			return
		}
		ex = next
	}
}

// watchOne은 한 문제의 watch loop를 실행하고,
// 다음 문제로 전환해야 하면 그 Exercise를 반환하고,
// 종료해야 하면 nil을 반환한다.
func watchOne(ex *exercise.Exercise) *exercise.Exercise {
	workDir := filepath.Join(baseDir(), ex.Slug)

	// 작업 폴더 없으면 setup
	if _, err := os.Stat(workDir); os.IsNotExist(err) {
		g := git.New(workDir)
		fmt.Printf("\n🔧 %s 문제를 준비하는 중...\n", ex.Name)
		if err2 := ex.Setup(workDir, g); err2 != nil {
			fmt.Printf("❌ setup 오류: %v\n", err2)
			os.Exit(1)
		}
		writeMarker(ex.ID, workDir, g)
		p, _ := progress.Load()
		p.LastStarted = ex.ID
		p.Save()
	}

	fd := int(os.Stdin.Fd())
	oldState, rawErr := term.MakeRaw(fd)
	if rawErr != nil {
		watchNoRaw(ex, workDir)
		return nil
	}
	defer term.Restore(fd, oldState)

	recheckCh := make(chan struct{}, 1)
	keyCh := make(chan byte, 4)

	// 파일 변경 감시 (폴링)
	go func() {
		var lastMod time.Time
		for {
			time.Sleep(600 * time.Millisecond)
			mod := latestModTime(workDir)
			if mod.After(lastMod) {
				lastMod = mod
				select {
				case recheckCh <- struct{}{}:
				default:
				}
			}
		}
	}()

	// 키보드 읽기
	go func() {
		buf := make([]byte, 1)
		for {
			n, err := os.Stdin.Read(buf)
			if n > 0 && err == nil {
				keyCh <- buf[0]
			}
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	recheckCh <- struct{}{} // 최초 체크

	var lastResult exercise.Result
	hintLevel := loadHintLevel(ex.ID)
	currentHint := hintAt(ex, hintLevel)

	for {
		select {
		case <-recheckCh:
			g := git.New(workDir)
			if _, e := os.Stat(filepath.Join(workDir, ".git")); e == nil {
				g.AddExclude(".gitlings_id")
			}
			result, err := ex.Check(workDir, g)
			if err != nil {
				result = exercise.Result{Message: fmt.Sprintf("채점 오류: %v", err)}
			}
			wasNotPassed := !lastResult.Passed
			lastResult = result
			if wasNotPassed && result.Passed {
				p, _ := progress.Load()
				p.MarkPassed(ex.ID)
				p.Save()
			}
			drawWatch(ex, workDir, lastResult, currentHint)

		case key := <-keyCh:
			switch key {
			case 'h', 'H':
				if hintLevel < len(ex.Hints) {
					hintLevel++
					saveHintLevel(ex.ID, hintLevel)
					currentHint = hintAt(ex, hintLevel)
					drawWatch(ex, workDir, lastResult, currentHint)
				}
			case 'n', 'N':
				nextEx := exercises.ByID(ex.ID + 1)
				if nextEx != nil {
					term.Restore(fd, oldState)
					return nextEx
				}
			case 'p', 'P':
				prevEx := exercises.ByID(ex.ID - 1)
				if prevEx != nil {
					term.Restore(fd, oldState)
					return prevEx
				}
			case 'r', 'R':
				term.Restore(fd, oldState)
				fmt.Print(ansiClear)
				doResetExercise(ex, workDir)
				hintLevel = 0
				currentHint = ""
				oldState, _ = term.MakeRaw(fd)
				lastResult = exercise.Result{}
				recheckCh <- struct{}{}
			case 'l', 'L':
				term.Restore(fd, oldState)
				fmt.Print(ansiClear)
				runList()
				fmt.Print("\n  문제 번호로 이동 (Enter = 현재 유지): ")
				var n int
				fmt.Scanf("%d", &n)
				if n > 0 {
					chosen := exercises.ByID(n)
					if chosen != nil {
						return chosen
					}
				}
				oldState, _ = term.MakeRaw(fd)
				recheckCh <- struct{}{}
			case 'q', 'Q', 3: // 3 = Ctrl+C
				term.Restore(fd, oldState)
				fmt.Print(ansiClear)
				fmt.Println("gitlings-ko를 종료합니다. 수고하셨습니다!")
				os.Exit(0)
			}
		case <-sigCh:
			term.Restore(fd, oldState)
			os.Exit(0)
		}
	}
}

const divLong = "------------------------------------------------------------"
const divShort = "----------------------------------------"

func drawWatch(ex *exercise.Exercise, workDir string, result exercise.Result, hint string) {
	p, _ := progress.Load()
	all := exercises.All()
	passed := 0
	for _, e := range all {
		if p.IsPassed(e.ID) {
			passed++
		}
	}

	var b strings.Builder
	fmt.Fprint(&b, ansiClear)

	// 헤더
	fmt.Fprintf(&b, "%sgitlings-ko%s   pass:%d  todo:%d\n",
		ansiBold, ansiReset, passed, len(all)-passed)
	fmt.Fprintln(&b, divLong)
	fmt.Fprintln(&b)

	// 작업 폴더
	home, _ := os.UserHomeDir()
	displayDir := strings.Replace(workDir, home, "~", 1)
	fmt.Fprintf(&b, "  %s%s%s\n\n", ansiCyan, displayDir, ansiReset)

	// 문제 제목
	fmt.Fprintf(&b, "%s  %s  (%d/%d)%s\n", ansiBold, ex.Name, ex.ID, len(all), ansiReset)
	fmt.Fprintln(&b, divShort)
	fmt.Fprintln(&b)

	// 문제 설명
	for _, l := range strings.Split(ex.Desc, "\n") {
		fmt.Fprintf(&b, "  %s\n", l)
	}
	fmt.Fprintln(&b)

	// 채점 결과
	fmt.Fprintln(&b, divShort)
	if result.Passed {
		fmt.Fprintf(&b, "%s  PASS!%s\n", ansiGreen+ansiBold, ansiReset)
	} else if result.Message != "" {
		for _, l := range strings.Split(result.Message, "\n") {
			fmt.Fprintf(&b, "  %s%s%s\n", ansiRed, l, ansiReset)
		}
	}
	fmt.Fprintln(&b, divShort)

	// 힌트
	if hint != "" {
		fmt.Fprintln(&b)
		fmt.Fprintf(&b, "%s  [힌트]%s\n", ansiYellow+ansiBold, ansiReset)
		for _, l := range strings.Split(hint, "\n") {
			fmt.Fprintf(&b, "  %s%s%s\n", ansiYellow, l, ansiReset)
		}
	}

	// 푸터
	fmt.Fprintln(&b)
	if result.Passed {
		if exercises.ByID(ex.ID+1) != nil {
			fmt.Fprintf(&b, "  [n]다음  [p]이전  [l]목록  [q]종료\n")
		} else {
			fmt.Fprintf(&b, "  모든 문제 완료!  [l]목록  [q]종료\n")
		}
	} else {
		hintKey := ""
		if len(ex.Hints) > 0 {
			hintKey = "[h]힌트  "
		}
		fmt.Fprintf(&b, "  [r]초기화  %s[n]다음  [l]목록  [q]종료   %s(감지중)%s\n",
			hintKey, ansiGray, ansiReset)
	}

	// raw 모드에서는 \n이 줄 맨 앞으로 돌아가지 않으므로 \r\n으로 바꿔 출력한다.
	fmt.Print(strings.ReplaceAll(b.String(), "\n", "\r\n"))
}

func hintAt(ex *exercise.Exercise, level int) string {
	if level <= 0 || level > len(ex.Hints) {
		return ""
	}
	return ex.Hints[level-1]
}

func latestModTime(dir string) time.Time {
	var latest time.Time
	filepath.Walk(dir, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.ModTime().After(latest) {
			latest = info.ModTime()
		}
		return nil
	})
	return latest
}

func loadHintLevel(id int) int {
	hintFile := filepath.Join(baseDir(), fmt.Sprintf(".hint_%d", id))
	data, err := os.ReadFile(hintFile)
	if err != nil {
		return 0
	}
	var level int
	fmt.Sscan(strings.TrimSpace(string(data)), &level)
	return level
}

func saveHintLevel(id, level int) {
	hintFile := filepath.Join(baseDir(), fmt.Sprintf(".hint_%d", id))
	os.WriteFile(hintFile, []byte(fmt.Sprintf("%d", level)), 0644)
}

func watchNoRaw(ex *exercise.Exercise, workDir string) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("\n감시 중: %s\n변경 감지 시 자동 채점합니다. Ctrl+C 로 종료.\n\n", workDir)
	var lastMod time.Time
	var lastMsg string
	for {
		select {
		case <-sigCh:
			os.Exit(0)
		default:
			time.Sleep(600 * time.Millisecond)
			mod := latestModTime(workDir)
			if mod.After(lastMod) {
				lastMod = mod
				g := git.New(workDir)
				result, err := ex.Check(workDir, g)
				if err != nil {
					result = exercise.Result{Message: fmt.Sprintf("채점 오류: %v", err)}
				}
				if result.Message != lastMsg {
					lastMsg = result.Message
					fmt.Print(ansiClear)
					drawWatch(ex, workDir, result, "")
				}
			}
		}
	}
}
