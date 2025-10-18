import type { VueWrapper } from "@vue/test-utils";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { mountWithPlugins } from "@/shared/test-utils";
import CreateGoalButton from "./CreateGoalButton.vue";

describe("CreateGoalButton", () => {
	let wrapper: VueWrapper;
	beforeEach(() => {
		wrapper = mountWithPlugins(CreateGoalButton);
	});

	afterEach(() => {
		vi.resetAllMocks();
		wrapper.unmount();
	});

	it("should render the component", () => {
		expect(wrapper.exists()).toBe(true);
		expect(wrapper.isVisible()).toBe(true);
	});

	it("should find the svg element", async () => {
		const svg = wrapper.findAll("svg");
		const path = wrapper.findAll("path");
		expect(svg.length).toBe(1);
		expect(path.length).toBe(1);
	});
});
