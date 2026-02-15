# 02 モデル層

## 目的

Java側の情報モデルをGoで明示的に定義し、変換器とライタの共通契約にする。

## 成果物

- `BookInfo`, `SectionInfo`, `ChapterInfo`, `ChapterLineInfo`, `ImageInfo`, `GaijiInfo`
- 列挙型/定数 (`PageBreakType`, `RubyCharType`, Title種別)

## チェックリスト

- [ ] Javaフィールド一覧を洗い出してGo structへ対応付け
- [ ] nullable相当フィールドのゼロ値方針を定義
- [ ] タイトル種別/章レベル等の定数を定義
- [ ] 出力ファイル名生成で必要なヘルパーを実装
- [ ] モデルのJSONスナップショットテストを作成
- [ ] 変換途中で更新される可変状態の責務を明文化

## 受け入れ条件

- 主要ユースケースの状態をモデルだけで表現できる。
