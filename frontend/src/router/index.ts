import { createRouter, createWebHistory } from "vue-router";
import Login from "../components/Login.vue";
import Register from "../components/Register.vue";
import Home from "@/components/Home.vue";
import { isLoggedIn } from "@/utils/user";
const routes = [
  { name: "Login", path: "/login", component: Login },
  { name: "Register", path: "/register", component: Register },
  { name: "Home", path: "/", component: Home },
];
const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, _) => {
  if (!isLoggedIn() && to.name !== "Login" && to.name !== "Register") {
    return { name: "Login" };
  }
});

export default router;
