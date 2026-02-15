# Compatibility Test Data

このディレクトリは、Java版AozoraEpub3とGo版azashの互換比較用データを管理します。

## Directory Layout

- `input/`: 比較対象の入力ファイル (txt/zip/txtz/cbz)
- `golden/java/`: Java版で生成したEPUBゴールデン
- `unpacked/`: EPUB展開/正規化後の比較用ファイル
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

## Example Command (replace as needed)

```bash
java -cp temp/AozoraEpub3/AozoraEpub3.jar AozoraEpub3 \
  -i AozoraEpub3.ini \
  -enc AUTO \
  -ext .epub \
  -d testdata/compat/golden/java \
  testdata/compat/input/<input-file>
```

## Normalization Notes

比較前に以下を正規化対象とする:

- UUID
- modified timestamp
- zip entry order (必要に応じて)

