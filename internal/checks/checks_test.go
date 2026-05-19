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
	cmdr := fakeCommander{
		outputs: map[string]result{
			"sonos discover": {output: "Living Room\t192.168.0.10\tRINCON_1\nKitchen\t192.168.0.11\tRINCON_2\n"},
		},
	}

	rooms, err := DiscoverSpeakers(cmdr)
	if err != nil {
		t.Fatalf("DiscoverSpeakers returned error: %v", err)
	}

	if len(rooms) != 2 || rooms[0] != "Living Room" || rooms[1] != "Kitchen" {
		t.Fatalf("unexpected rooms: %#v", rooms)
	}
}

func TestExtractCodeFindsAuthCode(t *testing.T) {
	code := extractCode("Open this URL\nCode: ABC123\n")
	if code != "ABC123" {
		t.Fatalf("expected code ABC123, got %q", code)
	}
}

func TestEnsureSonosCLIInstallsFromSteipeteRepo(t *testing.T) {
	cmdr := fakeCommander{
		outputs: map[string]result{
			"sonos --version": {err: errors.New("missing")},
			"go version":      {output: "go version go1.23.0"},
			"go install github.com/steipete/sonoscli@latest": {output: ""},
		},
		counts: map[string]int{},
	}
	output := new(bytes.Buffer)
	ask := func(string) (bool, error) { return true, nil }
	err := EnsureSonosCLI(cmdr, ask, output, true)
	if err != nil {
		t.Fatalf("EnsureSonosCLI returned error: %v", err)
	}
}
