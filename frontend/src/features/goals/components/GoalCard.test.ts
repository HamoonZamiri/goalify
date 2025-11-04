import type { VueWrapper } from "@vue/test-utils";
import { flushPromises } from "@vue/test-utils";
import { afterEach, describe, expect, it, vi } from "vitest";
import { goal } from "@/__mocks__/mocks";
import { mountWithPlugins, setupFetchSpies } from "@/shared/test-utils";
import { API_BASE } from "@/utils/constants";
import GoalCard from "./GoalCard.vue";
import type { Goal } from "../schemas";

function mountComponent() {
	return mountWithPlugins(GoalCard, {
		props: {
			goal: goal,
		},
	});
}

describe("GoalCard tests", () => {
	let wrapper: VueWrapper;

	afterEach(() => {
		vi.restoreAllMocks();
		wrapper?.unmount();
	});

	it("renders the component", () => {
		setupFetchSpies([
			{
				url: `${API_BASE}/goals/${goal.id}`,
				method: "PUT",
				response: goal,
			},
		]);
		wrapper = mountComponent();

		expect(wrapper.exists()).toBe(true);
		expect(wrapper.isVisible()).toBe(true);
	});

	it("toggles goal status from not_complete to complete when check icon is clicked", async () => {
		const expectedUrl = `${API_BASE}/goals/${goal.id}`;
		const fetchMock = setupFetchSpies([
			{
				url: expectedUrl,
				method: "PUT",
				response: {
					...goal,
					status: "complete",
					updated_at: new Date().toISOString(),
				},
			},
		]);
		wrapper = mountComponent();

		const checkIcon = wrapper.find("svg");

		await checkIcon.trigger("click");
		await flushPromises();

		expect(fetchMock.spy).toHaveBeenCalledWith(
			expectedUrl,
			expect.objectContaining({
				method: "PUT",
				body: JSON.stringify({ status: "complete" }),
			}),
		);
	});

	it("toggles goal status from complete to not_complete when check icon is clicked", async () => {
		const completedGoal: Goal = { ...goal, status: "complete" };
		const expectedUrl = `${API_BASE}/goals/${completedGoal.id}`;

		const fetchMock = setupFetchSpies([
			{
				url: expectedUrl,
				method: "PUT",
				response: {
					...completedGoal,
					status: "not_complete",
					updated_at: new Date().toISOString(),
				},
			},
		]);

		wrapper = mountWithPlugins(GoalCard, {
			props: { goal: completedGoal },
		});

		const checkIcon = wrapper.find("svg");

		await checkIcon.trigger("click");
		await flushPromises();

		// Verify the fetch was called with correct parameters
		expect(fetchMock.spy).toHaveBeenCalledWith(
			expectedUrl,
			expect.objectContaining({
				method: "PUT",
				body: JSON.stringify({ status: "not_complete" }),
			}),
		);

		// Verify the fetch call succeeded (did not throw an error)
		const callResult = fetchMock.spy.mock.results[0];
		expect(callResult?.type).toBe("return");

		// Verify the response was successfully parsed
		const response = await callResult?.value;
		expect(response).toBeInstanceOf(Response);
		expect(response.ok).toBe(true);
	});

	it("opens EditGoalForm dialog when goal card is clicked", async () => {
		setupFetchSpies([
			{
				url: `${API_BASE}/goals/${goal.id}`,
				method: "PUT",
				response: goal,
			},
		]);
		wrapper = mountComponent();

		const cardHeader = wrapper.find("header");
		expect(wrapper.findComponent({ name: "EditGoalForm" }).exists()).toBe(
			false,
		);

		await cardHeader.trigger("click");
		await flushPromises();

		expect(wrapper.findComponent({ name: "EditGoalForm" }).exists()).toBe(true);
	});
});
