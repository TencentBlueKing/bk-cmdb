// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"

	"github.com/mongodb/mongo-go-driver/bson"
)

// Cursor instances iterate a stream of documents. Each document is
// decoded into the result according to the rules of the bson package.
//
// A typical usage of the Cursor interface would be:
//
//		var cur Cursor
//		ctx := context.Background()
//		defer cur.Close(ctx)
//
// 		for cur.Next(ctx) {
//			elem := &bson.D{}
//			if err := cur.Decode(elem); err != nil {
// 				log.Fatal(err)
// 			}
//
// 			// do something with elem....
//		}
//
// 		if err := cur.Err(); err != nil {
//			log.Fatal(err)
//		}
//
type Cursor interface {
	// NOTE: Whenever ops.Cursor changes, this must be changed to match it.

	// Get the ID of the cursor.
	ID() int64

	// Get the next result from the cursor.
	// Returns true if there were no errors and there is a next result.
	Next(context.Context) bool

	Decode(interface{}) error

	DecodeBytes() (bson.Raw, error)

	// Returns the error status of the cursor
	Err() error

	// Close the cursor.
	Close(context.Context) error
}
