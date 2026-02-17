# 08 画像処理層

## 目的

画像メタ情報、リサイズ、単ページ化判定などを移植し、本文画像の出力品質を確保する。

## 成果物

- 画像情報リーダ
- 画像変換ユーティリティ
- 画像配置判定ロジック

## チェックリスト

- [ ] `BookInfo.loadCoverImage` の移植 (02-model-layer から延期、ImageUtils 依存)
- [ ] `ImageInfo.getImageInfo(File/Stream)` の移植 (02-model-layer から延期)
- [ ] 画像フォーマット判定を実装
- [ ] 寸法取得と向き判定を実装
- [ ] リサイズ/圧縮品質設定を実装
- [ ] 単ページ化判定を実装
- [ ] 回り込み判定を実装
- [ ] `test_png.zip` でJava版との差分確認

## Step 03 からの補足情報

### Config 画像関連フィールド

- `ImageSizeType`: 画像サイズ方式（デフォルト 2）
- `FitImage` / `SvgImage`: 画像フィット・SVG 出力
- `JpegQuality`: JPEG 圧縮品質（デフォルト 80）
- `ImageScale`: 画像スケール（デフォルト 1.0）
- `RotateImage`: `"1"` → 90°、`"2"` → -90°、それ以外 → 0°（INI 値からの変換済み）
- `ImageFloatType` / `ImageFloatW` / `ImageFloatH`: 回り込み設定
- `SinglePageSizeW` / `SinglePageSizeH` / `SinglePageWidth`: 単ページ化判定閾値

### Config リサイズフィールド

- `ResizeW` / `ResizeH`: 有効時のみ `ResizeNumW` / `ResizeNumH` がパースされる（条件付きロード）

### Config 自動余白フィールド

- `AutoMargin=0` のときサブプロパティはパースされない（条件付きロード）
- `AutoMarginNombreSize`: **Java 版とは意図的に非互換**。Java 版 L183 はタイポで `AutoMarginNombreSize` の値を `autoMarginPadding` に代入し、`nobreSize` は常に 0.03f のまま。Go 版はこのバグを再現せず、各フィールドに正しく値を格納する
- `Gamma`: `Gamma=1` のとき `GammaValue` がパースされる（条件付きロード）

## 受け入れ条件

- 代表サンプルで画像の欠落・崩れがなく出力される。
