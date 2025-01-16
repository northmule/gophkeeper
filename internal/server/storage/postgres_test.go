package storage

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

type MockDBQuery struct {
	mock.Mock
}

func (m *MockDBQuery) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	args = append([]interface{}{ctx, query}, args...)
	ret := m.Called(args...)
	return ret.Get(0).(sql.Result), ret.Error(1)
}

func (m *MockDBQuery) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	args = append([]interface{}{ctx, query}, args...)
	ret := m.Called(args...)
	return ret.Get(0).(*sql.Rows), ret.Error(1)
}

func (m *MockDBQuery) PingContext(ctx context.Context) error {
	ret := m.Called(ctx)
	return ret.Error(0)
}

func (m *MockDBQuery) Begin() (*sql.Tx, error) {
	ret := m.Called()
	return ret.Get(0).(*sql.Tx), ret.Error(1)
}

func (m *MockDBQuery) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	args = append([]interface{}{ctx, query}, args...)
	ret := m.Called(args...)
	return ret.Get(0).(*sql.Row)
}

func (m *MockDBQuery) Prepare(query string) (*sql.Stmt, error) {
	ret := m.Called(query)
	return ret.Get(0).(*sql.Stmt), ret.Error(1)
}

// MockTxDBQuery is a mock implementation of TxDBQuery interface
type MockTxDBQuery struct {
	mock.Mock
}

func (m *MockTxDBQuery) QueryRowContext(ctx context.Context, query string, args ...any) (*sql.Row, error) {
	args = append([]interface{}{ctx, query}, args...)
	ret := m.Called(args...)
	return ret.Get(0).(*sql.Row), ret.Error(1)
}

func (m *MockTxDBQuery) Rollback() error {
	ret := m.Called()
	return ret.Error(0)
}

func (m *MockTxDBQuery) Commit() error {
	ret := m.Called()
	return ret.Error(0)
}

func (m *MockTxDBQuery) Tx() *sql.Tx {
	ret := m.Called()
	return ret.Get(0).(*sql.Tx)
}

func (m *MockTxDBQuery) Error() []error {
	ret := m.Called()
	return ret.Get(0).([]error)
}

func (m *MockTxDBQuery) AddError(e error) {
	m.Called(e)
}

func TestNewPostgres(t *testing.T) {
	p, err := NewPostgres("dsn")
	assert.NoError(t, err)
	assert.NotNil(t, p)
}

func TestTransaction_Error(t *testing.T) {
	tx := &Transaction{e: []error{errors.New("error 1"), errors.New("error 2")}}

	errs := tx.Error()
	assert.Len(t, errs, 2)
	assert.Equal(t, "error 1", errs[0].Error())
	assert.Equal(t, "error 2", errs[1].Error())
}

func TestTransaction_AddError(t *testing.T) {
	tx := &Transaction{e: []error{}}

	tx.AddError(errors.New("error 1"))
	assert.Len(t, tx.e, 1)
	assert.Equal(t, "error 1", tx.e[0].Error())

	tx.AddError(errors.New("error 2"))
	assert.Len(t, tx.e, 2)
	assert.Equal(t, "error 2", tx.e[1].Error())
}

func TestTransaction_Tx(t *testing.T) {
	mockTx := new(sql.Tx)
	tx := &Transaction{t: mockTx}

	actualTx := tx.Tx()
	assert.Equal(t, mockTx, actualTx)
}
