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

// NewDB initializes and returns a new database instance.
// It creates the engine, sets up the WAL, and restores state from the log.
func Connect(filePath string) (*DB, error) {
	fmt.Println("🚀 Initializing DocStore...")

	// Initialize persistence (Write-Ahead Log)
	wal, err := persistence.NewWAL(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create WAL: %w", err)
	}

	// Initialize the core storage engine
	engine := core.NewEngine()

	// Restore state from the WAL. This is critical for persistence.
	if err := wal.Restore(engine); err != nil {
		// This is often not a fatal error, just means we're starting fresh.
		log.Printf("Note: Starting with a fresh database: %v", err)
	} else {
		fmt.Println("✅ Database state restored from log.")
	}

	return &DB{
		engine: engine,
		wal:    wal,
	}, nil
}

// Close gracefully shuts down the database, ensuring the WAL is closed.
func (db *DB) Close() error {
	fmt.Println("👋 Shutting down database...")
	return db.wal.Close()
}

// Execute parses and runs a single command against the database.
// It returns the result of the command (e.g., found documents, count) or an error.
func (db *DB) Execute(commandStr string) (any, error) {
	input := strings.TrimSpace(commandStr)
	if input == "" {
		return nil, nil // No command is a no-op, not an error.
	}

	parts := strings.SplitN(input, " ", 4)
	command := strings.ToUpper(parts[0])

	switch command {
	case "INSERT":
		if len(parts) < 3 {
			return nil, fmt.Errorf("❌ usage: INSERT <collection> <json_data>")
		}
		collection, jsonData := parts[1], parts[2]

		var data core.Document
		if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
			return nil, fmt.Errorf("❌ invalid JSON: %w", err)
		}

		// --- CRITICAL FIX ---
		// Generate the ID *before* writing to the WAL. This makes the operation
		// deterministic and ensures the restore process is reliable.
		id := core.GenerateID()
		cmd := core.Command{Op: "insert", Collection: collection, Data: data, ID: id}
		// --- END FIX ---

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
		return results, nil // Return the slice of documents

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
		return count, nil // Return the integer count

	case "SORT":
		if len(parts) < 3 {
			return nil, fmt.Errorf("❌ usage: SORT <collection> <sort_key>")
		}
		collection, key := parts[1], parts[2]
		docs := db.engine.Sort(collection, key)
		return docs, nil // Return the sorted slice of documents

	case "LIST_COLLECTIONS":
		return "📊 Collections feature coming soon!", nil

	case "EXIT", "QUIT":
		// The caller should handle this command to exit its own loop.
		return "Command 'QUIT' received.", nil

	default:
		return nil, fmt.Errorf("❌ unknown command: %s", command)
	}
}
