package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/storage"
)

type CardDataRepository struct {
	store storage.DBQuery

	sqlFindOneByUUID *sql.Stmt
}

func NewCardDataRepository(store storage.DBQuery) (*CardDataRepository, error) {
	var err error
	instance := new(CardDataRepository)
	instance.store = store
	instance.sqlFindOneByUUID, err = store.Prepare(`select id, value, object_type, name, uuid from card_data where uuid = $1 limit 1`)
	if err != nil {
		return nil, ErrorMsg(err)
	}
	return instance, nil
}

// FindOneByUUID поиск значения по UUID
func (r *CardDataRepository) FindOneByUUID(ctx context.Context, uuid string) (*models.CardData, error) {

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
	data := new(models.CardData)
	if rows.Next() {
		data.Value = models.CardDataValueV1{}
		var jsonbValue string
		err = rows.Scan(&data.ID, &jsonbValue, &data.ObjectType, &data.Name, &data.UUID)
		if err != nil {
			return nil, ErrorMsg(err)
		}
		err = json.Unmarshal([]byte(jsonbValue), &data.Value)
		if err != nil {
			return nil, ErrorMsg(err)
		}
	}

	return data, nil
}

// Add Новое значение
func (r *CardDataRepository) Add(ctx context.Context, data *models.CardData) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	rows := r.store.QueryRowContext(ctx, `insert into card_data (name, object_type, "value", uuid) values ($1, $2, $3, $4) returning id`, data.Name, data.ObjectType, data.Value, data.UUID)
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
func (r *CardDataRepository) Update(ctx context.Context, data *models.CardData) error {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	rows := r.store.QueryRowContext(ctx, `update card_data set name = $1, value = $2 where uuid = $3`, data.Name, data.Value, data.UUID)

	return rows.Err()
}
