package dbalancer_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/sync/errgroup"

	"github.com/rubenschmidt/dbalancer"
	"github.com/stretchr/testify/require"
)

func BenchmarkBalancer(b *testing.B) {
	ctx := context.TODO()
	db, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	require.NoError(b, err)
	setupDB(db)
	rep, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable")
	require.NoError(b, err)
	setupDB(rep)

	rep2, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5434/postgres?sslmode=disable")
	require.NoError(b, err)
	setupDB(rep2)

	// Create a new DBalancer with the master DB
	bl := dbalancer.NewDBalancer(db, rep, rep2)
	defer bl.Close()

	bl.SetMaxOpenConns(100)
	bl.SetMaxIdleConns(50)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		g := errgroup.Group{}
		for i := 0; i < 1000; i++ {
			g.Go(func() error {
				c, err := bl.ReadConn(ctx)
				if err != nil {
					panic(err)
				}
				r, err := c.QueryContext(ctx, "SELECT * FROM users where id>30")
				if err != nil {
					panic(err)
				}
				r.Close()
				c.Close()
				return nil
			})
		}
		err := g.Wait()
		if err != nil {
			panic(err)
		}
	}
}

func setupDB(db *sql.DB) {
	db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name text NOT NULL)")
	db.Exec("delete from users")
	for i := 0; i < 100; i++ {
		db.Exec("INSERT INTO users (name) VALUES ($1)", "user"+fmt.Sprintf("%d", i))
	}
}
