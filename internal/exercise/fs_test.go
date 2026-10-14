package exercise

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteFiles(t *testing.T) {
	dir := t.TempDir()
	err := WriteFiles(dir, map[string]string{"a/b/c.md": "내용\n", "empty/": ""})
	if err != nil {
		t.Fatal(err)
	}
	if !IsFile(dir, "a/b/c.md") || !IsDir(dir, "empty") {
		t.Fatal("파일과 빈 폴더가 생겨야 합니다")
	}
}

func TestExistsIsCaseExact(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "til"), 0o755); err != nil {
		t.Fatal(err)
	}
	if Exists(dir, "TIL") {
		t.Fatal("대소문자가 다르면 없는 것으로 봐야 합니다")
	}
	if !Exists(dir, "til") || !IsDir(dir, "til") || IsFile(dir, "til") {
		t.Fatal("정확한 이름은 있어야 합니다")
	}
}

func TestReadNormalizesCRLF(t *testing.T) {
	dir := t.TempDir()
	if err := WriteFiles(dir, map[string]string{"x.md": "가\r\n나\r\n"}); err != nil {
		t.Fatal(err)
	}
	got, ok := Read(dir, "x.md")
	if !ok || got != "가\n나\n" {
		t.Fatalf("CRLF를 LF로 바꿔야 합니다: %q %v", got, ok)
	}
	if _, ok := Read(dir, "없음.md"); ok {
		t.Fatal("없는 파일은 ok=false")
	}
}

func TestSameIgnoresTrailingNewline(t *testing.T) {
	if !Same("가\n나\n", "가\r\n나") {
		t.Fatal("끝 줄바꿈과 CRLF는 무시해야 합니다")
	}
	if Same("가\n나\n", "가\n") {
		t.Fatal("내용이 다르면 달라야 합니다")
	}
}

func TestNonEmptyLines(t *testing.T) {
	got := NonEmptyLines("\n가\n  \n나\n")
	if len(got) != 2 || got[0] != "가" || got[1] != "나" {
		t.Fatalf("빈 줄을 빼야 합니다: %q", got)
	}
}

func TestList(t *testing.T) {
	dir := t.TempDir()
	if err := WriteFiles(dir, map[string]string{"p/b.md": "", "p/a.md": ""}); err != nil {
		t.Fatal(err)
	}
	got := List(dir, "p")
	if len(got) != 2 || got[0] != "a.md" || got[1] != "b.md" {
		t.Fatalf("정렬된 이름이어야 합니다: %q", got)
	}
	if List(dir, "없음") != nil {
		t.Fatal("없는 폴더는 nil")
	}
}

func TestListSkipsHiddenFiles(t *testing.T) {
	dir := t.TempDir()
	if err := WriteFiles(dir, map[string]string{"p/a.md": "", "p/.DS_Store": ""}); err != nil {
		t.Fatal(err)
	}
	if got := List(dir, "p"); len(got) != 1 || got[0] != "a.md" {
		t.Fatalf("Finder가 만드는 .DS_Store 같은 숨김 파일은 빼야 합니다: %q", got)
	}
}

func TestMissingSuggestsPwd(t *testing.T) {
	if r := Missing(t.TempDir(), "TIL"); !strings.Contains(r.Message, "pwd") {
		t.Fatalf("엉뚱한 폴더에서 작업했을 수 있으니 pwd 확인을 안내해야 합니다: %s", r.Message)
	}
}

func TestMissingPointsOutCase(t *testing.T) {
	dir := t.TempDir()
	if err := WriteFiles(dir, map[string]string{"til/": ""}); err != nil {
		t.Fatal(err)
	}
	r := Missing(dir, "TIL/10-12.md")
	if r.Passed || !strings.Contains(r.Message, "대소문자") || !strings.Contains(r.Message, "til") {
		t.Fatalf("대소문자 안내가 있어야 합니다: %s", r.Message)
	}
	r = Missing(t.TempDir(), "TIL/10-12.md")
	if !strings.Contains(r.Message, "아직 없습니다: `TIL/10-12.md`") {
		t.Fatalf("없다는 안내가 있어야 합니다: %s", r.Message)
	}
}
