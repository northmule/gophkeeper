package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type OwnerRepositoryTestSuite struct {
	suite.Suite
	DB         *sql.DB
	mock       sqlmock.Sqlmock
	repository *OwnerRepository
}

func (s *OwnerRepositoryTestSuite) SetupTest() {
	var err error
	s.DB, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.mock.ExpectPrepare("select id, user_uuid, data_type, data_uuid from owner")
	s.repository, err = NewOwnerRepository(s.DB)
	require.NoError(s.T(), err)
}

func TestOwnerRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(OwnerRepositoryTestSuite))
}

func (s *OwnerRepositoryTestSuite) TestFindOneByUserUUIDAndDataUUIDAndDataType_ValidData() {
	userUUID := "user-uuid"
	dataUUID := "data-uuid"
	dataType := "data-type"
	expectedData := &models.Owner{
		ID:       1,
		UserUUID: userUUID,
		DataType: dataType,
		DataUUID: dataUUID,
	}

	s.mock.ExpectQuery("select id, user_uuid, data_type, data_uuid from owner").
		WithArgs(userUUID, dataUUID, dataType).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_uuid", "data_type", "data_uuid"}).
			AddRow(expectedData.ID, expectedData.UserUUID, expectedData.DataType, expectedData.DataUUID))

	data, err := s.repository.FindOneByUserUUIDAndDataUUIDAndDataType(context.Background(), userUUID, dataUUID, dataType)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), expectedData, data)
}

func (s *OwnerRepositoryTestSuite) TestFindOneByUserUUIDAndDataUUIDAndDataType_NoData() {
	userUUID := "user-uuid"
	dataUUID := "data-uuid"
	dataType := "data-type"

	s.mock.ExpectQuery("select id, user_uuid, data_type, data_uuid from owner").
		WithArgs(userUUID, dataUUID, dataType).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_uuid", "data_type", "data_uuid"}))

	data, err := s.repository.FindOneByUserUUIDAndDataUUIDAndDataType(context.Background(), userUUID, dataUUID, dataType)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), data)
}

func (s *OwnerRepositoryTestSuite) TestFindOneByUserUUIDAndDataUUIDAndDataType_Error() {
	userUUID := "user-uuid"
	dataUUID := "data-uuid"
	dataType := "data-type"

	s.mock.ExpectQuery("select id, user_uuid, data_type, data_uuid from owner").
		WithArgs(userUUID, dataUUID, dataType).
		WillReturnError(errors.New("query failed"))

	_, err := s.repository.FindOneByUserUUIDAndDataUUIDAndDataType(context.Background(), userUUID, dataUUID, dataType)
	require.Error(s.T(), err)
}

func (s *OwnerRepositoryTestSuite) TestFindOneByUserUUIDAndDataUUID_ValidData() {
	userUUID := "user-uuid"
	dataUUID := "data-uuid"
	expectedData := &models.Owner{
		ID:       1,
		UserUUID: userUUID,
		DataType: "data-type",
		DataUUID: dataUUID,
	}

	s.mock.ExpectQuery("select id, user_uuid, data_type, data_uuid from owner").
		WithArgs(userUUID, dataUUID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_uuid", "data_type", "data_uuid"}).
			AddRow(expectedData.ID, expectedData.UserUUID, expectedData.DataType, expectedData.DataUUID))

	data, err := s.repository.FindOneByUserUUIDAndDataUUID(context.Background(), userUUID, dataUUID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), expectedData, data)
}

func (s *OwnerRepositoryTestSuite) TestFindOneByUserUUIDAndDataUUID_NoData() {
	userUUID := "user-uuid"
	dataUUID := "data-uuid"

	s.mock.ExpectQuery("select id, user_uuid, data_type, data_uuid").
		WithArgs(userUUID, dataUUID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_uuid", "data_type", "data_uuid"}))

	data, err := s.repository.FindOneByUserUUIDAndDataUUID(context.Background(), userUUID, dataUUID)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), data)
}

func (s *OwnerRepositoryTestSuite) TestFindOneByUserUUIDAndDataUUID_Error() {
	userUUID := "user-uuid"
	dataUUID := "data-uuid"

	s.mock.ExpectQuery("select id, user_uuid, data_type, data_uuid").
		WithArgs(userUUID, dataUUID).
		WillReturnError(errors.New("query failed"))

	_, err := s.repository.FindOneByUserUUIDAndDataUUID(context.Background(), userUUID, dataUUID)
	require.Error(s.T(), err)
}

func (s *OwnerRepositoryTestSuite) TestAdd_ValidData() {
	data := &models.Owner{
		UserUUID: "new-uuid",
		DataType: "data-type",
		DataUUID: "new-data-uuid",
	}
	s.mock.ExpectQuery("insert into owner").
		WithArgs(data.UserUUID, data.DataType, data.DataUUID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := s.repository.Add(context.Background(), data)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), id)
}

func (s *OwnerRepositoryTestSuite) TestAdd_InvalidData() {
	data := &models.Owner{
		UserUUID: "",
		DataType: "data-type",
		DataUUID: "new-data-uuid",
	}
	id, err := s.repository.Add(context.Background(), data)
	require.Error(s.T(), err)
	assert.Equal(s.T(), int64(0), id)
}

func (s *OwnerRepositoryTestSuite) TestAdd_DuplicateUUID() {
	data := &models.Owner{
		UserUUID: "new-uuid",
		DataType: "data-type",
		DataUUID: "existing-uuid",
	}
	s.mock.ExpectQuery("insert into owner").
		WithArgs(data.UserUUID, data.DataType, data.DataUUID).
		WillReturnError(sql.ErrNoRows)

	id, err := s.repository.Add(context.Background(), data)
	require.Error(s.T(), err)
	assert.Equal(s.T(), int64(0), id)
}

func (s *OwnerRepositoryTestSuite) TestAllOwnerData_ValidData() {
	userUUID := "user-uuid"
	offset := 0
	limit := 10
	expectedData := []models.OwnerData{
		{
			DataType:     "data-type",
			DataUUID:     "data-uuid",
			UserUUID:     userUUID,
			DataName:     "data-name",
			DataTypeName: data_type.TranslateDataType("data-type"),
		},
	}

	s.mock.ExpectQuery("select o.data_type as data_type").
		WithArgs(userUUID, offset, limit).
		WillReturnRows(sqlmock.NewRows([]string{"data_type", "data_uuid", "user_uuid", "name"}).
			AddRow(expectedData[0].DataType, expectedData[0].DataUUID, expectedData[0].UserUUID, expectedData[0].DataName))

	data, err := s.repository.AllOwnerData(context.Background(), userUUID, offset, limit)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), expectedData, data)
}

func (s *OwnerRepositoryTestSuite) TestAllOwnerData_NoData() {
	userUUID := "user-uuid"
	offset := 0
	limit := 10

	s.mock.ExpectQuery("select o.data_type as data_type, o.data_uuid as data_uuid").
		WithArgs(userUUID, offset, limit).
		WillReturnRows(sqlmock.NewRows([]string{"data_type", "data_uuid", "user_uuid", "name"}))

	data, err := s.repository.AllOwnerData(context.Background(), userUUID, offset, limit)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), data)
}

func (s *OwnerRepositoryTestSuite) TestAllOwnerData_Error() {
	userUUID := "user-uuid"
	offset := 0
	limit := 10

	s.mock.ExpectQuery("select o.data_type as data_type, o.data_uuid as data_uuid").
		WithArgs(userUUID, offset, limit).
		WillReturnError(errors.New("query failed"))

	data, err := s.repository.AllOwnerData(context.Background(), userUUID, offset, limit)
	require.Error(s.T(), err)
	assert.Empty(s.T(), data)
}

func (s *OwnerRepositoryTestSuite) TestAllOwnerData_InvalidUserUUID() {
	userUUID := ""
	offset := 0
	limit := 10

	data, err := s.repository.AllOwnerData(context.Background(), userUUID, offset, limit)
	require.Error(s.T(), err)
	assert.Empty(s.T(), data)
}
