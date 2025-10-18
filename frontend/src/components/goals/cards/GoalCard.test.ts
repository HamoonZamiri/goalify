import { mount, type VueWrapper } from "@vue/test-utils";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { goal } from "@/__mocks__/mocks";
import GoalCard from "./GoalCard.vue";

function mountComponent() {
	return mount(GoalCard, {
		props: {
			goal: goal,
		},
	});
}

type GoalCardInstance = InstanceType<typeof GoalCard>;

describe("GoalCard tests", () => {
	let wrapper: VueWrapper<GoalCardInstance>;

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
