package mongobyc

// InsertOneResult  is a result of an InsertOne operation
type InsertOneResult struct {
	Count int64
	// the identifier that was inserted
	InsertedID interface{}
}

// InsertManyResult is a result of an InsertMany operation
type InsertManyResult struct {
	// Maps the indexes of inserted documents to their _id fields
	InsertedIDS []interface{}
}

// DeleteResult is a result of an DeleteOne operat
type DeleteResult struct {
	// The number of documents that were deleted.
	DeletedCount int64 `bson:"n"`
}

// UpdateResult is a result of an update operation.
type UpdateResult struct {
	// The number of documents that matched the filter.
	MatchedCount int64
	// The number of documents that were modified.
	ModifiedCount int64
	// The identifier of the inserted document if an upsert took place.
	UpsertedID interface{}
}

// DocumentResult represents a single document returned from an operation. If
// the operation returned an error, the Err method of DocumentResult will
// return that error.
type DocumentResult struct {
}
