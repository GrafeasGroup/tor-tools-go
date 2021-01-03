package shadowban

import (
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func initLogs() {
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
}

// Checker loop
func Checker(done chan bool) {
	initLogs()

	botUsernames := strings.Fields(os.Getenv("BOT_USERNAMES"))
	bannedUsernames := []string{}
	bans := make(chan string)
	notice := make(chan string)

	go sendSlackAlert(bans, notice)
	for _, username := range botUsernames {
		go watchUsername(username, bans)
	}
	log.Info("Initialized all go routines")

	var username string

	// Loop until everyone is found to be banned
	for {
		username = <-notice
		bannedUsernames = append(bannedUsernames, username)
		if len(bannedUsernames) == len(botUsernames) {
			break
		}
	}

	time.Sleep(3 * time.Second)

	done <- true
}
