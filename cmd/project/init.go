package project

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var (
	name        string
	projectType string
	framework   string
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "init",
		Short:            "The initailized project backend server",
		Example:          "gwb init or gwb init -n backend or gwb init -name backend",
		TraverseChildren: true,
		Run:              run,
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "set a project name.")
	cmd.Flags().StringVarP(&projectType, "type", "t", "", "set a project type.")
	return cmd
}

func run(cmd *cobra.Command, args []string) {
	if name == "" {
		inputProjectName()
	}

	// Select a project type
	if projectType == "" {
		selectType()
	}

	switch projectType {
	case "web":
		selectFramework()
	case "microservice":
		fmt.Println("make microservice project")
	default:
		fmt.Fprintln(os.Stderr, "error process")
		os.Exit(1)
	}
	// init module
	gomodinit()
	// install packages
	pkgs := getPackages(framework)
	if pkgs == nil {
		fmt.Fprintln(os.Stderr, "install packages error")
		os.Exit(1)
	}
}

func inputProjectName() {
	prompt := &survey.Input{
		Message: "Please enter the project name:",
	}
	survey.AskOne(prompt, &name)
}

func selectType() {
	prompt := &survey.Select{
		Message: "Select the type of project you want to create:",
		Options: []string{"web", "microservice"},
	}
	survey.AskOne(prompt, &projectType, survey.WithValidator(survey.Required))
}

func selectFramework() {
	prompt := &survey.Select{
		Message: "Please choose a web framework for your project:",
		Options: []string{"fiber", "gin"},
	}
	survey.AskOne(prompt, &framework, survey.WithValidator(survey.Required))
}

func getPackages(cond string) []string {
	switch cond {
	case "fiber":
		return []string{
			"github.com/gofiber/fiber/v2",
		}
	case "gin":
		return []string{
			"github.com/gin-gonic/gin",
			"github.com/gin-contrib/cors",
			"github.com/gin-contrib/gzip",
		}
	}
	return nil
}

func gomodinit() {
	cmd := exec.Command("go", "mod", "init", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
