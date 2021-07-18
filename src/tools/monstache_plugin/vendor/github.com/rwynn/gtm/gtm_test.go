package gtm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestInvalidCursorCheck(t *testing.T) {
	positionLost := mongo.CommandError{Code: 136}
	assert.True(t, invalidCursor(positionLost))
	err := mongo.CommandError{Code: 999}
	assert.False(t, invalidCursor(err))
	assert.False(t, invalidCursor(nil))
}
