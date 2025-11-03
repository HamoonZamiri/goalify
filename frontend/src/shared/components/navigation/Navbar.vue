<script setup lang="ts">
import { computed } from "vue";
import useAuth from "@/shared/hooks/auth/useAuth";
import { useLevelInfo } from "@/features/levels";
import ProfileMenu from "./ProfileMenu.vue";

const { isLoggedIn, getUser } = useAuth();
const user = getUser();

const { data: level } = useLevelInfo(user?.level_id ?? 0);

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
      <RouterLink class="text-xl hover:text-gray-300" to="/register"
        >Register</RouterLink
      >
    </div>
    <div v-else class="flex gap-2 items-center">
      <ProfileMenu />
      <section class="flex items-center gap-x-1.5 w-48">
        <span class="text-sm text-gray-300">{{
          user?.email.split("@")[0]
        }}</span>
        <div class="w-full h-5 bg-gray-200 rounded-md overflow-hidden">
          <div
            class="h-full bg-green-500 transition-all duration-1000 ease-in-out"
            :style="{ width: `${progressBar}%` }"
          ></div>
        </div>
        <span
          class="text-xs w-14 h-5 rounded-md flex items-center justify-center bg-gray-200 text-gray-700"
          >lvl {{ user?.level_id }}</span
        >
      </section>
    </div>
  </nav>
</template>
