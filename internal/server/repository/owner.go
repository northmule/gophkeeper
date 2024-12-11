package repository

import (
	"context"
	"database/sql"

	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/storage"
)

type OwnerRepository struct {
	store                                      storage.DBQuery
	sqlFindByDataUUID                          *sql.Stmt
	sqlFindByUserUUID                          *sql.Stmt
	sqlFindOneByUserUUIDAndDataUUIDAndDataType *sql.Stmt
}

func NewOwnerRepository(store storage.DBQuery) (*OwnerRepository, error) {
	var err error
	instance := new(OwnerRepository)
	instance.store = store

	instance.sqlFindOneByUserUUIDAndDataUUIDAndDataType, err = store.Prepare(`select id, user_uuid, data_type, data_uuid from owner where user_uuid = $1 and data_uuid = $2 and data_type = $3 limit 1`)
	if err != nil {
		return nil, ErrorMsg(err)
	}

	return instance, nil
}

func (r *OwnerRepository) FindOneByUserUUIDAndDataUUIDAndDataType(ctx context.Context, userUuid string, dataUuid string, dataType string) (*models.Owner, error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	rows, err := r.sqlFindOneByUserUUIDAndDataUUIDAndDataType.QueryContext(ctx, userUuid, dataUuid, dataType)
	if err != nil {
		return nil, ErrorMsg(err)
	}
	err = rows.Err()
	if err != nil {
		return nil, ErrorMsg(err)
	}
	data := new(models.Owner)
	if rows.Next() {
		err = rows.Scan(&data.ID, &data.UserUUID, &data.DataType, &data.DataUUID)
		if err != nil {
			return nil, ErrorMsg(err)
		}
	}

	return data, nil
}

// Add Новое значение
func (r *OwnerRepository) Add(ctx context.Context, data *models.Owner) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	rows := r.store.QueryRowContext(ctx, `insert into owner (user_uuid, data_type, data_uuid) values ($1, $2, $3) returning id`, data.UserUUID, data.DataType, data.DataUUID)
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
