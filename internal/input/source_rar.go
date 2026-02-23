package input

// RAR対応方針:
//
// 青空文庫は主にzip形式で配布されており、RAR形式は極めて稀である。
// 初期実装ではスタブとし、Open("xxx.rar") は ErrRarNotSupported を返す。
//
// 後続対応として github.com/nwaples/rardecode/v2 を使用して
// zipSource と同等のインタフェースで実装予定。
// Source インタフェースにより、追加しても他コードへの影響は最小限。
