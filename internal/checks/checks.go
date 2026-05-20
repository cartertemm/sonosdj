package checks

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const GoDownloadURL = "https://go.dev/dl/"
const SonosInstallTarget = "github.com/steipete/sonoscli/cmd/sonos@latest"

type Commander interface {
	Run(name string, args ...string) (string, error)
}

type ExecCommander struct{}
type AskFunc func(string) (bool, error)
type OpenBrowserFunc func(string) error

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

func DiscoverSpeakers(cmdr Commander) ([]string, string, error) {
	output, err := cmdr.Run("sonos", "discover")
	if err != nil {
		return nil, "", err
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
		return nil, "", errors.New("no Sonos speakers were discovered")
	}
	return rooms, output, nil
}

func CheckSpotifyAuth(cmdr Commander, room string) error {
	_, err := cmdr.Run("sonos", "smapi", "search", "--name", room, "--service", "Spotify", "--category", "playlists", "test")
	return err
}

func RunSpotifyAuthFlow(cmdr Commander, room string, ask AskFunc, openBrowser OpenBrowserFunc, output io.Writer, input io.Reader) error {
	beginOutput, err := cmdr.Run("sonos", "auth", "smapi", "begin", "--name", room, "--service", "Spotify")
	if err != nil {
		return fmt.Errorf("auth begin failed: %w", err)
	}
	authURL := extractAuthURL(beginOutput)
	if authURL == "" {
		return errors.New("could not find Spotify auth URL in sonos auth output")
	}
	code := extractCode(beginOutput)
	if code == "" {
		return errors.New("could not find Spotify auth code in sonos auth output")
	}
	fmt.Fprintln(output, "Complete Spotify linking in your browser.")
	fmt.Fprintln(output, authURL)
	open, err := ask("Want to open this in your default browser?")
	if err != nil {
		return err
	}
	if open {
		if err := openBrowser(authURL); err != nil {
			fmt.Fprintln(output, "Couldn't open automatically.")
		}
	}
	fmt.Fprint(output, "Press Enter after you finish Spotify login. SonosDJ will complete linking for you...")
	if _, err := bufio.NewReader(input).ReadString('\n'); err != nil {
		return err
	}
	if _, err := cmdr.Run("sonos", "auth", "smapi", "complete", "--service", "Spotify", "--code", code, "--wait", "5m"); err != nil {
		return fmt.Errorf("spotify auth completion failed: %w", err)
	}
	return nil
}

func OpenURLInBrowser(rawURL string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", rawURL)
	case "darwin":
		cmd = exec.Command("open", rawURL)
	default:
		cmd = exec.Command("xdg-open", rawURL)
	}
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	cmd.Stdin = os.Stdin
	return cmd.Start()
}

func extractCode(output string) string {
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "code:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
		if strings.Contains(line, "--code") {
			fields := strings.Fields(line)
			for i := 0; i < len(fields)-1; i++ {
				if fields[i] == "--code" {
					return strings.Trim(fields[i+1], "\"'")
				}
			}
		}
		if strings.Contains(line, "linkCode=") {
			parsed, err := url.Parse(line)
			if err == nil {
				if code := strings.TrimSpace(parsed.Query().Get("linkCode")); code != "" {
					return code
				}
			}
		}
	}
	return ""
}

func extractAuthURL(output string) string {
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
			parsed, err := url.Parse(line)
			if err == nil && parsed.Scheme != "" && parsed.Host != "" {
				return parsed.String()
			}
		}
	}
	return ""
}
