package checks

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

type fakeCommander struct {
	outputs map[string]result
	counts  map[string]int
}

type result struct {
	output string
	err    error
}

func (f fakeCommander) Run(name string, args ...string) (string, error) {
	key := name + " " + strings.Join(args, " ")
	if f.counts != nil {
		f.counts[key]++
		if key == "sonos --version" && f.counts[key] > 1 {
			return "sonos version", nil
		}
	}
	res, ok := f.outputs[key]
	if !ok {
		return "", nil
	}
	return res.output, res.err
}

func TestDiscoverSpeakersParsesRooms(t *testing.T) {
	rawOutput := "Living Room\t192.168.0.10\tRINCON_1\nKitchen\t192.168.0.11\tRINCON_2\n"
	cmdr := fakeCommander{
		outputs: map[string]result{
			"sonos discover": {output: rawOutput},
		},
	}

	rooms, output, err := DiscoverSpeakers(cmdr)
	if err != nil {
		t.Fatalf("DiscoverSpeakers returned error: %v", err)
	}

	if len(rooms) != 2 || rooms[0] != "Living Room" || rooms[1] != "Kitchen" {
		t.Fatalf("unexpected rooms: %#v", rooms)
	}
	if output != rawOutput {
		t.Fatalf("expected raw output %q, got %q", rawOutput, output)
	}
}

func TestExtractCodeFindsAuthCode(t *testing.T) {
	code := extractCode("Open this URL\nCode: ABC123\n")
	if code != "ABC123" {
		t.Fatalf("expected code ABC123, got %q", code)
	}
}

func TestExtractCodeFindsAuthCodeInLinkURL(t *testing.T) {
	code := extractCode("Open this URL and link your account:\n  https://spotify-v5.ws.sonos.com/deviceLink/home?linkCode=OA9ZKON\n")
	if code != "OA9ZKON" {
		t.Fatalf("expected code OA9ZKON, got %q", code)
	}
}

func TestExtractCodeFindsAuthCodeInCompletionCommand(t *testing.T) {
	code := extractCode("Then run:\n  sonos auth smapi complete --service \"Spotify\" --code OA9ZKON --wait 5m\n")
	if code != "OA9ZKON" {
		t.Fatalf("expected code OA9ZKON, got %q", code)
	}
}

func TestRunSpotifyAuthFlowCompletesAfterEnter(t *testing.T) {
	cmdr := fakeCommander{
		outputs: map[string]result{
			"sonos auth smapi begin --name Living Room --service Spotify": {
				output: "Service: Spotify\nOpen this URL and link your account:\n  https://spotify-v5.ws.sonos.com/deviceLink/home?linkCode=OA9ZKON\nThen run:\n  sonos auth smapi complete --service \"Spotify\" --code OA9ZKON --wait 5m\n",
			},
			"sonos auth smapi complete --service Spotify --code OA9ZKON --wait 5m": {output: "linked"},
		},
		counts: map[string]int{},
	}

	output := new(bytes.Buffer)
	input := strings.NewReader("\n")
	ask := func(string) (bool, error) { return false, nil }

	if err := RunSpotifyAuthFlow(cmdr, "Living Room", ask, func(string) error { return nil }, output, input); err != nil {
		t.Fatalf("RunSpotifyAuthFlow returned error: %v", err)
	}

	if cmdr.counts["sonos auth smapi complete --service Spotify --code OA9ZKON --wait 5m"] != 1 {
		t.Fatalf("expected auth completion command to run once, got counts %#v", cmdr.counts)
	}

	rendered := output.String()
	if strings.Contains(rendered, "Then run:") {
		t.Fatalf("expected raw sonos CLI instructions to be hidden, got output %q", rendered)
	}
	if !strings.Contains(rendered, "Complete Spotify linking in your browser.") {
		t.Fatalf("expected custom linking message, got output %q", rendered)
	}
	if !strings.Contains(rendered, "https://spotify-v5.ws.sonos.com/deviceLink/home?linkCode=OA9ZKON") {
		t.Fatalf("expected auth URL in output, got %q", rendered)
	}
}

func TestRunSpotifyAuthFlowOpensBrowserWhenAccepted(t *testing.T) {
	cmdr := fakeCommander{
		outputs: map[string]result{
			"sonos auth smapi begin --name Living Room --service Spotify": {
				output: "Open this URL and link your account:\n  https://spotify-v5.ws.sonos.com/deviceLink/home?linkCode=OA9ZKON\nThen run:\n  sonos auth smapi complete --service \"Spotify\" --code OA9ZKON --wait 5m\n",
			},
			"sonos auth smapi complete --service Spotify --code OA9ZKON --wait 5m": {output: "linked"},
		},
		counts: map[string]int{},
	}

	output := new(bytes.Buffer)
	input := strings.NewReader("\n")
	askCalls := 0
	var opened string

	ask := func(prompt string) (bool, error) {
		askCalls++
		if prompt != "Want to open this in your default browser?" {
			t.Fatalf("unexpected prompt %q", prompt)
		}
		return true, nil
	}

	openBrowser := func(rawURL string) error {
		opened = rawURL
		return nil
	}

	if err := RunSpotifyAuthFlow(cmdr, "Living Room", ask, openBrowser, output, input); err != nil {
		t.Fatalf("RunSpotifyAuthFlow returned error: %v", err)
	}

	if askCalls != 1 {
		t.Fatalf("expected one browser prompt, got %d", askCalls)
	}
	if opened != "https://spotify-v5.ws.sonos.com/deviceLink/home?linkCode=OA9ZKON" {
		t.Fatalf("expected browser to open auth URL, got %q", opened)
	}
}

func TestRunSpotifyAuthFlowPrintsFallbackWhenBrowserOpenFails(t *testing.T) {
	cmdr := fakeCommander{
		outputs: map[string]result{
			"sonos auth smapi begin --name Living Room --service Spotify": {
				output: "Open this URL and link your account:\n  https://spotify-v5.ws.sonos.com/deviceLink/home?linkCode=OA9ZKON\nThen run:\n  sonos auth smapi complete --service \"Spotify\" --code OA9ZKON --wait 5m\n",
			},
			"sonos auth smapi complete --service Spotify --code OA9ZKON --wait 5m": {output: "linked"},
		},
		counts: map[string]int{},
	}

	output := new(bytes.Buffer)
	input := strings.NewReader("\n")

	ask := func(string) (bool, error) { return true, nil }
	openBrowser := func(string) error { return errors.New("open failed") }

	if err := RunSpotifyAuthFlow(cmdr, "Living Room", ask, openBrowser, output, input); err != nil {
		t.Fatalf("RunSpotifyAuthFlow returned error: %v", err)
	}

	if !strings.Contains(output.String(), "Couldn't open automatically.") {
		t.Fatalf("expected browser-open fallback message, got output %q", output.String())
	}
}

func TestEnsureSonosCLIInstallsFromSteipeteRepo(t *testing.T) {
	installKey := "go install " + SonosInstallTarget
	cmdr := fakeCommander{
		outputs: map[string]result{
			"sonos --version": {err: errors.New("missing")},
			"go version":      {output: "go version go1.23.0"},
			installKey:        {output: ""},
		},
		counts: map[string]int{},
	}
	output := new(bytes.Buffer)
	ask := func(string) (bool, error) { return true, nil }
	err := EnsureSonosCLI(cmdr, ask, output, true)
	if err != nil {
		t.Fatalf("EnsureSonosCLI returned error: %v", err)
	}
	if cmdr.counts[installKey] != 1 {
		t.Fatalf("expected install command %q to run once, got counts %#v", installKey, cmdr.counts)
	}
}
