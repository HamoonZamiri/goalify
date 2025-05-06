import { createRouter, createWebHistory } from "vue-router";
import Login from "@/components/pages/Login.vue";
import Register from "@/components/pages/Register.vue";
import Home from "@/components/pages/Home.vue";

const routes = [
  { name: "Login", path: "/login", component: Login },
  { name: "Register", path: "/register", component: Register },
  { name: "Home", path: "/", component: Home },
];
const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach(async (to, _) => {
  const { default: useAuth } = await import("@/hooks/auth/useAuth");
  const { getUser } = useAuth();
  const user = getUser();

  if (!user && to.name !== "Login" && to.name !== "Register") {
    return { name: "Login" };
  } else if (user && (to.name === "Login" || to.name === "Register")) {
    return { name: "Home", path: "/" };
  }
});

export const RouteNames = {
  LOGIN: "Login",
  REGISTER: "Register",
  HOME: "Home",
} as const;

export default router;
