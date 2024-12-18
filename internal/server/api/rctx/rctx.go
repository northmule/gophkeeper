package rctx

type key int

// Наименования контекста
const (
	// UserCtxKey объект с пользователем
	UserCtxKey key = iota
	// TransactionCtxKey транзакция в рамках запроса
	TransactionCtxKey
)

const (
	MapKeyUserUUID = "user_uuid"
)
