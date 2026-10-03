namespace ChessSandbox.Ui;

internal static class Terminal
{
    public static IDisposable HideCursor()
    {
        SetCursorVisible(false);
        return new CursorRestorer();
    }

    private static void SetCursorVisible(bool visible)
    {
        try
        {
            Console.CursorVisible = visible;
        }
        catch (IOException)
        {
        }
    }

    private sealed class CursorRestorer : IDisposable
    {
        public void Dispose() => SetCursorVisible(true);
    }
}
