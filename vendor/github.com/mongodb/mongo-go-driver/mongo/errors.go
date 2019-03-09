// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/x/mongo/driver"
	"github.com/mongodb/mongo-go-driver/x/mongo/driver/topology"
	"github.com/mongodb/mongo-go-driver/x/network/command"
	"github.com/mongodb/mongo-go-driver/x/network/result"
)

// ErrUnacknowledgedWrite is returned from functions that have an unacknowledged
// write concern.
var ErrUnacknowledgedWrite = errors.New("unacknowledged write")

// ErrClientDisconnected is returned when a user attempts to call a method on a
// disconnected client
var ErrClientDisconnected = errors.New("client is disconnected")

// ErrNilDocument is returned when a user attempts to pass a nil document or filter
// to a function where the field is required.
var ErrNilDocument = errors.New("document is nil")

// ErrEmptySlice is returned when a user attempts to pass an empty slice as input
// to a function wehere the field is required.
var ErrEmptySlice = errors.New("must provide at least one element in input slice")

func replaceTopologyErr(err error) error {
	if err == topology.ErrTopologyClosed {
		return ErrClientDisconnected
	}
	return err
}

// WriteError is a non-write concern failure that occurred as a result of a write
// operation.
type WriteError struct {
	Index   int
	Code    int
	Message string
}

func (we WriteError) Error() string { return we.Message }

// WriteErrors is a group of non-write concern failures that occurred as a result
// of a write operation.
type WriteErrors []WriteError

func (we WriteErrors) Error() string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "write errors: [")
	for idx, err := range we {
		if idx != 0 {
			fmt.Fprintf(&buf, ", ")
		}
		fmt.Fprintf(&buf, "{%s}", err)
	}
	fmt.Fprint(&buf, "]")
	return buf.String()
}

func writeErrorsFromResult(rwes []result.WriteError) WriteErrors {
	wes := make(WriteErrors, 0, len(rwes))
	for _, err := range rwes {
		wes = append(wes, WriteError{Index: err.Index, Code: err.Code, Message: err.ErrMsg})
	}
	return wes
}

// WriteConcernError is a write concern failure that occurred as a result of a
// write operation.
type WriteConcernError struct {
	Code    int
	Message string
	Details bson.Raw
}

func (wce WriteConcernError) Error() string { return wce.Message }

func convertBulkWriteErrors(errors []driver.BulkWriteError) []BulkWriteError {
	bwErrors := make([]BulkWriteError, 0, len(errors))
	for _, err := range errors {
		bwErrors = append(bwErrors, BulkWriteError{
			WriteError{
				Index:   err.Index,
				Code:    err.Code,
				Message: err.ErrMsg,
			},
			dispatchToMongoModel(err.Model),
		})
	}

	return bwErrors
}

func convertWriteConcernError(wce *result.WriteConcernError) *WriteConcernError {
	if wce == nil {
		return nil
	}

	return &WriteConcernError{Code: wce.Code, Message: wce.ErrMsg, Details: wce.ErrInfo}
}

// BulkWriteError is an error for one operation in a bulk write.
type BulkWriteError struct {
	WriteError
	Request WriteModel
}

func (bwe BulkWriteError) Error() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "{%s}", bwe.WriteError)
	return buf.String()
}

// BulkWriteException is an error for a bulk write operation.
type BulkWriteException struct {
	WriteConcernError *WriteConcernError
	WriteErrors       []BulkWriteError
}

func (bwe BulkWriteException) Error() string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "bulk write error: [")
	fmt.Fprintf(&buf, "{%s}, ", bwe.WriteErrors)
	fmt.Fprintf(&buf, "{%s}]", bwe.WriteConcernError)
	return buf.String()
}

// returnResult is used to determine if a function calling processWriteError should return
// the result or return nil. Since the processWriteError function is used by many different
// methods, both *One and *Many, we need a way to differentiate if the method should return
// the result and the error.
type returnResult int

const (
	rrNone returnResult = 1 << iota // None means do not return the result ever.
	rrOne                           // One means return the result if this was called by a *One method.
	rrMany                          // Many means return the result is this was called by a *Many method.

	rrAll returnResult = rrOne | rrMany // All means always return the result.
)

// processWriteError handles processing the result of a write operation. If the retrunResult matches
// the calling method's type, it should return the result object in addition to the error.
// This function will wrap the errors from other packages and return them as errors from this package.
//
// WriteConcernError will be returned over WriteErrors if both are present.
func processWriteError(wce *result.WriteConcernError, wes []result.WriteError, err error) (returnResult, error) {
	switch {
	case err == command.ErrUnacknowledgedWrite:
		return rrAll, ErrUnacknowledgedWrite
	case err != nil:
		return rrNone, replaceTopologyErr(err)
	case wce != nil:
		return rrMany, WriteConcernError{Code: wce.Code, Message: wce.ErrMsg}
	case len(wes) > 0:
		return rrMany, writeErrorsFromResult(wes)
	default:
		return rrAll, nil
	}
}
