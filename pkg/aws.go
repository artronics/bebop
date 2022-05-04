package pkg

import (
	"fmt"
	"github.com/nhsdigital/bebop-cli/internal"
)

type AwsConfig struct {
	Profile   string
	AccountId string
	Username  string
	Mfa       string
}

type credentials struct {
	AccessKeyId     string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
	SessionToken    string `json:"SessionToken"`
	Expiration      string `json:"Expiration"`
}

type AwsLoginData struct {
	Credentials credentials `json:"Credentials"`
}

func Login(config AwsConfig) (_ string, err error) {
	mfa := fmt.Sprintf("arn:aws:iam::%s:mfa/%s", config.AccountId, config.Username)
	args := []string{"--profile", config.Profile, "sts", "get-session-token",
		"--serial-number", mfa, "--token-code", config.Mfa, "--duration-seconds", "129600", "--output", "json"}

	return internal.JustRun("aws", args)
}

func SetCredProfile(profile string, data AwsLoginData) error {
	args := []string{"configure", "--profile", profile, "set", "aws_access_key_id", data.Credentials.AccessKeyId}
	_, err := internal.JustRun("aws", args)
	if err != nil {
		return err
	}

	args = []string{"configure", "--profile", profile, "set", "aws_secret_access_key", data.Credentials.SecretAccessKey}
	_, err = internal.JustRun("aws", args)
	if err != nil {
		return err
	}

	args = []string{"configure", "--profile", profile, "set", "aws_session_token", data.Credentials.SessionToken}
	_, err = internal.JustRun("aws", args)
	if err != nil {
		return err
	}

	return nil
}
