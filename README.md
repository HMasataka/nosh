# nosh

macOSでローカル通知を送るシンプルなCLI。

## インストール

```bash
go install github.com/HMasataka/nosh@latest
```

または、リポジトリをクローンしてビルド:

```bash
git clone https://github.com/HMasataka/nosh.git
cd nosh
go install .
```

## 使い方

```bash
# 基本
nosh "メッセージ"

# タイトル付き
nosh -title "ビルド完了" "成功しました"

# サウンド指定
nosh -sound "Ping" "通知です"

# すべてのオプション
nosh -title "タイトル" -sound "Glass" "メッセージ"

# 対話的セッションのときだけ通知する (フック向け)
nosh -require-tty "セッションが終了しました"
```

## オプション

| オプション     | デフォルト | 説明                                                 |
| -------------- | ---------- | ---------------------------------------------------- |
| `-title`       | `nosh`     | 通知のタイトル                                       |
| `-message`     | (なし)     | 通知メッセージ                                       |
| `-sound`       | `default`  | サウンド名                                           |
| `-require-tty` | `false`    | 制御端末を持つ対話的セッションのときだけ通知する     |
| `-owner`       | `claude`   | `-require-tty` で制御端末を確認する祖先プロセス名    |

### `-require-tty` (対話的セッションだけ通知)

エージェント CLI のフックから使うと、裏で起動された CLI セッション
(例: pensieve のフィードバックループ、belt が起こす外部 CLI) でも一律に
通知が鳴ってしまう。`-require-tty` を付けると、制御端末を持つ対話的な
セッション (= 人間がターミナルで操作しているセッション) のときだけ通知する。

判定は祖先プロセスをたどって行う。`-owner` で指定した名前 (既定 `claude`)
を含む最も近い祖先プロセスの制御端末を見て、端末があれば通知、なければ黙る。
裏で起動されたセッションは制御端末を持たないため抑制される。`-owner` を含む
祖先が見つからない場合は通知する (スクリプトからの通常利用を妨げない)。

Claude Code の `settings.json` でのフック例:

```json
{
  "hooks": {
    "Stop": [
      {
        "matcher": ".*",
        "hooks": [
          {
            "type": "command",
            "command": "nosh --require-tty --title 'Claude Code' --message 'セッションが終了しました'"
          }
        ]
      }
    ]
  }
}
```

## 利用可能なサウンド

`/System/Library/Sounds/` にあるサウンドが使用可能:

- Basso
- Blow
- Bottle
- Frog
- Funk
- Glass
- Hero
- Morse
- Ping
- Pop
- Purr
- Sosumi
- Submarine
- Tink

## 使用例

長時間タスクの完了通知:

```bash
make build && nosh "ビルド完了"
```

```bash
go test ./... ; nosh -title "テスト" "完了しました"
```

## 動作環境

- macOS

## License

MIT
