package launch

import (
	"testing"

	"github.com/cartertemm/sonosdj/internal/agents"
)

func TestBuildCommandForClaude(t *testing.T) {
	agent := agents.Agent{Name: "claude", Path: "claude"}
	cmd := BuildCommand(agent, "hello")

	if got := cmd.Path; got != "claude" {
		t.Fatalf("expected claude path, got %q", got)
	}

	if len(cmd.Args) != 6 || cmd.Args[1] != "--allowedTools" || cmd.Args[2] != "Bash(sonos:*)" || cmd.Args[3] != "Bash(wait:*)" || cmd.Args[4] != "--" || cmd.Args[5] != "hello" {
		t.Fatalf("unexpected args: %#v", cmd.Args)
	}
}

func TestBuildCommandForCodex(t *testing.T) {
	agent := agents.Agent{Name: "codex", Path: "C:\\bin\\codex.cmd"}
	cmd := BuildCommand(agent, "hello")

	if got := cmd.Path; got != "C:\\bin\\codex.cmd" {
		t.Fatalf("expected codex path, got %q", got)
	}

	if len(cmd.Args) != 2 || cmd.Args[1] != "hello" {
		t.Fatalf("unexpected args: %#v", cmd.Args)
	}
}
