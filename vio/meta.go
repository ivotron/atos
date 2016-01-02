package main

import "github.com/spf13/cobra"

var revision string
var timestamp string

var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: "Show meta info.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var metaGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Obtain value for key.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var metaSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Assign value to key.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var metaListCmd = &cobra.Command{
	Use:   "list",
	Short: "List keys.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	RootCmd.AddCommand(metaCmd)
	meta.AddCommand(metaGetCmd)
	meta.AddCommand(metaSetCmd)
	meta.PersistentFlags().StringVarP()
	commitCmd.Flags().StringVarP(&meta,
		"revision", "r", "{}", "JSON-formatted string of key-value pairs.")
	commitCmd.Flags().StringVarP(&msg,
		"timestamp", "t", "", "Commit message.")
}
