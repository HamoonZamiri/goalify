<script setup lang="ts">
import { ref } from "vue";
import { Disclosure, DisclosureButton, DisclosurePanel } from "@headlessui/vue";
import type { GoalCategory } from "@/features/goals/schemas/goal.schema";
import { Box, Text } from "@/shared/components/ui";
import { ChevronUp } from "@/shared/components/icons";
import { ModalForm } from "@/shared/components/modals";
import { GoalCard, CreateGoalButton, CreateGoalForm } from "@/features/goals";

const props = defineProps<{
	goalCategory: GoalCategory;
}>();

const isCreateGoalDialogOpen = ref(false);
</script>

<template>
	<Disclosure as="div" v-slot="{ open }">
		<DisclosureButton class="w-full text-left">
			<Box
				flex-direction="col"
				padding="p-4"
				class="hover:bg-gray-800 transition-colors rounded-lg"
			>
				<header class="flex justify-between w-full items-center">
					<Text as="h3" weight="semibold">{{ goalCategory.title }}</Text>
					<Box class="items-center gap-2" flex-direction="row">
						<CreateGoalButton
							class="hover:cursor-pointer"
							@click.stop="isCreateGoalDialogOpen = true"
						/>
						<ModalForm v-model="isCreateGoalDialogOpen">
							<CreateGoalForm
								:category-id="props.goalCategory.id"
								@close="isCreateGoalDialogOpen = false"
							/>
						</ModalForm>
						<ChevronUp :class="open ? 'rotate-180 transform' : ''"/>
					</Box>
				</header>
				<DisclosurePanel class="transition w-full mt-4">
					<Box flex-direction="col" v-for="goal in goalCategory.goals">
						<GoalCard :goal="goal" :xp-per-goal="goalCategory.xp_per_goal"/>
					</Box>
				</DisclosurePanel>
			</Box>
		</DisclosureButton>
	</Disclosure>
</template>
