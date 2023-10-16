package slacknotification

import (
	"log"
	
	"github.com/CircleCI-Public/slack-orb-go/src/scripts/stringutils"
	"github.com/CircleCI-Public/slack-orb-go/src/scripts/jsonutils"
)

type SlackNotification struct {
	Status        string
	Branch        string
	Tag           string
	Event         string
	BranchPattern string
	TagPattern    string
	InvertMatch   bool
	Template      string
	InlineTemplate string
	EnvVarContainingTemplate string
}

func NewSlackNotification(status, branch, tag, event, branchPattern, tagPattern, inlineTemplate, envVarContainingTemplate string, invertMatch bool) *SlackNotification {
	return &SlackNotification{
		Status:        status,
		Branch:        branch,
		Tag:           tag,
		Event:         event,
		BranchPattern: branchPattern,
		TagPattern:    tagPattern,
		InvertMatch:   invertMatch,
		InlineTemplate: inlineTemplate,
		EnvVarContainingTemplate: envVarContainingTemplate,
	}
}

func (j *SlackNotification) IsEventMatchingStatus() bool {
	return stringutils.IsEventMatchingStatus(j.Event, j.Status)
}

func (j *SlackNotification) IsPostConditionMet() bool {
	branchMatches, _ := stringutils.IsPatternMatchingString(j.BranchPattern, j.Branch)
	tagMatches, _ := stringutils.IsPatternMatchingString(j.TagPattern, j.Tag)
	return stringutils.IsPostConditionMet(branchMatches, tagMatches, j.InvertMatch)
}

func (j *SlackNotification) BuildMessageBody() string {
	// Build the message body
	template, err := jsonutils.DetermineTemplate(j.InlineTemplate, j.Status, j.EnvVarContainingTemplate)
	if err != nil {
		log.Fatalf("%v", err)
	}
	if template == "" {
		log.Fatalf("the template %q is empty. Exiting without posting to Slack...", template)
	}

	// Expand environment variables in the template
	templateWithExpandedVars, err := jsonutils.ApplyFunctionToJSON(template, jsonutils.ExpandEnvVarsInInterface)
	if err != nil {
		log.Fatal(err)
	}

	// Add a "channel" property with a nested "myChannel" property
	modifiedJSON, err := jsonutils.ApplyFunctionToJSON(templateWithExpandedVars, jsonutils.AddRootProperty("channel", "my_channel"))
	if err != nil {
		log.Fatalf("%v", err)
	}

	return modifiedJSON
}

