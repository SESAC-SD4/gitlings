package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"gitlings-ko/internal/exercise"
	"gitlings-ko/internal/progress"
	"gitlings-ko/internal/workspace"
)

func GlamourRenderer(dark bool) Renderer {
	style := "light"
	if dark {
		style = "dark"
	}
	return func(md string, width int) string {
		if width < 20 {
			width = 20
		}
		r, err := glamour.NewTermRenderer(glamour.WithStandardStyle(style), glamour.WithWordWrap(width-2))
		if err != nil {
			return md
		}
		out, err := r.Render(md)
		if err != nil {
			return md
		}
		return out
	}
}

// 배경 밝기는 프로그램을 띄우기 전에 한 번만 묻는다. 실행 중에 물으면 응답 문자가 키 입력으로 섞일 수 있다.
func Run(exs []*exercise.Exercise, prog *progress.Progress) error {
	if err := workspace.EnsureAll(exs); err != nil {
		return err
	}
	render := GlamourRenderer(lipgloss.HasDarkBackground())
	_, err := tea.NewProgram(New(exs, prog, render), tea.WithAltScreen()).Run()
	return err
}
