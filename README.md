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