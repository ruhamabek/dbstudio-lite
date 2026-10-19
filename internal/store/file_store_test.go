package store_test

import (
	"app-dev/internal/domain"
	"app-dev/internal/store"
	"os"
	"path/filepath"
	"testing"
)

func TestFileConnectionStore(t *testing.T) {
	tempDir := t.TempDir()
	storePath := filepath.Join(tempDir, "connections.json")

	t.Run("return an empty list when file does not exist yet", func(t *testing.T) {
		s := store.NewFileConnectionStore(storePath)
		conns, err := s.List()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(conns) != 0 {
			t.Errorf("got %d connections, want 0", len(conns))
		}
	})

	t.Run("saved a new connection and generates and ID", func(t *testing.T) {
		s := store.NewFileConnectionStore(storePath)
		conn := domain.SavedConnection{
			Name: "Production DB",
			URL:  "postgres://localhost:5432/prod",
		}

		saved, err := s.Save(conn)
		if err != nil {
			t.Fatalf("failed to save connection: %v", err)
		}

		if saved.ID == "" {
			t.Errorf("expected generated ID, got empty string")
		}

		if saved.Name != conn.Name || saved.URL != conn.URL {
			t.Errorf("saved connection data mismatch: got %+v, want %+v", saved, conn)
		}

		list, err := s.List()
		if err != nil {
			t.Fatalf("failed to list connections: %v", err)
		}

		if len(list) != 1 {
			t.Fatalf("got %d connections, want 1", len(list))
		}
	})

	t.Run("deletes a connection by ID", func(t *testing.T) {
		s := store.NewFileConnectionStore(storePath)
		list, err := s.List()
		if err != nil {
			t.Fatalf("failed to list: %v", err)
		}
		if len(list) == 0 {
			t.Fatal("expected at least 1 connection to delete")
		}

		targetID := list[0].ID
		err = s.Delete(targetID)
		if err != nil {
			t.Fatalf("failed to delete connection: %v", err)
		}

		afterDelete, err := s.List()
		if err != nil {
			t.Fatalf("failed to list after delete: %v", err)
		}

		for _, c := range afterDelete {
			if c.ID == targetID {
				t.Errorf("connection with ID %q was not deleted", targetID)
			}
		}
	})

	t.Run("creates file with 0600 permissions", func(t *testing.T) {
		s := store.NewFileConnectionStore(storePath)
		_, err := s.Save(domain.SavedConnection{Name: "Perm Test", URL: "postgres://localhost/test"})
		if err != nil {
			t.Fatalf("failed to save: %v", err)
		}

		info, err := os.Stat(storePath)
		if err != nil {
			t.Fatalf("failed to stat file: %v", err)
		}

		perm := info.Mode().Perm()
		if perm != 0600 {
			t.Errorf("got file permissions %o, want 0600", perm)
		}
	})
}
