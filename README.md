# azash

[AozoraEpub3](https://github.com/kyukyunyorituryo/AozoraEpub3) の Go 移植です。青空文庫形式のテキストファイルを EPUB3 に変換します。

## ステータス

開発初期段階です。現在は Java 版との互換比較ハーネスが稼働しており、移植タスクを順次進めています。

## 必要環境

- Go 1.25 以上
- Java 21 以上（互換比較ハーネス実行時）

## ビルドとテスト

```bash
go test ./...
go vet ./...
gofmt -l -d .
go tool golangci-lint run ./...
```

## 互換比較ハーネス

Java 版の出力をゴールデンとし、Go 版の変換結果との差分を機械的に検証する仕組みです。

```bash
# Java 版ソースを取得
mkdir -p temp
git clone --depth 1 https://github.com/kyukyunyorituryo/AozoraEpub3.git temp/AozoraEpub3

# ハーネスを実行
go run ./cmd/azash-compat run \
  --sample-id sample-006 \
  --java-cp temp/AozoraEpub3/AozoraEpub3.jar \
  --ini temp/AozoraEpub3/AozoraEpub3.ini
```

レポートは `testdata/compat/reports/latest/` に出力されます。

## 移植ロードマップ

詳細は [docs/design/initial-migration/](docs/design/initial-migration/README.md) を参照してください。

| マイルストーン | 内容 |
|---|---|
| M1 | 変換パイプライン開通 + Java 互換リソース読込 + 互換比較ハーネス定常実行 |
| M2 | 注記/ルビ/外字/縦中横のコア互換 |
| M3 | zip/cbz/画像のみを含む実用互換 |

## ライセンス

GPL-3.0 — 詳細は [LICENSE](LICENSE) を参照してください。
