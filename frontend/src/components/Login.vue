<script setup lang="ts"></script>

<template>
  <div className="w-full flex flex-col items-center">
    <h3 className="text-3xl font-semibold mb-6 text-center">
      Sign in to your account
    </h3>
    <form className="w-4/5 sm:w-2/5">
      <div className="mb-4">
        <label className="">Email</label>
        <input
          className="block w-full rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300
                        placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-200 sm:text-sm sm:leading-6"
          type="email"
        />
      </div>
      <div className="mb-6">
        <label>Password</label>
        <input
          className="block w-full rounded-md border-0 px-1.5 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300
                        placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
          type="password"
        />
      </div>
      <button className="w-full h-8 bg-blue-100 hover:bg-blue-200">
        Login
      </button>
    </form>
  </div>
</template>
