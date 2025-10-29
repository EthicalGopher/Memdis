package cmd

import (
	"fmt"
	"strings"

	"github.com/EthicalGopher/Memdis/Mem"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update [collection] [filter_json] [update_json]",
	Short: "Update documents in a collection",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		collection := args[0]
		filterJson := args[1]
		updateJson := args[2]

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

		cmdStr := fmt.Sprintf("UPDATE %s %s %s", collection, filterJson, updateJson)
		result, err := DB.Execute(strings.TrimSpace(cmdStr))
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(result)
	},
}

func AddUpdateCommand(root *cobra.Command) {
	root.AddCommand(updateCmd)
}
