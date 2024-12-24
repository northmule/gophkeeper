package repository

import (
	"context"
	"database/sql"

	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/storage"
)

// FileDataRepository репозитарий файлов
type FileDataRepository struct {
	store storage.DBQuery

	sqlFindOneByUUID *sql.Stmt
}

// NewFileDataRepository конструктор
func NewFileDataRepository(store storage.DBQuery) (*FileDataRepository, error) {
	var err error
	instance := new(FileDataRepository)
	instance.store = store
	instance.sqlFindOneByUUID, err = store.Prepare(`select id, name, uuid, mime_type, path, path_tmp, extension, file_name, "size", storage, uploaded from file_data where uuid = $1 limit 1`)
	if err != nil {
		return nil, ErrorMsg(err)
	}
	return instance, nil
}

// FindOneByUUID поиск значения по UUID
func (r *FileDataRepository) FindOneByUUID(ctx context.Context, uuid string) (*models.FileData, error) {

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
	data := new(models.FileData)
	if rows.Next() {
		err = rows.Scan(&data.ID, &data.Name, &data.UUID, &data.MimeType, &data.Path, &data.PathTmp, &data.Extension, &data.FileName, &data.Size, &data.Storage, &data.Uploaded)
		if err != nil {
			return nil, ErrorMsg(err)
		}

	}

	return data, nil
}

// Add Новое значение
func (r *FileDataRepository) Add(ctx context.Context, data *models.FileData) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	rows := r.store.QueryRowContext(ctx, `insert into file_data (name, uuid, mime_type, path, path_tmp, extension, file_name, "size", storage, uploaded) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id`, data.Name, data.UUID, data.MimeType, data.Path, data.PathTmp, data.Extension, data.FileName, data.Size, data.Storage, data.Uploaded)
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
func (r *FileDataRepository) Update(ctx context.Context, data *models.FileData) error {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	rows := r.store.QueryRowContext(ctx, `update file_data set name = $1, mime_type = $2, path = $3, extension = $4, file_name = $5, size = $6, storage = $7, uploaded = $8 where uuid = $9`, data.Name, data.MimeType, data.Path, data.Extension, data.FileName, data.Size, data.Storage, data.Uploaded, data.UUID)

	return rows.Err()
}
