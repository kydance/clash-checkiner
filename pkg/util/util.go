package util

import (
	"os/exec"
	"runtime"

	"github.com/kydance/ziwi/pkg/log"
)

// SendNotify sends notification to system
func SendNotify(title string, level string, body string) {
	switch runtime.GOOS {
	case "linux": // Linux
		_ = exec.Command("notify-send", "-u", level, title, body).Run()
	case "darwin": // MAC
		str := "display notification \"" + body + "\" with title \"" + title + "\""
		_ = exec.Command("osascript", "-e", str).Run()
	case "windows":
		panic("Not implemented on Windows")
	default:
		panic("Unsupported OS")
	}

	log.Infoln(title, body)
}
