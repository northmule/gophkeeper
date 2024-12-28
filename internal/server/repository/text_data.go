package repository

import (
	"context"
	"database/sql"

	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/storage"
)

// TextDataRepository репозитарий текстовых данных
type TextDataRepository struct {
	store storage.DBQuery

	sqlFindOneByUUID *sql.Stmt
}

// NewTextDataRepository конструктор
func NewTextDataRepository(store storage.DBQuery) (*TextDataRepository, error) {
	var err error
	instance := new(TextDataRepository)
	instance.store = store
	instance.sqlFindOneByUUID, err = store.Prepare(`select id, name, value, uuid from text_data where uuid = $1 limit 1`)
	if err != nil {
		return nil, ErrorMsg(err)
	}
	return instance, nil
}

// FindOneByUUID поиск значения по UUID
func (r *TextDataRepository) FindOneByUUID(ctx context.Context, uuid string) (*models.TextData, error) {

	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	rows, err := r.sqlFindOneByUUID.QueryContext(ctx, uuid)
	if err != nil {
		return nil, ErrorMsg(err)
	}
	err = rows.Err()
	if err != nil {
		return nil, ErrorMsg(err)
	}
	data := new(models.TextData)
	if rows.Next() {
		err = rows.Scan(&data.ID, &data.Name, &data.Value, &data.UUID)
		if err != nil {
			return nil, ErrorMsg(err)
		}
	}

	return data, nil
}

// Add Новое значение
func (r *TextDataRepository) Add(ctx context.Context, data *models.TextData) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	rows := r.store.QueryRowContext(ctx, `insert into text_data (name, "value", uuid) values ($1, $2, $3) returning id`, data.Name, data.Value, data.UUID)
	err := rows.Err()
	if err != nil {
		return 0, ErrorMsg(err)
	}

	var id int64
	err = rows.Scan(&id)
	if err != nil {
		return 0, ErrorMsg(err)
	}
	return id, nil
}

// Update Обновление основных полей
func (r *TextDataRepository) Update(ctx context.Context, data *models.TextData) error {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	rows := r.store.QueryRowContext(ctx, `update text_data set name = $1, value = $2 where uuid = $3`, data.Name, data.Value, data.UUID)

	return rows.Err()
}
