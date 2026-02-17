# 05 文字変換コア

## 目的

外字、ラテン、JIS/文字種判定などの純粋変換ロジックを先行移植する。

## 成果物

- 外字変換モジュール
- ラテン変換モジュール
- 文字種/補助変換ユーティリティ

## チェックリスト

- [ ] `AozoraGaijiConverter` 相当を移植
- [ ] `LatinConverter` 相当を移植
- [ ] `JisConverter` / `CharUtils` 必要部を移植
- [ ] IVS/BMP/SSP出力オプションを実装
- [ ] 既存JavaテストケースをGoへ移植
- [ ] 文字単位ベンチマークを作成

## Step 03 からの補足情報

### GaijiMap (`internal/dict/gaiji.go`)

- `UtfMap` / `AltMap` ともに**重複キーは先勝ち**（Java の `loadChukiFile` 互換）
- ロード順序で優先度を制御: IVS → UTF → ALT の順に読むことで IVS 優先を実現
- `AltMap` は代替文字列（HTML タグを含む場合あり）を保持

### LatinMap (`internal/dict/latin.go`)

- `DecompToChar`: 分解表記文字列 → 拡張ラテン文字（rune）
- `CharToCID`: ラテン文字 → `[横CID, 縦CID]`
- 分解表記が空の行は `DecompToChar` に登録されない（CID のみ登録）

### ReplaceMap (`internal/dict/replace.go`)

- `Single` / `Double` の分類は **UTF-16 コードユニット数**で判定（Java の `String.length()` 互換）
- BMP 外文字（サロゲートペア = 2 ユニット）は `Double` に分類される
- `splitTabJava` で Java の `split("\t")` と互換のタブ分割を行う（末尾空要素を落とす）

### Config (`internal/config/config.go`)

- `IvsBMP` / `IvsSSP`: IVS 出力オプション。外字変換時の出力制御に使用

## 受け入れ条件

- 外字/ラテン変換の主要ケースでJava期待値と一致する。
