package client

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/iana/etypeID"
	"github.com/jcmturner/gokrb5/v8/keytab"
	"github.com/jcmturner/gokrb5/v8/test"
	"github.com/jcmturner/gokrb5/v8/test/testdata"
	"github.com/stretchr/testify/assert"
)

func TestMultiThreadedClientSession(t *testing.T) {
	test.Integration(t)

	b, _ := hex.DecodeString(testdata.KEYTAB_TESTUSER1_TEST_GOKRB5)
	kt := keytab.New()
	kt.Unmarshal(b)
	c, _ := config.NewFromString(testdata.KRB5_CONF)
	addr := os.Getenv("TEST_KDC_ADDR")
	if addr == "" {
		addr = testdata.KDC_IP_TEST_GOKRB5
	}
	c.Realms[0].KDC = []string{addr + ":" + testdata.KDC_PORT_TEST_GOKRB5}
	cl := NewWithKeytab("testuser1", "TEST.GOKRB5", kt, c)
	err := cl.Login()
	if err != nil {
		t.Fatalf("failed to log in: %v", err)
	}

	s, ok := cl.sessions.get("TEST.GOKRB5")
	if !ok {
		t.Fatal("error initially getting session")
	}
	go func() {
		for {
			err := cl.renewTGT(s)
			if err != nil {
				t.Logf("error renewing TGT: %v", err)
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			tgt, _, err := cl.sessionTGT("TEST.GOKRB5")
			if err != nil || tgt.Realm != "TEST.GOKRB5" {
				t.Logf("error getting session: %v", err)
			}
			_, _, _, r, _ := cl.sessionTimes("TEST.GOKRB5")
			fmt.Fprintf(ioutil.Discard, "%v", r)
		}()
		time.Sleep(time.Second)
	}
	wg.Wait()
}

func TestClient_AutoRenew_Goroutine(t *testing.T) {
	test.Integration(t)

	// Tests that the auto renew of client credentials is not spawning goroutines out of control.
	addr := os.Getenv("TEST_KDC_ADDR")
	if addr == "" {
		addr = testdata.KDC_IP_TEST_GOKRB5
	}
	b, _ := hex.DecodeString(testdata.KEYTAB_TESTUSER2_TEST_GOKRB5)
	kt := keytab.New()
	kt.Unmarshal(b)
	c, _ := config.NewFromString(testdata.KRB5_CONF)
	c.Realms[0].KDC = []string{addr + ":" + testdata.KDC_PORT_TEST_GOKRB5_SHORTTICKETS}
	c.LibDefaults.PreferredPreauthTypes = []int{int(etypeID.DES3_CBC_SHA1_KD)} // a preauth etype the KDC does not support. Test this does not cause renewal to fail.
	cl := NewWithKeytab("testuser2", "TEST.GOKRB5", kt, c)

	err := cl.Login()
	if err != nil {
		t.Errorf("error on logging in: %v\n", err)
	}
	n := runtime.NumGoroutine()
	for i := 0; i < 24; i++ {
		time.Sleep(time.Second * 5)
		_, endTime, _, _, err := cl.sessionTimes("TEST.GOKRB5")
		if err != nil {
			t.Errorf("could not get client's session: %v", err)
		}
		if time.Now().UTC().After(endTime) {
			t.Fatalf("session auto update failed")
		}
		spn := "HTTP/host.test.gokrb5"
		tkt, key, err := cl.GetServiceTicket(spn)
		if err != nil {
			t.Fatalf("error getting service ticket: %v\n", err)
		}
		b, _ := hex.DecodeString(testdata.HTTP_KEYTAB)
		skt := keytab.New()
		skt.Unmarshal(b)
		tkt.DecryptEncPart(skt, nil)
		assert.Equal(t, spn, tkt.SName.PrincipalNameString())
		assert.Equal(t, int32(18), key.KeyType)
		if runtime.NumGoroutine() > n {
			t.Fatalf("number of goroutines is increasing: should not be more than %d, is %d", n, runtime.NumGoroutine())
		}
	}
}

func TestSessions_JSON(t *testing.T) {
	s := &sessions{
		Entries: make(map[string]*session),
	}
	for i := 0; i < 3; i++ {
		realm := fmt.Sprintf("test%d", i)
		e := &session{
			realm:                realm,
			authTime:             time.Unix(int64(0+i), 0).UTC(),
			endTime:              time.Unix(int64(10+i), 0).UTC(),
			renewTill:            time.Unix(int64(20+i), 0).UTC(),
			sessionKeyExpiration: time.Unix(int64(30+i), 0).UTC(),
		}
		s.Entries[realm] = e
	}
	j, err := s.JSON()
	if err != nil {
		t.Errorf("error getting json: %v", err)
	}
	expected := `[
  {
    "Realm": "test0",
    "AuthTime": "1970-01-01T00:00:00Z",
    "EndTime": "1970-01-01T00:00:10Z",
    "RenewTill": "1970-01-01T00:00:20Z",
    "SessionKeyExpiration": "1970-01-01T00:00:30Z"
  },
  {
    "Realm": "test1",
    "AuthTime": "1970-01-01T00:00:01Z",
    "EndTime": "1970-01-01T00:00:11Z",
    "RenewTill": "1970-01-01T00:00:21Z",
    "SessionKeyExpiration": "1970-01-01T00:00:31Z"
  },
  {
    "Realm": "test2",
    "AuthTime": "1970-01-01T00:00:02Z",
    "EndTime": "1970-01-01T00:00:12Z",
    "RenewTill": "1970-01-01T00:00:22Z",
    "SessionKeyExpiration": "1970-01-01T00:00:32Z"
  }
]`
	assert.Equal(t, expected, j, "json output not as expected")
}
