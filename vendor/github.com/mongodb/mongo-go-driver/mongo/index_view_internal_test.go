// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/stretchr/testify/require"
)

var seed = time.Now().UnixNano()

type index struct {
	Key  map[string]int
	NS   string
	Name string
}

func getIndexableCollection(t *testing.T) (string, *Collection) {
	atomic.AddInt64(&seed, 1)
	rand.Seed(atomic.LoadInt64(&seed))

	client := createTestClient(t)
	db := client.Database(t.Name())

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	require.NoError(t, err)

	dbName := hex.EncodeToString(randomBytes)

	err = db.RunCommand(
		context.Background(),
		bsonx.Doc{{"create", bsonx.String(dbName)}},
	).Err()
	require.NoError(t, err)

	return dbName, db.Collection(dbName)
}

func TestIndexView_List(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip()
	}

	dbName, coll := getIndexableCollection(t)
	expectedNS := fmt.Sprintf("%s.%s", t.Name(), dbName)
	indexView := coll.Indexes()

	cursor, err := indexView.List(context.Background())
	require.NoError(t, err)

	var found bool

	for cursor.Next(context.Background()) {
		var idx index
		err := cursor.Decode(&idx)
		require.NoError(t, err)

		require.Equal(t, expectedNS, idx.NS)

		if idx.Name == "_id_" {
			require.Len(t, idx.Key, 1)
			require.Equal(t, 1, idx.Key["_id"])
			found = true
		}
	}
	require.NoError(t, cursor.Err())
	require.True(t, found)
}

func TestIndexView_CreateOne(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip()
	}

	dbName, coll := getIndexableCollection(t)
	expectedNS := fmt.Sprintf("%s.%s", t.Name(), dbName)
	indexView := coll.Indexes()

	indexName, err := indexView.CreateOne(
		context.Background(),
		IndexModel{
			Keys: bsonx.Doc{{"foo", bsonx.Int32(-1)}},
		},
	)
	require.NoError(t, err)

	cursor, err := indexView.List(context.Background())
	require.NoError(t, err)

	var found bool

	for cursor.Next(context.Background()) {
		var idx index
		err := cursor.Decode(&idx)
		require.NoError(t, err)

		require.Equal(t, expectedNS, idx.NS)

		if idx.Name == indexName {
			require.Len(t, idx.Key, 1)
			require.Equal(t, -1, idx.Key["foo"])
			found = true
		}
	}
	require.NoError(t, cursor.Err())
	require.True(t, found)
}

func TestIndexView_CreateOneWithNameOption(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip()
	}

	dbName, coll := getIndexableCollection(t)
	expectedNS := fmt.Sprintf("%s.%s", t.Name(), dbName)
	indexView := coll.Indexes()

	indexName, err := indexView.CreateOne(
		context.Background(),
		IndexModel{
			Keys:    bsonx.Doc{{"foo", bsonx.Int32(-1)}},
			Options: options.Index().SetName("testname"),
		},
	)
	require.NoError(t, err)
	require.Equal(t, "testname", indexName)

	cursor, err := indexView.List(context.Background())
	require.NoError(t, err)

	var found bool

	for cursor.Next(context.Background()) {
		var idx index
		err := cursor.Decode(&idx)
		require.NoError(t, err)

		require.Equal(t, expectedNS, idx.NS)

		if idx.Name == indexName {
			require.Len(t, idx.Key, 1)
			require.Equal(t, -1, idx.Key["foo"])
			found = true
		}
	}
	require.NoError(t, cursor.Err())
	require.True(t, found)
}

// Omits collation option because it's incompatible with version option
func TestIndexView_CreateOneWithAllOptions(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip()
	}

	_, coll := getIndexableCollection(t)
	indexView := coll.Indexes()

	_, err := indexView.CreateOne(
		context.Background(),
		IndexModel{
			Keys: bsonx.Doc{{"foo", bsonx.String("text")}},
			Options: options.Index().
				SetBackground(false).
				SetExpireAfterSeconds(10).
				SetName("a").
				SetSparse(false).
				SetUnique(false).
				SetVersion(1).
				SetDefaultLanguage("english").
				SetLanguageOverride("english").
				SetTextVersion(1).
				SetWeights(bsonx.Doc{}).
				SetSphereVersion(1).
				SetBits(2).
				SetMax(10).
				SetMin(1).
				SetBucketSize(1).
				SetPartialFilterExpression(bsonx.Doc{}).
				SetStorageEngine(bsonx.Doc{
					{"wiredTiger", bsonx.Document(bsonx.Doc{
						{"configString", bsonx.String("block_compressor=zlib")},
					})},
				}),
		},
	)
	require.NoError(t, err)
}

func TestIndexView_CreateOneWithCollationOption(t *testing.T) {
	skipIfBelow34(t, createTestDatabase(t, nil)) // collation invalid for server versions < 3.4
	t.Parallel()

	if testing.Short() {
		t.Skip()
	}

	_, coll := getIndexableCollection(t)
	indexView := coll.Indexes()

	_, err := indexView.CreateOne(
		context.Background(),
		IndexModel{
			Keys: bsonx.Doc{{"bar", bsonx.String("text")}},
			Options: options.Index().SetCollation(&options.Collation{
				Locale: "simple",
			}),
		},
	)
	require.NoError(t, err)
}

func TestIndexView_CreateOneWithNilKeys(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip()
	}

	_, coll := getIndexableCollection(t)
	indexView := coll.Indexes()

	_, err := indexView.CreateOne(
		context.Background(),
		IndexModel{
			Keys: nil,
		},
	)
	require.Error(t, err)
}

func TestIndexView_CreateMany(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip()
	}

	dbName, coll := getIndexableCollection(t)
	expectedNS := fmt.Sprintf("%s.%s", t.Name(), dbName)
	indexView := coll.Indexes()

	indexNames, err := indexView.CreateMany(
		context.Background(),
		[]IndexModel{
			{
				Keys: bsonx.Doc{{"foo", bsonx.Int32(-1)}},
			},
			{
				Keys: bsonx.Doc{
					{"bar", bsonx.Int32(1)},
					{"baz", bsonx.Int32(-1)},
				},
			},
		},
	)
	require.NoError(t, err)

	require.Len(t, indexNames, 2)

	fooName := indexNames[0]
	barBazName := indexNames[1]

	cursor, err := indexView.List(context.Background())
	require.NoError(t, err)

	fooFound := false
	barBazFound := false

	for cursor.Next(context.Background()) {
		var idx index
		err := cursor.Decode(&idx)
		require.NoError(t, err)

		require.Equal(t, expectedNS, idx.NS)

		if idx.Name == fooName {
			require.Len(t, idx.Key, 1)
			require.Equal(t, -1, idx.Key["foo"])
			fooFound = true
		}

		if idx.Name == barBazName {
			require.Len(t, idx.Key, 2)
			require.Equal(t, 1, idx.Key["bar"])
			require.Equal(t, -1, idx.Key["baz"])
			barBazFound = true
		}
	}
	require.NoError(t, cursor.Err())
	require.True(t, fooFound)
	require.True(t, barBazFound)
}

func TestIndexView_DropOne(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip()
	}

	dbName, coll := getIndexableCollection(t)
	expectedNS := fmt.Sprintf("%s.%s", t.Name(), dbName)
	indexView := coll.Indexes()

	indexNames, err := indexView.CreateMany(
		context.Background(),
		[]IndexModel{
			{
				Keys: bsonx.Doc{{"foo", bsonx.Int32(-1)}},
			},
			{
				Keys: bsonx.Doc{
					{"bar", bsonx.Int32(1)},
					{"baz", bsonx.Int32(-1)},
				},
			},
		},
	)
	require.NoError(t, err)

	require.Len(t, indexNames, 2)

	_, err = indexView.DropOne(
		context.Background(),
		indexNames[1],
	)
	require.NoError(t, err)

	cursor, err := indexView.List(context.Background())
	require.NoError(t, err)

	var idx index

	for cursor.Next(context.Background()) {
		err := cursor.Decode(&idx)
		require.NoError(t, err)

		require.Equal(t, expectedNS, idx.NS)
		require.NotEqual(t, indexNames[1], idx.Name)
	}
	require.NoError(t, cursor.Err())
}

func TestIndexView_DropAll(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip()
	}

	dbName, coll := getIndexableCollection(t)
	expectedNS := fmt.Sprintf("%s.%s", t.Name(), dbName)
	indexView := coll.Indexes()

	indexNames, err := indexView.CreateMany(
		context.Background(),
		[]IndexModel{
			{
				Keys: bsonx.Doc{{"foo", bsonx.Int32(-1)}},
			},
			{
				Keys: bsonx.Doc{
					{"bar", bsonx.Int32(1)},
					{"baz", bsonx.Int32(-1)},
				},
			},
		},
	)
	require.NoError(t, err)

	require.Len(t, indexNames, 2)

	_, err = indexView.DropAll(
		context.Background(),
	)
	require.NoError(t, err)

	cursor, err := indexView.List(context.Background())
	require.NoError(t, err)

	var idx index

	for cursor.Next(context.Background()) {
		err := cursor.Decode(&idx)
		require.NoError(t, err)

		require.Equal(t, expectedNS, idx.NS)
		require.NotEqual(t, indexNames[0], idx.Name)
		require.NotEqual(t, indexNames[1], idx.Name)
	}
	require.NoError(t, cursor.Err())
}

func TestIndexView_CreateIndexesOptioner(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip()
	}

	dbName, coll := getIndexableCollection(t)
	expectedNS := fmt.Sprintf("%s.%s", t.Name(), dbName)
	indexView := coll.Indexes()

	opts := options.CreateIndexes().SetMaxTime(1000)
	indexNames, err := indexView.CreateMany(
		context.Background(),
		[]IndexModel{
			{
				Keys: bsonx.Doc{{"foo", bsonx.Int32(-1)}},
			},
			{
				Keys: bsonx.Doc{
					{"bar", bsonx.Int32(1)},
					{"baz", bsonx.Int32(-1)},
				},
			},
		},
		opts,
	)
	require.NoError(t, err)
	require.NoError(t, err)

	require.Len(t, indexNames, 2)

	fooName := indexNames[0]
	barBazName := indexNames[1]

	cursor, err := indexView.List(context.Background())
	require.NoError(t, err)

	fooFound := false
	barBazFound := false

	for cursor.Next(context.Background()) {
		var idx index
		err := cursor.Decode(&idx)
		require.NoError(t, err)

		require.Equal(t, expectedNS, idx.NS)

		if idx.Name == fooName {
			require.Len(t, idx.Key, 1)
			require.Equal(t, -1, idx.Key["foo"])
			fooFound = true
		}

		if idx.Name == barBazName {
			require.Len(t, idx.Key, 2)
			require.Equal(t, 1, idx.Key["bar"])
			require.Equal(t, -1, idx.Key["baz"])
			barBazFound = true
		}
	}
	require.NoError(t, cursor.Err())
	require.True(t, fooFound)
	require.True(t, barBazFound)
	defer func() {
		_, err := indexView.DropAll(context.Background())
		require.NoError(t, err)
	}()
}

func TestIndexView_DropIndexesOptioner(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip()
	}

	dbName, coll := getIndexableCollection(t)
	expectedNS := fmt.Sprintf("%s.%s", t.Name(), dbName)
	indexView := coll.Indexes()

	opts := options.DropIndexes().SetMaxTime(1000)
	indexNames, err := indexView.CreateMany(
		context.Background(),
		[]IndexModel{
			{
				Keys: bsonx.Doc{{"foo", bsonx.Int32(-1)}},
			},
			{
				Keys: bsonx.Doc{
					{"bar", bsonx.Int32(1)},
					{"baz", bsonx.Int32(-1)},
				},
			},
		},
	)
	require.NoError(t, err)

	require.Len(t, indexNames, 2)

	_, err = indexView.DropAll(
		context.Background(),
		opts,
	)
	require.NoError(t, err)

	cursor, err := indexView.List(context.Background())
	require.NoError(t, err)

	for cursor.Next(context.Background()) {
		var idx index
		err := cursor.Decode(&idx)
		require.NoError(t, err)
		require.Equal(t, expectedNS, idx.NS)
		require.NotEqual(t, indexNames[0], idx.Name)
		require.NotEqual(t, indexNames[1], idx.Name)
	}
	require.NoError(t, cursor.Err())
}
