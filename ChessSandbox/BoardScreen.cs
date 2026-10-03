using System.Text;
using ChessSandbox.Chess;
using ChessSandbox.Ui;
using Spectre.Console;

namespace ChessSandbox;

// Input loop. Typed text is a SAN move; commands start with ':'. Arrow/Home/End keys navigate when the input is empty.
internal sealed class BoardScreen(MoveHistory history, SessionStore store)
{
    private const string CommandList = ":flip :fen :home :end :quit";

    private readonly StringBuilder _input = new();
    private bool _flipped;
    private bool _quit;
    private string? _message;

    public void Run()
    {
        using var _ = Terminal.HideCursor();

        while (!_quit)
        {
            Draw();
            HandleKey(Console.ReadKey(intercept: true));
        }

        Save();
    }

    private void HandleKey(ConsoleKeyInfo key)
    {
        switch (key.Key)
        {
            case ConsoleKey.Escape:
                _quit = true;
                return;
            case ConsoleKey.Enter:
                SubmitInput();
                return;
            case ConsoleKey.Backspace:
                RemoveLastInputChar();
                return;
        }

        if (_input.Length == 0 && TryNavigate(key.Key))
        {
            _message = null;
            Save();
            return;
        }

        if (!char.IsControl(key.KeyChar))
        {
            _input.Append(key.KeyChar);
        }
    }

    private void RemoveLastInputChar()
    {
        if (_input.Length > 0)
        {
            _input.Length--;
        }
    }

    private bool TryNavigate(ConsoleKey key)
    {
        switch (key)
        {
            case ConsoleKey.LeftArrow: history.Back(); return true;
            case ConsoleKey.RightArrow: history.Forward(); return true;
            case ConsoleKey.Home: history.GoTo(0); return true;
            case ConsoleKey.End: history.GoTo(history.Moves.Count); return true;
            default: return false;
        }
    }

    private void SubmitInput()
    {
        var text = _input.ToString().Trim();
        _input.Clear();

        if (text.Length == 0)
        {
            return;
        }

        _message = null;

        if (text.StartsWith(':'))
        {
            RunCommand(text[1..].Trim().ToLowerInvariant());
        }
        else
        {
            PlayMove(text);
        }

        Save();
    }

    private void RunCommand(string command)
    {
        switch (command)
        {
            case "f" or "flip": _flipped = !_flipped; break;
            case "fen": _message = "FEN: " + history.CurrentFen; break;
            case "home": history.GoTo(0); break;
            case "end": history.GoTo(history.Moves.Count); break;
            case "q" or "quit": _quit = true; break;
            default: _message = "Commands: " + CommandList; break;
        }
    }

    private void PlayMove(string san)
    {
        if (history.CanGoForward && !Confirm("This will discard the moves ahead. Continue? (y/N)"))
        {
            _message = "Cancelled";
            return;
        }

        if (!history.TryPlay(san))
        {
            _message = $"Illegal move: {san}";
        }
    }

    private bool Confirm(string question)
    {
        DrawBoard();
        AnsiConsole.MarkupLine($"[red]{Markup.Escape(question)}[/]");
        return Console.ReadKey(intercept: true).Key == ConsoleKey.Y;
    }

    private void Save()
    {
        if (!store.Save(history))
        {
            _message = "Could not save the session";
        }
    }

    private void Draw()
    {
        DrawBoard();
        AnsiConsole.Markup("> " + Markup.Escape(_input.ToString()) + "_");
    }

    private void DrawBoard()
    {
        var snapshot = history.CurrentSnapshot();

        AnsiConsole.Clear();
        AnsiConsole.Markup(BoardRenderer.Render(snapshot, _flipped));
        AnsiConsole.WriteLine();
        AnsiConsole.WriteLine();
        AnsiConsole.MarkupLine(Markup.Escape(DescribeLastMove(snapshot)));
        AnsiConsole.MarkupLine($"[grey]{Markup.Escape(Hint())}[/]");

        if (_message is not null)
        {
            AnsiConsole.MarkupLine($"[yellow]{Markup.Escape(_message)}[/]");
        }
    }

    private string Hint() =>
        $"{history.Position}/{history.Moves.Count}  ←→ move  Home/End  {CommandList}  Esc save and quit";

    private static string DescribeLastMove(BoardSnapshot snapshot)
    {
        if (snapshot.San is null)
        {
            return "(start position)";
        }

        var number = (snapshot.MoveNumber + 1) / 2;
        var dots = snapshot.MoveNumber % 2 == 1 ? "." : "...";
        var suffix = snapshot.IsCheckmate ? " checkmate" : snapshot.IsCheck ? " check" : "";
        return $"{number}{dots} {snapshot.San}{suffix}";
    }
}
