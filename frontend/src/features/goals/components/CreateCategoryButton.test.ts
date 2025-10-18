import type { VueWrapper } from "@vue/test-utils";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { mountWithPlugins } from "@/shared/test-utils";
import CreateCategoryButton from "./CreateCategoryButton.vue";

const mountComponent = () => {
	return mountWithPlugins(CreateCategoryButton, {
		props: {
			setIsOpen: () => {},
			isOpen: true,
		},
	});
};

describe("CreateCategoryButton", () => {
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

	it("should find the svg element", async () => {
		const svg = wrapper.findAll("svg");
		const path = wrapper.findAll("path");
		expect(svg.length).toBe(1);
		expect(path.length).toBe(1);
	});
});
