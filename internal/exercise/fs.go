package exercise

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func abs(dir, rel string) string { return filepath.Join(dir, filepath.FromSlash(rel)) }

func parts(rel string) []string {
	var out []string
	for _, p := range strings.Split(filepath.ToSlash(rel), "/") {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func WriteFiles(dir string, files map[string]string) error {
	for rel, body := range files {
		p := abs(dir, rel)
		if strings.HasSuffix(rel, "/") {
			if err := os.MkdirAll(p, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// mac·Windows 파일시스템은 대소문자를 구별하지 않아 os.Stat만으로는 `til`과 `TIL`을 가를 수 없다.
func Exists(dir, rel string) bool {
	cur := dir
	for _, p := range parts(rel) {
		entries, err := os.ReadDir(cur)
		if err != nil {
			return false
		}
		found := false
		for _, e := range entries {
			if e.Name() == p {
				found = true
				break
			}
		}
		if !found {
			return false
		}
		cur = filepath.Join(cur, p)
	}
	return true
}

func IsDir(dir, rel string) bool {
	if !Exists(dir, rel) {
		return false
	}
	info, err := os.Stat(abs(dir, rel))
	return err == nil && info.IsDir()
}

func IsFile(dir, rel string) bool {
	if !Exists(dir, rel) {
		return false
	}
	info, err := os.Stat(abs(dir, rel))
	return err == nil && info.Mode().IsRegular()
}

func Read(dir, rel string) (string, bool) {
	if !IsFile(dir, rel) {
		return "", false
	}
	b, err := os.ReadFile(abs(dir, rel))
	if err != nil {
		return "", false
	}
	return strings.ReplaceAll(string(b), "\r\n", "\n"), true
}

func norm(s string) string { return strings.TrimSpace(strings.ReplaceAll(s, "\r\n", "\n")) }

func Same(a, b string) bool { return norm(a) == norm(b) }

func NonEmptyLines(s string) []string {
	var out []string
	for _, line := range strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n") {
		if t := strings.TrimSpace(line); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func List(dir, rel string) []string {
	if !IsDir(dir, rel) {
		return nil
	}
	entries, err := os.ReadDir(abs(dir, rel))
	if err != nil {
		return nil
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		// Finder가 만드는 .DS_Store 같은 숨김 파일 때문에 정답이 실패하지 않도록 뺀다.
		if !strings.HasPrefix(e.Name(), ".") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names
}

func Missing(dir, rel string) Result {
	cur, done := dir, ""
	for _, p := range parts(rel) {
		entries, err := os.ReadDir(cur)
		if err != nil {
			break
		}
		exact := false
		for _, e := range entries {
			if e.Name() == p {
				exact = true
				break
			}
		}
		if exact {
			cur, done = filepath.Join(cur, p), done+p+"/"
			continue
		}
		for _, e := range entries {
			if strings.EqualFold(e.Name(), p) {
				return Fail(fmt.Sprintf("비슷한 이름이 있습니다: `%s%s` — 대소문자까지 똑같이 `%s%s`로 지어야 합니다", done, e.Name(), done, p))
			}
		}
		break
	}
	return Fail(fmt.Sprintf("아직 없습니다: `%s`. 다른 터미널에서 `pwd`로 연습 폴더 안에 있는지 확인하세요", rel))
}
