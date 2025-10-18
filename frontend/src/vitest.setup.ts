import { HttpResponse, http } from "msw";
import { setupServer } from "msw/node";
import { afterAll, afterEach, beforeAll, vi } from "vitest";
import { goal, goalCategory, levelOne } from "@/__mocks__/mocks";

/**
 * Global mocks for all tests
 */

// Mock useAuth to provide a test user with access token
vi.mock("@/shared/hooks/auth/useAuth", async () => {
	const { ref } = await import("vue");
	const mockUser = {
		id: "test-user-id",
		email: "test@example.com",
		access_token: "test-token",
		refresh_token: "test-refresh-token",
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

const API_BASE = "http://localhost:8080/api" as const;

export const restHandlers = [
	http.post(`${API_BASE}/goals`, () => {
		return HttpResponse.json(goal);
	}),
	http.post(`${API_BASE}/goals/categories`, () => {
		return HttpResponse.json(goalCategory);
	}),
	http.get(`${API_BASE}/levels/1`, () => {
		return HttpResponse.json(levelOne);
	}),
] as const;

export const server = setupServer(...restHandlers);

// Start server before all tests
beforeAll(() => server.listen({ onUnhandledRequest: "error" }));

// Close server after all tests
afterAll(() => server.close());

// Reset handlers after each test for test isolation
afterEach(() => server.resetHandlers());
