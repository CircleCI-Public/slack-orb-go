package slack

import (
	"fmt"
	"os"

	"github.com/CircleCI-Public/slack-orb-go/src/scripts/httputils"
	"github.com/CircleCI-Public/slack-orb-go/src/scripts/jsonutils"
)

type Message struct {
	Template     string
	AccessToken  string
	IgnoreErrors bool
	Channels     []string
}

func (s *Message) PostMessage(channel string) {
	jsonWithChannel, _ := jsonutils.ApplyFunctionToJSON(s.Template, jsonutils.AddRootProperty("channel", channel))
	fmt.Printf("Posting the following JSON to Slack:\n%s\n", jsonWithChannel)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + s.AccessToken,
	}
	response, statusCode, _ := httputils.SendHTTPRequest("POST", "https://slack.com/api/chat.postMessage", jsonWithChannel, headers)
	fmt.Printf("Slack API response:\nStatus Code: %d\n%s\n", statusCode, response)

	errorMsg, _ := jsonutils.ApplyFunctionToJSON(response, jsonutils.ExtractRootProperty("error"))
	if errorMsg != "" {
		fmt.Printf("Slack API returned an error message:\nStatus Code: %d\\n%s", statusCode, errorMsg)
		fmt.Println("\n\nView the Setup Guide: https://github.com/CircleCI-Public/slack-orb/wiki/Setup")
		if !s.IgnoreErrors {
			os.Exit(1)
		}
	}
}
