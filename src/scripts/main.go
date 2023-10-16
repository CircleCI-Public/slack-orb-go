package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/CircleCI-Public/slack-orb-go/src/scripts/ioutils"
	"github.com/CircleCI-Public/slack-orb-go/src/scripts/slackmessage"
	"github.com/CircleCI-Public/slack-orb-go/src/scripts/slacknotification"

	"github.com/a8m/envsubst"
	"github.com/joho/godotenv"
)

func main() {
	// Load the environment variables from the configuration file
	// This has to be done before anything else to ensure that the environment variables modified by the configuration file are available
	bashEnv := os.Getenv("BASH_ENV")
	if ioutils.FileExists(bashEnv) {
		fmt.Println("Loading BASH_ENV into the environment...")
		if err := godotenv.Load(bashEnv); err != nil {
			log.Fatal("Error loading BASH_ENV file:", err)
		}
	}

	// Load the job status from the configuration file
	jobStatusFile := "/tmp/SLACK_JOB_STATUS"
	if ioutils.FileExists(jobStatusFile) {
		fmt.Println("Loading SLACK_JOB_STATUS into the environment...")
		if err := godotenv.Load(jobStatusFile); err != nil {
			log.Fatal("Error loading SLACK_JOB_STATUS file:", err)
		}
	}

	// Fetch environment variables
	accessToken := os.Getenv("SLACK_ACCESS_TOKEN")
	branchPattern := os.Getenv("SLACK_PARAM_BRANCHPATTERN")
	channelsStr := os.Getenv("SLACK_PARAM_CHANNEL")
	envVarContainingTemplate := os.Getenv("SLACK_PARAM_TEMPLATE")
	eventToSendMessage := os.Getenv("SLACK_PARAM_EVENT")
	inlineTemplate := os.Getenv("SLACK_PARAM_CUSTOM")
	invertMatchStr := os.Getenv("SLACK_PARAM_INVERT_MATCH")
	jobBranch := os.Getenv("CIRCLE_BRANCH")
	jobStatus := os.Getenv("CCI_STATUS")
	jobTag := os.Getenv("CIRCLE_TAG")
	tagPattern := os.Getenv("SLACK_PARAM_TAGPATTERN")
	ignoreErrorsStr := os.Getenv("SLACK_PARAM_IGNORE_ERRORS")

	// Expand environment variables
	accessToken, _ = envsubst.String(accessToken)
	branchPattern, _ = envsubst.String(branchPattern)
	channelsStr, _ = envsubst.String(channelsStr)
	envVarContainingTemplate, _ = envsubst.String(envVarContainingTemplate)
	eventToSendMessage, _ = envsubst.String(eventToSendMessage)
	invertMatchStr, _ = envsubst.String(invertMatchStr)
	ignoreErrorsStr, _ = envsubst.String(ignoreErrorsStr)
	tagPattern, _ = envsubst.String(tagPattern)

	invertMatch, _ := strconv.ParseBool(invertMatchStr)
	ignoreErrors, _ := strconv.ParseBool(ignoreErrorsStr)
	channels := strings.Split(channelsStr, ",")
	
	slackNotification := slacknotification.NewSlackNotification(jobStatus, jobBranch, jobTag, eventToSendMessage, branchPattern, tagPattern, inlineTemplate, envVarContainingTemplate, invertMatch)
	modifiedJSON := slackNotification.BuildMessageBody()
	slackMessage := slackmessage.NewSlackMessage(modifiedJSON, accessToken, ignoreErrors, channels)
	
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
