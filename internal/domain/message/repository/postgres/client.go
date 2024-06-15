package postgres

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"os"
	"time"
)

type Client struct {
	db *pgx.Conn
}

func NewClient(connString string) *Client {
	return &Client{
		db: connectToPostgres(connString),
	}
}

func connectToPostgres(connString string) *pgx.Conn {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	makeMigrations(connString)

	return conn
}

func makeMigrations(connString string) {
	m, err := migrate.New("file://migrations", connString)
	if err != nil {
		slog.Error(err.Error())
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error(err.Error())
	}

	slog.Info("migrations completes successfully")
}

func (c *Client) GetSupportEmployee(username, password string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var count int
	err := c.db.QueryRow(ctx, "SELECT COUNT(*) FROM support_employee WHERE username = $1 AND password = $2;", username, password).Scan(&count)
	if err != nil {
		slog.Error(err.Error())
		return -1, err
	}

	return count, nil
}
