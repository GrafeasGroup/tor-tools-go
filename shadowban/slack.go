package shadowban

import (
	"fmt"
	"os"
	"strings"

	"github.com/slack-go/slack"
)

func getBannedMessage(username string) string {
	return fmt.Sprintf(strings.TrimSpace(`
:biohazard::derp::biohazard::derp::biohazard: SHADOWBAN ALERT! :biohazard::derp::biohazard::derp::biohazard:

:banhammer::banhammer: User <https://www.reddit.com/u/%s|%s> appears to be *SHADOWBANNED*! :banhammer::banhammer:
`), username, username)
}

func getSlackToken() string {
	token := os.Getenv("SLACK_TOKEN")

	if token == "" {
		panic("Missing environment variable: SLACK_TOKEN")
	}

	return token
}

func getChannelName() string {
	channel := os.Getenv("SLACK_CHANNEL")

	if channel == "" {
		panic("Missing environment variable: SLACK_CHANNEL")
	}

	return channel
}

func getChannelID(api *slack.Client, channel string) (string, error) {
	channels, _, err := api.GetConversations(&slack.GetConversationsParameters{Cursor: "", ExcludeArchived: "true", Limit: 900, Types: []string{"public_channel","private_channel"}})
	if err != nil {
		fmt.Println(err)
		panic("ERROR")
	}
	for _, focus := range channels {
		if focus.Name == channel {
			return focus.ID, nil
		}
	}

	return "", fmt.Errorf("Could not find #%s", channel)
}

func sendSlackAlert(bans <-chan string, notice chan<- string) {
	var username string
	var channelID string

	api := slack.New(getSlackToken())
	channelID, err := getChannelID(api, getChannelName())
	if err != nil {
		fmt.Println(err)
		panic("ERROR")
	}
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		username = <-bans
		rtm.SendMessage(rtm.NewOutgoingMessage(getBannedMessage(username), channelID))
		notice <- username
	}
}
