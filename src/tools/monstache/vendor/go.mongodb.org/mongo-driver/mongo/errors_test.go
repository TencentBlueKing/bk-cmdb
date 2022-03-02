package mongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestErrorMessages(t *testing.T) {
	details, err := bson.Marshal(bson.D{{"details", bson.D{{"operatorName", "$jsonSchema"}}}})
	require.Nil(t, err, "unexpected error marshaling BSON")

	cases := []struct {
		desc     string
		err      error
		expected string
	}{
		{
			desc: "WriteException error message should contain the WriteError Message and Details",
			err: WriteException{
				WriteErrors: WriteErrors{
					{
						Message: "test message 1",
						Details: details,
					},
					{
						Message: "test message 2",
						Details: details,
					},
				},
			},
			expected: `write exception: write errors: [test message 1: {"details": {"operatorName": "$jsonSchema"}}, test message 2: {"details": {"operatorName": "$jsonSchema"}}]`,
		},
		{
			desc: "BulkWriteException error message should contain the WriteError Message and Details",
			err: BulkWriteException{
				WriteErrors: []BulkWriteError{
					{
						WriteError: WriteError{
							Message: "test message 1",
							Details: details,
						},
					},
					{
						WriteError: WriteError{
							Message: "test message 2",
							Details: details,
						},
					},
				},
			},
			expected: `bulk write exception: write errors: [test message 1: {"details": {"operatorName": "$jsonSchema"}}, test message 2: {"details": {"operatorName": "$jsonSchema"}}]`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, tc.err.Error())
		})
	}
}
