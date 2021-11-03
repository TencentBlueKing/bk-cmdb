// Copyright (c) 2021 Klaus Post. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gzhttp

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/klauspost/compress/gzip"
)

func TestTransport(t *testing.T) {
	bin, err := ioutil.ReadFile("testdata/benchmark.json")
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(newTestHandler(bin))

	c := http.Client{Transport: Transport(http.DefaultTransport)}
	resp, err := c.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, bin) {
		t.Errorf("data mismatch")
	}
}

func TestTransportInvalid(t *testing.T) {
	bin, err := ioutil.ReadFile("testdata/benchmark.json")
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(newTestHandler(bin))

	c := http.Client{Transport: Transport(http.DefaultTransport)}
	// Serves json as gzippped...
	resp, err := c.Get(server.URL + "/gzipped")
	if err != nil {
		t.Fatal(err)
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDefaultTransport(t *testing.T) {
	bin, err := ioutil.ReadFile("testdata/benchmark.json")
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(newTestHandler(bin))

	// Not wrapped...
	c := http.Client{Transport: http.DefaultTransport}
	resp, err := c.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, bin) {
		t.Errorf("data mismatch")
	}
}

func BenchmarkTransport(b *testing.B) {
	bin, err := ioutil.ReadFile("testdata/benchmark.json")
	if err != nil {
		b.Fatal(err)
	}
	sz := len(bin)
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	zw.Write(bin)
	zw.Close()
	bin = buf.Bytes()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body.Close()
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(bin)
	}))
	b.Run("gzhttp", func(b *testing.B) {
		c := http.Client{Transport: Transport(http.DefaultTransport)}

		b.SetBytes(int64(sz))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			resp, err := c.Get(server.URL + "/gzipped")
			if err != nil {
				b.Fatal(err)
			}
			_, err = io.Copy(ioutil.Discard, resp.Body)
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()
		}
	})
	b.Run("stdlib", func(b *testing.B) {
		c := http.Client{Transport: http.DefaultTransport}
		b.SetBytes(int64(sz))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			resp, err := c.Get(server.URL + "/gzipped")
			if err != nil {
				b.Fatal(err)
			}
			_, err = io.Copy(ioutil.Discard, resp.Body)
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()
		}
	})
	b.Run("gzhttp-par", func(b *testing.B) {
		c := http.Client{
			Transport: Transport(&http.Transport{
				MaxConnsPerHost:     runtime.GOMAXPROCS(0),
				MaxIdleConnsPerHost: runtime.GOMAXPROCS(0),
			}),
		}

		b.SetBytes(int64(sz))
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				resp, err := c.Get(server.URL + "/gzipped")
				if err != nil {
					b.Fatal(err)
				}
				_, err = io.Copy(ioutil.Discard, resp.Body)
				if err != nil {
					b.Fatal(err)
				}
				resp.Body.Close()
			}
		})
	})
	b.Run("stdlib-par", func(b *testing.B) {
		c := http.Client{Transport: &http.Transport{
			MaxConnsPerHost:     runtime.GOMAXPROCS(0),
			MaxIdleConnsPerHost: runtime.GOMAXPROCS(0),
		}}
		b.SetBytes(int64(sz))
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				resp, err := c.Get(server.URL + "/gzipped")
				if err != nil {
					b.Fatal(err)
				}
				_, err = io.Copy(ioutil.Discard, resp.Body)
				if err != nil {
					b.Fatal(err)
				}
				resp.Body.Close()
			}
		})
	})
}
