package exercises

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

func TestRegistry(t *testing.T) {
	all := All()
	if len(all) == 0 {
		t.Fatal("등록된 문제가 없습니다")
	}
	for i, ex := range all {
		if ex.ID != i+1 {
			t.Errorf("번호가 이어지지 않습니다: %d번째 문제의 ID가 %d", i+1, ex.ID)
		}
		if !strings.HasPrefix(ex.Slug, fmt.Sprintf("%02d-", ex.ID)) || !strings.HasPrefix(ex.Title, fmt.Sprintf("%02d", ex.ID)) {
			t.Errorf("slug·제목이 번호로 시작해야 합니다: %s / %s", ex.Slug, ex.Title)
		}
		for _, r := range ex.Slug {
			if r > 127 {
				t.Errorf("slug는 ASCII여야 합니다: %s", ex.Slug)
			}
		}
		if strings.TrimSpace(ex.Prompt) == "" || len(ex.Hints) != 3 {
			t.Errorf("%s: 문제 설명과 힌트 3단계가 있어야 합니다", ex.Slug)
		}
		if ex.Optional != (ex.ID >= 12) || ex.Optional != strings.Contains(ex.Title, "★") {
			t.Errorf("%s: 12~13번만 선택(★)입니다", ex.Slug)
		}
	}
	if ByID(0) != nil || ByID(len(all)+1) != nil {
		t.Error("없는 번호는 nil")
	}
}

func TestRegistryHasThirteen(t *testing.T) {
	if n := len(All()); n != 13 {
		t.Fatalf("문제는 13개여야 합니다(지금 %d)", n)
	}
}

func TestEverySetupFailsBeforeSolving(t *testing.T) {
	for _, ex := range All() {
		dir := prepare(t, ex)
		if ex.Check(dir).Passed {
			t.Errorf("%s: 아무것도 안 했는데 통과합니다", ex.Slug)
		}
	}
}

func TestSetupNamesAreASCII(t *testing.T) {
	for _, ex := range All() {
		dir := prepare(t, ex)
		_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
			for _, r := range d.Name() {
				if r > 127 {
					t.Errorf("%s: setup이 만든 이름이 ASCII가 아닙니다: %s", ex.Slug, p)
				}
			}
			return nil
		})
	}
}
