package dungeon

import (
	"strings"
	"testing"
	"time"

	"jailbreak/internal/console"
	"jailbreak/internal/event"
	"jailbreak/internal/player"
	"jailbreak/internal/room"
)

// screenRows は VTerm の各行を文字列にする(全角の継続マーク R=0 は除き、右側の空白は落とす)。
func screenRows(v *console.VTerm) []string {
	cols, rows := v.Size()
	cells := v.Snapshot(nil)
	out := make([]string, rows)
	for y := 0; y < rows; y++ {
		var b strings.Builder
		for x := 0; x < cols; x++ {
			if r := cells[y*cols+x].R; r != 0 {
				b.WriteRune(r)
			}
		}
		out[y] = strings.TrimRight(b.String(), " ")
	}
	return out
}

func rowsOf(rows []string, a console.Area) []string {
	return rows[a.Top : a.Top+a.Rows]
}

// drawDungeonScreen はゲームループと同じ順(Display → CheckEvent → WaitAction)で通常の探索画面を描く。
// WaitAction のキー入力には無効なキーを渡す(画面下部にエラーが出るだけで状態は変わらない)。
func drawDungeonScreen(v *console.VTerm, d *Dungeon) []string {
	d.Display()
	d.CheckEvent()
	v.PushRune('x')
	d.WaitAction()
	return screenRows(v)
}

// console の区画定義(ViewArea / GuideArea / MapArea)が、実際の探索画面の行割りと一致すること。
// どの部屋・どの向きでも部屋のAAの行数は同じなので、行割りは変わらない。
func TestScreenAreasMatchDungeonScreen(t *testing.T) {
	v := console.NewVTerm(console.ScreenCols, console.ScreenRows)
	console.SetBackend(v)
	d := Create(*player.NewPlayer())

	for x := 1; x <= 6; x++ {
		for y := 1; y <= 4; y++ {
			for dir := room.North; dir <= room.West; dir++ {
				d.player.SetPosition(x, y, dir)
				rows := drawDungeonScreen(v, d)

				view := rowsOf(rows, console.ViewArea)
				if !strings.Contains(view[0], "＼") || !strings.Contains(view[len(view)-1], "＼") {
					t.Fatalf("(%d,%d,%d) 描写エリアの先頭・末尾が部屋の壁になっていない:\n%q\n%q", x, y, dir, view[0], view[len(view)-1])
				}
				for r := console.ViewArea.Top + console.ViewArea.Rows; r < console.GuideArea.Top; r++ {
					if rows[r] != "" {
						t.Fatalf("(%d,%d,%d) %d 行目は空行のはず: %q", x, y, dir, r, rows[r])
					}
				}
				if !strings.HasPrefix(rows[console.GuideArea.Top], "w：前に進む") {
					t.Fatalf("(%d,%d,%d) %d 行目が操作説明になっていない: %q", x, y, dir, console.GuideArea.Top, rows[console.GuideArea.Top])
				}
				m := rowsOf(rows, console.MapArea)
				if m[0] != "" || m[len(m)-4] == "" || m[len(m)-3] != "北" || m[len(m)-1] != "　→東" {
					t.Fatalf("(%d,%d,%d) マップの行割りがずれている: %q", x, y, dir, m)
				}
			}
		}
	}
}

// waitFor は cond を満たす画面になるまで待ち、その画面を返す。
func waitFor(t *testing.T, v *console.VTerm, what string, cond func([]string) bool) []string {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if rows := screenRows(v); cond(rows) {
			return rows
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("%s の画面にならない:\n%s", what, strings.Join(screenRows(v), "\n"))
	return nil
}

// 戦闘は一人称視点の描写エリアの中で行われ、戦闘中も下部にマップ(戦闘開始時点のもの)が残ること。
// 戦闘後は部屋が描き直されてから操作説明・マップが出ること(9/25 に直した不具合の再発防止)。
func TestBattleInViewAreaKeepsMap(t *testing.T) {
	v := console.NewVTerm(console.ScreenCols, console.ScreenRows)
	console.SetBackend(v)
	d := Create(*player.NewPlayer())

	// ゲームループと同じく部屋を描いてから、移動直後としてイベント判定する(必ずエンカウントさせる)
	d.Display()
	roomRows := rowsOf(screenRows(v), console.ViewArea)
	d.moveRoom = true
	d.encounterRate = 1000

	done := make(chan event.Event, 1)
	go func() { done <- d.CheckEvent() }()

	// (1) 登場表示: 描写エリアの上下中央(3行 → 10 行目から)に、左右中央で出る
	intro := waitFor(t, v, "登場表示", func(r []string) bool { return strings.Contains(r[10], "あらわれた") })
	if got, want := intro[10], strings.Repeat(" ", 28)+"警備員があらわれた！"; got != want {
		t.Errorf("登場表示: got %q, want %q", got, want)
	}
	if !strings.HasPrefix(intro[console.GuideArea.Top], "ローマじを 入力して 警備員を たおせ！") {
		t.Errorf("操作説明の行が戦闘用の案内になっていない: %q", intro[console.GuideArea.Top])
	}
	mapRows := rowsOf(intro, console.MapArea)
	if mapRows[len(mapRows)-4] == "" || mapRows[len(mapRows)-1] != "　→東" {
		t.Fatalf("登場表示の時点でマップが出ていない: %q", mapRows)
	}
	for y := console.ViewArea.Top; y < console.ViewArea.Top+console.ViewArea.Rows; y++ {
		if y != 10 && y != 12 && intro[y] != "" {
			t.Errorf("登場表示で描写エリアの %d 行目に部屋が残っている: %q", y, intro[y])
		}
	}

	// (2) 戦闘中: 問題文(12〜13 行目)が描写エリアに出て、マップはそのまま
	v.PushRune('x') // 開始(任意のキー)
	battle := waitFor(t, v, "戦闘中", func(r []string) bool { return strings.Contains(r[16], "残り時間") })
	romaji := strings.TrimSpace(battle[13])
	if romaji == "" || strings.TrimLeft(battle[13], " ") != romaji {
		t.Fatalf("ローマ字の行(13 行目)が読めない: %q", battle[13])
	}
	if !strings.Contains(battle[4], "∧") || !strings.Contains(battle[12], "「") {
		t.Errorf("敵のAA(4 行目)・問題文(12 行目)の位置がずれている: %q / %q", battle[4], battle[12])
	}
	if got := rowsOf(battle, console.MapArea); strings.Join(got, "\n") != strings.Join(mapRows, "\n") {
		t.Errorf("戦闘中にマップが変わった:\n got %q\nwant %q", got, mapRows)
	}
	if battle[console.GuideArea.Top] != intro[console.GuideArea.Top] {
		t.Errorf("戦闘中に案内の行が変わった: %q", battle[console.GuideArea.Top])
	}

	for _, r := range romaji {
		v.PushRune(r)
	}
	select {
	case ev := <-done:
		if ev != event.NoEvent {
			t.Fatalf("勝ったのにイベント %v が返った", ev)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("戦闘が終わらない")
	}

	// (3) 戦闘後: ゲームループは続けて WaitAction を呼ぶ。部屋が描き直された状態で操作説明とマップが出る
	v.PushRune('x')
	d.WaitAction()
	after := screenRows(v)
	if got := rowsOf(after, console.ViewArea); strings.Join(got, "\n") != strings.Join(roomRows, "\n") {
		t.Errorf("戦闘後に部屋が描き直されていない:\n%s", strings.Join(got, "\n"))
	}
	if !strings.HasPrefix(after[console.GuideArea.Top], "w：前に進む") {
		t.Errorf("戦闘後の操作説明の行: %q", after[console.GuideArea.Top])
	}
	if got := rowsOf(after, console.MapArea); strings.Join(got, "\n") != strings.Join(mapRows, "\n") {
		t.Errorf("戦闘中のマップが戦闘後(通常時)のマップと違う:\n got %q\nwant %q", mapRows, got)
	}
}
