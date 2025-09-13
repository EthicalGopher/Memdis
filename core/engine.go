package core

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"
)

// Document is a generic JSON-like document
type Document map[string]interface{}

// Command represents a database operation
type Command struct {
	Op         string   // "insert", "update", "delete"
	Collection string   // Like a table in SQL, collection in NoSQL
	Data       Document // The document data
	Filter     Document // For update/delete operations
	ID         string   // Optional specific ID
}

// Engine is our document database
type Engine struct {
	mu          sync.RWMutex
	collections map[string]map[string]Document // collection -> id -> document
}

// NewEngine creates a new document store
func NewEngine() *Engine {
	return &Engine{
		collections: make(map[string]map[string]Document),
	}
}

// ApplyCommand applies a command to the database
func (e *Engine) ApplyCommand(cmd Command) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Ensure the collection exists
	if _, exists := e.collections[cmd.Collection]; !exists {
		e.collections[cmd.Collection] = make(map[string]Document)
	}
	collection := e.collections[cmd.Collection]

	switch cmd.Op {
	case "insert":
		id := cmd.ID
		if id == "" {
			id = GenerateID()
		}
		collection[id] = cmd.Data
		collection[id]["_id"] = id

	case "update":
		for id, doc := range collection {
			if matchesFilter(doc, cmd.Filter) {
				for k, v := range cmd.Data {
					doc[k] = v
				}
				collection[id] = doc
			}
		}

	case "delete":
		for id, doc := range collection {
			if matchesFilter(doc, cmd.Filter) {
				delete(collection, id)
			}
		}
	}
}

// Serialize converts the entire engine state into a byte slice for snapshotting.
func (e *Engine) Serialize() ([]byte, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return json.Marshal(e.collections)
}

// Deserialize populates the engine from a byte slice when loading a snapshot.
func (e *Engine) Deserialize(data []byte) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	var collections map[string]map[string]Document
	if err := json.Unmarshal(data, &collections); err != nil {
		return err
	}
	e.collections = collections
	return nil
}

// Find documents in a collection
func (e *Engine) Find(collectionName string, filter Document) []Document {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var results []Document
	collection, exists := e.collections[collectionName]
	if !exists {
		return results
	}

	for _, doc := range collection {
		if matchesFilter(doc, filter) {
			results = append(results, doc)
		}
	}
	return results
}

// Count the number of items
func (e *Engine) Count(collectionName string, filter Document) int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	collection, exists := e.collections[collectionName]
	if !exists {
		return 0
	}
	if len(filter) == 0 {
		return len(collection)
	}
	count := 0
	for _, doc := range collection {
		if matchesFilter(doc, filter) {
			count++
		}
	}
	return count
}

// Sort documents in a collection by a specific key
func (e *Engine) Sort(collectionName string, sortKey string) []Document {
	e.mu.RLock()
	defer e.mu.RUnlock()

	collection, exists := e.collections[collectionName]
	if !exists {
		return []Document{}
	}

	docs := make([]Document, 0, len(collection))
	for _, doc := range collection {
		docs = append(docs, doc)
	}

	sort.Slice(docs, func(i, j int) bool {
		valI, iExists := docs[i][sortKey]
		valJ, jExists := docs[j][sortKey]

		if !iExists {
			return false
		}
		if !jExists {
			return true
		}

		switch vI := valI.(type) {
		case float64:
			if vJ, ok := valJ.(float64); ok {
				return vI < vJ
			}
		case string:
			if vJ, ok := valJ.(string); ok {
				return vI < vJ
			}
		case int:
			if vJ, ok := valJ.(int); ok {
				return vI < vJ
			}
		}
		return false
	})

	return docs
}

// Helper function to check if document matches filter
func matchesFilter(doc Document, filter Document) bool {
	if len(filter) == 0 {
		return true
	}

	for key, filterValue := range filter {
		docValue, exists := doc[key]
		if !exists || docValue != filterValue {
			return false
		}
	}
	return true
}

// GenerateID creates a new unique ID.
func GenerateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
