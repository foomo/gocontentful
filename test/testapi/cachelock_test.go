package testapi

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// These tests hold the cache update lock directly, which only code in this package can do.

func newLockTestClient(t *testing.T) *ContentfulClient {
	t.Helper()
	file, err := os.ReadFile("../test-space-export.json")
	if err != nil {
		t.Fatal(err)
	}
	cc, err := NewOfflineContentfulClient(file, nil, LogError, true, true)
	if err != nil {
		t.Fatal(err)
	}
	return cc
}

// returnsWithin fails the test if call is still running after d, and returns its error.
func returnsWithin(t *testing.T, d time.Duration, call func() error) error {
	t.Helper()
	done := make(chan error, 1)
	go func() { done <- call() }()
	select {
	case err := <-done:
		return err
	case <-time.After(d):
		t.Fatalf("call still blocked after %s", d)
		return nil
	}
}

func holdCacheUpdateLock(t *testing.T, cc *ContentfulClient) {
	t.Helper()
	if err := cc.cacheMutex.cacheUpdateGcLock.lock(context.Background()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(cc.cacheMutex.cacheUpdateGcLock.unlock)
}

func TestUpdateCacheForEntityStopsWaitingWhenContextEnds(t *testing.T) {
	cc := newLockTestClient(t)
	holdCacheUpdateLock(t, cc)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	err := returnsWithin(t, 2*time.Second, func() error {
		return cc.UpdateCacheForEntity(ctx, sysTypeEntry, ContentTypeBrand, "651CQ8rLoIYCeY6G0QG22q")
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected a deadline error, got %v", err)
	}
}

func TestSyncStopsWaitingWhenContextEnds(t *testing.T) {
	cc := newLockTestClient(t)
	holdCacheUpdateLock(t, cc)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	err := returnsWithin(t, 2*time.Second, func() error {
		_, _, err := cc.syncCache(ctx, spaceContentTypes)
		return err
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected a deadline error, got %v", err)
	}
}

func TestRebuildStopsWaitingAtCacheUpdateTimeout(t *testing.T) {
	cc := newLockTestClient(t)
	cc.cacheMutex.sharedDataGcLock.Lock()
	cc.offline = false // only online rebuilds have a deadline
	cc.cacheMutex.sharedDataGcLock.Unlock()
	cc.SetCacheUpdateTimeout(1)
	holdCacheUpdateLock(t, cc)
	err := returnsWithin(t, 3*time.Second, func() error {
		_, err := cc.runCacheJob(spaceContentTypes, true, false)
		return err
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected a deadline error, got %v", err)
	}
}

// rateLimitedOnce answers the first request with a rate-limit error asking to retry after four
// seconds, and every later one with an empty JSON object.
type rateLimitedOnce struct {
	requests atomic.Int32
}

func (rt *rateLimitedOnce) RoundTrip(r *http.Request) (*http.Response, error) {
	response := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{}`)),
		Request:    r,
	}
	if rt.requests.Add(1) == 1 {
		response.StatusCode = http.StatusTooManyRequests
		response.Header.Set("X-Contentful-Ratelimit-Reset", "4")
		response.Body = io.NopCloser(strings.NewReader(`{"sys":{"type":"Error","id":"RateLimitExceeded"},"message":"rate limit exceeded"}`))
	}
	return response, nil
}

// The SDK sleeps through Contentful's rate-limit reset without watching the context, so only the
// caller's own timeout keeps UpdateCache from waiting for such a rebuild.
func TestUpdateCacheWaitIsBoundedByTheTimeout(t *testing.T) {
	cc := newLockTestClient(t)
	cc.Client.SetHTTPClient(&http.Client{Transport: &rateLimitedOnce{}})
	cc.cacheMutex.sharedDataGcLock.Lock()
	cc.offline = false
	cc.cacheMutex.sharedDataGcLock.Unlock()
	cc.SetCacheUpdateTimeout(1)
	// 2 seconds of settling delay plus the 1 second timeout, with some slack.
	err := returnsWithin(t, 4*time.Second, func() error {
		_, _, err := cc.UpdateCache(context.Background(), nil, true)
		return err
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected a deadline error, got %v", err)
	}
}
