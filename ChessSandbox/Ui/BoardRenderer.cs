using System.Text;
using ChessSandbox.Chess;

namespace ChessSandbox.Ui;

public static class BoardRenderer
{
    // Both sides use filled glyphs; the side is distinguished by color only.
    private const string WhiteColor = "white";
    private const string BlackColor = "orange1";
    private const string HighlightBackground = "grey37";

    private static readonly Dictionary<PieceKind, char> Glyphs = new()
    {
        [PieceKind.King] = '♚',
        [PieceKind.Queen] = '♛',
        [PieceKind.Rook] = '♜',
        [PieceKind.Bishop] = '♝',
        [PieceKind.Knight] = '♞',
        [PieceKind.Pawn] = '♟',
    };

    // Returns Spectre.Console markup for the 8x8 board only (caller adds header/footer).
    public static string Render(BoardSnapshot snapshot, bool flipped)
    {
        var files = FileLabels(flipped);
        var builder = new StringBuilder();

        builder.Append("    ").Append(files).Append('\n');
        builder.Append("  ┌─────────────────┐\n");

        foreach (var rank in Ranks(flipped))
        {
            builder.Append(rank + 1).Append(" │ ");

            foreach (var file in Files(flipped))
            {
                var piece = snapshot.Board[file, rank];
                var highlighted = snapshot.FromSquare == (file, rank) || snapshot.ToSquare == (file, rank);

                builder.Append(Square(piece, highlighted));
                builder.Append(' ');
            }

            builder.Append("│\n");
        }

        builder.Append("  └─────────────────┘\n");
        builder.Append("    ").Append(files);

        return builder.ToString();
    }

    private static string Square(BoardPiece? piece, bool highlighted)
    {
        if (piece is null)
        {
            return highlighted ? $"[on {HighlightBackground}] [/]" : " ";
        }

        var color = piece.Value.Color == Side.White ? WhiteColor : BlackColor;
        var style = highlighted ? $"{color} on {HighlightBackground}" : color;
        return $"[{style}]{Glyphs[piece.Value.Kind]}[/]";
    }

    private static IEnumerable<int> Ranks(bool flipped) =>
        flipped ? Enumerable.Range(0, 8) : Enumerable.Range(0, 8).Reverse();

    private static IEnumerable<int> Files(bool flipped) =>
        flipped ? Enumerable.Range(0, 8).Reverse() : Enumerable.Range(0, 8);

    private static string FileLabels(bool flipped) =>
        string.Join(' ', Files(flipped).Select(f => (char)('a' + f)));
}
