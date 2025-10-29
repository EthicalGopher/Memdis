package cmd

import (
	"fmt"

	"github.com/EthicalGopher/Memdis/Mem"
	"github.com/spf13/cobra"
)

var listCollectionsCmd = &cobra.Command{
	Use:   "list-collections",
	Short: "List all collections",
	Run: func(cmd *cobra.Command, args []string) {
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

		result, err := DB.Execute("LIST_COLLECTIONS")
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(result)
	},
}

func AddListCollectionsCommand(root *cobra.Command) {
	root.AddCommand(listCollectionsCmd)
}
