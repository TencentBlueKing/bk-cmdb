// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// NOTE: Any time this file is modified, a WEBSITE ticket should be opened to sync the changes with
// the "What is MongoDB" webpage, which the example was originally added to as part of WEBSITE-5148.

package documentation_examples

import (
	"context"
	"fmt"
	"io/ioutil"
	logger "log"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func requireCursorLength(t *testing.T, cursor *mongo.Cursor, length int) {
	i := 0
	for cursor.Next(context.Background()) {
		i++
	}

	require.NoError(t, cursor.Err())
	require.Equal(t, i, length)
}

func containsKey(doc bson.Raw, key ...string) bool {
	_, err := doc.LookupErr(key...)
	if err != nil {
		return false
	}
	return true
}

func parseDate(t *testing.T, dateString string) time.Time {
	rfc3339MilliLayout := "2006-01-02T15:04:05.999Z07:00" // layout defined with Go reference time
	parsedDate, err := time.Parse(rfc3339MilliLayout, dateString)

	require.NoError(t, err)
	return parsedDate
}

// InsertExamples contains examples for insert operations.
func InsertExamples(t *testing.T, db *mongo.Database) {
	coll := db.Collection("inventory_insert")

	err := coll.Drop(context.Background())
	require.NoError(t, err)

	{
		// Start Example 1

		result, err := coll.InsertOne(
			context.Background(),
			bson.D{
				{"item", "canvas"},
				{"qty", 100},
				{"tags", bson.A{"cotton"}},
				{"size", bson.D{
					{"h", 28},
					{"w", 35.5},
					{"uom", "cm"},
				}},
			})

		// End Example 1

		require.NoError(t, err)
		require.NotNil(t, result.InsertedID)
	}

	{
		// Start Example 2

		cursor, err := coll.Find(
			context.Background(),
			bson.D{{"item", "canvas"}},
		)

		// End Example 2

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)

	}

	{
		// Start Example 3

		result, err := coll.InsertMany(
			context.Background(),
			[]interface{}{
				bson.D{
					{"item", "journal"},
					{"qty", int32(25)},
					{"tags", bson.A{"blank", "red"}},
					{"size", bson.D{
						{"h", 14},
						{"w", 21},
						{"uom", "cm"},
					}},
				},
				bson.D{
					{"item", "mat"},
					{"qty", int32(25)},
					{"tags", bson.A{"gray"}},
					{"size", bson.D{
						{"h", 27.9},
						{"w", 35.5},
						{"uom", "cm"},
					}},
				},
				bson.D{
					{"item", "mousepad"},
					{"qty", 25},
					{"tags", bson.A{"gel", "blue"}},
					{"size", bson.D{
						{"h", 19},
						{"w", 22.85},
						{"uom", "cm"},
					}},
				},
			})

		// End Example 3

		require.NoError(t, err)
		require.Len(t, result.InsertedIDs, 3)
	}
}

// QueryToplevelFieldsExamples contains examples for querying top-level fields.
func QueryToplevelFieldsExamples(t *testing.T, db *mongo.Database) {
	coll := db.Collection("inventory_query_top")

	err := coll.Drop(context.Background())
	require.NoError(t, err)

	{
		// Start Example 6

		docs := []interface{}{
			bson.D{
				{"item", "journal"},
				{"qty", 25},
				{"size", bson.D{
					{"h", 14},
					{"w", 21},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "notebook"},
				{"qty", 50},
				{"size", bson.D{
					{"h", 8.5},
					{"w", 11},
					{"uom", "in"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "paper"},
				{"qty", 100},
				{"size", bson.D{
					{"h", 8.5},
					{"w", 11},
					{"uom", "in"},
				}},
				{"status", "D"},
			},
			bson.D{
				{"item", "planner"},
				{"qty", 75},
				{"size", bson.D{
					{"h", 22.85},
					{"w", 30},
					{"uom", "cm"},
				}},
				{"status", "D"},
			},
			bson.D{
				{"item", "postcard"},
				{"qty", 45},
				{"size", bson.D{
					{"h", 10},
					{"w", 15.25},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
		}

		result, err := coll.InsertMany(context.Background(), docs)

		// End Example 6

		require.NoError(t, err)
		require.Len(t, result.InsertedIDs, 5)
	}

	{
		// Start Example 7

		cursor, err := coll.Find(
			context.Background(),
			bson.D{},
		)

		// End Example 7

		require.NoError(t, err)
		requireCursorLength(t, cursor, 5)
	}

	{
		// Start Example 9

		cursor, err := coll.Find(
			context.Background(),
			bson.D{{"status", "D"}},
		)

		// End Example 9

		require.NoError(t, err)
		requireCursorLength(t, cursor, 2)
	}

	{
		// Start Example 10

		cursor, err := coll.Find(
			context.Background(),
			bson.D{{"status", bson.D{{"$in", bson.A{"A", "D"}}}}})

		// End Example 10

		require.NoError(t, err)
		requireCursorLength(t, cursor, 5)
	}

	{
		// Start Example 11

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"status", "A"},
				{"qty", bson.D{{"$lt", 30}}},
			})

		// End Example 11

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)
	}

	{
		// Start Example 12

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"$or",
					bson.A{
						bson.D{{"status", "A"}},
						bson.D{{"qty", bson.D{{"$lt", 30}}}},
					}},
			})

		// End Example 12

		require.NoError(t, err)
		requireCursorLength(t, cursor, 3)
	}

	{
		// Start Example 13

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"status", "A"},
				{"$or", bson.A{
					bson.D{{"qty", bson.D{{"$lt", 30}}}},
					bson.D{{"item", primitive.Regex{Pattern: "^p", Options: ""}}},
				}},
			})

		// End Example 13

		require.NoError(t, err)
		requireCursorLength(t, cursor, 2)
	}

}

// QueryEmbeddedDocumentsExamples contains examples for querying embedded document fields.
func QueryEmbeddedDocumentsExamples(t *testing.T, db *mongo.Database) {
	coll := db.Collection("inventory_query_embedded")

	err := coll.Drop(context.Background())
	require.NoError(t, err)

	{
		// Start Example 14

		docs := []interface{}{
			bson.D{
				{"item", "journal"},
				{"qty", 25},
				{"size", bson.D{
					{"h", 14},
					{"w", 21},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "notebook"},
				{"qty", 50},
				{"size", bson.D{
					{"h", 8.5},
					{"w", 11},
					{"uom", "in"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "paper"},
				{"qty", 100},
				{"size", bson.D{
					{"h", 8.5},
					{"w", 11},
					{"uom", "in"},
				}},
				{"status", "D"},
			},
			bson.D{
				{"item", "planner"},
				{"qty", 75},
				{"size", bson.D{
					{"h", 22.85},
					{"w", 30},
					{"uom", "cm"},
				}},
				{"status", "D"},
			},
			bson.D{
				{"item", "postcard"},
				{"qty", 45},
				{"size", bson.D{
					{"h", 10},
					{"w", 15.25},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
		}

		result, err := coll.InsertMany(context.Background(), docs)

		// End Example 14

		require.NoError(t, err)
		require.Len(t, result.InsertedIDs, 5)
	}

	{
		// Start Example 15

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"size", bson.D{
					{"h", 14},
					{"w", 21},
					{"uom", "cm"},
				}},
			})

		// End Example 15

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)
	}

	{
		// Start Example 16

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"size", bson.D{
					{"w", 21},
					{"h", 14},
					{"uom", "cm"},
				}},
			})

		// End Example 16

		require.NoError(t, err)
		requireCursorLength(t, cursor, 0)
	}

	{
		// Start Example 17

		cursor, err := coll.Find(
			context.Background(),
			bson.D{{"size.uom", "in"}},
		)

		// End Example 17

		require.NoError(t, err)
		requireCursorLength(t, cursor, 2)
	}

	{
		// Start Example 18

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"size.h", bson.D{
					{"$lt", 15},
				}},
			})

		// End Example 18

		require.NoError(t, err)
		requireCursorLength(t, cursor, 4)
	}

	{
		// Start Example 19

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"size.h", bson.D{
					{"$lt", 15},
				}},
				{"size.uom", "in"},
				{"status", "D"},
			})

		// End Example 19

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)
	}

}

// QueryArraysExamples contains examples for querying array fields.
func QueryArraysExamples(t *testing.T, db *mongo.Database) {
	coll := db.Collection("inventory_query_array")

	err := coll.Drop(context.Background())
	require.NoError(t, err)

	{
		// Start Example 20

		docs := []interface{}{
			bson.D{
				{"item", "journal"},
				{"qty", 25},
				{"tags", bson.A{"blank", "red"}},
				{"dim_cm", bson.A{14, 21}},
			},
			bson.D{
				{"item", "notebook"},
				{"qty", 50},
				{"tags", bson.A{"red", "blank"}},
				{"dim_cm", bson.A{14, 21}},
			},
			bson.D{
				{"item", "paper"},
				{"qty", 100},
				{"tags", bson.A{"red", "blank", "plain"}},
				{"dim_cm", bson.A{14, 21}},
			},
			bson.D{
				{"item", "planner"},
				{"qty", 75},
				{"tags", bson.A{"blank", "red"}},
				{"dim_cm", bson.A{22.85, 30}},
			},
			bson.D{
				{"item", "postcard"},
				{"qty", 45},
				{"tags", bson.A{"blue"}},
				{"dim_cm", bson.A{10, 15.25}},
			},
		}

		result, err := coll.InsertMany(context.Background(), docs)

		// End Example 20

		require.NoError(t, err)
		require.Len(t, result.InsertedIDs, 5)
	}

	{
		// Start Example 21

		cursor, err := coll.Find(
			context.Background(),
			bson.D{{"tags", bson.A{"red", "blank"}}},
		)

		// End Example 21

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)
	}

	{
		// Start Example 22

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"tags", bson.D{{"$all", bson.A{"red", "blank"}}}},
			})

		// End Example 22

		require.NoError(t, err)
		requireCursorLength(t, cursor, 4)
	}

	{
		// Start Example 23

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"tags", "red"},
			})

		// End Example 23

		require.NoError(t, err)
		requireCursorLength(t, cursor, 4)
	}

	{
		// Start Example 24

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"dim_cm", bson.D{
					{"$gt", 25},
				}},
			})

		// End Example 24

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)
	}

	{
		// Start Example 25

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"dim_cm", bson.D{
					{"$gt", 15},
					{"$lt", 20},
				}},
			})

		// End Example 25

		require.NoError(t, err)
		requireCursorLength(t, cursor, 4)
	}

	{
		// Start Example 26

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"dim_cm", bson.D{
					{"$elemMatch", bson.D{
						{"$gt", 22},
						{"$lt", 30},
					}},
				}},
			})

		// End Example 26

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)
	}

	{
		// Start Example 27

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"dim_cm.1", bson.D{
					{"$gt", 25},
				}},
			})

		// End Example 27

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)
	}

	{
		// Start Example 28

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"tags", bson.D{
					{"$size", 3},
				}},
			})

		// End Example 28

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)
	}

}

// QueryArrayEmbeddedDocumentsExamples contains examples for querying fields with arrays and embedded documents.
func QueryArrayEmbeddedDocumentsExamples(t *testing.T, db *mongo.Database) {
	coll := db.Collection("inventory_query_array_embedded")

	err := coll.Drop(context.Background())
	require.NoError(t, err)

	{
		// Start Example 29

		docs := []interface{}{
			bson.D{
				{"item", "journal"},
				{"instock", bson.A{
					bson.D{
						{"warehouse", "A"},
						{"qty", 5},
					},
					bson.D{
						{"warehouse", "C"},
						{"qty", 15},
					},
				}},
			},
			bson.D{
				{"item", "notebook"},
				{"instock", bson.A{
					bson.D{
						{"warehouse", "C"},
						{"qty", 5},
					},
				}},
			},
			bson.D{
				{"item", "paper"},
				{"instock", bson.A{
					bson.D{
						{"warehouse", "A"},
						{"qty", 60},
					},
					bson.D{
						{"warehouse", "B"},
						{"qty", 15},
					},
				}},
			},
			bson.D{
				{"item", "planner"},
				{"instock", bson.A{
					bson.D{
						{"warehouse", "A"},
						{"qty", 40},
					},
					bson.D{
						{"warehouse", "B"},
						{"qty", 5},
					},
				}},
			},
			bson.D{
				{"item", "postcard"},
				{"instock", bson.A{
					bson.D{
						{"warehouse", "B"},
						{"qty", 15},
					},
					bson.D{
						{"warehouse", "C"},
						{"qty", 35},
					},
				}},
			},
		}

		result, err := coll.InsertMany(context.Background(), docs)

		// End Example 29

		require.NoError(t, err)
		require.Len(t, result.InsertedIDs, 5)
	}

	{
		// Start Example 30

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"instock", bson.D{
					{"warehouse", "A"},
					{"qty", 5},
				}},
			})

		// End Example 30

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)
	}

	{
		// Start Example 31

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"instock", bson.D{
					{"qty", 5},
					{"warehouse", "A"},
				}},
			})

		// End Example 31

		require.NoError(t, err)
		requireCursorLength(t, cursor, 0)
	}

	{
		// Start Example 32

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"instock.0.qty", bson.D{
					{"$lte", 20},
				}},
			})

		// End Example 32

		require.NoError(t, err)
		requireCursorLength(t, cursor, 3)
	}

	{
		// Start Example 33

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"instock.qty", bson.D{
					{"$lte", 20},
				}},
			})

		// End Example 33

		require.NoError(t, err)
		requireCursorLength(t, cursor, 5)
	}

	{
		// Start Example 34

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"instock", bson.D{
					{"$elemMatch", bson.D{
						{"qty", 5},
						{"warehouse", "A"},
					}},
				}},
			})

		// End Example 34

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)
	}

	{
		// Start Example 35

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"instock", bson.D{
					{"$elemMatch", bson.D{
						{"qty", bson.D{
							{"$gt", 10},
							{"$lte", 20},
						}},
					}},
				}},
			})

		// End Example 35

		require.NoError(t, err)
		requireCursorLength(t, cursor, 3)
	}

	{
		// Start Example 36

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"instock.qty", bson.D{
					{"$gt", 10},
					{"$lte", 20},
				}},
			})

		// End Example 36

		require.NoError(t, err)
		requireCursorLength(t, cursor, 4)
	}

	{
		// Start Example 37

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"instock.qty", 5},
				{"instock.warehouse", "A"},
			})

		// End Example 37

		require.NoError(t, err)
		requireCursorLength(t, cursor, 2)
	}
}

// QueryNullMissingFieldsExamples contains examples for querying fields that are null or missing.
func QueryNullMissingFieldsExamples(t *testing.T, db *mongo.Database) {
	coll := db.Collection("inventory_query_null_missing")

	err := coll.Drop(context.Background())
	require.NoError(t, err)

	{
		// Start Example 38

		docs := []interface{}{
			bson.D{
				{"_id", 1},
				{"item", nil},
			},
			bson.D{
				{"_id", 2},
			},
		}

		result, err := coll.InsertMany(context.Background(), docs)

		// End Example 38

		require.NoError(t, err)
		require.Len(t, result.InsertedIDs, 2)
	}

	{
		// Start Example 39

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"item", nil},
			})

		// End Example 39

		require.NoError(t, err)
		requireCursorLength(t, cursor, 2)
	}

	{
		// Start Example 40

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"item", bson.D{
					{"$type", 10},
				}},
			})

		// End Example 40

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)
	}

	{
		// Start Example 41

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"item", bson.D{
					{"$exists", false},
				}},
			})

		// End Example 41

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)
	}
}

// ProjectionExamples contains examples for specifying projections in find operations.
func ProjectionExamples(t *testing.T, db *mongo.Database) {
	coll := db.Collection("inventory_project")

	err := coll.Drop(context.Background())
	require.NoError(t, err)

	{
		// Start Example 42

		docs := []interface{}{
			bson.D{
				{"item", "journal"},
				{"status", "A"},
				{"size", bson.D{
					{"h", 14},
					{"w", 21},
					{"uom", "cm"},
				}},
				{"instock", bson.A{
					bson.D{
						{"warehouse", "A"},
						{"qty", 5},
					},
				}},
			},
			bson.D{
				{"item", "notebook"},
				{"status", "A"},
				{"size", bson.D{
					{"h", 8.5},
					{"w", 11},
					{"uom", "in"},
				}},
				{"instock", bson.A{
					bson.D{
						{"warehouse", "EC"},
						{"qty", 5},
					},
				}},
			},
			bson.D{
				{"item", "paper"},
				{"status", "D"},
				{"size", bson.D{
					{"h", 8.5},
					{"w", 11},
					{"uom", "in"},
				}},
				{"instock", bson.A{
					bson.D{
						{"warehouse", "A"},
						{"qty", 60},
					},
				}},
			},
			bson.D{
				{"item", "planner"},
				{"status", "D"},
				{"size", bson.D{
					{"h", 22.85},
					{"w", 30},
					{"uom", "cm"},
				}},
				{"instock", bson.A{
					bson.D{
						{"warehouse", "A"},
						{"qty", 40},
					},
				}},
			},
			bson.D{
				{"item", "postcard"},
				{"status", "A"},
				{"size", bson.D{
					{"h", 10},
					{"w", 15.25},
					{"uom", "cm"},
				}},
				{"instock", bson.A{
					bson.D{
						{"warehouse", "B"},
						{"qty", 15},
					},
					bson.D{
						{"warehouse", "EC"},
						{"qty", 35},
					},
				}},
			},
		}

		result, err := coll.InsertMany(context.Background(), docs)

		// End Example 42

		require.NoError(t, err)
		require.Len(t, result.InsertedIDs, 5)
	}

	{
		// Start Example 43

		cursor, err := coll.Find(
			context.Background(),
			bson.D{{"status", "A"}},
		)

		// End Example 43

		require.NoError(t, err)
		requireCursorLength(t, cursor, 3)
	}

	{
		// Start Example 44

		projection := bson.D{
			{"item", 1},
			{"status", 1},
		}

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"status", "A"},
			},
			options.Find().SetProjection(projection),
		)

		// End Example 44

		require.NoError(t, err)

		for cursor.Next(context.Background()) {
			doc := cursor.Current

			require.True(t, containsKey(doc, "_id"))
			require.True(t, containsKey(doc, "item"))
			require.True(t, containsKey(doc, "status"))
			require.False(t, containsKey(doc, "size"))
			require.False(t, containsKey(doc, "instock"))
		}

		require.NoError(t, cursor.Err())
	}

	{
		// Start Example 45

		projection := bson.D{
			{"item", 1},
			{"status", 1},
			{"_id", 0},
		}

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"status", "A"},
			},
			options.Find().SetProjection(projection),
		)

		// End Example 45

		require.NoError(t, err)

		for cursor.Next(context.Background()) {
			doc := cursor.Current

			require.False(t, containsKey(doc, "_id"))
			require.True(t, containsKey(doc, "item"))
			require.True(t, containsKey(doc, "status"))
			require.False(t, containsKey(doc, "size"))
			require.False(t, containsKey(doc, "instock"))
		}

		require.NoError(t, cursor.Err())
	}

	{
		// Start Example 46

		projection := bson.D{
			{"status", 0},
			{"instock", 0},
		}

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"status", "A"},
			},
			options.Find().SetProjection(projection),
		)

		// End Example 46

		require.NoError(t, err)

		for cursor.Next(context.Background()) {
			doc := cursor.Current

			require.True(t, containsKey(doc, "_id"))
			require.True(t, containsKey(doc, "item"))
			require.False(t, containsKey(doc, "status"))
			require.True(t, containsKey(doc, "size"))
			require.False(t, containsKey(doc, "instock"))
		}

		require.NoError(t, cursor.Err())
	}

	{
		// Start Example 47

		projection := bson.D{
			{"item", 1},
			{"status", 1},
			{"size.uom", 1},
		}

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"status", "A"},
			},
			options.Find().SetProjection(projection),
		)

		// End Example 47

		require.NoError(t, err)

		for cursor.Next(context.Background()) {
			doc := cursor.Current

			require.True(t, containsKey(doc, "_id"))
			require.True(t, containsKey(doc, "item"))
			require.True(t, containsKey(doc, "status"))
			require.True(t, containsKey(doc, "size"))
			require.False(t, containsKey(doc, "instock"))

			require.True(t, containsKey(doc, "size", "uom"))
			require.False(t, containsKey(doc, "size", "h"))
			require.False(t, containsKey(doc, "size", "w"))

		}

		require.NoError(t, cursor.Err())
	}

	{
		// Start Example 48

		projection := bson.D{
			{"size.uom", 0},
		}

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"status", "A"},
			},
			options.Find().SetProjection(projection),
		)

		// End Example 48

		require.NoError(t, err)

		for cursor.Next(context.Background()) {
			doc := cursor.Current

			require.True(t, containsKey(doc, "_id"))
			require.True(t, containsKey(doc, "item"))
			require.True(t, containsKey(doc, "status"))
			require.True(t, containsKey(doc, "size"))
			require.True(t, containsKey(doc, "instock"))

			require.False(t, containsKey(doc, "size", "uom"))
			require.True(t, containsKey(doc, "size", "h"))
			require.True(t, containsKey(doc, "size", "w"))

		}

		require.NoError(t, cursor.Err())
	}

	{
		// Start Example 49

		projection := bson.D{
			{"item", 1},
			{"status", 1},
			{"instock.qty", 1},
		}

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"status", "A"},
			},
			options.Find().SetProjection(projection),
		)

		// End Example 49

		require.NoError(t, err)

		for cursor.Next(context.Background()) {
			doc := cursor.Current

			require.True(t, containsKey(doc, "_id"))
			require.True(t, containsKey(doc, "item"))
			require.True(t, containsKey(doc, "status"))
			require.False(t, containsKey(doc, "size"))
			require.True(t, containsKey(doc, "instock"))

			instock, err := doc.LookupErr("instock")
			require.NoError(t, err)

			vals, err := instock.Array().Values()
			require.NoError(t, err)

			for _, val := range vals {
				require.Equal(t, bson.TypeEmbeddedDocument, val.Type)
				subdoc := val.Document()
				elems, err := subdoc.Elements()
				require.NoError(t, err)

				require.Equal(t, 1, len(elems))
				_, err = subdoc.LookupErr("qty")
				require.NoError(t, err)
			}
		}

		require.NoError(t, cursor.Err())
	}

	{
		// Start Example 50

		projection := bson.D{
			{"item", 1},
			{"status", 1},
			{"instock", bson.D{
				{"$slice", -1},
			}},
		}

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"status", "A"},
			},
			options.Find().SetProjection(projection),
		)

		// End Example 50

		require.NoError(t, err)

		for cursor.Next(context.Background()) {
			doc := cursor.Current

			require.True(t, containsKey(doc, "_id"))
			require.True(t, containsKey(doc, "item"))
			require.True(t, containsKey(doc, "status"))
			require.False(t, containsKey(doc, "size"))
			require.True(t, containsKey(doc, "instock"))

			instock, err := doc.LookupErr("instock")
			require.NoError(t, err)
			vals, err := instock.Array().Values()
			require.NoError(t, err)
			require.Equal(t, len(vals), 1)
		}

		require.NoError(t, cursor.Err())
	}
}

// UpdateExamples contains examples of update operations.
func UpdateExamples(t *testing.T, db *mongo.Database) {
	coll := db.Collection("inventory_update")

	err := coll.Drop(context.Background())
	require.NoError(t, err)

	{
		// Start Example 51

		docs := []interface{}{
			bson.D{
				{"item", "canvas"},
				{"qty", 100},
				{"size", bson.D{
					{"h", 28},
					{"w", 35.5},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "journal"},
				{"qty", 25},
				{"size", bson.D{
					{"h", 14},
					{"w", 21},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "mat"},
				{"qty", 85},
				{"size", bson.D{
					{"h", 27.9},
					{"w", 35.5},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "mousepad"},
				{"qty", 25},
				{"size", bson.D{
					{"h", 19},
					{"w", 22.85},
					{"uom", "in"},
				}},
				{"status", "P"},
			},
			bson.D{
				{"item", "notebook"},
				{"qty", 50},
				{"size", bson.D{
					{"h", 8.5},
					{"w", 11},
					{"uom", "in"},
				}},
				{"status", "P"},
			},
			bson.D{
				{"item", "paper"},
				{"qty", 100},
				{"size", bson.D{
					{"h", 8.5},
					{"w", 11},
					{"uom", "in"},
				}},
				{"status", "D"},
			},
			bson.D{
				{"item", "planner"},
				{"qty", 75},
				{"size", bson.D{
					{"h", 22.85},
					{"w", 30},
					{"uom", "cm"},
				}},
				{"status", "D"},
			},
			bson.D{
				{"item", "postcard"},
				{"qty", 45},
				{"size", bson.D{
					{"h", 10},
					{"w", 15.25},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "sketchbook"},
				{"qty", 80},
				{"size", bson.D{
					{"h", 14},
					{"w", 21},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "sketch pad"},
				{"qty", 95},
				{"size", bson.D{
					{"h", 22.85},
					{"w", 30.5},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
		}

		result, err := coll.InsertMany(context.Background(), docs)

		// End Example 51

		require.NoError(t, err)
		require.Len(t, result.InsertedIDs, 10)
	}

	{
		// Start Example 52

		result, err := coll.UpdateOne(
			context.Background(),
			bson.D{
				{"item", "paper"},
			},
			bson.D{
				{"$set", bson.D{
					{"size.uom", "cm"},
					{"status", "P"},
				}},
				{"$currentDate", bson.D{
					{"lastModified", true},
				}},
			},
		)

		// End Example 52

		require.NoError(t, err)
		require.Equal(t, int64(1), result.MatchedCount)
		require.Equal(t, int64(1), result.ModifiedCount)

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"item", "paper"},
			})

		require.NoError(t, err)

		for cursor.Next(context.Background()) {
			doc := cursor.Current

			uom, err := doc.LookupErr("size", "uom")
			require.NoError(t, err)
			require.Equal(t, uom.StringValue(), "cm")

			status, err := doc.LookupErr("status")
			require.NoError(t, err)
			require.Equal(t, status.StringValue(), "P")

			require.True(t, containsKey(doc, "lastModified"))
		}

		require.NoError(t, cursor.Err())
	}

	{
		// Start Example 53

		result, err := coll.UpdateMany(
			context.Background(),
			bson.D{
				{"qty", bson.D{
					{"$lt", 50},
				}},
			},
			bson.D{
				{"$set", bson.D{
					{"size.uom", "cm"},
					{"status", "P"},
				}},
				{"$currentDate", bson.D{
					{"lastModified", true},
				}},
			},
		)

		// End Example 53

		require.NoError(t, err)
		require.Equal(t, int64(3), result.MatchedCount)
		require.Equal(t, int64(3), result.ModifiedCount)

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"qty", bson.D{
					{"$lt", 50},
				}},
			})

		require.NoError(t, err)

		for cursor.Next(context.Background()) {
			doc := cursor.Current

			uom, err := doc.LookupErr("size", "uom")
			require.NoError(t, err)
			require.Equal(t, uom.StringValue(), "cm")

			status, err := doc.LookupErr("status")
			require.NoError(t, err)
			require.Equal(t, status.StringValue(), "P")

			require.True(t, containsKey(doc, "lastModified"))
		}

		require.NoError(t, cursor.Err())
	}

	{
		// Start Example 54

		result, err := coll.ReplaceOne(
			context.Background(),
			bson.D{
				{"item", "paper"},
			},
			bson.D{
				{"item", "paper"},
				{"instock", bson.A{
					bson.D{
						{"warehouse", "A"},
						{"qty", 60},
					},
					bson.D{
						{"warehouse", "B"},
						{"qty", 40},
					},
				}},
			},
		)

		// End Example 54

		require.NoError(t, err)
		require.Equal(t, int64(1), result.MatchedCount)
		require.Equal(t, int64(1), result.ModifiedCount)

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"item", "paper"},
			})

		require.NoError(t, err)

		for cursor.Next(context.Background()) {
			require.True(t, containsKey(cursor.Current, "_id"))
			require.True(t, containsKey(cursor.Current, "item"))
			require.True(t, containsKey(cursor.Current, "instock"))

			instock, err := cursor.Current.LookupErr("instock")
			require.NoError(t, err)
			vals, err := instock.Array().Values()
			require.NoError(t, err)
			require.Equal(t, len(vals), 2)

		}

		require.NoError(t, cursor.Err())
	}

}

// DeleteExamples contains examples of delete operations.
func DeleteExamples(t *testing.T, db *mongo.Database) {
	coll := db.Collection("inventory_delete")

	err := coll.Drop(context.Background())
	require.NoError(t, err)

	{
		// Start Example 55
		docs := []interface{}{
			bson.D{
				{"item", "journal"},
				{"qty", 25},
				{"size", bson.D{
					{"h", 14},
					{"w", 21},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "notebook"},
				{"qty", 50},
				{"size", bson.D{
					{"h", 8.5},
					{"w", 11},
					{"uom", "in"},
				}},
				{"status", "P"},
			},
			bson.D{
				{"item", "paper"},
				{"qty", 100},
				{"size", bson.D{
					{"h", 8.5},
					{"w", 11},
					{"uom", "in"},
				}},
				{"status", "D"},
			},
			bson.D{
				{"item", "planner"},
				{"qty", 75},
				{"size", bson.D{
					{"h", 22.85},
					{"w", 30},
					{"uom", "cm"},
				}},
				{"status", "D"},
			},
			bson.D{
				{"item", "postcard"},
				{"qty", 45},
				{"size", bson.D{
					{"h", 10},
					{"w", 15.25},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
		}

		result, err := coll.InsertMany(context.Background(), docs)

		// End Example 55

		require.NoError(t, err)
		require.Len(t, result.InsertedIDs, 5)
	}

	{
		// Start Example 57

		result, err := coll.DeleteMany(
			context.Background(),
			bson.D{
				{"status", "A"},
			},
		)

		// End Example 57

		require.NoError(t, err)
		require.Equal(t, int64(2), result.DeletedCount)
	}

	{
		// Start Example 58

		result, err := coll.DeleteOne(
			context.Background(),
			bson.D{
				{"status", "D"},
			},
		)

		// End Example 58

		require.NoError(t, err)
		require.Equal(t, int64(1), result.DeletedCount)

	}

	{
		// Start Example 56

		result, err := coll.DeleteMany(context.Background(), bson.D{})

		// End Example 56

		require.NoError(t, err)
		require.Equal(t, int64(2), result.DeletedCount)
	}
}

var log = logger.New(ioutil.Discard, "", logger.LstdFlags)

// Start Transactions Intro Example 1

// UpdateEmployeeInfo is an example function demonstrating transactions.
func UpdateEmployeeInfo(ctx context.Context, client *mongo.Client) error {
	employees := client.Database("hr").Collection("employees")
	events := client.Database("reporting").Collection("events")

	return client.UseSession(ctx, func(sctx mongo.SessionContext) error {
		err := sctx.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)
		if err != nil {
			return err
		}

		_, err = employees.UpdateOne(sctx, bson.D{{"employee", 3}}, bson.D{{"$set", bson.D{{"status", "Inactive"}}}})
		if err != nil {
			sctx.AbortTransaction(sctx)
			log.Println("caught exception during transaction, aborting.")
			return err
		}
		_, err = events.InsertOne(sctx, bson.D{{"employee", 3}, {"status", bson.D{{"new", "Inactive"}, {"old", "Active"}}}})
		if err != nil {
			sctx.AbortTransaction(sctx)
			log.Println("caught exception during transaction, aborting.")
			return err
		}

		for {
			err = sctx.CommitTransaction(sctx)
			switch e := err.(type) {
			case nil:
				return nil
			case mongo.CommandError:
				if e.HasErrorLabel("UnknownTransactionCommitResult") {
					log.Println("UnknownTransactionCommitResult, retrying commit operation...")
					continue
				}
				log.Println("Error during commit...")
				return e
			default:
				log.Println("Error during commit...")
				return e
			}
		}
	})
}

// End Transactions Intro Example 1

// Start Transactions Retry Example 1

// RunTransactionWithRetry is an example function demonstrating transaction retry logic.
func RunTransactionWithRetry(sctx mongo.SessionContext, txnFn func(mongo.SessionContext) error) error {
	for {
		err := txnFn(sctx) // Performs transaction.
		if err == nil {
			return nil
		}

		log.Println("Transaction aborted. Caught exception during transaction.")

		// If transient error, retry the whole transaction
		if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.HasErrorLabel("TransientTransactionError") {
			log.Println("TransientTransactionError, retrying transaction...")
			continue
		}
		return err
	}
}

// End Transactions Retry Example 1

// Start Transactions Retry Example 2

// CommitWithRetry is an example function demonstrating transaction commit with retry logic.
func CommitWithRetry(sctx mongo.SessionContext) error {
	for {
		err := sctx.CommitTransaction(sctx)
		switch e := err.(type) {
		case nil:
			log.Println("Transaction committed.")
			return nil
		case mongo.CommandError:
			// Can retry commit
			if e.HasErrorLabel("UnknownTransactionCommitResult") {
				log.Println("UnknownTransactionCommitResult, retrying commit operation...")
				continue
			}
			log.Println("Error during commit...")
			return e
		default:
			log.Println("Error during commit...")
			return e
		}
	}
}

// End Transactions Retry Example 2

// TransactionsExamples contains examples for transaction operations.
func TransactionsExamples(ctx context.Context, client *mongo.Client) error {
	_, err := client.Database("hr").Collection("employees").InsertOne(ctx, bson.D{{"pi", 3.14159}})
	if err != nil {
		return err
	}
	_, err = client.Database("hr").Collection("employees").DeleteOne(ctx, bson.D{{"pi", 3.14159}})
	if err != nil {
		return err
	}
	_, err = client.Database("reporting").Collection("events").InsertOne(ctx, bson.D{{"pi", 3.14159}})
	if err != nil {
		return err
	}
	_, err = client.Database("reporting").Collection("events").DeleteOne(ctx, bson.D{{"pi", 3.14159}})
	if err != nil {
		return err
	}
	// Start Transactions Retry Example 3

	runTransactionWithRetry := func(sctx mongo.SessionContext, txnFn func(mongo.SessionContext) error) error {
		for {
			err := txnFn(sctx) // Performs transaction.
			if err == nil {
				return nil
			}

			log.Println("Transaction aborted. Caught exception during transaction.")

			// If transient error, retry the whole transaction
			if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.HasErrorLabel("TransientTransactionError") {
				log.Println("TransientTransactionError, retrying transaction...")
				continue
			}
			return err
		}
	}

	commitWithRetry := func(sctx mongo.SessionContext) error {
		for {
			err := sctx.CommitTransaction(sctx)
			switch e := err.(type) {
			case nil:
				log.Println("Transaction committed.")
				return nil
			case mongo.CommandError:
				// Can retry commit
				if e.HasErrorLabel("UnknownTransactionCommitResult") {
					log.Println("UnknownTransactionCommitResult, retrying commit operation...")
					continue
				}
				log.Println("Error during commit...")
				return e
			default:
				log.Println("Error during commit...")
				return e
			}
		}
	}

	// Updates two collections in a transaction.
	updateEmployeeInfo := func(sctx mongo.SessionContext) error {
		employees := client.Database("hr").Collection("employees")
		events := client.Database("reporting").Collection("events")

		err := sctx.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)
		if err != nil {
			return err
		}

		_, err = employees.UpdateOne(sctx, bson.D{{"employee", 3}}, bson.D{{"$set", bson.D{{"status", "Inactive"}}}})
		if err != nil {
			sctx.AbortTransaction(sctx)
			log.Println("caught exception during transaction, aborting.")
			return err
		}
		_, err = events.InsertOne(sctx, bson.D{{"employee", 3}, {"status", bson.D{{"new", "Inactive"}, {"old", "Active"}}}})
		if err != nil {
			sctx.AbortTransaction(sctx)
			log.Println("caught exception during transaction, aborting.")
			return err
		}

		return commitWithRetry(sctx)
	}

	return client.UseSessionWithOptions(
		ctx, options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return runTransactionWithRetry(sctx, updateEmployeeInfo)
		},
	)
}

// End Transactions Retry Example 3

// Start Transactions withTxn API Example 1

// WithTransactionExample is an example of using the Session.WithTransaction function.
func WithTransactionExample() {
	ctx := context.Background()
	// For a replica set, include the replica set name and a seedlist of the members in the URI string; e.g.
	// uri := "mongodb://mongodb0.example.com:27017,mongodb1.example.com:27017/?replicaSet=myRepl"
	// For a sharded cluster, connect to the mongos instances; e.g.
	// uri := "mongodb://mongos0.example.com:27017,mongos1.example.com:27017/"
	var uri string

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}
	defer func() { _ = client.Disconnect(ctx) }()

	// Prereq: Create collections.
	wcMajority := writeconcern.New(writeconcern.WMajority(), writeconcern.WTimeout(1*time.Second))
	wcMajorityCollectionOpts := options.Collection().SetWriteConcern(wcMajority)
	fooColl := client.Database("mydb1").Collection("foo", wcMajorityCollectionOpts)
	barColl := client.Database("mydb1").Collection("bar", wcMajorityCollectionOpts)

	// Step 1: Define the callback that specifies the sequence of operations to perform inside the transaction.
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Important: You must pass sessCtx as the Context parameter to the operations for them to be executed in the
		// transaction.
		if _, err := fooColl.InsertOne(sessCtx, bson.D{{"abc", 1}}); err != nil {
			return nil, err
		}
		if _, err := barColl.InsertOne(sessCtx, bson.D{{"xyz", 999}}); err != nil {
			return nil, err
		}

		return nil, nil
	}

	// Step 2: Start a session and run the callback using WithTransaction.
	session, err := client.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(ctx)

	result, err := session.WithTransaction(ctx, callback)
	if err != nil {
		panic(err)
	}
	fmt.Printf("result: %v\n", result)
}

// End Transactions withTxn API Example 1

// ChangeStreamExamples contains examples of changestream operations.
func ChangeStreamExamples(t *testing.T, db *mongo.Database) {
	ctx := context.Background()

	coll := db.Collection("inventory_changestream")

	err := coll.Drop(context.Background())
	require.NoError(t, err)

	_, err = coll.InsertOne(ctx, bson.D{{"x", int32(1)}})
	require.NoError(t, err)

	var stop int32

	doInserts := func(coll *mongo.Collection) {
		for atomic.LoadInt32(&stop) == 0 {
			_, err = coll.InsertOne(ctx, bson.D{{"x", 1}})
			time.Sleep(10 * time.Millisecond)
			coll.DeleteOne(ctx, bson.D{{"x", 1}})
		}
	}

	go doInserts(coll)

	{
		// Start Changestream Example 1

		cs, err := coll.Watch(ctx, mongo.Pipeline{})
		require.NoError(t, err)
		defer cs.Close(ctx)

		ok := cs.Next(ctx)
		next := cs.Current

		// End Changestream Example 1

		require.True(t, ok)
		require.NoError(t, err)
		require.NotEqual(t, len(next), 0)
	}
	{
		// Start Changestream Example 2

		cs, err := coll.Watch(ctx, mongo.Pipeline{}, options.ChangeStream().SetFullDocument(options.UpdateLookup))
		require.NoError(t, err)
		defer cs.Close(ctx)

		ok := cs.Next(ctx)
		next := cs.Current

		// End Changestream Example 2

		require.True(t, ok)
		require.NoError(t, err)
		require.NotEqual(t, len(next), 0)
	}

	{
		original, err := coll.Watch(ctx, mongo.Pipeline{})
		require.NoError(t, err)
		defer original.Close(ctx)

		ok := original.Next(ctx)
		require.True(t, ok)

		// Start Changestream Example 3
		resumeToken := original.ResumeToken()

		cs, err := coll.Watch(ctx, mongo.Pipeline{}, options.ChangeStream().SetResumeAfter(resumeToken))
		require.NoError(t, err)
		defer cs.Close(ctx)

		ok = cs.Next(ctx)
		result := cs.Current

		// End Changestream Example 3

		require.True(t, ok)
		require.NoError(t, err)
		require.NotEqual(t, len(result), 0)
	}

	{
		// Start Changestream Example 4
		pipeline := mongo.Pipeline{bson.D{{"$match", bson.D{{"$or",
			bson.A{
				bson.D{{"fullDocument.username", "alice"}},
				bson.D{{"operationType", "delete"}}}}},
		}}}
		cs, err := coll.Watch(ctx, pipeline)
		require.NoError(t, err)
		defer cs.Close(ctx)

		ok := cs.Next(ctx)
		next := cs.Current

		// End Changestream Example 4

		require.True(t, ok)
		require.NoError(t, err)
		require.NotEqual(t, len(next), 0)
	}

	atomic.StoreInt32(&stop, 1)
}

// AggregationExamples contains examples of aggregation operations.
func AggregationExamples(t *testing.T, db *mongo.Database) {
	ctx := context.Background()

	salesColl := db.Collection("sales")
	airlinesColl := db.Collection("airlines")
	airAlliancesColl := db.Collection("air_alliances")

	err := salesColl.Drop(ctx)
	require.NoError(t, err)
	err = airlinesColl.Drop(ctx)
	require.NoError(t, err)
	err = airAlliancesColl.Drop(ctx)
	require.NoError(t, err)

	date20180208 := parseDate(t, "2018-02-08T09:00:00.000Z")
	date20180109 := parseDate(t, "2018-01-09T07:12:00.000Z")
	date20180127 := parseDate(t, "2018-01-27T09:13:00.000Z")
	date20180203 := parseDate(t, "2018-02-03T07:58:00.000Z")
	date20180205 := parseDate(t, "2018-02-05T06:03:00.000Z")
	date20180111 := parseDate(t, "2018-01-11T07:15:00.000Z")

	sales := []interface{}{
		bson.D{
			{"date", date20180208},
			{"items", bson.A{
				bson.D{
					{"fruit", "kiwi"},
					{"quantity", 2},
					{"price", 0.5},
				},
				bson.D{
					{"fruit", "apple"},
					{"quantity", 1},
					{"price", 1.0},
				},
			}},
		},
		bson.D{
			{"date", date20180109},
			{"items", bson.A{
				bson.D{
					{"fruit", "banana"},
					{"quantity", 8},
					{"price", 1.0},
				},
				bson.D{
					{"fruit", "apple"},
					{"quantity", 1},
					{"price", 1.0},
				},
				bson.D{
					{"fruit", "papaya"},
					{"quantity", 1},
					{"price", 4.0},
				},
			}},
		},
		bson.D{
			{"date", date20180127},
			{"items", bson.A{
				bson.D{
					{"fruit", "banana"},
					{"quantity", 1},
					{"price", 1.0},
				},
			}},
		},
		bson.D{
			{"date", date20180203},
			{"items", bson.A{
				bson.D{
					{"fruit", "banana"},
					{"quantity", 1},
					{"price", 1.0},
				},
			}},
		},
		bson.D{
			{"date", date20180205},
			{"items", bson.A{
				bson.D{
					{"fruit", "banana"},
					{"quantity", 1},
					{"price", 1.0},
				},
				bson.D{
					{"fruit", "mango"},
					{"quantity", 2},
					{"price", 2.0},
				},
				bson.D{
					{"fruit", "apple"},
					{"quantity", 1},
					{"price", 1.0},
				},
			}},
		},
		bson.D{
			{"date", date20180111},
			{"items", bson.A{
				bson.D{
					{"fruit", "banana"},
					{"quantity", 1},
					{"price", 1.0},
				},
				bson.D{
					{"fruit", "apple"},
					{"quantity", 1},
					{"price", 1.0},
				},
				bson.D{
					{"fruit", "papaya"},
					{"quantity", 3},
					{"price", 4.0},
				},
			}},
		},
	}
	airlines := []interface{}{
		bson.D{
			{"airline", 17},
			{"name", "Air Canada"},
			{"alias", "AC"},
			{"iata", "ACA"},
			{"icao", "AIR CANADA"},
			{"active", "Y"},
			{"country", "Canada"},
			{"base", "TAL"},
		},
		bson.D{
			{"airline", 18},
			{"name", "Turkish Airlines"},
			{"alias", "YK"},
			{"iata", "TRK"},
			{"icao", "TURKISH"},
			{"active", "Y"},
			{"country", "Turkey"},
			{"base", "AET"},
		},
		bson.D{
			{"airline", 22},
			{"name", "Saudia"},
			{"alias", "SV"},
			{"iata", "SVA"},
			{"icao", "SAUDIA"},
			{"active", "Y"},
			{"country", "Saudi Arabia"},
			{"base", "JSU"},
		},
		bson.D{
			{"airline", 29},
			{"name", "Finnair"},
			{"alias", "AY"},
			{"iata", "FIN"},
			{"icao", "FINNAIR"},
			{"active", "Y"},
			{"country", "Finland"},
			{"base", "JMZ"},
		},
		bson.D{
			{"airline", 34},
			{"name", "Afric'air Express"},
			{"alias", ""},
			{"iata", "AAX"},
			{"icao", "AFREX"},
			{"active", "N"},
			{"country", "Ivory Coast"},
			{"base", "LOK"},
		},
		bson.D{
			{"airline", 37},
			{"name", "Artem-Avia"},
			{"alias", ""},
			{"iata", "ABA"},
			{"icao", "ARTEM-AVIA"},
			{"active", "N"},
			{"country", "Ukraine"},
			{"base", "JBR"},
		},
		bson.D{
			{"airline", 38},
			{"name", "Lufthansa"},
			{"alias", "LH"},
			{"iata", "DLH"},
			{"icao", "LUFTHANSA"},
			{"active", "Y"},
			{"country", "Germany"},
			{"base", "CYS"},
		},
	}
	airAlliances := []interface{}{
		bson.D{
			{"name", "Star Alliance"},
			{"airlines", bson.A{
				"Air Canada",
				"Avianca",
				"Air China",
				"Air New Zealand",
				"Asiana Airlines",
				"Brussels Airlines",
				"Copa Airlines",
				"Croatia Airlines",
				"EgyptAir",
				"TAP Portugal",
				"United Airlines",
				"Turkish Airlines",
				"Swiss International Air Lines",
				"Lufthansa",
			}},
		},
		bson.D{
			{"name", "SkyTeam"},
			{"airlines", bson.A{
				"Aerolinias Argentinas",
				"Aeromexico",
				"Air Europa",
				"Air France",
				"Alitalia",
				"Delta Air Lines",
				"Garuda Indonesia",
				"Kenya Airways",
				"KLM",
				"Korean Air",
				"Middle East Airlines",
				"Saudia",
			}},
		},
		bson.D{
			{"name", "OneWorld"},
			{"airlines", bson.A{
				"Air Berlin",
				"American Airlines",
				"British Airways",
				"Cathay Pacific",
				"Finnair",
				"Iberia Airlines",
				"Japan Airlines",
				"LATAM Chile",
				"LATAM Brasil",
				"Malasya Airlines",
				"Canadian Airlines",
			}},
		},
	}

	salesResult, salesErr := salesColl.InsertMany(ctx, sales)
	airlinesResult, airlinesErr := airlinesColl.InsertMany(ctx, airlines)
	airAlliancesResult, airAlliancesErr := airAlliancesColl.InsertMany(ctx, airAlliances)

	require.NoError(t, salesErr)
	require.Len(t, salesResult.InsertedIDs, 6)
	require.NoError(t, airlinesErr)
	require.Len(t, airlinesResult.InsertedIDs, 7)
	require.NoError(t, airAlliancesErr)
	require.Len(t, airAlliancesResult.InsertedIDs, 3)

	{
		// Start Aggregation Example 1
		pipeline := mongo.Pipeline{
			{
				{"$match", bson.D{
					{"items.fruit", "banana"},
				}},
			},
			{
				{"$sort", bson.D{
					{"date", 1},
				}},
			},
		}

		cursor, err := salesColl.Aggregate(ctx, pipeline)

		// End Aggregation Example 1

		require.NoError(t, err)
		defer cursor.Close(ctx)
		requireCursorLength(t, cursor, 5)
	}
	{
		// Start Aggregation Example 2
		pipeline := mongo.Pipeline{
			{
				{"$unwind", "$items"},
			},
			{
				{"$match", bson.D{
					{"items.fruit", "banana"},
				}},
			},
			{
				{"$group", bson.D{
					{"_id", bson.D{
						{"day", bson.D{
							{"$dayOfWeek", "$date"},
						}},
					}},
					{"count", bson.D{
						{"$sum", "$items.quantity"},
					}},
				}},
			},
			{
				{"$project", bson.D{
					{"dayOfWeek", "$_id.day"},
					{"numberSold", "$count"},
					{"_id", 0},
				}},
			},
			{
				{"$sort", bson.D{
					{"numberSold", 1},
				}},
			},
		}

		cursor, err := salesColl.Aggregate(ctx, pipeline)

		// End Aggregation Example 2

		require.NoError(t, err)
		defer cursor.Close(ctx)
		requireCursorLength(t, cursor, 4)
	}
	{
		// Start Aggregation Example 3
		pipeline := mongo.Pipeline{
			{
				{"$unwind", "$items"},
			},
			{
				{"$group", bson.D{
					{"_id", bson.D{
						{"day", bson.D{
							{"$dayOfWeek", "$date"},
						}},
					}},
					{"items_sold", bson.D{
						{"$sum", "$items.quantity"},
					}},
					{"revenue", bson.D{
						{"$sum", bson.D{
							{"$multiply", bson.A{"$items.quantity", "$items.price"}},
						}},
					}},
				}},
			},
			{
				{"$project", bson.D{
					{"day", "$_id.day"},
					{"revenue", 1},
					{"items_sold", 1},
					{"discount", bson.D{
						{"$cond", bson.D{
							{"if", bson.D{
								{"$lte", bson.A{"$revenue", 250}},
							}},
							{"then", 25},
							{"else", 0},
						}},
					}},
				}},
			},
		}

		cursor, err := salesColl.Aggregate(ctx, pipeline)

		// End Aggregation Example 3

		require.NoError(t, err)
		defer cursor.Close(ctx)
		requireCursorLength(t, cursor, 4)
	}
	{
		// Start Aggregation Example 4
		pipeline := mongo.Pipeline{
			{
				{"$lookup", bson.D{
					{"from", "air_airlines"},
					{"let", bson.D{
						{"constituents", "$airlines"}},
					},
					{"pipeline", bson.A{bson.D{
						{"$match", bson.D{
							{"$expr", bson.D{
								{"$in", bson.A{"$name", "$$constituents"}},
							}},
						}},
					}}},
					{"as", "airlines"},
				}},
			},
			{
				{"$project", bson.D{
					{"_id", 0},
					{"name", 1},
					{"airlines", bson.D{
						{"$filter", bson.D{
							{"input", "$airlines"},
							{"as", "airline"},
							{"cond", bson.D{
								{"$eq", bson.A{"$$airline.country", "Canada"}},
							}},
						}},
					}},
				}},
			},
		}

		cursor, err := airAlliancesColl.Aggregate(ctx, pipeline)

		// End Aggregation Example 4

		require.NoError(t, err)
		defer cursor.Close(ctx)
		requireCursorLength(t, cursor, 3)
	}
}

// RunCommandExamples contains examples of RunCommand operations.
func RunCommandExamples(t *testing.T, db *mongo.Database) {
	ctx := context.Background()

	coll := db.Collection("restaurants")

	err := coll.Drop(ctx)
	require.NoError(t, err)

	restaurants := []interface{}{
		bson.D{
			{"name", "Chez Panisse"},
			{"city", "Oakland"},
			{"state", "California"},
			{"country", "United States"},
			{"rating", 4.4},
		},
		bson.D{
			{"name", "Central"},
			{"city", "Lima"},
			{"country", "Peru"},
			{"rating", 4.8},
		},
		bson.D{
			{"name", "Eleven Madison Park"},
			{"city", "New York City"},
			{"state", "New York"},
			{"country", "United States"},
			{"rating", 4.6},
		},
		bson.D{
			{"name", "Gaggan"},
			{"city", "Bangkok"},
			{"country", "Thailand"},
			{"rating", 4.3},
		},
		bson.D{
			{"name", "Dad's Grill"},
			{"city", "Oklahoma City"},
			{"state", "Oklahoma"},
			{"country", "United States"},
			{"rating", 2.1},
		},
	}

	result, err := coll.InsertMany(ctx, restaurants)
	require.NoError(t, err)
	require.Len(t, result.InsertedIDs, 5)

	{
		// Start RunCommand Example 1
		res := db.RunCommand(ctx, bson.D{{"buildInfo", 1}})

		// End RunCommand Example 1

		err := res.Err()
		require.NoError(t, err)
	}
	{
		// Start RunCommand Example 2
		res := db.RunCommand(ctx, bson.D{{"collStats", "restaurants"}})

		// End RunCommand Example 2

		err := res.Err()
		require.NoError(t, err)
	}
}

// IndexExamples contains examples of Index operations.
func IndexExamples(t *testing.T, db *mongo.Database) {
	ctx := context.Background()

	recordsColl := db.Collection("records")
	restaurantsColl := db.Collection("restaurants")

	err := recordsColl.Drop(ctx)
	require.NoError(t, err)
	err = restaurantsColl.Drop(ctx)
	require.NoError(t, err)

	records := []interface{}{
		bson.D{
			{"student", "Marty McFly"},
			{"classYear", 1986},
			{"school", "Hill Valley High"},
			{"score", 56.5},
		},
		bson.D{
			{"student", "Ferris F. Bueller"},
			{"classYear", 1987},
			{"school", "Glenbrook North High"},
			{"status", "Suspended"},
			{"score", 76.0},
		},
		bson.D{
			{"student", "Reynard Muldoon"},
			{"classYear", 2007},
			{"school", "Stonetown Middle"},
			{"score", 99.9},
		},
	}
	restaurants := []interface{}{
		bson.D{
			{"name", "Chez Panisse"},
			{"cuisine", "American/French"},
			{"city", "Oakland"},
			{"state", "California"},
			{"country", "United States"},
			{"rating", 4.9},
		},
		bson.D{
			{"name", "Central"},
			{"cuisine", "Peruvian"},
			{"city", "Lima"},
			{"country", "Peru"},
			{"rating", 5.8},
		},
		bson.D{
			{"name", "Eleven Madison Park"},
			{"cuisine", "French"},
			{"city", "New York City"},
			{"state", "New York"},
			{"country", "United States"},
			{"rating", 7.1},
		},
		bson.D{
			{"name", "Gaggan"},
			{"cuisine", "Thai Fusion"},
			{"city", "Bangkok"},
			{"country", "Thailand"},
			{"rating", 9.2},
		},
		bson.D{
			{"name", "Dad's Grill"},
			{"cuisine", "BBQ"},
			{"city", "Oklahoma City"},
			{"state", "Oklahoma"},
			{"country", "United States"},
			{"rating", 2.1},
		},
	}

	recordsResult, recordsErr := recordsColl.InsertMany(ctx, records)
	restaurantsResult, restaurantsErr := restaurantsColl.InsertMany(ctx, restaurants)

	require.NoError(t, recordsErr)
	require.Len(t, recordsResult.InsertedIDs, 3)
	require.NoError(t, restaurantsErr)
	require.Len(t, restaurantsResult.InsertedIDs, 5)

	{
		// Start Index Example 1
		indexModel := mongo.IndexModel{
			Keys: bson.D{
				{"score", 1},
			},
		}
		_, err := recordsColl.Indexes().CreateOne(ctx, indexModel)

		// End Index Example 1

		require.NoError(t, err)
	}
	{
		// Start Index Example 2
		partialFilterExpression := bson.D{
			{"rating", bson.D{
				{"$gt", 5},
			}},
		}
		indexModel := mongo.IndexModel{
			Keys: bson.D{
				{"cuisine", 1},
				{"name", 1},
			},
			Options: options.Index().SetPartialFilterExpression(partialFilterExpression),
		}

		_, err := restaurantsColl.Indexes().CreateOne(ctx, indexModel)

		// End Index Example 2

		require.NoError(t, err)
	}
}

// Start Versioned API Example 1

// VersionedAPIExample is an example of creating a client with versioned API.
func VersionedAPIExample() {
	ctx := context.Background()
	// For a replica set, include the replica set name and a seedlist of the members in the URI string; e.g.
	// uri := "mongodb://mongodb0.example.com:27017,mongodb1.example.com:27017/?replicaSet=myRepl"
	// For a sharded cluster, connect to the mongos instances; e.g.
	// uri := "mongodb://mongos0.example.com:27017,mongos1.example.com:27017/"
	uri := "mongodb://localhost:27017"

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOpts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPIOptions)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}
	defer func() { _ = client.Disconnect(ctx) }()
}

// End Versioned API Example 1

// Start Versioned API Example 2

// VersionedAPIStrictExample is an example of creating a client with strict versioned API.
func VersionedAPIStrictExample() {
	ctx := context.Background()
	// For a replica set, include the replica set name and a seedlist of the members in the URI string; e.g.
	// uri := "mongodb://mongodb0.example.com:27017,mongodb1.example.com:27017/?replicaSet=myRepl"
	// For a sharded cluster, connect to the mongos instances; e.g.
	// uri := "mongodb://mongos0.example.com:27017,mongos1.example.com:27017/"
	uri := "mongodb://localhost:27017"

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1).SetStrict(true)
	clientOpts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPIOptions)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}
	defer func() { _ = client.Disconnect(ctx) }()
}

// End Versioned API Example 2

// Start Versioned API Example 3

// VersionedAPINonStrictExample is an example of creating a client with non-strict versioned API.
func VersionedAPINonStrictExample() {
	ctx := context.Background()
	// For a replica set, include the replica set name and a seedlist of the members in the URI string; e.g.
	// uri := "mongodb://mongodb0.example.com:27017,mongodb1.example.com:27017/?replicaSet=myRepl"
	// For a sharded cluster, connect to the mongos instances; e.g.
	// uri := "mongodb://mongos0.example.com:27017,mongos1.example.com:27017/"
	uri := "mongodb://localhost:27017"

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1).SetStrict(false)
	clientOpts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPIOptions)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}
	defer func() { _ = client.Disconnect(ctx) }()
}

// End Versioned API Example 3

// Start Versioned API Example 4

// VersionedAPIDeprecationErrorsExample is an example of creating a client with versioned API
// with deprecation errors.
func VersionedAPIDeprecationErrorsExample() {
	ctx := context.Background()
	// For a replica set, include the replica set name and a seedlist of the members in the URI string; e.g.
	// uri := "mongodb://mongodb0.example.com:27017,mongodb1.example.com:27017/?replicaSet=myRepl"
	// For a sharded cluster, connect to the mongos instances; e.g.
	// uri := "mongodb://mongos0.example.com:27017,mongos1.example.com:27017/"
	uri := "mongodb://localhost:27017"

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1).SetDeprecationErrors(true)
	clientOpts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPIOptions)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}
	defer func() { _ = client.Disconnect(ctx) }()
}

// End Versioned API Example 4

// VersionedAPIExamples runs all versioned API examples.
func VersionedAPIExamples() {
	VersionedAPIExample()
	VersionedAPIStrictExample()
	VersionedAPINonStrictExample()
	VersionedAPIDeprecationErrorsExample()
}
