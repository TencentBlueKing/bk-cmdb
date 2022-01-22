package gzhttp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/klauspost/compress/gzip"
)

var (
	smallTestBody = []byte("aaabbcaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbc")
	testBody      = []byte("aaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbccc aaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbccc aaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbccc aaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbccc aaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbccc aaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbccc aaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbccc")
)

func TestParseEncodings(t *testing.T) {
	examples := map[string]codings{

		// Examples from RFC 2616
		"compress, gzip":                     {"compress": 1.0, "gzip": 1.0},
		",,,,":                               {},
		"":                                   {},
		"*":                                  {"*": 1.0},
		"compress;q=0.5, gzip;q=1.0":         {"compress": 0.5, "gzip": 1.0},
		"gzip;q=1.0, identity; q=0.5, *;q=0": {"gzip": 1.0, "identity": 0.5, "*": 0.0},

		// More random stuff
		"AAA;q=1":     {"aaa": 1.0},
		"BBB ; q = 2": {"bbb": 1.0},
	}

	for eg, exp := range examples {
		t.Run(eg, func(t *testing.T) {
			act, _ := parseEncodings(eg)
			assertEqual(t, exp, act)
			gz := parseEncodingGzip(eg)
			assertEqual(t, exp["gzip"], gz)
		})
	}
}

func TestMustNewGzipHandler(t *testing.T) {
	// This just exists to provide something for GzipHandler to wrap.
	handler := newTestHandler(testBody)

	// requests without accept-encoding are passed along as-is

	req1, _ := http.NewRequest("GET", "/whatever", nil)
	resp1 := httptest.NewRecorder()
	handler.ServeHTTP(resp1, req1)
	res1 := resp1.Result()

	assertEqual(t, 200, res1.StatusCode)
	assertEqual(t, "", res1.Header.Get("Content-Encoding"))
	assertEqual(t, "Accept-Encoding", res1.Header.Get("Vary"))
	assertEqual(t, testBody, resp1.Body.Bytes())

	// but requests with accept-encoding:gzip are compressed if possible

	req2, _ := http.NewRequest("GET", "/whatever", nil)
	req2.Header.Set("Accept-Encoding", "gzip")
	resp2 := httptest.NewRecorder()
	handler.ServeHTTP(resp2, req2)
	res2 := resp2.Result()

	assertEqual(t, 200, res2.StatusCode)
	assertEqual(t, "gzip", res2.Header.Get("Content-Encoding"))
	assertEqual(t, "Accept-Encoding", res2.Header.Get("Vary"))
	assertEqual(t, gzipStrLevel(testBody, gzip.DefaultCompression), resp2.Body.Bytes())

	// content-type header is correctly set based on uncompressed body

	req3, _ := http.NewRequest("GET", "/whatever", nil)
	req3.Header.Set("Accept-Encoding", "gzip")
	res3 := httptest.NewRecorder()
	handler.ServeHTTP(res3, req3)

	assertEqual(t, http.DetectContentType([]byte(testBody)), res3.Header().Get("Content-Type"))
}

func TestGzipHandlerSmallBodyNoCompression(t *testing.T) {
	handler := newTestHandler(smallTestBody)

	req, _ := http.NewRequest("GET", "/whatever", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	res := resp.Result()

	// with less than 1400 bytes the response should not be gzipped

	assertEqual(t, 200, res.StatusCode)
	assertEqual(t, "", res.Header.Get("Content-Encoding"))
	assertEqual(t, "Accept-Encoding", res.Header.Get("Vary"))
	assertEqual(t, smallTestBody, resp.Body.Bytes())

}

func TestGzipHandlerAlreadyCompressed(t *testing.T) {
	handler := newTestHandler(testBody)

	req, _ := http.NewRequest("GET", "/gzipped", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	assertEqual(t, testBody, res.Body.Bytes())
}

func TestGzipHandlerRangeReply(t *testing.T) {
	handler := GzipHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Range", "bytes 0-300/804")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(testBody))
		}))
	req, _ := http.NewRequest("GET", "/gzipped", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	res := resp.Result()
	assertEqual(t, 200, res.StatusCode)
	assertEqual(t, "", res.Header.Get("Content-Encoding"))
	assertEqual(t, testBody, resp.Body.Bytes())
}

func TestGzipHandlerAcceptRange(t *testing.T) {
	handler := GzipHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Accept-Ranges", "bytes")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(testBody))
		}))
	req, _ := http.NewRequest("GET", "/gzipped", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	res := resp.Result()
	assertEqual(t, 200, res.StatusCode)
	assertEqual(t, "gzip", res.Header.Get("Content-Encoding"))
	assertEqual(t, "", res.Header.Get("Accept-Ranges"))
	zr, err := gzip.NewReader(resp.Body)
	assertNil(t, err)
	got, err := ioutil.ReadAll(zr)
	assertNil(t, err)
	assertEqual(t, testBody, got)
}

func TestGzipHandlerKeepAcceptRange(t *testing.T) {
	wrapper, err := NewWrapper(KeepAcceptRanges())
	assertNil(t, err)
	handler := wrapper(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Accept-Ranges", "bytes")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(testBody))
		}))
	req, _ := http.NewRequest("GET", "/gzipped", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	res := resp.Result()
	assertEqual(t, 200, res.StatusCode)
	assertEqual(t, "gzip", res.Header.Get("Content-Encoding"))
	assertEqual(t, "bytes", res.Header.Get("Accept-Ranges"))
	zr, err := gzip.NewReader(resp.Body)
	assertNil(t, err)
	got, err := ioutil.ReadAll(zr)
	assertNil(t, err)
	assertEqual(t, testBody, got)
}

func TestNewGzipLevelHandler(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(testBody)
	})

	for lvl := gzip.StatelessCompression; lvl <= gzip.BestCompression; lvl++ {
		t.Run(fmt.Sprint(lvl), func(t *testing.T) {
			wrapper, err := NewWrapper(CompressionLevel(lvl))
			assertNil(t, err)

			req, _ := http.NewRequest("GET", "/whatever", nil)
			req.Header.Set("Accept-Encoding", "gzip")
			resp := httptest.NewRecorder()
			wrapper(handler).ServeHTTP(resp, req)
			res := resp.Result()

			assertEqual(t, 200, res.StatusCode)
			assertEqual(t, "gzip", res.Header.Get("Content-Encoding"))
			assertEqual(t, "Accept-Encoding", res.Header.Get("Vary"))
			got := gzipStrLevel(testBody, lvl)
			if lvl != gzip.StatelessCompression {
				assertEqual(t, got, resp.Body.Bytes())
			}
			t.Log(lvl, len(got))
		})
	}
}

func TestNewGzipLevelHandlerReturnsErrorForInvalidLevels(t *testing.T) {
	var err error
	_, err = NewWrapper(CompressionLevel(-42))
	assertNotNil(t, err)

	_, err = NewWrapper(CompressionLevel(42))
	assertNotNil(t, err)
}

func TestMustNewGzipLevelHandlerWillPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("panic was called with", r)
		}
	}()

	_ = GzipHandler(nil)
}

func TestGzipHandlerNoBody(t *testing.T) {
	tests := []struct {
		statusCode      int
		contentEncoding string
		emptyBody       bool
		body            []byte
	}{
		// Body must be empty.
		{http.StatusNoContent, "", true, nil},
		{http.StatusNotModified, "", true, nil},
		// Body is going to get gzip'd no matter what.
		{http.StatusOK, "", true, []byte{}},
		{http.StatusOK, "gzip", false, []byte(testBody)},
	}

	for num, test := range tests {
		t.Run(fmt.Sprintf("test-%d", num), func(t *testing.T) {
			handler := GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.statusCode)
				if test.body != nil {
					w.Write(test.body)
				}
			}))

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Accept-Encoding", "gzip")
			handler.ServeHTTP(rec, req)

			body, err := ioutil.ReadAll(rec.Body)
			if err != nil {
				t.Fatalf("Unexpected error reading response body: %v", err)
			}

			header := rec.Header()
			assertEqual(t, test.contentEncoding, header.Get("Content-Encoding"))
			assertEqual(t, "Accept-Encoding", header.Get("Vary"))
			if test.emptyBody {
				assertEqual(t, 0, len(body))
			} else {
				assertNotEqual(t, 0, len(body))
				assertNotEqual(t, test.body, body)
			}
		})

	}
}

func TestGzipHandlerContentLength(t *testing.T) {
	testBodyBytes := []byte(testBody)
	tests := []struct {
		bodyLen   int
		bodies    [][]byte
		emptyBody bool
	}{
		{len(testBody), [][]byte{testBodyBytes}, false},
		// each of these writes is less than the DefaultMinSize
		{len(testBody), [][]byte{testBodyBytes[:200], testBodyBytes[200:]}, false},
		// without a defined Content-Length it should still gzip
		{0, [][]byte{testBodyBytes[:200], testBodyBytes[200:]}, false},
		// simulate a HEAD request with an empty write (to populate headers)
		{len(testBody), [][]byte{nil}, true},
	}

	// httptest.NewRecorder doesn't give you access to the Content-Length
	// header so instead, we create a server on a random port and make
	// a request to that instead
	ln, err := net.Listen("tcp", "localhost:")
	if err != nil {
		t.Fatalf("failed creating listen socket: %v", err)
	}
	defer ln.Close()
	srv := &http.Server{
		Handler: nil,
	}
	go srv.Serve(ln)

	for num, test := range tests {
		t.Run(fmt.Sprintf("test-%d", num), func(t *testing.T) {
			srv.Handler = GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if test.bodyLen > 0 {
					w.Header().Set("Content-Length", strconv.Itoa(test.bodyLen))
				}
				for _, b := range test.bodies {
					w.Write(b)
				}
			}))
			req := &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: "/", Scheme: "http", Host: ln.Addr().String()},
				Header: make(http.Header),
				Close:  true,
			}
			req.Header.Set("Accept-Encoding", "gzip")
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Unexpected error making http request in test iteration %d: %v", num, err)
			}
			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Unexpected error reading response body in test iteration %d: %v", num, err)
			}

			l, err := strconv.Atoi(res.Header.Get("Content-Length"))
			if err != nil {
				t.Fatalf("Unexpected error parsing Content-Length in test iteration %d: %v", num, err)
			}
			if test.emptyBody {
				assertEqual(t, 0, len(body))
				assertEqual(t, 0, l)
			} else {
				assertEqual(t, len(body), l)
			}
			assertEqual(t, "gzip", res.Header.Get("Content-Encoding"))
			assertNotEqual(t, test.bodyLen, l)
		})
	}
}

func TestGzipHandlerMinSizeMustBePositive(t *testing.T) {
	_, err := NewWrapper(MinSize(-1))
	assertNotNil(t, err)
}

func TestGzipHandlerMinSize(t *testing.T) {
	responseLength := 0
	b := []byte{'x'}

	wrapper, _ := NewWrapper(MinSize(128))
	handler := wrapper(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Write responses one byte at a time to ensure that the flush
			// mechanism, if used, is working properly.
			for i := 0; i < responseLength; i++ {
				n, err := w.Write(b)
				assertEqual(t, 1, n)
				assertNil(t, err)
			}
		},
	))

	r, _ := http.NewRequest("GET", "/whatever", &bytes.Buffer{})
	r.Header.Add("Accept-Encoding", "gzip")

	// Short response is not compressed
	responseLength = 127
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	if w.Result().Header.Get(contentEncoding) == "gzip" {
		t.Error("Expected uncompressed response, got compressed")
	}

	// Long response is not compressed
	responseLength = 128
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	if w.Result().Header.Get(contentEncoding) != "gzip" {
		t.Error("Expected compressed response, got uncompressed")
	}
}

type panicOnSecondWriteHeaderWriter struct {
	http.ResponseWriter
	headerWritten bool
}

func (w *panicOnSecondWriteHeaderWriter) WriteHeader(s int) {
	if w.headerWritten {
		panic("header already written")
	}
	w.headerWritten = true
	w.ResponseWriter.WriteHeader(s)
}

func TestGzipHandlerDoubleWriteHeader(t *testing.T) {
	handler := GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "15000")
		// Specifically write the header here
		w.WriteHeader(304)
		// Ensure that after a Write the header isn't triggered again on close
		w.Write(nil)
	}))
	wrapper := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w = &panicOnSecondWriteHeaderWriter{
			ResponseWriter: w,
		}
		handler.ServeHTTP(w, r)
	})

	rec := httptest.NewRecorder()
	// TODO: in Go1.7 httptest.NewRequest was introduced this should be used
	// once 1.6 is not longer supported.
	req := &http.Request{
		Method:     "GET",
		URL:        &url.URL{Path: "/"},
		Proto:      "HTTP/1.1",
		ProtoMinor: 1,
		RemoteAddr: "192.0.2.1:1234",
		Header:     make(http.Header),
	}
	req.Header.Set("Accept-Encoding", "gzip")
	wrapper.ServeHTTP(rec, req)
	body, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("Unexpected error reading response body: %v", err)
	}
	assertEqual(t, 0, len(body))
	header := rec.Header()
	assertEqual(t, "gzip", header.Get("Content-Encoding"))
	assertEqual(t, "Accept-Encoding", header.Get("Vary"))
	assertEqual(t, 304, rec.Code)
}

func TestStatusCodes(t *testing.T) {
	handler := GzipHandler(http.NotFoundHandler())
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	result := w.Result()
	if result.StatusCode != 404 {
		t.Errorf("StatusCode should have been 404 but was %d", result.StatusCode)
	}
}

func TestFlushBeforeWrite(t *testing.T) {
	b := []byte(testBody)
	handler := GzipHandler(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
		rw.(http.Flusher).Flush()
		rw.Write(b)
	}))
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	res := w.Result()
	assertEqual(t, http.StatusNotFound, res.StatusCode)
	assertEqual(t, "gzip", res.Header.Get("Content-Encoding"))
	assertNotEqual(t, b, w.Body.Bytes())
}

func TestImplementCloseNotifier(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(acceptEncoding, "gzip")
	GzipHandler(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, ok := rw.(http.CloseNotifier)
		// response writer must implement http.CloseNotifier
		assertEqual(t, true, ok)
	})).ServeHTTP(&mockRWCloseNotify{}, request)
}

func TestImplementFlusherAndCloseNotifier(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(acceptEncoding, "gzip")
	GzipHandler(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, okCloseNotifier := rw.(http.CloseNotifier)
		// response writer must implement http.CloseNotifier
		assertEqual(t, true, okCloseNotifier)
		_, okFlusher := rw.(http.Flusher)
		// "response writer must implement http.Flusher"
		assertEqual(t, true, okFlusher)
	})).ServeHTTP(&mockRWCloseNotify{}, request)
}

func TestNotImplementCloseNotifier(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(acceptEncoding, "gzip")
	GzipHandler(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, ok := rw.(http.CloseNotifier)
		// response writer must not implement http.CloseNotifier
		assertEqual(t, false, ok)
	})).ServeHTTP(httptest.NewRecorder(), request)
}

type mockRWCloseNotify struct{}

func (m *mockRWCloseNotify) CloseNotify() <-chan bool {
	panic("implement me")
}

func (m *mockRWCloseNotify) Header() http.Header {
	return http.Header{}
}

func (m *mockRWCloseNotify) Write([]byte) (int, error) {
	panic("implement me")
}

func (m *mockRWCloseNotify) WriteHeader(int) {
	panic("implement me")
}

func TestIgnoreSubsequentWriteHeader(t *testing.T) {
	handler := GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.WriteHeader(404)
	}))
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	result := w.Result()
	if result.StatusCode != 500 {
		t.Errorf("StatusCode should have been 500 but was %d", result.StatusCode)
	}
}

func TestDontWriteWhenNotWrittenTo(t *testing.T) {
	// When using gzip as middleware without ANY writes in the handler,
	// ensure the gzip middleware doesn't touch the actual ResponseWriter
	// either.

	handler0 := GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	handler1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler0.ServeHTTP(w, r)
		w.WriteHeader(404) // this only works if gzip didn't do a WriteHeader(200)
	})

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	handler1.ServeHTTP(w, r)

	result := w.Result()
	if result.StatusCode != 404 {
		t.Errorf("StatusCode should have been 404 but was %d", result.StatusCode)
	}
}

var contentTypeTests = []struct {
	name                 string
	contentType          string
	acceptedContentTypes []string
	expectedGzip         bool
}{
	{
		name:                 "Always gzip when content types are empty",
		contentType:          "",
		acceptedContentTypes: []string{},
		expectedGzip:         true,
	},
	{
		name:                 "MIME match",
		contentType:          "application/json",
		acceptedContentTypes: []string{"application/json"},
		expectedGzip:         true,
	},
	{
		name:                 "MIME no match",
		contentType:          "text/xml",
		acceptedContentTypes: []string{"application/json"},
		expectedGzip:         false,
	},
	{
		name:                 "MIME match with no other directive ignores non-MIME directives",
		contentType:          "application/json; charset=utf-8",
		acceptedContentTypes: []string{"application/json"},
		expectedGzip:         true,
	},
	{
		name:                 "MIME match with other directives requires all directives be equal, different charset",
		contentType:          "application/json; charset=ascii",
		acceptedContentTypes: []string{"application/json; charset=utf-8"},
		expectedGzip:         false,
	},
	{
		name:                 "MIME match with other directives requires all directives be equal, same charset",
		contentType:          "application/json; charset=utf-8",
		acceptedContentTypes: []string{"application/json; charset=utf-8"},
		expectedGzip:         true,
	},
	{
		name:                 "MIME match with other directives requires all directives be equal, missing charset",
		contentType:          "application/json",
		acceptedContentTypes: []string{"application/json; charset=ascii"},
		expectedGzip:         false,
	},
	{
		name:                 "MIME match case insensitive",
		contentType:          "Application/Json",
		acceptedContentTypes: []string{"application/json"},
		expectedGzip:         true,
	},
	{
		name:                 "MIME match ignore whitespace",
		contentType:          "application/json;charset=utf-8",
		acceptedContentTypes: []string{"application/json;            charset=utf-8"},
		expectedGzip:         true,
	},
}

func TestContentTypes(t *testing.T) {
	for _, tt := range contentTypeTests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", tt.contentType)
				w.Write(testBody)
			})

			wrapper, err := NewWrapper(ContentTypes(tt.acceptedContentTypes))
			assertNil(t, err)

			req, _ := http.NewRequest("GET", "/whatever", nil)
			req.Header.Set("Accept-Encoding", "gzip")
			resp := httptest.NewRecorder()
			wrapper(handler).ServeHTTP(resp, req)
			res := resp.Result()

			assertEqual(t, 200, res.StatusCode)
			if tt.expectedGzip {
				assertEqual(t, "gzip", res.Header.Get("Content-Encoding"))
			} else {
				assertNotEqual(t, "gzip", res.Header.Get("Content-Encoding"))
			}
		})
		t.Run("not-"+tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", tt.contentType)
				w.Write(testBody)
			})

			wrapper, err := NewWrapper(ExceptContentTypes(tt.acceptedContentTypes))
			assertNil(t, err)

			req, _ := http.NewRequest("GET", "/whatever", nil)
			req.Header.Set("Accept-Encoding", "gzip")
			resp := httptest.NewRecorder()
			wrapper(handler).ServeHTTP(resp, req)
			res := resp.Result()

			assertEqual(t, 200, res.StatusCode)
			if !tt.expectedGzip {
				assertEqual(t, "gzip", res.Header.Get("Content-Encoding"))
			} else {
				assertNotEqual(t, "gzip", res.Header.Get("Content-Encoding"))
			}
		})
		t.Run("disable-"+tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", tt.contentType)
				w.Header().Set(HeaderNoCompression, "plz")
				w.Write(testBody)
			})

			wrapper, err := NewWrapper(ContentTypes(tt.acceptedContentTypes))
			assertNil(t, err)

			req, _ := http.NewRequest("GET", "/whatever", nil)
			req.Header.Set("Accept-Encoding", "gzip")
			resp := httptest.NewRecorder()
			wrapper(handler).ServeHTTP(resp, req)
			res := resp.Result()

			assertEqual(t, 200, res.StatusCode)
			assertNotEqual(t, "gzip", res.Header.Get("Content-Encoding"))
		})
	}
}

func TestFlush(t *testing.T) {
	for _, tt := range contentTypeTests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", tt.contentType)
				tb := testBody
				for len(tb) > 0 {
					// Write 100 bytes per run
					// Detection should not be affected (we send 100 bytes)
					toWrite := 100
					if toWrite > len(tb) {
						toWrite = len(tb)
					}
					_, err := w.Write(tb[:toWrite])
					if err != nil {
						t.Fatal(err)
					}
					// Flush between each write
					w.(http.Flusher).Flush()
					tb = tb[toWrite:]
				}
			})

			wrapper, err := NewWrapper(ContentTypes(tt.acceptedContentTypes))
			assertNil(t, err)

			req, _ := http.NewRequest("GET", "/whatever", nil)
			req.Header.Set("Accept-Encoding", "gzip")
			// This doesn't allow checking flushes, but we validate if content is correct.
			resp := httptest.NewRecorder()
			wrapper(handler).ServeHTTP(resp, req)
			res := resp.Result()

			assertEqual(t, 200, res.StatusCode)
			if tt.expectedGzip {
				assertEqual(t, "gzip", res.Header.Get("Content-Encoding"))
				zr, err := gzip.NewReader(resp.Body)
				assertNil(t, err)
				got, err := ioutil.ReadAll(zr)
				assertNil(t, err)
				assertEqual(t, testBody, got)

			} else {
				assertNotEqual(t, "gzip", res.Header.Get("Content-Encoding"))
				got, err := ioutil.ReadAll(resp.Body)
				assertNil(t, err)
				assertEqual(t, testBody, got)
			}
		})
	}
}

var contentTypeTest2 = []struct {
	name         string
	contentType  string
	expectedGzip bool
}{
	{
		name:         "Always gzip when content types are empty",
		contentType:  "",
		expectedGzip: true,
	},
	{
		name:         "MIME match",
		contentType:  "application/json",
		expectedGzip: true,
	},
	{
		name:         "MIME no match",
		contentType:  "text/xml",
		expectedGzip: true,
	},

	{
		name:         "MIME match case insensitive",
		contentType:  "Video/Something",
		expectedGzip: false,
	},
	{
		name:         "MIME match case insensitive",
		contentType:  "audio/Something",
		expectedGzip: false,
	},
	{
		name:         "MIME match ignore whitespace",
		contentType:  " video/mp4",
		expectedGzip: false,
	},
	{
		name:         "without prefix..",
		contentType:  "avideo/mp4",
		expectedGzip: true,
	},
	{
		name:         "application/zip",
		contentType:  "application/zip;lalala",
		expectedGzip: false,
	},
	{
		name:         "x-zip-compressed",
		contentType:  "application/x-zip-compressed",
		expectedGzip: false,
	},
	{
		name:         "application/x-gzip",
		contentType:  "application/x-gzip",
		expectedGzip: false,
	},
}

func TestDefaultContentTypes(t *testing.T) {
	for _, tt := range contentTypeTest2 {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", tt.contentType)
				w.Write(testBody)
			})

			wrapper, err := NewWrapper()
			assertNil(t, err)

			req, _ := http.NewRequest("GET", "/whatever", nil)
			req.Header.Set("Accept-Encoding", "gzip")
			resp := httptest.NewRecorder()
			wrapper(handler).ServeHTTP(resp, req)
			res := resp.Result()

			assertEqual(t, 200, res.StatusCode)
			if tt.expectedGzip {
				assertEqual(t, "gzip", res.Header.Get("Content-Encoding"))
			} else {
				assertNotEqual(t, "gzip", res.Header.Get("Content-Encoding"))
			}
		})
	}
}

var sniffTests = []struct {
	desc        string
	data        []byte
	contentType string
}{
	// Some nonsense.
	{"Empty", []byte{}, "text/plain; charset=utf-8"},
	{"Binary", []byte{1, 2, 3}, "application/octet-stream"},

	{"HTML document #1", []byte(`<HtMl><bOdY>blah blah blah</body></html>`), "text/html; charset=utf-8"},
	{"HTML document #2", []byte(`<HTML></HTML>`), "text/html; charset=utf-8"},
	{"HTML document #3 (leading whitespace)", []byte(`   <!DOCTYPE HTML>...`), "text/html; charset=utf-8"},
	{"HTML document #4 (leading CRLF)", []byte("\r\n<html>..."), "text/html; charset=utf-8"},

	{"Plain text", []byte(`This is not HTML. It has ☃ though.`), "text/plain; charset=utf-8"},

	{"XML", []byte("\n<?xml!"), "text/xml; charset=utf-8"},

	// Image types.
	{"Windows icon", []byte("\x00\x00\x01\x00"), "image/x-icon"},
	{"Windows cursor", []byte("\x00\x00\x02\x00"), "image/x-icon"},
	{"BMP image", []byte("BM..."), "image/bmp"},
	{"GIF 87a", []byte(`GIF87a`), "image/gif"},
	{"GIF 89a", []byte(`GIF89a...`), "image/gif"},
	{"WEBP image", []byte("RIFF\x00\x00\x00\x00WEBPVP"), "image/webp"},
	{"PNG image", []byte("\x89PNG\x0D\x0A\x1A\x0A"), "image/png"},
	{"JPEG image", []byte("\xFF\xD8\xFF"), "image/jpeg"},

	// Audio types.
	{"MIDI audio", []byte("MThd\x00\x00\x00\x06\x00\x01"), "audio/midi"},
	{"MP3 audio/MPEG audio", []byte("ID3\x03\x00\x00\x00\x00\x0f"), "audio/mpeg"},
	{"WAV audio #1", []byte("RIFFb\xb8\x00\x00WAVEfmt \x12\x00\x00\x00\x06"), "audio/wave"},
	{"WAV audio #2", []byte("RIFF,\x00\x00\x00WAVEfmt \x12\x00\x00\x00\x06"), "audio/wave"},
	{"AIFF audio #1", []byte("FORM\x00\x00\x00\x00AIFFCOMM\x00\x00\x00\x12\x00\x01\x00\x00\x57\x55\x00\x10\x40\x0d\xf3\x34"), "audio/aiff"},

	{"OGG audio", []byte("OggS\x00\x02\x00\x00\x00\x00\x00\x00\x00\x00\x7e\x46\x00\x00\x00\x00\x00\x00\x1f\xf6\xb4\xfc\x01\x1e\x01\x76\x6f\x72"), "application/ogg"},
	{"Must not match OGG", []byte("owow\x00"), "application/octet-stream"},
	{"Must not match OGG", []byte("oooS\x00"), "application/octet-stream"},
	{"Must not match OGG", []byte("oggS\x00"), "application/octet-stream"},

	// Video types.
	{"MP4 video", []byte("\x00\x00\x00\x18ftypmp42\x00\x00\x00\x00mp42isom<\x06t\xbfmdat"), "video/mp4"},
	{"AVI video #1", []byte("RIFF,O\n\x00AVI LISTÀ"), "video/avi"},
	{"AVI video #2", []byte("RIFF,\n\x00\x00AVI LISTÀ"), "video/avi"},

	// Font types.
	// {"MS.FontObject", []byte("\x00\x00")},
	{"TTF sample  I", []byte("\x00\x01\x00\x00\x00\x17\x01\x00\x00\x04\x01\x60\x4f"), "font/ttf"},
	{"TTF sample II", []byte("\x00\x01\x00\x00\x00\x0e\x00\x80\x00\x03\x00\x60\x46"), "font/ttf"},

	{"OTTO sample  I", []byte("\x4f\x54\x54\x4f\x00\x0e\x00\x80\x00\x03\x00\x60\x42\x41\x53\x45"), "font/otf"},

	{"woff sample  I", []byte("\x77\x4f\x46\x46\x00\x01\x00\x00\x00\x00\x30\x54\x00\x0d\x00\x00"), "font/woff"},
	{"woff2 sample", []byte("\x77\x4f\x46\x32\x00\x01\x00\x00\x00"), "font/woff2"},
	{"wasm sample", []byte("\x00\x61\x73\x6d\x01\x00"), "application/wasm"},

	// Archive types
	{"RAR v1.5-v4.0", []byte("Rar!\x1A\x07\x00"), "application/x-rar-compressed"},
	{"RAR v5+", []byte("Rar!\x1A\x07\x01\x00"), "application/x-rar-compressed"},
	{"Incorrect RAR v1.5-v4.0", []byte("Rar \x1A\x07\x00"), "application/octet-stream"},
	{"Incorrect RAR v5+", []byte("Rar \x1A\x07\x01\x00"), "application/octet-stream"},
}

func TestContentTypeDetect(t *testing.T) {
	for _, tt := range sniffTests {
		t.Run(tt.desc, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				for i := range tt.data {
					// Do one byte writes...
					w.Write([]byte{tt.data[i]})
				}
				w.Write(testBody)
			})

			wrapper, err := NewWrapper()
			assertNil(t, err)

			req, _ := http.NewRequest("GET", "/whatever", nil)
			req.Header.Set("Accept-Encoding", "gzip")
			resp := httptest.NewRecorder()
			wrapper(handler).ServeHTTP(resp, req)
			res := resp.Result()

			assertEqual(t, 200, res.StatusCode)
			assertEqual(t, tt.contentType, res.Header.Get("Content-Type"))
			shouldGZ := DefaultContentTypeFilter(tt.contentType)
			if shouldGZ {
				assertEqual(t, "gzip", res.Header.Get("Content-Encoding"))
			} else {
				assertNotEqual(t, "gzip", res.Header.Get("Content-Encoding"))
			}
		})
		t.Run(tt.desc+"empty", func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "")
				w.WriteHeader(http.StatusOK)
				for i := range tt.data {
					// Do one byte writes...
					w.Write([]byte{tt.data[i]})
				}
				w.Write(testBody)
			})

			wrapper, err := NewWrapper()
			assertNil(t, err)

			req, _ := http.NewRequest("GET", "/whatever", nil)
			req.Header.Set("Accept-Encoding", "gzip")
			resp := httptest.NewRecorder()
			wrapper(handler).ServeHTTP(resp, req)
			res := resp.Result()

			assertEqual(t, 200, res.StatusCode)
			// Is Content-Type still empty?
			assertEqual(t, "", res.Header.Get("Content-Type"))
			shouldGZ := DefaultContentTypeFilter(tt.contentType)
			if shouldGZ {
				assertEqual(t, "gzip", res.Header.Get("Content-Encoding"))
			} else {
				assertNotEqual(t, "gzip", res.Header.Get("Content-Encoding"))
			}
		})
		t.Run(tt.desc+"flush", func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "")
				w.WriteHeader(http.StatusOK)
				for i := range tt.data {
					// Do one byte writes...
					w.Write([]byte{tt.data[i]})
				}
				w.(http.Flusher).Flush()
				w.Write(testBody)
			})

			wrapper, err := NewWrapper()
			assertNil(t, err)

			req, _ := http.NewRequest("GET", "/whatever", nil)
			req.Header.Set("Accept-Encoding", "gzip")
			resp := httptest.NewRecorder()
			wrapper(handler).ServeHTTP(resp, req)
			res := resp.Result()

			assertEqual(t, 200, res.StatusCode)
			// Is Content-Type still empty?
			assertEqual(t, "", res.Header.Get("Content-Type"))
			shouldGZ := DefaultContentTypeFilter(tt.contentType)
			if shouldGZ {
				assertEqual(t, "gzip", res.Header.Get("Content-Encoding"))
			} else {
				assertNotEqual(t, "gzip", res.Header.Get("Content-Encoding"))
			}
		})
	}
}

// --------------------------------------------------------------------

func BenchmarkGzipHandler_S2k(b *testing.B)   { benchmark(b, false, 2048, gzip.DefaultCompression) }
func BenchmarkGzipHandler_S20k(b *testing.B)  { benchmark(b, false, 20480, gzip.DefaultCompression) }
func BenchmarkGzipHandler_S100k(b *testing.B) { benchmark(b, false, 102400, gzip.DefaultCompression) }
func BenchmarkGzipHandler_P2k(b *testing.B)   { benchmark(b, true, 2048, gzip.DefaultCompression) }
func BenchmarkGzipHandler_P20k(b *testing.B)  { benchmark(b, true, 20480, gzip.DefaultCompression) }
func BenchmarkGzipHandler_P100k(b *testing.B) { benchmark(b, true, 102400, gzip.DefaultCompression) }

func BenchmarkGzipBestSpeedHandler_S2k(b *testing.B)   { benchmark(b, false, 2048, gzip.BestSpeed) }
func BenchmarkGzipBestSpeedHandler_S20k(b *testing.B)  { benchmark(b, false, 20480, gzip.BestSpeed) }
func BenchmarkGzipBestSpeedHandler_S100k(b *testing.B) { benchmark(b, false, 102400, gzip.BestSpeed) }
func BenchmarkGzipBestSpeedHandler_P2k(b *testing.B)   { benchmark(b, true, 2048, gzip.BestSpeed) }
func BenchmarkGzipBestSpeedHandler_P20k(b *testing.B)  { benchmark(b, true, 20480, gzip.BestSpeed) }
func BenchmarkGzipBestSpeedHandler_P100k(b *testing.B) { benchmark(b, true, 102400, gzip.BestSpeed) }

// --------------------------------------------------------------------

func gzipStrLevel(s []byte, lvl int) []byte {
	var b bytes.Buffer
	w, _ := gzip.NewWriterLevel(&b, lvl)
	w.Write(s)
	w.Close()
	return b.Bytes()
}

func benchmark(b *testing.B, parallel bool, size, level int) {
	bin, err := ioutil.ReadFile("testdata/benchmark.json")
	if err != nil {
		b.Fatal(err)
	}

	req, _ := http.NewRequest("GET", "/whatever", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	handler := newTestHandlerLevel(bin[:size], level)

	b.ReportAllocs()
	b.SetBytes(int64(size))
	if parallel {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				runBenchmark(b, req, handler)
			}
		})
	} else {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			runBenchmark(b, req, handler)
		}
	}
}

func runBenchmark(b *testing.B, req *http.Request, handler http.Handler) {
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if code := res.Code; code != 200 {
		b.Fatalf("Expected 200 but got %d", code)
	} else if blen := res.Body.Len(); blen < 500 {
		b.Fatalf("Expected complete response body, but got %d bytes", blen)
	}
}

func newTestHandler(body []byte) http.Handler {
	return GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/gzipped":
			w.Header().Set("Content-Encoding", "gzip")
			w.Write(body)
		default:
			w.Write(body)
		}
	}))
}

func newTestHandlerLevel(body []byte, level int) http.Handler {
	wrapper, err := NewWrapper(CompressionLevel(level))
	if err != nil {
		panic(err)
	}
	return wrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/gzipped":
			w.Header().Set("Content-Encoding", "gzip")
			w.Write(body)
		default:
			w.Write(body)
		}
	}))
}

func TestGzipHandlerNilContentType(t *testing.T) {
	// This just exists to provide something for GzipHandler to wrap.
	handler := newTestHandler(testBody)

	// content-type header not set when provided nil

	req, _ := http.NewRequest("GET", "/whatever", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	res := httptest.NewRecorder()
	res.Header()["Content-Type"] = nil
	handler.ServeHTTP(res, req)

	assertEqual(t, "", res.Header().Get("Content-Type"))
}
