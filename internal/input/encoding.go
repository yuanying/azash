package input

import (
	"errors"
	"io"

	"github.com/saintfish/chardet"
)

const (
	detectBufSize = 4096
	maxProbeSize  = 64 * 1024 // 最大64KBまで読み取って判定
)

// DetectEncoding は入力ストリームから文字コードを検出する。
// Java版 Detector.getCharset の移植。4096バイトのチャンクで繰り返し読み取り、
// chardetライブラリで判定する。最大64KBまたはEOFまで読む。
// 検出できない場合は空文字列を返す。
//
// 注意: この関数はストリームを消費する。呼び出し後に同じストリームを
// 本文パースに使用する場合は、別途 Source.OpenText() で再オープンすること。
// Java版も検出用と本文読取用で別ストリームを開く前提の設計。
func DetectEncoding(r io.Reader) (string, error) {
	buf := make([]byte, detectBufSize)
	data := make([]byte, 0, maxProbeSize)

	for len(data) < maxProbeSize {
		n, err := r.Read(buf)
		if n > 0 {
			data = append(data, buf[:n]...)
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
	}

	if len(data) == 0 {
		return "", nil
	}

	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(data)
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
