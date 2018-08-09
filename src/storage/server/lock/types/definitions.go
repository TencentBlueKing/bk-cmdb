package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/rs/xid"
)

var (
	LockPermissionDenied = errors.New("Permission denied")
	LockNotFound         = errors.New("lock not found")
)

var (
	LockIDPrefix = "bkcc"
)

// lcok struct
type Lock struct {
	//  id of this transaction
	TxnID string `json:"txnID"`

	// sub  id of this lock
	SubTxnID string `json:"subTxnID"`

	// lock name is used to define the resources that this lock should be locked
	LockName string `json:"lockName"`

	// timeout means that the time of the client can bear to wait for the lock is locked.
	Timeout time.Duration `json:"timeout"`

	Createtime time.Time `json:"createTime"`
}

// LockResult lock check result
type LockResult struct {
	// the sub txn ID of the txn.
	SubTxnID string `json:"subTxnID"`

	// whether the resources has been locked or not
	Locked bool `json:"locked"`

	// first lock resources TxnID
	LockSubTxnID string `json:"lockSubTxnID"`
}

// GetID   lock tag  ID
func GetID(prefix string) string {
	id := xid.New()
	return fmt.Sprintf("%s-%s", prefix, id.String())
}
