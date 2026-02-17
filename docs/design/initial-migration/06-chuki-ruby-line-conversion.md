# 06 注記/ルビ/行変換

## 目的

本文変換の中心となる注記展開、ルビ処理、行単位変換を移植する。

## 成果物

- `convertGaijiChuki`
- `replaceChukiSufTag`
- `convertRubyText`
- `convertTextLineToEpub3` の中核

## チェックリスト

- [ ] 注記タグ変換テーブルを構築
- [ ] 前方参照注記処理を実装
- [ ] ルビ範囲判定 (`｜` あり/なし) を実装
- [ ] 縦中横判定 (`checkTcyPrev/Next`) を実装
- [ ] コメント/空行/字下げ関連オプションを適用
- [ ] Javaテスト (`AozoraEpub3ConverterTest`) を移植

## Step 03 からの補足情報

### ChukiTagMap (`internal/dict/chuki_tag.go`)

- `Tags`: 注記キー → `[tag, endTag]` のスライス。`endTag` が空の場合は 1 要素
- フラグマップ（`Flag*`）: 各注記キーの処理モードを `map[string]struct{}` で保持
  - `FlagNoBr` ('1'): 改行抑制
  - `FlagNoRubyStart` ('2'): ルビ開始抑制
  - `FlagNoRubyEnd` ('3'): ルビ終了抑制
  - `FlagPageBreak` ('P','M','L'): 改ページ
  - `FlagMiddle` ('M'): 中央寄せ改ページ
  - `FlagBottom` ('L'): 下寄せ改ページ
  - `FlagKunten` ('K'): 訓点

### ChukiSufMap (`internal/dict/chuki_tag.go`)

- `Tags`: キー → `[startTag, endTag]` のスライス
- 4 列目に別名がある場合、`別名+キー` でも追加登録される

### Config 関連フィールド

- `DakutenType`: 濁点処理方式
- `SpaceHyphenation`: 空白ハイフネーション方式
- `AutoYoko` / `AutoYokoNum1` / `AutoYokoNum3` / `AutoYokoEQ1`: 縦中横自動判定
- `RemoveEmptyLine` / `MaxEmptyLine`: 空行処理
- `CommentPrint` / `CommentConvert`: コメント行処理
- `BoldUseGothic` / `GothicUseBold`: ゴシック/太字変換
- `SameLineChapter`: 同行見出し処理（Step 07 とも関連）

## 受け入れ条件

- 代表的な注記・ルビケースでJava版と同等出力になる。
