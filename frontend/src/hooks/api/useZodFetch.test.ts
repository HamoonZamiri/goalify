import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { http, HttpResponse } from "msw";
import { server } from "@/vitest.setup";
import useZodFetch from "./useZodFetch";
import useAuth from "../auth/useAuth";
import { Schemas } from "@/utils/schemas";
import { user as mockUser } from "@/__mocks__/mocks";
import Cookies from "js-cookie";

const API_BASE = "http://localhost:8080/api";

vi.stubGlobal(
	"EventSource",
	vi.fn(() => ({
		addEventListener: vi.fn(),
		removeEventListener: vi.fn(),
		close: vi.fn(),
		readyState: 0,
		url: "",
		onopen: null,
		onerror: null,
		onmessage: null,
	})),
);

describe("useZodFetch - Token Refresh Flow", () => {
	let requestCount = 0;
	const initialUser = {
		...mockUser,
		access_token: "old-access-token",
		refresh_token: "026a715f-a023-4e6b-973e-4bb0e96562af",
	};

	const refreshedUser = {
		...mockUser,
		access_token: "new-access-token",
		refresh_token: "126a715f-a023-4e6b-973e-4bb0e96562af",
	};

	beforeEach(() => {
		requestCount = 0;
		Cookies.remove("user");

		server.use(
			http.get(`${API_BASE}/goals/categories`, ({ request }) => {
				requestCount++;
				const authHeader = request.headers.get("Authorization");

				if (requestCount === 1) {
					return HttpResponse.json(
						{ message: "Unauthorized" },
						{ status: 401 },
					);
				}

				if (authHeader === `Bearer ${refreshedUser.access_token}`) {
					return HttpResponse.json({
						data: [
							{
								id: "126a715f-a023-4e6b-973e-4bb0e96562af",
								title: "Test Category",
								xp_per_goal: 50,
								user_id: mockUser.id,
								goals: [],
							},
						],
					});
				}

				return HttpResponse.json({ message: "Invalid token" }, { status: 401 });
			}),

			http.post(`${API_BASE}/users/refresh`, async ({ request }) => {
				const contentType = request.headers.get("Content-Type");

				if (contentType !== "application/json") {
					return HttpResponse.json(
						{ message: "Content-Type must be application/json" },
						{ status: 400 },
					);
				}

				const body = (await request.json()) as {
					user_id: string;
					refresh_token: string;
				};

				if (
					body.user_id === initialUser.id &&
					body.refresh_token === initialUser.refresh_token
				) {
					return HttpResponse.json(refreshedUser);
				}

				return HttpResponse.json(
					{ message: "Invalid refresh token" },
					{ status: 401 },
				);
			}),
		);
	});

	afterEach(() => {
		vi.clearAllMocks();
		Cookies.remove("user");
	});

	it("refreshes token on 401 and retries request with new token", async () => {
		const { setUser } = useAuth();
		const { zodFetch } = useZodFetch();

		setUser(initialUser);

		const result = await zodFetch(
			`${API_BASE}/goals/categories`,
			Schemas.GoalCategoryResponseArraySchema,
			{
				headers: {
					Authorization: `Bearer ${initialUser.access_token}`,
				},
			},
		);

		expect(requestCount).toBe(2);

		expect(result).toHaveProperty("data");
		if ("data" in result) {
			expect(result.data).toHaveLength(1);
			expect(result.data[0]).toMatchObject({
				title: "Test Category",
				xp_per_goal: 50,
			});
		}

		const updatedUser = useAuth().getUser();
		expect(updatedUser?.access_token).toBe(refreshedUser.access_token);
		expect(updatedUser?.refresh_token).toBe(refreshedUser.refresh_token);
	});

	it("includes Content-Type header in refresh request", async () => {
		const { setUser } = useAuth();
		const { zodFetch } = useZodFetch();

		setUser(initialUser);

		const result = await zodFetch(
			`${API_BASE}/goals/categories`,
			Schemas.GoalCategoryResponseArraySchema,
			{
				headers: {
					Authorization: `Bearer ${initialUser.access_token}`,
				},
			},
		);

		expect(result).not.toHaveProperty("statusCode");
		if ("data" in result) {
			expect(result.data).toBeDefined();
		}
	});

	it("logs out and returns error when refresh fails", async () => {
		const { setUser, getUser } = useAuth();
		const { zodFetch } = useZodFetch();

		const invalidUser = {
			...initialUser,
			refresh_token: "226a715f-a023-4e6b-973e-4bb0e96562af",
		};
		setUser(invalidUser);

		const result = await zodFetch(
			`${API_BASE}/goals/categories`,
			Schemas.GoalCategoryResponseArraySchema,
			{
				headers: {
					Authorization: `Bearer ${invalidUser.access_token}`,
				},
			},
		);

		expect(result).toHaveProperty("message");
		if ("message" in result) {
			expect(result.message).toBe("Unauthorized");
		}

		expect(getUser()).toBeUndefined();
	});

	it("returns error when user is not logged in during 401", async () => {
		const { zodFetch } = useZodFetch();

		const result = await zodFetch(
			`${API_BASE}/goals/categories`,
			Schemas.GoalCategoryResponseArraySchema,
			{
				headers: {
					Authorization: "Bearer invalid-token",
				},
			},
		);

		expect(result).toHaveProperty("message");
		if ("message" in result) {
			expect(result.message).toBe("Unauthorized");
		}
	});
});
