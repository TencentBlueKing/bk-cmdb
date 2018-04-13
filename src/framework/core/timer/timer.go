package timer

import (
	"configcenter/src/framework/core/timer/regular"
	"configcenter/src/framework/core/timer/transaction"
	"container/list"
	"context"
)

// workerTransaction worker transaction list
//type workerTransaction map[WorkerKey]list.List

type timeLineWorkerTransaction map[WorkerKey]list.List

type workerTransaction struct {
	workerKey WorkerKey
}

// timer implements the Timer interface
type timer struct {
	regularTimer regular.Regular
}

// Regular get the regular timer. multiple calls return the same instance
func (cli *timer) Regular() regular.Regular {

	return cli.regularTimer
}

// CreateTransaction create a new transaction. each calls returns a new instance.
func (cli *timer) CreateTransaction(ctx context.Context, key WorkerKey) transaction.Transaction {
	tran := transaction.New(ctx)
	return tran
}

// DestoryTransaction destory the transaction
func (cli *timer) DestoryTransaction(tran transaction.Transaction) {

}
