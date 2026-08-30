package db_test


import (
	"path/filepath"
	"testing"
    "app-dev/internal/db"

 )


func TestFileConnectionStore(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "connections.json")
	    store := db.NewFileConnectionStore(tempFile)

		t.Run("return an empty list when file does not exist yet", func(t *testing.T) {
			conns, err := store.List()
            
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(conns) != 0 {
				t.Fatalf("expected 0 conns but got %v", len(conns))
			}
			
		})

		t.Run("saved a new connection and generates and ID", func(t *testing.T) {
			saved, err := store.Save(db.SavedConnection{
					Name: "Local Postgres",
			URL:  "postgres://postgres:pass@localhost:5432/mydb",
				})
				if err != nil {
					t.Fatalf("failed to save connection: %v", err)
				}
				if saved.ID == "" {
					t.Errorf("expected generated ID, got empty string")
				}
				
				conns, err := store.List()
				if err != nil {
					t.Fatalf("failed to list connections: %v", err)
				}
				
				if len(conns) != 1 {
					  t.Fatalf("got %v, want 1", len(conns))
				}
				if conns[0].Name != "Local Postgres" {
				t.Errorf("got name %q, want 'Local Postgres'", conns[0].Name)
		}
	})

    t.Run("deletes a connection by ID", func(t *testing.T) {
		conns, _ := store.List()
		if len(conns) == 0 {
			t.Fatal("expected at least 1 connection to delete")
		}
		idToDelete := conns[0].ID
		err := store.Delete(idToDelete)
		if err != nil {
			t.Fatalf("failed to delete connection: %v", err)
		}
		afterDelete, err := store.List()
		if err != nil {
			t.Fatalf("failed to list connections after delete: %v", err)
		}
		if len(afterDelete) != 0 {
			t.Errorf("expected 0 connections after deletion, got %d", len(afterDelete))
		}
	})

}