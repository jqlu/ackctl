package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/jqlu/ackctl/config"
	"github.com/spf13/cobra"
	"log"
)

var qs = []*survey.Question{
	{
		Name:     "accessKeyId",
		Prompt:   &survey.Password{Message: "Access Key Id:"},
		Validate: survey.Required,
	},
	{
		Name:     "accessKeySecret",
		Prompt:   &survey.Password{Message: "Access Key Secret:"},
		Validate: survey.Required,
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure access to Aliyun",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		input := config.AkConfig{}

		err := survey.Ask(qs, &input)
		if err != nil {
			log.Fatalf("Failed to parse input: %v", err)
		}

		if err := config.UpdateConfigFile(input); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Ackctl configured.")
	},
}