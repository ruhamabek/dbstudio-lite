package db_test

import ("testing"
	   "app-dev/internal/db"
)

func TestConfig_ConnectionString(t *testing.T){
   t.Run("generates a valid connection string", func(t *testing.T) {
	    cfg := db.Config{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "secretpassword",
			Database: "mydb",
			SSLMode:  "disable",
		}

		got := cfg.ConnectionString()
		want := "postgres://postgres:secretpassword@localhost:5432/mydb?sslmode=disable"

		if got != want {
			t.Errorf("got %v, want %v", got , want)
		}
   })

   t.Run("defaults to 5432 when port is not specified", func(t *testing.T) {
	        cfg := db.Config{
			Host:     "localhost",
			User:     "postgres",
			Database: "mydb",
			SSLMode:  "disable",
		}

		got := cfg.ConnectionString()
		want := "postgres://postgres@localhost:5432/mydb?sslmode=disable"

		if got != want {
			t.Errorf("got %v, want %v", got , want)
		}
   })
}

func TestPass_url(t *testing.T) {
	t.Run("successfully parses the postgres url", func(t *testing.T) {
		rawURL := "postgres://alice:secret123@db.example.com:5433/production?sslmode=require"
		got, err := db.ParseURL(rawURL)
      
		if err != nil {
			t.Errorf("failed to parse url: %v", err)
		}

       want := db.Config{
			Host:     "db.example.com",
			Port:     5433,
			User:     "alice",
			Password: "secret123",
			Database: "production",
			SSLMode:  "require",
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

	})

	t.Run("returns error on invalid url", func(t *testing.T) {
		  _, err := db.ParseURL("::invalid::")
       
		  if err == nil {
			t.Errorf("expected an error for invalid url got nil")
		  }
	})
}