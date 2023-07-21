package system

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
)

// 外部コマンド実行用
// C言語のsystem関数と同等の動作を行う
func System(cmd string) int {
	c := exec.Command("sh", "-c", cmd)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()

	if err == nil {
		return 0
	}

	// 終了ステータス(Exit code)を返す
	if ws, ok := c.ProcessState.Sys().(syscall.WaitStatus); ok {
		if ws.Exited() {
			return ws.ExitStatus()
		}

		if ws.Signaled() {
			return -int(ws.Signal())
		}
	}

	return -1
}

func PrintFile(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))
}

func LoadLineText(path string, line int) string {
	fp, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fp.Close() 

	scanner := bufio.NewScanner(fp)
	n := 1
	for scanner.Scan() { // 1行ずつ読み込み
		if n == line {
			return scanner.Text()
		}
		n++
	}
	return ""
}