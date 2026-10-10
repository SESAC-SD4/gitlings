package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Git은 특정 작업 디렉토리에서 git 명령을 실행하는 도우미다.
type Git struct {
	WorkDir string
	// Env은 git 명령 실행 시 추가할 환경 변수 (GIT_DIR 등)
	Env []string
}

// New는 주어진 디렉토리를 기본 작업 디렉토리로 하는 Git을 반환한다.
func New(workDir string) *Git {
	return &Git{WorkDir: workDir}
}

// Run은 git 명령을 실행하고 stdout을 반환한다. stderr도 에러 메시지에 포함한다.
func (g *Git) Run(args ...string) (string, error) {
	return g.RunEnv(nil, args...)
}

// RunEnv는 추가 환경 변수와 함께 git 명령을 실행한다.
func (g *Git) RunEnv(extraEnv []string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = g.WorkDir
	if len(extraEnv) > 0 || len(g.Env) > 0 {
		cmd.Env = append(cmd.Environ(), append(g.Env, extraEnv...)...)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w\n%s", err, strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(stdout.String()), nil
}

// RunWithAuthor는 커밋 작성자를 강제 지정하여 git 명령을 실행한다.
func (g *Git) RunWithAuthor(args ...string) (string, error) {
	allArgs := append([]string{
		"-c", "user.name=gitlings",
		"-c", "user.email=gitlings@example.com",
	}, args...)
	return g.Run(allArgs...)
}

// CheckVersion은 설치된 git이 2.23 이상인지 확인한다.
// 아니면 한국어 안내 메시지와 함께 에러를 반환한다.
func CheckVersion() error {
	out, err := exec.Command("git", "--version").Output()
	if err != nil {
		return fmt.Errorf("git을 찾을 수 없습니다. git을 먼저 설치해 주세요.\n설치: https://git-scm.com/")
	}
	// "git version 2.39.2" 형태
	parts := strings.Fields(string(out))
	if len(parts) < 3 {
		return nil
	}
	ver := parts[2]
	if !versionAtLeast(ver, 2, 23) {
		return fmt.Errorf("git 버전이 너무 낮습니다 (현재: %s, 필요: 2.23 이상).\ngit을 업데이트해 주세요: https://git-scm.com/", ver)
	}
	return nil
}

func versionAtLeast(ver string, major, minor int) bool {
	parts := strings.SplitN(ver, ".", 3)
	if len(parts) < 2 {
		return false
	}
	maj, err1 := strconv.Atoi(parts[0])
	min, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return false
	}
	if maj != major {
		return maj > major
	}
	return min >= minor
}

// FileAtHEAD는 HEAD에 파일이 있는지 확인한다.
func (g *Git) FileAtHEAD(path string) bool {
	_, err := g.Run("cat-file", "-e", "HEAD:"+path)
	return err == nil
}

// CommitCount는 현재 브랜치의 커밋 수를 반환한다.
func (g *Git) CommitCount() (int, error) {
	out, err := g.Run("rev-list", "--count", "HEAD")
	if err != nil {
		return 0, err
	}
	n, err := strconv.Atoi(out)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// IsClean은 작업 트리가 깨끗한지(변경 없음) 확인한다.
func (g *Git) IsClean() bool {
	out, err := g.Run("status", "--porcelain")
	return err == nil && out == ""
}

// HasGitDir은 workDir에 .git이 있는지 확인한다.
func (g *Git) HasGitDir() bool {
	_, err := g.Run("rev-parse", "--git-dir")
	return err == nil
}

// HeadCommit은 HEAD 커밋 해시를 반환한다.
func (g *Git) HeadCommit() (string, error) {
	return g.Run("rev-parse", "HEAD")
}

// HeadMessage는 HEAD 커밋 메시지를 반환한다.
func (g *Git) HeadMessage() (string, error) {
	return g.Run("log", "-1", "--format=%s")
}

// StagedFiles는 인덱스에 staged된 파일 목록을 반환한다.
func (g *Git) StagedFiles() ([]string, error) {
	out, err := g.Run("diff", "--cached", "--name-only")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

// IsTracked는 파일이 HEAD에서 추적되고 있는지 확인한다.
func (g *Git) IsTracked(path string) bool {
	out, err := g.Run("ls-files", "--error-unmatch", path)
	return err == nil && out != ""
}

// IsIgnored는 .gitignore에 의해 파일이 무시되는지 확인한다.
func (g *Git) IsIgnored(path string) bool {
	_, err := g.Run("check-ignore", "-q", path)
	return err == nil
}

// BranchExists는 브랜치가 존재하는지 확인한다.
func (g *Git) BranchExists(name string) bool {
	_, err := g.Run("rev-parse", "--verify", name)
	return err == nil
}

// CommitParentCount는 HEAD 커밋의 부모 수를 반환한다 (병합 커밋 확인용).
func (g *Git) CommitParentCount() (int, error) {
	out, err := g.Run("log", "-1", "--format=%P")
	if err != nil {
		return 0, err
	}
	parents := strings.Fields(out)
	return len(parents), nil
}

// FileContent는 HEAD에서 파일 내용을 읽는다.
func (g *Git) FileContent(path string) (string, error) {
	return g.Run("show", "HEAD:"+path)
}

// LocalMatchesRemote는 로컬 main이 origin/main과 같은 커밋을 가리키는지 확인한다.
func (g *Git) LocalMatchesRemote(branch, remote string) bool {
	local, err := g.Run("rev-parse", branch)
	if err != nil {
		return false
	}
	rem, err := g.Run("rev-parse", remote+"/"+branch)
	if err != nil {
		return false
	}
	return local == rem
}

// CommitContains는 브랜치의 이력에 특정 커밋이 포함되어 있는지 확인한다.
func (g *Git) CommitContains(branch, commit string) bool {
	out, err := g.Run("branch", "--contains", commit, branch)
	return err == nil && strings.TrimSpace(out) != ""
}

// HasConflictMarkers는 파일에 충돌 표시가 있는지 확인한다.
func (g *Git) HasConflictMarkers(path string) (bool, error) {
	out, err := g.Run("diff", "--check")
	if err != nil {
		// diff --check는 충돌 표시 발견 시 비정상 종료
		_ = out
		return true, nil
	}
	return false, nil
}

// StashList는 stash 목록을 반환한다.
func (g *Git) StashList() ([]string, error) {
	out, err := g.Run("stash", "list")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

// CurrentBranch는 현재 체크아웃된 브랜치 이름을 반환한다.
func (g *Git) CurrentBranch() (string, error) {
	return g.Run("rev-parse", "--abbrev-ref", "HEAD")
}
