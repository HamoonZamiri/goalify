import { afterAll, afterEach, beforeAll } from "vitest";
import { setupServer } from "msw/node";
import { http, HttpResponse } from "msw";
import { randomUUID } from "crypto";
import { goal, goalCategory, levelOne } from "@/__mocks__/mocks";

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

const server = setupServer(...restHandlers);

// Start server before all tests
beforeAll(() => server.listen({ onUnhandledRequest: "error" }));

// Close server after all tests
afterAll(() => server.close());

// Reset handlers after each test for test isolation
afterEach(() => server.resetHandlers());
