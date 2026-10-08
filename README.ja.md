# ChessSandbox

局面を試すための、ターミナル用のチェス盤です。両方の側を自分で指し、手を戻したり進めたりでき、前回の続きから再開できます。

<img width="733" height="322" alt="ChessSandbox の盤面" src="https://github.com/user-attachments/assets/71e49067-fc0c-4f3d-9991-ffea8b15b2ec" />

## インストール

**Windows**（PowerShell）

```powershell
irm https://raw.githubusercontent.com/takashi145/chess-sandbox/main/install.ps1 | iex
```

**macOS / Linux**

```sh
curl -fsSL https://raw.githubusercontent.com/takashi145/chess-sandbox/main/install.sh | sh
```

アップデートするときは、同じコマンドをもう一度実行してください。
[Releases](https://github.com/takashi145/chess-sandbox/releases/latest) ページからお使いの環境用のアーカイブをダウンロードして展開し、
`chess-sandbox` を自分で `PATH` の通った場所に置くこともできます。

### インストーラーが行うこと

- **Windows**：`chess-sandbox.exe` を `%LOCALAPPDATA%\Programs\chess-sandbox` に置き、そのフォルダをユーザーの `PATH` に追加します。
- **macOS / Linux**：`chess-sandbox` を `~/.local/bin` に置きます（`CHESS_SANDBOX_INSTALL_DIR` で変更できます）。

### アンインストール

**Windows**

1. インストール先のフォルダを削除します。

   ```powershell
   Remove-Item -Recurse "$env:LOCALAPPDATA\Programs\chess-sandbox"
   ```

2. スタートメニューから「**環境変数を編集**」を開き、`Path` を選択して
   `...\Programs\chess-sandbox` のエントリを削除します。

**macOS / Linux**

```sh
rm ~/.local/bin/chess-sandbox
```

保存したセッションも削除する場合は、[保存されるセッション](#保存されるセッション)に書かれた `ChessSandbox` フォルダを削除してください。

## 使い方

```
chess-sandbox [--fen "<FEN>" | --pgn <file>] [--version]
```

| オプション | 説明 |
|---|---|
| `--fen "<FEN>"` | この局面から新しく始める。保存されたセッションより優先される |
| `--pgn <file>` | PGN ファイルの対局を、最初の局面から開く。保存されたセッションより優先される。 |
| `--version` | バージョンを表示 |

オプションを指定しない場合は、標準の初期配置から始めます。保存されたセッションがあるときは、続きから再開する、
標準の初期配置から新しく始める、入力した FEN から新しく始める、のいずれかを選べます。

### 例

```
# 標準の初期配置から始める（保存があれば続きから）
chess-sandbox

# 指定した局面から始める
chess-sandbox --fen "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1"

# PGN ファイルの対局を開く
chess-sandbox --pgn game.pgn
```

## 操作方法

手を SAN（`e4`、`Nf3`、`O-O`、`exd5`、`e8=Q`）で入力して `Enter` を押します。両方の側を自分で指します。
指せない手は拒否され、盤面は変わりません。

| キー | 操作 |
|---|---|
| `←` / `→` | 前の手 / 次の手（入力欄が空のとき） |
| `Home` / `End` | 最初 / 最後の手へ移動（入力欄が空のとき） |
| `Enter` | 入力した手を指す、または入力したコマンドを実行する |
| `Esc` | 保存して終了 |

コマンドは `:` で始めます。

| コマンド | 操作 |
|---|---|
| `:flip`（`:f`） | 盤面を反転 |
| `:fen` | 現在の局面の FEN を表示 |
| `:home` / `:end` | 最初 / 最後の手へ移動 |
| `:pgn <file>` | 現在の手順の全体を PGN ファイルに書き出す |
| `:quit`（`:q`） | 保存して終了 |

手を戻した状態で別の手を指すと、それより先の手は消えます。消える前に確認が出ます
（`y` で続行、それ以外のキーでキャンセル）。

## 保存されるセッション

セッション（開始局面、手の列、現在位置）は、1手動かすたびに保存されます。
ターミナルを閉じても直前の状態が残ります。保存先は以下の JSON ファイルです。

- **Windows**：`%LOCALAPPDATA%\ChessSandbox\session.json`
- **macOS / Linux**：`~/.local/share/ChessSandbox/session.json`

## ライセンス

MIT
