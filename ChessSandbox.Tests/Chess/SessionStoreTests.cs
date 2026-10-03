using ChessSandbox.Chess;

namespace ChessSandbox.Tests.Chess;

public sealed class SessionStoreTests : IDisposable
{
    private const string StandardFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1";

    private readonly string _directory = Path.Combine(Path.GetTempPath(), "ChessSandboxTests-" + Guid.NewGuid().ToString("N"));

    private string SessionPath => Path.Combine(_directory, "nested", "session.json");

    public void Dispose()
    {
        if (Directory.Exists(_directory))
        {
            Directory.Delete(_directory, recursive: true);
        }
    }

    [Fact]
    public void SaveThenLoad_RestoresSamePosition()
    {
        var history = new MoveHistory();
        foreach (var move in new[] { "e4", "e5", "Nf3", "Nc6" })
        {
            Assert.True(history.TryPlay(move));
        }

        history.GoTo(2);

        var store = new SessionStore(SessionPath);
        store.Save(history);
        var loaded = store.Load();

        Assert.NotNull(loaded);
        Assert.Equal(history.StartFen, loaded.StartFen);
        Assert.Equal(history.Moves, loaded.Moves);
        Assert.Equal(history.Position, loaded.Position);
        Assert.Equal(history.CurrentFen, loaded.CurrentFen);
    }

    [Fact]
    public void SaveThenLoad_KeepsCustomStartFen()
    {
        const string fen = "7k/P7/8/8/8/8/8/K7 w - - 0 1";
        Assert.True(MoveHistory.TryCreate(fen, out var history));
        Assert.True(history!.TryPlay("a8=Q"));

        var store = new SessionStore(SessionPath);
        store.Save(history);
        var loaded = store.Load();

        Assert.NotNull(loaded);
        Assert.Equal(fen, loaded.StartFen);
        Assert.Equal(history.Moves, loaded.Moves);
    }

    [Fact]
    public void Save_CreatesMissingDirectory()
    {
        var store = new SessionStore(SessionPath);

        store.Save(new MoveHistory());

        Assert.True(File.Exists(SessionPath));
    }

    [Fact]
    public void Save_OverwritesPreviousSession()
    {
        var store = new SessionStore(SessionPath);
        var first = new MoveHistory();
        first.TryPlay("e4");
        store.Save(first);

        store.Save(new MoveHistory());

        Assert.Empty(store.Load()!.Moves);
    }

    [Fact]
    public void Load_MissingFile_ReturnsNull()
    {
        Assert.Null(new SessionStore(SessionPath).Load());
    }

    [Theory]
    [InlineData("")]
    [InlineData("not json")]
    [InlineData("{")]
    [InlineData("null")]
    [InlineData("{}")]
    [InlineData("{\"StartFen\":\"bad\",\"Moves\":[],\"Position\":0}")]
    [InlineData("{\"StartFen\":\"" + StandardFen + "\",\"Moves\":[\"e4\",\"Ke5\"],\"Position\":2}")]
    [InlineData("{\"StartFen\":\"" + StandardFen + "\",\"Moves\":[\"e4\"],\"Position\":5}")]
    public void Load_CorruptFile_ReturnsNull(string content)
    {
        Directory.CreateDirectory(Path.GetDirectoryName(SessionPath)!);
        File.WriteAllText(SessionPath, content);

        Assert.Null(new SessionStore(SessionPath).Load());
    }
}
