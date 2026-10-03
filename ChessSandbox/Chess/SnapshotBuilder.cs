using Chess;

namespace ChessSandbox.Chess;

public static class SnapshotBuilder
{
    // Snapshot of `board` as it stands, after its last executed move (or the start position if there is none).
    public static BoardSnapshot Capture(ChessBoard board)
    {
        var moves = board.ExecutedMoves;

        if (moves.Count == 0)
        {
            return new BoardSnapshot
            {
                Board = CaptureBoard(board),
                MoveNumber = 0,
                SideToMove = ToSide(board.Turn),
            };
        }

        var move = moves[^1];

        return new BoardSnapshot
        {
            Board = CaptureBoard(board),
            MoveNumber = moves.Count,
            SideToMove = ToSide(board.Turn),
            San = move.San,
            FromSquare = (move.OriginalPosition.X, move.OriginalPosition.Y),
            ToSquare = (move.NewPosition.X, move.NewPosition.Y),
            IsCheck = move.IsCheck,
            IsCheckmate = move.IsMate,
        };
    }

    private static BoardPiece?[,] CaptureBoard(ChessBoard board)
    {
        var grid = new BoardPiece?[8, 8];

        for (var file = 0; file < 8; file++)
        {
            for (var rank = 0; rank < 8; rank++)
            {
                var piece = board[file, rank];
                grid[file, rank] = piece is null ? null : ToBoardPiece(piece);
            }
        }

        return grid;
    }

    private static BoardPiece ToBoardPiece(Piece piece) => new(ToSide(piece.Color), ToKind(piece.Type));

    private static Side ToSide(PieceColor color) => color == PieceColor.White ? Side.White : Side.Black;

    private static PieceKind ToKind(PieceType type)
    {
        if (type == PieceType.Pawn) return PieceKind.Pawn;
        if (type == PieceType.Knight) return PieceKind.Knight;
        if (type == PieceType.Bishop) return PieceKind.Bishop;
        if (type == PieceType.Rook) return PieceKind.Rook;
        if (type == PieceType.Queen) return PieceKind.Queen;
        if (type == PieceType.King) return PieceKind.King;
        throw new ArgumentOutOfRangeException(nameof(type), type, "Unknown piece type.");
    }
}
