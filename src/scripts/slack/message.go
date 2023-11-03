package slack

import (
	"context"
	"fmt"
	"os"
	"time"
	"github.com/circleci/ex/httpclient"
	"github.com/CircleCI-Public/slack-orb-go/src/scripts/jsonutils"
)

type Message struct {
	Template     string
	AccessToken  string
	IgnoreErrors bool
	Channels     []string
}

type SlackError struct {
	ErrorMessage	string
	Err				error		
}

func (e *SlackError) Error () string {
	return fmt.Sprintf("Error posting to Slack: %v\nSlack API response: %s", e.Err, e.ErrorMessage)
}
type SlackResponse struct {
	Ok      bool   `json:"ok"`
	Error string `json:"error"`
}

func (s *Message) PostMessage(channel string) (int, error) {
	jsonWithChannel, err := jsonutils.ApplyFunctionToJSON(s.Template, jsonutils.AddRootProperty("channel", channel))
	if err != nil {
		return 0, err
	}
	fmt.Printf("Posting the following JSON to Slack:\n%s\n", jsonWithChannel)
	
	var response SlackResponse
	
	client := httpclient.New(httpclient.Config{
			Name:		"Slack Client",
			BaseURL: 	"https://slack.com",
			AuthToken:	s.AccessToken,
			AcceptType:	"application/json",
			Timeout:     time.Second * 10,
	})

	req := httpclient.NewRequest("POST", "/api/chat.postMessage",
		httpclient.Body(jsonWithChannel),
		httpclient.Header("Content-Type", "application/json"),
		httpclient.JSONDecoder(&response),
	)

	err = client.Call(context.Background(),req)
	fmt.Printf("This is the error: %v", err)
	if err != nil {
		if !s.IgnoreErrors{
			os.Exit(1)
		}

		httpErr, ok := err.(*httpclient.HTTPError)
		if ok && err != nil  {
			slackErr := &SlackError{ErrorMessage: httpErr.Error(), Err: err}
			return httpErr.Code(), slackErr
		}
		slackErr := &SlackError{ErrorMessage: response.Error, Err: err}
		return 0, slackErr
	}
	return 200, nil
}
