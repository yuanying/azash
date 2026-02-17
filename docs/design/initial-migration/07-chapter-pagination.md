# 07 章抽出/改ページ

## 目的

目次構築に必要な章抽出と、強制改ページロジックを移植する。

## 成果物

- 章抽出ロジック
- 改ページ判定ロジック
- TOC生成用データ整形

## チェックリスト

- [ ] `BookInfo.setMetaInfo` / `reloadMetadata` の移植 (02-model-layer から延期、CharUtils 依存)
- [ ] 章見出しパターン設定を実装
- [ ] 章レベル/除外/連番抑制オプションを実装
- [ ] 強制改ページサイズ条件を実装
- [ ] 空行起点改ページ条件を実装
- [ ] `test_data/test_chapter.txt` で差分確認
- [ ] TOC情報がEPUBライタへ渡る形に整理

## Step 03 からの補足情報

### Config 章設定フィールド

- `ChapterSection`: **キー不在時 `true`**（Java 版 L236 互換の特殊デフォルト）
- `ChapterH` / `ChapterH1` / `ChapterH2` / `ChapterH3`: 見出しレベル別の章判定
- `ChapterNameLength`: 章名の最大長（デフォルト 64）
- `ChapterExclude` / `ChapterUseNextLine` / `SameLineChapter`: 章抽出オプション
- `ChapterName` / `ChapterNumOnly` / `ChapterNumTitle` / `ChapterNumParen`: 章番号判定
- `ChapterNumParenTitle`: Java 版のタイポ `hapterNumParenTitle` も互換キーとして認識
- `ChapterPattern` / `ChapterPatternText`: 正規表現による章パターン（`ChapterPattern=0` のとき `ChapterPatternText` は無視）
- `TitleToc`: タイトルを目次に含めるか
- `NavNest` / `NcxNest`: 目次階層化

### Config 改ページフィールド

- `PageBreakSize` / `PageBreakEmptySize` / `PageBreakChapterSize`: **内部はバイト単位**（INI は KB 値 × 1024）
- `PageBreak=0` のときサブプロパティ（`PageBreakSize` 等）はパースされない（条件付きロード）
- `PageBreakEmpty=0` のとき `PageBreakEmptyLine` / `PageBreakEmptySize` はパースされない

## 受け入れ条件

- 章分割と目次構造が実用上Java版と同等。
