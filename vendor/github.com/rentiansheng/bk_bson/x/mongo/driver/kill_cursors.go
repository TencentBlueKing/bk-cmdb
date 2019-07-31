// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"context"
	"bk_bson/x/network/connection"
	"bk_bson/x/network/wiremessage"

	"bk_bson/x/mongo/driver/topology"
	"bk_bson/x/network/command"
	"bk_bson/x/network/result"
)

// KillCursors handles the full cycle dispatch and execution of an aggregate command against the provided
// topology.
func KillCursors(
	ctx context.Context,
	ns command.Namespace,
	server *topology.Server,
	cursorID int64,
) (result.KillCursors, error) {
	desc := server.SelectedDescription()
	conn, err := server.Connection(ctx)
	if err != nil {
		return result.KillCursors{}, err
	}
	defer conn.Close()

	if desc.WireVersion.Max < 4 {
		return result.KillCursors{}, legacyKillCursors(ctx, ns, cursorID, conn)
	}

	cmd := command.KillCursors{
		NS:  ns,
		IDs: []int64{cursorID},
	}

	return cmd.RoundTrip(ctx, desc, conn)
}

func legacyKillCursors(ctx context.Context, ns command.Namespace, cursorID int64, conn connection.Connection) error {
	kc := wiremessage.KillCursors{
		NumberOfCursorIDs: 1,
		CursorIDs:         []int64{cursorID},
		CollectionName:    ns.Collection,
		DatabaseName:      ns.DB,
	}

	return conn.WriteWireMessage(ctx, kc)
}
