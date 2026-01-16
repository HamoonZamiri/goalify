<script setup lang="ts">
import { computed } from "vue";
import useAuth from "@/shared/hooks/auth/useAuth";
import { useLevelInfoQuery } from "@/features/levels";
import ProfileMenu from "./ProfileMenu.vue";

const { isLoggedIn, getUser } = useAuth();
const user = getUser();

const { data: level } = useLevelInfoQuery({ levelId: user?.level_id ?? 0 });

const progressBar = computed(() => {
	if (!user || !level.value) return 0;
	return (user.xp / level.value.level_up_xp) * 100;
});
</script>

<template>
	<nav class="flex gap-4">
		<div v-if="!isLoggedIn()" class="flex gap-2">
			<RouterLink class="text-xl hover:text-gray-300" to="/login">
				Login
			</RouterLink>
			<RouterLink class="text-xl hover:text-gray-300" to="/register">
				Register
			</RouterLink>
		</div>
		<div v-else class="flex gap-2 items-center">
			<span class="text-sm text-gray-300">{{ user?.email.split("@")[0] }}</span>
		</div>
	</nav>
</template>
