import {
	createMemoryHistory,
	createRouter,
	createWebHistory,
	type RouteLocationNormalized,
	type RouteLocationRaw,
	type RouterHistory,
} from "vue-router";
import { LoginPage, RegisterPage } from "@/features/auth";
import type { User } from "@/features/auth/schemas";
import HomePage from "@/pages/HomePage.vue";
import NotFoundPage from "@/pages/NotFoundPage.vue";
import useAuth from "@/shared/hooks/auth/useAuth";

export const routes = [
	{ name: "Login", path: "/login", component: LoginPage },
	{ name: "Register", path: "/register", component: RegisterPage },
	{ name: "Home", path: "/", component: HomePage },
	{ name: "NotFound", path: "/:pathMatch(.*)*", component: NotFoundPage },
];

export const RouteNames = {
	LOGIN: "Login",
	REGISTER: "Register",
	HOME: "Home",
	NOT_FOUND: "NotFound",
} as const;

const publicRoutes = new Set<string>([
	RouteNames.LOGIN,
	RouteNames.REGISTER,
	RouteNames.NOT_FOUND,
]);

export function createAuthGuard(getUser: () => User | undefined) {
	return (to: RouteLocationNormalized): RouteLocationRaw | undefined => {
		const user = getUser();

		if (!user && !publicRoutes.has(String(to.name))) {
			return { name: RouteNames.LOGIN };
		}
		if (
			user &&
			(to.name === RouteNames.LOGIN || to.name === RouteNames.REGISTER)
		) {
			return { name: RouteNames.HOME };
		}
	};
}

export function createAppRouter(history: RouterHistory = createWebHistory()) {
	const router = createRouter({ history, routes });
	const { getUser } = useAuth();
	router.beforeEach(createAuthGuard(getUser));
	return router;
}

export { createMemoryHistory };
export default createAppRouter();
