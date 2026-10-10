package exercise

import (
	"gitlings-ko/internal/git"
	"strings"
)

// Result는 check 결과다.
type Result struct {
	Passed  bool
	Message string // 실패 메시지 또는 통과 메시지
}

// Exercise는 하나의 연습 문제를 나타낸다.
type Exercise struct {
	ID       int
	Slug     string // "01-first-commit"
	Name     string // "01 첫 사진"
	Optional bool
	Desc     string
	Hints    []string // 3단계 힌트
	Setup    func(workDir string, g *git.Git) error
	Check    func(workDir string, g *git.Git) (Result, error)
}

// WorkDir는 이 문제의 작업 디렉토리 절대 경로를 반환한다.
func (e *Exercise) WorkDir(baseDir string) string {
	return baseDir + "/" + e.Slug
}

// ParseHints는 빈 줄로 구분된 힌트 텍스트를 3단계 슬라이스로 파싱한다.
func ParseHints(raw string) []string {
	parts := strings.Split(raw, "\n\n")
	var hints []string
	for _, p := range parts {
		h := strings.TrimSpace(p)
		if h != "" {
			hints = append(hints, h)
		}
	}
	return hints
}
