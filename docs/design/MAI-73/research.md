# Research: MAI-73 — SSE Connection Failures and Interruptions

**Ticket:** MAI-73

**Date:** 2026-03-29

## Scope

Investigate why long-running browser tabs show CORS errors and Cloudflare 503 failures on the SSE endpoint (`GET /api/events`). Determine root cause and identify all contributing factors.

## Codebase Findings

### Relevant Files

| File | Purpose |
|------|---------|
| `backend/internal/events/sse.go` | SSE handler, connection struct, event writing |
| `backend/internal/events/events.go` | EventManager pub/sub system |
| `backend/internal/middleware/middleware.go` | CORS config and middleware chains |
| `backend/internal/middleware/auth.go` | QueryTokenAuth middleware (SSE auth) |
| `backend/internal/routes/router.go` | Route registration (`GET /api/events`) |
| `frontend/src/shared/hooks/events/useSse.ts` | EventSource connection, reconnection logic |
| `frontend/src/pages/HomePage.vue` | SSE initialization on mount |

### Existing Patterns

**Backend SSE handler (`sse.go:53-105`):**
- Sets standard SSE headers: `text/event-stream`, `no-cache`, `keep-alive`, `X-Accel-Buffering: no`
- Creates an `SSEConn` with a buffered channel (size 10)
- Subscribes to user-scoped events via `EventManager`
- Sends an initial `sse_connected` event, then enters an infinite `select` loop waiting on the event channel or `r.Context().Done()`
- Flushes after every event write
- **No heartbeat or keepalive ping is sent between real events**

**Frontend SSE hook (`useSse.ts`):**
- Uses native `EventSource` API with token passed as query parameter
- `onerror` handler: closes connection, attempts reconnect up to `MAX_RECONNECT_ATTEMPTS` (3) with linear backoff (1s, 2s, 3s)
- Listens for `DEFAULT_GOAL_CREATED` and `XP_UPDATED` events
- Connected/disconnected in `onMounted`/`onUnmounted` of `HomePage.vue`

**Middleware chain for SSE (`middleware.go:49`):**
- `QueryTokenAuthChain` = Logging -> CORS -> QueryTokenAuth
- CORS is configured with `AllowedOrigins` from env, `AllowCredentials: true`, `AllowedHeaders: ["*"]`

### Constraints & Gotchas

**1. No server-side heartbeat (PRIMARY ROOT CAUSE)**

The SSE handler only writes data when a real application event occurs. Between events, the TCP connection sits completely idle. Cloudflare enforces idle connection timeouts (~100 seconds on most plans). When no data flows for that period, Cloudflare terminates the connection with a 503.

Since the 503 originates from Cloudflare's edge (not the Go backend), it does **not** include the app's CORS headers. The browser sees a cross-origin response without `Access-Control-Allow-Origin` and reports it as a CORS error. This is why both symptoms (CORS errors AND 503s) appear together — they are the same event seen from different layers.

**2. Broken reconnection counter (SECONDARY BUG)**

In `useSse.ts:40-51`, the `onerror` handler calls `closeConnection()` before checking the attempt counter. But `closeConnection()` (line 28-34) resets `reconnectAttempts` to 0 unconditionally:

```typescript
es.onerror = () => {
    closeConnection();  // resets reconnectAttempts to 0

    if (reconnectAttempts.value < MAX_RECONNECT_ATTEMPTS) {  // always true: 0 < 3
        reconnectAttempts.value++;  // goes to 1
        setTimeout(() => { connect(...); }, 1000 * reconnectAttempts.value);
        return;
    }
    toast.error("...");  // never reached
};
```

Result: the connection enters an infinite reconnect loop with a ~1 second delay. The "please refresh" toast never appears. This generates a stream of console errors on long-running tabs.

**3. Token staleness on reconnect**

When reconnecting, the frontend calls `getUser()?.access_token` to build the URL. If the token has expired during the idle period, the backend returns 401. The `zodFetch` automatic refresh logic doesn't apply here because `EventSource` is a separate browser API — it doesn't go through `zodFetch`. A stale token would cause auth failures that compound the reconnection spam.

## External Research

**SSE keepalive best practice:** The SSE spec supports comment lines (lines starting with `:`) that are ignored by the browser's EventSource parser but keep the TCP connection alive. The standard pattern is to send `: keepalive\n\n` every 15-30 seconds.

**Cloudflare idle timeouts:** Cloudflare terminates connections that are idle for ~100 seconds (varies by plan). This is well-documented and affects any long-lived HTTP connection (SSE, WebSocket upgrade handshakes, long-polling). The recommended mitigation is application-level keepalive.

**EventSource reconnection:** The browser's native EventSource has built-in reconnection (default ~3 seconds), but the current code overrides this by calling `.close()` in onerror and managing reconnection manually. The manual approach is fine but the counter bug defeats it.

## Open Questions

None — root causes are identified and the fix approach is straightforward.
