<script setup lang="ts">
import router from "./router";
import { isLoggedIn, removeUser } from "./utils/user";

function handleLogout(e: MouseEvent) {
  e.preventDefault();
  removeUser();
  router.push({ name: "Login" });
}
</script>

<template>
  <header class="bg-emerald-100 mb-2">
    <div class="flex text-blue-400 justify-between p-6">
      <RouterLink to="/">
        <h1 class="font-semibold text-3xl hover:text-blue-500">Goalify</h1>
      </RouterLink>
      <nav class="flex gap-4">
        <div v-if="!isLoggedIn()" class="flex gap-2">
          <RouterLink
            class="font-semibold text-xl hover:text-blue-500"
            to="/login"
            >Login</RouterLink
          >
          <RouterLink
            class="font-semibold text-xl hover:text-blue-500"
            to="/register"
            >Register</RouterLink
          >
        </div>
        <div v-else>
          <button
            @click="handleLogout"
            class="font-semibold text-xl hover:text-blue-500"
          >
            Log Out
          </button>
        </div>
      </nav>
    </div>
  </header>
  <main><RouterView /></main>
</template>
