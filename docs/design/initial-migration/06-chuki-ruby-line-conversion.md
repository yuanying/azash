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

## 受け入れ条件

- 代表的な注記・ルビケースでJava版と同等出力になる。
