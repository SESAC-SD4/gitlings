package tui

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"gitlings-ko/exercises"
	"gitlings-ko/internal/progress"
	"gitlings-ko/internal/workspace"
)

func plain(md string, width int) string { return md }

func newModel(t *testing.T, prog *progress.Progress) Model {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	if prog == nil {
		prog = &progress.Progress{}
	}
	return New(exercises.All(), prog, plain)
}

func key(k string) tea.KeyMsg {
	switch k {
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "ctrl+c":
		return tea.KeyMsg{Type: tea.KeyCtrlC}
	case "ctrl+r":
		return tea.KeyMsg{Type: tea.KeyCtrlR}
	case "ctrl+x":
		return tea.KeyMsg{Type: tea.KeyCtrlX}
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "tab":
		return tea.KeyMsg{Type: tea.KeyTab}
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	case "pgdown":
		return tea.KeyMsg{Type: tea.KeyPgDown}
	}
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)}
}

func press(m Model, keys ...string) Model {
	var tm tea.Model = m
	for _, k := range keys {
		tm, _ = tm.Update(key(k))
	}
	return tm.(Model)
}

func TestStartsAtFirstUnpassed(t *testing.T) {
	m := newModel(t, &progress.Progress{Passed: []int{1, 2}})
	if m.Cursor() != 2 {
		t.Fatalf("3번 문제(인덱스 2)에서 시작해야 합니다: %d", m.Cursor())
	}
}

func TestMoveStaysInBounds(t *testing.T) {
	m := newModel(t, nil)
	m = press(m, "k", "up")
	if m.Cursor() != 0 {
		t.Fatalf("맨 위에서 더 올라가면 안 됩니다: %d", m.Cursor())
	}
	for i := 0; i < 30; i++ {
		m = press(m, "j")
	}
	if m.Cursor() != len(exercises.All())-1 {
		t.Fatalf("맨 아래에서 멈춰야 합니다: %d", m.Cursor())
	}
	m = press(m, "down", "k")
	if m.Cursor() != len(exercises.All())-2 {
		t.Fatalf("down과 k가 동작해야 합니다: %d", m.Cursor())
	}
}

func TestCheckWithoutFolder(t *testing.T) {
	m := press(newModel(t, nil), "c")
	if !strings.Contains(m.Status(), "Ctrl+R") {
		t.Fatalf("연습 폴더가 없으면 Ctrl+R로 만들라고 해야 합니다: %s", m.Status())
	}
}

func TestResetCreatesMissingFolderWithoutAsking(t *testing.T) {
	m := press(newModel(t, nil), "ctrl+r")
	if !workspace.Exists("01-first-commit") || m.Mode() != ModePrompt {
		t.Fatal("연습 폴더가 없으면 묻지 않고 바로 만들어야 합니다")
	}
	if !strings.Contains(m.Content(), "cd ~/gitlings-work/01-first-commit") {
		t.Fatal("오른쪽 칸에 연습 폴더 경로가 있어야 합니다")
	}
}

func TestResetAsksBeforeOverwriting(t *testing.T) {
	m := press(newModel(t, nil), "ctrl+r")
	marker := filepath.Join(workspace.Dir("01-first-commit"), "mine.txt")
	if err := os.WriteFile(marker, []byte("내 작업"), 0o644); err != nil {
		t.Fatal(err)
	}
	m = press(m, "ctrl+r")
	if m.Mode() != ModeConfirm || !strings.Contains(m.Status(), "(y/n)") {
		t.Fatalf("되돌릴지 물어야 합니다: %v %s", m.Mode(), m.Status())
	}
	m = press(m, "n")
	if _, err := os.Stat(marker); err != nil || !strings.Contains(m.Status(), "취소") {
		t.Fatal("n이면 학생 작업을 지키고 취소했다고 알려야 합니다")
	}
	m = press(m, "r", "enter")
	if _, err := os.Stat(marker); err != nil {
		t.Fatal("확인 질문에서 Enter는 '아니오'여야 합니다(실수로 지워지지 않게)")
	}
	m = press(m, "r", "y")
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("확인 질문에서 y는 '예'여야 합니다")
	}
	if m.Mode() != ModePrompt {
		t.Fatal("확인이 끝나면 문제 화면으로 돌아가야 합니다")
	}
}

func TestSKeyNoLongerExists(t *testing.T) {
	press(newModel(t, nil), "s")
	if workspace.Exists("01-first-commit") {
		t.Fatal("s 키는 없앴습니다(연습 폴더는 켤 때 자동으로 만든다)")
	}
}

func TestCheckFailThenPass(t *testing.T) {
	m := press(newModel(t, nil), "ctrl+r", "enter")
	if !strings.Contains(m.Banner(), "아직입니다") || m.Cursor() != 0 {
		t.Fatalf("실패하면 그 자리에서 이유를 보여 줘야 합니다: %s", m.Banner())
	}
	solveFirst(t, workspace.Dir("01-first-commit"))
	m = press(m, "enter")
	if !strings.Contains(m.Banner(), "통과") || m.Cursor() != 1 {
		t.Fatalf("통과하면 배너를 남기고 다음 문제로 가야 합니다: %s %d", m.Banner(), m.Cursor())
	}
	if !strings.Contains(m.Content(), "통과") {
		t.Fatal("다음 문제 화면 위에 통과 배너가 보여야 합니다")
	}
	saved, err := progress.Load()
	if err != nil || !saved.IsPassed(1) {
		t.Fatal("통과 기록이 파일에 저장되어야 합니다")
	}
	m = press(m, "j")
	if m.Banner() != "" {
		t.Fatal("직접 이동하면 배너를 지워야 합니다")
	}
}

func TestHintLevels(t *testing.T) {
	m := press(newModel(t, nil), "tab")
	if m.Mode() != ModeHint || m.HintLevel() != 1 || !strings.Contains(m.Content(), "힌트 1/3") {
		t.Fatalf("Tab 한 번이면 1단계: %v %d", m.Mode(), m.HintLevel())
	}
	m = press(m, "h", "tab", "tab")
	if m.HintLevel() != 3 || !strings.Contains(m.Content(), "힌트 3/3") || !strings.Contains(m.Content(), "힌트 1/3") {
		t.Fatalf("3단계에서 멈추고 앞 단계도 함께 보여야 합니다: %d", m.HintLevel())
	}
	m = press(m, "x")
	if m.Mode() != ModePrompt || m.HintLevel() != 3 {
		t.Fatal("다른 키면 문제 설명으로 돌아가되 단계는 기억해야 합니다")
	}
	m = press(m, "j")
	if m.HintLevel() != 0 {
		t.Fatal("문제를 바꾸면 힌트 단계는 처음으로")
	}
}

func TestKoreanInputKeys(t *testing.T) {
	m := press(newModel(t, nil), "ㄱ")
	if workspace.Exists("01-first-commit") || !strings.Contains(m.Status(), "한/영") {
		t.Fatal("ㄱ(r 자리)은 되돌리기로 쓰지 않아야 합니다. 입력기가 붙잡았다가 Enter와 함께 보내면 작업이 지워집니다")
	}
	m = press(m, "ctrl+r")
	marker := filepath.Join(workspace.Dir("01-first-commit"), "mine.txt")
	if err := os.WriteFile(marker, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	m = press(m, "ctrl+r", "ㅜ", "esc")
	if _, err := os.Stat(marker); err != nil || m.Mode() != ModePrompt {
		t.Fatal("확인 질문에서 한글(ㅜ 등)로는 지워지지 않고, Esc로 취소되어야 합니다")
	}
	m = press(m, "ctrl+r", "ㅛ")
	if _, err := os.Stat(marker); err != nil || m.Mode() != ModeConfirm || !strings.Contains(m.Status(), "영문") {
		t.Fatal("확인 질문에서 한글이 들어오면 지우지 말고, 질문을 연 채 영문 y를 안내해야 합니다")
	}
	m = press(m, "y")
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("영문 y는 '예'여야 합니다")
	}
	m = press(m, "ㅗ")
	if m.HintLevel() != 1 {
		t.Fatal("ㅗ는 h로 동작해야 합니다")
	}
	m = press(m, "가")
	if !strings.Contains(m.Status(), "한/영") {
		t.Fatalf("알 수 없는 한글 입력이면 한/영 키를 안내해야 합니다: %s", m.Status())
	}
}

func TestScrollIsNotUndone(t *testing.T) {
	m := newModel(t, nil)
	tm, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 12})
	m = press(tm.(Model), "pgdown")
	if m.vp.YOffset == 0 {
		t.Fatal("PgDn으로 내린 스크롤이 바로 맨 위로 되돌아가면 안 됩니다")
	}
}

func TestEscClosesHintAndConfirmWithoutQuitting(t *testing.T) {
	m := newModel(t, nil)
	tm, cmd := press(m, "tab").Update(key("esc"))
	m = tm.(Model)
	if cmd != nil || m.Mode() != ModePrompt {
		t.Fatal("힌트 화면에서 Esc는 문제 설명으로 돌아가야 하고, 프로그램을 끄면 안 됩니다")
	}
	m = press(m, "ctrl+r", "ctrl+r")
	tm, cmd = m.Update(key("esc"))
	m = tm.(Model)
	if cmd != nil || m.Mode() != ModePrompt || !strings.Contains(m.Status(), "취소") {
		t.Fatal("확인 질문에서 Esc는 취소여야 합니다")
	}
}

func TestOnlyCtrlCQuits(t *testing.T) {
	m := newModel(t, nil)
	for _, k := range []string{"q", "esc", "ㅂ"} {
		if _, cmd := m.Update(key(k)); cmd != nil {
			if _, ok := cmd().(tea.QuitMsg); ok {
				t.Fatalf("%s로는 종료되지 않아야 합니다", k)
			}
		}
	}
	_, cmd := m.Update(key("ctrl+c"))
	if cmd == nil {
		t.Fatal("Ctrl+C는 종료해야 합니다")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("Ctrl+C는 tea.Quit이어야 합니다")
	}
}

func TestResetClearsPassMark(t *testing.T) {
	prog := &progress.Progress{}
	m := press(newModel(t, prog), "ctrl+r")
	solveFirst(t, workspace.Dir("01-first-commit"))
	m = press(m, "enter", "k", "ctrl+r", "y")
	if prog.IsPassed(1) {
		t.Fatal("이 문제를 처음으로 되돌리면 ✅ 표시도 지워야 합니다")
	}
	saved, _ := progress.Load()
	if saved.IsPassed(1) {
		t.Fatal("지운 통과 기록이 파일에도 저장되어야 합니다")
	}
}

func TestFullReset(t *testing.T) {
	prog := &progress.Progress{Passed: []int{1, 5}}
	m := newModel(t, prog)
	m = press(m, "ctrl+r")
	marker := filepath.Join(workspace.Dir(m.exs[m.Cursor()].Slug), "mine.txt")
	if err := os.WriteFile(marker, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	m = press(m, "ctrl+x")
	if m.Mode() != ModeConfirm || !strings.Contains(m.Status(), "전체") {
		t.Fatalf("전체 초기화는 먼저 물어야 합니다: %s", m.Status())
	}
	m = press(m, "enter")
	if _, err := os.Stat(marker); err != nil || !prog.IsPassed(1) {
		t.Fatal("Enter로는 전체 초기화되지 않아야 합니다")
	}
	m = press(m, "ctrl+x", "y")
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("전체 초기화는 모든 연습 폴더를 처음 상태로 되돌려야 합니다")
	}
	for _, ex := range exercises.All() {
		if !workspace.Exists(ex.Slug) {
			t.Fatalf("전체 초기화 뒤 %s 연습 폴더가 있어야 합니다", ex.Slug)
		}
	}
	if len(prog.Passed) != 0 || m.Cursor() != 0 {
		t.Fatalf("✅ 기록을 모두 지우고 1번으로 돌아가야 합니다: %v %d", prog.Passed, m.Cursor())
	}
}

// solveFirst는 1번 문제(첫 사진)를 학생처럼 진짜 git으로 푼다.
func solveFirst(t *testing.T, dir string) {
	t.Helper()
	cmd := exec.Command("bash", "-c", `git init -b main && git add words.md && git commit -m "첫 번째 커밋"`)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL="+os.DevNull,
		"GIT_AUTHOR_NAME=학생", "GIT_AUTHOR_EMAIL=student@example.com",
		"GIT_COMMITTER_NAME=학생", "GIT_COMMITTER_EMAIL=student@example.com")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("1번 풀이 실패: %v\n%s", err, out)
	}
}
