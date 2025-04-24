<script setup lang="ts">
import useApi from "@/hooks/api/useApi";
import useAuth from "@/hooks/auth/useAuth";
import { onMounted, ref } from "vue";
import ProfileMenu from "@/components/navbar/ProfileMenu.vue";

const { isLoggedIn, getUser } = useAuth();
const { getLevel, isError } = useApi();
const user = getUser();

const progressBar = ref<number | undefined>(0);

onMounted(async () => {
  if (!user) {
    return;
  }
  const res = await getLevel(user.level_id);
  if (isError(res)) {
    alert("Error getting level");
    return;
  }
  const level = res;
  const currXp = user.xp;
  setTimeout(() => {
    progressBar.value = (currXp / level.level_up_xp) * 100;
  }, 100);
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
