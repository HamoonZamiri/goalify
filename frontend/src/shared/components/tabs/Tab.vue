<script setup lang="ts">
import { Tab as HeadlessTab } from "@headlessui/vue";

defineOptions({
	inheritAttrs: false,
});

withDefaults(
	defineProps<{
		as?: string | object;
	}>(),
	{
		as: "template",
	},
);
</script>

<template>
	<HeadlessTab :as="as" v-slot="{ selected }">
		<slot :selected="selected"/>
	</HeadlessTab>
</template>
