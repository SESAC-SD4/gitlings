package workspace

import (
	"os"
	"path/filepath"
	"testing"

	"gitlings-ko/internal/exercise"
)

func fake(id int, slug string) *exercise.Exercise {
	return &exercise.Exercise{ID: id, Slug: slug, Setup: func(dir string) error {
		return exercise.WriteFiles(dir, map[string]string{"seed.txt": "씨앗\n"})
	}}
}

func TestBaseFollowsHOME(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if Base() != filepath.Join(home, "gitlings-work") {
		t.Fatalf("HOME 아래여야 합니다: %s", Base())
	}
	if Display("01-a") != "~/gitlings-work/01-a" {
		t.Fatalf("표시 경로가 다릅니다: %s", Display("01-a"))
	}
}

func TestPrepareRecreates(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	ex := fake(1, "01-a")
	if Exists("01-a") {
		t.Fatal("처음에는 없어야 합니다")
	}
	if err := Prepare(ex); err != nil {
		t.Fatal(err)
	}
	junk := filepath.Join(Dir("01-a"), "junk.txt")
	if err := os.WriteFile(junk, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Prepare(ex); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(junk); !os.IsNotExist(err) {
		t.Fatal("다시 만들면 학생이 만든 파일은 사라져야 합니다")
	}
	if !exercise.IsFile(Dir("01-a"), "seed.txt") {
		t.Fatal("setup이 실행되어야 합니다")
	}
}

func TestPrepareKeepsFolderForOpenShell(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	ex := fake(1, "01-a")
	if err := Prepare(ex); err != nil {
		t.Fatal(err)
	}
	before, _ := os.Stat(Dir("01-a"))
	if err := os.WriteFile(filepath.Join(Dir("01-a"), ".hidden"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Prepare(ex); err != nil {
		t.Fatal(err)
	}
	after, _ := os.Stat(Dir("01-a"))
	if !os.SameFile(before, after) {
		t.Fatal("다른 터미널이 들어가 있는 폴더 자체를 지우면 안 됩니다. 안의 내용만 비워야 합니다")
	}
	if _, err := os.Stat(filepath.Join(Dir("01-a"), ".hidden")); !os.IsNotExist(err) {
		t.Fatal("숨김 파일까지 비워야 합니다")
	}
}

func TestEnsureAllCreatesOnlyMissing(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	exs := []*exercise.Exercise{fake(1, "01-a"), fake(2, "02-b")}
	if err := Prepare(exs[0]); err != nil {
		t.Fatal(err)
	}
	mine := filepath.Join(Dir("01-a"), "mine.txt")
	if err := os.WriteFile(mine, []byte("풀던 내용"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := EnsureAll(exs); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(mine); err != nil {
		t.Fatal("이미 있는 연습 폴더는 건드리지 않아야 합니다")
	}
	if !exercise.IsFile(Dir("02-b"), "seed.txt") {
		t.Fatal("없는 연습 폴더는 만들어야 합니다")
	}
}

func TestPrepareRefusesEmptySlug(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	if err := Prepare(fake(1, "")); err == nil {
		t.Fatal("slug가 비면 연습 폴더 전체를 지울 위험이 있으므로 거부해야 합니다")
	}
}

func TestDetect(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	exs := []*exercise.Exercise{fake(1, "01-a"), fake(2, "02-b")}
	inside := filepath.Join(Base(), "02-b", "TIL", "deep")
	if err := os.MkdirAll(inside, 0o755); err != nil {
		t.Fatal(err)
	}
	if got := Detect(inside, exs, 1); got == nil || got.ID != 2 {
		t.Fatalf("현재 경로의 문제를 골라야 합니다: %v", got)
	}
	if got := Detect(home, exs, 1); got == nil || got.ID != 1 {
		t.Fatalf("밖에서는 last_started를 써야 합니다: %v", got)
	}
	if got := Detect(home, exs, 0); got != nil {
		t.Fatal("아무것도 시작하지 않았으면 nil")
	}
	link := filepath.Join(t.TempDir(), "link")
	if err := os.Symlink(filepath.Join(Base(), "02-b"), link); err != nil {
		t.Fatal(err)
	}
	if got := Detect(link, exs, 1); got == nil || got.ID != 2 {
		t.Fatalf("심볼릭 링크로 들어와도 같은 문제여야 합니다: %v", got)
	}
}

func TestPrepareRemovesSiblingRepos(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	ex := fake(1, "01-a")
	if err := Prepare(ex); err != nil {
		t.Fatal(err)
	}
	remote, other := Dir("01-a")+".git-remote", Dir("01-ab")
	for _, d := range []string{remote, other} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := Prepare(ex); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(remote); !os.IsNotExist(err) {
		t.Fatal("연습 폴더 옆의 원격 저장소도 지워야 합니다")
	}
	if _, err := os.Stat(other); err != nil {
		t.Fatal("이름이 비슷한 다른 폴더는 건드리지 않아야 합니다")
	}
}
