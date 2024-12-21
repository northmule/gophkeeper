package db

import (
	"context"
	"database/sql"
	"embed"
	"time"

	"github.com/pressly/goose/v3"
)

type Migrations struct {
	mFS   embed.FS
	sqlDB *sql.DB
}

//go:embed migrations/*.sql
var migrationsFS embed.FS

func NewMigrations(db *sql.DB) *Migrations {
	instance := Migrations{}
	instance.mFS = migrationsFS
	instance.sqlDB = db
	return &instance
}

func (m *Migrations) Up(ctx context.Context) error {
	goose.SetBaseFS(m.mFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	ctxMigrations, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := goose.UpContext(ctxMigrations, m.sqlDB, "migrations"); err != nil {
		return err
	}

	return nil
}
