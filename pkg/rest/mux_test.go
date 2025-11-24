package rest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, urlPath string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+urlPath, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestDefaultRouter(t *testing.T) {
	r := NewRouter()
	// r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("index"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/", nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, "index")

	resp, body = testRequest(t, ts, "GET", "/xxx", nil)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound)
	assert.Equal(t, "404 page not found\n", body)

	resp, body = testRequest(t, ts, "POST", "/", nil)
	assert.Equal(t, resp.StatusCode, http.StatusMethodNotAllowed)
	assert.Equal(t, "GET, HEAD", resp.Header.Get("Allow"))
	assert.Equal(t, "Method Not Allowed\n", body)
}

func TestConflictsRouter(t *testing.T) {
	defer func() {
		r := recover()
		require.NotNil(t, r)
		err, ok := r.(error)
		require.True(t, ok)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "conflicts with pattern")
	}()

	r := NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("index"))
	})

	r.Mount("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("index"))
	}))

	r.Get("/{bizID}/user", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("index"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()
}

func TestInvalidURLRouter(t *testing.T) {
	r := NewRouter()

	checkPanic := func(errMsg string) {
		r := recover()
		require.NotNil(t, r)
		err, ok := r.(error)
		require.True(t, ok)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errMsg)
	}

	t.Run("empty url", func(t *testing.T) {
		defer checkPanic("host/path missing /")

		r.Get("", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("index"))
		})
	})

	t.Run("not start with /", func(t *testing.T) {
		defer checkPanic("host/path missing /")

		r.Get("abc", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("index"))
		})
	})

	t.Run("mount not start with /", func(t *testing.T) {
		defer checkPanic("pattern must begin with /")

		r.Mount("api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("index"))
		}))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()
}

func TestParamsRouter(t *testing.T) {
	r := NewRouter()
	// r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("index"))
	})

	r.Get("/{bizID}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.PathValue("bizID")))
	})

	r.Get("/{bizID}/module/{moduleID}", func(w http.ResponseWriter, r *http.Request) {
		body := r.PathValue("bizID") + ":" + r.PathValue("moduleID")
		w.Write([]byte(body))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/", nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, body, "index")

	resp, body = testRequest(t, ts, "GET", "/xxx", nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, "xxx", body)

	resp, body = testRequest(t, ts, "GET", "/xxx/module/abc", nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, "xxx:abc", body)

	resp, body = testRequest(t, ts, "GET", "/xxx/module/abc/xxx", nil)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound)
	assert.Equal(t, "404 page not found\n", body)
}

func TestMountRouter(t *testing.T) {
	r := NewRouter()
	// r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("index"))
	})

	r.Mount("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := "mount index:" + r.URL.Path
		w.Write([]byte(body))
	}))

	subR := NewRouter()
	subR.Get("/user", func(w http.ResponseWriter, r *http.Request) {
		body := "api index:" + r.URL.Path
		w.Write([]byte(body))
	})
	r.Mount("/api", subR)

	r.Get("/biz/{bizID}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.PathValue("bizID")))
	})

	r.Get("/biz/{bizID}/module/{moduleID}", func(w http.ResponseWriter, r *http.Request) {
		body := r.PathValue("bizID") + ":" + r.PathValue("moduleID")
		w.Write([]byte(body))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, body, "index")

	resp, body = testRequest(t, ts, "GET", "/xxx", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "mount index:/xxx", body)

	resp, body = testRequest(t, ts, "GET", "/biz/xxx/module/abc", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "xxx:abc", body)

	resp, body = testRequest(t, ts, "GET", "/biz/xxx/module/abc/xxx", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "mount index:/biz/xxx/module/abc/xxx", body)

	resp, body = testRequest(t, ts, "GET", "/api/user", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "api index:/user", body)

}

func TestMuxPattern(t *testing.T) {
	r := NewRouter()
	// r := chi.NewRouter()

	r.Route("/org/{org}", func(r Router) {
		r.Get("/user/pattern", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(r.Pattern))
		})
	})

	subR := NewRouter()
	subR.Get("/user/{user}/pattern", func(w http.ResponseWriter, r *http.Request) {
		// pattern := chi.RouteContext(r.Context()).RoutePattern()
		w.Write([]byte(r.Pattern))
	})

	subR.Get("/user/{user}/pattern/ctx", func(w http.ResponseWriter, r *http.Request) {
		// pattern := chi.RouteContext(r.Context()).RoutePattern()
		w.Write([]byte(RoutePattern(r)))
	})

	subR2 := NewRouter()
	subR2.Mount("/pattern/{pattern}", subR)

	r.Mount("/", subR)
	r.Mount("/org", subR)
	r.Mount("/mount/{org}", subR)
	r.Mount("/mount1/{org}", subR2)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/org/1/user/pattern", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "GET /org/{org}/user/pattern", body)

	resp, body = testRequest(t, ts, "GET", "/mount/1/user/2/pattern", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "GET /user/{user}/pattern", body)

	resp, body = testRequest(t, ts, "GET", "/mount/1/user/2/pattern/ctx", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "/mount/{org}/user/{user}/pattern/ctx", body)

	resp, body = testRequest(t, ts, "GET", "/mount1/1/pattern/11/user/2/pattern/ctx", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "/mount1/{org}/pattern/{pattern}/user/{user}/pattern/ctx", body)

	resp, body = testRequest(t, ts, "GET", "/user/2/pattern", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "GET /user/{user}/pattern", body)

	resp, body = testRequest(t, ts, "GET", "/org/user/2/pattern", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "GET /user/{user}/pattern", body)
}

func TestRoutePattern(t *testing.T) {

	r := NewRouter()
	// r := chi.NewRouter()

	r.Route("/org/{org}", func(r Router) {
		r.Get("/user/pattern", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(r.Pattern))
		})
	})

	subR := NewRouter()
	subR.Get("/user/{user}/pattern", func(w http.ResponseWriter, r *http.Request) {
		// pattern := chi.RouteContext(r.Context()).RoutePattern()
		w.Write([]byte(r.Pattern))
	})

	subR.Get("/user/{user}/pattern/ctx", func(w http.ResponseWriter, r *http.Request) {
		// pattern := chi.RouteContext(r.Context()).RoutePattern()
		w.Write([]byte(RoutePattern(r)))
	})

	r.Mount("/mount/{org}", subR)
	r.Mount("/", subR)

	r.Mount("/org", subR)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/org/1/user/pattern", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "GET /org/{org}/user/pattern", body)

	resp, body = testRequest(t, ts, "GET", "/mount/1/user/2/pattern", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "GET /user/{user}/pattern", body)

	resp, body = testRequest(t, ts, "GET", "/mount/1/user/2/pattern/ctx", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "/mount/{org}/user/{user}/pattern/ctx", body)

	resp, body = testRequest(t, ts, "GET", "/user/2/pattern", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "GET /user/{user}/pattern", body)

	resp, body = testRequest(t, ts, "GET", "/org/user/2/pattern", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "GET /user/{user}/pattern", body)

	resp, body = testRequest(t, ts, "GET", "/org/user/2/pattern", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "GET /user/{user}/pattern", body)
}
func TestNotFoundRouter(t *testing.T) {
	r := NewRouter()
	// r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("index"))
	})

	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("404"))
	}))

	r.MethodNotAllowed(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not allowed"))
	}))

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/", nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, "index", body)

	resp, body = testRequest(t, ts, "GET", "/xxx", nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, "404", body)

	resp, body = testRequest(t, ts, "POST", "/", nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, "GET, HEAD", resp.Header.Get("Allow"))
	assert.Equal(t, "not allowed", body)
}

func TestSubRouter(t *testing.T) {
	r := NewRouter()
	// r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("index"))
	})

	r.Route("/api", func(r Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("api index"))
		})

		r.Get("/user", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("api user"))
		})
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/", nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, "index", body)

	resp, body = testRequest(t, ts, "GET", "/xxx", nil)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound)
	assert.Equal(t, "404 page not found\n", body)

	resp, body = testRequest(t, ts, "GET", "/api", nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, "api index", body)

	resp, body = testRequest(t, ts, "GET", "/api/user", nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, "api user", body)
}

func TestGroupRouter(t *testing.T) {
	r := NewRouter()
	// r := chi.NewRouter()

	r.Group(func(r Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("group index"))
		})

		r.Get("/user", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("group user"))
		})
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/", nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, "group index", body)

	resp, body = testRequest(t, ts, "GET", "/xxx", nil)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound)
	assert.Equal(t, "404 page not found\n", body)

	resp, body = testRequest(t, ts, "GET", "/user", nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, "group user", body)
}

func TestGroupConflictsRouter(t *testing.T) {
	defer func() {
		r := recover()
		require.NotNil(t, r)
		err, ok := r.(error)
		require.True(t, ok)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "conflicts with pattern")
	}()

	r := NewRouter()
	// r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("index"))
	})

	r.Group(func(r Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("group index"))
		})

		r.Get("/user", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("group user"))
		})
	})

	ts := httptest.NewServer(r)
	defer ts.Close()
}

func TestUseRouter(t *testing.T) {
	r := NewRouter()
	// r := chi.NewRouter()

	h := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("index"))
	}

	m := func(name string) func(next http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(name + ":"))
				next.ServeHTTP(w, r)
			})
		}
	}

	r.Use(m("m1"))
	r.Use(m("m2"))

	r.With(m("m3")).Get("/user", h)
	r.Get("/", h)

	r.Group(func(r Router) {
		r.Use(m("m4"))
		r.Get("/user2", h)
	})

	r.Route("/api", func(r Router) {
		r.Use(m("m5"))
		r.Get("/user3", h)
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "m1:m2:index", body)

	resp, body = testRequest(t, ts, "GET", "/user", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "m1:m2:m3:index", body)

	resp, body = testRequest(t, ts, "GET", "/api/user3", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "m1:m2:m5:index", body)
}
