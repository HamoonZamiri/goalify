<script setup lang="ts">
import useAuth from "@/hooks/auth/useAuth";
import router from "@/router/index";
import Navbar from "@/components/navbar/Navbar.vue";

const { isLoggedIn } = useAuth();
if (!isLoggedIn()) {
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
        <Navbar />
      </div>
    </header>
    <div class="w-full h-full"><RouterView /></div>
  </div>
</template>
