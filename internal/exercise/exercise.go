package exercise

import (
	"regexp"
	"strings"
)

type Result struct {
	Passed  bool
	Message string
}

func Pass(msg string) Result { return Result{Passed: true, Message: msg} }
func Fail(msg string) Result { return Result{Message: msg} }

type Exercise struct {
	ID       int
	Slug     string // "01-make-til"
	Title    string // "01 TIL 폴더"
	Optional bool
	Prompt   string
	Hints    []string
	Setup    func(dir string) error
	Check    func(dir string) Result
}

var hintSep = regexp.MustCompile(`(?m)^---[ \t]*$`)

func ParseHints(raw string) []string {
	var out []string
	for _, part := range hintSep.Split(strings.ReplaceAll(raw, "\r\n", "\n"), -1) {
		if s := strings.TrimSpace(part); s != "" {
			out = append(out, s)
		}
	}
	return out
}
