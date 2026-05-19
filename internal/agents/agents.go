package agents

import (
	"errors"
	"fmt"
)

type Agent struct {
	Name string
	Path string
}

type LookPather interface {
	LookPath(name string) (string, error)
}

func DetectAvailable(finder LookPather, goos string) []Agent {
	available := make([]Agent, 0, 2)

	if path, err := finder.LookPath("claude"); err == nil {
		available = append(available, Agent{Name: "claude", Path: path})
	}

	codexNames := []string{"codex"}
	if goos == "windows" {
		codexNames = []string{"codex.cmd", "codex"}
	}

	for _, name := range codexNames {
		if path, err := finder.LookPath(name); err == nil {
			available = append(available, Agent{Name: "codex", Path: path})
			break
		}
	}

	return available
}

func ResolveSelection(available []Agent, explicit string) (Agent, bool, error) {
	if len(available) == 0 {
		return Agent{}, false, errors.New("neither claude nor codex was found on PATH")
	}

	if explicit != "" {
		for _, agent := range available {
			if agent.Name == explicit {
				return agent, false, nil
			}
		}

		return Agent{}, false, fmt.Errorf("%s was requested but is not available on PATH", explicit)
	}

	if len(available) == 1 {
		return available[0], false, nil
	}

	return Agent{}, true, nil
}
