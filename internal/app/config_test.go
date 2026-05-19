package app

import "testing"

func TestParseArgsDefaults(t *testing.T) {
	cfg, err := ParseArgs([]string{})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}

	if cfg.Agent != "" {
		t.Fatalf("expected empty agent, got %q", cfg.Agent)
	}

	if cfg.Permissiveness != "medium" {
		t.Fatalf("expected default permissiveness medium, got %q", cfg.Permissiveness)
	}
}

func TestParseArgsExplicitAgentAndPermissiveness(t *testing.T) {
	cfg, err := ParseArgs([]string{"--codex", "--permissive", "high"})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}

	if cfg.Agent != "codex" {
		t.Fatalf("expected codex agent, got %q", cfg.Agent)
	}

	if cfg.Permissiveness != "high" {
		t.Fatalf("expected permissiveness high, got %q", cfg.Permissiveness)
	}
}

func TestParseArgsRejectsConflictingAgents(t *testing.T) {
	_, err := ParseArgs([]string{"--claude", "--codex"})
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
}

func TestParseArgsRejectsInvalidPermissiveness(t *testing.T) {
	_, err := ParseArgs([]string{"-p", "wild"})
	if err == nil {
		t.Fatal("expected invalid permissiveness error, got nil")
	}
}

func TestParseArgsVerboseFlag(t *testing.T) {
	cfg, err := ParseArgs([]string{"-V"})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}
	if !cfg.Verbose {
		t.Fatal("expected verbose to be true")
	}
}
