# Research: MAI-66 — fix: routing /login does not work on UI

**Ticket:** MAI-66

**Date:** 2026-03-08

## Scope
Investigate why direct navigation to `/login` and `/signup` breaks in production (Netlify 404) and behaves incorrectly locally (wrong content rendered, empty page).

## Codebase Findings

### Relevant Files
- `frontend/src/router/index.ts` — route definitions + global `beforeEach` guard
- `frontend/src/App.vue` — root layout, inline auth redirect in setup script
- `frontend/src/features/auth/pages/LoginPage.vue` — login form page
- `frontend/src/features/auth/pages/RegisterPage.vue` — register form page
- `frontend/src/shared/hooks/auth/useAuth.ts` — auth state (cookie-based)
- `frontend/vite.config.ts` — build config, no SPA fallback configured
- `frontend/public/` — no `_redirects` file present
- No `netlify.toml` in repo root

### Existing Patterns

**Router:** `createWebHistory()` — uses the browser History API, correct for SPA. Three routes defined:
- `/login` → `LoginPage`
- `/register` → `RegisterPage`
- `/` → `HomePage`

**Navigation guard (`beforeEach`):**
```typescript
router.beforeEach(async (to, _) => {
  const { default: useAuth } = await import("@/shared/hooks/auth/useAuth");
  const { getUser } = useAuth();
  const user = getUser();

  if (!user && to.name !== "Login" && to.name !== "Register") {
    return { name: "Login" };
  }
  if (user && (to.name === "Login" || to.name === "Register")) {
    return { name: "Home", path: "/" };
  }
});
```

**Auth state:** Cookie-based. `getUser()` parses a `user` cookie and validates with Zod. Reactive `authState` ref caches it in memory.

**App.vue layout:** Always renders a full-screen layout with a `<Navbar>`, conditional `<Sidebar v-if="isLoggedIn()">`, and `<RouterView>`. The setup script also runs an inline redirect:
```typescript
if (!isLoggedIn()) {
  router.push({ name: "Login" });
}
```

### Constraints & Gotchas

**Root cause 1 — Netlify 404:**
No `_redirects` file in `public/` and no `netlify.toml`. Netlify serves static files; when a user hits `/login` directly, it looks for a literal `/login` file and returns 404. SPAs need a catch-all redirect rule: `/* /index.html 200`.

**Root cause 2 — Localhost `/login` shows dashboard + URL not updating:**
`App.vue` runs `router.push({ name: "Login" })` synchronously in setup if not logged in. But if the user *is* logged in and navigates to `/login`, the `beforeEach` guard correctly redirects to `Home` — however the guard is async (dynamic import of `useAuth`) which can cause a timing gap. The redundant push in `App.vue` and the guard fire in an inconsistent order.

**URL not reflecting redirects** is a direct consequence of this: `router.push()` called inside a component's setup script fires *after* navigation has already settled, initiating a second navigation rather than replacing the current one. This causes the URL and the rendered content to fall out of sync. The correct approach is `return { name: "..." }` from `beforeEach`, which intercepts navigation before it completes and updates the URL atomically.

**Root cause 3 — Localhost `/signup` shows empty page with sidebar:**
`/signup` is not a defined route — the route is named `Register` and lives at `/register`, not `/signup`. Navigating to `/signup` hits a route that doesn't exist. With `createWebHistory()` and Vite's dev server (which implicitly serves `index.html` as fallback), the app loads but Vue Router finds no matching route, rendering nothing in `<RouterView>`. The sidebar appears because the layout in `App.vue` always renders.

**No 404/catch-all route defined** in the router, so unmatched paths silently render an empty `<RouterView>`.

**The App.vue inline redirect is redundant** — the `beforeEach` guard already handles this. Having both creates potential conflicts and makes auth flow harder to reason about.

## External Research
Not required — all issues are diagnosable from the codebase.

## Decisions
- **`/signup`** — no alias or redirect; unmatched routes return a 404 page via catch-all route
- **Catch-all route** — add `{ path: "/:pathMatch(.*)*", name: "NotFound", component: NotFoundPage }` to router
- **App.vue inline redirect** — remove entirely; industry best practice is for the router guard (`beforeEach`) to be the single source of truth for auth-based navigation. Component-level `router.push()` in setup fires after navigation settles, which is what causes the URL desync
- **URL not updating** — addressed by removing the `router.push()` from App.vue and relying solely on `beforeEach` returning `{ name: "..." }`, which intercepts navigation before completion and keeps URL in sync
