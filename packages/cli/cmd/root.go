package cmd

import (
	"errors"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/CircleCI-Public/slack-orb-go/packages/cli/config"
)

var SlackConfig config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "slack-orb-cli",
	Short: "The slack-orb-cli interface for the CircleCI Slack orb",
	Long:  `The slack-orb-cli by CircleCI is a command-line tool for sending slack notifications as a part of a CI/CD workflow.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	viper.SetConfigName("config")
	SlackConfig, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error loading environment configuration: \n%v\n", err)
	}
	if err := SlackConfig.Validate(); err != nil {
		handleConfigurationError(err)
	}

}

func handleConfigurationError(err error) {
	var envVarError *config.EnvVarError
	if errors.As(err, &envVarError) {
		switch envVarError.VarName {
		case "SLACK_ACCESS_TOKEN":
			log.Fatalf(
				"In order to use the Slack Orb an OAuth token must be present via the SLACK_ACCESS_TOKEN environment variable." +
					"\nFollow the setup guide available in the wiki: https://github.com/CircleCI-Public/slack-orb/wiki/Setup.",
			)
		case "SLACK_PARAM_CHANNEL":
			//nolint:lll // user message
			log.Fatalf(
				`No channel was provided. Please provide one or more channels using the "SLACK_PARAM_CHANNEL" environment variable or the "channel" parameter.`,
			)
		default:
			log.Fatalf("Configuration validation failed: Environment variable not set: %s", envVarError.VarName)
		}
	}

	log.Fatal(err)
}