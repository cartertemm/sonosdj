package prompt

import (
	"strings"
	"testing"
)

func TestBuildIncludesPermissiveness(t *testing.T) {
	text := Build("high", "")

	if !strings.Contains(text, "Default permissiveness for this session: high.") {
		t.Fatalf("expected prompt to include permissiveness, got:\n%s", text)
	}
}

func TestBuildIncludesWelcomeExamples(t *testing.T) {
	text := Build("medium", "")

	if !strings.Contains(text, `You can say things like "play something hazy and nocturnal"`) {
		t.Fatalf("expected prompt welcome examples, got:\n%s", text)
	}
}
