package slacknotification

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSlackNotification(t *testing.T) {
	status := "success"
	branch := "main"
	tag := "v1.0"
	event := "push"
	branchPattern := "main"
	tagPattern := "v1.*"
	inlineTemplate := "template"
	envVarContainingTemplate := "TEMPLATE_ENV"
	invertMatch := false

	notification := NewSlackNotification(status, branch, tag, event, branchPattern, tagPattern, inlineTemplate, envVarContainingTemplate, invertMatch)

	assert.Equal(t, status, notification.Status)
	assert.Equal(t, branch, notification.Branch)
	assert.Equal(t, tag, notification.Tag)
	assert.Equal(t, event, notification.Event)
	assert.Equal(t, branchPattern, notification.BranchPattern)
	assert.Equal(t, tagPattern, notification.TagPattern)
	assert.Equal(t, inlineTemplate, notification.InlineTemplate)
	assert.Equal(t, envVarContainingTemplate, notification.EnvVarContainingTemplate)
	assert.Equal(t, invertMatch, notification.InvertMatch)
}
func TestIsEventMatchingStatus(t *testing.T) {
	tests := []struct {
		name   string
		event  string
		status string
		want   bool
	}{
		{
			name:   "matching event and status",
			event:  "push",
			status: "push",
			want:   true,
		},
		{
			name:   "non-matching event and status",
			event:  "push",
			status: "pull",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sn := NewSlackNotification(tt.status, "", "", tt.event, "", "","","", false)
			got := sn.IsEventMatchingStatus()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsPostConditionMet(t *testing.T) {
	tests := []struct {
		name          string
		branch        string
		tag           string
		branchPattern string
		tagPattern    string
		invertMatch   bool
		want          bool
	}{
		{
			name:          "matching branch and tag patterns",
			branch:        "main",
			tag:           "v1.0",
			branchPattern: "main",
			tagPattern:    "v1.*",
			invertMatch:   false,
			want:          true,
		},
		{
			name:          "non-matching branch and tag patterns",
			branch:        "dev",
			tag:           "v2.0",
			branchPattern: "main",
			tagPattern:    "v1.*",
			invertMatch:   false,
			want:          false,
		},
		{
			name:          "invert match",
			branch:        "dev",
			tag:           "v2.0",
			branchPattern: "main",
			tagPattern:    "v1.*",
			invertMatch:   true,
			want:          true,
		},
		{
			name:          "empty branch and tag",
			branch:        "",
			tag:           "",
			branchPattern: "main",
			tagPattern:    "v1.*",
			invertMatch:   false,
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sn := NewSlackNotification("", tt.branch, tt.tag, "", tt.branchPattern, tt.tagPattern,"","", tt.invertMatch)
			got := sn.IsPostConditionMet()
			assert.Equal(t, tt.want, got)
		})
	}
}
