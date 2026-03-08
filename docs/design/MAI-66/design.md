# Fix: Routing — /login does not work on UI

**Status:** in-progress

**Ticket:** MAI-66

**Branch:** hamoondev/mai-66-fix-routing-login-does-not-work-on-ui

**Created:** 2026-03-08

Status definitions:
  - draft: Planning/designing, not yet approved for implementation
  - approved: Design approved, ready to start implementation
  - in-progress: Actively being worked on
  - review: Implementation complete, undergoing self-review or awaiting PR review
  - done: PR merged, acceptance criteria met

## Goal
Fix SPA routing so that direct navigation to `/login` works in production on Netlify, auth redirects correctly update the URL, and unmatched routes show a 404 page.

## Context
Three distinct bugs identified in research:
1. **Netlify 404** — no `_redirects` or `netlify.toml`, so Netlify can't serve `index.html` for deep links
2. **URL desync on redirect** — `App.vue` calls `router.push()` in setup after navigation has settled, causing the URL and rendered content to fall out of sync
3. **Unmatched routes render blank** — no catch-all route, so `/signup` (and any unknown path) silently renders an empty `<RouterView>` with the sidebar still visible

## Out of Scope
- Any changes to the login or register form logic
- Backend auth changes
- Adding a `/signup` alias route

## Approach

### 1. Fix Netlify SPA routing
Create `frontend/public/_redirects`:
```
/* /index.html 200
```
This tells Netlify to serve `index.html` for any path, letting Vue Router handle routing client-side.

### 2. Remove inline redirect from App.vue
Delete the `router.push({ name: "Login" })` call from `App.vue`'s setup script. The `beforeEach` guard in `router/index.ts` already handles this correctly — it intercepts navigation before it completes and returns `{ name: "Login" }`, which atomically updates the URL. Having both causes a race condition and the URL desync bug.

### 3. Add catch-all 404 route
Add a `NotFoundPage.vue` component and register a catch-all route as the last entry in the router:
```typescript
{
  path: "/:pathMatch(.*)*",
  name: "NotFound",
  component: NotFoundPage,
}
```
The `beforeEach` guard must also be updated to allow the `NotFound` route through for unauthenticated users (currently any non-Login/Register route redirects to Login, which would swallow 404s):
```typescript
if (!user && to.name !== "Login" && to.name !== "Register" && to.name !== "NotFound") {
  return { name: "Login" };
}
```

### 4. NotFoundPage component
Simple page, consistent with existing layout. No sidebar (unauthenticated users can hit 404s). Just a centred message with a link back to `/`.

## Approach (amended)

### 5. Refactor guard to static import
The `beforeEach` guard used a dynamic `await import()` for `useAuth` — no circular dependency exists so this is unnecessary. The async nature introduced a timing gap where navigation could partially settle before the redirect fired, causing the URL to not update. Switch to a static top-level import and make the guard synchronous.

### 6. 404 full-page takeover
`NotFoundPage` should take over the full viewport — no navbar, no sidebar. A 404 is a dead end; showing authenticated navigation alongside it is misleading. Gate the App.vue layout wrapper in `<RouterView>` with `v-if="route.name !== 'NotFound'"`.

### 7. Router unit tests
Write `frontend/src/router/index.test.ts` covering all guard scenarios as executable success criteria:
- Unauthenticated → `/` redirects to `/login`
- Unauthenticated → `/login` stays on `/login`
- Unauthenticated → `/register` stays on `/register`
- Unauthenticated → unknown path lands on `NotFound` (no redirect to login)
- Authenticated → `/login` redirects to `/`
- Authenticated → `/register` redirects to `/`
- Authenticated → `/` stays on `/`

## Tasks
- [x] Create `frontend/public/_redirects` with `/* /index.html 200`
- [x] Remove `router.push({ name: "Login" })` from `App.vue` setup script
- [x] Create `frontend/src/pages/NotFoundPage.vue`
- [x] Add catch-all route to `frontend/src/router/index.ts`
- [x] Update `beforeEach` guard to allow `NotFound` through for unauthenticated users
- [ ] Refactor `beforeEach` guard to use static import (fix URL desync)
- [ ] Update `App.vue` to hide layout when on `NotFound` route
- [ ] Write router unit tests (`frontend/src/router/index.test.ts`)
- [ ] Verify locally: direct navigation to `/login` loads login page with correct URL
- [ ] Verify locally: direct navigation to `/signup` shows 404 page
- [ ] Verify locally: auth redirect from `/login` (when logged in) updates URL to `/`
- [ ] Verify in production (Netlify): direct navigation to `/login` no longer 404s

## Acceptance Criteria
- [ ] Navigating directly to `/login` in production does not return 404
- [ ] Navigating directly to `/login` locally loads the login page with the URL correctly showing `/login`
- [ ] When a logged-in user navigates to `/login`, the URL updates to `/` (not stuck on `/login`)
- [ ] Navigating to any undefined route (e.g. `/signup`) shows a 404 page, not a blank page
- [ ] 404 page takes over the full viewport — no navbar or sidebar rendered
- [ ] Auth guard redirects keep the URL and rendered content in sync
- [ ] All router unit tests pass

## Open Questions
None — all decisions resolved in research.

## Decisions Log
<!-- Append-only: key choices made -->

### 2026-03-08
- Resolved in research: see `docs/design/MAI-66/research.md` Decisions section

## Session Log
<!-- Append-only: progress updates for context recovery -->
