<script setup lang="ts">
import { useRoute } from "vue-router";
import { Navbar, Sidebar } from "@/shared/components/navigation";
import { Box } from "@/shared/components/ui";
import useAuth from "@/shared/hooks/auth/useAuth";
import { RouteNames } from "@/router/index";

const { isLoggedIn } = useAuth();
const route = useRoute();
</script>

<template>
	<RouterView
		v-if="route.name === RouteNames.NOT_FOUND"
		class="h-screen w-screen"
	/>
	<Box v-else bg="darkest" class="h-screen w-screen">
		<header class="bg-gray-800">
			<Box flex-direction="row" class="justify-between p-6">
				<RouterLink to="/">
					<h1 class="font-semibold text-3xl hover:text-gray-300">Goalify</h1>
				</RouterLink>
				<Navbar />
			</Box>
		</header>
		<Box flex-direction="row" class="w-full h-full">
			<Sidebar v-if="isLoggedIn()" class="rounded-none text-nowrap" />
			<RouterView class="rounded-none w-full h-full" />
		</Box>
	</Box>
</template>
