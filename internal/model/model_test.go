package model

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestModel(t *testing.T) {
	// This test is running postgresql with testcontainers
	// If you want to use other instance set dsn below manually
	dsn := startPostgresql(t)

	assert.NoError(t, InitDBPool(dsn))

	repo, e := NewDBRepo(context.Background())
	assert.NoError(t, e)
	assert.NoError(t, repo.Init())
	repo.Close()

	getTr := func(id string) *Transaction {
		return &Transaction{
			TransactionId:   id,
			TransactionType: Sale,
			OrderId:         "c66oxMaisTwJQXjD",
			Amount:          10,
			Currency:        Eur,
			Description:     "Test transaction",
			AccountId:       "001",
		}
	}

	t.Run("Insert", func(t *testing.T) {
		repo, e = NewDBRepo(context.Background())
		assert.NoError(t, e)
		defer repo.Close()

		p := getTr("tqZi6QapS41zcEHy")
		assert.NoError(t, repo.Insert(p))
		assert.NoError(t, repo.Empty())
	})

	t.Run("Find", func(t *testing.T) {
		repo, e = NewDBRepo(context.Background())
		assert.NoError(t, e)
		defer repo.Close()

		p := getTr("tqZi6QapS41zcEHy")
		assert.NoError(t, repo.Insert(p))
		r, e := repo.Find(p.TransactionId)
		assert.NoError(t, e)
		expTr := getTr("tqZi6QapS41zcEHy")
		expTr.Id = 1
		assert.Equal(t, expTr, r)
		assert.NoError(t, repo.Empty())
	})

	t.Run("FetchAll", func(t *testing.T) {
		repo, e = NewDBRepo(context.Background())
		assert.NoError(t, e)
		defer repo.Close()

		p := getTr("tqZi6QapS41zcEHy1")
		p2 := getTr("tqZi6QapS41zcEHy2")
		assert.NoError(t, repo.Insert(p))
		assert.NoError(t, repo.Insert(p2))
		r, err := repo.FetchAll()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(r))
		assert.NoError(t, repo.Empty())
	})

	t.Run("MultipleInsert - the same transaction_id", func(t *testing.T) {
		wg := sync.WaitGroup{}
		runs := 10
		wg.Add(runs)
		ctx := context.Background()
		for i := 0; i < runs; i++ {
			go func(i int) {
				// Must have only one transaction with this transaction_id
				p := getTr("tqZi6QapS41zcEHy")
				repo, e = NewDBRepo(ctx)
				defer repo.Close()
				repo.Insert(p)
				wg.Done()
			}(i)
		}
		wg.Wait()
		repo, e = NewDBRepo(context.Background())
		defer repo.Close()
		r, err := repo.FetchAll()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(r))
	})
}

const (
	dbUsername string = "user"
	dbPassword string = "password"
	dbName     string = "test"
)

func startPostgresql(t *testing.T) string {
	ctx := context.Background()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:14-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     dbUsername,
				"POSTGRES_PASSWORD": dbPassword,
				"POSTGRES_DB":       dbName,
			},
			WaitingFor: wait.ForListeningPort("5432/tcp"),
		},
		Started: true,
	})

	if err != nil {
		t.Error(err)
	}

	t.Cleanup(func() {
		container.Terminate(ctx) // nolint: errcheck
	})

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432/tcp")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUsername, dbPassword, host, port.Port(), dbName)
}
