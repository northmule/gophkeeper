package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type FileDataRepositoryTestSuite struct {
	suite.Suite
	DB         *sql.DB
	mock       sqlmock.Sqlmock
	repository *FileDataRepository
}

func (s *FileDataRepositoryTestSuite) SetupTest() {
	var err error
	s.DB, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.mock.ExpectPrepare("select id, name")
	s.repository, err = NewFileDataRepository(s.DB)
	require.NoError(s.T(), err)
}

func TestFileDataRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(FileDataRepositoryTestSuite))
}

func (s *FileDataRepositoryTestSuite) TestFindOneByUUID_ValidUUID() {
	uuid := "valid-uuid"
	expectedData := &models.FileData{
		Name:      "File Name",
		MimeType:  "application/pdf",
		Path:      "/path/to/file",
		PathTmp:   "/tmp/path/to/file",
		Extension: ".pdf",
		FileName:  "file.pdf",
		Size:      1024,
		Storage:   "local",
		Uploaded:  true,
	}
	expectedData.ID = 1
	expectedData.UUID = uuid
	s.mock.ExpectQuery("select").
		WithArgs(uuid).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "uuid", "mime_type", "path", "path_tmp", "extension", "file_name", "size", "storage", "uploaded"}).
			AddRow(expectedData.ID, expectedData.Name, expectedData.UUID, expectedData.MimeType, expectedData.Path, expectedData.PathTmp, expectedData.Extension, expectedData.FileName, expectedData.Size, expectedData.Storage, expectedData.Uploaded))

	data, err := s.repository.FindOneByUUID(context.Background(), uuid)
	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedData, data)
}

func (s *FileDataRepositoryTestSuite) TestFindOneByUUID_InvalidUUID() {
	uuid := "invalid-uuid"

	s.mock.ExpectQuery("select").
		WithArgs(uuid).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "uuid", "mime_type", "path", "path_tmp", "extension", "file_name", "size", "storage", "uploaded"}))

	data, err := s.repository.FindOneByUUID(context.Background(), uuid)
	require.NoError(s.T(), err)
	require.Empty(s.T(), data)
}

func (s *FileDataRepositoryTestSuite) TestFindOneByUUID_EmptyUUID() {
	uuid := ""

	data, err := s.repository.FindOneByUUID(context.Background(), uuid)
	require.Error(s.T(), err)
	require.Empty(s.T(), data)
}

func (s *FileDataRepositoryTestSuite) TestAdd_ValidData() {
	data := &models.FileData{
		Name:      "File Name",
		MimeType:  "application/pdf",
		Path:      "/path/to/file",
		PathTmp:   "/tmp/path/to/file",
		Extension: ".pdf",
		FileName:  "file.pdf",
		Size:      1024,
		Storage:   "local",
		Uploaded:  true,
	}
	data.UUID = "new-uuid"
	s.mock.ExpectQuery("insert into").
		WithArgs(data.Name, data.UUID, data.MimeType, data.Path, data.PathTmp, data.Extension, data.FileName, data.Size, data.Storage, data.Uploaded).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := s.repository.Add(context.Background(), data)
	require.NoError(s.T(), err)
	require.Equal(s.T(), int64(1), id)
}

func (s *FileDataRepositoryTestSuite) TestAdd_InvalidData() {
	data := &models.FileData{
		Name:      "",
		MimeType:  "application/pdf",
		Path:      "/path/to/file",
		PathTmp:   "/tmp/path/to/file",
		Extension: ".pdf",
		FileName:  "file.pdf",
		Size:      1024,
		Storage:   "local",
		Uploaded:  true,
	}
	data.UUID = "new-uuid"
	id, err := s.repository.Add(context.Background(), data)
	require.Error(s.T(), err)
	require.Equal(s.T(), int64(0), id)
}

func (s *FileDataRepositoryTestSuite) TestAdd_DuplicateUUID() {
	data := &models.FileData{
		Name:      "File Name",
		MimeType:  "application/pdf",
		Path:      "/path/to/file",
		PathTmp:   "/tmp/path/to/file",
		Extension: ".pdf",
		FileName:  "file.pdf",
		Size:      1024,
		Storage:   "local",
		Uploaded:  true,
	}
	data.UUID = "existing-uuid"
	s.mock.ExpectQuery("insert into").
		WithArgs(data.Name, data.UUID, data.MimeType, data.Path, data.PathTmp, data.Extension, data.FileName, data.Size, data.Storage, data.Uploaded).
		WillReturnError(sql.ErrNoRows)

	id, err := s.repository.Add(context.Background(), data)
	require.Error(s.T(), err)
	require.Equal(s.T(), int64(0), id)
}

func (s *FileDataRepositoryTestSuite) TestUpdate_ValidData() {
	data := &models.FileData{
		Name:      "Updated File Name",
		MimeType:  "application/pdf",
		Path:      "/path/to/file",
		PathTmp:   "/tmp/path/to/file",
		Extension: ".pdf",
		FileName:  "file.pdf",
		Size:      1024,
		Storage:   "local",
		Uploaded:  true,
	}
	data.UUID = "existing-uuid"
	s.mock.ExpectQuery("update file_data").
		WithArgs(data.Name, data.MimeType, data.Path, data.Extension, data.FileName, data.Size, data.Storage, data.Uploaded, data.UUID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err := s.repository.Update(context.Background(), data)
	require.NoError(s.T(), err)
}
