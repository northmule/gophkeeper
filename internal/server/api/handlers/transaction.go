package handlers

import (
	"context"

	"net/http"

	"github.com/northmule/gophkeeper/internal/server/api/rctx"
	"github.com/northmule/gophkeeper/internal/server/logger"
	"github.com/northmule/gophkeeper/internal/server/storage"
)

type TransactionHandler struct {
	db  storage.DBQuery
	log *logger.Logger
}

func NewTransactionHandler(db storage.DBQuery, log *logger.Logger) *TransactionHandler {
	instance := &TransactionHandler{
		db:  db,
		log: log,
	}

	return instance
}

func (th *TransactionHandler) Transaction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var err error

		transaction, err := storage.NewTransaction(th.db)
		th.log.Info("The beginning of the transaction")
		if err != nil {
			th.log.Errorf("The transaction is not open: %s", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer func() {
			if r := recover(); r != nil {
				th.log.Info("An application error has occurred. Rollback the transaction.", r)
				if err = transaction.Rollback(); err != nil {
					th.log.Errorf("Commit request error: %s", err)
				}
			}
			transaction = nil
		}()
		ctx := context.WithValue(req.Context(), rctx.TransactionCtxKey, transaction)
		req = req.WithContext(ctx)

		next.ServeHTTP(res, req)

		if len(transaction.Error()) > 0 {
			th.log.Info("Errors occurred during the execution of the transaction", transaction.Error())
			if err = transaction.Rollback(); err != nil {
				th.log.Errorf("Commit request error: %s", err)

			}
			th.log.Info("The transaction has been rolled back, the data has not been changed")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = transaction.Commit()
		if err != nil {
			th.log.Errorf("Commit request error: %s", err)
			res.WriteHeader(http.StatusInternalServerError)
		}
		th.log.Info("The transaction is completed. All changes are saved.")
	})
}
