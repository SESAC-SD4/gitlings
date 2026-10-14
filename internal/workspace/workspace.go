package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gitlings-ko/internal/exercise"
)

// Git Bash에서 HOME을 따로 정한 학생도 `cd ~/gitlings-work`와 같은 곳을 보도록 HOME을 먼저 쓴다.
func Base() string {
	home := os.Getenv("HOME")
	if home == "" {
		home, _ = os.UserHomeDir()
	}
	return filepath.Join(home, "gitlings-work")
}

func Dir(slug string) string     { return filepath.Join(Base(), slug) }
func Display(slug string) string { return "~/gitlings-work/" + slug }

func Exists(slug string) bool {
	info, err := os.Stat(Dir(slug))
	return err == nil && info.IsDir()
}

func Prepare(ex *exercise.Exercise) error {
	if ex.Slug == "" {
		return errors.New("문제 slug가 비어 있습니다")
	}
	d := Dir(ex.Slug)
	// 원격 저장소(`.git-remote`)·짝의 저장소(`.partner`)처럼 연습 폴더 옆에 만든 폴더도 함께 지운다.
	siblings, err := filepath.Glob(d + ".*")
	if err != nil {
		return err
	}
	for _, s := range siblings {
		if err := os.RemoveAll(s); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(d, 0o755); err != nil {
		return err
	}
	// 학생의 다른 터미널이 이 폴더 안에 있을 수 있으므로 폴더 자체는 남기고 안만 비운다.
	entries, err := os.ReadDir(d)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if err := os.RemoveAll(filepath.Join(d, e.Name())); err != nil {
			return err
		}
	}
	return ex.Setup(d)
}

// EnsureAll은 아직 없는 연습 폴더만 만든다. 이미 있는 폴더(학생이 풀던 것)는 건드리지 않는다.
func EnsureAll(exs []*exercise.Exercise) error {
	for _, ex := range exs {
		if Exists(ex.Slug) {
			continue
		}
		if err := Prepare(ex); err != nil {
			return err
		}
	}
	return nil
}

func real(p string) string {
	if r, err := filepath.EvalSymlinks(p); err == nil {
		return r
	}
	return p
}

func Detect(cwd string, exs []*exercise.Exercise, lastStarted int) *exercise.Exercise {
	if rel, err := filepath.Rel(real(Base()), real(cwd)); err == nil && rel != "." && !strings.HasPrefix(rel, "..") {
		first := strings.Split(filepath.ToSlash(rel), "/")[0]
		for _, ex := range exs {
			if ex.Slug == first {
				return ex
			}
		}
	}
	for _, ex := range exs {
		if ex.ID == lastStarted {
			return ex
		}
	}
	return nil
}
