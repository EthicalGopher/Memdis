package core

import (
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
		// The ID should now be pre-generated, but we keep this for compatibility
		// with old WAL entries that might still be processed.
		id := cmd.ID
		if id == "" {
			id = GenerateID() // This was the source of the non-deterministic restore
		}
		collection[id] = cmd.Data
		// Add the _id field to the document
		collection[id]["_id"] = id

	case "update":
		// Simple implementation: update all documents that match filter
		for id, doc := range collection {
			if matchesFilter(doc, cmd.Filter) {
				// Merge existing document with update data
				for k, v := range cmd.Data {
					doc[k] = v
				}
				collection[id] = doc
			}
		}

	case "delete":
		// Remove documents that match filter
		for id, doc := range collection {
			if matchesFilter(doc, cmd.Filter) {
				delete(collection, id)
			}
		}
	}
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

	// Convert map to slice for sorting
	docs := make([]Document, 0, len(collection))
	for _, doc := range collection {
		docs = append(docs, doc)
	}

	// Sort the slice based on the sortKey
	sort.Slice(docs, func(i, j int) bool {
		valI, iExists := docs[i][sortKey]
		valJ, jExists := docs[j][sortKey]

		// Documents with the key missing will be at the end
		if !iExists {
			return false
		}
		if !jExists {
			return true
		}

		// Type switch to handle different data types
		switch vI := valI.(type) {
		case float64: // JSON numbers are decoded as float64
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
		// Default case if types are different or not sortable
		return false
	})

	return docs
}

// Helper function to check if document matches filter
func matchesFilter(doc Document, filter Document) bool {
	if len(filter) == 0 {
		return true // No filter means match all
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
// It's exported so the access layer can pre-generate IDs before logging.
func GenerateID() string {
	// In a real-world scenario, a more robust UUID library would be better
	// to guarantee uniqueness across machines and time.
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
