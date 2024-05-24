<script setup lang="ts"></script>

<template>
  <div class="text-lg font-semibold">Goalify</div>
  <nav>
    <RouterLink to="/login">Login</RouterLink>
    <RouterLink to="/register">Register</RouterLink>
  </nav>
  <Main> <RouterView /> </Main>
</template>
