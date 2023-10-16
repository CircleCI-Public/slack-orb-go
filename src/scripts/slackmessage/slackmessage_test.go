package slackmessage

import (
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestNewSlackMessage(t *testing.T) {
	template := "template"
	accessToken := "token"
	ignoreErrors := false
	channels := []string{"channel1", "channel2"}

	expected := &SlackMessage{
		Template:     template,
		AccessToken:  accessToken,
		IgnoreErrors: ignoreErrors,
		Channels:     channels,
	}

	actual := NewSlackMessage(template, accessToken, ignoreErrors, channels)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected and actual do not match. Expected: %v, Actual: %v", expected, actual)
	}
}

func TestPostMessage(t *testing.T) {
	// Create a new SlackMessage
	s := &SlackMessage{
		AccessToken:  "test_token",
		Template:     "{}",
		IgnoreErrors: true,
	}

	// Activate the httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock the Slack API response
	httpmock.RegisterResponder("POST", "https://slack.com/api/chat.postMessage",
		httpmock.NewStringResponder(200, `{"ok": true}`))

	// Capture the stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function
	s.PostMessage("test_channel")

	// Close the pipe and restore the original stdout
	w.Close()
	os.Stdout = oldStdout

	// Read the captured stdout
	out, _ := io.ReadAll(r)

	// Assert that the function printed the correct output
	assert.Contains(t, string(out), "Posting the following JSON to Slack:")
	assert.Contains(t, string(out), "Slack API response:")
}
