package spnego

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/jcmturner/goidentity/v6"
	"github.com/jcmturner/gokrb5/v8/client"
	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/keytab"
	"github.com/jcmturner/gokrb5/v8/service"
	"github.com/jcmturner/gokrb5/v8/test"
	"github.com/jcmturner/gokrb5/v8/test/testdata"
	"github.com/stretchr/testify/assert"
)

func TestClient_SetSPNEGOHeader(t *testing.T) {
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
	l := log.New(os.Stderr, "SPNEGO Client:", log.LstdFlags)
	cl := client.NewWithKeytab("testuser1", "TEST.GOKRB5", kt, c, client.Logger(l))

	err := cl.Login()
	if err != nil {
		t.Fatalf("error on AS_REQ: %v\n", err)
	}
	urls := []string{
		"http://cname.test.gokrb5",
		"http://host.test.gokrb5",
	}
	paths := []string{
		"/modkerb/index.html",
		//"/modgssapi/index.html",
	}
	for _, url := range urls {
		for _, p := range paths {
			r, _ := http.NewRequest("GET", url+p, nil)
			httpResp, err := http.DefaultClient.Do(r)
			if err != nil {
				t.Fatalf("%s request error: %v", url+p, err)
			}
			assert.Equal(t, http.StatusUnauthorized, httpResp.StatusCode, "Status code in response to client with no SPNEGO not as expected")

			err = SetSPNEGOHeader(cl, r, "")
			if err != nil {
				t.Fatalf("error setting client SPNEGO header: %v", err)
			}

			httpResp, err = http.DefaultClient.Do(r)
			if err != nil {
				t.Fatalf("%s request error: %v\n", url+p, err)
			}
			assert.Equal(t, http.StatusOK, httpResp.StatusCode, "Status code in response to client SPNEGO request not as expected")
		}
	}
}

func TestSPNEGOHTTPClient(t *testing.T) {
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
	l := log.New(os.Stderr, "SPNEGO Client:", log.LstdFlags)
	cl := client.NewWithKeytab("testuser1", "TEST.GOKRB5", kt, c, client.Logger(l))

	err := cl.Login()
	if err != nil {
		t.Fatalf("error on AS_REQ: %v\n", err)
	}
	urls := []string{
		"http://cname.test.gokrb5",
		"http://host.test.gokrb5",
	}
	// This path issues a redirect which the http client will automatically follow.
	// It should cause a replay issue if the negInit token is sent in the first instance.
	paths := []string{
		"/modgssapi", // This issues a redirect which the http client will automatically follow. Could cause a replay issue
		"/redirect",
	}
	for _, url := range urls {
		for _, p := range paths {
			r, _ := http.NewRequest("GET", url+p, nil)
			httpCl := http.DefaultClient
			httpCl.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				t.Logf("http client redirect: %+v", *req)
				return nil
			}
			spnegoCl := NewClient(cl, httpCl, "")
			httpResp, err := spnegoCl.Do(r)
			if err != nil {
				t.Fatalf("%s request error: %v", url+p, err)
			}
			assert.Equal(t, http.StatusOK, httpResp.StatusCode, "Status code in response to client SPNEGO request not as expected")
		}
	}
}

func TestService_SPNEGOKRB_NoAuthHeader(t *testing.T) {
	s := httpServer()
	defer s.Close()
	r, _ := http.NewRequest("GET", s.URL, nil)
	httpResp, err := http.DefaultClient.Do(r)
	if err != nil {
		t.Fatalf("Request error: %v\n", err)
	}
	assert.Equal(t, http.StatusUnauthorized, httpResp.StatusCode, "Status code in response to client with no SPNEGO not as expected")
	assert.Equal(t, "Negotiate", httpResp.Header.Get("WWW-Authenticate"), "Negotiation header not set by server.")
}

func TestService_SPNEGOKRB_ValidUser(t *testing.T) {
	test.Integration(t)

	s := httpServer()
	defer s.Close()
	r, _ := http.NewRequest("GET", s.URL, nil)

	cl := getClient()
	err := SetSPNEGOHeader(cl, r, "HTTP/host.test.gokrb5")
	if err != nil {
		t.Fatalf("error setting client's SPNEGO header: %v", err)
	}

	httpResp, err := http.DefaultClient.Do(r)
	if err != nil {
		t.Fatalf("Request error: %v\n", err)
	}
	assert.Equal(t, http.StatusOK, httpResp.StatusCode, "Status code in response to client SPNEGO request not as expected")
}

func TestService_SPNEGOKRB_ValidUser_RawKRB5Token(t *testing.T) {
	test.Integration(t)

	s := httpServer()
	defer s.Close()
	r, _ := http.NewRequest("GET", s.URL, nil)

	cl := getClient()
	sc := SPNEGOClient(cl, "HTTP/host.test.gokrb5")
	err := sc.AcquireCred()
	if err != nil {
		t.Fatalf("could not acquire client credential: %v", err)
	}
	st, err := sc.InitSecContext()
	if err != nil {
		t.Fatalf("could not initialize context: %v", err)
	}

	// Use the raw KRB5 context token
	nb := st.(*SPNEGOToken).NegTokenInit.MechTokenBytes
	hs := "Negotiate " + base64.StdEncoding.EncodeToString(nb)
	r.Header.Set(HTTPHeaderAuthRequest, hs)

	httpResp, err := http.DefaultClient.Do(r)
	if err != nil {
		t.Fatalf("Request error: %v\n", err)
	}
	assert.Equal(t, http.StatusOK, httpResp.StatusCode, "Status code in response to client SPNEGO request not as expected")
}

func TestService_SPNEGOKRB_Replay(t *testing.T) {
	test.Integration(t)

	s := httpServerWithoutSessionManager()
	defer s.Close()
	r1, _ := http.NewRequest("GET", s.URL, nil)

	cl := getClient()
	err := SetSPNEGOHeader(cl, r1, "HTTP/host.test.gokrb5")
	if err != nil {
		t.Fatalf("error setting client's SPNEGO header: %v", err)
	}

	// First request with this ticket should be accepted
	httpResp, err := http.DefaultClient.Do(r1)
	if err != nil {
		t.Fatalf("Request error: %v\n", err)
	}
	assert.Equal(t, http.StatusOK, httpResp.StatusCode, "Status code in response to client SPNEGO request not as expected")

	// Use ticket again should be rejected
	httpResp, err = http.DefaultClient.Do(r1)
	if err != nil {
		t.Fatalf("Request error: %v\n", err)
	}
	assert.Equal(t, http.StatusUnauthorized, httpResp.StatusCode, "Status code in response to client with no SPNEGO not as expected. Expected a replay to be detected.")

	// Form a 2nd ticket
	r2, _ := http.NewRequest("GET", s.URL, nil)

	err = SetSPNEGOHeader(cl, r2, "HTTP/host.test.gokrb5")
	if err != nil {
		t.Fatalf("error setting client's SPNEGO header: %v", err)
	}

	// First use of 2nd ticket should be accepted
	httpResp, err = http.DefaultClient.Do(r2)
	if err != nil {
		t.Fatalf("Request error: %v\n", err)
	}
	assert.Equal(t, http.StatusOK, httpResp.StatusCode, "Status code in response to client SPNEGO request not as expected")

	// Using the 1st ticket again should still be rejected
	httpResp, err = http.DefaultClient.Do(r1)
	if err != nil {
		t.Fatalf("Request error: %v\n", err)
	}
	assert.Equal(t, http.StatusUnauthorized, httpResp.StatusCode, "Status code in response to client with no SPNEGO not as expected. Expected a replay to be detected.")

	// Using the 2nd again should be rejected as replay
	httpResp, err = http.DefaultClient.Do(r2)
	if err != nil {
		t.Fatalf("Request error: %v\n", err)
	}
	assert.Equal(t, http.StatusUnauthorized, httpResp.StatusCode, "Status code in response to client with no SPNEGO not as expected. Expected a replay to be detected.")
}

func TestService_SPNEGOKRB_ReplayCache_Concurrency(t *testing.T) {
	test.Integration(t)

	s := httpServerWithoutSessionManager()
	defer s.Close()
	r1, _ := http.NewRequest("GET", s.URL, nil)

	cl := getClient()
	err := SetSPNEGOHeader(cl, r1, "HTTP/host.test.gokrb5")
	if err != nil {
		t.Fatalf("error setting client's SPNEGO header: %v", err)
	}
	r1h := r1.Header.Get(HTTPHeaderAuthRequest)

	r2, _ := http.NewRequest("GET", s.URL, nil)

	err = SetSPNEGOHeader(cl, r2, "HTTP/host.test.gokrb5")
	if err != nil {
		t.Fatalf("error setting client's SPNEGO header: %v", err)
	}
	r2h := r2.Header.Get(HTTPHeaderAuthRequest)

	// Concurrent 1st requests should be OK
	var wg sync.WaitGroup
	wg.Add(2)
	go httpGet(r1, &wg)
	go httpGet(r2, &wg)
	wg.Wait()

	// A number of concurrent requests with the same ticket should be rejected due to replay
	var wg2 sync.WaitGroup
	noReq := 10
	wg2.Add(noReq * 2)
	for i := 0; i < noReq; i++ {
		rr1, _ := http.NewRequest("GET", s.URL, nil)
		rr1.Header.Set(HTTPHeaderAuthRequest, r1h)
		rr2, _ := http.NewRequest("GET", s.URL, nil)
		rr2.Header.Set(HTTPHeaderAuthRequest, r2h)
		go httpGet(rr1, &wg2)
		go httpGet(rr2, &wg2)
	}
	wg2.Wait()
}

func TestService_SPNEGOKRB_Upload(t *testing.T) {
	test.Integration(t)

	s := httpServer()
	defer s.Close()

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", "testfile.bin")
	if err != nil {
		t.Fatalf("error writing to buffer: %v", err)
	}

	data := make([]byte, 10240)
	rand.Read(data)
	br := bytes.NewReader(data)
	_, err = io.Copy(fileWriter, br)
	if err != nil {
		t.Fatalf("error copying bytes: %v", err)
	}
	bodyWriter.Close()

	r, _ := http.NewRequest("POST", s.URL, bodyBuf)
	r.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	cl := getClient()
	cookieJar, _ := cookiejar.New(nil)
	httpCl := http.DefaultClient
	httpCl.Jar = cookieJar
	spnegoCl := NewClient(cl, httpCl, "HTTP/host.test.gokrb5")
	httpResp, err := spnegoCl.Do(r)
	if err != nil {
		t.Fatalf("Request error: %v\n", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(httpResp.Body)
		bodyString := string(bodyBytes)
		httpResp.Body.Close()
		t.Errorf("unexpected code from http server (%d): %s", httpResp.StatusCode, bodyString)
	}
}

func httpGet(r *http.Request, wg *sync.WaitGroup) {
	defer wg.Done()
	http.DefaultClient.Do(r)
}

func httpServerWithoutSessionManager() *httptest.Server {
	l := log.New(os.Stderr, "GOKRB5 Service Tests: ", log.LstdFlags)
	b, _ := hex.DecodeString(testdata.HTTP_KEYTAB)
	kt := keytab.New()
	kt.Unmarshal(b)
	th := http.HandlerFunc(testAppHandler)
	s := httptest.NewServer(SPNEGOKRB5Authenticate(th, kt, service.Logger(l)))
	return s
}

func httpServer() *httptest.Server {
	l := log.New(os.Stderr, "GOKRB5 Service Tests: ", log.LstdFlags)
	b, _ := hex.DecodeString(testdata.HTTP_KEYTAB)
	kt := keytab.New()
	kt.Unmarshal(b)
	th := http.HandlerFunc(testAppHandler)
	s := httptest.NewServer(SPNEGOKRB5Authenticate(th, kt, service.Logger(l), service.SessionManager(NewSessionMgr("gokrb5"))))
	return s
}

func testAppHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		maxUploadSize := int64(11240)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			http.Error(w, fmt.Sprintf("cannot parse multipart form: %v", err), http.StatusBadRequest)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		file, _, err := r.FormFile("uploadfile")
		if err != nil {
			http.Error(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// write out to /dev/null
		_, err = io.Copy(ioutil.Discard, file)
		if err != nil {
			http.Error(w, "WRITE_ERR", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	id := goidentity.FromHTTPRequestContext(r)
	fmt.Fprintf(w, "<html>\nTEST.GOKRB5 Handler\nAuthenticed user: %s\nUser's realm: %s\n</html>",
		id.UserName(),
		id.Domain())
	return
}

func getClient() *client.Client {
	b, _ := hex.DecodeString(testdata.KEYTAB_TESTUSER1_TEST_GOKRB5)
	kt := keytab.New()
	kt.Unmarshal(b)
	c, _ := config.NewFromString(testdata.KRB5_CONF)
	c.LibDefaults.NoAddresses = true
	addr := os.Getenv("TEST_KDC_ADDR")
	if addr == "" {
		addr = testdata.KDC_IP_TEST_GOKRB5
	}
	c.Realms[0].KDC = []string{addr + ":" + testdata.KDC_PORT_TEST_GOKRB5}
	c.Realms[0].KPasswdServer = []string{addr + ":464"}
	cl := client.NewWithKeytab("testuser1", "TEST.GOKRB5", kt, c)
	return cl
}

type SessionMgr struct {
	skey       []byte
	store      sessions.Store
	cookieName string
}

func NewSessionMgr(cookieName string) SessionMgr {
	skey := []byte("thisistestsecret") // Best practice is to load this key from a secure location.
	return SessionMgr{
		skey:       skey,
		store:      sessions.NewCookieStore(skey),
		cookieName: cookieName,
	}
}

func (smgr SessionMgr) Get(r *http.Request, k string) ([]byte, error) {
	s, err := smgr.store.Get(r, smgr.cookieName)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, errors.New("nil session")
	}
	b, ok := s.Values[k].([]byte)
	if !ok {
		return nil, fmt.Errorf("could not get bytes held in session at %s", k)
	}
	return b, nil
}

func (smgr SessionMgr) New(w http.ResponseWriter, r *http.Request, k string, v []byte) error {
	s, err := smgr.store.New(r, smgr.cookieName)
	if err != nil {
		return fmt.Errorf("could not get new session from session manager: %v", err)
	}
	s.Values[k] = v
	return s.Save(r, w)
}
