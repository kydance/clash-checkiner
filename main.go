package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"Checkiner/pkg/checkin"
	"Checkiner/pkg/util"
)

var (
	//>>>> flags
	h bool
	// THY
	web string
	// /home/username/...
	path string
	// web : path
	webs map[string]string

	// time interval
	INTEVAL  time.Duration = time.Minute
	interval float64

	// The log file path
	LOG_FILE string = "./checkiner.log"
	// flags <<<<
)

// setWebMap sets website and cookie path
// map: webMap[website]cookie_path
func setWebMap(web string, path string) map[string]string {
	webMap := make(map[string]string)
	webs := strings.Split(web, util.DELEMITER)
	paths := strings.Split(path, util.DELEMITER)

	for idx, w := range webs {
		webMap[w] = paths[idx]
	}

	return webMap
}

func checkinRun(webs map[string]string) (string, error) {
	checkers := make([]*checkin.Checkin, 0)

	for webName, webCfg := range webs {
		// NewCheckiner()
		fmt.Println(webName, webCfg)
		if webName[0] == 'T' {
			THY_checker := checkin.NewCheckiner(webName,
				util.LOGIN_HEADER_ACCEPT, util.LOGIN_HEADER_CONTENT_TYPE,
				util.LOGIN_HEADER_METHOD, util.THY_URL_LOGIN,
				util.CHECKIN_HEADER_METHOD, util.THY_URL_CHECKIN, webs[webName])
			checkers = append(checkers, THY_checker)
		} else {
			CUTECLOUD_checker := checkin.NewCheckiner(webName,
				util.CUTECLOUD_LOGIN_HEADER_ACCEPT, util.LOGIN_HEADER_CONTENT_TYPE,
				util.LOGIN_HEADER_METHOD, util.CUTECLOUD_URL_LOGIN,
				util.CHECKIN_HEADER_METHOD, util.CUTECLOUD_URL_CHECKIN, webs[webName])

			checkers = append(checkers, CUTECLOUD_checker)
		}
	}

	// Create channel
	ch := make(chan struct{})

	// Timer
	go func(ch chan<- struct{}) {
		// Create timer
		timer := time.NewTicker(INTEVAL)
		defer func() {
			timer.Stop()
			close(ch)
		}()
		for {
			if _, ok := <-timer.C; !ok {
				log.Println("Timer error")
				return
			}
			ch <- struct{}{}
		}
	}(ch)

	// Checkin
	for {
		if _, ok := <-ch; ok {
			// fmt.Println("It's time to checkin")
			wg := sync.WaitGroup{}
			wg.Add(len(checkers))

			curr_day := time.Time.Day(time.Now())

			for _, checker := range checkers {
				fmt.Printf("curr day: %v %v", curr_day, checker)
				go func(checker *checkin.Checkin) {
					defer wg.Done()

					if checker.LastDay != curr_day {
						if checker.Flag_checkined {
							return
						}

						// thy
						log.Printf(
							"%s last_day: %d, curr_day: %d\n",
							checker.Whoami,
							checker.LastDay,
							curr_day,
						)
						if _, ok := webs[checker.Whoami]; ok {
							var err error = nil
							if checker.Whoami[0] == 'T' {
								err = checker.Checkin(util.CHECKIN_HEADER_ACCEPT,
									util.HEADER_CONTENT_LENGTH, util.THY_URL_ORIGIN)
							} else {
								err = checker.Checkin(util.CHECKIN_HEADER_ACCEPT,
									util.HEADER_CONTENT_LENGTH, util.CUTECLOUD_URL_ORIGIN)
							}
							if err != nil {
								util.NotifySend(
									"Checkiner",
									"critical",
									checker.Whoami+" Check in Failed: "+err.Error(),
								)
								return
							}
							checker.Flag_checkined = true
							checker.LastDay = curr_day
						} else {
							log.Printf("%s does not exist\n", checker.Whoami)
						}
					} else {
						// fmt.Printf("Checkined tody: %v", curr_day)
						log.Printf("Checkined tody: %v", curr_day)
						checker.Flag_checkined = false
					}
				}(checker)
			}

			wg.Wait()
		}
	}
}

func init() {
	flag.BoolVar(&h, "h", false, "help")

	flag.StringVar(&web, "w", `THY@THY1@CUTECLOUD`, "set target webs ("+
		util.DELEMITER+" is split char) support: [THY, CUTECLOUD]")
	flag.StringVar(&path, "p",
		`/home/tianen/go/src/Checkiner/config/THY
@/home/tianen/go/src/Checkiner/config/THY0@
/home/tianen/go/src/Checkiner/config/CUTECLOUD`,
		"set target webs cookie ("+util.DELEMITER+" is split char) support: [THY, CUTECLOUD]")
	flag.Float64Var(&interval, "i", 10, "set checkin interval (minute)")
	flag.StringVar(&LOG_FILE, "l", "./checkiner.log", "set log file path")

	flag.Usage = usage
}

func main() {
	flag.Parse()
	webs = setWebMap(web, path)
	INTEVAL = time.Duration(float64(time.Minute) * interval)

	if h || web == "" || path == "" || interval <= 0 {
		flag.Usage()
		return
	}

	// Logger
	log_file, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Open log file error: ", err)
		return
	}
	defer log_file.Close()

	// Set log information to log file
	log.SetOutput(log_file)

	// Welcome
	util.NotifySend("Checkiner", "normal", "Welcome to enjoy your time with Checkiner")

	// It's time to checkin
	for {
		who, err := checkinRun(webs)
		// Checkiner failed
		if err != nil {
			util.NotifySend("Checkiner", "critical", who+" Check in Failed: "+err.Error())
		}
	}
}

func usage() {
	log.Printf(`Checkiner version: checkiner/1.3.0
Usage: checkiner [-h] [-w web]

Example: checkiner -i 120 -w THY@CUTECLOUD 
-p /home/tianen/go/src/Checkiner/config/THY@/home/tianen/go/src/Checkiner/config/CUTECLOUD

Options:
`)
	flag.PrintDefaults()
}
