package repository

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	DB         *sql.DB
	mock       sqlmock.Sqlmock
	repository *UserRepository
	TX         *storage.Transaction
}

func (s *UserRepositoryTestSuite) SetupTest() {
	var err error
	s.DB, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.mock.ExpectPrepare("select")
	s.mock.ExpectPrepare("insert")
	s.mock.ExpectPrepare("select")
	s.repository, err = NewUserRepository(s.DB)

	transaction, _ := storage.NewTransaction(s.DB)
	s.TX = transaction
	require.NoError(s.T(), err)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (s *UserRepositoryTestSuite) TestFindOneByLogin_ValidData() {
	login := "test-login"
	expectedUser := &models.User{
		Login:            login,
		Password:         "test-password",
		CreatedAt:        time.Now(),
		Email:            "test-email",
		PublicKey:        "test-public-key",
		PrivateClientKey: "test-private-client-key",
	}
	expectedUser.ID = 1
	expectedUser.UUID = "test-uuid"

	s.mock.ExpectQuery("select id, login, password").
		WithArgs(login).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password", "created_at", "uuid", "email", "public_key", "private_client_key"}).
			AddRow(expectedUser.ID, expectedUser.Login, expectedUser.Password, expectedUser.CreatedAt, expectedUser.UUID, expectedUser.Email, expectedUser.PublicKey, expectedUser.PrivateClientKey))

	user, err := s.repository.FindOneByLogin(context.Background(), login)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), expectedUser, user)
}

func (s *UserRepositoryTestSuite) TestFindOneByLogin_NoData() {
	login := "test-login"

	s.mock.ExpectQuery("select id, login").
		WithArgs(login).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password", "created_at", "uuid", "email", "public_key", "private_client_key"}))

	user, err := s.repository.FindOneByLogin(context.Background(), login)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), user)
}

func (s *UserRepositoryTestSuite) TestFindOneByLogin_Error() {
	login := "test-login"

	s.mock.ExpectQuery("select id, login").
		WithArgs(login).
		WillReturnError(errors.New("query failed"))

	_, err := s.repository.FindOneByLogin(context.Background(), login)
	require.Error(s.T(), err)
}

func (s *UserRepositoryTestSuite) TestFindOneByUUID_ValidData() {
	uuid := "test-uuid"
	expectedUser := &models.User{
		Login:            "test-login",
		Password:         "test-password",
		CreatedAt:        time.Now(),
		Email:            "test-email",
		PublicKey:        "test-public-key",
		PrivateClientKey: "test-private-client-key",
	}
	expectedUser.ID = 1
	expectedUser.UUID = uuid
	s.mock.ExpectQuery("select id, login, password").
		WithArgs(uuid).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password", "created_at", "uuid", "email", "public_key", "private_client_key"}).
			AddRow(expectedUser.ID, expectedUser.Login, expectedUser.Password, expectedUser.CreatedAt, expectedUser.UUID, expectedUser.Email, expectedUser.PublicKey, expectedUser.PrivateClientKey))

	user, err := s.repository.FindOneByUUID(context.Background(), uuid)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), expectedUser, user)
}

func (s *UserRepositoryTestSuite) TestFindOneByUUID_NoData() {
	uuid := "test-uuid"

	s.mock.ExpectQuery("select id, login, password").
		WithArgs(uuid).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password", "created_at", "uuid", "email", "public_key", "private_client_key"}))

	user, err := s.repository.FindOneByUUID(context.Background(), uuid)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), user)
}

func (s *UserRepositoryTestSuite) TestFindOneByUUID_Error() {
	uuid := "test-uuid"

	s.mock.ExpectQuery("select id, login, password").
		WithArgs(uuid).
		WillReturnError(errors.New("query failed"))

	_, err := s.repository.FindOneByUUID(context.Background(), uuid)
	require.Error(s.T(), err)
}

func (s *UserRepositoryTestSuite) TestCreateNewUser_ValidData() {
	user := models.User{
		Login:    "new-login",
		Password: "new-password",
		Email:    "new-email",
	}
	user.UUID = "new-uuid"
	s.mock.ExpectQuery("insert into users").
		WithArgs(user.Login, user.Password, user.UUID, user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := s.repository.CreateNewUser(context.Background(), user)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), id)
}

func (s *UserRepositoryTestSuite) TestCreateNewUser_InvalidData() {
	user := models.User{
		Login:    "",
		Password: "new-password",
		Email:    "new-email",
	}
	user.UUID = "new-uuid"
	id, err := s.repository.CreateNewUser(context.Background(), user)
	require.Error(s.T(), err)
	assert.Equal(s.T(), int64(0), id)
}

func (s *UserRepositoryTestSuite) TestCreateNewUser_DuplicateUUID() {
	user := models.User{
		Login:    "new-login",
		Password: "new-password",
		Email:    "new-email",
	}
	user.UUID = "existing-uuid"
	s.mock.ExpectQuery("insert into users").
		WithArgs(user.Login, user.Password, user.UUID, user.Email).
		WillReturnError(sql.ErrNoRows)

	id, err := s.repository.CreateNewUser(context.Background(), user)
	require.Error(s.T(), err)
	assert.Equal(s.T(), int64(0), id)
}

func (s *UserRepositoryTestSuite) TestSetPublicKey_ValidData() {
	data := "new-public-key"
	userUUID := "test-uuid"

	s.mock.ExpectQuery("update users").
		WithArgs(data, userUUID).WillReturnRows(sqlmock.NewRows([]string{"id"}))

	err := s.repository.SetPublicKey(context.Background(), data, userUUID)
	require.NoError(s.T(), err)
}

func (s *UserRepositoryTestSuite) TestSetPublicKey_InvalidData() {
	data := ""
	userUUID := "test-uuid"

	err := s.repository.SetPublicKey(context.Background(), data, userUUID)
	require.Error(s.T(), err)
}

func (s *UserRepositoryTestSuite) TestSetPublicKey_Error() {
	data := "new-public-key"
	userUUID := "test-uuid"

	s.mock.ExpectExec("update users set").
		WithArgs(data, userUUID).
		WillReturnError(errors.New("update failed"))

	err := s.repository.SetPublicKey(context.Background(), data, userUUID)
	require.Error(s.T(), err)
}

func (s *UserRepositoryTestSuite) TestSetPrivateClientKey_ValidData() {
	data := "new-private-client-key"
	userUUID := "test-uuid"

	s.mock.ExpectQuery("update users set").
		WithArgs(data, userUUID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	err := s.repository.SetPrivateClientKey(context.Background(), data, userUUID)
	require.NoError(s.T(), err)
}

func (s *UserRepositoryTestSuite) TestSetPrivateClientKey_InvalidData() {
	data := ""
	userUUID := "test-uuid"

	err := s.repository.SetPrivateClientKey(context.Background(), data, userUUID)
	require.Error(s.T(), err)
}

func (s *UserRepositoryTestSuite) TestSetPrivateClientKey_Error() {
	data := "new-private-client-key"
	userUUID := "test-uuid"

	s.mock.ExpectExec("update users set").
		WithArgs(data, userUUID).
		WillReturnError(errors.New("update failed"))

	err := s.repository.SetPrivateClientKey(context.Background(), data, userUUID)
	require.Error(s.T(), err)
}

func (s *UserRepositoryTestSuite) TestCreateNewUser_InvalidUUID() {
	user := models.User{
		Login:    "new-login",
		Password: "new-password",
		Email:    "new-email",
	}
	id, err := s.repository.CreateNewUser(context.Background(), user)
	require.Error(s.T(), err)
	assert.Equal(s.T(), int64(0), id)
}

func (s *UserRepositoryTestSuite) TestCreateNewUser_InvalidEmail() {
	user := models.User{
		Login:    "new-login",
		Password: "new-password",
		Email:    "",
	}
	user.UUID = "new-uuid"
	id, err := s.repository.CreateNewUser(context.Background(), user)
	require.Error(s.T(), err)
	assert.Equal(s.T(), int64(0), id)
}

func (s *UserRepositoryTestSuite) TestFindOneByLogin_InvalidLogin() {
	login := ""

	_, err := s.repository.FindOneByLogin(context.Background(), login)
	require.Error(s.T(), err)
}

func (s *UserRepositoryTestSuite) TestFindOneByUUID_InvalidUUID() {
	uuid := ""

	_, err := s.repository.FindOneByUUID(context.Background(), uuid)
	require.Error(s.T(), err)
}
