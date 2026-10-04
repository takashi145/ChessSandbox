# ChessSandbox

English | [日本語](README.ja.md)

A terminal chess board for trying out positions. Play both sides, step back and forward through your moves, and pick up where you left off.

<img width="733" height="322" alt="ChessSandbox board" src="https://github.com/user-attachments/assets/71e49067-fc0c-4f3d-9991-ffea8b15b2ec" />

## Installation

**Windows** (PowerShell)

```powershell
irm https://raw.githubusercontent.com/takashi145/chess-sandbox/main/install.ps1 | iex
```

**macOS / Linux**

```sh
curl -fsSL https://raw.githubusercontent.com/takashi145/chess-sandbox/main/install.sh | sh
```

Run the same command again to update. You can also download the archive for your platform from the
[Releases](https://github.com/takashi145/chess-sandbox/releases/latest) page, extract it, and put `chess-sandbox` on your `PATH` yourself.

### What the installer does

- **Windows**: puts `chess-sandbox.exe` in `%LOCALAPPDATA%\Programs\chess-sandbox` and adds that folder to your user `PATH`.
- **macOS / Linux**: puts `chess-sandbox` in `~/.local/bin` (set `CHESS_SANDBOX_INSTALL_DIR` to change it).

### Uninstall

**Windows**

1. Delete the install folder:

   ```powershell
   Remove-Item -Recurse "$env:LOCALAPPDATA\Programs\chess-sandbox"
   ```

2. Open **Edit environment variables for your account** from the Start menu, select `Path`, and remove the
   `...\Programs\chess-sandbox` entry.

**macOS / Linux**

```sh
rm ~/.local/bin/chess-sandbox
```

To also remove your saved session, delete the `ChessSandbox` folder described under [Saved session](#saved-session).

## Usage

```
chess-sandbox [--fen "<FEN>"]
```

| Option | Description |
|---|---|
| `--fen "<FEN>"` | Start a new game from this position. Takes priority over a saved session |

With no options, it starts from the standard position. If a saved session exists, you can choose to continue it,
start a new game from the standard position, or start a new game from a FEN you type in.

### Examples

```
# Start from the standard position (or continue a saved session)
chess-sandbox

# Start from a specific position
chess-sandbox --fen "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1"
```

## Controls

Type a move in SAN (`e4`, `Nf3`, `O-O`, `exd5`, `e8=Q`) and press `Enter`. You play both sides. Illegal moves are
rejected and the board does not change.

| Key | Action |
|---|---|
| `←` / `→` | Previous / next move (when the input is empty) |
| `Home` / `End` | Jump to start / end of the line (when the input is empty) |
| `Enter` | Play the typed move or run the typed command |
| `Esc` | Save and quit |

Commands start with `:`.

| Command | Action |
|---|---|
| `:flip` (`:f`) | Flip the board |
| `:fen` | Show the FEN of the current position |
| `:home` / `:end` | Jump to start / end of the line |
| `:quit` (`:q`) | Save and quit |

If you step back and play a different move, the moves ahead are discarded. You are asked to confirm first
(`y` to continue, any other key to cancel).

## Saved session

The session (start position, moves and current position) is saved after every move.
Closing the terminal keeps the latest state. It is saved to the following JSON file:

- **Windows**: `%LOCALAPPDATA%\ChessSandbox\session.json`
- **macOS / Linux**: `~/.local/share/ChessSandbox/session.json`

## License

MIT
