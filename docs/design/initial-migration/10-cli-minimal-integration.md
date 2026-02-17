# 10 CLI最小統合

## 目的

変換パイプラインをCLIから実行可能にする。初期は最小フラグ互換に留める。

## 成果物

- `cmd/azash/main.go`
- 最小オプション解析
- 変換実行フロー統合

## チェックリスト

- [ ] `-i`, `-enc`, `-ext`, `-d` を最優先で対応
- [ ] `-h`, `--help` を実装
- [ ] 入力ファイル複数処理を実装
- [ ] 変換ログ/エラーログを統一
- [ ] 出力ファイル命名規則を実装
- [ ] 既存ディレクトリ互換を壊さない実行導線を確認
- [ ] 後続で追加するフラグの拡張ポイントを設計

## Step 03 からの補足情報

### Config ロード

- `config.LoadFromReader(r)` は `(*Config, []string, error)` を返す。第 2 返値は未知キーの警告リスト
- INI ファイルのパスは CLI フラグまたはデフォルト（jar パス + `AozoraEpub3.ini`）から決定
- INI が存在しない場合は `config.NewConfig()` でデフォルト値を使用

### 辞書ロード順序

辞書ファイルは以下の順序でロードする（Java 版の初期化順序と互換）:

1. `chuki_tag.txt` → `ChukiTagMap`
2. `chuki_tag_suf.txt` → `ChukiSufMap`
3. `chuki_ivs.txt` → `GaijiMap.UtfMap`（最優先）
4. `chuki_utf.txt` → `GaijiMap.UtfMap`（IVS と重複するキーはスキップ）
5. `chuki_alt.txt` → `GaijiMap.AltMap`
6. `chuki_latin.txt` → `LatinMap`
7. `replace.txt` → `ReplaceMap`（ファイル不在時は nil、optional）

## 受け入れ条件

- 必須フラグのみで実用変換が回る。
