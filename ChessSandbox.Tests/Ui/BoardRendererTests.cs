using ChessSandbox.Chess;
using ChessSandbox.Ui;

namespace ChessSandbox.Tests.Ui;

public class BoardRendererTests
{
    private const string Pawn = "♙";

    [Fact]
    public void Render_WhiteAndBlackPiecesUseDifferentColors()
    {
        var markup = BoardRenderer.Render(new MoveHistory().CurrentSnapshot(), flipped: false);

        Assert.Contains($"[white]{Pawn}[/]", markup);
        Assert.Contains($"[orange1]{Pawn}[/]", markup);
    }

    [Fact]
    public void Render_DoesNotUseTheEmojiPawn()
    {
        var markup = BoardRenderer.Render(new MoveHistory().CurrentSnapshot(), flipped: false);

        Assert.DoesNotContain("♟", markup, StringComparison.Ordinal);
        Assert.DoesNotContain("︎", markup, StringComparison.Ordinal);
        Assert.DoesNotContain("️", markup, StringComparison.Ordinal);
    }

    [Fact]
    public void Render_CapturingPiece_KeepsItsOwnColorOnHighlight()
    {
        var history = new MoveHistory();
        foreach (var move in new[] { "e4", "d5", "exd5" })
        {
            Assert.True(history.TryPlay(move));
        }

        var markup = BoardRenderer.Render(history.CurrentSnapshot(), flipped: false);

        Assert.Contains($"[white on grey37]{Pawn}[/]", markup);
        Assert.DoesNotContain($"[orange1 on grey37]{Pawn}[/]", markup);
    }

    [Fact]
    public void Render_HighlightsFromAndToSquares()
    {
        var history = new MoveHistory();
        Assert.True(history.TryPlay("e4"));

        var markup = BoardRenderer.Render(history.CurrentSnapshot(), flipped: false);

        Assert.Equal(2, CountOf(markup, "on grey37"));
    }

    [Fact]
    public void Render_AtStart_HasNoHighlight()
    {
        var markup = BoardRenderer.Render(new MoveHistory().CurrentSnapshot(), flipped: false);

        Assert.DoesNotContain("on grey37", markup);
    }

    [Fact]
    public void Render_Flipped_ReversesFilesAndRanks()
    {
        var snapshot = new MoveHistory().CurrentSnapshot();

        var normal = BoardRenderer.Render(snapshot, flipped: false);
        var flipped = BoardRenderer.Render(snapshot, flipped: true);

        Assert.StartsWith("    a b c d e f g h", normal);
        Assert.StartsWith("    h g f e d c b a", flipped);
        Assert.True(normal.IndexOf("8 │", StringComparison.Ordinal) < normal.IndexOf("1 │", StringComparison.Ordinal));
        Assert.True(flipped.IndexOf("1 │", StringComparison.Ordinal) < flipped.IndexOf("8 │", StringComparison.Ordinal));
    }

    private static int CountOf(string text, string value)
    {
        var count = 0;
        for (var i = text.IndexOf(value, StringComparison.Ordinal); i >= 0; i = text.IndexOf(value, i + value.Length, StringComparison.Ordinal))
        {
            count++;
        }

        return count;
    }
}
