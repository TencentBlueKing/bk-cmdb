// Copyright (c) 2021 Klaus Post. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gzhttp

import (
	"io"
	"net/http"
	"sync"

	"github.com/klauspost/compress/gzip"
)

// Transport will wrap a transport with a custom gzip handler
// that will request gzip and automatically decompress it.
// Using this is significantly faster than using the default transport.
func Transport(parent http.RoundTripper) http.RoundTripper {
	return gzRoundtripper{parent: parent}
}

type gzRoundtripper struct {
	parent http.RoundTripper
}

func (g gzRoundtripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var requestedGzip bool
	if req.Header.Get("Accept-Encoding") == "" &&
		req.Header.Get("Range") == "" &&
		req.Method != "HEAD" {
		// Request gzip only, not deflate. Deflate is ambiguous and
		// not as universally supported anyway.
		// See: https://zlib.net/zlib_faq.html#faq39
		//
		// Note that we don't request this for HEAD requests,
		// due to a bug in nginx:
		//   https://trac.nginx.org/nginx/ticket/358
		//   https://golang.org/issue/5522
		//
		// We don't request gzip if the request is for a range, since
		// auto-decoding a portion of a gzipped document will just fail
		// anyway. See https://golang.org/issue/8923
		requestedGzip = true
		req.Header.Set("Accept-Encoding", "gzip")
	}
	resp, err := g.parent.RoundTrip(req)
	if err != nil || !requestedGzip {
		return resp, err
	}
	if asciiEqualFold(resp.Header.Get("Content-Encoding"), "gzip") {
		resp.Body = &gzipReader{body: resp.Body}
		resp.Header.Del("Content-Encoding")
		resp.Header.Del("Content-Length")
		resp.ContentLength = -1
		resp.Uncompressed = true
	}
	return resp, nil
}

var gzReaderPool sync.Pool

// gzipReader wraps a response body so it can lazily
// call gzip.NewReader on the first call to Read
type gzipReader struct {
	body io.ReadCloser // underlying HTTP/1 response body framing
	zr   *gzip.Reader  // lazily-initialized gzip reader
	zerr error         // any error from gzip.NewReader; sticky
}

func (gz *gzipReader) Read(p []byte) (n int, err error) {
	if gz.zr == nil {
		if gz.zerr == nil {
			zr, ok := gzReaderPool.Get().(*gzip.Reader)
			if ok {
				gz.zr, gz.zerr = zr, zr.Reset(gz.body)
			} else {
				gz.zr, gz.zerr = gzip.NewReader(gz.body)
			}
		}
		if gz.zerr != nil {
			return 0, gz.zerr
		}
	}

	return gz.zr.Read(p)
}

func (gz *gzipReader) Close() error {
	if gz.zr != nil {
		gzReaderPool.Put(gz.zr)
		gz.zr = nil
	}
	return gz.body.Close()
}

// asciiEqualFold is strings.EqualFold, ASCII only. It reports whether s and t
// are equal, ASCII-case-insensitively.
func asciiEqualFold(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if lower(s[i]) != lower(t[i]) {
			return false
		}
	}
	return true
}

// lower returns the ASCII lowercase version of b.
func lower(b byte) byte {
	if 'A' <= b && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}
