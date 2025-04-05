package cmd

import (
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of grawl",
	Long:  `All software has versions. This is grawl's`,
	Run: func(cmd *cobra.Command, args []string) {
		// Print the version number
		// This is a placeholder, replace with actual version retrieval logic
		version := "v0.1.0"
		println("grawl version:", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
