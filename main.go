package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/cartertemm/sonosdj/internal/agents"
	"github.com/cartertemm/sonosdj/internal/app"
	"github.com/cartertemm/sonosdj/internal/checks"
	"github.com/cartertemm/sonosdj/internal/launch"
	"github.com/cartertemm/sonosdj/internal/prompt"
	"github.com/cartertemm/sonosdj/internal/ui"
)

type execLookPath struct{}

func (execLookPath) LookPath(name string) (string, error) {
	return exec.LookPath(name)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := app.ParseArgs(os.Args[1:])
	if err != nil {
		return err
	}
	keyReader := ui.NewTerminalKeyReader(os.Stdin)
	available := agents.DetectAvailable(execLookPath{}, runtime.GOOS)
	selected, needsMenu, err := agents.ResolveSelection(available, cfg.Agent)
	if err != nil {
		return err
	}
	if needsMenu {
		selected, err = ui.PromptMenu(keyReader, os.Stdout, available)
		if err != nil {
			return err
		}
	}
	if cfg.Verbose {
		fmt.Fprintf(os.Stdout, "Using agent: %s\n", selected.Name)
	}
	cmdr := checks.ExecCommander{}
	askInstall := func(message string) (bool, error) {
		return ui.PromptYesNo(keyReader, os.Stdout, message)
	}
	if err := checks.EnsureSonosCLI(cmdr, askInstall, os.Stdout, cfg.Verbose); err != nil {
		return err
	}
	rooms, err := checks.DiscoverSpeakers(cmdr)
	if err != nil {
		return err
	}
	if cfg.Verbose {
		fmt.Fprintf(os.Stdout, "Discovered Sonos rooms: %s\n", strings.Join(rooms, ", "))
	}
	authRoom := cfg.Room
	if authRoom == "" {
		authRoom = rooms[0]
	}
	if err := checks.CheckSpotifyAuth(cmdr, authRoom); err != nil {
		if flowErr := checks.RunSpotifyAuthFlow(cmdr, authRoom, askInstall, checks.OpenURLInBrowser, os.Stdout, os.Stdin); flowErr != nil {
			return flowErr
		}
	}
	if cfg.Verbose {
		fmt.Fprintf(os.Stdout, "Spotify auth check passed for %s.\n", authRoom)
	}
	cmd := launch.BuildCommand(selected, prompt.Build(cfg.Permissiveness, cfg.Room))
	return cmd.Run()
}
