package messages

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/jcmturner/gokrb5/v8/iana"
	"github.com/jcmturner/gokrb5/v8/iana/addrtype"
	"github.com/jcmturner/gokrb5/v8/iana/msgtype"
	"github.com/jcmturner/gokrb5/v8/iana/nametype"
	"github.com/jcmturner/gokrb5/v8/test/testdata"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalKRBCred(t *testing.T) {
	t.Parallel()
	var a KRBCred
	b, err := hex.DecodeString(testdata.MarshaledKRB5cred)
	if err != nil {
		t.Fatalf("Test vector read error: %v", err)
	}
	err = a.Unmarshal(b)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	assert.Equal(t, iana.PVNO, a.PVNO, "PVNO not as expected")
	assert.Equal(t, msgtype.KRB_CRED, a.MsgType, "Message type not as expected")
	assert.Equal(t, 2, len(a.Tickets), "Number of tickets not as expected")
	for i, tkt := range a.Tickets {
		assert.Equal(t, iana.PVNO, tkt.TktVNO, fmt.Sprintf("Ticket (%v) ticket-vno not as expected", i+1))
		assert.Equal(t, testdata.TEST_REALM, tkt.Realm, fmt.Sprintf("Ticket (%v) realm not as expected", i+1))
		assert.Equal(t, nametype.KRB_NT_PRINCIPAL, tkt.SName.NameType, fmt.Sprintf("Ticket (%v) SName NameType not as expected", i+1))
		assert.Equal(t, len(testdata.TEST_PRINCIPALNAME_NAMESTRING), len(tkt.SName.NameString), fmt.Sprintf("Ticket (%v) SName does not have the expected number of NameStrings", i+1))
		assert.Equal(t, testdata.TEST_PRINCIPALNAME_NAMESTRING, tkt.SName.NameString, fmt.Sprintf("Ticket (%v) SName name string entries not as expected", i+1))
		assert.Equal(t, testdata.TEST_ETYPE, tkt.EncPart.EType, fmt.Sprintf("Ticket (%v) encPart etype not as expected", i+1))
		assert.Equal(t, iana.PVNO, tkt.EncPart.KVNO, fmt.Sprintf("Ticket (%v) encPart KVNO not as expected", i+1))
		assert.Equal(t, []byte(testdata.TEST_CIPHERTEXT), tkt.EncPart.Cipher, fmt.Sprintf("Ticket (%v) encPart cipher not as expected", i+1))
	}
	assert.Equal(t, testdata.TEST_ETYPE, a.EncPart.EType, "encPart etype not as expected")
	assert.Equal(t, iana.PVNO, a.EncPart.KVNO, "encPart KVNO not as expected")
	assert.Equal(t, []byte(testdata.TEST_CIPHERTEXT), a.EncPart.Cipher, "encPart cipher not as expected")
}

func TestUnmarshalEncCredPart(t *testing.T) {
	t.Parallel()
	var a EncKrbCredPart
	b, err := hex.DecodeString(testdata.MarshaledKRB5enc_cred_part)
	if err != nil {
		t.Fatalf("Test vector read error: %v", err)
	}
	err = a.Unmarshal(b)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	//Parse the test time value into a time.Time type
	tt, _ := time.Parse(testdata.TEST_TIME_FORMAT, testdata.TEST_TIME)

	assert.Equal(t, 2, len(a.TicketInfo), "Number of ticket info items not as expected")
	for i, tkt := range a.TicketInfo {
		assert.Equal(t, int32(1), tkt.Key.KeyType, fmt.Sprintf("Key type not as expected in ticket info item %d", i+1))
		assert.Equal(t, []byte("12345678"), tkt.Key.KeyValue, fmt.Sprintf("Key value not as expected in ticket info item %d", i+1))
		assert.Equal(t, testdata.TEST_REALM, tkt.PRealm, fmt.Sprintf("PRealm not as expected on ticket info item %d", i+1))
		assert.Equal(t, nametype.KRB_NT_PRINCIPAL, tkt.PName.NameType, fmt.Sprintf("Ticket info (%v) PName NameType not as expected", i+1))
		assert.Equal(t, len(testdata.TEST_PRINCIPALNAME_NAMESTRING), len(tkt.PName.NameString), fmt.Sprintf("Ticket info (%v) PName does not have the expected number of NameStrings", i+1))
		assert.Equal(t, testdata.TEST_PRINCIPALNAME_NAMESTRING, tkt.PName.NameString, fmt.Sprintf("Ticket info (%v) PName name string entries not as expected", i+1))
		assert.Equal(t, "fedcba98", hex.EncodeToString(tkt.Flags.Bytes), fmt.Sprintf("Flags not as expected on ticket info %d", i+1))
		assert.Equal(t, tt, tkt.AuthTime, fmt.Sprintf("Auth time value not as expected for ticket info %d", i+1))
		assert.Equal(t, tt, tkt.StartTime, fmt.Sprintf("Start time value not as expected for ticket info %d", i+1))
		assert.Equal(t, tt, tkt.EndTime, fmt.Sprintf("End time value not as expected for ticket info %d", i+1))
		assert.Equal(t, tt, tkt.RenewTill, fmt.Sprintf("Renew Till time value not as expected for ticket info %d", i+1))
		assert.Equal(t, nametype.KRB_NT_PRINCIPAL, tkt.SName.NameType, fmt.Sprintf("Ticket info (%v) PName NameType not as expected", i+1))
		assert.Equal(t, len(testdata.TEST_PRINCIPALNAME_NAMESTRING), len(tkt.SName.NameString), fmt.Sprintf("Ticket info (%v) PName does not have the expected number of NameStrings", i+1))
		assert.Equal(t, testdata.TEST_PRINCIPALNAME_NAMESTRING, tkt.SName.NameString, fmt.Sprintf("Ticket info (%v) PName name string entries not as expected", i+1))
		assert.Equal(t, 2, len(tkt.CAddr), "Number of client addresses not as expected")
		for j, addr := range tkt.CAddr {
			assert.Equal(t, addrtype.IPv4, addr.AddrType, fmt.Sprintf("Host address type not as expected for address item %d within ticket info %d", j+1, i+1))
			assert.Equal(t, "12d00023", hex.EncodeToString(addr.Address), fmt.Sprintf("Host address not as expected for address item %d within ticket info %d", j+1, i+1))
		}
	}
	assert.Equal(t, testdata.TEST_NONCE, a.Nouce, "Nouce not as expected")
	assert.Equal(t, tt, a.Timestamp, "Timestamp not as expected")
	assert.Equal(t, 123456, a.Usec, "Microseconds not as expected")
	assert.Equal(t, addrtype.IPv4, a.SAddress.AddrType, "SAddress type not as expected")
	assert.Equal(t, "12d00023", hex.EncodeToString(a.SAddress.Address), "Address not as expected for SAddress")
	assert.Equal(t, addrtype.IPv4, a.RAddress.AddrType, "RAddress type not as expected")
	assert.Equal(t, "12d00023", hex.EncodeToString(a.RAddress.Address), "Address not as expected for RAddress")
}

func TestUnmarshalEncCredPart_optionalsNULL(t *testing.T) {
	t.Parallel()
	var a EncKrbCredPart
	b, err := hex.DecodeString(testdata.MarshaledKRB5enc_cred_partOptionalsNULL)
	if err != nil {
		t.Fatalf("Test vector read error: %v", err)
	}
	err = a.Unmarshal(b)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	//Parse the test time value into a time.Time type
	tt, _ := time.Parse(testdata.TEST_TIME_FORMAT, testdata.TEST_TIME)

	assert.Equal(t, 2, len(a.TicketInfo), "Number of ticket info items not as expected")
	//1st Ticket
	i := 0
	assert.Equal(t, int32(1), a.TicketInfo[i].Key.KeyType, fmt.Sprintf("Key type not as expected in ticket info item %d", i+1))
	assert.Equal(t, []byte("12345678"), a.TicketInfo[i].Key.KeyValue, fmt.Sprintf("Key value not as expected in ticket info item %d", i+1))

	//2nd Ticket
	i = 1
	assert.Equal(t, int32(1), a.TicketInfo[i].Key.KeyType, fmt.Sprintf("Key type not as expected in ticket info item %d", i+1))
	assert.Equal(t, []byte("12345678"), a.TicketInfo[i].Key.KeyValue, fmt.Sprintf("Key value not as expected in ticket info item %d", i+1))
	assert.Equal(t, testdata.TEST_REALM, a.TicketInfo[i].PRealm, fmt.Sprintf("PRealm not as expected on ticket info item %d", i+1))
	assert.Equal(t, nametype.KRB_NT_PRINCIPAL, a.TicketInfo[i].PName.NameType, fmt.Sprintf("Ticket info (%v) PName NameType not as expected", i+1))
	assert.Equal(t, len(testdata.TEST_PRINCIPALNAME_NAMESTRING), len(a.TicketInfo[i].PName.NameString), fmt.Sprintf("Ticket info (%v) PName does not have the expected number of NameStrings", i+1))
	assert.Equal(t, testdata.TEST_PRINCIPALNAME_NAMESTRING, a.TicketInfo[i].PName.NameString, fmt.Sprintf("Ticket info (%v) PName name string entries not as expected", i+1))
	assert.Equal(t, "fedcba98", hex.EncodeToString(a.TicketInfo[i].Flags.Bytes), fmt.Sprintf("Flags not as expected on ticket info %d", i+1))
	assert.Equal(t, tt, a.TicketInfo[i].AuthTime, fmt.Sprintf("Auth time value not as expected for ticket info %d", i+1))
	assert.Equal(t, tt, a.TicketInfo[i].StartTime, fmt.Sprintf("Start time value not as expected for ticket info %d", i+1))
	assert.Equal(t, tt, a.TicketInfo[i].EndTime, fmt.Sprintf("End time value not as expected for ticket info %d", i+1))
	assert.Equal(t, tt, a.TicketInfo[i].RenewTill, fmt.Sprintf("Renew Till time value not as expected for ticket info %d", i+1))
	assert.Equal(t, nametype.KRB_NT_PRINCIPAL, a.TicketInfo[i].SName.NameType, fmt.Sprintf("Ticket info (%v) PName NameType not as expected", i+1))
	assert.Equal(t, len(testdata.TEST_PRINCIPALNAME_NAMESTRING), len(a.TicketInfo[i].SName.NameString), fmt.Sprintf("Ticket info (%v) PName does not have the expected number of NameStrings", i+1))
	assert.Equal(t, testdata.TEST_PRINCIPALNAME_NAMESTRING, a.TicketInfo[i].SName.NameString, fmt.Sprintf("Ticket info (%v) PName name string entries not as expected", i+1))
	assert.Equal(t, 2, len(a.TicketInfo[i].CAddr), "Number of client addresses not as expected")
	for j, addr := range a.TicketInfo[i].CAddr {
		assert.Equal(t, addrtype.IPv4, addr.AddrType, fmt.Sprintf("Host address type not as expected for address item %d within ticket info %d", j+1, i+1))
		assert.Equal(t, "12d00023", hex.EncodeToString(addr.Address), fmt.Sprintf("Host address not as expected for address item %d within ticket info %d", j+1, i+1))
	}
}
