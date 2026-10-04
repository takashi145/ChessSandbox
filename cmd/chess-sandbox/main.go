package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/takashi145/chess-sandbox/internal/chess"
	"github.com/takashi145/chess-sandbox/internal/ui"
)

const usage = `Usage: chess-sandbox [--fen "<FEN>" | --pgn <file>]`

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	fen, pgnPath := "", ""

	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--fen" && i+1 < len(args):
			i++
			fen = args[i]
		case args[i] == "--pgn" && i+1 < len(args):
			i++
			pgnPath = args[i]
		default:
			fmt.Fprintln(os.Stderr, usage)
			return 1
		}
	}

	if fen != "" && pgnPath != "" {
		fmt.Fprintln(os.Stderr, "--fen and --pgn cannot be used together.")
		return 1
	}

	store := chess.DefaultSessionStore()
	var model ui.Model
	if pgnPath != "" {
		history, err := loadPGN(pgnPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Cannot load PGN:", err)
			return 1
		}
		model = ui.NewFromHistory(store, history)
	} else {
		var err error
		model, err = ui.New(store, fen)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Invalid FEN.")
			return 1
		}
	}

	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func loadPGN(path string) (*chess.History, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	games, err := chess.ParsePGNGames(file)
	if err != nil {
		return nil, err
	}
	// A PGN is opened to read a game from the beginning, so start at the first position rather than the last move.
	return chess.RestoreHistory(games[0].StartFEN, games[0].Moves, 0)
}
