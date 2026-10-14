package git

import "os"

// InitMain은 workDir에 main 브랜치로 git 저장소를 초기화한다.
// git 2.23+에서 동작하며, init.defaultBranch 전역 설정과 무관하다.
func InitMain(workDir string) (*Git, error) {
	g := New(workDir)
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return nil, err
	}
	if _, err := g.Run("init"); err != nil {
		return nil, err
	}
	// 커밋 전 HEAD가 가리키는 브랜치를 main으로 강제 설정
	if _, err := g.Run("symbolic-ref", "HEAD", "refs/heads/main"); err != nil {
		return nil, err
	}
	return g, nil
}

// SetupUserConfig는 로컬 저장소에 커밋 작성자를 설정한다.
// 학생의 전역 설정 여부와 무관하게 동작한다.
func (g *Git) SetupUserConfig() error {
	if _, err := g.Run("config", "user.name", "gitlings"); err != nil {
		return err
	}
	if _, err := g.Run("config", "user.email", "gitlings@example.com"); err != nil {
		return err
	}
	return nil
}

// AddExclude는 .git/info/exclude에 패턴을 추가한다 (이미 있으면 무시).
func (g *Git) AddExclude(pattern string) error {
	excludePath := g.WorkDir + "/.git/info/exclude"

	data, err := os.ReadFile(excludePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	content := string(data)
	for _, line := range splitLines(content) {
		if line == pattern {
			return nil // 이미 있음
		}
	}

	f, err := os.OpenFile(excludePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString("\n" + pattern + "\n")
	return err
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	result := []string{}
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		result = append(result, s[start:])
	}
	return result
}

// InitBareMain은 가짜 원격으로 쓸 bare 저장소를 만들고 HEAD를 main으로 고정한다.
// 학생의 init.defaultBranch가 master면 원격 HEAD가 없는 master를 가리켜,
// 그 원격을 clone한 저장소가 main을 받지 못한다.
func InitBareMain(dir string) (*Git, error) {
	g := New(dir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	if _, err := g.Run("init", "--bare"); err != nil {
		return nil, err
	}
	if _, err := g.Run("symbolic-ref", "HEAD", "refs/heads/main"); err != nil {
		return nil, err
	}
	return g, nil
}
