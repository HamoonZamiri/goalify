<script setup lang="ts">
import router from "@/router";
import { API_BASE } from "@/utils/constants";
import { Schemas, type ErrorResponse } from "@/utils/schemas";
import { ref } from "vue";
import useAuth from "@/hooks/auth/useAuth";
import Box from "@/components/primitives/Box.vue";
import Text from "@/components/primitives/Text.vue";
import InputField from "@/components/primitives/InputField.vue";
import Button from "@/components/primitives/Button.vue";

const emit = defineEmits(["submit"]);
const { setUser } = useAuth();
const error = ref<ErrorResponse>();
const formData = ref<{ email: string; password: string }>({
	email: "",
	password: "",
});

async function login() {
	emit("submit", { ...formData.value });
	const res = await fetch(`${API_BASE}/users/login`, {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(formData.value),
	});
	const json: unknown = await res.json();
	if (!res.ok) {
		error.value = json as ErrorResponse;
		return;
	}
	const parsed = Schemas.UserSchema.parse(json);
	setUser(parsed);
	error.value = undefined;
	router.push({ name: "Home" });
}
</script>

<template>
  <Box bg="darkest" width="w-full" height="h-full" class="items-center">
    <Text as="h3" size="3xl"> Sign in to your account </Text>
    <form @submit.prevent="login" class="w-4/5 sm:w-2/5 flex flex-col gap-4">
      <InputField
        bg="primary"
        text-color="dark"
        type="email"
        v-model="formData.email"
        errorslot
      >
        <template #label><Text>Email</Text></template>
        <template v-if="error?.errors?.email" #error>
          <Text color="error">{{ error?.errors?.email }}</Text>
        </template>
      </InputField>
      <InputField
        text-color="dark"
        bg="primary"
        type="password"
        v-model="formData.password"
        errorslot
      >
        <template #label><Text>Password</Text></template>
        <template v-if="error?.errors?.password" #error>
          <Text color="error">{{ error?.errors?.password }}</Text>
        </template>
      </InputField>
      <Button class="mt-4" height="h-10" width="w-full">
        <Text>Login</Text>
      </Button>
    </form>
    <Text v-if="error" as="p" size="sm" color="error">
      {{ error?.message }}
    </Text>
  </Box>
</template>
