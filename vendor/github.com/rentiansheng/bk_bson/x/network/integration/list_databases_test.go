// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package integration

import (
	"context"
	"testing"

	"github.com/rentiansheng/bk_bsoninternal/testutil"
	"github.com/rentiansheng/bk_bsonmongo/writeconcern"
	"github.com/rentiansheng/bk_bson/x/bsonx"
	"github.com/rentiansheng/bk_bsonx/network/command"
	"github.com/rentiansheng/bk_bsonx/network/description"
)

func TestListDatabases(t *testing.T) {
	noerr := func(t *testing.T, err error) {
		// t.Helper()
		if err != nil {
			t.Errorf("Unepexted error: %v", err)
			t.FailNow()
		}
	}
	server, err := testutil.Topology(t).SelectServer(context.Background(), description.WriteSelector())
	noerr(t, err)
	conn, err := server.Connection(context.Background())
	noerr(t, err)

	wc := writeconcern.New(writeconcern.WMajority())
	testutil.AutoDropCollection(t)
	testutil.AutoInsertDocs(t, wc, bsonx.Doc{{"_id", bsonx.Int32(1)}})

	res, err := (&command.ListDatabases{}).RoundTrip(context.Background(), server.SelectedDescription(), conn)
	noerr(t, err)
	var found bool
	for _, db := range res.Databases {
		if db.Name == dbName {
			found = true
		}
	}
	if !found {
		t.Error("Should have found database in listDatabases result, but didn't.")
	}
}
