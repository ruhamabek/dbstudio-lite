package relay_test

import (
	"app-dev/internal/relay"
	"testing"
)

func TestQueryGuard_Validate(t *testing.T) {
	guard := relay.NewQueryGuard()

	t.Run("view permision rejects any SQL execution", func(t *testing.T) {
		      err := guard.Validate(relay.PermissionView, "SELECT * FROM users")
			  if err == nil {
				t.Fatal("expected error for SQL exceution but got nil")
			  }
	})

	t.Run("Read Permission allows SELECT queries", func(t *testing.T) {
		validQueries := []string{
			"SELECT * FROM users",
			"select id, name from products where price > 10",
			"WITH active_users AS (SELECT * FROM users) SELECT * FROM active_users",
			"EXPLAIN SELECT 1",
		}

		for _, q := range validQueries{
			if err := guard.Validate(relay.PermissionRead, q); err != nil{
				t.Errorf("expected no error for query '%s' but got: %v", q, err)
			}
		}
	})

	t.Run("Read permission rejects mutating and DDL queries", func(t *testing.T) {
		blockedQueries := []string{
			"INSERT INTO users (name) VALUES ('alice')",
			"UPDATE users SET name = 'bob'",
			"DELETE FROM users WHERE id = 1",
			"DROP TABLE users",
			"ALTER TABLE users ADD COLUMN age INT",
			"TRUNCATE TABLE users",
		}

		for _,q := range blockedQueries{
			if err := guard.Validate(relay.PermissionRead, q); err == nil{
				t.Errorf("expected error for query '%s' but got nil", q)
			}
		}
	})

	t.Run("update permission allows DML but rejcts destructive DDL", func(t *testing.T) {
		allowed := []string{
			"SELECT * FROM items",
			"INSERT INTO items (name) VALUES ('item1')",
			"UPDATE items SET price = 20 WHERE id = 1",
			"DELETE FROM items WHERE id = 1",
		}

		for _,q := range allowed {
			if err := guard.Validate(relay.PermissionUpdate, q); err != nil {
				t.Errorf("expected %q to be allowed under Update permission, got: %v", q, err)
			}
		}

		destructive := []string{
			"DROP TABLE items",
			"TRUNCATE items",
			"ALTER TABLE items DROP COLUMN price",
		}
		for _, q := range destructive {
			if err := guard.Validate(relay.PermissionUpdate, q); err == nil {
				t.Errorf("expected destructive query %q to be rejected, got nil", q)
			}
		}
	})
}