package redir

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Helper to create a test server that performs redirections.
func createRedirectServer(redirects []string, finalStatus int) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Pop first URL from slice
		if len(redirects) > 0 {
			target := redirects[0]
			redirects = redirects[1:]
			w.Header().Set("Location", target)
			w.WriteHeader(http.StatusFound)
			return
		}
		w.WriteHeader(finalStatus)
	})
	return httptest.NewServer(handler)
}

func TestFollowRedirects_NoRedirect(t *testing.T) {
	// Create a server that always returns OK.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	steps, err := FollowRedirects(ts.URL, 5)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(steps) != 1 {
		t.Fatalf("Expected 1 step, got %d", len(steps))
	}

	if steps[0].StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, steps[0].StatusCode)
	}
}

func TestFollowRedirects_MultipleRedirects(t *testing.T) {
	// Create a chain of test servers.
	// We'll simulate a chain like: ts1 -> ts2 -> ts3 -> final (200 OK).
	ts3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts3.Close()

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", ts3.URL)
		w.WriteHeader(http.StatusFound)
	}))
	defer ts2.Close()

	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", ts2.URL)
		w.WriteHeader(http.StatusFound)
	}))
	defer ts1.Close()

	steps, err := FollowRedirects(ts1.URL, 10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(steps) != 3 {
		t.Fatalf("Expected 3 steps, got %d", len(steps))
	}

	// Check final status.
	final := steps[len(steps)-1]
	if final.StatusCode != http.StatusOK {
		t.Errorf("Expected final status %d, got %d", http.StatusOK, final.StatusCode)
	}
}

func TestFollowRedirects_NoLocationHeader(t *testing.T) {
	// Server that returns a redirect status without Location header.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusFound)
	}))
	defer ts.Close()

	_, err := FollowRedirects(ts.URL, 5)
	if err == nil {
		t.Fatal("Expected error due to missing Location header, got nil")
	}
	if !errors.Is(err, errors.New("redirection status code received but no Location header found")) {
		// We cannot compare errors directly with errors.New so we can check error string.
		if err.Error() != "redirection status code received but no Location header found" {
			t.Errorf("Unexpected error message: %v", err)
		}
	}
}

func TestFollowRedirects_TimeoutSimulation(t *testing.T) {
	// Create a server that delays response to simulate a slow redirection.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	steps, err := FollowRedirects(ts.URL, 5)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(steps) != 1 {
		t.Fatalf("Expected 1 step, got %d", len(steps))
	}

	// Check that duration is at least 100ms.
	if steps[0].Duration < 100*time.Millisecond {
		t.Errorf("Expected duration >= 100ms, got %v", steps[0].Duration)
	}
}
