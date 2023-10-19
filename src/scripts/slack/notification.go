package slack

import (
	"log"

	"github.com/CircleCI-Public/slack-orb-go/src/scripts/jsonutils"
	"github.com/CircleCI-Public/slack-orb-go/src/scripts/stringutils"
)

type Notification struct {
	Status                   string
	Branch                   string
	Tag                      string
	Event                    string
	BranchPattern            string
	TagPattern               string
	InvertMatch              bool
	Template                 string
	InlineTemplate           string
	EnvVarContainingTemplate string
}

func (j *Notification) IsEventMatchingStatus() bool {
	return j.Status == j.Event || j.Event == "always"
}

func (j *Notification) IsPostConditionMet() bool {
	branchMatches, _ := stringutils.IsPatternMatchingString(j.BranchPattern, j.Branch)
	tagMatches, _ := stringutils.IsPatternMatchingString(j.TagPattern, j.Tag)
	return (branchMatches || tagMatches) != j.InvertMatch

}

func (j *Notification) BuildMessageBody() string {
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
