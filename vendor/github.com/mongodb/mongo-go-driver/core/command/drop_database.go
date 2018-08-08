// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package command

import (
	"context"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/core/description"
	"github.com/mongodb/mongo-go-driver/core/wiremessage"
	"github.com/mongodb/mongo-go-driver/core/writeconcern"
)

// DropDatabase represents the DropDatabase command.
//
// The DropDatabases command drops database.
type DropDatabase struct {
	DB           string
	WriteConcern *writeconcern.WriteConcern

	result bson.Reader
	err    error
}

// Encode will encode this command into a wire message for the given server description.
func (dd *DropDatabase) Encode(desc description.SelectedServer) (wiremessage.WireMessage, error) {
	cmd, err := dd.encode(desc)
	if err != nil {
		return nil, err
	}

	return cmd.Encode(desc)
}

func (dd *DropDatabase) encode(desc description.SelectedServer) (*Write, error) {
	cmd := bson.NewDocument(
		bson.EC.Int32("dropDatabase", 1),
	)

	return &Write{
		DB:           dd.DB,
		Command:      cmd,
		WriteConcern: dd.WriteConcern,
	}, nil
}

// Decode will decode the wire message using the provided server description. Errors during decoding
// are deferred until either the Result or Err methods are called.
func (dd *DropDatabase) Decode(desc description.SelectedServer, wm wiremessage.WireMessage) *DropDatabase {
	dd.result, dd.err = (&Write{}).Decode(desc, wm).Result()
	return dd
}

// Result returns the result of a decoded wire message and server description.
func (dd *DropDatabase) Result() (bson.Reader, error) {
	if dd.err != nil {
		return nil, dd.err
	}
	return dd.result, nil
}

// Err returns the error set on this command.
func (dd *DropDatabase) Err() error { return dd.err }

// RoundTrip handles the execution of this command using the provided wiremessage.ReadWriter.
func (dd *DropDatabase) RoundTrip(ctx context.Context, desc description.SelectedServer, rw wiremessage.ReadWriter) (bson.Reader, error) {
	cmd, err := dd.encode(desc)
	if err != nil {
		return nil, err
	}

	dd.result, err = cmd.RoundTrip(ctx, desc, rw)
	if err != nil {
		return nil, err
	}

	return dd.Result()
}
