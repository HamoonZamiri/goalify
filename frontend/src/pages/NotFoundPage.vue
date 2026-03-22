<script setup lang="ts">
import { RouterLink } from "vue-router";
import { Box, Button, Text } from "@/shared/components/ui";
</script>

<template>
	<Box height="h-full" bg="darkest" class="items-center justify-center gap-4">
		<Text as="h1" size="3xl" weight="bold">404</Text>
		<Text as="p" size="xl"
			>We couldn't find the page you were looking for.</Text
		>
		<RouterLink to="/">
			<Button variant="primary" padding="p-4"> <Text>Go home</Text> </Button>
		</RouterLink>
	</Box>
</template>
