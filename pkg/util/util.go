package util

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	HEADERS = map[string]string{
		"Accept-Encoding":    "gzip, deflate, br",
		"Accept-Language":    "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
		"Sec-Ch-Ua":          `"Chromium";v="118", "Google Chrome";v="118", "Not=A?Brand";v="99"`,
		"Sec-Ch-Ua-Mobile":   "?0",
		"Sec-Ch-Ua-Platform": "Linux",
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-origin",
		//nolint
		"User-Agent":       `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36`,
		"X-Requested-With": "XMLHttpRequest",
	}

	// Web site
	// FIXME Due to the change of web site, please configure it with file.

	THY_URL_ORIGIN  string = "https://cloud.ssthy.us/"
	THY_URL_LOGIN   string = "https://cloud.ssthy.us/auth/login"
	THY_URL_CHECKIN string = "https://cloud.ssthy.us/user/checkin"

	CUTECLOUD_URL_ORIGIN  string = "https://www.cute-cloud.top"
	CUTECLOUD_URL_LOGIN   string = "https://www.cute-cloud.top/auth/login"
	CUTECLOUD_URL_CHECKIN string = "https://www.cute-cloud.top/user/checkin"

	CHECKIN_HEADER_ACCEPT     string = "application/json, text/javascript, */*; q=0.01"
	LOGIN_HEADER_ACCEPT       string = "application/json, text/javascript, */*; q=0.01"
	LOGIN_HEADER_CONTENT_TYPE string = "application/x-www-form-urlencoded; charset=UTF-8"
	LOGIN_HEADER_METHOD       string = "POST"
	CHECKIN_HEADER_METHOD     string = "POST"
	HEADER_CONTENT_LENGTH     string = "0"

	CUTECLOUD_LOGIN_HEADER_ACCEPT string = "*/*;"

	DELEMITER string = "@"
)

/**
 * TAG Read config file
 *	email: first line
 *	passwd: second line
 */
func ReadConfigFromFile(path string) (string, string, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		// fmt.Println("Read config file error: ", err)
		log.Println("Read config file error: ", err)
		return "", "", err
	}
	buf_str := strings.Split(string(buf), "\n")
	email, passwd := buf_str[0], buf_str[1]

	return email, passwd, nil
}

func NotifySend(title string, level string, body string) {
	switch runtime.GOOS {
	case "linux":
		_ = exec.Command("notify-send", "-u", level, title, body).Run()
	case "darwin": // MAC
		str := "display notification \"" + body + "\" with title \"" + title + "\""
		_ = exec.Command("osascript", "-e", str).Run()
	case "windows":
		panic("Not implemented on Windows")
	default:
		panic("Unsupported OS")
	}
	log.Println(title, body)
}
