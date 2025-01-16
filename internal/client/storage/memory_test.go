package storage

import (
	"testing"

	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/stretchr/testify/assert"
)

func TestNewMemoryStorage_Success(t *testing.T) {
	storage := NewMemoryStorage()
	assert.NotNil(t, storage)
	assert.Empty(t, storage.token)
	assert.Empty(t, storage.cardDataList)
	assert.Empty(t, storage.metaDataList)
	assert.Empty(t, storage.textDataList)
	assert.Empty(t, storage.fileDataList)
}

func TestMemoryStorage_SetToken(t *testing.T) {
	storage := NewMemoryStorage()
	storage.SetToken("test_token")
	assert.Equal(t, "test_token", storage.Token())
}

func TestMemoryStorage_ResetToken(t *testing.T) {
	storage := NewMemoryStorage()
	storage.SetToken("test_token")
	storage.ResetToken()
	assert.Empty(t, storage.Token())
}

func TestMemoryStorage_AddCardDataList_Success(t *testing.T) {
	storage := NewMemoryStorage()
	cardData := models.CardData{
		Name: "card1123",
		Value: models.CardDataValueV1{
			CardNumber:     "1234567890123456",
			FullNameHolder: "John Doe",
		},
	}
	cardData.UUID = "card1"
	err := storage.AddCardDataList(cardData)
	assert.NoError(t, err)

	retrievedCardData, exists := storage.cardDataList["card1"]
	assert.True(t, exists)
	assert.Equal(t, cardData, retrievedCardData)
}

func TestMemoryStorage_AddMetaDataList_Success(t *testing.T) {
	storage := NewMemoryStorage()
	metaData := models.MetaData{
		DataUUID: "meta1",
		MetaName: "card",
	}

	err := storage.AddMetaDataList(metaData)
	assert.NoError(t, err)

	var retrievedMetaData models.MetaData
	for _, item := range storage.metaDataList {
		if item.DataUUID == "meta1" {
			retrievedMetaData = item
			break
		}
	}

	assert.Equal(t, metaData, retrievedMetaData)
}

func TestMemoryStorage_AddTextData_Success(t *testing.T) {
	storage := NewMemoryStorage()
	textData := models.TextData{
		Name:  "text123",
		Value: "This is a test text",
	}
	textData.UUID = "text1"
	err := storage.AddTextData(textData)
	assert.NoError(t, err)

	retrievedTextData, exists := storage.textDataList["text1"]
	assert.True(t, exists)
	assert.Equal(t, textData, retrievedTextData)
}

func TestMemoryStorage_AddFileData_Success(t *testing.T) {
	storage := NewMemoryStorage()
	fileData := models.FileData{
		FileName: "testfile.txt",
	}
	fileData.UUID = "file1"
	err := storage.AddFileData(fileData)
	assert.NoError(t, err)

	retrievedFileData, exists := storage.fileDataList["file1"]
	assert.True(t, exists)
	assert.Equal(t, fileData, retrievedFileData)
}

func TestMemoryStorage_AddMetaDataList_ReplaceExisting(t *testing.T) {
	storage := NewMemoryStorage()
	metaData1 := models.MetaData{
		DataUUID: "meta1",
		MetaName: "Test Meta 1",
	}

	err := storage.AddMetaDataList(metaData1)
	assert.NoError(t, err)

	metaData2 := models.MetaData{
		DataUUID: "meta1",
		MetaName: "Test Meta 2",
	}

	err = storage.AddMetaDataList(metaData2)
	assert.NoError(t, err)

	var retrievedMetaData models.MetaData
	for _, item := range storage.metaDataList {
		if item.DataUUID == "meta1" {
			retrievedMetaData = item
			break
		}
	}

	assert.Equal(t, metaData2, retrievedMetaData)
}
