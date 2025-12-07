<script setup lang="ts">
import { TabPanel as HeadlessTabPanel } from "@headlessui/vue";

defineOptions({
	inheritAttrs: false,
});
</script>

<template>
	<HeadlessTabPanel v-slot="{ selected }" :class="['tab-content', $attrs.class]">
		<slot :selected="selected" />
	</HeadlessTabPanel>
</template>
