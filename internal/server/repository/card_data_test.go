package repository

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type CardDataRepositoryTestSuite struct {
	suite.Suite
	DB         *sql.DB
	mock       sqlmock.Sqlmock
	repository *CardDataRepository
}

func (s *CardDataRepositoryTestSuite) SetupTest() {
	var err error
	s.DB, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	s.mock.ExpectPrepare("select id, value")
	s.repository, err = NewCardDataRepository(s.DB)
	require.NoError(s.T(), err)
}

func TestCardDataRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CardDataRepositoryTestSuite))
}

func (s *CardDataRepositoryTestSuite) TestFindOneByUUID_ValidUUID() {
	uuid := "valid-uuid"
	expectedData := &models.CardData{
		Value:      models.CardDataValueV1{},
		ObjectType: "card",
		Name:       "Card Name",
	}
	expectedData.UUID = uuid
	expectedData.ID = 1
	jsonValue, err := json.Marshal(expectedData.Value)
	require.NoError(s.T(), err)

	s.mock.ExpectQuery("select").
		WithArgs(uuid).
		WillReturnRows(sqlmock.NewRows([]string{"id", "value", "object_type", "name", "uuid"}).
			AddRow(expectedData.ID, string(jsonValue), expectedData.ObjectType, expectedData.Name, expectedData.UUID))

	data, err := s.repository.FindOneByUUID(context.Background(), uuid)
	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedData, data)
}

func (s *CardDataRepositoryTestSuite) TestFindOneByUUID_InvalidUUID() {
	uuid := "invalid-uuid"

	s.mock.ExpectQuery("select").
		WithArgs(uuid).
		WillReturnRows(sqlmock.NewRows([]string{"id", "value", "object_type", "name", "uuid"}))

	data, err := s.repository.FindOneByUUID(context.Background(), uuid)
	require.NoError(s.T(), err)
	require.Empty(s.T(), data)
}

func (s *CardDataRepositoryTestSuite) TestFindOneByUUID_EmptyUUID() {
	uuid := ""

	data, err := s.repository.FindOneByUUID(context.Background(), uuid)
	require.Error(s.T(), err)
	require.Nil(s.T(), data)
}

func (s *CardDataRepositoryTestSuite) TestAdd_ValidData() {
	data := &models.CardData{
		Value:      models.CardDataValueV1{},
		ObjectType: "card",
		Name:       "Card Name",
	}
	data.UUID = "new-uuid"

	s.mock.ExpectQuery("insert into").
		WithArgs(data.Name, data.ObjectType, data.Value, data.UUID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	_, _ = s.repository.Add(context.Background(), data)

}

func (s *CardDataRepositoryTestSuite) TestAdd_InvalidData() {
	data := &models.CardData{
		Value:      models.CardDataValueV1{},
		ObjectType: "card",
		Name:       "",
	}
	data.UUID = "new-uuid"
	id, err := s.repository.Add(context.Background(), data)
	require.Error(s.T(), err)
	require.Equal(s.T(), int64(0), id)
}
func (s *CardDataRepositoryTestSuite) TestAdd_DuplicateUUID() {
	data := &models.CardData{
		Value:      models.CardDataValueV1{},
		ObjectType: "card",
		Name:       "Card Name",
	}
	data.UUID = "new-uuid"
	jsonValue, err := json.Marshal(data.Value)
	require.NoError(s.T(), err)

	s.mock.ExpectQuery("insert into").
		WithArgs(data.Name, data.ObjectType, string(jsonValue), data.UUID).
		WillReturnError(sql.ErrNoRows)

	id, err := s.repository.Add(context.Background(), data)
	require.Error(s.T(), err)
	require.Equal(s.T(), int64(0), id)
}

func (s *CardDataRepositoryTestSuite) TestUpdate_ValidData() {
	data := &models.CardData{
		Value:      models.CardDataValueV1{},
		ObjectType: "card",
		Name:       "Updated Card Name",
	}
	data.UUID = "existing-uuid"
	jsonValue, err := json.Marshal(data.Value)
	require.NoError(s.T(), err)

	s.mock.ExpectQuery("update card_data set name = $1, value = $2 where uuid = $3").
		WithArgs(data.Name, string(jsonValue), data.UUID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.repository.Update(context.Background(), data)
}

func (s *CardDataRepositoryTestSuite) TestUpdate_InvalidData() {
	data := &models.CardData{
		Value:      models.CardDataValueV1{},
		ObjectType: "card",
		Name:       "",
	}
	data.UUID = "existing-uuid"
	err := s.repository.Update(context.Background(), data)
	require.Error(s.T(), err)
}

func (s *CardDataRepositoryTestSuite) TestUpdate_NonExistentUUID() {
	data := &models.CardData{
		Value:      models.CardDataValueV1{},
		ObjectType: "card",
		Name:       "Updated Card Name",
	}
	data.UUID = "non-existent-uuid"
	jsonValue, err := json.Marshal(data.Value)
	require.NoError(s.T(), err)

	s.mock.ExpectQuery("update card_data set name = $1, value = $2 where uuid = $3").
		WithArgs(data.Name, string(jsonValue), data.UUID).
		WillReturnError(sql.ErrNoRows)

	err = s.repository.Update(context.Background(), data)
	require.Error(s.T(), err)
}
