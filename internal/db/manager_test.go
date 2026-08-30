package db_test

import (
	"context"
	"errors"
	"testing"

	"app-dev/internal/db"
)

 type FakeClient struct {
	*DynamicFakeQuerier
	pingErr   error
	closeErr  error
	isClosed  bool
	pingCount int
}

func (fc *FakeClient) Ping(ctx context.Context) error {
	fc.pingCount++
	return fc.pingErr
}

func (fc *FakeClient) Close() error {
	fc.isClosed = true
	return fc.closeErr
}

func TestConnectionManager(t *testing.T) {
	ctx := context.Background()

	t.Run("returns ErrNotConnected when no connection is active", func(t *testing.T) {
		manager := db.NewConnectionManager(nil)

		_, err := manager.Inspector()
		if !errors.Is(err, db.ErrNotConnected) {
			t.Errorf("got %v, want %v", err, db.ErrNotConnected)
		}
	})

	t.Run("successfully connects and provides inspector", func(t *testing.T) {
		fakeClient := &FakeClient{
			DynamicFakeQuerier: &DynamicFakeQuerier{},
		}
      
		fakeConnector := func(ctx context.Context, dsn string)(db.Client, error){
			return fakeClient, nil
		}

		manager := db.NewConnectionManager(fakeConnector)

		cfg := db.Config{Host: "localhost", Database: "testdb"}

       
		err := manager.Connect(ctx, cfg)

		inspector, err := manager.Inspector()
		if err != nil {
			t.Fatalf("expected inspector, got error: %v", err)
		}
		if inspector == nil {
			t.Fatal("expected non-nil inspector")
		}
	})

	t.Run("closes previous connection when connecting to a new one", func(t *testing.T) {
		firstClient := &FakeClient{DynamicFakeQuerier: &DynamicFakeQuerier{}}
		secondClient := &FakeClient{DynamicFakeQuerier: &DynamicFakeQuerier{}}

		callCount := 0

		fakeConnector := func(ctx context.Context, dsn string)(db.Client, error){
			 callCount++

			 if callCount == 1 {
				return firstClient, nil
			 }
			 return secondClient, nil
		}

		manager := db.NewConnectionManager(fakeConnector)

		_= manager.Connect(ctx, db.Config{Host: "db1"})
		if firstClient.isClosed{
			t.Error("First client should still be open")
		}

		_= manager.Connect(ctx, db.Config{Host: "db2"})

			if !firstClient.isClosed {
			t.Error("first client should have been closed when second connection was opened")
		}
	})
}