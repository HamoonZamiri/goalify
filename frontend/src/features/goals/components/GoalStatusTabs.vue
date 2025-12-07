<script setup lang="ts">
import { computed } from "vue";
import type { Goal } from "@/features/goals/schemas";
import { Tabs, Tab } from "@/shared/components/tabs";
import { Button } from "@/shared/components/ui";

const props = defineProps<{
	status: Goal["status"];
}>();

const emit = defineEmits<{
	"update:status": [status: Goal["status"]];
}>();

const selectedIndex = computed(() => (props.status === "complete" ? 1 : 0));

function handleTabChange(index: number) {
	const newStatus: Goal["status"] = index === 1 ? "complete" : "not_complete";
	emit("update:status", newStatus);
}
</script>

<template>
	<Tabs
		variant="bordered"
		:modelValue="selectedIndex"
		@update:modelValue="handleTabChange"
	>
		<template #tabs>
			<Tab v-slot="{ selected }">
				<Button
					:variant="selected ? 'secondary' : 'ghost'"
					width="w-auto"
					class="tab flex-1"
				>
					In Progress
				</Button>
			</Tab>
			<Tab v-slot="{ selected }">
				<Button
					:variant="selected ? 'secondary' : 'ghost'"
					width="w-auto"
					class="tab flex-1"
				>
					Completed
				</Button>
			</Tab>
		</template>
	</Tabs>
</template>
