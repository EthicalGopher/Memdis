package cmd

import (
	"fmt"
	"strings"

	"github.com/EthicalGopher/Memdis/Mem"
	"github.com/spf13/cobra"
)

var insertCmd = &cobra.Command{
	Use:   "insert [collection] [json_data]",
	Short: "Insert a document into a collection",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		collection := args[0]
		jsonData := args[1]

		DB, err := Mem.Connect("data.mem")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer func() {
			err := DB.Close()
			if err != nil {
				fmt.Println(err)
			}
		}()

		cmdStr := fmt.Sprintf("INSERT %s %s", collection, jsonData)
		result, err := DB.Execute(strings.TrimSpace(cmdStr))
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(result)
	},
}

func AddInsertCommand(root *cobra.Command) {
	root.AddCommand(insertCmd)
}
