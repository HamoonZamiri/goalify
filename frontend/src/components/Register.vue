<script setup lang="ts">
import router from "@/router";
import { API_BASE } from "@/utils/constants";
import {
  UserSchema,
  createServerResponseSchema,
  type User,
} from "@/utils/schemas";
import { setUser } from "@/utils/user";
import { ref } from "vue";

const error = ref<string | null>(null);
const formData = ref<{
  email: string;
  password: string;
  confirmPassword: string;
}>({
  email: "",
  password: "",
  confirmPassword: "",
});

async function signup(payload: MouseEvent) {
  payload.preventDefault();
  const res = await fetch(`${API_BASE}/users/signup`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(formData.value),
  });
  const json: unknown = await res.json();
  if (!res.ok) {
    error.value = (json as { message: string }).message;
    return;
  }
  const parsed = createServerResponseSchema(UserSchema).parse(json);
  setUser(parsed.data as User);
  error.value = null;
  router.push({ name: "Home" });
}
</script>

<template>
  <div class="w-full flex flex-col items-center">
    <h3 class="text-3xl font-semibold mb-6">Sign up for a new account</h3>
    <form class="w-4/5 sm:w-2/5">
      <div class="mb-4">
        <label class="">Email</label>
        <input
          v-model="formData.email"
          class="block w-full rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 sm:text-sm sm:leading-6"
          type="email"
        />
      </div>
      <div class="mb-6">
        <label>Password</label>
        <input
          v-model="formData.password"
          class="block w-full rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-50 sm:text-sm sm:leading-6"
          type="password"
        />
      </div>
      <div class="mb-6">
        <label>Confirm Password</label>
        <input
          v-model="formData.confirmPassword"
          class="block w-full rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-400 sm:text-sm sm:leading-6"
          type="password"
        />
      </div>
      <button @click="signup" class="w-full h-8 bg-blue-100 hover:bg-blue-200">
        Signup
      </button>
    </form>
  </div>
</template>
