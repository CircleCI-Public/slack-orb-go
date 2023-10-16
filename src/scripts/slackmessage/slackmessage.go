package slackmessage

import (
	"fmt"
	"os"
	"github.com/CircleCI-Public/slack-orb-go/src/scripts/httputils"
	"github.com/CircleCI-Public/slack-orb-go/src/scripts/jsonutils"
)

type SlackMessage struct {
	Template      string
	AccessToken   string
	IgnoreErrors  bool
	Channels      []string
}

func NewSlackMessage (template string, accessToken string, ignoreErrors bool, channels []string) *SlackMessage{
	return &SlackMessage{
		Template: template,
		AccessToken: accessToken,
		IgnoreErrors: ignoreErrors,
		Channels: channels,
	}
}

func (s *SlackMessage) PostMessage(channel string) {
	jsonWithChannel, _ := jsonutils.ApplyFunctionToJSON(s.Template, jsonutils.AddRootProperty("channel", channel))
	fmt.Printf("Posting the following JSON to Slack:\n%s\n", jsonWithChannel)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + s.AccessToken,
	}
	response, _ := httputils.SendHTTPRequest("POST", "https://slack.com/api/chat.postMessage", jsonWithChannel, headers)
	fmt.Printf("Slack API response:\n%s\n", response)

	errorMsg, _ := jsonutils.ApplyFunctionToJSON(response, jsonutils.ExtractRootProperty("error"))
	if errorMsg != "" {
		fmt.Printf("Slack API returned an error message:\n%s", errorMsg)
		fmt.Println("\n\nView the Setup Guide: https://github.com/CircleCI-Public/slack-orb/wiki/Setup")
		if !s.IgnoreErrors {
			os.Exit(1)
		}
	}
}
