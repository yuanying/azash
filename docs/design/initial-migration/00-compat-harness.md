# 00 互換比較ハーネス

## 目的

Java版をオラクルにして、Go版との差分を機械比較できる基盤を最初に作る。

## 成果物

- `cmd/azash-compat` (`golden-java` / `compare` / `run`)
- Java出力ゴールデン生成処理
- EPUB展開・正規化ツール (UUID/日時/zip順序の正規化)
- XML/HTML/CSS差分レポート (Markdown + JSON)
- 差分カテゴリ分類 (`ruby/chuki/chapter/image/encoding/metadata/packaging/unknown`)

## チェックリスト

- [x] `test_data/*.txt` と `test_png.zip` を対象にゴールデン生成対象を定義
- [x] 追加の青空文庫サンプルを収集して管理台帳を作成 (URL/取得日/文字コード/用途)
- [x] Java版実行コマンドを固定化 (入力/INI/出力先)
- [x] EPUB展開処理を実装
- [x] 非決定要素の正規化ルールを実装
- [x] XML/HTML/CSS比較ルールを定義
- [x] 差分レポートをCIで保存できる形にする
- [x] 失敗時に「どの変換機能差分か」を分類表示する

## CLI仕様 (初期)

- `azash-compat golden-java`: Javaゴールデン生成
- `azash-compat compare`: Java/Go EPUB比較
- `azash-compat run`: 生成 + 比較の一括実行

主要フラグ:

- `--samples-file` (default: `testdata/compat/samples.csv`)
- `--input-dir` (default: `testdata/compat/input`)
- `--golden-dir` (default: `testdata/compat/golden/java`)
- `--go-out-dir` (default: `testdata/compat/output/go`)
- `--work-dir` (default: `testdata/compat/unpacked`)
- `--report-dir` (default: `testdata/compat/reports/latest`)
- `--sample-id` (複数指定可)
- `--java-cmd` / `--java-cp` / `--java-main` / `--ini`
- `--fail-on-diff` (default: `false`)

## 差分判定ポリシー

- 初期E2Eは非ゲート運用。
- 差分があっても失敗扱いにしない。
- 失敗扱いにするのは以下のみ:
  - Java生成失敗
  - EPUB展開失敗
  - 正規化失敗
  - 比較処理失敗
  - レポート生成失敗
- 将来ゲート化時は `--fail-on-diff=true` を使用する。

## レポート仕様

- `testdata/compat/reports/latest/summary.md`
- `testdata/compat/reports/latest/report.json`
- 比較用中間成果物:
  - `testdata/compat/unpacked/<sample-id>/java/{raw,normalized}`
  - `testdata/compat/unpacked/<sample-id>/go/{raw,normalized}`

`report.json` 主要項目:

- `run` (開始/終了時刻、mode、version)
- `sample_count`
- `samples[]` (sampleごとのstatus, diff_summary, failures)
- `diff_summary` (files/node/text/attr)
- `categories`
- `failures[]`

## 人手実施手順 (サンプル収集込み)

### 1. サンプル方針を決める

- ルビ/注記が多い作品
- 挿絵あり (zip同梱画像) 作品
- 長編 (章分割・改ページ確認)
- 記号/外字が多い作品
- 文字コード差分確認用 (UTF-8/BOM, Shift_JIS)

### 2. サンプルをダウンロードして固定化する

- 青空文庫から対象作品のテキスト(zip/txt)を取得
- `testdata/compat/input/` に保存
- ファイル名は一意にする (`aozora_<作品ID>_<variant>.zip` など)
- 再取得時に差し替えないよう、原本は更新禁止にする

### 3. 収集台帳を作る

`testdata/compat/samples.csv` に次を記録:

- 作品名
- 取得元URL
- 取得日
- 入力ファイル名
- 想定文字コード
- 用途タグ (ruby/chuki/image/chapter/encoding)

### 4. Java版ゴールデンを生成する

- Java実行環境を準備
- 入力・INI・出力先を固定してJava版でEPUB生成
- 出力を `testdata/compat/golden/java/` に保存
- 実行コマンドとJavaバージョンを `testdata/compat/README.md` に記録

### 5. 正規化して比較可能状態にする

- EPUBを展開して `testdata/compat/unpacked/` に配置
- UUID/日時/zipエントリ順の差分を正規化
- 正規化後の比較結果をレポート化

### 6. 回帰運用に組み込む

- 新規サンプル追加時は台帳更新を必須化
- 既存サンプルを置換する場合は理由を記録
- CIで比較結果をアーティファクト保存

## 受け入れ条件

- 同一入力に対し、Java/Go差分を再現性高く表示できる。
- 人手でEPUBを開かなくても差分原因を追える。
- 初期フェーズでは「差分ゼロ」ではなく「差分を安定して可視化できる」ことを完了条件とする。
