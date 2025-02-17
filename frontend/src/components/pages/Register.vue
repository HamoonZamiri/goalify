<script setup lang="ts">
import router from "@/router";
import { API_BASE } from "@/utils/constants";
import { Schemas, type ErrorResponse, type User } from "@/utils/schemas";
import { ref } from "vue";
import useAuth from "@/hooks/auth/useAuth";

const { setUser } = useAuth();
const error = ref<ErrorResponse | null>(null);
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
  const reqBody = {
    email: formData.value.email,
    password: formData.value.password,
    confirm_password: formData.value.confirmPassword,
  };

  const res = await fetch(`${API_BASE}/users/signup`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(reqBody),
  });
  const json: unknown = await res.json();
  if (!res.ok) {
    error.value = json as ErrorResponse;
    return;
  }
  const parsed = Schemas.UserResponseSchema.parse(json);
  setUser(parsed.data as User);
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
          class="block h-10 w-full bg-gray-300 rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm placeholder:text-gray-400 sm:text-sm sm:leading-6 focus:outline-none"
          type="email"
        />
        <p class="text-red-400" v-if="error?.errors?.email">
          {{ error.errors.email }}
        </p>
      </div>
      <div class="">
        <label class="text-gray-200">Password</label>
        <input
          v-model="formData.password"
          class="block h-10 w-full bg-gray-300 rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm placeholder:text-gray-400 sm:text-sm sm:leading-6 focus:outline-none"
          type="password"
        />
        <p class="text-red-400" v-if="error?.errors?.password">
          {{ error.errors.password }}
        </p>
      </div>
      <div class="mb-2">
        <label class="text-gray-200">Confirm Password</label>
        <input
          v-model="formData.confirmPassword"
          class="block h-10 w-full bg-gray-300 rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm placeholder:text-gray-400 sm:text-sm sm:leading-6 focus:outline-none"
          type="password"
        />
        <p class="text-red-400" v-if="error?.errors?.confirm_password">
          {{ error.errors.confirm_password }}
        </p>
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
