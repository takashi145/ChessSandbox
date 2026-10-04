package ui

import (
	"strings"
	"testing"

	"github.com/takashi145/chess-sandbox/internal/chess"
)

const pawn = "♙"

func snapshotAfter(t *testing.T, moves ...string) chess.Snapshot {
	t.Helper()
	h, err := chess.NewHistory(chess.StandardFEN)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range moves {
		if !h.TryPlay(m) {
			t.Fatalf("expected %s to be legal", m)
		}
	}
	return h.Snapshot()
}

func TestRenderWhiteAndBlackPiecesUseDifferentColors(t *testing.T) {
	out := renderBoard(snapshotAfter(t), false)

	if !strings.Contains(out, "\x1b[37m"+pawn+reset) {
		t.Error("white pawn color missing")
	}
	if !strings.Contains(out, "\x1b[38;5;208m"+pawn+reset) {
		t.Error("black pawn color missing")
	}
}

func TestRenderDoesNotUseTheEmojiPawn(t *testing.T) {
	out := renderBoard(snapshotAfter(t), false)

	for _, bad := range []string{"♟", "︎", "️"} {
		if strings.Contains(out, bad) {
			t.Errorf("output contains %q", bad)
		}
	}
}

func TestRenderCapturingPieceKeepsItsOwnColorOnHighlight(t *testing.T) {
	out := renderBoard(snapshotAfter(t, "e4", "d5", "exd5"), false)

	if !strings.Contains(out, "\x1b[37;48;5;59m"+pawn+reset) {
		t.Error("white pawn on highlight missing")
	}
	if strings.Contains(out, "\x1b[38;5;208;48;5;59m"+pawn+reset) {
		t.Error("captured pawn color leaked onto the highlight")
	}
}

func TestRenderHighlightsFromAndToSquares(t *testing.T) {
	out := renderBoard(snapshotAfter(t, "e4"), false)

	if got := strings.Count(out, "48;5;59"); got != 2 {
		t.Errorf("highlighted squares = %d, want 2", got)
	}
}

func TestRenderAtStartHasNoHighlight(t *testing.T) {
	if strings.Contains(renderBoard(snapshotAfter(t), false), "48;5;59") {
		t.Error("start position should not highlight anything")
	}
}

func TestRenderFlippedReversesFilesAndRanks(t *testing.T) {
	s := snapshotAfter(t)

	normal := renderBoard(s, false)
	flipped := renderBoard(s, true)

	if !strings.HasPrefix(normal, "    a b c d e f g h") {
		t.Error("normal file labels wrong")
	}
	if !strings.HasPrefix(flipped, "    h g f e d c b a") {
		t.Error("flipped file labels wrong")
	}
	if strings.Index(normal, "8 │") > strings.Index(normal, "1 │") {
		t.Error("normal board should list rank 8 first")
	}
	if strings.Index(flipped, "1 │") > strings.Index(flipped, "8 │") {
		t.Error("flipped board should list rank 1 first")
	}
}
