package progress

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Progress struct {
	Passed      []int `json:"passed"`
	LastStarted int   `json:"last_started"`
}

func dir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "gitlings-work")
}

func file() string {
	return filepath.Join(dir(), ".progress.json")
}

func Load() (*Progress, error) {
	data, err := os.ReadFile(file())
	if os.IsNotExist(err) {
		return &Progress{}, nil
	}
	if err != nil {
		return nil, err
	}
	var p Progress
	if err := json.Unmarshal(data, &p); err != nil {
		return &Progress{}, nil
	}
	return &p, nil
}

func (p *Progress) Save() error {
	if err := os.MkdirAll(dir(), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file(), data, 0644)
}

func (p *Progress) MarkPassed(id int) {
	for _, v := range p.Passed {
		if v == id {
			return
		}
	}
	p.Passed = append(p.Passed, id)
}

func (p *Progress) IsPassed(id int) bool {
	for _, v := range p.Passed {
		if v == id {
			return true
		}
	}
	return false
}
