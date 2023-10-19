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

func TestPostMessage_ErrorCases(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name           string
		mockResponse   string
		expectedOutput string
	}{
		{"invalid_auth", `{"ok": false, "error": "invalid_auth"}`, "invalid_auth"},
		{"account_inactive", `{"ok": false, "error": "account_inactive"}`, "account_inactive"},
		{"channel_not_found", `{"ok": false, "error": "channel_not_found"}`, "channel_not_found"},
		{"is_archived", `{"ok": false, "error": "is_archived"}`, "is_archived"},
		{"msg_too_long", `{"ok": false, "error": "msg_too_long"}`, "msg_too_long"},
		{"no_text", `{"ok": false, "error": "no_text"}`, "no_text"},
		{"rate_limited", `{"ok": false, "error": "rate_limited"}`, "rate_limited"},
		{"not_authed", `{"ok": false, "error": "not_authed"}`, "not_authed"},
		{"not_in_channel", `{"ok": false, "error": "not_in_channel"}`, "not_in_channel"},
		{"user_is_bot", `{"ok": false, "error": "user_is_bot"}`, "channel_not_found"},
	}

	// Create a new Message
	s := &Message{
		AccessToken:  "test_token",
		Template:     "{}",
		IgnoreErrors: true,
	}

	// Activate the httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Iterate over the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock the Slack API response
			httpmock.RegisterResponder("POST", "https://slack.com/api/chat.postMessage",
				httpmock.NewStringResponder(200, tc.mockResponse))

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
			assert.Check(t, cmp.Contains(string(out), tc.expectedOutput))
		})
	}
}