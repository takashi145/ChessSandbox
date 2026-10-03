using ChessSandbox;
using ChessSandbox.Chess;
using Spectre.Console;

Console.OutputEncoding = System.Text.Encoding.UTF8;

var store = SessionStore.CreateDefault();
string? fen = null;

for (var i = 0; i < args.Length; i++)
{
    if (args[i] == "--fen" && i + 1 < args.Length)
    {
        fen = args[++i];
    }
    else
    {
        Console.Error.WriteLine("Usage: chess-sandbox [--fen \"<FEN>\"]");
        return 1;
    }
}

MoveHistory? history;

if (fen is not null)
{
    if (!MoveHistory.TryCreate(fen, out history))
    {
        Console.Error.WriteLine("Invalid FEN.");
        return 1;
    }
}
else
{
    history = ChooseStart(store);
}

new BoardScreen(history, store).Run();
return 0;

static MoveHistory ChooseStart(SessionStore store)
{
    var saved = store.Load();
    if (saved is null)
    {
        return new MoveHistory();
    }

    const string Continue = "Continue";
    const string Standard = "New game (standard position)";
    const string Custom = "New game (from FEN)";

    var choice = AnsiConsole.Prompt(
        new SelectionPrompt<string>().Title("A saved session was found").AddChoices(Continue, Standard, Custom));

    if (choice == Continue)
    {
        return saved;
    }

    if (choice == Standard)
    {
        return new MoveHistory();
    }

    while (true)
    {
        var text = AnsiConsole.Ask<string>("FEN:");
        if (MoveHistory.TryCreate(text.Trim(), out var created))
        {
            return created;
        }

        AnsiConsole.MarkupLine("[red]Invalid FEN.[/]");
    }
}
