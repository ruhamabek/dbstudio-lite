package store

import (
	"app-dev/internal/domain"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileConnectionStore struct {
	mu       sync.RWMutex
	filePath string
}

func NewFileConnectionStore(filePath string) *FileConnectionStore {
	return &FileConnectionStore{
		filePath: filePath,
	}
}

func (s *FileConnectionStore) List() ([]domain.SavedConnection, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.filePath)
	if errors.Is(err, os.ErrNotExist) || len(data) == 0 {
		return []domain.SavedConnection{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read connections file: %w", err)
	}

	var conns []domain.SavedConnection
	if err := json.Unmarshal(data, &conns); err != nil {
		return nil, fmt.Errorf("failed to parse connections file: %w", err)
	}

	return conns, nil
}

func (s *FileConnectionStore) Save(conn domain.SavedConnection) (domain.SavedConnection, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	conns, err := s.readAllUnsafe()
	if err != nil {
		return domain.SavedConnection{}, err
	}

	if conn.ID == "" {
		conn.ID = fmt.Sprintf("conn_%d", time.Now().UnixNano())
		conn.CreatedAt = time.Now()
		conns = append(conns, conn)
	} else {
		found := false
		for i, c := range conns {
			if c.ID == conn.ID {
				conns[i] = conn
				found = true
				break
			}
		}
		if !found {
			conns = append(conns, conn)
		}
	}

	if err := s.writeAllUnsafe(conns); err != nil {
		return domain.SavedConnection{}, err
	}

	return conn, nil
}

func (s *FileConnectionStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	conns, err := s.readAllUnsafe()
	if err != nil {
		return err
	}

	filtered := make([]domain.SavedConnection, 0, len(conns))
	for _, c := range conns {
		if c.ID != id {
			filtered = append(filtered, c)
		}
	}

	return s.writeAllUnsafe(filtered)
}

func (s *FileConnectionStore) readAllUnsafe() ([]domain.SavedConnection, error) {
	data, err := os.ReadFile(s.filePath)
	if errors.Is(err, os.ErrNotExist) || len(data) == 0 {
		return []domain.SavedConnection{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read connections file: %w", err)
	}

	var conns []domain.SavedConnection
	if err := json.Unmarshal(data, &conns); err != nil {
		return nil, fmt.Errorf("failed to parse connections file: %w", err)
	}
	return conns, nil
}

func (s *FileConnectionStore) writeAllUnsafe(conns []domain.SavedConnection) error {
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(conns, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize connections: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write connections file: %w", err)
	}

	return nil
}
