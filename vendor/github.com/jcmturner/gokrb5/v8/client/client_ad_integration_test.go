package client

import (
	"bytes"
	"encoding/hex"
	"log"
	"testing"

	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/iana/etypeID"
	"github.com/jcmturner/gokrb5/v8/iana/nametype"
	"github.com/jcmturner/gokrb5/v8/keytab"
	"github.com/jcmturner/gokrb5/v8/test"
	"github.com/jcmturner/gokrb5/v8/test/testdata"
	"github.com/jcmturner/gokrb5/v8/types"
	"github.com/stretchr/testify/assert"
)

func TestClient_SuccessfulLogin_AD(t *testing.T) {
	test.AD(t)

	b, _ := hex.DecodeString(testdata.KEYTAB_TESTUSER1_USER_GOKRB5)
	kt := keytab.New()
	kt.Unmarshal(b)
	c, _ := config.NewFromString(testdata.KRB5_CONF_AD)
	cl := NewWithKeytab("testuser1", "USER.GOKRB5", kt, c, DisablePAFXFAST(true))

	err := cl.Login()
	if err != nil {
		t.Fatalf("Error on login: %v\n", err)
	}
}

func TestClient_SuccessfulLogin_AD_Without_PreAuth(t *testing.T) {
	test.AD(t)

	b, _ := hex.DecodeString(testdata.KEYTAB_TESTUSER3_USER_GOKRB5)
	kt := keytab.New()
	kt.Unmarshal(b)
	c, _ := config.NewFromString(testdata.KRB5_CONF_AD)
	cl := NewWithKeytab("testuser3", "USER.GOKRB5", kt, c, DisablePAFXFAST(true))

	err := cl.Login()
	if err != nil {
		t.Fatalf("Error on login: %v\n", err)
	}
}

func TestClient_GetServiceTicket_AD(t *testing.T) {
	test.AD(t)

	b, _ := hex.DecodeString(testdata.KEYTAB_TESTUSER1_USER_GOKRB5)
	kt := keytab.New()
	kt.Unmarshal(b)
	c, _ := config.NewFromString(testdata.KRB5_CONF_AD)
	cl := NewWithKeytab("testuser1", "USER.GOKRB5", kt, c)

	err := cl.Login()
	if err != nil {
		t.Fatalf("Error on login: %v\n", err)
	}
	spn := "HTTP/user2.user.gokrb5"
	tkt, key, err := cl.GetServiceTicket(spn)
	if err != nil {
		t.Fatalf("Error getting service ticket: %v\n", err)
	}
	assert.Equal(t, spn, tkt.SName.PrincipalNameString())
	assert.Equal(t, int32(18), key.KeyType)

	b, _ = hex.DecodeString(testdata.KEYTAB_TESTUSER2_USER_GOKRB5)
	skt := keytab.New()
	skt.Unmarshal(b)
	sname := types.PrincipalName{NameType: nametype.KRB_NT_PRINCIPAL, NameString: []string{"testuser2"}}
	err = tkt.DecryptEncPart(skt, &sname)
	if err != nil {
		t.Errorf("could not decrypt service ticket: %v", err)
	}
	w := bytes.NewBufferString("")
	l := log.New(w, "", 0)
	isPAC, pac, err := tkt.GetPACType(skt, &sname, l)
	if err != nil {
		t.Log(w.String())
		t.Errorf("error getting PAC: %v", err)
	}
	assert.True(t, isPAC, "should have PAC")
	assert.Equal(t, "USER", pac.KerbValidationInfo.LogonDomainName.String(), "domain name in PAC not correct")
}

func TestClient_GetServiceTicket_AD_TRUST_USER_DOMAIN(t *testing.T) {
	test.AD(t)

	b, _ := hex.DecodeString(testdata.KEYTAB_TESTUSER1_USER_GOKRB5)
	kt := keytab.New()
	kt.Unmarshal(b)
	c, _ := config.NewFromString(testdata.KRB5_CONF_AD)
	c.LibDefaults.Canonicalize = true
	c.LibDefaults.DefaultTktEnctypes = []string{"rc4-hmac"}
	c.LibDefaults.DefaultTktEnctypeIDs = []int32{etypeID.ETypesByName["rc4-hmac"]}
	c.LibDefaults.DefaultTGSEnctypes = []string{"rc4-hmac"}
	c.LibDefaults.DefaultTGSEnctypeIDs = []int32{etypeID.ETypesByName["rc4-hmac"]}
	cl := NewWithKeytab("testuser1", "USER.GOKRB5", kt, c, DisablePAFXFAST(true))
	err := cl.Login()

	if err != nil {
		t.Fatalf("Error on login: %v\n", err)
	}
	spn := "HTTP/host.res.gokrb5"
	tkt, key, err := cl.GetServiceTicket(spn)
	if err != nil {
		t.Fatalf("Error getting service ticket: %v\n", err)
	}
	assert.Equal(t, spn, tkt.SName.PrincipalNameString())
	assert.Equal(t, etypeID.ETypesByName["rc4-hmac"], key.KeyType)

	b, _ = hex.DecodeString(testdata.KEYTAB_SYSHTTP_RES_GOKRB5)
	skt := keytab.New()
	skt.Unmarshal(b)
	sname := types.PrincipalName{NameType: nametype.KRB_NT_PRINCIPAL, NameString: []string{"sysHTTP"}}
	err = tkt.DecryptEncPart(skt, &sname)
	if err != nil {
		t.Errorf("error decrypting ticket with service keytab: %v", err)
	}
	w := bytes.NewBufferString("")
	l := log.New(w, "", 0)
	isPAC, pac, err := tkt.GetPACType(skt, &sname, l)
	if err != nil {
		t.Log(w.String())
		t.Errorf("error getting PAC: %v", err)
	}
	assert.True(t, isPAC, "Did not find PAC in service ticket")
	assert.Equal(t, "testuser1", pac.KerbValidationInfo.EffectiveName.Value, "PAC value not parsed")

}

func TestClient_GetServiceTicket_AD_USER_DOMAIN(t *testing.T) {
	test.AD(t)

	b, _ := hex.DecodeString(testdata.KEYTAB_TESTUSER1_USER_GOKRB5)
	kt := keytab.New()
	kt.Unmarshal(b)
	c, _ := config.NewFromString(testdata.KRB5_CONF_AD)
	c.LibDefaults.Canonicalize = true
	c.LibDefaults.DefaultTktEnctypes = []string{"rc4-hmac"}
	c.LibDefaults.DefaultTktEnctypeIDs = []int32{etypeID.ETypesByName["rc4-hmac"]}
	c.LibDefaults.DefaultTGSEnctypes = []string{"rc4-hmac"}
	c.LibDefaults.DefaultTGSEnctypeIDs = []int32{etypeID.ETypesByName["rc4-hmac"]}
	cl := NewWithKeytab("testuser1", "USER.GOKRB5", kt, c, DisablePAFXFAST(true))

	err := cl.Login()

	if err != nil {
		t.Fatalf("Error on login: %v\n", err)
	}
	spn := "HTTP/user2.user.gokrb5"
	tkt, _, err := cl.GetServiceTicket(spn)
	if err != nil {
		t.Fatalf("Error getting service ticket: %v\n", err)
	}
	assert.Equal(t, spn, tkt.SName.PrincipalNameString())
	//assert.Equal(t, etypeID.ETypesByName["rc4-hmac"], key.KeyType)

	b, _ = hex.DecodeString(testdata.KEYTAB_TESTUSER2_USER_GOKRB5)
	skt := keytab.New()
	skt.Unmarshal(b)
	sname := types.PrincipalName{NameType: nametype.KRB_NT_PRINCIPAL, NameString: []string{"testuser2"}}
	err = tkt.DecryptEncPart(skt, &sname)
	if err != nil {
		t.Errorf("error decrypting ticket with service keytab: %v", err)
	}
	w := bytes.NewBufferString("")
	l := log.New(w, "", 0)
	isPAC, pac, err := tkt.GetPACType(skt, &sname, l)
	if err != nil {
		t.Log(w.String())
		t.Errorf("error getting PAC: %v", err)
	}
	assert.True(t, isPAC, "Did not find PAC in service ticket")
	assert.Equal(t, "testuser1", pac.KerbValidationInfo.EffectiveName.Value, "PAC value not parsed")

}
