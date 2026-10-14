package exercises

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// 문제마다 "정답 명령 줄"로 볼 수 있는 조각. 힌트에 이것이 보이면 답을 준 것이다.
var answerLines = map[int][]string{
	1:  {"git add words.md", "git commit -m"},
	2:  {"git add words.md", "git add [파일"},
	3:  {"restore --staged"},
	4:  {"git restore words.md", "git restore [파일"},
	5:  {"--amend"},
	6:  {"rm --cached", "> .gitignore"},
	7:  {"--soft", "HEAD~1"},
	8:  {"git revert HEAD", "--no-edit"},
	9:  {"switch -c", "git merge new-word"},
	10: {"git merge partner"},
	11: {"pull.rebase false", "--no-rebase", "git push origin main"},
	12: {"reset --hard [해시]", "HEAD@{1}"},
	13: {"git stash pop"},
}

// 답이 되는 git 명령을 인자까지 완성해서 쓴 꼴. `git add`처럼 명령 이름만 쓴 것과
// `git status`·`git log`처럼 상태를 보는 명령은 괜찮다.
var fullCommand = regexp.MustCompile("`git (init|add|commit|restore|reset|revert|push|pull|switch|checkout|merge|rm|stash|config|branch) [^`]*`")

var chapterRef = regexp.MustCompile(`교안|\d+절`)

func TestHintsDoNotGiveAnswers(t *testing.T) {
	for _, ex := range All() {
		for i, h := range ex.Hints {
			for _, ans := range answerLines[ex.ID] {
				if strings.Contains(h, ans) {
					t.Errorf("%s 힌트 %d단계에 정답 명령 %q가 들어 있습니다", ex.Slug, i+1, ans)
				}
			}
			// push는 수업에서 항상 `git push origin main`으로 쓰는 고정 형태(D35)라 답이 아니다.
			// 11번에서는 answerLines가 따로 막는다.
			if m := fullCommand.FindString(strings.ReplaceAll(h, "`git push origin main`", "")); m != "" {
				t.Errorf("%s 힌트 %d단계에 완성된 명령 %s가 들어 있습니다", ex.Slug, i+1, m)
			}
			if strings.Contains(h, "```") {
				t.Errorf("%s 힌트 %d단계에 명령 코드 블록이 있습니다(그대로 따라 치게 됩니다)", ex.Slug, i+1)
			}
			if chapterRef.MatchString(h) {
				t.Errorf("%s 힌트 %d단계가 교안 절을 가리킵니다(교안은 바뀔 수 있음)", ex.Slug, i+1)
			}
		}
	}
}

func TestMessagesDoNotGiveCommands(t *testing.T) {
	files, _ := filepath.Glob("e*.go")
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if m := fullCommand.FindString(string(b)); m != "" {
			t.Errorf("%s의 채점 메시지에 완성된 명령 %s가 들어 있습니다", f, m)
		}
		for _, ans := range answerLines[fileID(f)] {
			if strings.Contains(string(b), ans) {
				t.Errorf("%s의 채점 메시지에 정답 조각 %q가 들어 있습니다", f, ans)
			}
		}
	}
}

// "e07_reset_soft.go" → 7
func fileID(name string) int {
	n := 0
	for _, r := range strings.TrimPrefix(name, "e")[:2] {
		n = n*10 + int(r-'0')
	}
	return n
}
