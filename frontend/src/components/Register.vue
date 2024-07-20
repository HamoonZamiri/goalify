<script setup lang="ts">
import router from "@/router";
import { API_BASE } from "@/utils/constants";
import {
  UserSchema,
  createServerResponseSchema,
  type User,
} from "@/utils/schemas";
import authState from "@/state/auth";
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
  authState.setUser(parsed.data as User);
  error.value = null;
  router.push({ name: "Home" });
}
</script>

<template>
  <div class="w-full flex flex-col items-center">
    <h3 class="text-3xl text-white mb-6">Sign up for a new account</h3>
    <form class="w-4/5 sm:w-2/5 grid grid-cols-1 gap-4">
      <div class="">
        <label class="text-gray-200">Email</label>
        <input
          v-model="formData.email"
          class="block h-10 w-full bg-gray-300 rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm placeholder:text-gray-400 sm:text-sm sm:leading-6"
          type="email"
        />
      </div>
      <div class="">
        <label class="text-gray-200">Password</label>
        <input
          v-model="formData.password"
          class="block h-10 w-full bg-gray-300 rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm placeholder:text-gray-400 sm:text-sm sm:leading-6"
          type="password"
        />
      </div>
      <div class="mb-2">
        <label class="text-gray-200">Confirm Password</label>
        <input
          v-model="formData.confirmPassword"
          class="block h-10 w-full bg-gray-300 rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm placeholder:text-gray-400 sm:text-sm sm:leading-6"
          type="password"
        />
      </div>
      <button
        @click="signup"
        class="bg-blue-400 mt-4 rounded-lg text-gray-300 hover:bg-blue-500 h-10 py-1.5"
      >
        Signup
      </button>
    </form>
  </div>
</template>
