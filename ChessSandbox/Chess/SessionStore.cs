using System.Text.Json;

namespace ChessSandbox.Chess;

internal sealed record SessionData(string StartFen, List<string> Moves, int Position);

public sealed class SessionStore(string path)
{
    public static string DefaultPath { get; } = Path.Combine(
        Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData),
        "ChessSandbox",
        "session.json");

    public static SessionStore CreateDefault() => new(DefaultPath);

    // Returns false if the session could not be written. The previous file is left intact in that case.
    public bool Save(MoveHistory history)
    {
        var data = new SessionData(history.StartFen, [.. history.Moves], history.Position);
        var temp = path + ".tmp";

        try
        {
            Directory.CreateDirectory(Path.GetDirectoryName(path)!);
            File.WriteAllText(temp, JsonSerializer.Serialize(data));
            File.Move(temp, path, overwrite: true);
            return true;
        }
        catch (Exception ex) when (ex is IOException or UnauthorizedAccessException)
        {
            return false;
        }
    }

    // Returns null if the file is missing, unreadable or does not describe a legal game.
    public MoveHistory? Load()
    {
        try
        {
            if (!File.Exists(path))
            {
                return null;
            }

            var data = JsonSerializer.Deserialize<SessionData>(File.ReadAllText(path));
            if (data?.StartFen is null || data.Moves is null)
            {
                return null;
            }

            return MoveHistory.TryRestore(data.StartFen, data.Moves, data.Position, out var history) ? history : null;
        }
        catch (Exception ex) when (ex is IOException or JsonException or UnauthorizedAccessException)
        {
            return null;
        }
    }
}
