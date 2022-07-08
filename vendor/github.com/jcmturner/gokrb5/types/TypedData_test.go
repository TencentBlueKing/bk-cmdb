package types

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jcmturner/gokrb5.v7/iana/patype"
	"gopkg.in/jcmturner/gokrb5.v7/test/testdata"
)

func TestUnmarshalTypedData(t *testing.T) {
	t.Parallel()
	var a TypedDataSequence
	b, err := hex.DecodeString(testdata.MarshaledKRB5typed_data)
	if err != nil {
		t.Fatalf("Test vector read error: %v", err)
	}
	err = a.Unmarshal(b)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	assert.Equal(t, 2, len(a), "Number of typed data elements not as expected")
	for i, d := range a {
		assert.Equal(t, patype.PA_SAM_RESPONSE, d.DataType, fmt.Sprintf("Data type of element %d not as expected", i+1))
		assert.Equal(t, []byte(testdata.TEST_PADATA_VALUE), d.DataValue, fmt.Sprintf("Data value of element %d not as expected", i+1))
	}
}
