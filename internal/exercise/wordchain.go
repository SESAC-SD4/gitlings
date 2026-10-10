package exercise

import (
	"strings"
	"unicode/utf8"
)

// WordChainWord는 단어장의 한 항목을 나타낸다.
type WordChainWord struct {
	Number int
	Word   string
}

// ParseWordsMd는 words.md 내용을 파싱해 단어 목록을 반환한다.
// 형식: "N. 단어" (숫자와 점 뒤에 단어)
func ParseWordsMd(content string) []WordChainWord {
	var words []WordChainWord
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// "1. 기차" 형태
		dotIdx := strings.Index(line, ".")
		if dotIdx < 0 {
			continue
		}
		word := strings.TrimSpace(line[dotIdx+1:])
		if word == "" {
			continue
		}
		var num int
		numStr := strings.TrimSpace(line[:dotIdx])
		for _, c := range numStr {
			if c >= '0' && c <= '9' {
				num = num*10 + int(c-'0')
			}
		}
		if num > 0 {
			words = append(words, WordChainWord{Number: num, Word: word})
		}
	}
	return words
}

// ValidateWordChain은 단어 목록이 끝말잇기 규칙에 맞는지 확인한다.
// 규칙: 다음 단어의 첫 글자 == 앞 단어의 끝 글자, 두 글자 이상.
// 유효하면 nil을 반환한다.
func ValidateWordChain(words []WordChainWord) error {
	for i := 1; i < len(words); i++ {
		prev := []rune(words[i-1].Word)
		curr := []rune(words[i].Word)
		if len(curr) < 2 {
			return &WordChainError{
				Index:   i,
				Message: "단어가 너무 짧습니다 (두 글자 이상이어야 합니다): " + words[i].Word,
			}
		}
		if len(prev) == 0 {
			continue
		}
		lastRune := prev[len(prev)-1]
		firstRune := curr[0]
		if lastRune != firstRune {
			return &WordChainError{
				Index: i,
				Message: "끝말잇기 규칙 위반: \"" + words[i-1].Word + "\"의 끝 글자 '" +
					string(lastRune) + "'와 \"" + words[i].Word + "\"의 첫 글자 '" +
					string(firstRune) + "'가 다릅니다",
			}
		}
	}
	return nil
}

// WordChainError는 끝말잇기 규칙 위반 에러다.
type WordChainError struct {
	Index   int
	Message string
}

func (e *WordChainError) Error() string {
	return e.Message
}

// LastRune은 문자열의 마지막 룬을 반환한다.
func LastRune(s string) rune {
	if s == "" {
		return 0
	}
	r, _ := utf8.DecodeLastRuneInString(s)
	return r
}

// FirstRune은 문자열의 첫 룬을 반환한다.
func FirstRune(s string) rune {
	if s == "" {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(s)
	return r
}
