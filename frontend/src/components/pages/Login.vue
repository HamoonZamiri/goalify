<script setup lang="ts">
import router from "@/router";
import { API_BASE } from "@/utils/constants";
import { Schemas, type ErrorResponse, type User } from "@/utils/schemas";
import { ref } from "vue";
import useAuth from "@/hooks/auth/useAuth";

const { setUser } = useAuth();
const error = ref<ErrorResponse>();
const formData = ref<{ email: string; password: string }>({
  email: "",
  password: "",
});

async function login(payload: MouseEvent) {
  payload.preventDefault();
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
  <div class="w-full flex flex-col items-center">
    <h3 class="text-3xl text-gray-300 mb-6 text-center">
      Sign in to your account
    </h3>
    <form class="w-4/5 sm:w-2/5 grid grid-cols-1 gap-4">
      <div class="">
        <label class="text-gray-200">Email</label>
        <input
          v-model="formData.email"
          class="block h-10 w-full bg-gray-300 rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm placeholder:text-gray-400 sm:text-sm sm:leading-6 focus:outline-none"
          type="email"
        />
      </div>
      <div class="mb-4">
        <label class="text-gray-200">Password</label>
        <input
          v-model="formData.password"
          class="block h-10 w-full bg-gray-300 rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm placeholder:text-gray-400 sm:text-sm sm:leading-6 focus:outline-none"
          type="password"
        />
      </div>
      <button
        @click="login"
        class="bg-blue-400 mt-4 rounded-lg text-gray-300 hover:bg-blue-500 h-10 py-1.5"
      >
        Login
      </button>
    </form>
    <p class="text-red-400 font-semibold">{{ error?.message }}</p>
  </div>
</template>
