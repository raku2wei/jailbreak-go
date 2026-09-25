PYTHON ?= python3
FONT_DIR := cmd/jailbreak-ebiten/fonts

.PHONY: font wasm serve

# 埋め込みフォント(JailbreakMono = HackGen Console のサブセット)を再生成する。
# 必要: Python 3 + fonttools (`pip install fonttools`)
# 画面に出す文字を追加・変更したら実行し、生成物(chars.txt / JailbreakMono-Regular.ttf)をコミットする。
font:
	go run ./tools/fontchars > $(FONT_DIR)/chars.txt
	$(PYTHON) tools/font/subset.py

# ブラウザ版の WASM をビルドする(-s -w でシンボル・DWARF を落として縮小)
wasm:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o web/jailbreak.wasm ./cmd/jailbreak-ebiten

serve: wasm
	go run ./cmd/serve
