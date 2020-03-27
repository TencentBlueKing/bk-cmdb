// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package session

import (
	"errors"

	"github.com/rentiansheng/bk_bsonbson"
	"github.com/rentiansheng/bk_bson/bson/primitive"
	"github.com/rentiansheng/bk_bsonmongo/readconcern"
	"github.com/rentiansheng/bk_bsonmongo/readpref"
	"github.com/rentiansheng/bk_bsonmongo/writeconcern"
	"github.com/rentiansheng/bk_bsonx/mongo/driver/uuid"
)

// ErrSessionEnded is returned when a client session is used after a call to endSession().
var ErrSessionEnded = errors.New("ended session was used")

// ErrNoTransactStarted is returned if a transaction operation is called when no transaction has started.
var ErrNoTransactStarted = errors.New("no transaction started")

// ErrTransactInProgress is returned if startTransaction() is called when a transaction is in progress.
var ErrTransactInProgress = errors.New("transaction already in progress")

// ErrAbortAfterCommit is returned when abort is called after a commit.
var ErrAbortAfterCommit = errors.New("cannot call abortTransaction after calling commitTransaction")

// ErrAbortTwice is returned if abort is called after transaction is already aborted.
var ErrAbortTwice = errors.New("cannot call abortTransaction twice")

// ErrCommitAfterAbort is returned if commit is called after an abort.
var ErrCommitAfterAbort = errors.New("cannot call commitTransaction after calling abortTransaction")

// ErrUnackWCUnsupported is returned if an unacknowledged write concern is supported for a transaciton.
var ErrUnackWCUnsupported = errors.New("transactions do not support unacknowledged write concerns")

// Type describes the type of the session
type Type uint8

// These constants are the valid types for a client session.
const (
	Explicit Type = iota
	Implicit
)

// State indicates the state of the FSM.
type state uint8

// Client Session states
const (
	None state = iota
	Starting
	InProgress
	Committed
	Aborted
)

// Client is a session for clients to run commands.
type Client struct {
	*Server
	ClientID       uuid.UUID
	ClusterTime    bson.Raw
	Consistent     bool // causal consistency
	OperationTime  *primitive.Timestamp
	SessionType    Type
	Terminated     bool
	RetryingCommit bool
	Committing     bool
	Aborting       bool
	RetryWrite     bool

	// options for the current transaction
	// most recently set by transactionopt
	CurrentRc *readconcern.ReadConcern
	CurrentRp *readpref.ReadPref
	CurrentWc *writeconcern.WriteConcern

	// default transaction options
	transactionRc *readconcern.ReadConcern
	transactionRp *readpref.ReadPref
	transactionWc *writeconcern.WriteConcern

	pool  *Pool
	state state
}

func getClusterTime(clusterTime bson.Raw) (uint32, uint32) {
	if clusterTime == nil {
		return 0, 0
	}

	clusterTimeVal, err := clusterTime.LookupErr("$clusterTime")
	if err != nil {
		return 0, 0
	}

	timestampVal, err := bson.Raw(clusterTimeVal.Value).LookupErr("clusterTime")
	if err != nil {
		return 0, 0
	}

	return timestampVal.Timestamp()
}

// MaxClusterTime compares 2 clusterTime documents and returns the document representing the highest cluster time.
func MaxClusterTime(ct1, ct2 bson.Raw) bson.Raw {
	epoch1, ord1 := getClusterTime(ct1)
	epoch2, ord2 := getClusterTime(ct2)

	if epoch1 > epoch2 {
		return ct1
	} else if epoch1 < epoch2 {
		return ct2
	} else if ord1 > ord2 {
		return ct1
	} else if ord1 < ord2 {
		return ct2
	}

	return ct1
}

// NewClientSession creates a Client.
func NewClientSession(pool *Pool, clientID uuid.UUID, sessionType Type, opts ...*ClientOptions) (*Client, error) {
	c := &Client{
		Consistent:  true, // set default
		ClientID:    clientID,
		SessionType: sessionType,
		pool:        pool,
	}

	mergedOpts := mergeClientOptions(opts...)
	if mergedOpts.CausalConsistency != nil {
		c.Consistent = *mergedOpts.CausalConsistency
	}
	if mergedOpts.DefaultReadPreference != nil {
		c.transactionRp = mergedOpts.DefaultReadPreference
	}
	if mergedOpts.DefaultReadConcern != nil {
		c.transactionRc = mergedOpts.DefaultReadConcern
	}
	if mergedOpts.DefaultWriteConcern != nil {
		c.transactionWc = mergedOpts.DefaultWriteConcern
	}

	servSess, err := pool.GetSession()
	if err != nil {
		return nil, err
	}

	c.Server = servSess

	return c, nil
}

// AdvanceClusterTime updates the session's cluster time.
func (c *Client) AdvanceClusterTime(clusterTime bson.Raw) error {
	if c.Terminated {
		return ErrSessionEnded
	}
	c.ClusterTime = MaxClusterTime(c.ClusterTime, clusterTime)
	return nil
}

// AdvanceOperationTime updates the session's operation time.
func (c *Client) AdvanceOperationTime(opTime *primitive.Timestamp) error {
	if c.Terminated {
		return ErrSessionEnded
	}

	if c.OperationTime == nil {
		c.OperationTime = opTime
		return nil
	}

	if opTime.T > c.OperationTime.T {
		c.OperationTime = opTime
	} else if (opTime.T == c.OperationTime.T) && (opTime.I > c.OperationTime.I) {
		c.OperationTime = opTime
	}

	return nil
}

// UpdateUseTime updates the session's last used time.
// Must be called whenver this session is used to send a command to the server.
func (c *Client) UpdateUseTime() error {
	if c.Terminated {
		return ErrSessionEnded
	}
	c.updateUseTime()
	return nil
}

// EndSession ends the session.
func (c *Client) EndSession() {
	if c.Terminated {
		return
	}

	c.Terminated = true
	c.pool.ReturnSession(c.Server)

	return
}

// TransactionInProgress returns true if the client session is in an active transaction.
func (c *Client) TransactionInProgress() bool {
	return c.state == InProgress
}

// TransactionStarting returns true if the client session is starting a transaction.
func (c *Client) TransactionStarting() bool {
	return c.state == Starting
}

// TransactionRunning returns true if the client session has started the transaction
// and it hasn't been committed or aborted
func (c *Client) TransactionRunning() bool {
	return c != nil && (c.state == Starting || c.state == InProgress)
}

// TransactionCommitted returns true of the client session just committed a transaciton.
func (c *Client) TransactionCommitted() bool {
	return c.state == Committed
}

// CheckStartTransaction checks to see if allowed to start transaction and returns
// an error if not allowed
func (c *Client) CheckStartTransaction() error {
	if c.state == InProgress || c.state == Starting {
		return ErrTransactInProgress
	}
	return nil
}

// StartTransaction initializes the transaction options and advances the state machine.
// It does not contact the server to start the transaction.
func (c *Client) StartTransaction(opts *TransactionOptions) error {
	err := c.CheckStartTransaction()
	if err != nil {
		return err
	}

	c.IncrementTxnNumber()
	c.RetryingCommit = false

	if opts != nil {
		c.CurrentRc = opts.ReadConcern
		c.CurrentRp = opts.ReadPreference
		c.CurrentWc = opts.WriteConcern
	}

	if c.CurrentRc == nil {
		c.CurrentRc = c.transactionRc
	}

	if c.CurrentRp == nil {
		c.CurrentRp = c.transactionRp
	}

	if c.CurrentWc == nil {
		c.CurrentWc = c.transactionWc
	}

	if !writeconcern.AckWrite(c.CurrentWc) {
		c.clearTransactionOpts()
		return ErrUnackWCUnsupported
	}

	c.state = Starting
	return nil
}

// CheckCommitTransaction checks to see if allowed to commit transaction and returns
// an error if not allowed.
func (c *Client) CheckCommitTransaction() error {
	if c.state == None {
		return ErrNoTransactStarted
	} else if c.state == Aborted {
		return ErrCommitAfterAbort
	}
	return nil
}

// CommitTransaction updates the state for a successfully committed transaction and returns
// an error if not permissible.  It does not actually perform the commit.
func (c *Client) CommitTransaction() error {
	err := c.CheckCommitTransaction()
	if err != nil {
		return err
	}
	c.state = Committed
	return nil
}

// CheckAbortTransaction checks to see if allowed to abort transaction and returns
// an error if not allowed.
func (c *Client) CheckAbortTransaction() error {
	if c.state == None {
		return ErrNoTransactStarted
	} else if c.state == Committed {
		return ErrAbortAfterCommit
	} else if c.state == Aborted {
		return ErrAbortTwice
	}
	return nil
}

// AbortTransaction updates the state for a successfully committed transaction and returns
// an error if not permissible.  It does not actually perform the abort.
func (c *Client) AbortTransaction() error {
	err := c.CheckAbortTransaction()
	if err != nil {
		return err
	}
	c.state = Aborted
	c.clearTransactionOpts()
	return nil
}

// ApplyCommand advances the state machine upon command execution.
func (c *Client) ApplyCommand() {
	if c.Committing {
		// Do not change state if committing after already committed
		return
	}
	if c.state == Starting {
		c.state = InProgress
	} else if c.state == Committed || c.state == Aborted {
		c.clearTransactionOpts()
		c.state = None
	}
}

func (c *Client) clearTransactionOpts() {
	c.RetryingCommit = false
	c.Aborting = false
	c.Committing = false
	c.CurrentWc = nil
	c.CurrentRp = nil
	c.CurrentRc = nil
}
