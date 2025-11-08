<script setup lang="ts">
const props = defineProps<{
	onclick?: (e: MouseEvent) => void;
	fill?: string;
}>();

function handleClick(e: MouseEvent) {
	props?.onclick?.(e);
}
</script>
<template>
	<svg
		xmlns="http://www.w3.org/2000/svg"
		:fill="props.fill || 'none'"
		viewBox="0 0 24 24"
		stroke-width="1.5"
		class="size-6 stroke-gray-300"
		@click.stop="handleClick"
	>
		<path
			stroke-linecap="round"
			stroke-linejoin="round"
			d="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
		/>
	</svg>
</template>
