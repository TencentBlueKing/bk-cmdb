// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"context"

	"time"

	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/mongodb/mongo-go-driver/x/bsonx/bsoncore"
	"github.com/mongodb/mongo-go-driver/x/mongo/driver/session"
	"github.com/mongodb/mongo-go-driver/x/mongo/driver/topology"
	"github.com/mongodb/mongo-go-driver/x/mongo/driver/uuid"
	"github.com/mongodb/mongo-go-driver/x/network/command"
	"github.com/mongodb/mongo-go-driver/x/network/connection"
	"github.com/mongodb/mongo-go-driver/x/network/description"
)

// ListIndexes handles the full cycle dispatch and execution of a listIndexes command against the provided
// topology.
func ListIndexes(
	ctx context.Context,
	cmd command.ListIndexes,
	topo *topology.Topology,
	selector description.ServerSelector,
	clientID uuid.UUID,
	pool *session.Pool,
	opts ...*options.ListIndexesOptions,
) (*BatchCursor, error) {

	ss, err := topo.SelectServer(ctx, selector)
	if err != nil {
		return nil, err
	}

	conn, err := ss.Connection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if ss.Description().WireVersion.Max < 3 {
		return legacyListIndexes(ctx, cmd, ss, conn, opts...)
	}

	lio := options.MergeListIndexesOptions(opts...)
	if lio.BatchSize != nil {
		elem := bsonx.Elem{"batchSize", bsonx.Int32(*lio.BatchSize)}
		cmd.Opts = append(cmd.Opts, elem)
		cmd.CursorOpts = append(cmd.CursorOpts, elem)
	}
	if lio.MaxTime != nil {
		cmd.Opts = append(cmd.Opts, bsonx.Elem{"maxTimeMS", bsonx.Int64(int64(*lio.MaxTime / time.Millisecond))})
	}

	// If no explicit session and deployment supports sessions, start implicit session.
	if cmd.Session == nil && topo.SupportsSessions() {
		cmd.Session, err = session.NewClientSession(pool, clientID, session.Implicit)
		if err != nil {
			return nil, err
		}
	}

	res, err := cmd.RoundTrip(ctx, ss.Description(), conn)
	if err != nil {
		closeImplicitSession(cmd.Session)
		return nil, err
	}

	return NewBatchCursor(bsoncore.Document(res), cmd.Session, cmd.Clock, ss.Server, cmd.CursorOpts...)
}

func legacyListIndexes(
	ctx context.Context,
	cmd command.ListIndexes,
	ss *topology.SelectedServer,
	conn connection.Connection,
	opts ...*options.ListIndexesOptions,
) (*BatchCursor, error) {
	lio := options.MergeListIndexesOptions(opts...)
	ns := cmd.NS.DB + "." + cmd.NS.Collection

	findCmd := command.Find{
		NS: command.NewNamespace(cmd.NS.DB, "system.indexes"),
		Filter: bsonx.Doc{
			{"ns", bsonx.String(ns)},
		},
	}

	findOpts := options.Find()
	if lio.BatchSize != nil {
		findOpts.SetBatchSize(*lio.BatchSize)
	}
	if lio.MaxTime != nil {
		findOpts.SetMaxTime(*lio.MaxTime)
	}

	return legacyFind(ctx, findCmd, nil, ss, conn, findOpts)
}
