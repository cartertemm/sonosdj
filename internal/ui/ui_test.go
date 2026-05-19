package ui

import (
	"bytes"
	"errors"
	"testing"

	"github.com/cartertemm/sonosdj/internal/agents"
)

type fakeKeyReader struct {
	keys []byte
	err  error
}

func (f *fakeKeyReader) ReadKey() (byte, error) {
	if f.err != nil {
		return 0, f.err
	}
	if len(f.keys) == 0 {
		return 0, errors.New("no keys left")
	}
	key := f.keys[0]
	f.keys = f.keys[1:]
	return key, nil
}

func TestPromptYesNoAcceptsSingleKey(t *testing.T) {
	reader := &fakeKeyReader{keys: []byte{'y'}}
	output := new(bytes.Buffer)

	accepted, err := PromptYesNo(reader, output, "Install?")
	if err != nil {
		t.Fatalf("PromptYesNo returned error: %v", err)
	}
	if !accepted {
		t.Fatal("expected accepted to be true")
	}
}

func TestPromptMenuSelectsSingleDigit(t *testing.T) {
	reader := &fakeKeyReader{keys: []byte{'2'}}
	output := new(bytes.Buffer)
	options := []agents.Agent{
		{Name: "claude"},
		{Name: "codex"},
	}

	selected, err := PromptMenu(reader, output, options)
	if err != nil {
		t.Fatalf("PromptMenu returned error: %v", err)
	}
	if selected.Name != "codex" {
		t.Fatalf("expected codex, got %q", selected.Name)
	}
}
