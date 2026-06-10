# SSE Reliability Audit

**Status:** research / proposal — no code changed yet.
**Verified against the actually-wired files** (not assumptions):
`frontend/src/pages/HomePage.vue` → `frontend/src/shared/hooks/events/useSSE.ts` → `GET /api/events` (`backend/internal/routes/router.go`) → `backend/internal/events/sse.go` + `backend/internal/events/events.go`. Server setup: `backend/cmd/app/app.go`. Auth: `backend/internal/middleware/auth.go`.
**Symptoms (prod, behind Cloudflare):** the XP/progress number sometimes bumps on goal-complete, sometimes doesn't; "bad access token" errors on the SSE endpoint; the number sits stagnant a while, then jumps by more than one.

---

## TL;DR — root causes, ranked

1. **The backend fan-out blocks the entire hub on one slow/dead client.** `SSEConn.HandleEvent` does a *blocking* send on a 10-slot channel with no `select/default` (`events/sse.go:30-32`), and `processEvents` runs the broadcast **while holding the global `EventManager` mutex** (`events.go:157-175`). There is a **single** `processEvents` goroutine (`events.go:89`) draining one shared queue under one global lock, so a stall doesn't just hurt the slow connection — it freezes event delivery for **every** connection at once, including your other tabs and a brand-new login. One stalled connection (a backgrounded tab, or a Cloudflare-dropped socket) fills its 10-slot buffer, blocks the broadcast, and the held mutex jams everything. This single bug explains "stagnant for a while," "then bumps by more than one" (the backlog flushes in a burst when it unblocks), **and** why a fresh login still shows nothing while an older tab is wedged.
2. **No write deadline anywhere, so a half-dead socket hangs forever.** `app.go:105-108` builds `http.Server{Addr, Handler}` with no timeouts. That's correct for *keeping* a stream open, but it means a write to a Cloudflare-dropped-but-not-yet-closed connection blocks ~indefinitely, which is exactly what wedges the hub in #1.
3. **No token refresh on the SSE path.** The URL is built once with the current `access_token` (`HomePage.vue:20`); `es.onerror` can't see HTTP status and just counts (`useSSE.ts:43-51`). An expired token on (re)connect → backend 401 "invalid access token" (`auth.go:73-84`) → browser retries the *same* expired URL → after 3 errors it gives up with a "refresh the page" toast. This is the "bad access token" error and a permanent stall until manual refresh.
4. **The connection is scoped to the Home page, not the session.** `useSSE().connect()` runs in `HomePage.vue` `onMounted` and `closeConnection()` in `onUnmounted` (`HomePage.vue:19-25`). Leave the dashboard → no live updates; the number can't move or reconcile while you're elsewhere.
5. **State updates from events only, with no refetch fallback.** `XP_UPDATED` carries an absolute `{xp, level_id}` and `setUser` *merges* it (`useSSE.ts:57-64`) — so duplicates don't double-count (good), but a *missed* event leaves the number stale until the next received event or a re-login. Nothing reconciles against the server.

Items 1–2 are the biggest reliability bugs and are backend-only. 3–5 explain the "bad token" errors and the "sometimes it just doesn't update."

---

## What's already correct (don't touch)

- **Heartbeat exists:** a 30s ticker writes `: heartbeat\n\n` and flushes (`sse.go:91-118`). 30s < Cloudflare's ~100s idle timeout, so idle-timeout drops are *not* the main issue — as long as writes actually flush (see #1/#2).
- **Headers are mostly right:** `text/event-stream`, `no-cache`, `keep-alive`, `X-Accel-Buffering: no`, and `retry: 3000` (`sse.go:60-77`).
- **Flusher is used** after every write/heartbeat (`sse.go:89,107,118`).
- **Absolute-value payloads** mean duplicate delivery won't double-count the number.

---

## How it works today (real flow)

1. `HomePage.vue` mounts → `connect(${API_BASE}/events?token=${access_token})`.
2. `useSSE.ts` opens a native `EventSource`; listeners: `DEFAULT_GOAL_CREATED` → invalidate category queries; `XP_UPDATED` → `setUser({...user, xp, level_id})`.
3. Backend `QueryTokenAuth` validates `?token=` once at connect (`auth.go:57-90`), sets `user_id`, hands off to `SSEHandler`.
4. `SSEHandler` registers the conn in `userSubs[userID]` and loops on its 10-slot channel + 30s heartbeat (`sse.go:54-123`).
5. Completing a goal: `PATCH /api/goals/{id}` → `GoalUpdated` event → user service sees status→`complete`, does `newXp = user.Xp + 1`, publishes `XPUpdated` → `processEvents` → `broadcastByUser` → `conn.HandleEvent` → written to the stream → client `setUser`.

---

## Symptom → cause

| You observe | Most likely cause |
|---|---|
| Number sometimes doesn't bump | Missed `xp_updated` (connection dropped/reconnecting, buffer dropped under #1, off the Home page, or token gave up) with no refetch to self-heal (#3,#4,#5) |
| "Bad access token" on SSE | Expired access token at (re)connect, no SSE refresh, then permanent give-up after 3 retries (#3) |
| Stagnant, then jumps by >1 | Hub frozen by one stalled client (#1/#2); backlog flushes at once; absolute payload jumps straight to the latest value |

---

## Why this hits you specifically (single user, tabs left open)

You're the only user in prod, so "freezes for every user" really means **"freezes across all of *your* tabs and sessions"** — because all of them are served by the one `processEvents` goroutine and one global lock (#1).

The trigger is leaving tabs open. Browsers throttle and eventually suspend background tabs; a suspended tab stops reading its `EventSource` socket, so the server-side 10-slot buffer for that connection fills and the next broadcast blocks on it. With no write deadline (#2), that blocked write doesn't fail for a long time (only OS TCP timeouts), so the freeze persists.

**That is why a *fresh login* now breaks too:** the new tab connects fine and registers a new subscriber, but the single shared event loop is still stuck delivering to the wedged old tab, so the new tab receives nothing. Closing all other tabs (or restarting the backend) "fixes" it temporarily — a tell-tale sign this is the cause. Fixes #1 + #2 address this directly, and they're the highest priority for your usage pattern.

---

## Does Cloudflare change things? — Yes, but it's not the root cause

- **Buffering:** Cloudflare can buffer `text/event-stream` (flushing in ~100 KB chunks, or whenever it compresses), which clumps events. `X-Accel-Buffering: no` is an *nginx* hint CF may ignore — the CF-reliable lever is **`Cache-Control: no-cache, no-transform`** plus ensuring compression is **off** for the stream.
- **~100s idle timeout / 524s:** covered by the existing 30s heartbeat *as long as writes flush*.
- **Reconnect churn:** CF drops/edge-restarts cause more reconnects than localhost ever does, which is what *exposes* the token-refresh gap (#3) and the hub stall (#1).

The reason it's fine on localhost is that local connections never stall or half-die, so #1/#2 never trigger and tokens rarely expire during a short test.

---

## What robust SSE setups do

- **Never block the producer on a slow consumer** — non-blocking send (`select { case ch<-e: default: drop/close }`), and don't hold a global lock across delivery (snapshot subscribers under lock, deliver outside it).
- **Per-connection write deadlines** (`http.NewResponseController(w).SetWriteDeadline(...)`, refreshed each write/heartbeat) so dead sockets fail fast and the goroutine exits + unsubscribes — instead of hanging.
- **Heartbeats** every 15–30s (already done ✅).
- **Reconnect with exponential backoff + jitter and effectively unlimited retries**, giving up only on a definitive auth failure — not a hard cap of 3.
- **Auth that survives token rotation:** refresh before (re)connecting, or move off native `EventSource` to a fetch-based SSE client (e.g. `@microsoft/fetch-event-source`) that can send an `Authorization` header and do refresh-then-retry on 401 (also removes the token-in-URL, which currently leaks into CF/access logs).
- **Reconcile, don't trust the stream alone:** refetch authoritative state on (re)connect and on window-focus so a dropped event self-heals. Keep absolute-value payloads.
- **Session-scoped connection:** open once for the authenticated session, not per page.

---

## Proposed fixes, in priority order

| # | Fix | File | Effort | Impact |
|---|-----|------|--------|--------|
| 1 | Non-blocking fan-out: `select { case s.eventQueue <- e: default: log + mark conn for close }`; stop holding `em.mu` across delivery | `events/sse.go`, `events/events.go` | M | Removes the hub-wide stall + burst — biggest win |
| 2 | Per-write deadline in `SSEHandler` via `http.NewResponseController`, so a dead socket errors in ~5–10s and the goroutine exits | `events/sse.go` | S | Dead clients can no longer hang the hub |
| 3 | Add `Cache-Control: no-cache, no-transform` and confirm compression is off for `/api/events` in Cloudflare | `events/sse.go` + CF config | S | Stops Cloudflare buffering/clumping |
| 4 | Move the connection to app/session scope (open on login in `App.vue`/an auth-aware provider, close on logout) | `App.vue`, `useSSE.ts` | M | Live updates persist across navigation |
| 5 | Refresh the access token before (re)connecting if expired/near-expiry; rebuild URL with the fresh token | `useSSE.ts` | M | Fixes "bad access token"; remove the hard 3-retry cap → backoff+jitter |
| 6 | On (re)connect and window-focus, refetch the user/levels query so a missed `xp_updated` self-heals | `useSSE.ts` | S | Number becomes eventually-consistent even with drops |
| 7 | (Longer term) replace native `EventSource` with a fetch-based SSE client to use an `Authorization` header + clean 401→refresh→retry | `useSSE.ts` | L | Removes token-in-URL and the refresh gap entirely |

**Sequencing:** #1–#3 are small backend/config changes that kill the worst symptom (stall→burst) and the CF buffering — do these first and validate in prod. Then #4–#6 for the "doesn't update / bad token" class. #7 is the clean long-term auth story.

### Progress

- **#1 — done** (`events/sse.go`, `events/events.go`): `SSEConn.HandleEvent` now does a non-blocking send (drops + logs on a full buffer); `processEvents` snapshots subscribers under the lock and delivers off the lock, so no subscriber can freeze the hub.
- **#2 — done** (`events/sse.go`): each write is armed with a `SetWriteDeadline` via `http.NewResponseController` (`sseWriteTimeout = 10s`); a stalled socket now errors fast and the handler unsubscribes, instead of hanging. Degrades gracefully where deadlines are unsupported.
- **#3 — partial:** origin header now `Cache-Control: no-cache, no-transform` (done). Remaining: confirm/disable compression + buffering for `/api/events` in the Cloudflare console (manual).
- **Code-side fix for the "fresh login freezes while an old tab is wedged" symptom is #1 + #2 — both landed.**
- Verified green: `go build ./...`, `go test ./internal/events/... -race`, `go vet`, `make unit`, `golangci-lint` (0 issues).
- **Pending:** #4 (session-scoped connection), #5 (token refresh / drop the hard 3-retry cap), #6 (refetch-on-reconnect), #7 (cookie or fetch-based auth) — all frontend.
- **Note:** `events/websocket.go` has the same blocking-send pattern in `webSocketConn.HandleEvent`. It's not on the wired path (frontend uses SSE), so left as-is; flagging for if WS is ever used.

---

## Decisions (resolved)

- **The number = XP** (`xp_updated` → `ProgressBar`). No separate tasks-completed stat to trace.
- **API is same-origin** as the frontend (both behind Cloudflare). This unlocks the cleanest auth fix: the SSE request can authenticate via an **HttpOnly session cookie** instead of `?token=` — same-origin `EventSource` sends cookies automatically, so it survives token rotation, removes the token-from-URL leak, and avoids the "rebuild URL with a fresh token" dance. (If we'd rather not introduce a cookie, the fetch-based-SSE option in #7 is the alternative.) This makes #5/#7 simpler and is worth doing once #1–#3 land.
- **Refetch-on-reconnect (#6) is approved** as a self-healing safety net.
- **Cloudflare:** we have web-console access and will confirm/disable compression + buffering for `/api/events` there, in addition to the origin `no-transform` header (#3).

---

## Note on how this was produced

A few file reads during this session returned bogus content for paths that don't exist on disk (`frontend/src/composables/useSSE.ts`, `frontend/src/features/notifications/*`, `backend/internal/handlers/sse.go`, `backend/internal/server/server.go`). Existence was re-checked with `test -f`/`ls` and a clean sub-agent; this audit reflects only the verified, wired files listed at the top.
