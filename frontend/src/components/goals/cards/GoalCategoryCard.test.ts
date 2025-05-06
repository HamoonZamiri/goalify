import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { VueWrapper, mount } from "@vue/test-utils";
import GoalCategoryCard from "./GoalCategoryCard.vue";
import { goalCategory } from "@/__mocks__/mocks";

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
