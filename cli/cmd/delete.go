package cmd

import (
	"fmt"
	"strings"

	"github.com/EthicalGopher/Memdis/Mem"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [collection] [filter_json]",
	Short: "Delete documents from a collection",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		collection := args[0]
		filterJson := args[1]

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

		cmdStr := fmt.Sprintf("DELETE %s %s", collection, filterJson)
		result, err := DB.Execute(strings.TrimSpace(cmdStr))
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(result)
	},
}

func AddDeleteCommand(root *cobra.Command) {
	root.AddCommand(deleteCmd)
}
