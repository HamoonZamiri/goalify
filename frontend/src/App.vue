<script setup lang="ts">
import router from "./router";
import authState from "./state/auth";

function handleLogout(e: MouseEvent) {
  e.preventDefault();
  authState.logout();
  router.push({ name: "Login" });
}
</script>

<template>
  <div class="bg-gray-900 h-screen w-screen flex flex-col">
    <header class="bg-gray-800 mb-2 h-auto">
      <div class="flex text-gray-200 justify-between p-6">
        <RouterLink to="/">
          <h1 class="font-semibold text-3xl hover:text-gray-300">Goalify</h1>
        </RouterLink>
        <nav class="flex gap-4">
          <div v-if="!authState.getUser" class="flex gap-2">
            <RouterLink class="text-xl hover:text-gray-300" to="/login"
              >Login</RouterLink
            >
            <RouterLink class="text-xl hover:text-gray-300" to="/register"
              >Register</RouterLink
            >
          </div>
          <div v-else>
            <button @click="handleLogout" class="text-xl hover:text-gray-100">
              Log Out
            </button>
          </div>
        </nav>
      </div>
    </header>
    <div class="w-full h-full"><RouterView /></div>
  </div>
</template>
