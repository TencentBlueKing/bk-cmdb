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
	"github.com/mongodb/mongo-go-driver/core/option"
	"github.com/mongodb/mongo-go-driver/core/result"
	"github.com/mongodb/mongo-go-driver/core/wiremessage"
	"github.com/mongodb/mongo-go-driver/core/writeconcern"
)

// CreateIndexes represents the createIndexes command.
//
// The createIndexes command creates indexes for a namespace.
type CreateIndexes struct {
	NS           Namespace
	Indexes      *bson.Array
	Opts         []option.CreateIndexesOptioner
	WriteConcern *writeconcern.WriteConcern

	result result.CreateIndexes
	err    error
}

// Encode will encode this command into a wire message for the given server description.
func (ci *CreateIndexes) Encode(desc description.SelectedServer) (wiremessage.WireMessage, error) {
	cmd, err := ci.encode(desc)
	if err != nil {
		return nil, err
	}

	return cmd.Encode(desc)
}

func (ci *CreateIndexes) encode(desc description.SelectedServer) (*Write, error) {
	cmd := bson.NewDocument(
		bson.EC.String("createIndexes", ci.NS.Collection),
		bson.EC.Array("indexes", ci.Indexes),
	)

	for _, opt := range ci.Opts {
		if opt == nil {
			continue
		}
		err := opt.Option(cmd)
		if err != nil {
			return nil, err
		}
	}

	return &Write{
		DB:           ci.NS.DB,
		Command:      cmd,
		WriteConcern: ci.WriteConcern,
	}, nil
}

// Decode will decode the wire message using the provided server description. Errors during decoding
// are deferred until either the Result or Err methods are called.
func (ci *CreateIndexes) Decode(desc description.SelectedServer, wm wiremessage.WireMessage) *CreateIndexes {
	rdr, err := (&Write{}).Decode(desc, wm).Result()
	if err != nil {
		ci.err = err
		return ci
	}

	return ci.decode(desc, rdr)
}

func (ci *CreateIndexes) decode(desc description.SelectedServer, rdr bson.Reader) *CreateIndexes {
	ci.err = bson.Unmarshal(rdr, &ci.result)
	return ci
}

// Result returns the result of a decoded wire message and server description.
func (ci *CreateIndexes) Result() (result.CreateIndexes, error) {
	if ci.err != nil {
		return result.CreateIndexes{}, ci.err
	}
	return ci.result, nil
}

// Err returns the error set on this command.
func (ci *CreateIndexes) Err() error { return ci.err }

// RoundTrip handles the execution of this command using the provided wiremessage.ReadWriter.
func (ci *CreateIndexes) RoundTrip(ctx context.Context, desc description.SelectedServer, rw wiremessage.ReadWriter) (result.CreateIndexes, error) {
	cmd, err := ci.encode(desc)
	if err != nil {
		return result.CreateIndexes{}, err
	}

	rdr, err := cmd.RoundTrip(ctx, desc, rw)
	if err != nil {
		return result.CreateIndexes{}, err
	}

	return ci.decode(desc, rdr).Result()
}
