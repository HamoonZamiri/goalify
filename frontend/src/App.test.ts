import type { VueWrapper } from "@vue/test-utils";
import { afterEach, describe, expect, it } from "vitest";
import { mountWithPlugins } from "@/shared/test-utils";
import useAuth from "@/shared/hooks/auth/useAuth";
import App from "./App.vue";

describe("App", () => {
	let wrapper: VueWrapper;

	afterEach(() => {
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
