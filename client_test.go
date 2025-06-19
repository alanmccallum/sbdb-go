package sbdb

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

// roundTripperFunc allows customizing http.Client behavior in tests.
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

type recordCloser struct {
	io.ReadCloser
	closed bool
}

func (rc *recordCloser) Close() error {
	rc.closed = true
	if rc.ReadCloser != nil {
		return rc.ReadCloser.Close()
	}
	return nil
}

func TestClient_GetURL(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := &Client{Endpoint: "http://example.com/api"}
		f := Filter{Fields: NewFieldSet(SpkID), Limit: 5}
		u, err := c.GetURL(f)
		if err != nil {
			t.Fatalf("GetURL() error = %v", err)
		}
		want, err := url.Parse("http://example.com/api?fields=spkid&limit=5")
		if err != nil {
			t.Fatalf("failed to parse want url: %v", err)
		}
		if diff := cmp.Diff(want.Host, u.Host); diff != "" {
			t.Errorf("host mismatch (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff(want.Path, u.Path); diff != "" {
			t.Errorf("path mismatch (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff(want.Query(), u.Query()); diff != "" {
			t.Errorf("query mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("filter error", func(t *testing.T) {
		c := &Client{Endpoint: "http://example.com/api"}
		if _, err := c.GetURL(Filter{}); err == nil {
			t.Fatal("expected error for invalid filter")
		}
	})

	t.Run("bad endpoint", func(t *testing.T) {
		c := &Client{Endpoint: "://bad url"}
		_, err := c.GetURL(Filter{Fields: NewFieldSet(SpkID)})
		if err == nil {
			t.Fatal("expected error for bad endpoint")
		}
	})
}

func TestClient_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var got *http.Request
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got = r
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		c := &Client{Endpoint: srv.URL}
		resp, err := c.Get(context.Background(), Filter{Fields: NewFieldSet(SpkID), Limit: 1})
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		resp.Body.Close()
		if got == nil {
			t.Fatal("expected request to server")
		}
		if diff := cmp.Diff(url.Values{"fields": []string{"spkid"}, "limit": []string{"1"}}, got.URL.Query()); diff != "" {
			t.Errorf("query mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("filter error", func(t *testing.T) {
		called := false
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true }))
		defer srv.Close()

		c := &Client{Endpoint: srv.URL}
		if _, err := c.Get(context.Background(), Filter{}); err == nil {
			t.Fatal("expected error for invalid filter")
		}
		if called {
			t.Error("server should not have been called")
		}
	})

	t.Run("bad endpoint", func(t *testing.T) {
		c := &Client{Endpoint: "://bad url"}
		if _, err := c.Get(context.Background(), Filter{Fields: NewFieldSet(SpkID)}); err == nil {
			t.Fatal("expected error for bad endpoint")
		}
	})

	t.Run("context canceled", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(50 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		c := &Client{Endpoint: srv.URL}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		if _, err := c.Get(ctx, Filter{Fields: NewFieldSet(SpkID), Limit: 1}); err == nil {
			t.Fatal("expected error when context is canceled")
		}
	})

	t.Run("http 500", func(t *testing.T) {
		var rc *recordCloser
		rt := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			rc = &recordCloser{ReadCloser: io.NopCloser(strings.NewReader("boom"))}
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       rc,
				Header:     make(http.Header),
				Request:    r,
			}, nil
		})

		c := &Client{Client: http.Client{Transport: rt}, Endpoint: "http://example.com"}
		_, err := c.Get(context.Background(), Filter{Fields: NewFieldSet(SpkID), Limit: 1})
		if err == nil {
			t.Fatal("expected error for http 500")
		}
		if rc != nil && !rc.closed {
			t.Error("expected response body to be closed")
		}
	})
}
