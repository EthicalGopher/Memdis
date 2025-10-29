package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/EthicalGopher/Memdis/Mem"
	"github.com/EthicalGopher/Memdis/core"
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find [collection] [filter]",
	Short: "Find documents in a collection",
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

		cmdStr := fmt.Sprintf("FIND %s %s", collection, filter)
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

func AddFindCommand(root *cobra.Command) {
	root.AddCommand(findCmd)
}
