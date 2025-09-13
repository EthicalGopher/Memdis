package persistence

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/EthicalGopher/Memdis/core"
)

// WAL handles both the Write-Ahead Log and snapshotting.
type WAL struct {
	file         *os.File
	snapshotPath string
}

// NewWAL creates a new WAL and determines the path for its snapshot file.
func NewWAL(walPath string) (*WAL, error) {
	file, err := os.OpenFile(walPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	// Derive snapshot path from WAL path (e.g., data.mem -> data.snapshot)
	snapshotPath := strings.TrimSuffix(walPath, ".mem") + ".snapshot"

	return &WAL{
		file:         file,
		snapshotPath: snapshotPath,
	}, nil
}

// Write appends a command to the WAL file.
func (w *WAL) Write(cmd core.Command) error {
	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	_, err = w.file.Write(append(data, '\n'))
	return err
}

// Close closes the WAL file.
func (w *WAL) Close() error {
	return w.file.Close()
}

// Truncate clears the WAL file. This is called after a successful snapshot.
func (w *WAL) Truncate() error {
	// To truncate, we close the current file handle,
	// and re-open the same file with the Truncate flag, which clears it.
	if err := w.file.Close(); err != nil {
		return err
	}
	file, err := os.OpenFile(w.file.Name(), os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	w.file = file
	return nil
}

// SaveSnapshot serializes the engine state and writes it to the snapshot file.
func (w *WAL) SaveSnapshot(engine *core.Engine) error {
	data, err := engine.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize engine state: %w", err)
	}

	// Write to a temporary file first to prevent corruption if the app crashes.
	tempPath := w.snapshotPath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary snapshot: %w", err)
	}

	// Atomically rename the temporary file to the final snapshot file.
	return os.Rename(tempPath, w.snapshotPath)
}

// Restore loads the database state, preferring the snapshot and then replaying the WAL.
func (w *WAL) Restore(engine *core.Engine) error {
	// 1. Attempt to load from snapshot.
	snapshotData, err := os.ReadFile(w.snapshotPath)
	if err == nil { // Snapshot exists and is readable.
		if err := engine.Deserialize(snapshotData); err != nil {
			log.Printf("⚠️ Warning: could not deserialize snapshot, it may be corrupt: %v. Attempting full WAL restore.", err)
		} else {
			log.Println("✅ State restored from snapshot.")
		}
	} else if !os.IsNotExist(err) {
		// An error other than "file not found" occurred.
		log.Printf("⚠️ Warning: could not read snapshot file: %v. Attempting full WAL restore.", err)
	}

	// 2. Replay any commands in the WAL that occurred after the snapshot.
	if _, err := w.file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek WAL file for restore: %w", err)
	}

	scanner := bufio.NewScanner(w.file)
	lines := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var cmd core.Command
		if err := json.Unmarshal(line, &cmd); err != nil {
			log.Printf("⚠️ Warning: skipping corrupt line in WAL: %v", err)
			continue
		}
		engine.ApplyCommand(cmd)
		lines++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading WAL file: %w", err)
	}

	if lines > 0 {
		log.Printf("✅ Replayed %d commands from WAL.", lines)
	}

	return nil
}
