<script setup lang="ts">
/**
 * ClickSurface — container for the overlay button pattern.
 *
 * Establishes a `position: relative` stacking context so that an OverlayButton
 * placed in `#overlay` covers the entire surface area.
 *
 * Slots:
 *   #overlay  — place an <OverlayButton> (or equivalent) here. Renders as
 *               `absolute inset-0 z-[1]`, covering the full surface.
 *   #actions  — secondary interactive elements (e.g. checkbox, icon buttons).
 *               Automatically wrapped in `relative z-10` so they sit above the
 *               overlay and intercept their own clicks. Renders after the
 *               default slot (rightmost in a flex row).
 *   default   — non-interactive content (text, icons). Falls through to the
 *               overlay button — clicks here trigger the primary action.
 *
 * All other attrs (class, event listeners, data-* etc.) are forwarded to the
 * root element so callers control layout and visual styling normally.
 */
const slots = defineSlots<{
	overlay?(): unknown;
	actions?(): unknown;
	default?(): unknown;
}>();
</script>

<template>
	<div class="relative">
		<slot name="overlay" />
		<slot />
		<div v-if="slots.actions" class="relative z-10">
			<slot name="actions" />
		</div>
	</div>
</template>
