package checkin

const (
	CheckinHeaderAccept string = "application/json, text/javascript, */*; q=0.01"

	LoginHeaderAccept string = "application/json, text/javascript, */*; q=0.01"

	LoginHeaderContentType string = "application/x-www-form-urlencoded; charset=UTF-8"
	HeaderMethod           string = "POST"
	HeaderContentLength    string = "0"

	CutecloudLoginHeaderAccept string = "*/*;"
)

var Headers = map[string]string{
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
