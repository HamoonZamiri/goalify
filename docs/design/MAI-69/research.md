# Research: MAI-69 — feat(ui): make ui vimium friendly

**Ticket:** MAI-69
**Date:** 2026-03-22

## Scope

Explored the entire frontend `src/` tree to identify every interactive/clickable element and determine which ones are non-semantic (i.e., not rendered as `<button>` or `<a>`) and therefore invisible to Vimium's `f` hint mode. Two codebase scans confirmed completeness.

---

## Codebase Findings

### Relevant Files

| File | Issue |
|------|-------|
| `features/goals/components/GoalCategoryCard/GoalCategoryCard.vue` | `<DisclosureButton as="div">` renders a `<div>` — Vimium cannot target it |
| `features/goals/components/GoalCard.vue` | Whole card is a `<Box>` (div) with `@click.stop` — not Vimium accessible |
| `shared/components/ui/Box.vue` | Always renders `<div>` — no way to opt into semantic element at callsite |
| `shared/components/ui/Button.vue` | Correct — renders `<button>` |
| `shared/components/icons/IconButton.vue` | Correct — wraps `Button.vue` |
| `shared/components/navigation/ProfileMenu.vue` | Correct — Headless UI `<MenuButton>` and `<MenuItem as="button">` |
| `shared/components/navigation/SideBar.vue` | Correct — `<RouterLink>` for nav items, `<Button>` for sign out |
| `features/goals/components/GoalStatusTabs.vue` | Correct — Headless UI tabs render as `<button>` |

**Scope note:** Settings and Rewards pages do not exist yet. This ticket is scoped entirely to the Goal Category accordion (`GoalCategoryCard`) and the goal rows within it (`GoalCard`).

---

### Problem 1: GoalCategoryCard accordion header

**File:** `features/goals/components/GoalCategoryCard/GoalCategoryCard.vue:125`

```vue
<DisclosureButton
  as="div"    ← overrides Headless UI default, renders as <div>
  class="w-full text-left hover:bg-gray-800 ..."
  @click="(e: MouseEvent) => handleDisclosureClick(e)"
>
```

Headless UI's `DisclosureButton` renders as `<button>` by default. The `as="div"` override was likely added to avoid default button styles/focus rings. Removing it restores the native `<button>` rendering, which Vimium can target.

**Decision:** Remove `as="div"` and rely on the Headless UI default. Verify visually that existing Tailwind classes handle any style regressions — no preemptive style changes until confirmed broken.

---

### Problem 2: GoalCard click target

**File:** `features/goals/components/GoalCard.vue:47`

```vue
<!-- Current: entire card is a div with a click handler -->
<Box
  @click.stop="isEditing = true"
  class="hover:cursor-pointer ..."
  flex-direction="row" gap="gap-4"
>
  <IconButton icon="check-outline" @click.stop="handleCheckClick" />
  <Text>{{ goal.title }}</Text>
</Box>
```

The desired UX: clicking the check icon completes the goal; clicking anywhere else on the card opens the edit dialog. Cannot simply make the outer Box a `<button>` because `<IconButton>` is already a `<button>` — nested buttons are invalid HTML.

**Decision:** Restructure as two sibling flex children inside a non-interactive outer Box:

```vue
<!-- After -->
<Box flex-direction="row" gap="gap-0" ...>   ← no click handler, no cursor
  <IconButton @click.stop="handleCheckClick" />
  <Box as="button" @click.stop="isEditing = true" class="flex-1 ...">
    <Text>{{ goal.title }}</Text>
  </Box>
</Box>
```

- Outer `<Box>` is a plain `<div>` — purely for layout
- `<IconButton>` handles check/complete (already correct)
- `<Box as="button">` takes up all remaining space (`flex-1`) — clicking anywhere on the card except the icon hits this button and opens the edit dialog
- No nested buttons — the two interactive elements are siblings

This relies on the polymorphic `Box` change below.

---

### Problem 3: Box.vue — always renders `<div>`

**File:** `shared/components/ui/Box.vue`

`Box.vue` hardcodes `<div>`. Since it uses `v-bind="$attrs"`, click handlers pass through functionally, but the DOM element is never semantic — making any "clickable Box" invisible to Vimium.

**Decision:** Add an optional `as` prop to `Box.vue` (default `"div"`) using Vue's `<component :is>` pattern. Callsites that need a semantic element opt in explicitly with `as="button"` or `as="a"`. No existing callsites change behavior unless explicitly updated.

```vue
<!-- Box.vue after change -->
<script setup lang="ts">
const props = withDefaults(defineProps<{ as?: string; /* ...existing props */ }>(), {
  as: "div",
  // ...existing defaults
})
const tag = computed(() => props.as)
</script>
<template>
  <component :is="tag" :class="[...]" v-bind="$attrs">
    <slot />
  </component>
</template>
```

**Why default stays `"div"`:** Zero risk of regression — all existing `<Box>` usage is unchanged. Only new/updated callsites with `as="button"` get button semantics.

---

### Existing Patterns

- **Headless UI** is already in use for Disclosure, Menu, and Tab — all render as `<button>` by default
- `Button.vue` and `IconButton.vue` are the canonical button primitives
- `Box.vue` uses `v-bind="$attrs"` — the polymorphic `as` prop fits naturally alongside this

### Constraints & Gotchas

- **`DisclosureButton` style regression risk:** Removing `as="div"` may introduce browser-default button styles (border, background, padding). `w-full`, `text-left`, and existing classes likely cover it but needs visual verification.
- **Focus rings:** Switching to native `<button>` may introduce default browser focus outlines. Deferring focus ring styling decisions until we see what the actual result looks like.
- **`Box as="button"` default styles:** Like DisclosureButton, a `<button>` rendered via `Box` may pick up browser defaults. The `<component :is>` approach means Box's existing Tailwind classes apply, which should override most defaults — verify after implementation.

---

## External Research

**Vimium `f` hint behavior:**
- Generates hints for natively interactive elements: `<a href>`, `<button>`, `<input>`, `<select>`, `<textarea>`, and elements with `role="button"` or `tabindex`
- `<div @click>` is NOT hinted — Vimium only sees the DOM element type
- Native `<button>` is always preferred over `role="button"` + `tabindex="0"` ARIA fallback

**Headless UI DisclosureButton:**
- Default render target is `<button>` — manages `aria-expanded` automatically
- The `as` prop is an escape hatch, not the intended default usage

---

## Decisions Log

| # | Decision | Rationale |
|---|----------|-----------|
| 1 | Remove `as="div"` from `DisclosureButton` | Rely on Headless UI default (`<button>`); verify styles before any preemptive fixes |
| 2 | Add `as` prop to `Box.vue`, default `"div"` | Systemic fix; zero regression risk; callers opt in explicitly |
| 3 | GoalCard: outer div → `<IconButton>` + `<Box as="button">` siblings | Solves nested-button invalidity; preserves "click anywhere except check = open edit" UX |
| 4 | Scope: GoalCategoryCard + GoalCard only | Settings/Rewards pages not yet built |
| 5 | Focus ring styling: defer | Iterate after seeing actual rendered output |
