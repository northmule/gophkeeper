package handlers

//func TestTransactionHandler_Transaction_Success(t *testing.T) {
//	mockDB := new(appMock.MockDBQuery)
//	mockTxDBQuery := new(appMock.MockTxDBQuery)
//	mockLogger, _ := logger.NewLogger("info")
//
//	mockTxDBQuery.On("Commit").Return(nil)
//	mockTxDBQuery.On("Rollback").Return(nil)
//	mockTxDBQuery.On("Error").Return([]error{})
//	mockDB.On("Begin").Return(mockTxDBQuery, nil)
//	//transaction, _ := storage.NewTransaction(mockDB)
//
//	handler := NewTransactionHandler(mockDB, mockLogger)
//
//	next := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
//		//
//	})
//
//	req := httptest.NewRequest("POST", "/decrypt", nil)
//	req.Header.Set("Content-Type", "application/json")
//	rr := httptest.NewRecorder()
//
//	handler.Transaction(next).ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusOK, rr.Code)
//	mockDB.AssertExpectations(t)
//}
