<script setup lang="ts">
import { computed } from "vue";
import { Box, Text } from "@/shared/components/ui";
import { useAuth } from "@/shared/hooks";
import { useLevelInfoQuery } from "../queries";

const { authState } = useAuth();
const { data: level } = useLevelInfoQuery(
	{ levelId: authState.value?.level_id ?? 0 },
	{
		enabled: !!authState.value,
	},
);

const progressPercent = computed(() => {
	if (!authState.value || !level.value) return 0;
	const user = authState.value;
	return ((user.xp / level.value.level_up_xp) * 100).toFixed(0);
});
</script>
<template>
	<Box flex-direction="col" gap="gap-2" class="bg-inherit">
		<Box
			flex-direction="row"
			class="justify-between"
			width="w-full"
			bg="darkest"
		>
			<Text size="base">Level {{ level?.id }}</Text>
			<Text size="sm">{{ progressPercent }}%</Text>
		</Box>
		<Box bg="darkest" width="w-full">
			<progress
				class="rounded-md transition-all duration-1000 ease-in-out"
				:value="progressPercent"
				max="100"
			/>
		</Box>
		<Box flex-direction="row" width="w-full" bg="darkest">
			<Text size="base">
				{{ `${authState?.xp} / ${level?.level_up_xp} XP` }}
			</Text>
		</Box>
	</Box>
</template>
