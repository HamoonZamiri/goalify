<script setup lang="ts">
import useAuth from "@/hooks/auth/useAuth";
import router from "@/router/index";
import Navbar from "@/components/navbar/Navbar.vue";
import Box from "@/components/primitives/Box.vue";

const { isLoggedIn } = useAuth();
if (!isLoggedIn()) {
	router.push({ name: "Login" });
}
</script>

<template>
  <Box class="bg-gray-900 h-screen w-screen">
    <header class="bg-gray-800 mb-2 h-auto">
      <Box flex-direction="row" class="justify-between p-6">
        <RouterLink to="/">
          <h1 class="font-semibold text-3xl hover:text-gray-300">Goalify</h1>
        </RouterLink>
        <Navbar />
      </Box>
    </header>
    <RouterView />
  </Box>
</template>
