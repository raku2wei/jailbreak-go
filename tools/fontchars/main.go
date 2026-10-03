// fontchars はゲームが画面に表示しうる文字を実装から機械的に収集し、
// 重複なし・コードポイント順の1行テキストとして標準出力へ書き出す。
// 埋め込みフォントのサブセット化(tools/font/subset.py)の入力に使う。
//
// 収集対象:
//   - assets/ 配下の全テキストアセット(AA・問題文・ルール等。*.go を除く)
//   - internal/ pkg/ cmd/ 配下の *.go(テストを除く)の文字列リテラル・rune リテラル
//     (コメントは対象外)
//
// 制御文字(ESC 等)は除外する。ASCII 印字可能文字やかな全体などの
// 「常に含める文字」は subset.py 側で追加する。
//
// 実行: go run ./tools/fontchars > cmd/jailbreak-ebiten/fonts/chars.txt
// (リポジトリのルートで実行すること)
package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {
	chars, err := collect(".")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(chars)
}

// collect はリポジトリのルート root 以下から表示しうる文字を集め、
// コードポイント順に並べた文字列を返す。
func collect(root string) (string, error) {
	set := map[rune]bool{}
	add := func(s string) {
		for _, r := range s {
			if r == utf8.RuneError || unicode.IsControl(r) {
				continue
			}
			set[r] = true
		}
	}

	// テキストアセット
	err := walk(filepath.Join(root, "assets"), func(path string) error {
		if strings.HasSuffix(path, ".go") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		add(string(b))
		return nil
	})
	if err != nil {
		return "", err
	}

	// Go ソースの文字列・rune リテラル
	for _, dir := range []string{"internal", "pkg", "cmd"} {
		err := walk(filepath.Join(root, dir), func(path string) error {
			if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			lits, err := goLiterals(path)
			if err != nil {
				return err
			}
			for _, l := range lits {
				add(l)
			}
			return nil
		})
		if err != nil {
			return "", err
		}
	}

	runes := make([]rune, 0, len(set))
	for r := range set {
		runes = append(runes, r)
	}
	sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })
	return string(runes), nil
}

func walk(root string, fn func(path string) error) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		return fn(path)
	})
}

// goLiterals は Go ソースファイル中の文字列リテラル・rune リテラルの値を返す。
func goLiterals(path string) ([]string, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	file := fset.AddFile(path, fset.Base(), len(src))
	var s scanner.Scanner
	var scanErr error
	s.Init(file, src, func(pos token.Position, msg string) {
		scanErr = fmt.Errorf("%s: %s", pos, msg)
	}, 0) // コメントは読み飛ばす

	var out []string
	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		switch tok {
		case token.STRING:
			v, err := strconv.Unquote(lit)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", path, err)
			}
			out = append(out, v)
		case token.CHAR:
			v, err := strconv.Unquote(lit)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", path, err)
			}
			out = append(out, v)
		}
	}
	return out, scanErr
}
