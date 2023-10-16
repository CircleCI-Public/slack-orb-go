package config

import (
	"os"
	"testing"
)

func TestExpandEnvVariables(t *testing.T) {
	tests := []struct {
		configVar   string
		envVars     map[string]string
		expectedErr string
		expectedVal string
		fieldName   string
		id          string
	}{
		{
			// This test case checks the basic functionality of variable expansion.
			// An environment variable "TEST_VAR" is set, and it is expected that
			// the variable in the configuration gets expanded to the value of "TEST_VAR".
			configVar:   "${TEST_VAR}",
			envVars:     map[string]string{"TEST_VAR": "value"},
			expectedErr: "",
			expectedVal: "value",
			fieldName:   "AccessToken",
			id:          "1",
		},
		{
			// This test case checks whether a suffix can be successfully appended
			// to the value of an environment variable after it gets expanded.
			configVar:   "${TEST_VAR}_suffix",
			envVars:     map[string]string{"TEST_VAR": "value"},
			expectedErr: "",
			expectedVal: "value_suffix",
			fieldName:   "AccessToken",
			id:          "2",
		},
		{
			// This test case checks whether two environment variables can be
			// concatenated successfully after they get expanded.
			configVar:   "${TEST_VAR}_${ANOTHER_VAR}",
			envVars:     map[string]string{"TEST_VAR": "value", "ANOTHER_VAR": "another_value"},
			expectedErr: "",
			expectedVal: "value_another_value",
			fieldName:   "AccessToken",
			id:          "3",
		},
		{
			// This test case checks whether the default value is used when the environment
			// variable is set but empty. The default value "default_value" should be used.
			configVar:   "${TEST_VAR:-default_value}",
			envVars:     map[string]string{"TEST_VAR": ""},
			expectedErr: "",
			expectedVal: "default_value",
			fieldName:   "AccessToken",
			id:          "4",
		},
		{
			// This test case checks the behavior when an environment variable is unset.
			// The configuration should have an empty value as "UNSET_VAR" is not set.
			configVar:   "${UNSET_VAR}",
			envVars:     map[string]string{},
			expectedErr: "",
			expectedVal: "",
			fieldName:   "AccessToken",
			id:          "5",
		},
	}

	for _, test := range tests {
		t.Run(test.id, func(t *testing.T) {
			// Setting environment variables
			for varName, val := range test.envVars {
				os.Setenv(varName, val)
				defer os.Unsetenv(varName)

			}

			config := &Config{AccessToken: test.configVar}
			err := config.ExpandEnvVariables()

			if err != nil {
				expErr, ok := err.(*ExpansionError)
				if ok {
					if expErr.FieldName != test.expectedErr {
						t.Errorf("Expected error field name: %q, got: %s", test.expectedErr, expErr.FieldName)
					}
				} else {
					t.Errorf("Expected ExpansionError, got: %v", err)
				}
			} else if test.expectedErr != "" {
				t.Errorf("Expected error for field name: %q, but got nil", test.expectedErr)
			} else {
				actualVal := config.AccessToken
				if actualVal != test.expectedVal {
					t.Errorf("Expected value %q, but got %s", test.expectedVal, actualVal)
				}
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		config      *Config
		expectedErr string // This holds the name of the field expected to error
		id          string
	}{
		{
			// This test case checks the behavior when the access token is missing.
			config:      &Config{AccessToken: "", ChannelsStr: "channel"},
			expectedErr: "SLACK_ACCESS_TOKEN",
			id:          "1",
		},
		{
			// This test case checks the behavior when the channel string is missing.
			config:      &Config{AccessToken: "token", ChannelsStr: ""},
			expectedErr: "SLACK_PARAM_CHANNEL",
			id:          "2",
		},
		{
			// This test case checks the behavior when nothing is missing.
			config:      &Config{AccessToken: "token", ChannelsStr: "channel"},
			expectedErr: "",
			id:          "3",
		},
	}

	for _, test := range tests {
		t.Run(test.id, func(t *testing.T) {
			err := test.config.Validate()

			if err != nil {
				envErr, ok := err.(*EnvVarError)
				if ok {
					if envErr.VarName != test.expectedErr {
						t.Errorf("Expected error var name: %s, got: %s", test.expectedErr, envErr.VarName)
					}
				} else {
					t.Errorf("Expected EnvVarError, got: %v", err)
				}
			} else if test.expectedErr != "" {
				t.Errorf("Expected error for field name: %s, but got nil", test.expectedErr)
			}
		})
	}
}

func TestLoadEnvFromFile(t *testing.T) {
	tests := []struct {
		envVarName  string
		envVarValue string
		expectedErr bool
		filePath    string
		id          string
	}{
		{
			// This test case checks the behavior when the file does not exist.
			filePath:    "/path/that/does/not/exist",
			envVarName:  "",
			envVarValue: "",
			expectedErr: false,
			id:          "1",
		},
		{
			// This test case checks the successful loading of environment variables from a file.
			envVarName:  "TEST_VAR",
			envVarValue: "potato",
			expectedErr: false,
			filePath:    "testdata/valid_env_file",
			id:          "2",
		},
		{
			// This test case checks the behavior when the file is invalid.
			envVarName:  "TEST_VAR",
			envVarValue: "potato",
			expectedErr: true,
			filePath:    "testdata/invalid_env_file",
			id:          "3",
		},
	}

	for _, test := range tests {
		t.Run(test.id, func(t *testing.T) {

			err := LoadEnvFromFile(test.filePath)

			if (err != nil) != test.expectedErr {
				t.Errorf("Expected error: %v, got: %v", test.expectedErr, err)
			}

			if !test.expectedErr && test.envVarName != "" {
				val, present := os.LookupEnv(test.envVarName)
				if !present || val != test.envVarValue {
					t.Errorf("Expected env var value: %q, got: %q", test.envVarValue, val)
				}
			}
		})
	}
}
