package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/nhsdigital/bebop-cli/pkg"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		data, err := pkg.Login(config)
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
	loginCmd.Flags().StringP("account", "a", "", "aws account ID. It should be a 12-digit number")
	_ = loginCmd.MarkFlagRequired("account")

	loginCmd.Flags().StringP("username", "u", "", "aws username. it should be your nhs email username by default")
	_ = loginCmd.MarkFlagRequired("username")

	loginCmd.Flags().StringP("mfa", "m", "", "MFA (Multi-Factor Authentication) code")
	_ = loginCmd.MarkFlagRequired("mfa")

	loginCmd.Flags().String("out-profile", "", `If set, then login will create a new profile or update
existing one. If this value is not set then it will fallback to aws cli default format`)

	awsCmd.AddCommand(loginCmd)
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
