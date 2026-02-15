# 12 Web小説変換 (narou含む)

## 目的

`WebAozoraConverter` 相当を移植し、Web小説URLから青空文庫テキストを生成してEPUB変換へ接続する。

対象には、少なくとも `ncode.syosetu.com` (小説家になろう) を含める。

## 成果物

- Web抽出エンジン (URL取得・DOM抽出・置換)
- サイト定義ローダ (`web/*/extract.txt`, `replace.txt`)
- narou向け回帰テスト

## 対象サイト

初期対象:

- `ncode.syosetu.com` (必須)

拡張対象 (段階対応):

- `novel18.syosetu.com`
- `kakuyomu.jp`
- `novelup.plus`
- `novel.fc2.com`
- `novelist.jp`
- `2.novelist.jp`
- `www.akatsuki-novels.com`
- `www.mai-net.net`
- `syosetu.org`

## チェックリスト

- [ ] `web/*/extract.txt` パーサを実装
- [ ] `web/*/replace.txt` パーサを実装
- [ ] URL正規化とサイト判定を実装
- [ ] HTTP取得層 (User-Agent, interval, retry) を実装
- [ ] 目次ページ抽出を実装
- [ ] 各話ページ抽出を実装
- [ ] 本文ノードの青空文庫形式への変換を実装
- [ ] 画像URL/外部リンクの扱いを定義
- [ ] 更新差分取得 (modified only) の方針を実装
- [ ] narou (`ncode.syosetu.com`) のゴールデンテストを作成
- [ ] 取得失敗時のリトライ/中断ポリシーを実装
- [ ] サイト仕様変更検知用の回帰テストを追加

## 互換要件

- Java版の `WebAozoraConverter` で定義される抽出ルールを再利用できること
- 既存 `web/` ディレクトリ資産を大きく変更せず流用できること
- 生成される青空文庫テキストが、既存の本文変換パイプラインでEPUB化できること

## 受け入れ条件

- narou作品URLから青空文庫テキスト生成 → EPUB出力まで一連で成功する
- 代表作品で章順序・本文欠落・重複取得がない
- 主要失敗パターン (アクセス制限、HTML変化、一時エラー) で適切に失敗/再試行できる
