package auth

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var (
	username string
	password string
	email    string
)

func CreateUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "createuser",
		Short:            "Create your account",
		TraverseChildren: true,
		Run:              run,
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) {
	inputUsername()
	inputPassword()
	inputEmail()
}

func inputUsername() {
	prompt := &survey.Input{
		Message: "Please enter your username:",
	}
	survey.AskOne(prompt, &username, survey.WithValidator(survey.Required))
}

func inputPassword() {
	prompt := &survey.Password{
		Message: "Please enter your password:",
	}
	survey.AskOne(prompt, &password, survey.WithValidator(survey.Required))
}

func inputEmail() {
	prompt := &survey.Password{
		Message: "Please enter your email:",
	}
	survey.AskOne(prompt, &email)
}
