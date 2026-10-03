namespace ChessSandbox.Chess;

public enum Side
{
    White,
    Black,
}

public enum PieceKind
{
    Pawn,
    Knight,
    Bishop,
    Rook,
    Queen,
    King,
}

public readonly record struct BoardPiece(Side Color, PieceKind Kind);

public sealed class BoardSnapshot
{
    public required BoardPiece?[,] Board { get; init; }
    public required int MoveNumber { get; init; }
    public required Side SideToMove { get; init; }
    public string? San { get; init; }
    public (int File, int Rank)? FromSquare { get; init; }
    public (int File, int Rank)? ToSquare { get; init; }
    public bool IsCheck { get; init; }
    public bool IsCheckmate { get; init; }
}
