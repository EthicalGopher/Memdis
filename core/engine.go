package core

import (
	"fmt"
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
		// Generate ID if not provided
		id := cmd.ID
		if id == "" {
			id = generateID() // Simple ID generator
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
	count := 0
	collection, exits := e.collections[collectionName]
	if exits {
		return count
	}
	for _, doc := range collection {
		if matchesFilter(doc, filter) {
			count++
		}
	}
	return count
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

// Simple ID generator
func generateID() string {
	// In real implementation, use UUID or better ID generation
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
