package main

import (
	"encoding/json"
	"fmt"
	"github.com/EthicalGopher/Memdis/Mem"
	"github.com/EthicalGopher/Memdis/core"
)

type Data struct {
	Name string `json:"name"`
}

func main() {

	DB, err := Mem.Connect("data.mem")
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		err := DB.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	data := Data{
		Name: "Ayush",
	}
	_, err = json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	cmd := `FIND test `
	result, err := DB.Execute(cmd)
	if err != nil {
		fmt.Println(err)
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

}
