package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"

	"gitlings-ko/exercises"
	"gitlings-ko/internal/console"
	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/git"
	"gitlings-ko/internal/progress"
	"gitlings-ko/internal/tui"
	"gitlings-ko/internal/workspace"
)

func main() {
	console.Init()
	if err := git.CheckVersion(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout))
}

func self() string {
	if runtime.GOOS == "windows" {
		return "./gitlings.exe"
	}
	return "./gitlings"
}

func run(args []string, in io.Reader, out io.Writer) int {
	prog, err := progress.Load()
	if err != nil {
		fmt.Fprintf(out, "❌ 진행 상황을 읽지 못했습니다: %v\n", err)
		return 1
	}
	cmd := "watch"
	if len(args) > 0 {
		cmd = args[0]
	}
	switch cmd {
	case "watch":
		if err := tui.Run(exercises.All(), prog); err != nil {
			fmt.Fprintf(out, "❌ 화면을 띄우지 못했습니다: %v\n대신 `%s list`, `%s start N`을 쓸 수 있습니다\n", err, self(), self())
			return 1
		}
		return 0
	case "list":
		return cmdList(out, prog)
	case "start":
		return cmdStart(args[1:], in, out, prog)
	case "check":
		return cmdCheck(out, prog)
	case "hint":
		return cmdHint(out, prog)
	case "reset":
		return cmdReset(out, prog)
	case "help", "-h", "--help":
		cmdHelp(out)
		return 0
	}
	fmt.Fprintf(out, "❌ 모르는 명령입니다: %s\n\n", cmd)
	cmdHelp(out)
	return 1
}

func cmdHelp(out io.Writer) {
	s := self()
	fmt.Fprintf(out, `gitlings — Git 연습 문제

  %[1]s            분할 화면으로 시작 (= %[1]s watch)
  %[1]s list       문제 목록
  %[1]s start N    N번 연습 폴더를 만들고 문제를 보여 줍니다
  %[1]s check      지금 문제를 채점합니다
  %[1]s hint       힌트를 한 단계씩 보여 줍니다
  %[1]s reset      지금 문제를 처음 상태로 되돌립니다
  %[1]s help       이 도움말
`, s)
}

func cmdList(out io.Writer, prog *progress.Progress) int {
	all := exercises.All()
	n := 0
	for _, ex := range all {
		mark := "⬜"
		if prog.IsPassed(ex.ID) {
			mark = "✅"
			n++
		}
		fmt.Fprintf(out, "%s %s\n", mark, ex.Title)
	}
	fmt.Fprintf(out, "\n%d/%d 통과\n", n, len(all))
	return 0
}

func cmdStart(args []string, in io.Reader, out io.Writer, prog *progress.Progress) int {
	all := exercises.All()
	if len(args) != 1 {
		fmt.Fprintf(out, "사용법: %s start N (N은 1~%d)\n", self(), len(all))
		return 1
	}
	id, err := strconv.Atoi(args[0])
	ex := exercises.ByID(id)
	if err != nil || ex == nil {
		fmt.Fprintf(out, "❌ %s번 문제가 없습니다 (1~%d)\n", args[0], len(all))
		return 1
	}
	overwrite := true
	if workspace.Exists(ex.Slug) {
		fmt.Fprintf(out, "이미 연습 폴더가 있습니다: %s\n처음 상태로 덮어쓸까요? (y/n) ", workspace.Display(ex.Slug))
		ans, _ := bufio.NewReader(in).ReadString('\n')
		overwrite = strings.EqualFold(strings.TrimSpace(ans), "y")
		if !overwrite {
			fmt.Fprintln(out, "그대로 두었습니다.")
		}
	}
	if overwrite {
		if err := workspace.Prepare(ex); err != nil {
			fmt.Fprintf(out, "❌ 연습 폴더를 만들지 못했습니다: %v\n", err)
			return 1
		}
		prog.ResetHint(ex.ID)
	}
	prog.LastStarted = ex.ID
	if err := prog.Save(); err != nil {
		fmt.Fprintf(out, "❌ 진행 상황을 저장하지 못했습니다: %v\n", err)
		return 1
	}
	fmt.Fprintf(out, "\n# %s\n\n%s\n\n이제 다른 터미널에서:\n  cd %s\n\n풀었으면: %s check   ·   막히면: %s hint\n",
		ex.Title, strings.TrimSpace(ex.Prompt), workspace.Display(ex.Slug), self(), self())
	return 0
}

func current(out io.Writer, prog *progress.Progress) *exercise.Exercise {
	cwd, _ := os.Getwd()
	ex := workspace.Detect(cwd, exercises.All(), prog.LastStarted)
	if ex == nil {
		fmt.Fprintf(out, "먼저 `%s start N`으로 문제를 시작하세요\n", self())
	}
	return ex
}

func cmdCheck(out io.Writer, prog *progress.Progress) int {
	ex := current(out, prog)
	if ex == nil {
		return 1
	}
	if !workspace.Exists(ex.Slug) {
		fmt.Fprintf(out, "연습 폴더가 없습니다: %s\n`%s start %d`로 만드세요\n", workspace.Display(ex.Slug), self(), ex.ID)
		return 1
	}
	r := ex.Check(workspace.Dir(ex.Slug))
	if !r.Passed {
		fmt.Fprintf(out, "❌ %s — 아직입니다\n%s\n\n막히면: %s hint\n", ex.Title, r.Message, self())
		return 1
	}
	prog.MarkPassed(ex.ID)
	if err := prog.Save(); err != nil {
		fmt.Fprintf(out, "❌ 진행 상황을 저장하지 못했습니다: %v\n", err)
		return 1
	}
	fmt.Fprintf(out, "✅ %s 통과!\n%s\n", ex.Title, r.Message)
	if next := exercises.ByID(ex.ID + 1); next != nil {
		fmt.Fprintf(out, "\n다음: %s start %d\n", self(), next.ID)
	}
	return 0
}

func cmdHint(out io.Writer, prog *progress.Progress) int {
	ex := current(out, prog)
	if ex == nil {
		return 1
	}
	level := prog.NextHint(ex.ID, len(ex.Hints))
	if err := prog.Save(); err != nil {
		fmt.Fprintf(out, "❌ 진행 상황을 저장하지 못했습니다: %v\n", err)
		return 1
	}
	for i := 0; i < level; i++ {
		fmt.Fprintf(out, "힌트 %d/%d\n%s\n\n", i+1, len(ex.Hints), ex.Hints[i])
	}
	if level == len(ex.Hints) {
		fmt.Fprintln(out, "마지막 힌트입니다.")
	}
	return 0
}

func cmdReset(out io.Writer, prog *progress.Progress) int {
	ex := current(out, prog)
	if ex == nil {
		return 1
	}
	if err := workspace.Prepare(ex); err != nil {
		fmt.Fprintf(out, "❌ 되돌리지 못했습니다: %v\n", err)
		return 1
	}
	prog.ResetHint(ex.ID)
	prog.Unpass(ex.ID)
	if err := prog.Save(); err != nil {
		fmt.Fprintf(out, "❌ 진행 상황을 저장하지 못했습니다: %v\n", err)
		return 1
	}
	fmt.Fprintf(out, "%s를 처음 상태로 되돌렸습니다.\n  cd %s\n", ex.Title, workspace.Display(ex.Slug))
	return 0
}
