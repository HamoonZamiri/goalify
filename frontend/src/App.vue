<script setup lang="ts">
import useAuth from "@/shared/hooks/auth/useAuth";
import router from "@/router/index";
import { Box } from "@/shared/components/ui";
import { Navbar, Sidebar } from "@/shared/components/navigation";

const { isLoggedIn } = useAuth();
if (!isLoggedIn()) {
	router.push({ name: "Login" });
}
</script>

<template>
	<Box bg="darkest" class="h-screen w-screen">
		<header class="bg-gray-800">
			<Box flex-direction="row" class="justify-between p-6">
				<RouterLink to="/">
					<h1 class="font-semibold text-3xl hover:text-gray-300">Goalify</h1>
				</RouterLink>
				<Navbar/>
			</Box>
		</header>
		<Box flex-direction="row" class="w-full">
			<Sidebar class="rounded-none text-nowrap"/>
			<RouterView class="rounded-none"/>
		</Box>
	</Box>
</template>
