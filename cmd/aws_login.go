package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/nhsdigital/bebop-cli/pkg"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var awsLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Assumes an AWS role based on provide profile and give MFA token",
	Long:  `TODO:`,
	Run: func(cmd *cobra.Command, args []string) {
		config := pkg.AwsConfig{}

		profileFlag := cmd.Flags().Lookup("profile")
		if profileFlag.Changed { // user explicitly sets the --profile flag
			config.Profile = profileFlag.Value.String()
		} else if profile, ok := os.LookupEnv("AWS_PROFILE"); ok {
			config.Profile = profile
		} else { // whatever the default value is
			config.Profile = profileFlag.Value.String()
		}
		config.AccountId = cmd.Flags().Lookup("account").Value.String()
		config.Username = cmd.Flags().Lookup("username").Value.String()
		config.Mfa = cmd.Flags().Lookup("mfa").Value.String()

		data, err := pkg.AwsLogin(config)
		if err != nil {
			log.Fatalln(err.Error())
		}

		outProfile := cmd.Flags().Lookup("out-profile")
		if outProfile.Changed {

			err = overwriteProfile(outProfile.Value.String(), data)
			if err != nil {
				log.Fatalln(err.Error())
			}

		} else {
			fmt.Println(data)
		}
	},
}

func init() {
	awsLoginCmd.Flags().StringP("account", "a", "", "aws account ID. It should be a 12-digit number")
	_ = awsLoginCmd.MarkFlagRequired("account")

	awsLoginCmd.Flags().StringP("username", "u", "", "aws username. it should be your nhs email username by default")
	_ = awsLoginCmd.MarkFlagRequired("username")

	awsLoginCmd.Flags().StringP("mfa", "m", "", "MFA (Multi-Factor Authentication) code")
	_ = awsLoginCmd.MarkFlagRequired("mfa")

	awsLoginCmd.Flags().String("out-profile", "", `If set, then login will create a new profile or update
existing one. If this value is not set then it will fallback to aws cli default format`)

	awsCmd.AddCommand(awsLoginCmd)
}

func overwriteProfile(profile string, data string) error {
	cred := pkg.AwsLoginData{}
	err := json.Unmarshal([]byte(data), &cred)
	if err != nil {
		return err
	}

	err = pkg.SetCredProfile(profile, cred)
	if err != nil {
		return err
	}

	return nil
}
