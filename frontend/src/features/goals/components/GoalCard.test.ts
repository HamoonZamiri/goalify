import type { VueWrapper } from "@vue/test-utils";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { goal } from "@/__mocks__/mocks";
import { mountWithPlugins } from "@/shared/test-utils";
import GoalCard from "./GoalCard.vue";

function mountComponent() {
	return mountWithPlugins(GoalCard, {
		props: {
			goal: goal,
		},
	});
}

describe("GoalCard tests", () => {
	let wrapper: VueWrapper;

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
