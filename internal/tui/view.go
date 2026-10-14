package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	listWidth = 30 // 가장 긴 줄 "  ⬜ 12★ 날린 커밋 되살리기"(27칸) + 좌우 여백 2
	paneExtra = 4  // 테두리 2 + 좌우 여백 2
)

var (
	pane     = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63")).Padding(0, 1)
	selected = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("226"))
	passed   = lipgloss.NewStyle().Foreground(lipgloss.Color("78"))
	keysTxt  = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

const keyHelp = "Enter 채점 · Tab 힌트 · ↑↓ 이동 · Ctrl+R 이 문제 처음으로 · Ctrl+X 전체 초기화 · Ctrl+C 종료 · 스페이스/b 스크롤"

func (m *Model) layout() {
	w := m.width - listWidth - paneExtra*2
	if w < 20 {
		w = 20
	}
	h := m.height - 2 - 2 // 테두리 2, 상태 줄 2
	if h < 5 {
		h = 5
	}
	m.vp.Width, m.vp.Height = w, h
}

func (m Model) View() string {
	if m.width == 0 {
		return "불러오는 중..."
	}
	var list strings.Builder
	list.WriteString("문제\n\n")
	for i, ex := range m.exs {
		mark := "⬜"
		if m.prog.IsPassed(ex.ID) {
			mark = "✅"
		}
		line := "  " + mark + " " + ex.Title
		switch {
		case i == m.cursor:
			line = selected.Render("> " + mark + " " + ex.Title)
		case m.prog.IsPassed(ex.ID):
			line = passed.Render(line)
		}
		list.WriteString(line + "\n")
	}
	left := pane.Width(listWidth).Height(m.vp.Height).Render(list.String())
	// lipgloss의 Width는 좌우 여백(padding)을 포함하므로 +2를 해야 viewport 줄이 다시 접히지 않는다.
	right := pane.Width(m.vp.Width + 2).Height(m.vp.Height).Render(m.vp.View())
	body := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	status := lipgloss.NewStyle().MaxWidth(m.width).Render(m.status)
	help := keysTxt.MaxWidth(m.width).Render(keyHelp)
	return lipgloss.JoinVertical(lipgloss.Left, body, status, help)
}
