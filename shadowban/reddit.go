package shadowban

import (
	"fmt"
	// "io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

const toolAuthor string = "u/personal_opinions"

func getUserAgent() string {
	return fmt.Sprintf("golang:tor-tools-go/shadowbanned:v%s (by %s)", toolVersion, toolAuthor)
}

func getMinPollTime() int {
	value := os.Getenv("MIN_POLL_TIME")

	if value == "" {
		return 30  // Default min: 30 seconds
	}

	i, _ := strconv.Atoi(value)

	return i
}

func getMaxPollTime() int {
	value := os.Getenv("MAX_POLL_TIME")

	if value == "" {
		return 300 // Default max: 5 minutes
	}

	i, _ := strconv.Atoi(value)

	return i
}

func randomJitter() {
	minPoll := getMinPollTime()
	maxPoll := getMaxPollTime()
	randTop := maxPoll - minPoll

	if randTop < 0 {
		panic("given minimum poll time is greater than given maximum poll time")
	}

	// Default: sleep for anywhere from 30 seconds to 5 minutes
	time.Sleep(time.Duration(rand.Intn(randTop) + minPoll) * time.Second)
}

func watchUsername(username string, banned chan<- string) {
	client := &http.Client{}
	var isBanned bool
	var err error

	for {
		randomJitter()

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
	// 	log.Error(err)
	// 	panic("ERROR")
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
