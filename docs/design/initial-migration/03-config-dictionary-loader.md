# 03 設定/辞書ローダ

## 目的

`AozoraEpub3.ini` と各種辞書ファイルを読み込み、変換設定へ統合する。

## 成果物

- INIローダ
- 辞書ローダ (`chuki_utf`, `chuki_ivs`, `chuki_latin`, `chuki_alt`, `replace`)
- バリデーション/警告ロガー

## チェックリスト

- [x] Java版INIキーの互換マッピング表を作成
- [x] デフォルト値をJava実装準拠で定義
- [x] 未対応キーを警告し処理継続できるようにする
- [x] 数値/真偽値/列挙型の変換エラー処理を統一
- [x] 辞書ファイルのエンコーディング/行形式を吸収
- [ ] 変換器・ライタへ設定反映するアダプタを実装

## 実装済みパッケージ

- `internal/config` — Config 構造体 + INI パーサ (`LoadFromReader`)
- `internal/dict` — 辞書ローダ群
  - `chuki_tag.go` — `ChukiTagMap` / `ChukiSufMap` (chuki_tag.txt, chuki_tag_suf.txt)
  - `gaiji.go` — `GaijiMap` (chuki_utf.txt, chuki_ivs.txt, chuki_alt.txt)
  - `latin.go` — `LatinMap` (chuki_latin.txt)
  - `replace.go` — `ReplaceMap` (replace.txt)

## 未実装（後続ステップで対応）

- 変換器・ライタへの設定反映アダプタ（step 05, 09 で実装予定）

## 受け入れ条件

- 同じINI入力でJava版と同傾向の設定が適用される。
