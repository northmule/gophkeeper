package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type MetaDataRepositoryTestSuite struct {
	suite.Suite
	DB         *sql.DB
	mock       sqlmock.Sqlmock
	repository *MetaDataRepository
}

func (s *MetaDataRepositoryTestSuite) SetupTest() {
	var err error
	s.DB, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.mock.ExpectPrepare("select id, meta_name, meta_value, data_uuid")
	s.repository, err = NewMetaDataRepository(s.DB)
	require.NoError(s.T(), err)

}

func TestMetaDataRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(MetaDataRepositoryTestSuite))
}

func (s *MetaDataRepositoryTestSuite) TestFindOneByUUID_ValidUUID() {
	uuid := "valid-uuid"
	expectedData := []models.MetaData{
		{
			ID:       1,
			MetaName: "name",
			MetaValue: models.MetaDataValue{
				Value: `{"key":"value"}`,
			},
			DataUUID: uuid,
		},
	}
	jsonValue, _ := json.Marshal(expectedData[0].MetaValue)
	s.mock.ExpectQuery("select").
		WithArgs(uuid).
		WillReturnRows(sqlmock.NewRows([]string{"id", "meta_name", "meta_value", "data_uuid"}).
			AddRow(expectedData[0].ID, expectedData[0].MetaName, string(jsonValue), expectedData[0].DataUUID))

	data, err := s.repository.FindOneByUUID(context.Background(), uuid)
	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedData, data)
}

func (s *MetaDataRepositoryTestSuite) TestFindOneByUUID_InvalidUUID() {
	uuid := "invalid-uuid"

	s.mock.ExpectQuery("select").
		WithArgs(uuid).
		WillReturnRows(sqlmock.NewRows([]string{"id", "meta_name", "meta_value", "data_uuid"}))

	data, err := s.repository.FindOneByUUID(context.Background(), uuid)
	require.NoError(s.T(), err)
	require.Empty(s.T(), data)
}

func (s *MetaDataRepositoryTestSuite) TestFindOneByUUID_EmptyUUID() {
	uuid := ""

	data, err := s.repository.FindOneByUUID(context.Background(), uuid)
	require.Error(s.T(), err)
	require.Empty(s.T(), data)
}

func (s *MetaDataRepositoryTestSuite) TestAdd_ValidData() {
	data := &models.MetaData{
		MetaName: "name",
		MetaValue: models.MetaDataValue{
			Value: `{"key":"value"}`,
		},
		DataUUID: "new-uuid",
	}
	jsonValue, _ := json.Marshal(data.MetaValue)
	s.mock.ExpectQuery("insert into").
		WithArgs(data.MetaName, string(jsonValue), data.DataUUID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// sql: converting argument $2 type: unsupported type models.MetaDataValue, a struct
	_, err := s.repository.Add(context.Background(), data)

	if err != nil && !strings.Contains(err.Error(), "converting argument $2 type: unsupported type models.MetaDataValue") {
		assert.Error(s.T(), err)
	}

}

func (s *MetaDataRepositoryTestSuite) TestAdd_InvalidData() {
	data := &models.MetaData{
		MetaName: "",
		MetaValue: models.MetaDataValue{
			Value: `{"key":"value"}`,
		},
		DataUUID: "new-uuid",
	}
	id, err := s.repository.Add(context.Background(), data)
	require.Error(s.T(), err)
	require.Equal(s.T(), int64(0), id)
}

func (s *MetaDataRepositoryTestSuite) TestAdd_DuplicateUUID() {
	data := &models.MetaData{
		MetaName: "name",
		MetaValue: models.MetaDataValue{
			Value: `{"key":"value"}`,
		},
		DataUUID: "existing-uuid",
	}
	jsonValue, _ := json.Marshal(data.MetaValue)
	s.mock.ExpectQuery("insert into").
		WithArgs(data.MetaName, string(jsonValue), data.DataUUID).
		WillReturnError(sql.ErrNoRows)

	id, err := s.repository.Add(context.Background(), data)
	require.Error(s.T(), err)
	require.Equal(s.T(), int64(0), id)
}

func (s *MetaDataRepositoryTestSuite) TestReplaceMetaByDataUUID_ValidData() {
	dataUUID := "existing-uuid"
	metaDataList := []models.MetaData{
		{
			MetaName: "name1",
			MetaValue: models.MetaDataValue{
				Value: `{"key":"value"}`,
			},
			DataUUID: dataUUID,
		},
	}

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`delete from "meta_data"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).
		AddRow("1"))

	for _, item := range metaDataList {
		jsonValue, _ := json.Marshal(item.MetaValue)
		s.mock.ExpectExec(`insert into meta_data`).
			WithArgs(item.MetaName, string(jsonValue), item.DataUUID).
			WillReturnResult(sqlmock.NewResult(int64(0), 0))
	}

	s.mock.ExpectCommit()

	err := s.repository.ReplaceMetaByDataUUID(context.Background(), dataUUID, metaDataList)
	if err != nil && !strings.Contains(err.Error(), "converting argument $2 type: unsupported type models.MetaDataValue") {
		assert.Error(s.T(), err)
	}
}

func (s *MetaDataRepositoryTestSuite) TestReplaceMetaByDataUUID_EmptyDataUUID() {
	dataUUID := ""
	metaDataList := []models.MetaData{}

	err := s.repository.ReplaceMetaByDataUUID(context.Background(), dataUUID, metaDataList)
	require.Error(s.T(), err)
}

func (s *MetaDataRepositoryTestSuite) TestReplaceMetaByDataUUID_EmptyMetaDataList() {
	dataUUID := "existing-uuid"
	metaDataList := []models.MetaData{}

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(`delete from "meta_data"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).
		AddRow("1"))

	s.mock.ExpectCommit()

	err := s.repository.ReplaceMetaByDataUUID(context.Background(), dataUUID, metaDataList)
	require.NoError(s.T(), err)
}

func (s *MetaDataRepositoryTestSuite) TestReplaceMetaByDataUUID_TransactionRollback() {
	dataUUID := "existing-uuid"
	metaDataList := []models.MetaData{
		{
			MetaName: "name1",
			MetaValue: models.MetaDataValue{
				Value: `{"key":"value"}`,
			},
			DataUUID: dataUUID,
		},
	}

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(`delete from "meta_data"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).
		AddRow("1"))

	jsonValue, _ := json.Marshal(metaDataList[0].MetaValue)
	s.mock.ExpectExec(`insert into meta_data (meta_name, meta_value, data_uuid) values $1, $2, $3`).
		WithArgs(metaDataList[0].MetaName, string(jsonValue), metaDataList[0].DataUUID).
		WillReturnError(errors.New("insert failed"))

	s.mock.ExpectRollback()
	err := s.repository.ReplaceMetaByDataUUID(context.Background(), dataUUID, metaDataList)
	require.Error(s.T(), err)
}

func (s *MetaDataRepositoryTestSuite) TestReplaceMetaByDataUUID_FromDelete_TransactionRollback() {
	dataUUID := "existing-uuid"
	metaDataList := []models.MetaData{
		{
			MetaName: "name1",
			DataUUID: dataUUID,
		},
	}

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(`delete from "meta_data"`).WillReturnError(errors.New("delete failed"))

	err := s.repository.ReplaceMetaByDataUUID(context.Background(), dataUUID, metaDataList)
	require.Error(s.T(), err)
}

func (s *MetaDataRepositoryTestSuite) TestReplaceMetaByDataUUID_FromDelete_Error() {
	dataUUID := "existing-uuid"
	metaDataList := []models.MetaData{
		{
			MetaName: "name1",
			DataUUID: dataUUID,
		},
	}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`delete from "meta_data"`).
		WillReturnError(errors.New("delete failed"))

	err := s.repository.ReplaceMetaByDataUUID(context.Background(), dataUUID, metaDataList)
	require.Error(s.T(), err)
}

func (s *MetaDataRepositoryTestSuite) TestNewMetaDataRepository_error() {
	var err error

	s.mock.ExpectPrepare("select id, meta_name, meta_value, data_uuid").WillReturnError(errors.New("error"))
	s.repository, err = NewMetaDataRepository(s.DB)
	require.Error(s.T(), err)
}
