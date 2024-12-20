package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/storage"
)

type UserRepository struct {
	store          storage.DBQuery
	sqlFindByLogin *sql.Stmt
	sqlCreateUser  *sql.Stmt
	sqlFindByUUID  *sql.Stmt
}

const timeOut = 5000000 * time.Second // todo only debug

type ErrorMsg error

func NewUserRepository(store storage.DBQuery) (*UserRepository, error) {
	instance := UserRepository{
		store: store,
	}
	var err error
	instance.sqlFindByLogin, err = store.Prepare(`select id, login, password, created_at, uuid, email, public_key, private_client_key from users where login = $1 limit 1`)
	if err != nil {
		return nil, ErrorMsg(err)
	}

	instance.sqlCreateUser, err = store.Prepare(`insert into users (login, password, email, uuid) values ($1, $2, $3, $4) returning id`)
	if err != nil {
		return nil, ErrorMsg(err)
	}

	instance.sqlFindByUUID, err = store.Prepare(`select id, login, password, created_at, uuid, email, public_key, private_client_key from users where uuid = $1 limit 1`)
	if err != nil {
		return nil, ErrorMsg(err)
	}

	return &instance, nil
}

func (r *UserRepository) FindOneByLogin(ctx context.Context, login string) (*models.User, error) {
	user := models.User{}
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	rows, err := r.sqlFindByLogin.QueryContext(ctx, login)
	if err != nil {
		return nil, ErrorMsg(err)
	}
	err = rows.Err()
	if err != nil {
		return nil, ErrorMsg(err)
	}

	if rows.Next() {
		err = rows.Scan(&user.ID, &user.Login, &user.Password, &user.CreatedAt, &user.UUID, &user.Email, &user.PublicKey, &user.PrivateClientKey)
		if err != nil {
			return nil, ErrorMsg(err)
		}
	}

	return &user, nil
}

func (r *UserRepository) FindOneByUUID(ctx context.Context, uuid string) (*models.User, error) {
	user := models.User{}
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	rows, err := r.sqlFindByUUID.QueryContext(ctx, uuid)
	if err != nil {
		return nil, ErrorMsg(err)
	}
	err = rows.Err()
	if err != nil {
		return nil, ErrorMsg(err)
	}

	if rows.Next() {
		err = rows.Scan(&user.ID, &user.Login, &user.Password, &user.CreatedAt, &user.UUID, &user.Email, &user.PublicKey, &user.PrivateClientKey)
		if err != nil {
			return nil, ErrorMsg(err)
		}
	}

	return &user, nil
}

func (r *UserRepository) CreateNewUser(ctx context.Context, user models.User) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	rows := r.sqlCreateUser.QueryRowContext(ctx, user.Login, user.Password, user.UUID, user.Email)
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

func (r *UserRepository) TxCreateNewUser(ctx context.Context, tx storage.TxDBQuery, user models.User) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	rows := tx.Tx().QueryRowContext(ctx, `insert into users (login, password, uuid, email) values ($1, $2, $3, $4) returning id`, user.Login, user.Password, user.UUID, user.Email)
	err := rows.Err()
	if err != nil {
		tx.AddError(ErrorMsg(err))
		return 0, ErrorMsg(err)
	}

	var id int64
	err = rows.Scan(&id)
	if err != nil {
		tx.AddError(err)
		return 0, ErrorMsg(err)
	}
	return id, nil
}

// SetPublicKey вставка значения публичного ключа
func (r *UserRepository) SetPublicKey(ctx context.Context, data string, userUUID string) error {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	rows := r.store.QueryRowContext(ctx, `update users set public_key = $1 where uuid = $2`, data, userUUID)
	err := rows.Err()
	if err != nil {
		return ErrorMsg(err)
	}

	return nil
}

// SetPrivateClientKey вставка значения публичного ключа
func (r *UserRepository) SetPrivateClientKey(ctx context.Context, data string, userUUID string) error {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()
	rows := r.store.QueryRowContext(ctx, `update users set private_client_key = $1 where uuid = $2`, data, userUUID)
	err := rows.Err()
	if err != nil {
		return ErrorMsg(err)
	}

	return nil
}
