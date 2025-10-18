import { mount, type VueWrapper } from "@vue/test-utils";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { goalCategory } from "@/__mocks__/mocks";
import GoalCategoryCard from "./GoalCategoryCard.vue";

function mountComponent() {
	return mount(GoalCategoryCard, {
		props: {
			goalCategory: goalCategory,
		},
	});
}

type GoalCategoryCardInstance = InstanceType<typeof GoalCategoryCard>;

describe("GoalCategoryCard tests", () => {
	let wrapper: VueWrapper<GoalCategoryCardInstance>;

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
});
