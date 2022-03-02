package messages

import (
	"encoding/hex"
	"testing"

	"github.com/jcmturner/gokrb5/v8/iana"
	"github.com/jcmturner/gokrb5/v8/iana/msgtype"
	"github.com/jcmturner/gokrb5/v8/iana/nametype"
	"github.com/jcmturner/gokrb5/v8/test/testdata"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalAPReq(t *testing.T) {
	t.Parallel()
	var a APReq
	b, err := hex.DecodeString(testdata.MarshaledKRB5ap_req)
	if err != nil {
		t.Fatalf("Test vector read error: %v", err)
	}
	err = a.Unmarshal(b)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	assert.Equal(t, iana.PVNO, a.PVNO, "PVNO not as expected")
	assert.Equal(t, msgtype.KRB_AP_REQ, a.MsgType, "MsgType is not as expected")
	assert.Equal(t, "fedcba98", hex.EncodeToString(a.APOptions.Bytes), "AP Options not as expected")
	assert.Equal(t, iana.PVNO, a.Ticket.TktVNO, "Ticket VNO not as expected")
	assert.Equal(t, testdata.TEST_REALM, a.Ticket.Realm, "Ticket realm not as expected")
	assert.Equal(t, nametype.KRB_NT_PRINCIPAL, a.Ticket.SName.NameType, "Ticket SName NameType not as expected")
	assert.Equal(t, len(testdata.TEST_PRINCIPALNAME_NAMESTRING), len(a.Ticket.SName.NameString), "Ticket SName does not have the expected number of NameStrings")
	assert.Equal(t, testdata.TEST_PRINCIPALNAME_NAMESTRING, a.Ticket.SName.NameString, "Ticket SName name string entries not as expected")
	assert.Equal(t, testdata.TEST_ETYPE, a.Ticket.EncPart.EType, "Ticket encPart etype not as expected")
	assert.Equal(t, iana.PVNO, a.Ticket.EncPart.KVNO, "Ticket encPart KVNO not as expected")
	assert.Equal(t, []byte(testdata.TEST_CIPHERTEXT), a.Ticket.EncPart.Cipher, "Ticket encPart cipher not as expected")
}

func TestMarshalAPReq(t *testing.T) {
	t.Parallel()
	var a APReq
	b, err := hex.DecodeString(testdata.MarshaledKRB5ap_req)
	if err != nil {
		t.Fatalf("Test vector read error: %v", err)
	}
	err = a.Unmarshal(b)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	mb, err := a.Marshal()
	if err != nil {
		t.Fatalf("Marshal of ticket errored: %v", err)
	}
	assert.Equal(t, b, mb, "Marshal bytes of Authenticator not as expected")
}
