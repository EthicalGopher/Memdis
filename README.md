# Memdis

Memdis is a simple, in-memory document database built with Go. It provides basic CRUD operations (Create, Read, Update, Delete) on JSON-like documents, along with features like counting, sorting, and snapshotting for persistence.

## Features

- **In-Memory:** Fast data access and manipulation.
- **Document-Oriented:** Stores JSON-like documents in collections.
- **Write-Ahead Log (WAL):** Ensures data durability and recovery.
- **Snapshotting:** Periodically saves the database state to disk for faster recovery.
- **CLI Interface:** Interact with the database using a command-line interface.

## Installation

To get started with Memdis, follow these steps:

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/EthicalGopher/Memdis.git
    cd Memdis
    ```

2.  **Build the project:**

    ```bash
    go build
    ```

    This will create an executable named `Memdis` in the project root directory.

## Usage

Memdis provides a command-line interface (CLI) to interact with the database. Below are the available commands and their usage.

### General CLI Structure

All commands follow the structure: `./Memdis [command] [arguments...]`

### Commands

#### `insert`

Inserts a new document into a specified collection.

-   **Usage:** `./Memdis insert [collection] '[json_data]'
-   **Arguments:**
    -   `collection`: The name of the collection to insert the document into.
    -   `json_data`: The document data in JSON format. **Must be enclosed in single quotes.**
-   **Example:**

    ```bash
    ./Memdis insert users '{"name":"Alice", "age":30}'
    ```

#### `find`

Finds documents in a specified collection, optionally filtered by a JSON query.

-   **Usage:** `./Memdis find [collection] ['[filter_json]']'
-   **Arguments:**
    -   `collection`: The name of the collection to search.
    -   `filter_json` (optional): A JSON object specifying the filter criteria. **Must be enclosed in single quotes.**
-   **Examples:**

    ```bash
    ./Memdis find users
    ./Memdis find users '{"age":30}'
    ```

#### `update`

Updates documents in a specified collection that match a filter with new data.

-   **Usage:** `./Memdis update [collection] '[filter_json]' '[update_json]'
-   **Arguments:**
    -   `collection`: The name of the collection to update.
    -   `filter_json`: A JSON object specifying the filter criteria for documents to update. **Must be enclosed in single quotes.**
    -   `update_json`: A JSON object containing the new data to apply to the matching documents. **Must be enclosed in single quotes.**
-   **Example:**

    ```bash
    ./Memdis update users '{"name":"Alice"}' '{"age":31}'
    ```

#### `delete`

Deletes documents from a specified collection that match a filter.

-   **Usage:** `./Memdis delete [collection] '[filter_json]'
-   **Arguments:**
    -   `collection`: The name of the collection to delete from.
    -   `filter_json`: A JSON object specifying the filter criteria for documents to delete. **Must be enclosed in single quotes.**
-   **Example:**

    ```bash
    ./Memdis delete users '{"age":31}'
    ```

#### `count`

Counts the number of documents in a specified collection, optionally filtered by a JSON query.

-   **Usage:** `./Memdis count [collection] ['[filter_json]']'
-   **Arguments:**
    -   `collection`: The name of the collection to count.
    -   `filter_json` (optional): A JSON object specifying the filter criteria. **Must be enclosed in single quotes.**
-   **Examples:**

    ```bash
    ./Memdis count users
    ./Memdis count users '{"name":"Alice"}'
    ```

#### `sort`

Sorts documents in a specified collection by a given key.

-   **Usage:** `./Memdis sort [collection] [sort_key]'
-   **Arguments:**
    -   `collection`: The name of the collection to sort.
    -   `sort_key`: The key by which to sort the documents.
-   **Example:**

    ```bash
    ./Memdis sort users age
    ```

#### `save`

Saves the current state of the database to a snapshot file.

-   **Usage:** `./Memdis save`
-   **Example:**

    ```bash
    ./Memdis save
    ```

#### `list-collections`

Lists all available collections in the database.

-   **Usage:** `./Memdis list-collections`
-   **Example:**

    ```bash
    ./Memdis list-collections
    ```

## Using Memdis as a Go Package

You can integrate Memdis directly into your Go applications as a library. This allows you to programmatically interact with the database without using the CLI.

### 1. Import the Package

First, ensure you have the Memdis package imported in your Go project:

```go
import (
    "fmt"
    "log"
    "github.com/EthicalGopher/Memdis/Mem" // Adjust import path if necessary
    "github.com/EthicalGopher/Memdis/core"
)
```

### 2. Connect to the Database

Initialize a database connection using `Mem.Connect`. You need to provide a file path for the Write-Ahead Log (WAL).

```go
db, err := Mem.Connect("my_database.mem")
if err != nil {
    log.Fatalf("Failed to connect to database: %v", err)
}
deffer db.Close() // Ensure the database connection is closed when your application exits
```

### 3. Execute Commands

You can execute any database command using the `db.Execute()` method. This method takes a command string (similar to the CLI commands) and returns the result and an error, if any.

#### Example: Inserting a Document

```go
insertCmd := `INSERT users {"name":"Bob", "email":"bob@example.com"}`
result, err := db.Execute(insertCmd)
if err != nil {
    fmt.Printf("Error inserting document: %v\n", err)
} else {
    fmt.Printf("Insert result: %v\n", result)
}
```

#### Example: Finding Documents

```go
findCmd := `FIND users {"name":"Bob"}`
result, err = db.Execute(findCmd)
if err != nil {
    fmt.Printf("Error finding documents: %v\n", err)
} else {
    // The result for FIND operations is typically a slice of core.Document
    if docs, ok := result.([]core.Document); ok {
        for _, doc := range docs {
            fmt.Printf("Found document: %+v\n", doc)
        }
    } else {
        fmt.Printf("Find result: %v\n", result)
    }
}
```

### Full Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/EthicalGopher/Memdis/Mem"
	"github.com/EthicalGopher/Memdis/core"
)

func main() {
	// 1. Connect to the database
	db, err := Mem.Connect("my_application_data.mem")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close() // Ensure the database connection is closed

	// 2. Insert a document
	insertCmd := `INSERT products {"name":"Laptop", "price":1200.00, "inStock":true}`
	insertResult, err := db.Execute(insertCmd)
	if err != nil {
		fmt.Printf("Error inserting document: %v\n", err)
	} else {
		fmt.Printf("Insert result: %v\n", insertResult)
	}

	// 3. Find documents
	findCmd := `FIND products {"inStock":true}`
	findResult, err := db.Execute(findCmd)
	if err != nil {
		fmt.Printf("Error finding documents: %v\n", err)
	} else {
		if docs, ok := findResult.([]core.Document); ok {
			fmt.Println("\nFound products:")
			for _, doc := range docs {
				fmt.Printf("- %+v\n", doc)
			}
		} else {
			fmt.Printf("Find result: %v\n", findResult)
		}
	}

	// 4. Update a document
	updateCmd := `UPDATE products {"name":"Laptop"} {"price":1150.00}`
	updateResult, err := db.Execute(updateCmd)
	if err != nil {
		fmt.Printf("Error updating document: %v\n", err)
	} else {
		fmt.Printf("\nUpdate result: %v\n", updateResult)
	}

	// 5. Count documents
	countCmd := `COUNT products {"inStock":true}`
	countResult, err := db.Execute(countCmd)
	if err != nil {
		fmt.Printf("Error counting documents: %v\n", err)
	} else {
		fmt.Printf("\nCount of in-stock products: %v\n", countResult)
	}

	// 6. Save snapshot
	saveResult, err := db.Execute("SAVE")
	if err != nil {
		fmt.Printf("Error saving snapshot: %v\n", err)
	} else {
		fmt.Printf("\nSave result: %v\n", saveResult)
	}
}
