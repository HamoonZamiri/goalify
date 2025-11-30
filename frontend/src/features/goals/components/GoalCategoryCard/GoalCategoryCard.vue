<script setup lang="ts">
import { Disclosure, DisclosureButton, DisclosurePanel } from "@headlessui/vue";
import { usePointerSwipe } from "@vueuse/core";
import { computed, ref } from "vue";
import { toast } from "vue3-toastify";
import { CreateGoalButton, CreateGoalForm, GoalCard } from "@/features/goals";
import {
	useDeleteGoalCategory,
	useResetGoalCategory,
} from "@/features/goals/queries";
import type { GoalCategory } from "@/features/goals/schemas/goal.schema";
import { ChevronUp, Trash, ArrowPath } from "@/shared/components/icons";
import { ModalForm } from "@/shared/components/modals";
import { Box, Text, Button } from "@/shared/components/ui";

const props = defineProps<{
	goalCategory: GoalCategory;
}>();

const isCreateGoalDialogOpen = ref(false);

/**
 * Swipe-to-delete logic
 */
const swipeTarget = ref<HTMLElement>();
const dragDistance = ref(0);
const didSwipe = ref(false);
const DELETE_THRESHOLD = 0.65; // 65% of card width triggers delete

const { mutateAsync: deleteCategory } = useDeleteGoalCategory();
const { mutateAsync: resetCategory, isPending: isPendingReset } =
	useResetGoalCategory();

// Calculate how far we've swiped as a percentage
const swipeRatio = computed(() => {
	const cardWidth = swipeTarget.value?.offsetWidth ?? 0;
	return cardWidth > 0 ? Math.abs(dragDistance.value) / cardWidth : 0;
});

const { distanceX, isSwiping } = usePointerSwipe(swipeTarget, {
	onSwipe: () => {
		if (distanceX.value >= 0) {
			dragDistance.value = 0;
			return;
		}

		didSwipe.value = true;
		const maxSwipe = (swipeTarget.value?.offsetWidth ?? 0) * 0.8;
		dragDistance.value = Math.max(distanceX.value, -maxSwipe);
	},
	onSwipeEnd: async () => {
		if (swipeRatio.value >= DELETE_THRESHOLD) {
			try {
				await deleteCategory(props.goalCategory.id);
				toast.success(`Deleted category: ${props.goalCategory.title}`);
			} catch (error) {
				toast.error(
					`Failed to delete category: ${error instanceof Error ? error.message : "Unknown error"}`,
				);
			}
		}
		dragDistance.value = 0;
		// Reset swipe flag after a short delay (after snap animation)
		setTimeout(() => {
			didSwipe.value = false;
		}, 250);
	},
});

/**
 * Prevent disclosure toggle if user dragged
 * Use distanceX directly - it persists after swipe ends
 */
function handleDisclosureClick(e: MouseEvent) {
	// If any drag movement occurred (>5px to ignore micro-movements), prevent toggle
	if (didSwipe.value) {
		e.preventDefault();
		e.stopPropagation();
	}
}

async function handleResetCategory() {
	try {
		await resetCategory({ category_id: props.goalCategory.id });
	} catch (error) {
		toast.error(
			`Failed to reset category: ${error instanceof Error ? error.message : "Unknown error"}`,
			{ autoClose: 1500 },
		);
	}
}
</script>

<template>
	<Disclosure as="div" v-slot="{ open }">
		<!-- Swipe container -->
		<div class="relative overflow-hidden rounded-lg">
			<!-- Red background layer (only show when disclosure is closed) -->
			<div
				v-if="!open"
				class="absolute inset-0 bg-red-600 flex items-center pr-8"
				:class="swipeRatio >= DELETE_THRESHOLD ? 'justify-start' : 'justify-end'"
			>
				<Trash/>
			</div>

			<!-- Draggable card layer -->
			<div
				ref="swipeTarget"
				:style="{
					transform: `translateX(${dragDistance}px)`,
					willChange: isSwiping ? 'transform' : 'auto',
				}"
				:class="[
					'bg-gray-900',
					!isSwiping
						? 'transition-transform duration-200 ease-out'
						: '',
				]"
			>
				<Box
					flex-direction="col"
					padding="p-4"
					class="transition-colors rounded-lg"
				>
					<DisclosureButton
						as="div"
						class="w-full text-left hover:bg-gray-800 rounded-lg transition-colors hover:cursor-pointer"
						@click="(e: MouseEvent) => handleDisclosureClick(e)"
					>
						<header class="flex justify-between w-full items-center">
							<Text as="h3" weight="semibold">{{ goalCategory.title }}</Text>
							<Box class="items-center gap-2" flex-direction="row">
								<Button
									variant="ghost"
									width="w-auto"
									class="p-0"
									@click.stop="isCreateGoalDialogOpen = true"
								>
									<CreateGoalButton/>
								</Button>
								<Button
									variant="ghost"
									width="w-auto"
									class="p-0"
									:disabled="isPendingReset"
									@click.stop="handleResetCategory"
								>
									<ArrowPath :class="isPendingReset ? 'animate-spin' : ''"/>
								</Button>
								<ModalForm v-model="isCreateGoalDialogOpen">
									<CreateGoalForm
										:category-id="props.goalCategory.id"
										@close="isCreateGoalDialogOpen = false"
									/>
								</ModalForm>
								<ChevronUp :class="open ? 'rotate-180 transform' : ''"/>
							</Box>
						</header>
					</DisclosureButton>
					<DisclosurePanel class="transition w-full mt-4">
						<Box flex-direction="col" v-for="goal in goalCategory.goals">
							<GoalCard :goal="goal" :xp-per-goal="goalCategory.xp_per_goal"/>
						</Box>
					</DisclosurePanel>
				</Box>
			</div>
		</div>
	</Disclosure>
</template>
