using ChessSandbox.Chess;

namespace ChessSandbox.Tests.Chess;

public class MoveHistoryTests
{
    private static MoveHistory Play(params string[] moves)
    {
        var history = new MoveHistory();
        foreach (var move in moves)
        {
            Assert.True(history.TryPlay(move), $"expected {move} to be legal");
        }

        return history;
    }

    [Fact]
    public void NewHistory_StartsAtStartPosition()
    {
        var history = new MoveHistory();

        Assert.Equal(0, history.Position);
        Assert.Empty(history.Moves);
        Assert.False(history.CanGoBack);
        Assert.False(history.CanGoForward);
        Assert.Equal(MoveHistory.StandardFen, history.CurrentFen);
    }

    [Fact]
    public void TryPlay_AppendsMoveAndAdvances()
    {
        var history = Play("e4", "e5");

        Assert.Equal(["e4", "e5"], history.Moves);
        Assert.Equal(2, history.Position);
    }

    [Fact]
    public void BackAndForward_MoveAlongTheLine()
    {
        var history = Play("e4", "e5", "Nf3");

        history.Back();
        Assert.Equal(2, history.Position);

        history.Back();
        history.Forward();
        Assert.Equal(2, history.Position);
        Assert.Equal(3, history.Moves.Count);
    }

    [Fact]
    public void Back_AtStart_StaysAtStart()
    {
        var history = Play("e4");

        history.Back();
        history.Back();
        history.Back();

        Assert.Equal(0, history.Position);
        Assert.False(history.CanGoBack);
        Assert.Equal(MoveHistory.StandardFen, history.CurrentFen);
    }

    [Fact]
    public void Forward_AtEnd_StaysAtEnd()
    {
        var history = Play("e4", "e5");
        var fen = history.CurrentFen;

        history.Forward();
        history.Forward();

        Assert.Equal(2, history.Position);
        Assert.False(history.CanGoForward);
        Assert.Equal(fen, history.CurrentFen);
    }

    [Fact]
    public void GoTo_ClampsToValidRange()
    {
        var history = Play("e4", "e5");

        history.GoTo(-5);
        Assert.Equal(0, history.Position);

        history.GoTo(99);
        Assert.Equal(2, history.Position);
    }

    [Fact]
    public void PlayingAfterGoingBack_DiscardsMovesAhead()
    {
        var history = Play("e4", "e5", "Nf3", "Nc6");

        history.GoTo(2);
        Assert.True(history.CanGoForward);
        Assert.True(history.TryPlay("d4"));

        Assert.Equal(["e4", "e5", "d4"], history.Moves);
        Assert.Equal(3, history.Position);
        Assert.False(history.CanGoForward);
    }

    [Fact]
    public void CanGoForward_IsFalseAtEndOfLine()
    {
        Assert.False(Play("e4", "e5").CanGoForward);
    }

    [Fact]
    public void IllegalMove_IsRejectedAndNothingChanges()
    {
        var history = Play("e4", "e5");
        var fen = history.CurrentFen;

        Assert.False(history.TryPlay("Ke5"));
        Assert.False(history.TryPlay("nonsense"));
        Assert.False(history.TryPlay(""));

        Assert.Equal(["e4", "e5"], history.Moves);
        Assert.Equal(2, history.Position);
        Assert.Equal(fen, history.CurrentFen);
    }

    [Fact]
    public void IllegalMove_AfterGoingBack_DoesNotDiscardMovesAhead()
    {
        var history = Play("e4", "e5", "Nf3");
        history.GoTo(1);

        Assert.False(history.TryPlay("Ke5"));

        Assert.Equal(3, history.Moves.Count);
        Assert.Equal(1, history.Position);
    }

    [Fact]
    public void MoveThatLeavesKingInCheck_IsRejected()
    {
        // White king e1 is in check from the rook on e8; the pawn move does not resolve it.
        Assert.True(MoveHistory.TryCreate("4r2k/8/8/8/8/8/P7/4K3 w - - 0 1", out var history));

        Assert.False(history!.TryPlay("a3"));
        Assert.True(history.TryPlay("Kd1"));
    }

    [Fact]
    public void CurrentSnapshot_ReflectsLastMove()
    {
        var history = Play("e4", "e5");

        var snapshot = history.CurrentSnapshot();

        Assert.Equal("e5", snapshot.San);
        Assert.Equal(2, snapshot.MoveNumber);
        Assert.Equal(Side.White, snapshot.SideToMove);
    }

    [Fact]
    public void CurrentSnapshot_AtStart_HasNoLastMove()
    {
        var history = Play("e4");
        history.GoTo(0);

        var snapshot = history.CurrentSnapshot();

        Assert.Null(snapshot.San);
        Assert.Null(snapshot.FromSquare);
    }

    [Fact]
    public void Check_IsReportedInSnapshot()
    {
        var history = Play("e4", "f5", "Qh5");

        Assert.True(history.CurrentSnapshot().IsCheck);
    }

    [Fact]
    public void Promotion_CanChooseAnyPiece()
    {
        Assert.True(MoveHistory.TryCreate("7k/P7/8/8/8/8/8/K7 w - - 0 1", out var history));

        Assert.True(history!.TryPlay("a8=N"));

        var piece = history.CurrentSnapshot().Board[0, 7];
        Assert.Equal(new BoardPiece(Side.White, PieceKind.Knight), piece);
    }

    [Fact]
    public void TryCreate_FromFen_StartsAtThatPosition()
    {
        const string fen = "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1";

        Assert.True(MoveHistory.TryCreate(fen, out var history));
        Assert.Equal(fen, history!.StartFen);
        Assert.Equal(Side.Black, history.CurrentSnapshot().SideToMove);
    }

    [Theory]
    [InlineData("")]
    [InlineData("bad")]
    [InlineData("8/8/8 w - - 0 1")]
    public void TryCreate_InvalidFen_ReturnsFalse(string fen)
    {
        Assert.False(MoveHistory.TryCreate(fen, out var history));
        Assert.Null(history);
    }

    [Fact]
    public void TryRestore_RebuildsLineAndPosition()
    {
        Assert.True(MoveHistory.TryRestore(MoveHistory.StandardFen, ["e4", "e5", "Nf3"], 2, out var history));

        Assert.Equal(3, history!.Moves.Count);
        Assert.Equal(2, history.Position);
    }

    [Fact]
    public void TryRestore_IllegalMove_ReturnsFalse()
    {
        Assert.False(MoveHistory.TryRestore(MoveHistory.StandardFen, ["e4", "Ke5"], 2, out var history));
        Assert.Null(history);
    }

    [Theory]
    [InlineData(-1)]
    [InlineData(3)]
    public void TryRestore_PositionOutOfRange_ReturnsFalse(int position)
    {
        Assert.False(MoveHistory.TryRestore(MoveHistory.StandardFen, ["e4", "e5"], position, out _));
    }
}
