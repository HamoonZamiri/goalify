<script setup lang="ts">
import { onMounted, onUnmounted, ref } from "vue";
import { GoalCategoryCard } from "@/features/goals";
import CreateGoalCategoryDialog from "@/features/goals/components/CreateGoalCategoryDialog.vue";
import { useGoalCategories } from "@/features/goals/queries";
import { ProgressBar } from "@/features/levels";
import { Box, Text, Button } from "@/shared/components/ui";
import useAuth from "@/shared/hooks/auth/useAuth";
import { useSSE } from "@/shared/hooks/events/useSse";
import { API_BASE } from "@/utils/constants";

const isCreateCategoryDialogOpen = ref(false);

const { getUser } = useAuth();
const { data: categories, isLoading, error } = useGoalCategories();

const { connect, closeConnection } = useSSE();

onMounted(() => {
	connect(`${API_BASE}/events?token=${getUser()?.access_token}`);
});

onUnmounted(() => {
	closeConnection();
});
</script>

<template>
	<div class="skeleton h-full w-full rounded-full" v-if="isLoading"/>
	<Box
		v-else-if="error"
		height="h-full"
		bg="darkest"
		class="items-center justify-center"
	>
		<p class="text-red-500">Error loading categories: {{ error.message }}</p>
	</Box>
	<Box
		v-else
		height="h-full"
		bg="darkest"
		class="items-center sm:items-start px-6 w-full overflow-hidden"
	>
		<Box flex-direction="row" width="w-full" gap="gap-4">
			<Box gap="gap-4" bg="darkest" flex-direction="col" width="w-full">
				<Text as="h1" size="3xl" weight="bold">Dashboard</Text>
				<Text as="h2" size="xl" weight="semibold">My Goals </Text>
				<ProgressBar/>
				<Box
					width="w-full"
					bg="darkest"
					class="w-full"
					v-for="cat in categories"
					:key="cat.id"
				>
					<GoalCategoryCard :goalCategory="cat"/>
				</Box>
				<Box bg="darkest" flex-direction="row">
					<Button @click="isCreateCategoryDialogOpen = true" variant="primary">
						<Text>Add Goal Category</Text>
					</Button>
				</Box>
				<CreateGoalCategoryDialog v-model="isCreateCategoryDialogOpen"/>
			</Box>
		</Box>
	</Box>
</template>
