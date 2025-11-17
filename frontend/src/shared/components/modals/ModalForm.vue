<script setup lang="ts">
import { Dialog, DialogPanel } from "@headlessui/vue";

defineProps<{
	modelValue: boolean;
}>();
const emit = defineEmits(["update:modelValue"]);
</script>
<template>
	<Dialog
		class="absolute inset-0 h-screen flex justify-center items-center hover:cursor-pointer bg-gray-600/50 z-10 w-screen  rounded-lg p-4"
		:open="modelValue"
		@close="() => emit('update:modelValue', false)"
	>
		<DialogPanel>
			<slot/>
		</DialogPanel>
	</Dialog>
</template>
