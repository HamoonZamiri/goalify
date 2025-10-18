import type { VueWrapper } from "@vue/test-utils";
import { flushPromises } from "@vue/test-utils";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { mountWithPlugins, setupFetchSpies } from "@/shared/test-utils";
import CreateGoalForm from "./CreateGoalForm.vue";

function mountComponent() {
	return mountWithPlugins(CreateGoalForm, {
		props: {
			categoryId: "123e4567-e89b-12d3-a456-426614174003",
		},
	});
}

describe("CreateGoalForm", () => {
	let wrapper: VueWrapper;

	beforeEach(() => {
		setupFetchSpies([
			{
				url: "/goals",
				method: "POST",
				response: {
					id: "123e4567-e89b-12d3-a456-426614174002",
					title: "Test Title",
					description: "Test Description",
					category_id: "123e4567-e89b-12d3-a456-426614174003",
					user_id: "123e4567-e89b-12d3-a456-426614174001",
					status: "not_complete",
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

	it("fills in the fields and creates a new goal", async () => {
		const titleInput = wrapper.find("input");
		const descriptionInput = wrapper.find("textarea");

		await titleInput.setValue("Test Title");
		await descriptionInput.setValue("Test Description");
		await wrapper.find("form").trigger("submit");
		await flushPromises();

		expect(wrapper.emitted("submit")?.[0][0]).toMatchObject({
			title: "Test Title",
			description: "Test Description",
			category_id: "123e4567-e89b-12d3-a456-426614174003",
		});
	});
});
