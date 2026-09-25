# jailbreak-go

This is a simple game to escape from the dungeon.<br>
This game runs on UNIX terminal, desktop window (Ebitengine), and web browser (WebAssembly).

## How to play

**ブラウザで遊ぶ: https://raku2wei.github.io/jailbreak-go/**

```sh
docker compose run --rm main
```
```sh
go run ./cmd/jailbreak
```
- Go 1.22 or later is required (Ebitengine v2.8 requirement).

## ブラウザ版 / デスクトップ版 (Ebitengine)

[Ebitengine](https://ebitengine.org/) 製のフロントエンドで、ターミナルと同じテキスト画面を
等幅フォント JailbreakMono(下記「フォントについて」参照)で描画します。
文字は Ebitengine の text パッケージではなく `golang.org/x/image/font` で CPU 描画しています
(text パッケージが含む go-text/typesetting を外して WASM を約3MB小さくするため)。

### デスクトップウィンドウで起動

```sh
go run ./cmd/jailbreak-ebiten
```

### ブラウザ (WebAssembly) で起動

```sh
# 1. WASMバイナリをビルド(`make wasm` でも可。-s -w でシンボル・DWARF を落として縮小)
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o web/jailbreak.wasm ./cmd/jailbreak-ebiten

# 2. ローカルサーバーを起動
go run ./cmd/serve

# 3. ブラウザで http://localhost:8080 を開く
```

### GitHub Pages(公開版)の更新

1. 上記の手順1(`-ldflags="-s -w"` 付き)で `web/jailbreak.wasm` をビルドする(`web/wasm_exec.js` は `$(go env GOROOT)/misc/wasm/wasm_exec.js` と同じバージョンにする)
2. `gh-pages` ブランチに `web/` の `index.html` / `wasm_exec.js` / `jailbreak.wasm` をコピーしてコミットする
3. `git push origin gh-pages` で push すると、数分で https://raku2wei.github.io/jailbreak-go/ に反映される

### フォントについて

埋め込みフォント `cmd/jailbreak-ebiten/fonts/JailbreakMono-Regular.ttf` は、
[白源 / HackGen](https://github.com/yuru7/HackGen) v2.10.0 の HackGen Console Regular
(Copyright (c) 2019, Yuko OTAWARA / SIL Open Font License 1.1)を、ゲームが表示する文字だけに
サブセット化したものです(WASM のサイズ削減のため。元 10.7MB → 約 0.34MB)。

- OFL の Reserved Font Name("白源", "HackGen")に従い、フォント内部名を **JailbreakMono** に変更しています。
  著作権表示(name テーブルの nameID 0)は元のまま残しています。
- ライセンス: `cmd/jailbreak-ebiten/fonts/LICENSE_HackGen`(SIL OFL 1.1 全文。JailbreakMono も同ライセンス)
- 収録文字: 実装(`assets/` のテキスト、`internal/` `pkg/` `cmd/` の Go 文字列リテラル)から機械的に集めた文字
  (`cmd/jailbreak-ebiten/fonts/chars.txt`)に、ASCII・ひらがな・カタカナ・全角英数・罫線・矢印・図形記号等を加えたもの。
  漢字は実装で使っているものだけです。

画面に出す文字を追加・変更したら、フォントを再生成して生成物をコミットしてください
(`go test ./...` がフォントに無い文字や `chars.txt` の更新漏れを検出します)。

```sh
pip install fonttools   # 初回のみ
make font               # 元フォントを公式リリースから取得(sha256 検証)してサブセットを再生成
```

同じ入力からは同一バイナリが生成されます。漢字を JIS 第1水準まで入れたい場合は
`python3 tools/font/subset.py --kanji jis1`(約 1.9MB)。

### 操作方法

- `w` / `s`: メニュー選択・前進/後退(矢印キーも使用可)
- `a` / `d`: 左/右を向く
- `Enter`: 決定
- 戦闘はタイピング(ローマ字入力)

*ブラウザ版はページ内をクリックしてフォーカスを与えてからキー入力してください。*

## Recommended Terminal Settings

- **Font:** SF Mono Regular  
  - *Note: Fonts that visualize spaces are not recommended.*
- **Font Size:** 12
- **Window Size:** At least 80×55

*These settings are provided as a reference and may vary depending on your PC environment. If the display is distorted, please adjust accordingly.*

## License

MIT

## Author

raku2wei
