package prompt

import (
	"strings"
	"testing"
)

func TestBuildIncludesPermissiveness(t *testing.T) {
	text := Build("high", "", "")

	if !strings.Contains(text, "Default permissiveness for this session: high.") {
		t.Fatalf("expected prompt to include permissiveness, got:\n%s", text)
	}
}

func TestBuildIncludesWelcomeExamples(t *testing.T) {
	text := Build("medium", "", "")

	if !strings.Contains(text, `You can say things like "play something hazy and nocturnal"`) {
		t.Fatalf("expected prompt welcome examples, got:\n%s", text)
	}
}

func TestBuildIncludesRawDiscoverOutput(t *testing.T) {
	rawOutput := "Living Room\t192.168.0.10\tRINCON_1\nKitchen\t192.168.0.11\tRINCON_2\n"
	text := Build("medium", "", rawOutput)

	if !strings.Contains(text, "Current sonos discover output:\n"+rawOutput) {
		t.Fatalf("expected prompt to include raw discover output, got:\n%s", text)
	}
}

func TestBuildOnlySuggestsRediscoveryForDeviceErrors(t *testing.T) {
	text := Build("medium", "", "Living Room\t192.168.0.10\tRINCON_1\n")

	if strings.Contains(text, "Discover speakers with: sonos discover") {
		t.Fatalf("expected prompt not to instruct unconditional rediscovery, got:\n%s", text)
	}
	if !strings.Contains(text, "Run `sonos discover` again only if") {
		t.Fatalf("expected prompt to mention conditional rediscovery, got:\n%s", text)
	}
}
