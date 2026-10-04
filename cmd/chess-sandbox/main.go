package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/takashi145/chess-sandbox/internal/chess"
	"github.com/takashi145/chess-sandbox/internal/ui"
)

const usage = `Usage: chess-sandbox [--fen "<FEN>"]`

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	fen := ""

	for i := 0; i < len(args); i++ {
		if args[i] == "--fen" && i+1 < len(args) {
			i++
			fen = args[i]
			continue
		}
		fmt.Fprintln(os.Stderr, usage)
		return 1
	}

	model, err := ui.New(chess.DefaultSessionStore(), fen)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid FEN.")
		return 1
	}

	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
