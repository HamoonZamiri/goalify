# Feature: Vimium-friendly UI

**Status:** done

**Ticket:** MAI-69

**Branch:** hamoondev/mai-69-featui-make-ui-vimium-friendly

**Created:** 2026-03-22

## Goal

Make all clickable elements in the Goal Category accordion and goal rows render as native `<button>` elements so Vimium's `f` hint mode can target them.

## Context

Vimium only generates hints for natively interactive DOM elements (`<button>`, `<a>`, etc.). Two elements in the goals UI are rendered as `<div>` despite being clickable — the Goal Category accordion header (`DisclosureButton as="div"`) and the GoalCard row (`Box` with `@click`) — making them unreachable via keyboard-driven Vimium navigation.

## Out of Scope

- Settings and Rewards pages (not yet built)
- Full keyboard tab-stop audit
- Focus ring styling (defer — iterate after seeing rendered output)

## Approach

Three targeted changes:

**1. Box.vue — add polymorphic `as` prop**

Add an optional `as` prop (default `"div"`) using `<component :is>`. All existing callsites are unaffected. New callsites opt in explicitly with `as="button"`.

```vue
const props = withDefaults(defineProps<{ as?: string; /* existing */ }>(), {
  as: "div",
  // existing defaults
})
```

```vue
<component :is="props.as" :class="[...]" v-bind="$attrs">
  <slot />
</component>
```

**2. GoalCategoryCard.vue — remove `as="div"` from DisclosureButton**

```vue
<!-- Before -->
<DisclosureButton as="div" ...>

<!-- After -->
<DisclosureButton ...>
```

Headless UI defaults to `<button>` and manages `aria-expanded` automatically. Verify visually — fix styles only if regressions appear.

**3. GoalCard.vue — restructure into two sibling interactive elements**

Outer `<Box>` becomes a layout-only `<div>`. The check icon and the card body become siblings, not parent/child — avoiding invalid nested `<button>` HTML.

```vue
<!-- Before: entire card is a clickable div -->
<Box @click.stop="isEditing = true" class="hover:cursor-pointer hover:bg-gray-700 ...">
  <IconButton @click.stop="handleCheckClick" />
  <Text>{{ goal.title }}</Text>
</Box>

<!-- After: two sibling interactive elements, outer Box owns hover state via Tailwind group -->
<Box flex-direction="row" class="group rounded-xl ...">
  <IconButton @click.stop="handleCheckClick" />
  <Box as="button" @click.stop="isEditing = true" class="flex-1 text-left ...">
    <Text>{{ goal.title }}</Text>
  </Box>
</Box>
```

**Hover state ownership:** The outer `<Box>` is a plain `<div>` with no click handler. It gets Tailwind's `group` class and `group-hover:bg-gray-700`. Hovering anywhere over either child — the icon or the title — triggers the background highlight on the outer container, preserving the existing whole-card hover behaviour. Neither child owns a `hover:bg-*` class directly.

- Clicking the check icon → completes goal (unchanged)
- Clicking anywhere else on the row → opens edit dialog
- Both are native `<button>` elements — Vimium can hint both

## Tasks

- [x] Add `as` prop to `Box.vue`
- [x] Remove `as="div"` from `DisclosureButton` in `GoalCategoryCard.vue`
- [x] Create `OverlayButton` and `ClickSurface` primitives in `shared/components/ui/`
- [x] Restructure `GoalCard.vue` using `ClickSurface` + `OverlayButton`
- [x] Restructure `GoalCategoryCard.vue` header using `ClickSurface` + `DisclosureButton` as overlay
- [x] Update `GoalCard.test.ts` to target `OverlayButton` by aria-label
- [x] Run `pnpm format:fix` and `pnpm build`
- [x] Manual Vimium test: confirm `f` hints appear on accordion header and goal rows

## Acceptance Criteria

- [ ] As a Vimium user, pressing `f` shows a hint on the Goal Category accordion header that expands/collapses it when activated
- [ ] As a Vimium user, pressing `f` shows a hint on each goal row that opens the edit dialog when activated
- [ ] As a Vimium user, pressing `f` shows a hint on the check icon that toggles goal completion when activated
- [ ] No visual regressions in the accordion or goal card layout

## Open Questions

None — all questions resolved in research.

## Decisions Log

| # | Decision | Rationale |
|---|----------|-----------|
| 1 | Remove `as="div"` from `DisclosureButton` | Headless UI default is `<button>`; `as="div"` was an unnecessary override |
| 2 | Add `as` prop to `Box.vue`, default `"div"` | Systemic fix; zero regression risk; callers opt in explicitly |
| 3 | GoalCard: outer div + `<IconButton>` sibling + `<Box as="button" flex-1>` sibling | Avoids nested buttons (invalid HTML); preserves "click anywhere except check = open edit" UX |
| 4 | Scope: GoalCategoryCard + GoalCard only | Settings/Rewards not yet built |
| 5 | Focus ring styling: defer | Iterate after seeing actual rendered output |

## Session Log

- **2026-03-22:** Research completed. Full codebase scan confirmed only two non-semantic clickable elements exist. Decisions finalized. Design doc written. Awaiting approval.
- **2026-03-22:** Implementation complete. Introduced `ClickSurface` + `OverlayButton` primitives to handle the overlay button pattern. `Box.vue` made polymorphic via `as` prop. `GoalCard` restructured with overlay + `#actions` slot for the check icon (now on the right). `GoalCategoryCard` header restructured with `DisclosureButton` as overlay and action buttons in `#actions` slot. All 24 tests passing. Pending manual Vimium verification.
