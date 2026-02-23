package input

import (
	"io"

	"github.com/saintfish/chardet"
)

const detectBufSize = 4096

// DetectEncoding は入力ストリームから文字コードを検出する。
// Java版 Detector.getCharset の移植。4096バイトのチャンクで読み取り、
// chardetライブラリで判定する。
// 検出できない場合は空文字列を返す。
func DetectEncoding(r io.Reader) (string, error) {
	buf := make([]byte, detectBufSize)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}
	if n == 0 {
		return "", nil
	}

	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(buf[:n])
	if err != nil {
		return "", nil
	}

	return result.Charset, nil
}

// ResolveEncoding は設定値(encType)と検出結果(detected)から最終的なエンコーディングを決定する。
// encType が "AUTO" の場合は detected を使用し、Shift_JIS は MS932 にマッピングする。
// encType が固定値の場合はそちらを優先する。
// 検出結果が空（検出不能）の場合はフォールバックとして UTF-8 を返す。
func ResolveEncoding(encType, detected string) string {
	if encType != "AUTO" {
		return encType
	}

	if detected == "" {
		return "UTF-8"
	}

	// Shift_JIS → MS932 マッピング (Windows拡張Shift_JIS)
	if detected == "Shift_JIS" {
		return "MS932"
	}

	return detected
}
