package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/progress"
	"gitlings-ko/internal/workspace"
)

type Renderer func(markdown string, width int) string

type Mode int

const (
	ModePrompt Mode = iota
	ModeHint
	ModeConfirm
)

const readyStatus = "다른 터미널에서 연습 폴더로 가서 풀고, Enter로 채점하세요"

type Model struct {
	exs        []*exercise.Exercise
	prog       *progress.Progress
	render     Renderer
	cursor     int
	mode       Mode
	confirmAll bool
	hintLevel  int
	status     string
	banner     string
	width      int
	height     int
	vp         viewport.Model
}

func New(exs []*exercise.Exercise, prog *progress.Progress, render Renderer) Model {
	m := Model{exs: exs, prog: prog, render: render, vp: viewport.New(0, 0), status: readyStatus}
	for i, ex := range exs {
		if !prog.IsPassed(ex.ID) {
			m.cursor = i
			break
		}
	}
	m.refresh()
	return m
}

func (m Model) Cursor() int             { return m.cursor }
func (m Model) Mode() Mode              { return m.mode }
func (m Model) HintLevel() int          { return m.hintLevel }
func (m Model) Status() string          { return m.status }
func (m Model) Banner() string          { return m.banner }
func (m Model) Init() tea.Cmd           { return nil }
func (m Model) cur() *exercise.Exercise { return m.exs[m.cursor] }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.layout()
		m.refresh()
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	var cmd tea.Cmd
	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

// 한글 입력 상태(두벌식)로 누른 키를 원래 영문 키로 돌려 놓는다. ㄱ(r)은 넣지 않는다: 입력기가 붙잡았다가 Enter와 함께 보내면 되돌리기가 '예'로 실행된다.
var hangulKeys = map[string]string{
	"ㅗ": "h", "ㅓ": "j", "ㅏ": "k", "ㅊ": "c",
}

// 한글 입력기는 자모를 조합 중으로 붙잡아 두므로, 입력기를 거치지 않는 키를 함께 둔다.
// 종료(Ctrl+C)·초기화(Ctrl+R·Ctrl+X)는 실수로 누르기 어려운 Ctrl 조합만 쓴다.
var aliasKeys = map[string]string{"enter": "c", "tab": "h", "ctrl+r": "r"}

func isASCII(s string) bool {
	for _, r := range s {
		if r > 127 {
			return false
		}
	}
	return true
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	k := msg.String()
	if k == "ctrl+c" {
		return m, tea.Quit
	}
	// 확인 질문의 '예'는 영문 y 하나뿐이다. Enter나 다른 키는 취소, 한글은 질문을 연 채 안내만 한다.
	if m.mode == ModeConfirm {
		switch {
		case k == "y":
			m.mode = ModePrompt
			if m.confirmAll {
				m.resetAll()
			} else {
				m.prepare()
			}
		case !isASCII(k):
			m.status = "한글로 들어오고 있습니다. 한/영 키로 영문으로 바꾼 뒤 y를 누르세요 (취소는 Esc)"
		default:
			m.mode = ModePrompt
			m.status = "취소했습니다. 연습 폴더는 그대로입니다"
		}
		m.refresh()
		return m, nil
	}
	if k == "esc" {
		if m.mode == ModeHint {
			m.mode, m.status = ModePrompt, readyStatus
			m.refresh()
		}
		return m, nil
	}
	if mapped, ok := hangulKeys[k]; ok {
		k = mapped
	} else if !isASCII(k) {
		m.status = "키가 한글로 들어오고 있습니다. 한/영 키를 눌러 영문으로 바꾸세요"
		return m, nil
	}
	if alias, ok := aliasKeys[k]; ok {
		k = alias
	}
	if m.mode == ModeHint && k != "h" {
		m.mode = ModePrompt
	}
	switch k {
	case "up", "k":
		m.move(-1)
	case "down", "j":
		m.move(1)
	case "r":
		m.askOrPrepare()
	case "ctrl+x":
		m.mode, m.confirmAll = ModeConfirm, true
		m.status = fmt.Sprintf("전체 초기화: %d문제의 연습 폴더와 ✅ 기록을 모두 처음으로 되돌릴까요? (y/n)", len(m.exs))
	case "c":
		m.check()
	case "h":
		m.hint()
	default:
		// 내용이 바뀌지 않았으므로 refresh(맨 위로 이동)를 부르지 않는다. 부르면 스크롤이 바로 되돌아간다.
		var cmd tea.Cmd
		m.vp, cmd = m.vp.Update(msg)
		return m, cmd
	}
	m.refresh()
	return m, nil
}

func (m *Model) move(d int) {
	n := m.cursor + d
	if n < 0 || n >= len(m.exs) {
		return
	}
	m.cursor, m.hintLevel, m.banner, m.mode, m.status = n, 0, "", ModePrompt, readyStatus
}

func (m *Model) askOrPrepare() {
	if workspace.Exists(m.cur().Slug) {
		m.mode, m.confirmAll = ModeConfirm, false
		m.status = "이 문제를 처음 상태로 되돌릴까요? 지금까지 한 작업과 ✅ 표시가 지워집니다 (y/n)"
		return
	}
	m.prepare()
}

func (m *Model) prepare() {
	ex := m.cur()
	if err := workspace.Prepare(ex); err != nil {
		m.status = "❌ 연습 폴더를 만들지 못했습니다: " + err.Error()
		return
	}
	m.prog.LastStarted = ex.ID
	m.prog.ResetHint(ex.ID)
	m.prog.Unpass(ex.ID)
	m.hintLevel, m.banner = 0, ""
	if err := m.prog.Save(); err != nil {
		m.status = "❌ 진행 상황을 저장하지 못했습니다: " + err.Error()
		return
	}
	m.status = "✅ 처음 상태로 만들었습니다. 다른 터미널에서: cd " + workspace.Display(ex.Slug)
}

func (m *Model) resetAll() {
	for _, ex := range m.exs {
		if err := workspace.Prepare(ex); err != nil {
			m.status = "❌ " + ex.Title + " 연습 폴더를 되돌리지 못했습니다: " + err.Error()
			return
		}
	}
	m.prog.Passed, m.prog.HintLevel, m.prog.LastStarted = nil, nil, 0
	m.cursor, m.hintLevel, m.banner = 0, 0, ""
	if err := m.prog.Save(); err != nil {
		m.status = "❌ 진행 상황을 저장하지 못했습니다: " + err.Error()
		return
	}
	m.status = fmt.Sprintf("✅ %d문제를 모두 처음 상태로 되돌렸습니다", len(m.exs))
}

func (m *Model) check() {
	ex := m.cur()
	if !workspace.Exists(ex.Slug) {
		m.status = "연습 폴더가 없습니다. Ctrl+R을 누르면 만듭니다"
		return
	}
	r := ex.Check(workspace.Dir(ex.Slug))
	if !r.Passed {
		m.banner = "❌ 아직입니다 — " + r.Message
		m.status = "아직입니다. Tab으로 힌트를 볼 수 있습니다"
		return
	}
	m.prog.MarkPassed(ex.ID)
	if err := m.prog.Save(); err != nil {
		m.status = "❌ 진행 상황을 저장하지 못했습니다: " + err.Error()
		return
	}
	m.banner = "✅ " + ex.Title + " 통과! " + r.Message
	m.status = "통과!"
	if m.cursor < len(m.exs)-1 {
		m.cursor++
		m.hintLevel = 0
	}
}

func (m *Model) hint() {
	ex := m.cur()
	if m.hintLevel < len(ex.Hints) {
		m.hintLevel++
	}
	m.mode = ModeHint
	m.status = fmt.Sprintf("힌트 %d/%d — Tab을 다시 누르면 다음 단계, Esc로 닫기", m.hintLevel, len(ex.Hints))
}

func (m Model) Content() string {
	ex := m.cur()
	var b strings.Builder
	if m.banner != "" {
		b.WriteString("> " + strings.ReplaceAll(m.banner, "\n", "\n> ") + "\n\n")
	}
	b.WriteString("# " + ex.Title + "\n\n")
	if m.mode == ModeHint {
		for i := 0; i < m.hintLevel; i++ {
			fmt.Fprintf(&b, "**힌트 %d/%d**\n\n%s\n\n", i+1, len(ex.Hints), ex.Hints[i])
		}
	} else {
		b.WriteString(strings.TrimSpace(ex.Prompt) + "\n\n")
	}
	b.WriteString("---\n\n연습 폴더: `cd " + workspace.Display(ex.Slug) + "`\n")
	if !workspace.Exists(ex.Slug) {
		b.WriteString("\n연습 폴더가 없습니다. Ctrl+R을 누르면 만듭니다.\n")
	}
	return b.String()
}

func (m *Model) refresh() {
	m.vp.SetContent(m.render(m.Content(), m.vp.Width))
	m.vp.GotoTop()
}
