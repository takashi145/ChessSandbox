package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/takashi145/chess-sandbox/internal/chess"
	"github.com/takashi145/chess-sandbox/internal/ui"
)

const usage = `Usage: chess-sandbox [--fen "<FEN>" | --pgn <file>] [--version]`

// Set at build time by goreleaser.
var version = "dev"

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	fen, pgnPath := "", ""

	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--version":
			fmt.Println(version)
			return 0
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
		games, err := loadPGN(pgnPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Cannot load PGN:", err)
			return 1
		}
		model, err = ui.NewFromGames(store, games)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Cannot load PGN:", err)
			return 1
		}
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

func loadPGN(path string) ([]chess.PGNGame, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return chess.ParsePGNGames(file)
}
