package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/EthicalGopher/Memdis/Mem"
	"github.com/EthicalGopher/Memdis/core"
	"github.com/spf13/cobra"
)

var sortCmd = &cobra.Command{
	Use:   "sort [collection] [sort_key]",
	Short: "Sort documents in a collection",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		collection := args[0]
		sortKey := args[1]

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

		cmdStr := fmt.Sprintf("SORT %s %s", collection, sortKey)
		result, err := DB.Execute(strings.TrimSpace(cmdStr))
		if err != nil {
			fmt.Println(err)
			return
		}

		if _, ok := result.(string); ok {
			fmt.Println(result.(string))
			return
		}

		docs := result.([]core.Document)
		var jsonByte []byte
		for _, doc := range docs {
			jsonByte, err = json.MarshalIndent(doc, " ", " ")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(jsonByte))
		}
	},
}

func AddSortCommand(root *cobra.Command) {
	root.AddCommand(sortCmd)
}
