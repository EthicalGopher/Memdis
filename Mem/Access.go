package Mem

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/EthicalGopher/Memdis/core"
	"github.com/EthicalGopher/Memdis/persistence"
)

// DB represents the database instance, holding the engine and persistence layer.
type DB struct {
	engine *core.Engine
	wal    *persistence.WAL
}

// Connect initializes and returns a new database instance.
func Connect(filePath string) (*DB, error) {
	fmt.Println("🚀 Initializing DocStore...")

	wal, err := persistence.NewWAL(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create WAL: %w", err)
	}

	engine := core.NewEngine()

	if err := wal.Restore(engine); err != nil {
		log.Printf("Note: Starting with a fresh database: %v", err)
	}

	return &DB{
		engine: engine,
		wal:    wal,
	}, nil
}

// Close gracefully shuts down the database.
func (db *DB) Close() error {
	fmt.Println("👋 Shutting down database...")
	return db.wal.Close()
}

// Execute parses and runs a single command.
func (db *DB) Execute(commandStr string) (any, error) {
	input := strings.TrimSpace(commandStr)
	if input == "" {
		return nil, nil
	}

	parts := strings.SplitN(input, " ", 4)
	command := strings.ToUpper(parts[0])

	switch command {
	// ... (INSERT, FIND, UPDATE, DELETE, COUNT, SORT cases remain the same) ...
	case "INSERT":
		if len(parts) < 3 {
			return nil, fmt.Errorf("❌ usage: INSERT <collection> <json_data>")
		}
		collection, jsonData := parts[1], parts[2]

		var data core.Document
		if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
			return nil, fmt.Errorf("❌ invalid JSON: %w", err)
		}

		id := core.GenerateID()
		cmd := core.Command{Op: "insert", Collection: collection, Data: data, ID: id}

		if err := db.wal.Write(cmd); err != nil {
			return nil, fmt.Errorf("❌ failed to persist command: %w", err)
		}

		db.engine.ApplyCommand(cmd)
		return fmt.Sprintf("✅ Document inserted into '%s'", collection), nil

	case "FIND":
		if len(parts) < 2 {
			return nil, fmt.Errorf("❌ usage: FIND <collection> [filter_json]")
		}
		collection := parts[1]
		var filter core.Document

		if len(parts) >= 3 {
			if err := json.Unmarshal([]byte(parts[2]), &filter); err != nil {
				return nil, fmt.Errorf("❌ invalid filter JSON: %w", err)
			}
		}

		results := db.engine.Find(collection, filter)
		return results, nil

	case "UPDATE":
		if len(parts) < 4 {
			return nil, fmt.Errorf("❌ usage: UPDATE <collection> <filter_json> <update_json>")
		}
		collection, filterJson, updateJson := parts[1], parts[2], parts[3]

		var filter, updateData core.Document
		if err := json.Unmarshal([]byte(filterJson), &filter); err != nil {
			return nil, fmt.Errorf("❌ invalid filter JSON: %w", err)
		}
		if err := json.Unmarshal([]byte(updateJson), &updateData); err != nil {
			return nil, fmt.Errorf("❌ invalid update JSON: %w", err)
		}

		cmd := core.Command{Op: "update", Collection: collection, Filter: filter, Data: updateData}
		if err := db.wal.Write(cmd); err != nil {
			return nil, fmt.Errorf("❌ failed to persist command: %w", err)
		}

		db.engine.ApplyCommand(cmd)
		return fmt.Sprintf("✅ Documents updated in '%s'", collection), nil

	case "DELETE":
		if len(parts) < 3 {
			return nil, fmt.Errorf("❌ usage: DELETE <collection> <filter_json>")
		}
		collection, filterJson := parts[1], parts[2]

		var filter core.Document
		if err := json.Unmarshal([]byte(filterJson), &filter); err != nil {
			return nil, fmt.Errorf("❌ invalid filter JSON: %w", err)
		}

		cmd := core.Command{Op: "delete", Collection: collection, Filter: filter}
		if err := db.wal.Write(cmd); err != nil {
			return nil, fmt.Errorf("❌ failed to persist command: %w", err)
		}

		db.engine.ApplyCommand(cmd)
		return fmt.Sprintf("✅ Documents deleted from '%s'", collection), nil

	case "COUNT":
		if len(parts) < 2 {
			return nil, fmt.Errorf("❌ usage: COUNT <collection> [filter_json]")
		}
		collection := parts[1]
		var filter core.Document

		if len(parts) >= 3 {
			if err := json.Unmarshal([]byte(parts[2]), &filter); err != nil {
				return nil, fmt.Errorf("❌ invalid filter JSON: %w", err)
			}
		}

		count := db.engine.Count(collection, filter)
		return count, nil

	case "SORT":
		if len(parts) < 3 {
			return nil, fmt.Errorf("❌ usage: SORT <collection> <sort_key>")
		}
		collection, key := parts[1], parts[2]
		docs := db.engine.Sort(collection, key)
		return docs, nil

	case "SAVE":
		log.Println("⚙️ Starting database snapshot...")

		// 1. Save the current engine state to the snapshot file.
		if err := db.wal.SaveSnapshot(db.engine); err != nil {
			return nil, fmt.Errorf("❌ snapshot failed: %w", err)
		}

		// 2. If snapshot is successful, truncate the WAL.
		if err := db.wal.Truncate(); err != nil {
			// This is non-fatal for the user, but should be logged.
			// The next snapshot will just have to cover more data.
			log.Printf("⚠️ Warning: snapshot successful, but failed to truncate WAL: %v", err)
		}

		log.Println("✅ Snapshot created successfully.")
		return "✅ Snapshot created successfully.", nil

	case "LIST_COLLECTIONS":
		return "📊 Collections feature coming soon!", nil

	case "EXIT", "QUIT":
		return "Command 'QUIT' received.", nil

	default:
		return nil, fmt.Errorf("❌ unknown command: %s", command)
	}
}
