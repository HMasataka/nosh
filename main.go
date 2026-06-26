package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	title := flag.String("title", "nosh", "通知のタイトル")
	message := flag.String("message", "", "通知メッセージ")
	sound := flag.String("sound", "default", "サウンド名")
	requireTTY := flag.Bool("require-tty", false, "制御端末を持つ対話的セッションのときだけ通知する")
	owner := flag.String("owner", "claude", "--require-tty で制御端末を確認する祖先プロセス名 (空文字なら祖先のいずれかが端末を持てば対話的とみなす)")
	flag.Parse()

	// 引数でもメッセージを受け取れるようにする
	msg := *message
	if msg == "" && flag.NArg() > 0 {
		msg = strings.Join(flag.Args(), " ")
	}

	if msg == "" {
		fmt.Fprintln(os.Stderr, "Usage: nosh [options] <message>")
		fmt.Fprintln(os.Stderr, "       nosh -message <message>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// --require-tty 指定時は、対話的セッションでなければ黙って終了する。
	// 裏で起動された claude (pensieve のフィードバックループ、belt が起こす
	// 外部 CLI など) は制御端末を持たないため、ここで抑制される。
	if *requireTTY && !isInteractive(*owner) {
		return
	}

	if err := notify(*title, msg, *sound); err != nil {
		fmt.Fprintln(os.Stderr, "通知の送信に失敗:", err)
		os.Exit(1)
	}
}

func notify(title, message, sound string) error {
	script := fmt.Sprintf(`display notification %q with title %q sound name %q`, message, title, sound)

	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}

// isInteractive は、呼び出し元が制御端末を持つ対話的セッションかを判定する。
//
// フックから起動された nosh 自身は端末から切り離されている (tty = "??") ため、
// 自分の tty ではなく祖先プロセスをたどって判定する。
//
//   - owner が空でない場合: 祖先のうち名前に owner を含む最も近いプロセス
//     (= このフックを所有する claude) の制御端末を見る。裏 claude は端末を
//     持たない (tty = "??") のでそこで false を返し、対話 claude まで遡って
//     誤判定することはない。owner が祖先に見つからなければ通知する (fail-open)。
//   - owner が空の場合: 祖先のいずれかが制御端末を持てば対話的とみなす。
func isInteractive(owner string) bool {
	pid := os.Getpid()
	for i := 0; i < 16 && pid > 1; i++ {
		ppid, tty, comm := procInfo(pid)
		if owner != "" {
			if strings.Contains(comm, owner) {
				return hasTTY(tty)
			}
		} else if hasTTY(tty) {
			return true
		}
		if ppid <= 1 || ppid == pid {
			break
		}
		pid = ppid
	}
	// owner を持つ祖先が見つからなかった場合は抑制しない。
	return owner != ""
}

// procInfo は ps で指定 PID の親 PID・制御端末・コマンド名を取得する。
func procInfo(pid int) (ppid int, tty, comm string) {
	out, err := exec.Command("ps", "-o", "ppid=,tty=,comm=", "-p", strconv.Itoa(pid)).Output()
	if err != nil {
		return 0, "", ""
	}
	fields := strings.Fields(string(out))
	if len(fields) < 3 {
		return 0, "", ""
	}
	ppid, _ = strconv.Atoi(fields[0])
	tty = fields[1]
	comm = strings.Join(fields[2:], " ")
	return ppid, tty, comm
}

// hasTTY は ps の tty 列が実在の制御端末を指すかを判定する。
// 制御端末を持たないプロセスでは "??" や "?" になる。
func hasTTY(tty string) bool {
	tty = strings.TrimSpace(tty)
	return tty != "" && tty != "??" && tty != "?"
}
