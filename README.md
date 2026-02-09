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
```

## オプション

| オプション | デフォルト | 説明 |
|-----------|-----------|------|
| `-title`  | `nosh`    | 通知のタイトル |
| `-message`| (なし)    | 通知メッセージ |
| `-sound`  | `default` | サウンド名 |

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
