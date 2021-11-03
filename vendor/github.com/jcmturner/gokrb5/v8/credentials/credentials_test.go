package credentials

import (
	"testing"

	"github.com/jcmturner/goidentity/v6"
	"github.com/stretchr/testify/assert"
)

func TestImplementsInterface(t *testing.T) {
	t.Parallel()
	u := new(Credentials)
	i := new(goidentity.Identity)
	assert.Implements(t, i, u, "Credentials type does not implement the Identity interface")
}

func TestCredentials_Marshal(t *testing.T) {
	var cred Credentials
	b, err := cred.Marshal()
	if err != nil {
		t.Fatalf("could not marshal credetials: %v", err)
	}
	var credum Credentials
	err = credum.Unmarshal(b)
	if err != nil {
		t.Fatalf("could not unmarshal credetials: %v", err)
	}
}
