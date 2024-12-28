package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/storage"
)

type MetaDataRepository struct {
	store storage.DBQuery

	sqlAllFindByDataUUID *sql.Stmt
}

func NewMetaDataRepository(store storage.DBQuery) (*MetaDataRepository, error) {
	var err error
	instance := new(MetaDataRepository)
	instance.store = store
	instance.sqlAllFindByDataUUID, err = store.Prepare(`select id, meta_name, meta_value, data_uuid from meta_data where data_uuid = $1`)
	if err != nil {
		return nil, ErrorMsg(err)
	}
	return instance, nil
}

// FindOneByUUID поиск значения по UUID
func (r *MetaDataRepository) FindOneByUUID(ctx context.Context, uuid string) ([]models.MetaData, error) {

	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	rows, err := r.sqlAllFindByDataUUID.QueryContext(ctx, uuid)
	if err != nil {
		return nil, ErrorMsg(err)
	}
	err = rows.Err()
	if err != nil {
		return nil, ErrorMsg(err)
	}

	var metaDataList []models.MetaData
	for rows.Next() {
		data := models.MetaData{}
		var jsonbValue string
		err = rows.Scan(&data.ID, &data.MetaName, &jsonbValue, &data.DataUUID)
		if err != nil {
			return nil, ErrorMsg(err)
		}
		err = json.Unmarshal([]byte(jsonbValue), &data.MetaValue)
		if err != nil {
			return nil, ErrorMsg(err)
		}
		metaDataList = append(metaDataList, data)
	}

	return metaDataList, nil
}

// Add Новое значение
func (r *MetaDataRepository) Add(ctx context.Context, data *models.MetaData) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	rows := r.store.QueryRowContext(ctx, `insert into meta_data (meta_name, meta_value, data_uuid) values ($1, $2, $3) returning id`, data.MetaName, data.MetaValue, data.DataUUID)
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

// ReplaceMetaByDataUUID удалит и затем добавит новые
func (r *MetaDataRepository) ReplaceMetaByDataUUID(ctx context.Context, dataUUID string, metaDataList []models.MetaData) error {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	var err error
	var tx *sql.Tx
	if tx, err = r.store.Begin(); err != nil {
		return ErrorMsg(err)
	}
	_, err = tx.QueryContext(ctx, `delete from "meta_data" where data_uuid = $1`, dataUUID)
	if err != nil {
		return ErrorMsg(tx.Rollback())
	}
	// insert новых
	for _, item := range metaDataList {
		insert := tx.QueryRowContext(ctx, `insert into meta_data (meta_name, meta_value, data_uuid) values ($1, $2, $3)`, item.MetaName, item.MetaValue, item.DataUUID)
		err = insert.Err()
		if err != nil {
			return ErrorMsg(errors.Join(err, tx.Rollback()))
		}
	}

	if err = tx.Commit(); err != nil {
		return ErrorMsg(err)
	}
	return nil
}
