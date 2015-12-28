package main

import (
	"log"

	"github.com/ivotron/vio"
	"github.com/spf13/cobra"
)

var snapPath string
var backend string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initializes the vio repo",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := vio.Init(snapPath, backend); err != nil {
			log.Fatalln(err.Error())
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&snapPath,
		"snapshots", "s", ".snapshots", "Path to where snapshots are stored")
	initCmd.Flags().StringVarP(&backend,
		"backend", "b", "posix", "Backend to manage snapshots (one of 'git', 'posix' or 'git-lfs')")
}
