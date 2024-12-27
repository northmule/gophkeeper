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

//func TestTransactionHandler_Transaction_NewTransactionError(t *testing.T) {
//	mockDB := new(MockDBQuery)
//	mockLogger, _ := logger.NewLogger("info")
//
//	mockDB.On("NewTransaction").Return(nil, fmt.Errorf("transaction error"))
//
//	handler := NewTransactionHandler(mockDB, mockLogger)
//
//	next := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
//		// This should not be called
//	})
//
//	req, _ := http.NewRequest("GET", "/test", nil)
//	rr := httptest.NewRecorder()
//
//	handler.Transaction(next).ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusInternalServerError, rr.Code)
//	mockDB.AssertExpectations(t)
//}
//
//func TestTransactionHandler_Transaction_TransactionError(t *testing.T) {
//	mockDB := new(MockDBQuery)
//	mockTransaction := new(MockTransaction)
//	mockLogger, _ := logger.NewLogger("info")
//
//	mockDB.On("NewTransaction").Return(mockTransaction, nil)
//	mockTransaction.On("Rollback").Return(nil)
//	mockTransaction.On("Error").Return([]error{fmt.Errorf("transaction error")})
//
//	handler := NewTransactionHandler(mockDB, mockLogger)
//
//	next := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
//		// Simulate error during processing
//		panic("test error")
//	})
//
//	req, _ := http.NewRequest("GET", "/test", nil)
//	rr := httptest.NewRecorder()
//
//	handler.Transaction(next).ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusInternalServerError, rr.Code)
//	mockDB.AssertExpectations(t)
//	mockTransaction.AssertExpectations(t)
//}
//
//func TestTransactionHandler_Transaction_CommitError(t *testing.T) {
//	mockDB := new(MockDBQuery)
//	mockTransaction := new(MockTransaction)
//	mockLogger, _ := logger.NewLogger("info")
//
//	mockDB.On("NewTransaction").Return(mockTransaction, nil)
//	mockTransaction.On("Commit").Return(fmt.Errorf("commit error"))
//	mockTransaction.On("Rollback").Return(nil)
//	mockTransaction.On("Error").Return([]error{})
//
//	handler := NewTransactionHandler(mockDB, mockLogger)
//
//	next := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
//		// Simulate successful processing
//	})
//
//	req, _ := http.NewRequest("GET", "/test", nil)
//	rr := httptest.NewRecorder()
//
//	handler.Transaction(next).ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusInternalServerError, rr.Code)
//	mockDB.AssertExpectations(t)
//	mockTransaction.AssertExpectations(t)
//}
//
//func TestTransactionHandler_Transaction_RollbackError(t *testing.T) {
//	mockDB := new(MockDBQuery)
//	mockTransaction := new(MockTransaction)
//	mockLogger, _ := logger.NewLogger("info")
//
//	mockDB.On("NewTransaction").Return(mockTransaction, nil)
//	mockTransaction.On("Rollback").Return(fmt.Errorf("rollback error"))
//	mockTransaction.On("Error").Return([]error{fmt.Errorf("transaction error")})
//
//	handler := NewTransactionHandler(mockDB, mockLogger)
//
//	next := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
//		// Simulate error during processing
//		panic("test error")
//	})
//
//	req, _ := http.NewRequest("GET", "/test", nil)
//	rr := httptest.NewRecorder()
//
//	handler.Transaction(next).ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusInternalServerError, rr.Code)
//	mockDB.AssertExpectations(t)
//	mockTransaction.AssertExpectations(t)
//}
//
//func TestTransactionHandler_Transaction_PanicDuringProcessing(t *testing.T) {
//	mockDB := new(MockDBQuery)
//	mockTransaction := new(MockTransaction)
//	mockLogger, _ := logger.NewLogger("info")
//
//	mockDB.On("NewTransaction").Return(mockTransaction, nil)
//	mockTransaction.On("Rollback").Return(nil)
//	mockTransaction.On("Error").Return([]error{})
//
//	handler := NewTransactionHandler(mockDB, mockLogger)
//
//	next := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
//		// Simulate panic during processing
//		panic("test error")
//	})
//
//	req, _ := http.NewRequest("GET", "/test", nil)
//	rr := httptest.NewRecorder()
//
//	handler.Transaction(next).ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusInternalServerError, rr.Code)
//	mockDB.AssertExpectations(t)
//	mockTransaction.AssertExpectations(t)
//}
//
//func TestTransactionHandler_Transaction_ContextPropagation(t *testing.T) {
//	mockDB := new(MockDBQuery)
//	mockTransaction := new(MockTransaction)
//	mockLogger, _ := logger.NewLogger("info")
//
//	mockDB.On("NewTransaction").Return(mockTransaction, nil)
//	mockTransaction.On("Commit").Return(nil)
//	mockTransaction.On("Rollback").Return(nil)
//	mockTransaction.On("Error").Return([]error{})
//
//	handler := NewTransactionHandler(mockDB, mockLogger)
//
//	next := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
//		// Check if the transaction is in the context
//		tx := rctx.GetTransaction(req.Context())
//		assert.NotNil(t, tx)
//	})
//
//	req, _ := http.NewRequest("GET", "/test", nil)
//	rr := httptest.NewRecorder()
//
//	handler.Transaction(next).ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusOK, rr.Code)
//	mockDB.AssertExpectations(t)
//	mockTransaction.AssertExpectations(t)
//}
