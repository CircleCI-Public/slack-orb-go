package slack

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostMessage(t *testing.T) {
	// Create a test server that always responds with a successful SlackResponse
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok": true}`))
	}))
	defer ts.Close()

	msg := &Message{
		Template:    `{"text": "Hello, world!"}`,
		AccessToken: "",
		IgnoreErrors: false,
		Channels:    []string{"test-channel"},
	}


	statusCode, err := msg.PostMessage("test-channel")

	fmt.Printf("%d", statusCode)

	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
}



