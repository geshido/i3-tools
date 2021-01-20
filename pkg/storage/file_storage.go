package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type FileStorage struct {
	f        *os.File
	filename string
	mux      sync.RWMutex
}

func NewFileStorage(filename string) (*FileStorage, error) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("cannot open file for writing: %w", err)
	}

	return &FileStorage{f: f, filename: filename}, nil
}

func (f *FileStorage) Store(data interface{}) error {
	f.mux.Lock()
	defer f.mux.Unlock()
	if err := f.f.Truncate(0); err != nil {
		return fmt.Errorf("cannot truncate file: %w", err)
	}
	if _, err := f.f.Seek(0, 0); err != nil {
		return fmt.Errorf("cannot seek to beginning of the file: %w", err)
	}

	if err := json.NewEncoder(f.f).Encode(data); err != nil {
		return fmt.Errorf("cannot encode data: %w", err)
	}

	return nil
}

func (f *FileStorage) Read(data interface{}) error {
	f.mux.RLock()
	defer f.mux.RUnlock()

	if err := json.NewDecoder(f.f).Decode(&data); err != nil {
		return fmt.Errorf("cannot decode: %w", err)
	}
	return nil
}

func (f *FileStorage) Cleanup() error {
	f.mux.Lock()
	defer f.mux.Unlock()
	return os.Remove(f.filename)
}
