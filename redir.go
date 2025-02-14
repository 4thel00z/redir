package redir

import (
	"errors"
	"net/http"
	"net/url"
	"time"
)

// Redirection holds information about a single redirection step.
type Redirection struct {
	URL        string        `json:"url"`
	StatusCode int           `json:"status_code"`
	Duration   time.Duration `json:"duration"`
}

// FollowRedirects performs a GET request starting at startURL and follows
// redirections up to maxRedirects. It returns a slice of Redirection detailing
// each step (including timing).
func FollowRedirects(startURL string, maxRedirects int) ([]Redirection, error) {
	var redirects []Redirection

	// Create a client that does NOT follow redirects automatically.
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	currentURL := startURL
	for i := 0; i < maxRedirects; i++ {
		req, err := http.NewRequest("GET", currentURL, nil)
		if err != nil {
			return redirects, err
		}
		startTime := time.Now()
		resp, err := client.Do(req)
		duration := time.Since(startTime)
		if err != nil {
			return redirects, err
		}
		// It's important to close the body to avoid resource leaks.
		resp.Body.Close()

		step := Redirection{
			URL:        currentURL,
			StatusCode: resp.StatusCode,
			Duration:   duration,
		}
		redirects = append(redirects, step)

		// Check if the status code indicates a redirection.
		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			location := resp.Header.Get("Location")
			if location == "" {
				return redirects, errors.New("redirection status code received but no Location header found")
			}
			// Resolve relative URLs.
			parsedCurrent, err := url.Parse(currentURL)
			if err != nil {
				return redirects, err
			}
			parsedLocation, err := url.Parse(location)
			if err != nil {
				return redirects, err
			}
			newURL := parsedCurrent.ResolveReference(parsedLocation)
			currentURL = newURL.String()
		} else {
			// Final destination reached.
			break
		}
	}

	return redirects, nil
}
