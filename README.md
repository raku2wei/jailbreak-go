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
等幅フォント([HackGen Console](https://github.com/yuru7/HackGen) / SIL OFL 1.1)で描画します。

### デスクトップウィンドウで起動

```sh
go run ./cmd/jailbreak-ebiten
```

### ブラウザ (WebAssembly) で起動

```sh
# 1. WASMバイナリをビルド
GOOS=js GOARCH=wasm go build -o web/jailbreak.wasm ./cmd/jailbreak-ebiten

# 2. ローカルサーバーを起動
go run ./cmd/serve

# 3. ブラウザで http://localhost:8080 を開く
```

### GitHub Pages(公開版)の更新

1. 上記の手順1で `web/jailbreak.wasm` をビルドする(`web/wasm_exec.js` は `$(go env GOROOT)/misc/wasm/wasm_exec.js` と同じバージョンにする)
2. `gh-pages` ブランチに `web/` の `index.html` / `wasm_exec.js` / `jailbreak.wasm` をコピーしてコミットする
3. `git push origin gh-pages` で push すると、数分で https://raku2wei.github.io/jailbreak-go/ に反映される

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
