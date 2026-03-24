<script setup lang="ts">
/**
 * OverlayButton — a full-area invisible button for use inside ClickSurface.
 *
 * Renders as `position: absolute; inset: 0; z-index: 1`, covering the entire
 * parent ClickSurface. z-[1] ensures it sits above sibling flex items (which
 * stack at z: auto in a flex context) while remaining below any secondary
 * interactive children that carry `relative z-10`.
 *
 * All attrs (aria-label, @click, etc.) are forwarded to the underlying
 * <button> element. Always provide an `aria-label` for accessibility.
 */
</script>

<template>
	<button
		class="absolute inset-0 z-[1] cursor-pointer rounded-[inherit]"
		v-bind="$attrs"
		@pointerdown.stop
	/>
</template>
