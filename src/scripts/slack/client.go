package slack

import (
	"errors"
	"context"
	"fmt"
	"time"

	"github.com/circleci/ex/config/secret"
	"github.com/circleci/ex/httpclient"
	"github.com/CircleCI-Public/slack-orb-go/src/scripts/jsonutils"

)

const defaultSlackURL = "https://slack.com/api"

type  Client struct {
	hc		 *httpclient.Client
}

type ClientOptions struct {
	BaseURL		string
	SlackToken	secret.String
}

type APIResponse struct {
	Error string `json:"error"`
}

func NewClient(options ClientOptions) *Client {
	baseURL := defaultSlackURL
	if options.BaseURL != "" {
		baseURL = options.BaseURL
	}
	hc := httpclient.New(httpclient.Config{
		Name:       "Slack Client",
		BaseURL:    baseURL,
		AuthToken:  options.SlackToken.Value(),
		AcceptType: httpclient.JSON,
		Timeout:    time.Second * 10,
	})

	return &Client{hc}
}

func (c *Client) PostMessage (ctx context.Context, message, channel string) error {
	jsonWithChannel, err := jsonutils.ApplyFunctionToJSON(message, jsonutils.AddRootProperty("channel", channel))
	if err != nil {
		return err
	}
	
	var response APIResponse

	req := httpclient.NewRequest("POST", "/chat.postMessage",
		httpclient.Body(jsonWithChannel),
		httpclient.JSONDecoder(&response),
	)

	err = c.hc.Call(context.Background(), req)

	fmt.Printf("This is the error: %v", err)
	if err != nil {
		return err
	}

	if response.Error != "" {
		return errors.New(response.Error)
	}
	return nil
}