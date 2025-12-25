import type { VueWrapper } from "@vue/test-utils";
import { flushPromises } from "@vue/test-utils";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { mountWithPlugins, setupFetchSpies } from "@/shared/test-utils";
import CreateGoalCategoryForm from "./CreateGoalCategoryForm.vue";

function mountComponent() {
	return mountWithPlugins(CreateGoalCategoryForm);
}

describe("CreateGoalCategoryForm", () => {
	let wrapper: VueWrapper;
	let _fetchMock: ReturnType<typeof setupFetchSpies>;

	beforeEach(() => {
		_fetchMock = setupFetchSpies([
			{
				url: "/goals/categories",
				method: "POST",
				response: {
					id: "123e4567-e89b-12d3-a456-426614174000",
					title: "Test Title",
					xp_per_goal: 50,
					user_id: "123e4567-e89b-12d3-a456-426614174001",
					goals: [],
					created_at: new Date().toISOString(),
					updated_at: new Date().toISOString(),
				},
			},
		]);
		wrapper = mountComponent();
	});

	afterEach(() => {
		vi.restoreAllMocks();
		wrapper.unmount();
	});

	it("renders the component", () => {
		expect(wrapper.exists()).toBe(true);
		expect(wrapper.isVisible()).toBe(true);
	});

	it("fills in the fields and creates a new goal category", async () => {
		const titleInput = wrapper.find("input[name=title]");
		const xpInput = wrapper.find("input[name=xp_per_goal]");

		await titleInput.setValue("Test Title");
		await xpInput.setValue("50");
		await wrapper.find("form").trigger("submit");
		await flushPromises();

		expect(wrapper.emitted("submit")?.[0]?.[0]).toMatchObject({
			title: "Test Title",
			xp_per_goal: 50,
		});
	});
});
