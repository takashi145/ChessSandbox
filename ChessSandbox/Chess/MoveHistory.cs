using System.Diagnostics.CodeAnalysis;
using Chess;

namespace ChessSandbox.Chess;

// One line of play: start position + SAN moves + current position.
public sealed class MoveHistory
{
    public const string StandardFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1";

    private readonly List<string> _moves = [];
    private ChessBoard _board;

    public MoveHistory(string startFen = StandardFen)
    {
        StartFen = startFen;
        _board = ChessBoard.LoadFromFen(startFen);
    }

    public string StartFen { get; }
    public IReadOnlyList<string> Moves => _moves;
    public int Position { get; private set; }
    public bool CanGoBack => Position > 0;

    // Playing a move while this is true discards the moves ahead of the current position.
    public bool CanGoForward => Position < _moves.Count;

    public string CurrentFen => _board.ToFen();

    public static bool TryCreate(string fen, [NotNullWhen(true)] out MoveHistory? history)
    {
        try
        {
            history = new MoveHistory(fen);
            return true;
        }
        catch (Exception)
        {
            history = null;
            return false;
        }
    }

    public static bool TryRestore(string startFen, IEnumerable<string> moves, int position, [NotNullWhen(true)] out MoveHistory? history)
    {
        history = null;
        if (!TryCreate(startFen, out var created))
        {
            return false;
        }

        foreach (var san in moves)
        {
            if (!created.TryPlayOn(created._board, san))
            {
                return false;
            }
        }

        if (position < 0 || position > created._moves.Count)
        {
            return false;
        }

        created.GoTo(position);
        history = created;
        return true;
    }

    // Plays a move at the current position, dropping any moves ahead. State is unchanged on failure.
    public bool TryPlay(string san)
    {
        var probe = Rebuild(Position);
        if (!TryPlayOn(probe, san))
        {
            return false;
        }

        _board = probe;
        return true;
    }

    public void Back() => GoTo(Position - 1);

    public void Forward() => GoTo(Position + 1);

    public void GoTo(int position)
    {
        position = Math.Clamp(position, 0, _moves.Count);
        if (position == Position)
        {
            return;
        }

        _board = Rebuild(position);
        Position = position;
    }

    public BoardSnapshot CurrentSnapshot() => SnapshotBuilder.Capture(_board);

    // Plays `san` on `board`, which must be at the current position, and records it.
    private bool TryPlayOn(ChessBoard board, string san)
    {
        var move = san.Trim();
        if (!TryMove(board, move))
        {
            return false;
        }

        _moves.RemoveRange(Position, _moves.Count - Position);
        _moves.Add(board.ExecutedMoves[^1].San ?? move);
        Position++;
        return true;
    }

    private ChessBoard Rebuild(int position)
    {
        var board = ChessBoard.LoadFromFen(StartFen);
        for (var i = 0; i < position; i++)
        {
            if (!board.Move(_moves[i]))
            {
                throw new InvalidOperationException($"Recorded move '{_moves[i]}' is no longer legal.");
            }
        }

        return board;
    }

    private static bool TryMove(ChessBoard board, string san)
    {
        try
        {
            return board.IsValidMove(san) && board.Move(san);
        }
        catch (Exception)
        {
            return false;
        }
    }
}
