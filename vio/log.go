package main

import (
	"fmt"
	"log"

	"github.com/ivotron/vio"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show log info.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		logstr, err := vio.Log()
		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Print(logstr)
	},
}

func init() {
	RootCmd.AddCommand(logCmd)
}
