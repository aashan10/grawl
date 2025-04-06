package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Version = "0.1.3"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the build details for grawl",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("Grawl version: %s", Version)

	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
