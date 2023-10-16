package config

import (
	"os"
	"testing"
)

func TestExpandEnvVariables(t *testing.T) {
	tests := []struct {
		envVar      string
		envValue    string
		configVar   string
		fieldName   string
		expectedErr string
	}{
		{
			envVar:      "TEST_VARIABLE",
			envValue:    "${UNSET_VARIABLE}",
			configVar:   "${TEST_VARIABLE}",
			fieldName:   "AccessToken",
			expectedErr: "AccessToken",
		},
		{
			envVar:    "TEST_VARIABLE",
			envValue:  "value",
			configVar: "${TEST_VARIABLE}",
			fieldName: "AccessToken",
		},
	}

	for _, test := range tests {
		os.Setenv(test.envVar, test.envValue)

		config := &Config{AccessToken: test.configVar}

		err := config.ExpandEnvVariables()

		if err != nil && test.expectedErr != "" {
			expErr, ok := err.(*ExpansionError)
			if ok {
				if expErr.FieldName != test.expectedErr {
					t.Errorf("Expected error field name: %s, got: %s", test.expectedErr, expErr.FieldName)
				}
			} else {
				t.Errorf("Expected ExpansionError, got: %v", err)
			}
		}

		os.Unsetenv(test.envVar)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		accessToken string
		channelStr  string
		expectedErr string
	}{
		{
			accessToken: "",
			channelStr:  "channel",
			expectedErr: "SLACK_ACCESS_TOKEN",
		},
		{
			accessToken: "token",
			channelStr:  "",
			expectedErr: "SLACK_PARAM_CHANNEL",
		},
	}

	for _, test := range tests {
		config := &Config{AccessToken: test.accessToken, ChannelsStr: test.channelStr}

		err := config.Validate()

		if err != nil {
			envErr, ok := err.(*EnvVarError)
			if ok {
				if envErr.VarName != test.expectedErr {
					t.Errorf("Expected error var name: %s, got: %s", test.expectedErr, envErr.VarName)
				}
			} else {
				t.Errorf("Expected EnvVarError, got: %v", err)
			}
		}
	}
}
