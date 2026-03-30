# Fix: SSE Connection Failures and Interruptions

**Status:** done

**Ticket:** MAI-73

**Branch:** hamoondev/mai-73-fix-sse-connection-failures-and-interuptions

**Created:** 2026-03-29

## Goal

Eliminate SSE connection drops on long-running tabs by adding server-side heartbeats and fixing the frontend reconnection logic.

## Context

Long-running browser tabs show CORS errors and Cloudflare 503 failures on the SSE endpoint. Root cause: no server-side heartbeat allows Cloudflare to kill idle connections, and a bug in the frontend reconnect counter creates an infinite retry loop. See `research.md` for full analysis.

## Out of Scope

- Migration to River/ergo for concurrency (tracked separately under MAI-70)
- WebSocket migration
- Token refresh integration with EventSource (would require switching to a fetch-based SSE client)

## Approach

### 1. Backend: Add heartbeat to SSE handler

Add a `time.Ticker` in the SSE event loop (`sse.go`) that sends an SSE comment (`: heartbeat\n\n`) every 30 seconds. SSE comments are lines starting with `:` — the browser's EventSource parser silently ignores them, but they keep the TCP connection alive through Cloudflare and any other intermediary proxies.

```go
// In SSEHandler, after flusher check:
heartbeat := time.NewTicker(30 * time.Second)
defer heartbeat.Stop()

for {
    select {
    case event := <-conn.eventQueue:
        // ... existing event handling ...
    case <-heartbeat.C:
        fmt.Fprintf(w, ": heartbeat\n\n")
        flusher.Flush()
    case <-r.Context().Done():
        return
    }
}
```

**Why 30 seconds:** Well under Cloudflare's ~100s idle timeout, with enough margin to survive occasional network jitter. Not so frequent that it adds meaningful overhead.

### 2. Frontend: Use native EventSource reconnection

The current code calls `es.close()` inside `onerror`, which permanently kills the EventSource and defeats the browser's built-in auto-reconnect. Instead, let the native retry handle transient failures and only intervene after repeated failures.

**How native EventSource reconnection works:** When the connection drops, the browser automatically retries after a delay (default ~3s, configurable via `retry:` field from server). It fires `onerror` on failure and `onopen` on successful reconnect — all without closing the EventSource.

**Changes:**
- Remove `closeConnection()` and manual `setTimeout` reconnect from `onerror`
- Track consecutive errors: increment on `onerror`, reset on `onopen`
- Only `close()` and show toast after `MAX_RECONNECT_ATTEMPTS` consecutive errors with no successful `onopen` in between
- Backend sends `retry: 3000\n` in initial SSE response to set the browser's retry interval

```typescript
const setupEventListeners = (es: EventSource) => {
    es.onopen = () => {
        reconnectAttempts.value = 0;  // Reset on successful connection
    };

    es.onerror = () => {
        // Native EventSource auto-reconnects — just count consecutive failures
        reconnectAttempts.value++;
        if (reconnectAttempts.value >= MAX_RECONNECT_ATTEMPTS) {
            closeConnection();  // Only close after repeated failures
            toast.error("Failed to connect to the server. Please refresh the page.");
        }
    };
    // ... rest unchanged
};
```

`closeConnection()` becomes a simple cleanup function (no counter reset):

```typescript
const closeConnection = () => {
    if (eventSource.value) {
        eventSource.value.close();
        eventSource.value = undefined;
    }
};
```

### 3. Backend: Handle heartbeat write errors

If `fmt.Fprintf` for the heartbeat fails (client disconnected), log and return to clean up the connection gracefully rather than looping on a dead writer.

### 4. Backend: Set retry interval in initial SSE response

Send `retry: 3000\n` before the first event so the browser uses a 3-second retry interval instead of the default.

## Tasks

- [ ] Add heartbeat ticker to `SSEHandler` in `backend/internal/events/sse.go`
- [ ] Handle write errors on heartbeat (return to trigger defer cleanup)
- [ ] Fix `closeConnection()` in `frontend/src/shared/hooks/events/useSse.ts` — stop resetting `reconnectAttempts`
- [ ] Add `onopen` handler to reset `reconnectAttempts` on successful connection
- [ ] Test: verify heartbeat appears in SSE stream (manual or integration test)
- [ ] Test: verify reconnection gives up after 3 attempts (manual or unit test)

## Acceptance Criteria

- [ ] SSE connections survive idle periods > 100 seconds without Cloudflare 503 it was actually 524 timeout but i think this doesnt change anything
- [ ] Server sends heartbeat comments every 30 seconds on idle connections
- [ ] Frontend reconnection counter works correctly — gives up after 3 failed attempts
- [ ] No regressions in existing SSE event delivery (goal created, XP updated)

## Open Questions

None.

## Decisions Log

- 2026-03-29: Chose SSE comment heartbeat (`: heartbeat\n\n`) over empty `data:` events — comments are invisible to EventSource listeners so no frontend handling needed
- 2026-03-29: Chose 30s interval — safe margin under Cloudflare's ~100s timeout without being noisy

## Session Log

- 2026-03-29: Investigated MAI-73. Identified two root causes: (1) no server-side heartbeat causes Cloudflare to kill idle SSE connections with 503, which browsers report as CORS errors since the 503 lacks CORS headers; (2) `closeConnection()` resets `reconnectAttempts` to 0 on every error, making the max-attempts guard unreachable and creating an infinite reconnect loop. Created research.md and design.md.
