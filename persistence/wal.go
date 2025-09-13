package persistence

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/EthicalGopher/Memdis/core"
	"os"
)

type WAL struct {
	file    *os.File
	encoder *json.Encoder
}

func NewWAL(filename string) (*WAL, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return &WAL{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (w *WAL) Write(cmd core.Command) error {
	return w.encoder.Encode(cmd)
}

func (w *WAL) Restore(engine *core.Engine) error {
	if _, err := w.file.Seek(0, 0); err != nil {
		return err
	}

	scanner := bufio.NewScanner(w.file)
	for scanner.Scan() {
		var cmd core.Command
		if err := json.Unmarshal(scanner.Bytes(), &cmd); err != nil {
			return fmt.Errorf("corrupted log entry: %v", err)
		}
		engine.ApplyCommand(cmd)
	}

	return scanner.Err()
}

func (w *WAL) Close() error {
	return w.file.Close()
}
