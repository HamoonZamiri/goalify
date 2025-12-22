<script setup lang="ts">
import Box from "./Box.vue";
import Button from "./Button.vue";
import Text from "./Text.vue";

withDefaults(
	defineProps<{
		title?: string;
		message?: string;
	}>(),
	{
		title: "Are you sure you want to delete this?",
		message: "This action cannot be undone.",
	},
);

const emit = defineEmits<{
	cancel: [];
	confirm: [];
}>();
</script>

<template>
	<Box
		flex-direction="col"
		gap="gap-4"
		class="border border-red-500 rounded-lg p-4"
	>
		<Text size="lg" weight="semibold" color="error">{{ title }}</Text>
		<Text size="sm" color="light">{{ message }}</Text>
		<Box flex-direction="row" gap="gap-4" class="justify-end">
			<Button variant="ghost" @click="emit('cancel')">
				<Text>Cancel</Text>
			</Button>
			<Button
				variant="secondary"
				class="hover:bg-red-600"
				@click="emit('confirm')"
			>
				<Text>Delete</Text>
			</Button>
		</Box>
	</Box>
</template>
