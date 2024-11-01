package checkin

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/kydance/ziwi/pkg/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"

	"Checkiner/pkg/util"
)

type WebInfo struct {
	Whoami string

	OriginURL  string
	LoginURL   string
	CheckinURL string

	email  string
	passwd string
}

func NewWebInfo(whoami, originURL, loginURL, checkinURL, email, passwd string) *WebInfo {
	return &WebInfo{
		Whoami: whoami,

		OriginURL:  originURL,
		LoginURL:   loginURL,
		CheckinURL: checkinURL,

		email:  email,
		passwd: passwd,
	}
}

type Checkiner struct {
	*WebInfo

	LoginHeaderAccpet      string
	LoginHeaderContentType string
	HeaderMethod           string

	// flag
	FlagCheckined bool
	LastDay       int
}

func NewCheckiner(loginHeaderAccpet, loginHeaderContentType,
	loginHeaderMethod string, webInfo *WebInfo,
) *Checkiner {
	return &Checkiner{
		WebInfo: webInfo,

		LoginHeaderAccpet:      loginHeaderAccpet,
		LoginHeaderContentType: loginHeaderContentType,
		HeaderMethod:           loginHeaderMethod,

		FlagCheckined: false,
		LastDay:       -1,
	}
}

func (c *Checkiner) setRequestHeader(req *http.Request) {
	header := map[string]string{
		"Accept":             c.LoginHeaderAccpet,
		"Content-Type":       c.LoginHeaderContentType,
		"Referer":            c.LoginURL,
		"Sec-Ch-Ua":          Headers["Sec-Ch-Ua"],
		"Sec-Ch-Ua-Mobile":   Headers["Sec-Ch-Ua-Mobile"],
		"Sec-Ch-Ua-Platform": Headers["Sec-Ch-Ua-Platform"],
		"User-Agent":         Headers["User-Agent"],
		"X-Requested-With":   Headers["X-Requested-With"],
	}
	// Add header
	for key, value := range header {
		req.Header.Set(key, value)
	}
}

func (c *Checkiner) setRequestBody(req *http.Request) {
	data := []byte("email=" + c.email + "&passwd=" + c.passwd)
	req.Body = io.NopCloser(bytes.NewBuffer(data))
}

func (c *Checkiner) handleLoginResponse(resp *http.Response, cookie *string) error {
	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Errorln("Status Code Error: ", resp.StatusCode)
		return errors.New("Status Code: " + string(rune(resp.StatusCode)))
	}

	// Handle response body
	buffer, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Read body failed: %v", err)
		return err
	}

	// fmt.Println(string(buffer))
	buf := map[string]any{}
	err = json.Unmarshal(buffer, &buf)
	if err != nil {
		log.Errorf("JSON parse failed: %v", err)
		return err
	}

	for k, v := range buf {
		if k == "ret" {
			log.Infoln(k, ":", v.(float64))
		} else if k == "msg" {
			util.SendNotify("Checkiner", "normal", ">>> "+c.Whoami+" "+v.(string))
		} else {
			util.SendNotify("Checkiner", "critical", "Unknown key: "+k)
		}
	}

	// TAG Get the lastest cookie
	for k, v := range resp.Header {
		if k == "Set-Cookie" {
			for _, val := range v {
				str := strings.Split(val, ";")
				*cookie += (str[0] + "; ")
			}
		}
	}
	return nil
}

// Display resoponse for JSON
func (c *Checkiner) handleResponse(reader io.Reader) error {
	body, err := io.ReadAll(reader)
	if err != nil {
		log.Errorf("Read body failed: %v", err)
		return err
	}

	var dat map[string]interface{}
	if err := json.Unmarshal(body, &dat); err != nil {
		log.Errorf("JSON parse failed: %v", err)
		return err
	}

	// Debug: response body
	for k, v := range dat {
		log.Infoln(k, ": ", v)
	}

	// TAG Level uses `critical` is to ensure checkin successfully for human.
	util.SendNotify("Checkiner", "critical",
		">>> "+c.Whoami+" checkin success: "+dat["msg"].(string))
	return nil
}

func (c *Checkiner) login() (string, error) {
	cookie := ""

	// Create request
	req, err := http.NewRequest(c.HeaderMethod, c.LoginURL, nil)
	if err != nil {
		log.Errorln(">>> "+c.Whoami+" Creating request failed: ", err)
		return cookie, err
	}
	c.setRequestHeader(req)
	c.setRequestBody(req)

	// Create HTTP client and send requset
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorln(">>> "+c.Whoami+" "+c.HeaderMethod+" request err: ", err)
		return cookie, err
	}
	defer resp.Body.Close()

	err = c.handleLoginResponse(resp, &cookie)
	if err != nil {
		return cookie, err
	}

	return cookie, nil
}

func (c *Checkiner) Checkin(headerAccpet, headerContentLength string) error {
	cookie, err := c.login()
	if err != nil {
		log.Errorln(">>> "+c.Whoami+" Login error: ", err)
		return err
	}

	header := map[string]string{
		"Accept":             headerAccpet,
		"Accept-Encoding":    Headers["Accept-Encoding"],
		"Accept-Language":    Headers["Accept-Language"],
		"Content-Length":     headerContentLength,
		"Cookie":             cookie,
		"Origin":             c.OriginURL,
		"Referer":            c.CheckinURL,
		"Sec-Ch-Ua":          Headers["Sec-Ch-Ua"],
		"Sec-Ch-Ua-Mobile":   Headers["Sec-Ch-Ua-Mobile"],
		"Sec-Ch-Ua-Platform": Headers["Sec-Ch-Ua-Platform"],
		"Sec-Fetch-Dest":     Headers["Sec-Fetch-Dest"],
		"Sec-Fetch-Mode":     Headers["Sec-Fetch-Mode"],
		"Sec-Fetch-Site":     Headers["Sec-Fetch-Site"],
		"User-Agent":         Headers["User-Agent"],
		"X-Requested-With":   Headers["X-Requested-With"],
	}

	// Create HTTP request
	req, err := http.NewRequest(c.HeaderMethod, c.CheckinURL, nil)
	if err != nil {
		log.Errorln(">>> "+c.Whoami+" Creating request failed: ", err)
		return err
	}

	// Add header
	for key, value := range header {
		req.Header.Set(key, value)
	}

	// Create HTTP client and send requset
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorln(">>> "+c.Whoami+" POST request err: ", err)
		return err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Errorln(">>> "+c.Whoami+" Status Code Error: ", resp.StatusCode)
		return err
	}

	// Debug: response header
	for k, v := range resp.Header {
		log.Infoln(k, ":", v[0])
	}

	// br 压缩
	// Cookie Expired
	if resp.Header.Get("Content-Type") == "text/html; charset=UTF-8" {
		return errors.New("cookie Expired")
	}

	if resp.Header.Get("Content-Encoding") == "br" {
		reader := brotli.NewReader(resp.Body)
		err := c.handleResponse(reader)
		if err != nil {
			log.Errorln(">>> "+c.Whoami+" Handle response failed: ", err)
			return err
		}
	} else if resp.Header.Get("Content-Encoding") == "gzip" {
		log.Infoln("gzip")
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			log.Errorln(">>> "+c.Whoami+" Create gzip reader error: ", err)
			return err
		}
		err = c.handleResponse(reader)
		if err != nil {
			log.Errorln(">>> "+c.Whoami+" Handle response failed: ", err)
			return err
		}
	} else if resp.Header.Get("Content-Encoding") == "deflate" {
		reader := flate.NewReader(resp.Body)
		err := c.handleResponse(reader)
		if err != nil {
			log.Errorln(">>> "+c.Whoami+" Handle response failed: ", err)
			return err
		}
	} else {
		log.Infoln("Not supported Content-Encoding")
		err := c.handleResponse(resp.Body)
		if err != nil {
			log.Errorln(">>> "+c.Whoami+" Handle response failed: ", err)
			return err
		}
	}
	return nil
}

func CheckinRun() (string, error) {
	// Create channel
	ch := make(chan struct{})

	// Timer
	go func(ch chan<- struct{}) {
		// Create timer
		timer := time.NewTicker(viper.GetDuration("common.interval"))
		defer func() {
			timer.Stop()
			close(ch)
		}()
		for {
			if _, ok := <-timer.C; !ok {
				log.Error("Timer error")
				return
			}
			ch <- struct{}{}
		}
	}(ch)

	// Checkin
	checkers := checker()
	for {
		if _, ok := <-ch; ok {
			log.Info("It's time to checkin")

			wg := sync.WaitGroup{}
			wg.Add(len(checkers))

			currDay := time.Time.Day(time.Now())

			for _, checker := range checkers {
				log.Infof("curr day: %v %v", currDay, checker)
				go func(checker *Checkiner) {
					defer wg.Done()

					if checker.LastDay != currDay {
						if checker.FlagCheckined {
							return
						}

						log.Infof("%s last_day: %d, curr_day: %d\n",
							checker.Whoami, checker.LastDay, currDay)

						if err := checker.Checkin(
							CheckinHeaderAccept, HeaderContentLength); err != nil {
							util.SendNotify("Checkiner", "critical",
								checker.Whoami+" Check in Failed: "+err.Error())
							return
						}
						checker.FlagCheckined = true
						checker.LastDay = currDay
					} else {
						log.Infof("Checkined tody: %v", currDay)
						checker.FlagCheckined = false
					}
				}(checker)
			}

			wg.Wait()
		}
	}
}

func webNames() []string {
	keys := viper.GetViper().AllKeys()
	names := make([]string, 0, len(keys)-2)

	for _, key := range keys {
		val := strings.Split(key, ".")[0]
		if val == "log" || val == "common" {
			continue
		}

		names = append(names, val)
	}

	return lo.Uniq(names)
}

func checker() []*Checkiner {
	checkers := make([]*Checkiner, 0)

	webNames := webNames()
	for _, webName := range webNames {
		infos := viper.GetStringMapString(webName)
		accountCnt := (len(infos) - 3) / 2
		for i := 0; i < accountCnt; i++ {
			checkers = append(checkers, NewCheckiner(
				LoginHeaderAccept,
				LoginHeaderContentType,
				HeaderMethod,
				NewWebInfo(
					webName,
					infos["origin"],
					infos["login"],
					infos["checkin"],
					infos[fmt.Sprintf("email_%d", i)],
					infos[fmt.Sprintf("passwd_%d", i)],
				),
			))
		}
	}

	return checkers
}
