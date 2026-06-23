package agents

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

type fakeLookPath map[string]string

var _ LookPather = fakeLookPath{}

func (f fakeLookPath) LookPath(name string) (string, error) {
	path, ok := f[name]
	if !ok {
		return "", fmt.Errorf("%s not found", name)
	}

	return path, nil
}

func TestDetectAvailableAgentsFindsClaudeAndCodex(t *testing.T) {
	finder := fakeLookPath{
		"claude": filepath.Join("bin", "claude"),
		"codex":  filepath.Join("bin", "codex"),
	}

	available := DetectAvailable(finder, runtime.GOOS)
	if len(available) != 2 {
		t.Fatalf("expected 2 agents, got %d", len(available))
	}

	if available[0].Name != "claude" || available[1].Name != "codex" {
		t.Fatalf("unexpected available agents: %+v", available)
	}
}

func TestDetectAvailableAgentsPrefersCodexCmdOnWindows(t *testing.T) {
	finder := fakeLookPath{
		"codex.cmd": filepath.Join("bin", "codex.cmd"),
	}

	available := DetectAvailable(finder, "windows")
	if len(available) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(available))
	}

	if available[0].Path != filepath.Join("bin", "codex.cmd") {
		t.Fatalf("expected codex.cmd path, got %q", available[0].Path)
	}
}

func TestResolveSelectionUsesExplicitAgent(t *testing.T) {
	available := []Agent{
		{Name: "claude", Path: "claude"},
		{Name: "codex", Path: "codex"},
	}

	selected, needsMenu, err := ResolveSelection(available, "codex")
	if err != nil {
		t.Fatalf("ResolveSelection returned error: %v", err)
	}

	if needsMenu {
		t.Fatal("expected menu to be unnecessary")
	}

	if selected.Name != "codex" {
		t.Fatalf("expected codex, got %q", selected.Name)
	}
}

func TestResolveSelectionRequestsMenuWhenBothAvailable(t *testing.T) {
	available := []Agent{
		{Name: "claude", Path: "claude"},
		{Name: "codex", Path: "codex"},
	}

	selected, needsMenu, err := ResolveSelection(available, "")
	if err != nil {
		t.Fatalf("ResolveSelection returned error: %v", err)
	}

	if !needsMenu {
		t.Fatal("expected menu to be needed")
	}

	if selected.Name != "" {
		t.Fatalf("expected empty selection, got %q", selected.Name)
	}
}
