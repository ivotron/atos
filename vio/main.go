package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "vio",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", cmd.UsageString())
	},
}

func main() {
	// check if have access to dependencies
	if _, err := exec.Command("git", "--version").Output(); err != nil {
		log.Fatalln("Unable to execute 'git': " + err.Error())
	}
	if _, err := exec.Command("rsync", "--version").Output(); err != nil {
		log.Fatalln("Unable to execute 'rsync': " + err.Error())
	}
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
