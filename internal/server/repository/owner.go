package repository

import (
	"context"
	"database/sql"

	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/storage"
)

// OwnerRepository репозитарий владельца данных
type OwnerRepository struct {
	store                                      storage.DBQuery
	sqlFindByDataUUID                          *sql.Stmt
	sqlFindByUserUUID                          *sql.Stmt
	sqlFindOneByUserUUIDAndDataUUIDAndDataType *sql.Stmt
}

// NewOwnerRepository конструктор
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

// FindOneByUserUUIDAndDataUUIDAndDataType запрос данных
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

// FindOneByUserUUIDAndDataUUID запрос данных
func (r *OwnerRepository) FindOneByUserUUIDAndDataUUID(ctx context.Context, userUuid string, dataUuid string) (*models.Owner, error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	rows, err := r.store.QueryContext(ctx, `select id, user_uuid, data_type, data_uuid from owner where user_uuid = $1 and data_uuid = $2 limit 1`, userUuid, dataUuid)
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

func (r *OwnerRepository) AllOwnerData(ctx context.Context, userUUID string, offset int, limit int) ([]models.OwnerData, error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	query := `select 
o.data_type as data_type,
o.data_uuid as data_uuid,
o.user_uuid as user_uuid,
coalesce(cd."name", fd."name", td."name") as "name" 
from owner o
left join card_data cd on cd."uuid"  = o.data_uuid 
left join file_data fd on fd."uuid"  = o.data_uuid 
left join text_data td on td."uuid"  = o.data_uuid 
where o.user_uuid  = $1
order by o.id asc
offset $2 limit $3
`

	rows, err := r.store.QueryContext(ctx, query, userUUID, offset, limit)
	if err != nil {
		return nil, ErrorMsg(err)
	}
	err = rows.Err()
	if err != nil {
		return nil, ErrorMsg(err)
	}

	var dataList []models.OwnerData
	for rows.Next() {
		data := models.OwnerData{}
		err = rows.Scan(&data.DataType, &data.DataUUID, &data.UserUUID, &data.DataName)
		if err != nil {
			return nil, ErrorMsg(err)
		}
		data.DataTypeName = data_type.TranslateDataType(data.DataType)
		dataList = append(dataList, data)
	}

	return dataList, nil
}
