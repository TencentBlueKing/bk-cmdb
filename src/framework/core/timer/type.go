package timer

import (
	"configcenter/src/framework/core/timer/regular"
	"configcenter/src/framework/core/timer/transaction"
	"context"
)

// WorkerKey same to the api.WorkerKey
type WorkerKey string

// Timer is the interface  that defines the creation of the method of the normal Timer and the transactional Timer.
type Timer interface {

	// Regular return the regular timer. multiple calls return the same instance
	Regular() regular.Regular

	// CreateTransaction create a new transaction. each calls returns a new instance.
	CreateTransaction(ctx context.Context, key WorkerKey) transaction.Transaction

	// DestoryTransaction destory the transaction
	DestoryTransaction(tran transaction.Transaction)
}
