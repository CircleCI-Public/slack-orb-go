package slack

import (
	"io"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

// Check Slack API docs to come up with a list of responses we support and then add test here
// ie Channel does not exist etc. 
// validation for slack notification package
func TestPostMessage(t *testing.T) {
	// Create a new Message
	s := &Message{
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
	assert.Check(t, cmp.Contains(string(out), "Posting the following JSON to Slack:"))
	assert.Check(t, cmp.Contains(string(out), "Slack API response:"))
}
