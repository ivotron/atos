package main

import (
	"log"

	"github.com/ivotron/vio"
	"github.com/spf13/cobra"
)

var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "Checks out a version.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalln("Expecting revision ID")
		}
		if err := vio.Checkout(args[0]); err != nil {
			log.Fatalln(err.Error())
		}
	},
}

func init() {
	RootCmd.AddCommand(checkoutCmd)
}
