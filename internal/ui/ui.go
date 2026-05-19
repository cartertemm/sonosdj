package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/cartertemm/sonosdj/internal/agents"
	"golang.org/x/term"
)

type KeyReader interface {
	ReadKey() (byte, error)
}

type TerminalKeyReader struct {
	reader *bufio.Reader
	file   *os.File
}

func NewTerminalKeyReader(file *os.File) *TerminalKeyReader {
	return &TerminalKeyReader{reader: bufio.NewReader(file), file: file}
}

func (r *TerminalKeyReader) ReadKey() (byte, error) {
	if term.IsTerminal(int(r.file.Fd())) {
		state, err := term.MakeRaw(int(r.file.Fd()))
		if err != nil {
			return 0, err
		}
		defer term.Restore(int(r.file.Fd()), state)
	}
	return r.reader.ReadByte()
}

func PromptYesNo(reader KeyReader, output io.Writer, prompt string) (bool, error) {
	fmt.Fprintf(output, "%s (y/n) ", prompt)
	for {
		key, err := reader.ReadKey()
		if err != nil {
			return false, err
		}
		switch strings.ToLower(string(key)) {
		case "y":
			fmt.Fprintln(output, "y")
			return true, nil
		case "n":
			fmt.Fprintln(output, "n")
			return false, nil
		}
	}
}

func PromptMenu(reader KeyReader, output io.Writer, options []agents.Agent) (agents.Agent, error) {
	fmt.Fprintln(output, "Which agent do you want to use?")
	for i, agent := range options {
		fmt.Fprintf(output, "%d. %s\n", i+1, agent.Name)
	}
	fmt.Fprint(output, "> ")
	for {
		key, err := reader.ReadKey()
		if err != nil {
			return agents.Agent{}, err
		}
		fmt.Fprintln(output, string(key))
		index, err := strconv.Atoi(string(key))
		if err != nil || index < 1 || index > len(options) {
			continue
		}
		return options[index-1], nil
	}
}
