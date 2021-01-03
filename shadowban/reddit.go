package shadowban

import (
	"fmt"
	// golog "log"
	// "io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const toolVersion string = "0.2.0"
const toolAuthor string = "u/personal_opinions"

func getUserAgent() string {
	return fmt.Sprintf("golang:tor-tools-go/shadowbanned:v%s (by %s)", toolVersion, toolAuthor)
}

func watchUsername(username string, banned chan<- string) {
	client := &http.Client{}
	var isBanned bool
	var err error

	for {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)

		isBanned, err = isShadowBanned(client, username)
		if err != nil {
			if err, ok := err.(*RedditError); ok {
				continue
			} else {
				fmt.Println(err)
				panic("ERROR")
			}
		}

		if isBanned {
			banned <- username
			break
		}
	}
}

func isShadowBanned(client *http.Client, username string) (bool, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://www.reddit.com/user/%s/about.json", username), nil)
	req.Header.Set("User-Agent", getUserAgent())
	req.Header.Set("Accept", "application/json")
	res, _ := client.Do(req)

	// body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	// if err != nil {
	// 	golog.Fatal(err)
	// }

	log.WithFields(logrus.Fields{
		"username":   username,
		"httpStatus": fmt.Sprintf("%d", res.StatusCode),
	}).Info("Response from Reddit")

	if res.StatusCode >= 500 {
		// Do nothing. Sometimes Reddit barfs. It's a feature. Don't trigger the red alert in this event.
		return false, &RedditError{message: "Reddit is down", StatusCode: res.StatusCode}
	}
	return false, nil
}
