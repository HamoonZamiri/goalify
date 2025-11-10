<script setup lang="ts">
import router from "@/router";
import { Gear, Gift, Home, Signout } from "@/shared/components/icons";
import { Box, Text } from "@/shared/components/ui";
import { useAuth } from "@/shared/hooks";

const { logout } = useAuth();

function handleLogout(e: MouseEvent) {
	e.preventDefault();
	logout();
	router.push({ name: "Login" });
}
</script>
<template>
	<Box flex-direction="col" bg="darkest" gap="gap-2" padding="p-2">
		<RouterLink
			class="hover:bg-gray-500 rounded-md flex flex-row items-center gap-2"
			to="/"
		>
			<Home/>
			<Text size="lg">Dashboard</Text>
		</RouterLink>
		<RouterLink
			class="hover:bg-gray-500 rounded-md flex flex-row items-center gap-2"
			to="/"
		>
			<Gift/>
			<Text size="lg">Rewards (Coming Soon)</Text>
		</RouterLink>
		<RouterLink
			class="hover:bg-gray-500 rounded-md flex flex-row items-center gap-2"
			to="/"
		>
			<Gear/>
			<Text size="lg">Settings (Coming Soon)</Text>
		</RouterLink>
		<RouterLink
			class="hover:bg-gray-500 rounded-md flex flex-row items-center gap-2"
			to="/"
			:onclick="handleLogout"
		>
			<Signout :onclick="handleLogout"/>
			<Text size="lg">Sign out</Text>
		</RouterLink>
	</Box>
</template>
