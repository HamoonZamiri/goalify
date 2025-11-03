import { vi } from "vitest";

type FetchMockConfig = {
	url: string | RegExp;
	method?: string;
	response?: unknown;
	responseFn?: (url: string, init?: RequestInit) => unknown;
	status?: number;
};

/**
 * Extracts URL string from fetch input parameter
 */
function getUrlFromInput(input: RequestInfo | URL): string {
	if (typeof input === "string") {
		return input;
	}
	if (input instanceof URL) {
		return input.href;
	}
	// Request object
	return input.url;
}

/**
 * Sets up fetch spies for testing API calls without MSW
 */
export function setupFetchSpies(configs: FetchMockConfig[]) {
	const fetchSpy = vi.spyOn(global, "fetch");

	fetchSpy.mockImplementation(async (input, init) => {
		const url = getUrlFromInput(input);
		const method = init?.method?.toUpperCase() || "GET";

		const config = configs.find((c) => {
			const urlMatches =
				typeof c.url === "string" ? url.includes(c.url) : c.url.test(url);
			const methodMatches = !c.method || c.method.toUpperCase() === method;
			return urlMatches && methodMatches;
		});

		if (!config) {
			throw new Error(`No mock found for ${method} ${url}`);
		}

		const data = config.responseFn
			? config.responseFn(url, init)
			: config.response;

		return new Response(JSON.stringify(data), {
			status: config.status ?? 200,
			headers: { "Content-Type": "application/json" },
		});
	});

	return {
		spy: fetchSpy,
		getRequestBody: (url: string) => {
			const call = fetchSpy.mock.calls.find(([callUrl]) => {
				const callUrlStr = getUrlFromInput(callUrl);
				return callUrlStr.includes(url);
			});
			return call?.[1]?.body ? JSON.parse(call[1].body as string) : null;
		},
	};
}
