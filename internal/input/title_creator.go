package input

import (
	"regexp"
	"strings"
)

var (
	// 拡張子除去パターン: 英数字の拡張子を2回除去 (Java版互換)
	reExt = regexp.MustCompile(`\.([A-Za-z0-9])+$`)

	// 青空文庫パターン: (青空...)
	reAozora = regexp.MustCompile(`\(青空[^)]*\)`)

	// 校正情報パターン: (...校正|軽量|表紙|挿絵|補正|修正|ルビ|Rev|rev...)
	reCorrection = regexp.MustCompile(`\([^)]*(校正|軽量|表紙|挿絵|補正|修正|ルビ|Rev|rev)[^)]*\)`)

	// パターン1: [著者名] タイトル or ［著者名］ タイトル
	reBracketAuthor = regexp.MustCompile(`[[\[［](.+?)[\]］][ |　]*(.*)[ |　]*$`)

	// パターン2: タイトル (著者 or タイトル （著者
	reParenTitle = regexp.MustCompile(`^(.*?)( |　)*(\(|（)`)
)

// FileTitleCreator はファイル名からタイトルと著者名を抽出する。
// Java版 BookInfo.getFileTitleCreator の移植。
func FileTitleCreator(fileName string) (title, creator string) {
	// 拡張子を2回除去 (二重拡張子対応)
	noExtName := reExt.ReplaceAllString(fileName, "")
	noExtName = reExt.ReplaceAllString(noExtName, "")

	// 全角括弧を半角に正規化
	noExtName = strings.ReplaceAll(noExtName, "（", "(")
	noExtName = strings.ReplaceAll(noExtName, "）", ")")

	// 青空文庫パターン除去
	noExtName = reAozora.ReplaceAllString(noExtName, "")

	// 校正情報パターン除去
	noExtName = reCorrection.ReplaceAllString(noExtName, "")

	// パターン1: [著者名] タイトル
	if m := reBracketAuthor.FindStringSubmatch(noExtName); m != nil {
		title = strings.TrimSpace(m[2])
		creator = strings.TrimSpace(m[1])
	} else if m := reParenTitle.FindStringSubmatch(noExtName); m != nil {
		// パターン2: タイトル (著者名
		title = strings.TrimSpace(m[1])
	} else {
		// パターン3: 全体をタイトル
		title = strings.TrimSpace(noExtName)
	}

	return title, creator
}
