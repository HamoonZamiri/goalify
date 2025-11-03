import { createRouter, createWebHistory } from "vue-router";
import { LoginPage, RegisterPage } from "@/features/auth";
import HomePage from "@/pages/HomePage.vue";

const routes = [
	{ name: "Login", path: "/login", component: LoginPage },
	{ name: "Register", path: "/register", component: RegisterPage },
	{ name: "Home", path: "/", component: HomePage },
];
const router = createRouter({
	history: createWebHistory(),
	routes,
});

router.beforeEach(async (to, _) => {
	const { default: useAuth } = await import("@/shared/hooks/auth/useAuth");
	const { getUser } = useAuth();
	const user = getUser();

	if (!user && to.name !== "Login" && to.name !== "Register") {
		return { name: "Login" };
	}
	if (user && (to.name === "Login" || to.name === "Register")) {
		return { name: "Home", path: "/" };
	}
});

export const RouteNames = {
	LOGIN: "Login",
	REGISTER: "Register",
	HOME: "Home",
} as const;

export default router;
