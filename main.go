package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	title := flag.String("title", "nosh", "通知のタイトル")
	message := flag.String("message", "", "通知メッセージ")
	sound := flag.String("sound", "default", "サウンド名")
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
