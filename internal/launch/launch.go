package launch

import (
	"os"
	"os/exec"

	"github.com/cartertemm/sonosdj/internal/agents"
)

func BuildCommand(agent agents.Agent, prompt string) *exec.Cmd {
	var args []string
	if agent.Name == "claude" {
		args = []string{agent.Path, "--allowedTools", "Bash(sonos:*)", "Bash(wait:*)", "--", prompt}
	} else {
		args = []string{agent.Path, prompt}
	}
	cmd := &exec.Cmd{
		Path: agent.Path,
		Args: args,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
