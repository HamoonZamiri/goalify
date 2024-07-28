import { createRouter, createWebHistory } from "vue-router";
import Login from "../components/Login.vue";
import Register from "../components/Register.vue";
import Home from "@/components/Home.vue";
import authState from "@/state/auth";
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
  if (!authState.user && to.name !== "Login" && to.name !== "Register") {
    return { name: "Login" };
  } else if (
    authState.user &&
    (to.name === "Login" || to.name === "Register")
  ) {
    return { name: "Home", path: "/" };
  }
});

export const RouteNames = {
  LOGIN: "Login",
  REGISTER: "Register",
  HOME: "Home",
} as const;

export default router;
