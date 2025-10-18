import { mount, type VueWrapper } from "@vue/test-utils";
import {
	afterAll,
	afterEach,
	beforeEach,
	describe,
	expect,
	it,
	vi,
} from "vitest";
import { user } from "@/__mocks__/mocks";
import useAuth from "@/hooks/auth/useAuth";
import Navbar from "./Navbar.vue";

describe("Navbar tests", () => {
	let wrapper: VueWrapper;
	const { authState } = useAuth();
	authState.value = user;

	beforeEach(() => {
		wrapper = mount(Navbar, {
			global: {
				stubs: {
					RouterLink: true,
				},
			},
		});
	});

	afterEach(() => {
		vi.resetAllMocks();
		wrapper.unmount();
	});

	afterAll(() => {
		authState.value = undefined;
	});

	it("should render the Navbar", () => {
		expect(wrapper.exists()).toBe(true);
	});
});
