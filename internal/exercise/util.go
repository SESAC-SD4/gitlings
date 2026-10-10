package exercise

import "strings"

// Lines는 문자열을 줄바꿈 기준으로 분리한다.
func Lines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(strings.TrimRight(s, "\n"), "\n")
}

// Contains는 슬라이스에 문자열이 있는지 확인한다.
func Contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

// ContainsSubstr는 슬라이스의 어느 항목에든 substr이 있는지 확인한다.
func ContainsSubstr(slice []string, substr string) bool {
	for _, v := range slice {
		if strings.Contains(v, substr) {
			return true
		}
	}
	return false
}
