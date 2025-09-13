package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/EthicalGopher/Memdis/core"
	"github.com/EthicalGopher/Memdis/persistence"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Println("🚀 Starting DocStore (MongoDB-like NoSQL Database)...")

	// Initialize persistence
	wal, err := persistence.NewWAL("data.mem")
	if err != nil {
		log.Fatalf("Failed to create WAL: %v", err)
	}
	defer wal.Close()

	// Initialize engine
	engine := core.NewEngine()

	// Restore state
	if err := wal.Restore(engine); err != nil {
		log.Printf("Note: Starting fresh database: %v", err)
	} else {
		fmt.Println("✅ Database state restored from log")
	}

	fmt.Println("💡 Commands:")
	fmt.Println("  INSERT <collection> <json>")
	fmt.Println("  FIND <collection> [filter_json]")
	fmt.Println("  UPDATE <collection> <filter_json> <update_json>")
	fmt.Println("  DELETE <collection> <filter_json>")
	fmt.Println("  LIST_COLLECTIONS")
	fmt.Println("  EXIT")
	fmt.Println("---")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.SplitN(input, " ", 4)
		command := strings.ToUpper(parts[0])

		switch command {
		case "INSERT":
			if len(parts) < 3 {
				fmt.Println("❌ Usage: INSERT <collection> <json_data>")
				continue
			}
			collection, jsonData := parts[1], parts[2]

			var data core.Document
			if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
				fmt.Printf("❌ Invalid JSON: %v\n", err)
				continue
			}

			cmd := core.Command{Op: "insert", Collection: collection, Data: data}
			if err := wal.Write(cmd); err != nil {
				fmt.Printf("❌ Failed to persist: %v\n", err)
				continue
			}

			engine.ApplyCommand(cmd)
			fmt.Printf("✅ Document inserted into '%s'\n", collection)

		case "FIND":
			if len(parts) < 2 {
				fmt.Println("❌ Usage: FIND <collection> [filter_json]")
				continue
			}
			collection := parts[1]
			var filter core.Document

			if len(parts) >= 3 {
				if err := json.Unmarshal([]byte(parts[2]), &filter); err != nil {
					fmt.Printf("❌ Invalid filter JSON: %v\n", err)
					continue
				}
			}

			results := engine.Find(collection, filter)
			if len(results) == 0 {
				fmt.Println("📭 No documents found")
			} else {
				fmt.Printf("📋 Found %d documents:\n", len(results))
				for i, doc := range results {
					jsonBytes, _ := json.MarshalIndent(doc, "  ", "  ")
					fmt.Printf("%d: %s\n", i+1, string(jsonBytes))
				}
			}

		case "UPDATE":
			if len(parts) < 4 {
				fmt.Println("❌ Usage: UPDATE <collection> <filter_json> <update_json>")
				continue
			}
			collection, filterJson, updateJson := parts[1], parts[2], parts[3]

			var filter, updateData core.Document
			if err := json.Unmarshal([]byte(filterJson), &filter); err != nil {
				fmt.Printf("❌ Invalid filter JSON: %v\n", err)
				continue
			}
			if err := json.Unmarshal([]byte(updateJson), &updateData); err != nil {
				fmt.Printf("❌ Invalid update JSON: %v\n", err)
				continue
			}

			cmd := core.Command{Op: "update", Collection: collection, Filter: filter, Data: updateData}
			if err := wal.Write(cmd); err != nil {
				fmt.Printf("❌ Failed to persist: %v\n", err)
				continue
			}

			engine.ApplyCommand(cmd)
			fmt.Printf("✅ Documents updated in '%s'\n", collection)

		case "DELETE":
			if len(parts) < 3 {
				fmt.Println("❌ Usage: DELETE <collection> <filter_json>")
				continue
			}
			collection, filterJson := parts[1], parts[2]

			var filter core.Document
			if err := json.Unmarshal([]byte(filterJson), &filter); err != nil {
				fmt.Printf("❌ Invalid filter JSON: %v\n", err)
				continue
			}

			cmd := core.Command{Op: "delete", Collection: collection, Filter: filter}
			if err := wal.Write(cmd); err != nil {
				fmt.Printf("❌ Failed to persist: %v\n", err)
				continue
			}

			engine.ApplyCommand(cmd)
			fmt.Printf("✅ Documents deleted from '%s'\n", collection)

		case "COUNT":
			if len(parts) < 3 {
				fmt.Println("❌ Usage: Count <collection> <filter_json>")
				continue
			}
			collection, filterJson := parts[1], parts[2]

			count := core.Count(collection, filterJson)
			fmt.Printf("Number of items inside %s : %d\n", collection, count)

		case "LIST_COLLECTIONS":
			// This would need to be implemented in the engine
			fmt.Println("📊 Collections feature coming soon!")

		case "EXIT", "QUIT":
			fmt.Println("👋 Goodbye!")
			return

		default:
			fmt.Println("❌ Unknown command")
		}
	}
}
