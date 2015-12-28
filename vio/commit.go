package main

import (
	"log"

	"github.com/ivotron/vio"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Create a commit for unversioned files.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := vio.Commit(); err != nil {
			log.Fatalln(err.Error())
		}
	},
}

func init() {
	RootCmd.AddCommand(commitCmd)
}
