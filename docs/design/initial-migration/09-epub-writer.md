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

## 受け入れ条件

- 生成EPUBが主要リーダで開け、構造検証に通る。
