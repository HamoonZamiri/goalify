<script setup lang="ts">
import router from "@/router";
import { Icon } from "@/shared/components/icons";
import { Box, Text, Button } from "@/shared/components/ui";
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
			class="hover:bg-gray-800 rounded-md flex flex-row items-center gap-2 p-0.5"
			to="/"
		>
			<Icon name="home"/>
			<Text size="lg">Dashboard</Text>
		</RouterLink>
		<RouterLink
			class="hover:bg-gray-800 rounded-md flex flex-row items-center gap-2 p-0.5"
			to="/"
		>
			<Icon name="gift"/>
			<Text size="lg">Rewards (Coming Soon)</Text>
		</RouterLink>
		<RouterLink
			class="hover:bg-gray-800 rounded-md flex flex-row items-center gap-2 p-0.5"
			to="/"
		>
			<Icon name="gear"/>
			<Text size="lg">Settings (Coming Soon)</Text>
		</RouterLink>
		<Button
			variant="ghost"
			height="h-auto"
			width="w-auto"
			class="justify-start gap-2 rounded-md hover:bg-gray-500 p-0.5"
			@click="handleLogout"
		>
			<Icon name="signout" size="sm"/>
			<Text size="lg">Sign out</Text>
		</Button>
	</Box>
</template>
