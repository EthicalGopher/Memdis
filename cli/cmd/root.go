package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "memdis",
	Short: "Memdis is a simple in-memory database",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Memdis!")
	},
}

func Execute() {
	AddFindCommand(rootCmd)
	AddInsertCommand(rootCmd)
	AddUpdateCommand(rootCmd)
	AddDeleteCommand(rootCmd)
	AddCountCommand(rootCmd)
	AddSortCommand(rootCmd)
	AddSaveCommand(rootCmd)
	AddListCollectionsCommand(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
