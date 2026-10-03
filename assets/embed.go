// Package assets はゲームで使用するテキストアセット(AA・問題文など)を
// go:embed でバイナリに埋め込んで提供する。
// WASM(ブラウザ)実行時はファイルシステムが使えないため、
// ターミナル版・Ebitengine版ともに埋め込みアセットを参照する。
package assets

import (
	"bufio"
	"bytes"
	"embed"
	"strings"
)

//go:embed battle enemy rooms title
var fs embed.FS

// MustRead は "assets/rooms/door0" のような従来のパス表記で
// 埋め込みアセットを読み込む。読み込み失敗時は panic する
// (従来の system.PrintFile と同じ挙動)。
func MustRead(path string) []byte {
	p := strings.TrimPrefix(path, "assets/")
	b, err := fs.ReadFile(p)
	if err != nil {
		panic(err)
	}
	return b
}

// LoadLineText は埋め込みアセットから指定行(1始まり)のテキストを返す。
// 従来の system.LoadLineText の埋め込みアセット版。
func LoadLineText(path string, line int) string {
	scanner := bufio.NewScanner(bytes.NewReader(MustRead(path)))
	n := 1
	for scanner.Scan() { // 1行ずつ読み込み
		if n == line {
			return scanner.Text()
		}
		n++
	}
	return ""
}
