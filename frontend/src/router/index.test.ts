import {
	createMemoryHistory,
	createRouter,
	type RouteLocationNormalized,
} from "vue-router";
import { describe, expect, it } from "vitest";
import { user } from "@/__mocks__/mocks";
import { createAuthGuard, routes, RouteNames } from "./index";

function makeRoute(name: string): RouteLocationNormalized {
	return {
		name,
		path: "/",
		params: {},
		query: {},
		hash: "",
		fullPath: "/",
		matched: [],
		meta: {},
		redirectedFrom: undefined,
	};
}

describe("createAuthGuard", () => {
	describe("unauthenticated user", () => {
		const guard = createAuthGuard(() => undefined);

		it("navigating to / redirects to /login", () => {
			expect(guard(makeRoute(RouteNames.HOME))).toEqual({
				name: RouteNames.LOGIN,
			});
		});

		it("navigating to /login stays on /login", () => {
			expect(guard(makeRoute(RouteNames.LOGIN))).toBeUndefined();
		});

		it("navigating to /register stays on /register", () => {
			expect(guard(makeRoute(RouteNames.REGISTER))).toBeUndefined();
		});

		it("navigating to unknown path stays on NotFound", () => {
			expect(guard(makeRoute(RouteNames.NOT_FOUND))).toBeUndefined();
		});
	});

	describe("authenticated user", () => {
		const guard = createAuthGuard(() => user);

		it("navigating to /login redirects to /", () => {
			expect(guard(makeRoute(RouteNames.LOGIN))).toEqual({
				name: RouteNames.HOME,
			});
		});

		it("navigating to /register redirects to /", () => {
			expect(guard(makeRoute(RouteNames.REGISTER))).toEqual({
				name: RouteNames.HOME,
			});
		});

		it("navigating to / stays on /", () => {
			expect(guard(makeRoute(RouteNames.HOME))).toBeUndefined();
		});
	});
});

describe("createAppRouter integration", () => {
	it("unauthenticated user navigating to / is redirected to /login", async () => {
		const router = createRouter({ history: createMemoryHistory(), routes });
		router.beforeEach(createAuthGuard(() => undefined));
		await router.push("/");
		expect(router.currentRoute.value.name).toBe(RouteNames.LOGIN);
	});

	it("authenticated user navigating to /login is redirected to /", async () => {
		const router = createRouter({ history: createMemoryHistory(), routes });
		router.beforeEach(createAuthGuard(() => user));
		await router.push("/login");
		expect(router.currentRoute.value.name).toBe(RouteNames.HOME);
	});
});
