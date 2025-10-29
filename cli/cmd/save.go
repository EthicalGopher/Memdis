package cmd

import (
	"fmt"

	"github.com/EthicalGopher/Memdis/Mem"
	"github.com/spf13/cobra"
)

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save the database snapshot",
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

		result, err := DB.Execute("SAVE")
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(result)
	},
}

func AddSaveCommand(root *cobra.Command) {
	root.AddCommand(saveCmd)
}
