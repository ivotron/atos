package main

import (
	"log"

	"github.com/ivotron/vio"
	"github.com/spf13/cobra"
)

var meta string

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Create a commit for unversioned files.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := vio.Commit(meta); err != nil {
			log.Fatalln(err.Error())
		}
	},
}

func init() {
	RootCmd.AddCommand(commitCmd)
	initCmd.Flags().StringVarP(&meta,
		"meta", "", "{}", "JSON-formatted string of key-value pairs")
}
