package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version info.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("vio version 0.1.0\n")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
