package cmd

import (
	"fmt"
	"strings"

	"github.com/EthicalGopher/Memdis/Mem"
	"github.com/spf13/cobra"
)

var countCmd = &cobra.Command{
	Use:   "count [collection] [filter_json]",
	Short: "Count documents in a collection",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		collection := args[0]
		filter := ""
		if len(args) > 1 {
			filter = args[1]
		}

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

		cmdStr := fmt.Sprintf("COUNT %s %s", collection, filter)
		result, err := DB.Execute(strings.TrimSpace(cmdStr))
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(result)
	},
}

func AddCountCommand(root *cobra.Command) {
	root.AddCommand(countCmd)
}
