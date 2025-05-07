import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { type VueWrapper, mount } from "@vue/test-utils";
import CreateGoalForm from "./CreateGoalForm.vue";

function mountComponent() {
	return mount(CreateGoalForm, {
		props: {
			categoryId: "test-id",
		},
	});
}

type CreateGoalFormInstance = InstanceType<typeof CreateGoalForm>;

describe("CreateGoalForm", () => {
	let wrapper: VueWrapper<CreateGoalFormInstance>;

	beforeEach(() => {
		wrapper = mountComponent();
	});

	afterEach(() => {
		vi.resetAllMocks();
		wrapper.unmount();
	});

	it("should render the component", () => {
		expect(wrapper.exists()).toBe(true);
		expect(wrapper.isVisible()).toBe(true);
	});

	it("should fill in the fields and create a new goal", async () => {
		const titleInput = wrapper.find("input");
		const descriptionInput = wrapper.find("textarea");

		await titleInput.setValue("Test Title");
		await descriptionInput.setValue("Test Description");
		await wrapper.find("form").trigger("submit.prevent");

		expect(wrapper.emitted("submit")?.[0][0]).toStrictEqual({
			title: "Test Title",
			description: "Test Description",
		});
	});
});
