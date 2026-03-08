<script setup lang="ts">
import { RouterLink } from "vue-router";
import { Box, Text, Button } from "@/shared/components/ui";
</script>

<template>
	<Box height="h-full" bg="darkest" class="items-center justify-center gap-4">
		<Text as="h1" size="3xl" weight="bold">404</Text>
		<Text as="p" size="xl">Page not found.</Text>
		<RouterLink to="/">
			<Button variant="primary">
				<Text>Go home</Text>
			</Button>
		</RouterLink>
	</Box>
</template>
