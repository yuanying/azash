# 02 モデル層

## 目的

Java側の情報モデルをGoで明示的に定義し、変換器とライタの共通契約にする。

## 成果物

- `BookInfo`, `SectionInfo`, `ChapterInfo`, `ChapterLineInfo`, `ImageInfo`, `GaijiInfo`
- 列挙型/定数 (`PageBreakType`, `RubyCharType`, Title種別)

## チェックリスト

- [x] Javaフィールド一覧を洗い出してGo structへ対応付け
- [x] nullable相当フィールドのゼロ値方針を定義
- [x] タイトル種別/章レベル等の定数を定義
- [x] 出力ファイル名生成で必要なヘルパーを実装
- [x] モデルのユニットテストを作成
- [x] 変換途中で更新される可変状態の責務を明文化

## 受け入れ条件

- 主要ユースケースの状態をモデルだけで表現できる。

## 実装済みファイル

| ファイル | 内容 |
|---|---|
| `internal/model/book_info.go` | BookInfo + TitleType enum + タイトルページ定数 |
| `internal/model/section_info.go` | SectionInfo + IMAGE_SIZE_TYPE 定数 |
| `internal/model/chapter.go` | ChapterInfo + ChapterLineInfo + SetTocNestLevel |
| `internal/model/image_info.go` | ImageInfo + Format() |
| `internal/model/gaiji_info.go` | GaijiInfo + FileName() |
| `internal/model/page_break_type.go` | PageBreakType + 画像/ページ種別定数 |
| `internal/model/ruby_char_type.go` | RubyCharType enum |
| `internal/model/cover_edit_info.go` | CoverEditInfo |

## 延期メソッド

以下のメソッドは依存関係のため後続ステップに延期する。

| メソッド | 移行先 | 理由 |
|---|---|---|
| `BookInfo.getFileTitleCreator` | 04 (入力抽象化) | ファイル名からタイトル/著者を抽出する入力処理 |
| `BookInfo.setMetaInfo` / `reloadMetadata` | 07 (章抽出/改ページ) | CharUtils (step 05) に依存 |
| `BookInfo.loadCoverImage` | 08 (画像処理) | ImageUtils.loadImage に依存 |
| `ImageInfo.getImageInfo(File/Stream)` | 08 (画像処理) | 画像フォーマット判定・寸法取得 |

## optional フィールド方針

- 行番号系 (int): Java 版の `-1` センチネル値を維持
- 文字列: Go のゼロ値 `""` を「未設定」として扱う
- `CoverFileName`: `*string` (nil=表紙無し, ""=先頭の挿絵)
- `CoverEditInfo`: `*CoverEditInfo` (nil=編集情報なし)
- `time.Time`: ゼロ値 = 未設定
