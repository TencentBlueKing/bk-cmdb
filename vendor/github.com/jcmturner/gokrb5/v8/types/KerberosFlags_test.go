package types

import (
	"github.com/jcmturner/gofork/encoding/asn1"
	"github.com/jcmturner/gokrb5/v8/iana/flags"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKerberosFlags_SetFlag(t *testing.T) {
	t.Parallel()
	b := []byte{byte(64), byte(0), byte(0), byte(16)}
	var f asn1.BitString
	SetFlag(&f, flags.Forwardable)
	SetFlag(&f, flags.RenewableOK)
	assert.Equal(t, b, f.Bytes, "Flag bytes not as expected")
}

func TestKerberosFlags_UnsetFlag(t *testing.T) {
	t.Parallel()
	b := []byte{byte(64), byte(0), byte(0), byte(0)}
	var f asn1.BitString
	SetFlag(&f, flags.Forwardable)
	SetFlag(&f, flags.RenewableOK)
	UnsetFlag(&f, flags.RenewableOK)
	assert.Equal(t, b, f.Bytes, "Flag bytes not as expected")
}

func TestKerberosFlags_IsFlagSet(t *testing.T) {
	t.Parallel()
	var f asn1.BitString
	SetFlag(&f, flags.Forwardable)
	SetFlag(&f, flags.RenewableOK)
	UnsetFlag(&f, flags.Proxiable)
	assert.True(t, IsFlagSet(&f, flags.Forwardable))
	assert.True(t, IsFlagSet(&f, flags.RenewableOK))
	assert.False(t, IsFlagSet(&f, flags.Proxiable))
}
