package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/CircleCI-Public/slack-orb-go/src/scripts/config"
	"github.com/CircleCI-Public/slack-orb-go/src/scripts/slack"
)

func main() {
	// Load environment variables from BASH_ENV and SLACK_JOB_STATUS files
	// This has to be done before loading the configuration because the configuration
	// depends on the environment variables loaded from these files
	if err := config.LoadEnvFromFile(os.Getenv("BASH_ENV")); err != nil {
		log.Fatal(err)
	}
	if err := config.LoadEnvFromFile("/tmp/SLACK_JOB_STATUS"); err != nil {
		log.Fatal(err)
	}

	conf := config.NewConfig()

	if err := conf.ExpandEnvVariables(); err != nil {
		log.Fatalf("Error expanding environment variables: %v", err)
	}

	if err := conf.Validate(); err != nil {
		if envVarError, ok := err.(*config.EnvVarError); ok {
			switch envVarError.VarName {
			case "SLACK_ACCESS_TOKEN":
				log.Fatalf(
					"In order to use the Slack Orb an OAuth token must be present via the SLACK_ACCESS_TOKEN environment variable." +
						"\nFollow the setup guide available in the wiki: https://github.com/CircleCI-Public/slack-orb/wiki/Setup.",
				)
			case "SLACK_PARAM_CHANNEL":
				log.Fatalf(
					`No channel was provided. Please provide one or more channels using the "SLACK_PARAM_CHANNEL" environment variable or the "channel" parameter.`,
				)
			default:
				log.Fatalf("Configuration validation failed: Environment variable not set: %s", envVarError.VarName)
			}
		} else {
			log.Fatalf("Configuration validation failed: %v", err)
		}
	}

	invertMatch, _ := strconv.ParseBool(conf.InvertMatchStr)
	ignoreErrors, _ := strconv.ParseBool(conf.IgnoreErrorsStr)
	channels := strings.Split(conf.ChannelsStr, ",")

	slackNotification := slack.Notification{
		Status:                   conf.JobStatus,
		Branch:                   conf.JobBranch,
		Tag:                      conf.JobTag,
		Event:                    conf.EventToSendMessage,
		BranchPattern:            conf.BranchPattern,
		TagPattern:               conf.TagPattern,
		InvertMatch:              invertMatch,
		InlineTemplate:           conf.InlineTemplate,
		EnvVarContainingTemplate: conf.EnvVarContainingTemplate,
	}

	modifiedJSON := slackNotification.BuildMessageBody()
	
	slackMessage := slack.Message{
		Template: modifiedJSON,
		AccessToken: conf.AccessToken,
		IgnoreErrors: ignoreErrors,
		Channels: channels,
	}

	if !slackNotification.IsEventMatchingStatus() {
		message := fmt.Sprintf("The job status %q does not match the status set to send alerts %q.", slackNotification.Status, slackNotification.Event)
		fmt.Println(message)
		fmt.Println("Exiting without posting to Slack...")
		os.Exit(0)
	}

	if !slackNotification.IsPostConditionMet() {
		fmt.Println("The post condition is not met. Neither the branch nor the tag matches the pattern or the match is inverted.")
		fmt.Println("Exiting without posting to Slack...")
		os.Exit(0)
	}

	for _, channel := range slackMessage.Channels {
		slackMessage.PostMessage(channel)
	}
}
