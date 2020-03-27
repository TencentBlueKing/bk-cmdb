// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"context"

	"github.com/rentiansheng/bk_bsonmongo/options"
	"github.com/rentiansheng/bk_bson/x/bsonx"

	"github.com/rentiansheng/bk_bsonmongo/writeconcern"
	"github.com/rentiansheng/bk_bsonx/mongo/driver/session"
	"github.com/rentiansheng/bk_bsonx/mongo/driver/topology"
	"github.com/rentiansheng/bk_bsonx/mongo/driver/uuid"
	"github.com/rentiansheng/bk_bsonx/network/command"
	"github.com/rentiansheng/bk_bsonx/network/description"
	"github.com/rentiansheng/bk_bsonx/network/result"
)

// Delete handles the full cycle dispatch and execution of a delete command against the provided
// topology.
func Delete(
	ctx context.Context,
	cmd command.Delete,
	topo *topology.Topology,
	selector description.ServerSelector,
	clientID uuid.UUID,
	pool *session.Pool,
	retryWrite bool,
	opts ...*options.DeleteOptions,
) (result.Delete, error) {

	ss, err := topo.SelectServer(ctx, selector)
	if err != nil {
		return result.Delete{}, err
	}

	// If no explicit session and deployment supports sessions, start implicit session.
	if cmd.Session == nil && topo.SupportsSessions() && writeconcern.AckWrite(cmd.WriteConcern) {
		cmd.Session, err = session.NewClientSession(pool, clientID, session.Implicit)
		if err != nil {
			return result.Delete{}, err
		}
		defer cmd.Session.EndSession()
	}

	deleteOpts := options.MergeDeleteOptions(opts...)
	if deleteOpts.Collation != nil {
		if ss.Description().WireVersion.Max < 5 {
			return result.Delete{}, ErrCollation
		}
		cmd.Opts = append(cmd.Opts, bsonx.Elem{"collation", bsonx.Document(deleteOpts.Collation.ToDocument())})
	}

	// Execute in a single trip if retry writes not supported, or retry not enabled
	if !retrySupported(topo, ss.Description(), cmd.Session, cmd.WriteConcern) || !retryWrite {
		if cmd.Session != nil {
			cmd.Session.RetryWrite = false // explicitly set to false to prevent encoding transaction number
		}
		return delete(ctx, cmd, ss, nil)
	}

	cmd.Session.RetryWrite = retryWrite
	cmd.Session.IncrementTxnNumber()

	res, originalErr := delete(ctx, cmd, ss, nil)

	// Retry if appropriate
	if cerr, ok := originalErr.(command.Error); ok && cerr.Retryable() ||
		res.WriteConcernError != nil && command.IsWriteConcernErrorRetryable(res.WriteConcernError) {
		ss, err := topo.SelectServer(ctx, selector)

		// Return original error if server selection fails or new server does not support retryable writes
		if err != nil || !retrySupported(topo, ss.Description(), cmd.Session, cmd.WriteConcern) {
			return res, originalErr
		}

		return delete(ctx, cmd, ss, cerr)
	}
	return res, originalErr
}

func delete(
	ctx context.Context,
	cmd command.Delete,
	ss *topology.SelectedServer,
	oldErr error,
) (result.Delete, error) {
	desc := ss.Description()

	conn, err := ss.Connection(ctx)
	if err != nil {
		if oldErr != nil {
			return result.Delete{}, oldErr
		}
		return result.Delete{}, err
	}

	if !writeconcern.AckWrite(cmd.WriteConcern) {
		go func() {
			defer func() { _ = recover() }()
			defer conn.Close()

			_, _ = cmd.RoundTrip(ctx, desc, conn)
		}()

		return result.Delete{}, command.ErrUnacknowledgedWrite
	}
	defer conn.Close()

	return cmd.RoundTrip(ctx, desc, conn)
}
