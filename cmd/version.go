package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of unseal",
	Long:  `All software has versions. This is unseal's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("unseal v0.2-beta")
	},
}
