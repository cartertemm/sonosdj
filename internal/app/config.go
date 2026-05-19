package app

import (
	"errors"
	"flag"
	"fmt"
)

type Config struct {
	Agent          string
	Permissiveness string
	Room           string
	Verbose        bool
}

func ParseArgs(args []string) (Config, error) {
	fs := flag.NewFlagSet("sonosdj", flag.ContinueOnError)
	var claude, codex, verbose bool
	var permissive string
	fs.BoolVar(&claude, "claude", false, "launch Claude")
	fs.BoolVar(&codex, "codex", false, "launch Codex")
	fs.StringVar(&permissive, "p", "medium", "default permissiveness")
	fs.StringVar(&permissive, "permissive", "medium", "default permissiveness")
	var room string
	fs.StringVar(&room, "r", "", "default Sonos room")
	fs.StringVar(&room, "room", "", "default Sonos room")
	fs.BoolVar(&verbose, "V", false, "increase verbosity")
	fs.BoolVar(&verbose, "verbose", false, "increase verbosity")
	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}
	if claude && codex {
		return Config{}, errors.New("choose either --claude or --codex, not both")
	}
	switch permissive {
	case "low", "medium", "high":
	default:
		return Config{}, fmt.Errorf("invalid permissiveness %q (expected low, medium, or high)", permissive)
	}
	cfg := Config{Permissiveness: permissive, Room: room, Verbose: verbose}
	if claude {
		cfg.Agent = "claude"
	}
	if codex {
		cfg.Agent = "codex"
	}
	return cfg, nil
}
