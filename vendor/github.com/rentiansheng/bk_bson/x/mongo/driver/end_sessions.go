// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"context"

	"bk_bson/x/mongo/driver/topology"
	"bk_bson/x/network/command"
	"bk_bson/x/network/description"
	"bk_bson/x/network/result"
)

// EndSessions handles the full cycle dispatch and execution of an endSessions command against the provided
// topology.
func EndSessions(
	ctx context.Context,
	cmd command.EndSessions,
	topo *topology.Topology,
	selector description.ServerSelector,
) ([]result.EndSessions, []error) {

	ss, err := topo.SelectServer(ctx, selector)
	if err != nil {
		return nil, []error{err}
	}

	desc := ss.Description()
	conn, err := ss.Connection(ctx)
	if err != nil {
		return nil, []error{err}
	}
	defer conn.Close()

	return cmd.RoundTrip(ctx, desc, conn)
}
