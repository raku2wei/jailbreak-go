# CLAUDE.md

このリポジトリで作業する AI エージェント(Claude Code 等)向けのガイド。
遊び方・操作方法・推奨ターミナル設定などの利用者向け情報は [README.md](README.md) を参照。

## プロジェクト概要

Go 製のテキストベースのダンジョン脱出ゲーム。同じゲームロジックを 3 つのフロントエンドで動かす。

- ターミナル版: `cmd/jailbreak`(標準出力 + go-tty。UNIX ターミナル専用)
- デスクトップ版: `cmd/jailbreak-ebiten`([Ebitengine](https://ebitengine.org/) のウィンドウに仮想端末を描画)
- ブラウザ版: `cmd/jailbreak-ebiten` を WebAssembly でビルドしたもの(GitHub Pages で公開)

Go のバージョンは `go.mod` の指定に従う(Ebitengine v2.8 の要件で Go 1.22 以上)。

## ディレクトリ構成

```
cmd/
  jailbreak/          ターミナル版エントリポイント(//go:build !js)
  jailbreak-ebiten/   Ebitengine 版エントリポイント(デスクトップ / WASM 共通)
    fonts/            埋め込みフォント HackGen Console(SIL OFL 1.1)
  serve/              WASM 版をローカル配信する簡易 HTTP サーバー(//go:build !js)
internal/
  console/            入出力の抽象化層(Backend / TerminalBackend / VTerm)
  game/               ゲーム進行の本体
  dungeon/            ダンジョンマップ・移動・エンカウント・罠
  room/ enemy/ player/ event/ title/  部屋描画・戦闘・プレイヤー・イベント・タイトル
pkg/system/           旧来のターミナル補助関数(ターミナル版のみが使用)
assets/               AA・問題文などのテキストアセット(go:embed で埋め込み)
web/                  index.html と wasm_exec.js(jailbreak.wasm はビルド成果物で git 管理外)
```

## ビルド・実行

実行方法の詳細は README の「How to play」「ブラウザ版 / デスクトップ版」を参照。要点のみ:

```sh
go run ./cmd/jailbreak            # ターミナル版
go run ./cmd/jailbreak-ebiten     # デスクトップ版

# WASM 版のビルドとローカル配信(http://localhost:8080)
GOOS=js GOARCH=wasm go build -o web/jailbreak.wasm ./cmd/jailbreak-ebiten
go run ./cmd/serve
```

GitHub Pages(公開版)の更新手順は README の「GitHub Pages(公開版)の更新」を参照。
`gh-pages` ブランチに `web/` の `index.html` / `wasm_exec.js` / `jailbreak.wasm` をコピーして push する。

## 検証コマンド

変更後は以下をすべて通すこと(CI `.github/workflows/ci.yml` と同じ内容)。

```sh
gofmt -l .                                         # 出力が空であること
go vet ./...
go build ./...
go test ./...
GOOS=js GOARCH=wasm go build -o /dev/null ./cmd/jailbreak-ebiten
```

Linux でネイティブビルドする場合は Ebitengine の依存パッケージ
(`libgl1-mesa-dev xorg-dev libasound2-dev` 等)が必要。

## 設計上の約束

- ゲームロジックは `internal/` に置き、ターミナル版と Ebitengine 版で共通にする。フロントエンド固有のコードを持ち込まない。
- ゲームロジックからの入出力は必ず `internal/console` 経由(`console.Printf` / `console.Println` / `console.Clear` / `console.ReadKey` 等)。`fmt.Printf` や go-tty を直接使わない。
- 各エントリポイントは起動時に `console.SetBackend` でバックエンドを設定する。
  - ターミナル版: `TerminalBackend`(実ターミナルへ出力)
  - Ebitengine 版: `VTerm`(仮想端末)
- ターミナル専用のコード(go-tty など WASM でビルドできないもの)はファイル先頭に `//go:build !js` を付けて分離する。
- `VTerm` はゲームが出力する ANSI エスケープ(SGR: 0,1,7,30-37,39,40-47,49 など)を解釈してセル格子(`[]Cell`)に反映する。Ebitengine 側は毎フレーム `Snapshot()` を描画するだけにする。全角文字は 2 セルを占有し、2 セル目は `R=0` の継続マーク。
- アセットは `assets/` に置き `go:embed` で埋め込む(WASM ではファイルシステムが使えないため)。読み込みは `assets.MustRead` / `assets.LoadLineText` を使う。

## ハマりどころ

- 全角スペース(U+3000): HackGen Console は全角スペースを点線の四角で描くため、Ebitengine 版の描画では半角スペースとともに描画をスキップしている。
- 文字幅: 判定を実行環境のロケール(`LANG` 等)に依存させない。★ などの曖昧幅文字は 2 セル扱い(HackGen Console の字形と AA のレイアウトに合わせる)。ロケール依存にすると環境によって表示がずれる。
- Ebitengine のキー入力: `inpututil.IsKeyJustPressed` は押された 1 フレームだけ true になる。`Update` で毎フレーム判定し `VTerm.PushRune` で入力チャネルに積むこと(判定を間引いたり条件付きにすると取りこぼす)。
- `web/wasm_exec.js` はビルドに使う Go のバージョンと一致させる。場所は Go 1.23 以前が `$(go env GOROOT)/misc/wasm/wasm_exec.js`、Go 1.24 以降は `$(go env GOROOT)/lib/wasm/wasm_exec.js`。不一致だとブラウザで起動に失敗する。
- WASM バイナリは約 24.8MB。主因は埋め込みフォント(HackGen Console 約 10MB)。サイズ削減を検討するならまずフォント(サブセット化等)を見る。

## コミット規約

- コミットメッセージは日本語で、変更内容が分かる 1 行目を書く。
- コミットの作者は `raku2wei` にする。
- `gofmt` 済みで、上記の検証コマンドが通る状態でコミットする。
