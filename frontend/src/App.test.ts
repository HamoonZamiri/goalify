import type { VueWrapper } from "@vue/test-utils";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { levelOne } from "@/__mocks__/mocks";
import { mountWithPlugins, setupFetchSpies } from "@/shared/test-utils";
import useAuth from "@/shared/hooks/auth/useAuth";
import { API_BASE } from "@/utils/constants";
import App from "./App.vue";

describe("App", () => {
	let wrapper: VueWrapper;

	beforeEach(() => {
		// Mock GET /api/levels/1 called by Navbar component
		setupFetchSpies([
			{
				url: `${API_BASE}/levels/1`,
				method: "GET",
				response: levelOne,
			},
		]);
	});

	afterEach(() => {
		vi.restoreAllMocks();
		wrapper?.unmount();
	});

	describe("Sidebar visibility", () => {
		it("shows Sidebar when user is logged in", () => {
			// Default mock state is logged in
			wrapper = mountWithPlugins(App, {
				global: {
					stubs: {
						RouterLink: true,
						RouterView: true,
					},
				},
			});

			const sidebar = wrapper.findComponent({ name: "SideBar" });
			expect(sidebar.exists()).toBe(true);
		});

		it("hides Sidebar when user is logged out", () => {
			// Simulate logged-out state by clearing authState
			const { authState } = useAuth();
			authState.value = undefined;

			wrapper = mountWithPlugins(App, {
				global: {
					stubs: {
						RouterLink: true,
						RouterView: true,
					},
				},
			});

			const sidebar = wrapper.findComponent({ name: "SideBar" });
			expect(sidebar.exists()).toBe(false);
		});
	});
});
