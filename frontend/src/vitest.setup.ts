import { vi } from "vitest";

/**
 * Global mocks for all tests
 */

// Mock ResizeObserver for HeadlessUI Dialog components
globalThis.ResizeObserver = class ResizeObserver {
	observe() {}
	unobserve() {}
	disconnect() {}
};

// Mock useAuth to provide a test user with access token
vi.mock("@/shared/hooks/auth/useAuth", async () => {
	const { ref } = await import("vue");
	const mockUser = {
		id: "test-user-id",
		email: "test@example.com",
		access_token: "test-token",
		refresh_token: "test-refresh-token",
		level_id: 1,
	};
	const authState = ref<typeof mockUser | undefined>(mockUser);

	return {
		default: () => ({
			authState,
			getUser: () => authState.value,
			setUser: vi.fn((user) => {
				authState.value = user;
			}),
			logout: vi.fn(() => {
				authState.value = undefined;
			}),
			isLoggedIn: () => !!authState.value,
		}),
	};
});

// Mock toast notifications
vi.mock("vue3-toastify", () => ({
	toast: {
		success: vi.fn(),
		error: vi.fn(),
		warning: vi.fn(),
		info: vi.fn(),
	},
}));
