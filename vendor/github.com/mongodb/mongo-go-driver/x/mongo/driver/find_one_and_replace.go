// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"context"

	"time"

	"github.com/mongodb/mongo-go-driver/bson/bsoncodec"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/mongo/writeconcern"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/mongodb/mongo-go-driver/x/mongo/driver/session"
	"github.com/mongodb/mongo-go-driver/x/mongo/driver/topology"
	"github.com/mongodb/mongo-go-driver/x/mongo/driver/uuid"
	"github.com/mongodb/mongo-go-driver/x/network/command"
	"github.com/mongodb/mongo-go-driver/x/network/description"
	"github.com/mongodb/mongo-go-driver/x/network/result"
)

// FindOneAndReplace handles the full cycle dispatch and execution of a FindOneAndReplace command against the provided
// topology.
func FindOneAndReplace(
	ctx context.Context,
	cmd command.FindOneAndReplace,
	topo *topology.Topology,
	selector description.ServerSelector,
	clientID uuid.UUID,
	pool *session.Pool,
	retryWrite bool,
	registry *bsoncodec.Registry,
	opts ...*options.FindOneAndReplaceOptions,
) (result.FindAndModify, error) {

	ss, err := topo.SelectServer(ctx, selector)
	if err != nil {
		return result.FindAndModify{}, err
	}

	// If no explicit session and deployment supports sessions, start implicit session.
	if cmd.Session == nil && topo.SupportsSessions() {
		cmd.Session, err = session.NewClientSession(pool, clientID, session.Implicit)
		if err != nil {
			return result.FindAndModify{}, err
		}
		defer cmd.Session.EndSession()
	}

	ro := options.MergeFindOneAndReplaceOptions(opts...)
	if ro.BypassDocumentValidation != nil {
		cmd.Opts = append(cmd.Opts, bsonx.Elem{"byapssDocumentValidation", bsonx.Boolean(*ro.BypassDocumentValidation)})
	}
	if ro.Collation != nil {
		if ss.Description().WireVersion.Max < 5 {
			return result.FindAndModify{}, ErrCollation
		}
		cmd.Opts = append(cmd.Opts, bsonx.Elem{"collation", bsonx.Document(ro.Collation.ToDocument())})
	}
	if ro.MaxTime != nil {
		cmd.Opts = append(cmd.Opts, bsonx.Elem{"maxTimeMS", bsonx.Int64(int64(*ro.MaxTime / time.Millisecond))})
	}
	if ro.Projection != nil {
		maxElem, err := interfaceToElement("fields", ro.Projection, registry)
		if err != nil {
			return result.FindAndModify{}, err
		}

		cmd.Opts = append(cmd.Opts, maxElem)
	}
	if ro.ReturnDocument != nil {
		cmd.Opts = append(cmd.Opts, bsonx.Elem{"new", bsonx.Boolean(*ro.ReturnDocument == options.After)})
	}
	if ro.Sort != nil {
		sortElem, err := interfaceToElement("sort", ro.Sort, registry)
		if err != nil {
			return result.FindAndModify{}, err
		}

		cmd.Opts = append(cmd.Opts, sortElem)
	}
	if ro.Upsert != nil {
		cmd.Opts = append(cmd.Opts, bsonx.Elem{"upsert", bsonx.Boolean(*ro.Upsert)})
	}

	// Execute in a single trip if retry writes not supported, or retry not enabled
	if !retrySupported(topo, ss.Description(), cmd.Session, cmd.WriteConcern) || !retryWrite {
		if cmd.Session != nil {
			cmd.Session.RetryWrite = false // explicitly set to false to prevent encoding transaction number
		}
		return findOneAndReplace(ctx, cmd, ss, nil)
	}

	cmd.Session.RetryWrite = retryWrite
	cmd.Session.IncrementTxnNumber()

	res, originalErr := findOneAndReplace(ctx, cmd, ss, nil)

	// Retry if appropriate
	if cerr, ok := originalErr.(command.Error); ok && cerr.Retryable() {
		ss, err := topo.SelectServer(ctx, selector)

		// Return original error if server selection fails or new server does not support retryable writes
		if err != nil || !retrySupported(topo, ss.Description(), cmd.Session, cmd.WriteConcern) {
			return result.FindAndModify{}, originalErr
		}

		return findOneAndReplace(ctx, cmd, ss, cerr)
	}

	return res, originalErr
}

func findOneAndReplace(
	ctx context.Context,
	cmd command.FindOneAndReplace,
	ss *topology.SelectedServer,
	oldErr error,
) (result.FindAndModify, error) {
	desc := ss.Description()
	conn, err := ss.Connection(ctx)
	if err != nil {
		if oldErr != nil {
			return result.FindAndModify{}, oldErr
		}
		return result.FindAndModify{}, err
	}

	if !writeconcern.AckWrite(cmd.WriteConcern) {
		go func() {
			defer func() { _ = recover() }()
			defer conn.Close()

			_, _ = cmd.RoundTrip(ctx, desc, conn)
		}()

		return result.FindAndModify{}, command.ErrUnacknowledgedWrite
	}
	defer conn.Close()

	return cmd.RoundTrip(ctx, desc, conn)
}
