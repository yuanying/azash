# Compatibility Test Data

このディレクトリは、Java版AozoraEpub3とGo版azashの互換比較用データを管理します。

## Directory Layout

- `input/`: 比較対象の入力ファイル (txt/zip/txtz/cbz)
- `golden/java/`: Java版で生成したEPUBゴールデン
- `output/go/`: Go版で生成したEPUB
- `unpacked/`: EPUB展開/正規化後の比較用ファイル
  - `unpacked/<sample-id>/java/raw`
  - `unpacked/<sample-id>/java/normalized`
  - `unpacked/<sample-id>/go/raw`
  - `unpacked/<sample-id>/go/normalized`
- `reports/latest/`: 比較レポート
  - `summary.md`
  - `report.json`
  - `logs/`
- `samples.csv`: サンプル管理台帳

## Golden Generation Policy

- Java版をオラクルとしてゴールデンを生成する
- 入力ファイルは原則immutable (差し替え禁止)
- 変換条件を固定して再現性を担保する

## Java Execution Record Template

以下を実行ごとに追記:

- Date:
- Java version:
- AozoraEpub3 version:
- Command:
- INI file:
- Input set:
- Output dir:
- Notes:

## Commands

### Java only (manual example)

```bash
java -cp temp/AozoraEpub3/AozoraEpub3.jar AozoraEpub3 \
  -i AozoraEpub3.ini \
  -enc AUTO \
  -ext .epub \
  -d testdata/compat/golden/java \
  testdata/compat/input/<input-file>
```

### Harness commands

```bash
# Java golden generation
go run ./cmd/azash-compat golden-java

# Compare existing java/go outputs
go run ./cmd/azash-compat compare

# Generate + compare in one shot
go run ./cmd/azash-compat run
```

## Diff Policy

- 初期フェーズは非ゲート運用。
- 差分があってもコマンドは成功扱い。
- 失敗扱いは、生成/展開/正規化/比較/レポートの処理失敗のみ。
- 将来ゲート化する場合は `--fail-on-diff=true` を指定する。

## Normalization Notes

比較前に以下を正規化対象とする:

- UUID (`UUID_NORMALIZED` へ置換)
- 日時 (`DATE_NORMALIZED` へ置換)
- zip entry order (比較前のファイル列挙順で吸収)
- XML属性順と空白のゆらぎ (意味を維持して正規化)

## Artifact policy

CIに載せる場合、最低限次を保存する:

- `testdata/compat/reports/latest/summary.md`
- `testdata/compat/reports/latest/report.json`
- `testdata/compat/unpacked/<sample-id>/*/normalized/`
