package cmd

import (
	"fmt"
	"gwb/cmd/auth"
	"gwb/cmd/project"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gwb",
		Short: "GWB is a very fast go backend web server generator",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func Execute() {
	rootCmd.AddCommand(project.Command())
	rootCmd.AddCommand(auth.CreateUserCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
