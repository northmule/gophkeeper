package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type TextDataRepositoryTestSuite struct {
	suite.Suite
	DB         *sql.DB
	mock       sqlmock.Sqlmock
	repository *TextDataRepository
}

func (s *TextDataRepositoryTestSuite) SetupTest() {
	var err error
	s.DB, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	s.mock.ExpectPrepare("select id, name, value")
	s.repository, err = NewTextDataRepository(s.DB)
	require.NoError(s.T(), err)
}

func TestTextDataRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TextDataRepositoryTestSuite))
}

func (s *TextDataRepositoryTestSuite) TestFindOneByUUID_ValidData() {
	uuid := "test-uuid"
	expectedData := &models.TextData{
		Name:  "test-name",
		Value: "test-value",
	}
	expectedData.UUID = uuid
	expectedData.ID = 1
	s.mock.ExpectQuery("select id, name").
		WithArgs(uuid).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "value", "uuid"}).
			AddRow(expectedData.ID, expectedData.Name, expectedData.Value, expectedData.UUID))

	data, err := s.repository.FindOneByUUID(context.Background(), uuid)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), expectedData, data)
}

func (s *TextDataRepositoryTestSuite) TestFindOneByUUID_NoData() {
	uuid := "test-uuid"

	s.mock.ExpectQuery("select id, name").
		WithArgs(uuid).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "value", "uuid"}))

	data, err := s.repository.FindOneByUUID(context.Background(), uuid)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), data)
}

func (s *TextDataRepositoryTestSuite) TestFindOneByUUID_Error() {
	uuid := "test-uuid"

	s.mock.ExpectQuery("select id, name").
		WithArgs(uuid).
		WillReturnError(errors.New("query failed"))

	_, err := s.repository.FindOneByUUID(context.Background(), uuid)
	require.Error(s.T(), err)
}

func (s *TextDataRepositoryTestSuite) TestAdd_ValidData() {
	data := &models.TextData{
		Name:  "test-name",
		Value: "test-value",
	}
	data.UUID = "test-uuid"

	s.mock.ExpectQuery("insert into text_data").
		WithArgs(data.Name, data.Value, data.UUID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := s.repository.Add(context.Background(), data)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), id)
}

func (s *TextDataRepositoryTestSuite) TestAdd_InvalidData() {
	data := &models.TextData{
		Name:  "",
		Value: "test-value",
	}
	data.UUID = "test-uuid"
	id, err := s.repository.Add(context.Background(), data)
	require.Error(s.T(), err)
	assert.Equal(s.T(), int64(0), id)
}

func (s *TextDataRepositoryTestSuite) TestAdd_DuplicateUUID() {
	data := &models.TextData{
		Name:  "test-name",
		Value: "test-value",
	}
	data.UUID = "existing-uuid"
	s.mock.ExpectQuery("insert into text_data").
		WithArgs(data.Name, data.Value, data.UUID).
		WillReturnError(sql.ErrNoRows)

	id, err := s.repository.Add(context.Background(), data)
	require.Error(s.T(), err)
	assert.Equal(s.T(), int64(0), id)
}

func (s *TextDataRepositoryTestSuite) TestUpdate_ValidData() {
	data := &models.TextData{
		Name:  "updated-name",
		Value: "updated-value",
	}
	data.ID = 1
	data.UUID = "test-uuid"
	s.mock.ExpectQuery("update text_data").
		WithArgs(data.Name, data.Value, data.UUID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err := s.repository.Update(context.Background(), data)
	require.NoError(s.T(), err)
}

func (s *TextDataRepositoryTestSuite) TestUpdate_InvalidData() {
	data := &models.TextData{
		Name:  "",
		Value: "updated-value",
	}
	data.ID = 1
	data.UUID = "test-uuid"
	err := s.repository.Update(context.Background(), data)
	require.Error(s.T(), err)
}

func (s *TextDataRepositoryTestSuite) TestUpdate_Error() {
	data := &models.TextData{
		Name:  "updated-name",
		Value: "updated-value",
	}
	data.ID = 1
	data.UUID = "test-uuid"
	s.mock.ExpectExec("update text_data").
		WithArgs(data.Name, data.Value, data.UUID).
		WillReturnError(errors.New("update failed"))

	err := s.repository.Update(context.Background(), data)
	require.Error(s.T(), err)
}

func (s *TextDataRepositoryTestSuite) TestFindOneByUUID_InvalidUUID() {
	uuid := ""

	_, err := s.repository.FindOneByUUID(context.Background(), uuid)
	require.Error(s.T(), err)
}

func (s *TextDataRepositoryTestSuite) TestAdd_EmptyUUID() {
	data := &models.TextData{
		Name:  "test-name",
		Value: "test-value",
	}

	id, err := s.repository.Add(context.Background(), data)
	require.Error(s.T(), err)
	assert.Equal(s.T(), int64(0), id)
}

func (s *TextDataRepositoryTestSuite) TestAdd_DuplicateData() {
	data := &models.TextData{
		Name:  "test-name",
		Value: "test-value",
	}
	data.UUID = "existing-uuid"
	s.mock.ExpectQuery("insert into text_data").
		WithArgs(data.Name, data.Value, data.UUID).
		WillReturnError(sql.ErrNoRows)

	id, err := s.repository.Add(context.Background(), data)
	require.Error(s.T(), err)
	assert.Equal(s.T(), int64(0), id)
}
