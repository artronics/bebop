package cmd

import (
	"fmt"
	"github.com/nhsdigital/bebop-cli/pkg"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var apigeeLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := pkg.ApigeeConfig{}

		config.Username = cmd.Flags().Lookup("username").Value.String()

		passwordF := cmd.Flags().Lookup("password")
		if passwordF.Changed {
			config.Password = passwordF.Value.String()
		} else if password, ok := os.LookupEnv("APIGEE_PASSWORD"); ok {
			config.Password = password
		} else {
			log.Fatalln("password is not provided. Either use --password flag or provide APIGEE_PASSWORD environment variable")
		}

		if cmd.Flags().Lookup("mfa").Changed {
			config.Mfa = cmd.Flags().Lookup("mfa").Value.String()
		}
		config.OauthTokenUrl = cmd.Flags().Lookup("oauth-token-url").Value.String()

		token, err := pkg.ApigeeLogin(config)
		if err != nil {
			log.Fatalln(err.Error())
		}

		fmt.Println(token)
	},
}

func init() {
	// TODO: change common flags to Persistent and define them in root. Currently we can redefine short name for username because it's taken by aws
	apigeeLoginCmd.Flags().String("username", "", `apigee account username.
	It should be an email address.`)
	_ = apigeeLoginCmd.MarkFlagRequired("username")

	apigeeLoginCmd.Flags().StringP("password", "p", "", `apigee account password. Also can be provided
via environment variable APIGEE_PASSWORD`)

	//apigeeLoginCmd.Flags().StringP("access_token", "t", "", `apigee access token.
	//If present there is no need for username/password. It can also be provided via environment variable APIGEE_ACCESS_TOKEN`)

	apigeeLoginCmd.Flags().StringP("mfa", "m", "", "MFA (Multi-Factor Authentication) code")

	apigeeLoginCmd.Flags().String("oauth-token-url", "https://login.apigee.com/oauth/token", "oauth token endpoint url.")

	apigeeCmd.AddCommand(apigeeLoginCmd)
}
