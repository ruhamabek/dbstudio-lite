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

   t.Run("generates a valid MYSQL connection string", func(t *testing.T) {
	    cfg := db.Config{
			Driver:   db.DriverMySQL,
			Host:     "localhost",
			Port:     3306,
			User:     "root",
			Password: "secretpassword",
			Database: "mydb",
		} 

		got := cfg.ConnectionString()

		want := "root:secretpassword@tcp(localhost:3306)/mydb"

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
   })


	t.Run("generates a valid SQLite connection string", func(t *testing.T) {
		cfg := db.Config{
			Driver:   db.DriverSQLite,
			FilePath: "/home/user/database.db",
		}
		got := cfg.ConnectionString()
		want := "/home/user/database.db"
		if got != want {
			t.Errorf("got %v, want %v", got, want)
     }})
}

func TestPass_url(t *testing.T) {
	t.Run("successfully parses the postgres url", func(t *testing.T) {
		rawURL := "postgres://alice:secret123@db.example.com:5433/production?sslmode=require"
		got, err := db.ParseURL(rawURL)
      
		if err != nil {
			t.Errorf("failed to parse url: %v", err)
		}

       want := db.Config{
		    Driver:   db.DriverPostgres,
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

	t.Run("succesfully parses mysql url", func(t *testing.T) {
		rawURL := "mysql://root:secret123@localhost:3306/ecommerce"
		got, err := db.ParseURL(rawURL)
		if err != nil {
			t.Fatalf("failed to parse mysql url: %v", err)
		}
		want := db.Config{
			Driver:   db.DriverMySQL,
			Host:     "localhost",
			Port:     3306,
			User:     "root",
			Password: "secret123",
			Database: "ecommerce",
		}
		if got != want {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("successfully parses sqlite uri", func(t *testing.T) {
		rawURL := "sqlite:///home/user/app.db"
		got, err := db.ParseURL(rawURL)
		if err != nil {
			t.Fatalf("failed to parse sqlite uri: %v", err)
		}
		want := db.Config{
			Driver:   db.DriverSQLite,
			FilePath: "/home/user/app.db",
		}
		if got != want {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})
	t.Run("successfully parses direct sqlite file path", func(t *testing.T) {
		rawURL := "/home/user/mydata.sqlite3"
		got, err := db.ParseURL(rawURL)
		if err != nil {
			t.Fatalf("failed to parse direct sqlite path: %v", err)
		}
		want := db.Config{
			Driver:   db.DriverSQLite,
			FilePath: "/home/user/mydata.sqlite3",
		}
		if got != want {
			t.Errorf("got %+v, want %+v", got, want)
		} })

	t.Run("returns error on invalid url", func(t *testing.T) {
		  _, err := db.ParseURL("::invalid::")
       
		  if err == nil {
			t.Errorf("expected an error for invalid url got nil")
		  }
	})
}