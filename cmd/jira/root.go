package jira

import (
	"fmt"
	"os"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const clientTimeout = 15 * time.Second

var (
	config     string
	jiraClient *jira.Client

	rootCmd = &cobra.Command{
		Use:   "jira",
		Short: "Jira cli designed for developers",
		Long:  `A jira command line designed for developers to help with frequent jira chores.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
)

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig, initJiraClient)

	rootCmd.PersistentFlags().StringVar(&config, "config", "", "config file (default is $HOME/.jira.yml)")
}

func initConfig() {
	if config != "" {
		viper.SetConfigFile(config)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			exitWithError(err)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".jira")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("jira")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initJiraClient() {
	config := jira.Config{
		Server:   viper.Get("server").(string),
		Login:    viper.Get("login").(string),
		APIToken: viper.Get("api_token").(string),
	}

	jiraClient = jira.NewClient(config, jira.WithTimeout(clientTimeout))
}

func exitWithError(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
	os.Exit(1)
}