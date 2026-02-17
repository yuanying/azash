# 09 EPUBライタ

## 目的

テンプレート適用とzipパッケージングを行い、EPUB3を生成する。

## 成果物

- EPUB writer
- image-only EPUB writer
- opf/nav/ncx/xhtml/css生成

## チェックリスト

- [ ] `mimetype` をSTOREDで先頭出力
- [ ] `META-INF/container.xml` を配置
- [ ] manifest/spine/nav/toc生成を実装
- [ ] セクション分割と連番命名を実装
- [ ] 画像/フォント/gaiji配置を実装
- [ ] image-only入力の出力を実装
- [ ] EPUB構造検証テストを追加

## Step 03 からの補足情報

### Config スタイルフィールド

- `PageMargin` / `BodyMargin`: カンマ区切り CSV 文字列（例: `"1,2,3,4"`）
  - **4 要素に分割できない場合、デフォルト `"0,0,0,0"` にリセットされる**（Java 版互換）
- `PageMarginUnit` / `BodyMarginUnit`: `"em"` または `"%"`
  - INI の `PageMarginUnit=0` → `"em"`、それ以外（不在含む） → `"%"`
  - unit は PageMargin が 4 要素の場合のみ適用
- `LineHeight`: 行高（デフォルト 1.8）
- `FontSize`: フォントサイズ（デフォルト 100）
- `CoverPage` / `CoverW` / `CoverH`: 表紙設定
- `TitlePage` / `TitlePageWrite`: タイトルページ設定（`TitlePageWrite=0` のとき `TitlePage` はパースされない）
- `TocPage` / `TocVertical` / `CoverPageToc`: 目次ページ設定

## 受け入れ条件

- 生成EPUBが主要リーダで開け、構造検証に通る。
