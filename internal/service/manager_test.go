package service_test

import (
	"app-dev/internal/domain"
	"app-dev/internal/service"
	"context"
	"errors"
	"testing"
)

type DynamicFakeQuerier struct {
	recordedSQL  string
	recordedArgs []any
}

func (q *DynamicFakeQuerier) Query(ctx context.Context, sql string, args ...any) (domain.Rows, error) {
	q.recordedSQL = sql
	q.recordedArgs = args
	return nil, nil
}

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
		manager := service.NewConnectionManager(nil)

		_, err := manager.Inspector()
		if !errors.Is(err, service.ErrNotConnected) {
			t.Errorf("got %v, want %v", err, service.ErrNotConnected)
		}
	})

	t.Run("successfully connects and provides inspector", func(t *testing.T) {
		fakeClient := &FakeClient{
			DynamicFakeQuerier: &DynamicFakeQuerier{},
		}

		fakeConnector := func(ctx context.Context, cfg domain.Config) (domain.Client, error) {
			return fakeClient, nil
		}

		manager := service.NewConnectionManager(fakeConnector)
		cfg := domain.Config{Host: "localhost", Database: "testdb"}

		err := manager.Connect(ctx, cfg)
		if err != nil {
			t.Fatalf("unexpected connect error: %v", err)
		}

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
		fakeConnector := func(ctx context.Context, cfg domain.Config) (domain.Client, error) {
			callCount++
			if callCount == 1 {
				return firstClient, nil
			}
			return secondClient, nil
		}

		manager := service.NewConnectionManager(fakeConnector)

		_ = manager.Connect(ctx, domain.Config{Host: "db1"})
		if firstClient.isClosed {
			t.Error("first client should still be open")
		}

		_ = manager.Connect(ctx, domain.Config{Host: "db2"})
		if !firstClient.isClosed {
			t.Error("first client should have been closed when second connection was opened")
		}
	})
}
