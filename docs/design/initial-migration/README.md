# AozoraEpub3 → Go (azash) 初期移植計画

このディレクトリは、初期移植計画をタスク単位に分割した設計ドキュメント群です。

## 方針

- 最優先: ディレクトリ/リソース互換 (`template/`, `chuki_*.txt`, `replace*.txt`, `gaiji/`, `presets/`, `web/`)
- 次点: 変換結果互換 (本文/注記/ルビ/目次/改ページ/画像配置)
- CLIは初期フェーズで最小互換に留める
- Java版との差分は機械比較ハーネスで検証する

## タスク一覧

- [00 互換比較ハーネス](./00-compat-harness.md)
- [01 リソース解決層](./01-resource-resolution.md)
- [02 モデル層](./02-model-layer.md)
- [03 設定/辞書ローダ](./03-config-dictionary-loader.md)
- [04 入力抽象化層](./04-input-abstraction.md)
- [05 文字変換コア](./05-text-conversion-core.md)
- [06 注記/ルビ/行変換](./06-chuki-ruby-line-conversion.md)
- [07 章抽出/改ページ](./07-chapter-pagination.md)
- [08 画像処理層](./08-image-processing.md)
- [09 EPUBライタ](./09-epub-writer.md)
- [10 CLI最小統合](./10-cli-minimal-integration.md)
- [11 E2E互換チューニング](./11-e2e-compat-tuning.md)
- [12 Web小説変換 (narou含む)](./12-web-narou-conversion.md)

## マイルストーン

- M1: 変換パイプライン開通 + Java互換リソース読込
- M2: 注記/ルビ/外字/縦中横のコア互換
- M3: zip/cbz/画像のみを含む実用互換
