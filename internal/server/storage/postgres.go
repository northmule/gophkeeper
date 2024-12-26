package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// DBQuery интерфейс запросов
type DBQuery interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	PingContext(ctx context.Context) error
	Begin() (*sql.Tx, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
}

// TxDBQuery интерфейс запросов с транзакцими
type TxDBQuery interface {
	QueryRowContext(ctx context.Context, query string, args ...any) (*sql.Row, error)
	Rollback() error
	Commit() error
	Tx() *sql.Tx
	Error() []error
	AddError(e error)
}

// Postgres тип хранилища
type Postgres struct {
	DB    DBQuery
	RawDB *sql.DB
	tx    *sql.Tx
}

// NewPostgres Postgres настройка подключения к БД
func NewPostgres(dsn string) (*Postgres, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	instance := &Postgres{
		DB:    db,
		RawDB: db,
	}

	return instance, nil
}

// Ping доступность БД
func (p *Postgres) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return p.DB.PingContext(ctx)
}

// Transaction транзакции
type Transaction struct {
	t *sql.Tx
	e []error
}

// NewTransaction Открывает новую транзакцию
func NewTransaction(db DBQuery) (*Transaction, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	instance := Transaction{
		t: tx,
		e: make([]error, 0),
	}
	return &instance, nil
}

// QueryRowContext запрос в рамках транзакции
func (t *Transaction) QueryRowContext(ctx context.Context, query string, args ...any) (*sql.Row, error) {
	rows := t.t.QueryRowContext(ctx, query, args...)
	err := rows.Err()
	if err != nil {
		err = errors.Join(err, t.t.Rollback())
		return nil, err
	}

	return rows, nil

}

// Tx транзакция
func (t *Transaction) Tx() *sql.Tx {
	return t.t
}

// Rollback откат
func (t *Transaction) Rollback() error {
	return t.t.Rollback()
}

// Commit сохранить изменения
func (t *Transaction) Commit() error {
	return t.t.Commit()
}

// Error получить ошибки
func (t *Transaction) Error() []error {
	return t.e
}

// AddError добавить ошибку
func (t *Transaction) AddError(e error) {
	t.e = append(t.e, e)
}
