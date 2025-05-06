import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { VueWrapper, mount } from "@vue/test-utils";
import CreateGoalCategoryForm from "./CreateGoalCategoryForm.vue";

function mountComponent() {
  return mount(CreateGoalCategoryForm);
}

type CreateGoalFormInstance = InstanceType<typeof CreateGoalCategoryForm>;

describe("CreateGoalCategoryForm", () => {
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

  it("should fill in the fields and create a new goal category", async () => {
    await wrapper.find("input[type=text]").setValue("Test Title");

    await wrapper.find("input[type=number]").setValue(50);
    await wrapper.find("form").trigger("submit.prevent");

    expect(wrapper.emitted("submit")?.[0][0]).toStrictEqual({
      title: "Test Title",
      xp_per_goal: 50,
    });
  });
});
