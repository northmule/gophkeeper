package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ManagerTestSuite struct {
	suite.Suite
	DB         *sql.DB
	mock       sqlmock.Sqlmock
	repository Repository
}

func (s *ManagerTestSuite) SetupTest() {
	var err error
	s.DB, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.mock.ExpectPrepare("select")
	s.mock.ExpectPrepare("insert")
	s.mock.ExpectPrepare("select")
	s.mock.ExpectPrepare("select")
	s.mock.ExpectPrepare("select")
	s.mock.ExpectPrepare("select")
	s.mock.ExpectPrepare("select")
	s.mock.ExpectPrepare("select")
	s.repository, err = NewManager(s.DB)
	require.NoError(s.T(), err)
}

func TestManagerTestSuite(t *testing.T) {
	suite.Run(t, new(ManagerTestSuite))
}

func (s *ManagerTestSuite) TestNewManager_UserRepositoryError() {
	s.mock.ExpectPrepare("select id, login, password_hash, public_key, private_client_key from users").WillReturnError(errors.New("user repo error"))
	_, err := NewManager(s.DB)
	require.Error(s.T(), err)
}

func (s *ManagerTestSuite) TestNewManager_CardDataRepositoryError() {
	s.mock.ExpectPrepare("select id, card_number, card_holder, expiration_date, cvv, data_uuid from card_data").WillReturnError(errors.New("card data repo error"))
	_, err := NewManager(s.DB)
	require.Error(s.T(), err)
}

func (s *ManagerTestSuite) TestNewManager_OwnerRepositoryError() {
	s.mock.ExpectPrepare("select id, user_uuid, data_type, data_uuid from owner").WillReturnError(errors.New("owner repo error"))
	_, err := NewManager(s.DB)
	require.Error(s.T(), err)
}

func (s *ManagerTestSuite) TestNewManager_MetaDataRepositoryError() {
	s.mock.ExpectPrepare("select id, meta_name, meta_value, data_uuid").WillReturnError(errors.New("meta data repo error"))
	_, err := NewManager(s.DB)
	require.Error(s.T(), err)
}

func (s *ManagerTestSuite) TestNewManager_TextDataRepositoryError() {
	s.mock.ExpectPrepare("select id, text, data_uuid from text_data").WillReturnError(errors.New("text data repo error"))
	_, err := NewManager(s.DB)
	require.Error(s.T(), err)
}

func (s *ManagerTestSuite) TestNewManager_FileDataRepositoryError() {
	s.mock.ExpectPrepare("select id, file_name, file_data, data_uuid from file_data").WillReturnError(errors.New("file data repo error"))
	_, err := NewManager(s.DB)
	require.Error(s.T(), err)
}

func (s *ManagerTestSuite) TestUser_ReturnsUserRepository() {
	userRepo := s.repository.User()
	assert.NotNil(s.T(), userRepo)
}

func (s *ManagerTestSuite) TestCardData_ReturnsCardDataRepository() {
	cardDataRepo := s.repository.CardData()
	assert.NotNil(s.T(), cardDataRepo)
}

func (s *ManagerTestSuite) TestOwner_ReturnsOwnerRepository() {
	ownerRepo := s.repository.Owner()
	assert.NotNil(s.T(), ownerRepo)
}

func (s *ManagerTestSuite) TestMetaData_ReturnsMetaDataRepository() {
	metaDataRepo := s.repository.MetaData()
	assert.NotNil(s.T(), metaDataRepo)
}

func (s *ManagerTestSuite) TestTextData_ReturnsTextDataRepository() {
	textDataRepo := s.repository.TextData()
	assert.NotNil(s.T(), textDataRepo)
}

func (s *ManagerTestSuite) TestFileData_ReturnsFileDataRepository() {
	fileDataRepo := s.repository.FileData()
	assert.NotNil(s.T(), fileDataRepo)
}
