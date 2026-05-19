package checks

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

const GoDownloadURL = "https://go.dev/dl/"
const SonosInstallTarget = "github.com/steipete/sonoscli/cmd/sonos@latest"

type Commander interface {
	Run(name string, args ...string) (string, error)
}

type ExecCommander struct{}

func (ExecCommander) Run(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func EnsureSonosCLI(cmdr Commander, ask func(string) (bool, error), output io.Writer, verbose bool) error {
	if _, err := cmdr.Run("sonos", "--version"); err == nil {
		if verbose {
			fmt.Fprintln(output, "Found sonos backend.")
		}
		return nil
	}
	if _, err := cmdr.Run("go", "version"); err != nil {
		return fmt.Errorf("the SonosCLI backend requires golang to be available and present on your path. Get it from %q", GoDownloadURL)
	}
	install, err := ask("Want to install the SonosCLI backend?")
	if err != nil {
		return err
	}
	if !install {
		return errors.New("sonoscli backend not installed")
	}
	if verbose {
		fmt.Fprintf(output, "Installing sonos backend from %s...\n", SonosInstallTarget)
	}
	if out, err := cmdr.Run("go", "install", SonosInstallTarget); err != nil {
		return fmt.Errorf("installing sonoscli failed: %w\n%s", err, strings.TrimSpace(out))
	}
	if _, err := cmdr.Run("sonos", "--version"); err != nil {
		return fmt.Errorf("sonos is still unavailable after install: %w", err)
	}
	return nil
}

func DiscoverSpeakers(cmdr Commander) ([]string, error) {
	output, err := cmdr.Run("sonos", "discover")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	rooms := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		room := strings.TrimSpace(fields[0])
		if room == "" {
			continue
		}
		rooms = append(rooms, room)
	}

	if len(rooms) == 0 {
		return nil, errors.New("no Sonos speakers were discovered")
	}

	return rooms, nil
}

func CheckSpotifyAuth(cmdr Commander, room string) error {
	_, err := cmdr.Run("sonos", "smapi", "search", "--name", room, "--service", "Spotify", "--category", "playlists", "test")
	return err
}

func RunSpotifyAuthFlow(cmdr Commander, room string, output io.Writer, input io.Reader) error {
	beginOutput, err := cmdr.Run("sonos", "auth", "smapi", "begin", "--name", room, "--service", "Spotify")
	if err != nil {
		return fmt.Errorf("auth begin failed: %w", err)
	}
	fmt.Fprintln(output, strings.TrimSpace(beginOutput))
	fmt.Fprint(output, "Press Enter after you finish Spotify login...")
	if _, err := bufio.NewReader(input).ReadString('\n'); err != nil {
		return err
	}
	code := extractCode(beginOutput)
	if code == "" {
		return errors.New("could not find Spotify auth code in sonos auth output")
	}
	if _, err := cmdr.Run("sonos", "auth", "smapi", "complete", "--service", "Spotify", "--code", code, "--wait", "5m"); err != nil {
		return fmt.Errorf("spotify auth completion failed: %w", err)
	}
	return nil
}

func extractCode(output string) string {
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "code:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Code:"))
		}
	}

	return ""
}
