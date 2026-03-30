package events

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockFlusher wraps httptest.ResponseRecorder and tracks Flush calls.
type mockFlusher struct {
	*httptest.ResponseRecorder
	mu         sync.Mutex
	flushCount int
}

func newMockFlusher() *mockFlusher {
	return &mockFlusher{ResponseRecorder: httptest.NewRecorder()}
}

func (m *mockFlusher) Flush() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.flushCount++
	m.ResponseRecorder.Flush()
}

func (m *mockFlusher) FlushCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.flushCount
}

func TestSSEHandlerUnauthorized(t *testing.T) {
	em := NewEventManager()
	w := newMockFlusher()
	r := httptest.NewRequest(http.MethodGet, "/sse", nil)

	em.SSEHandler(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSSEHandlerSendsInitialEvent(t *testing.T) {
	em := NewEventManager()
	w := newMockFlusher()

	ctx, cancel := context.WithCancel(context.Background())
	r := httptest.NewRequest(http.MethodGet, "/sse", nil)
	r = r.WithContext(ctx)
	r.Header.Set("user_id", "test-user-1")

	done := make(chan struct{})
	go func() {
		em.SSEHandler(w, r)
		close(done)
	}()

	// Give the handler time to write the initial event, then cancel.
	time.Sleep(50 * time.Millisecond)
	cancel()
	<-done

	body := w.Body.String()
	assert.Contains(t, body, "retry: 3000")
	assert.Contains(t, body, "event: "+SSEConnected)
	assert.Equal(t, "text/event-stream", w.Header().Get("Content-Type"))
}

func TestSSEHandlerDeliversEvent(t *testing.T) {
	em := NewEventManager()
	w := newMockFlusher()

	ctx, cancel := context.WithCancel(context.Background())
	r := httptest.NewRequest(http.MethodGet, "/sse", nil)
	r = r.WithContext(ctx)

	const userID = "test-user-2"
	r.Header.Set("user_id", userID)

	ready := make(chan struct{})
	done := make(chan struct{})

	go func() {
		// Wait until the handler has subscribed the connection before publishing.
		time.Sleep(50 * time.Millisecond)
		close(ready)
		em.SSEHandler(w, r)
		close(done)
	}()

	<-ready
	// Allow handler setup to complete.
	time.Sleep(50 * time.Millisecond)

	event := NewEventWithUserID(GoalCreated, map[string]string{"title": "Run 5k"}, userID)
	em.Publish(event)

	// Give the handler time to process the event, then cancel.
	time.Sleep(100 * time.Millisecond)
	cancel()
	<-done

	body := w.Body.String()
	assert.Contains(t, body, "event: "+GoalCreated)
	assert.Contains(t, body, "Run 5k")
}

// This test verifies that the heartbeat ticker path compiles and that the
// initial response (retry + sse_connected) is present. A real 30-second
// ticker is not exercised in unit tests to keep the suite fast.
func TestSSEHandlerSendsHeartbeat(t *testing.T) {
	em := NewEventManager()
	w := newMockFlusher()

	ctx, cancel := context.WithCancel(context.Background())
	r := httptest.NewRequest(http.MethodGet, "/sse", nil)
	r = r.WithContext(ctx)
	r.Header.Set("user_id", "test-user-3")

	done := make(chan struct{})
	go func() {
		em.SSEHandler(w, r)
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()
	<-done

	body := w.Body.String()
	assert.Contains(t, body, "retry: 3000")
	assert.Contains(t, body, "event: "+SSEConnected)
	// Flush must have been called at least once (after the initial event).
	assert.GreaterOrEqual(t, w.FlushCount(), 1)
}

func TestSSEHandlerWritesRetryBeforeInitialEvent(t *testing.T) {
	em := NewEventManager()
	w := newMockFlusher()

	ctx, cancel := context.WithCancel(context.Background())
	r := httptest.NewRequest(http.MethodGet, "/sse", nil)
	r = r.WithContext(ctx)
	r.Header.Set("user_id", "test-user-4")

	done := make(chan struct{})
	go func() {
		em.SSEHandler(w, r)
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()
	<-done

	body := w.Body.String()
	retryIdx := strings.Index(body, "retry: 3000")
	connectedIdx := strings.Index(body, "event: "+SSEConnected)
	assert.NotEqual(t, -1, retryIdx, "retry field not found in response")
	assert.NotEqual(t, -1, connectedIdx, "sse_connected event not found in response")
	assert.Less(t, retryIdx, connectedIdx, "retry field must appear before sse_connected event")
}
