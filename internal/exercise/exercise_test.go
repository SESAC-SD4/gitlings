package exercise

import "testing"

func TestParseHints(t *testing.T) {
	raw := "첫째\n---\n둘째\r\n---  \r\n셋째 줄 하나\n셋째 줄 둘\n"
	got := ParseHints(raw)
	if len(got) != 3 {
		t.Fatalf("3단계여야 합니다: %q", got)
	}
	if got[0] != "첫째" || got[1] != "둘째" || got[2] != "셋째 줄 하나\n셋째 줄 둘" {
		t.Fatalf("나눈 결과가 다릅니다: %q", got)
	}
}

func TestPassFail(t *testing.T) {
	if !Pass("좋다").Passed || Fail("아직").Passed {
		t.Fatal("Pass는 통과, Fail은 실패여야 합니다")
	}
}
