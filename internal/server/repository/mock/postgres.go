package mock

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

// MockDBQuery мок
type MockDBQuery struct {
	mock.Mock
}

// ExecContext мок
func (m *MockDBQuery) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	ret := m.Called(ctx, query, args)

	var r0 sql.Result
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(sql.Result)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

// QueryContext мок
func (m *MockDBQuery) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	ret := m.Called(ctx, query, args)

	var r0 *sql.Rows
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*sql.Rows)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

// PingContext мок
func (m *MockDBQuery) PingContext(ctx context.Context) error {
	ret := m.Called(ctx)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// Begin мок
func (m *MockDBQuery) Begin() (*sql.Tx, error) {
	ret := m.Called()

	var r0 *sql.Tx
	if ret.Get(0) != nil {
		//r0 = ret.Get(0).(*sql.Tx)
		r0 = nil // услоно создали транзакцию
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

// QueryRowContext мок
func (m *MockDBQuery) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	ret := m.Called(ctx, query, args)

	var r0 *sql.Row
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*sql.Row)
	}

	return r0
}

// Prepare мок
func (m *MockDBQuery) Prepare(query string) (*sql.Stmt, error) {
	ret := m.Called(query)

	var r0 *sql.Stmt
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*sql.Stmt)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

// MockTxDBQuery мок
type MockTxDBQuery struct {
	mock.Mock
}

// QueryRowContext мок
func (m *MockTxDBQuery) QueryRowContext(ctx context.Context, query string, args ...any) (*sql.Row, error) {
	ret := m.Called(ctx, query, args)

	var r0 *sql.Row
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*sql.Row)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

// Rollback мок
func (m *MockTxDBQuery) Rollback() error {
	ret := m.Called()

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// Commit мок
func (m *MockTxDBQuery) Commit() error {
	ret := m.Called()

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// Tx мок
func (m *MockTxDBQuery) Tx() *sql.Tx {
	ret := m.Called()

	var r0 *sql.Tx
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*sql.Tx)
	}

	return r0
}

// Error мок
func (m *MockTxDBQuery) Error() []error {
	ret := m.Called()

	var r0 []error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]error)
	}

	return r0
}

// AddError мок
func (m *MockTxDBQuery) AddError(e error) {
	m.Called(e)
}
