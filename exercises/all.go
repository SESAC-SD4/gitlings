package exercises

import (
	"embed"
	"fmt"
	"sort"

	"gitlings-ko/internal/exercise"
)

//go:embed */prompt.md */hint.md
var files embed.FS

var registry []*exercise.Exercise

func register(id int, slug, title string, optional bool, setup func(string) error, check func(string) exercise.Result) {
	prompt, err := files.ReadFile(slug + "/prompt.md")
	if err != nil {
		panic(fmt.Sprintf("%s/prompt.md를 읽지 못했습니다: %v", slug, err))
	}
	raw, err := files.ReadFile(slug + "/hint.md")
	if err != nil {
		panic(fmt.Sprintf("%s/hint.md를 읽지 못했습니다: %v", slug, err))
	}
	hints := exercise.ParseHints(string(raw))
	if len(hints) != 3 {
		panic(fmt.Sprintf("%s/hint.md의 힌트는 3단계여야 합니다(지금 %d)", slug, len(hints)))
	}
	registry = append(registry, &exercise.Exercise{
		ID: id, Slug: slug, Title: title, Optional: optional,
		Prompt: string(prompt), Hints: hints, Setup: setup, Check: check,
	})
}

func All() []*exercise.Exercise {
	out := append([]*exercise.Exercise(nil), registry...)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func ByID(id int) *exercise.Exercise {
	for _, ex := range registry {
		if ex.ID == id {
			return ex
		}
	}
	return nil
}
